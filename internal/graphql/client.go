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

package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/interlynk-io/lynk-mcp/internal/retry"
)

// HTTPError represents an HTTP error response with a status code.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("unexpected status code %d: %s", e.StatusCode, e.Body)
}

// Client is a GraphQL client for the Lynk API
type Client struct {
	endpoint    string
	token       string
	httpClient  *http.Client
	retryConfig retry.Config
}

// NewClient creates a new GraphQL client
func NewClient(endpoint, token string, timeout time.Duration) *Client {
	return &Client{
		endpoint: endpoint,
		token:    token,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		retryConfig: retry.DefaultTransientConfig(),
	}
}

// Request represents a GraphQL request
type Request struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// Response represents a GraphQL response
type Response struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message    string                 `json:"message"`
	Locations  []ErrorLocation        `json:"locations,omitempty"`
	Path       []interface{}          `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// ErrorLocation represents the location of a GraphQL error
type ErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// isTransientError returns true for errors that are worth retrying:
// 5xx server errors, 429 rate limiting, and network-level errors.
func isTransientError(err error) bool {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode >= 500 || httpErr.StatusCode == http.StatusTooManyRequests
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	return false
}

// Execute executes a GraphQL query with automatic retry for transient errors.
func (c *Client) Execute(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	return retry.Do(ctx, c.retryConfig, func() error {
		return c.executeOnce(ctx, query, variables, result)
	}, isTransientError, nil)
}

// executeOnce performs a single GraphQL request.
func (c *Client) executeOnce(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	req := Request{
		Query:     query,
		Variables: variables,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)
	httpReq.Header.Set("User-Agent", "lynk-mcp/1.0")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &HTTPError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	var graphQLResp Response
	if err := json.Unmarshal(respBody, &graphQLResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(graphQLResp.Errors) > 0 {
		return fmt.Errorf("GraphQL error: %s", graphQLResp.Errors[0].Message)
	}

	if result != nil && graphQLResp.Data != nil {
		if err := json.Unmarshal(graphQLResp.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	return nil
}
