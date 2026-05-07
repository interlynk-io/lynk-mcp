// Copyright 2025 Interlynk.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/interlynk-io/lynk-mcp/internal/api"
	"github.com/interlynk-io/lynk-mcp/internal/config"
	"github.com/interlynk-io/lynk-mcp/internal/mcp"
	"github.com/interlynk-io/lynk-mcp/internal/retry"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"sigs.k8s.io/release-utils/version"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "lynk-mcp",
		Short: "MCP server for Lynk version management API",
		Long:  `lynk-mcp is an MCP (Model Context Protocol) server that bridges AI assistants with the Lynk API for version management, vulnerability tracking, and compliance checking.`,
	}

	// Configure command
	configureCmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure API token and settings",
		Long:  `Interactive setup to configure your Lynk API token. The token will be stored securely in your system keychain.`,
		RunE:  runConfigure,
	}

	// Verify command
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify connection and show organization info",
		Long:  `Test the API connection using the stored token and display your organization information.`,
		RunE:  runVerify,
	}

	// Serve command
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server",
		Long:  `Start the MCP server in stdio mode for integration with AI assistants like Claude Desktop.`,
		RunE:  runServe,
	}

	// Add commands
	rootCmd.AddCommand(configureCmd, verifyCmd, serveCmd)
	rootCmd.AddCommand(version.Version())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runConfigure(cmd *cobra.Command, args []string) error {
	fmt.Println("Lynk MCP Configuration")
	fmt.Println("======================")
	fmt.Println()

	// Load existing config or create new
	cfg, err := config.Load()
	if err != nil {
		// Config doesn't exist, create default
		cfg = config.DefaultConfig()
	}

	// Prompt for API endpoint
	fmt.Printf("API Endpoint [%s]: ", cfg.API.Endpoint)
	var endpoint string
	fmt.Scanln(&endpoint)
	if endpoint != "" {
		cfg.API.Endpoint = endpoint
	}

	// Prompt for API token
	fmt.Print("API Token (lynk_*): ")
	var token string
	fmt.Scanln(&token)

	if token == "" {
		return fmt.Errorf("API token is required")
	}

	// Validate token format
	if !config.ValidateTokenFormat(token) {
		return fmt.Errorf("invalid token format: must start with one of: %s", config.ValidTokenPrefixesDescription())
	}

	// Store token in keyring
	if err := config.StoreToken(token); err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println()
	fmt.Println("Configuration saved successfully!")
	fmt.Println("Token stored securely in system keychain.")
	fmt.Println()
	fmt.Println("Run 'lynk-mcp verify' to test your connection.")

	return nil
}

func runVerify(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'lynk-mcp configure' first)", err)
	}

	token, err := config.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %w (run 'lynk-mcp configure' or set %s environment variable)", err, config.EnvTokenKey)
	}

	fmt.Printf("Connecting to %s...\n", cfg.API.Endpoint)

	client := api.NewClient(cfg.API.Endpoint, token, cfg.API.Timeout)

	var org *api.Organization
	start := time.Now()

	err = retry.Do(cmd.Context(), retry.DefaultVerifyConfig(), func() error {
		var getErr error
		org, getErr = client.GetOrganization(cmd.Context())
		return getErr
	}, nil, func(attempt int, retryErr error, delay time.Duration) {
		elapsed := time.Since(start).Truncate(time.Second)
		fmt.Printf("  Attempt %d failed (%s elapsed): %v\n", attempt, elapsed, retryErr)
		fmt.Printf("  Token may still be propagating. Retrying in %s...\n", delay.Truncate(time.Second))
	})
	if err != nil {
		return fmt.Errorf("failed to connect after %s: %w", time.Since(start).Truncate(time.Second), err)
	}

	fmt.Println()
	fmt.Println("Connection successful!")
	fmt.Println()
	fmt.Printf("Organization: %s\n", org.Name)
	fmt.Printf("ID: %s\n", org.ID)
	if org.Metrics != nil {
		fmt.Printf("Environments: %d\n", org.Metrics.ProjectCount)
		fmt.Printf("Versions: %d\n", org.Metrics.VersionCount)
		fmt.Printf("Components: %d\n", org.Metrics.ComponentCount)
	}

	return nil
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'lynk-mcp configure' first)", err)
	}

	token, err := config.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %w (run 'lynk-mcp configure' or set %s environment variable)", err, config.EnvTokenKey)
	}

	// Setup logger
	logger, err := setupLogger(cfg.Logging.Level)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}
	defer logger.Sync()

	// Create API client
	client := api.NewClient(cfg.API.Endpoint, token, cfg.API.Timeout)

	// Create and start MCP server
	server := mcp.NewServer(client, logger)

	logger.Info("starting MCP server", zap.String("endpoint", cfg.API.Endpoint))

	return server.Serve()
}

func setupLogger(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stderr"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return cfg.Build()
}
