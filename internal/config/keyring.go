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
	"strings"

	"github.com/99designs/keyring"
)

const (
	serviceName = "lynk-mcp"
	tokenKey    = "api-token"
)

// ValidateTokenFormat checks if the token has a valid format
func ValidateTokenFormat(token string) bool {
	validPrefixes := []string{
		"lynk_live_",
		"lynk_staging_",
		"lynk_test_",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(token, prefix) {
			return true
		}
	}

	return false
}

// StoreToken securely stores the API token in the system keychain
func StoreToken(token string) error {
	ring, err := openKeyring()
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}

	err = ring.Set(keyring.Item{
		Key:  tokenKey,
		Data: []byte(token),
	})
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	return nil
}

// EnvTokenKey is the environment variable name for the API token
const EnvTokenKey = "LYNK_API_TOKEN"

// GetToken retrieves the API token, checking environment variable first,
// then falling back to the system keychain. This allows Docker/CI usage
// where a keyring is not available.
func GetToken() (string, error) {
	// Check environment variable first (useful for Docker/CI environments)
	if token := os.Getenv(EnvTokenKey); token != "" {
		if !ValidateTokenFormat(token) {
			return "", fmt.Errorf("invalid token format in %s: must start with lynk_live_, lynk_staging_, or lynk_test_", EnvTokenKey)
		}
		return token, nil
	}

	// Fall back to keyring
	ring, err := openKeyring()
	if err != nil {
		return "", fmt.Errorf("failed to open keyring: %w", err)
	}

	item, err := ring.Get(tokenKey)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	return string(item.Data), nil
}

// DeleteToken removes the API token from the system keychain
func DeleteToken() error {
	ring, err := openKeyring()
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}

	err = ring.Remove(tokenKey)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

// openKeyring opens or creates the keyring for storing credentials
func openKeyring() (keyring.Keyring, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return keyring.Open(keyring.Config{
		ServiceName: serviceName,

		// macOS Keychain
		KeychainName:             "login",
		KeychainTrustApplication: true,

		// Linux Secret Service or file-based fallback
		FileDir:          configDir,
		FilePasswordFunc: keyring.TerminalPrompt,

		// Windows Credential Manager - no extra config needed
	})
}
