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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

const (
	DefaultEndpoint = "https://api.interlynk.io/lynkapi"
	DefaultTimeout  = 30 * time.Second
	DefaultLogLevel = "info"

	configDir  = ".lynk-mcp"
	configFile = "config.yaml"
)

// Config holds the application configuration
type Config struct {
	API     APIConfig     `mapstructure:"api"`
	Logging LoggingConfig `mapstructure:"logging"`
}

// APIConfig holds API-related configuration
type APIConfig struct {
	Endpoint string        `mapstructure:"endpoint"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

// DefaultConfig returns a config with default values
func DefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			Endpoint: DefaultEndpoint,
			Timeout:  DefaultTimeout,
		},
		Logging: LoggingConfig{
			Level: DefaultLogLevel,
		},
	}
}

// Load loads configuration from file, environment variables, and defaults
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("api.endpoint", DefaultEndpoint)
	viper.SetDefault("api.timeout", DefaultTimeout)
	viper.SetDefault("logging.level", DefaultLogLevel)

	// Environment variable overrides
	viper.SetEnvPrefix("LYNK_MCP")
	viper.AutomaticEnv()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Return error only if it's not a "file not found" error
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to read config: %w", err)
			}
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.Set("api.endpoint", c.API.Endpoint)
	viper.Set("api.timeout", c.API.Timeout)
	viper.Set("logging.level", c.Logging.Level)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, configDir, configFile), nil
}

// GetConfigDir returns the config directory path
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, configDir), nil
}
