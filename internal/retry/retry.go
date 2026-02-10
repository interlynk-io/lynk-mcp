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

package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Config controls retry behavior.
type Config struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	BackoffFactor float64
}

// DefaultTransientConfig returns a config for retrying transient errors
// at the GraphQL client level (3 retries, 1s initial, 10s max, 2x backoff).
func DefaultTransientConfig() Config {
	return Config{
		MaxRetries:    3,
		InitialDelay:  1 * time.Second,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
	}
}

// DefaultVerifyConfig returns a config for the verify command retry loop
// (24 retries, 10s initial, 15s max, 1.5x backoff, ~6 min total).
func DefaultVerifyConfig() Config {
	return Config{
		MaxRetries:    24,
		InitialDelay:  10 * time.Second,
		MaxDelay:      15 * time.Second,
		BackoffFactor: 1.5,
	}
}

// Do executes fn with retries according to cfg.
//
// shouldRetry decides whether a given error is retryable. If nil, all errors
// are retried.
//
// onRetry is called before each retry sleep with the attempt number (1-based),
// the error from the last attempt, and the delay before the next attempt.
// It may be nil.
func Do(ctx context.Context, cfg Config, fn func() error, shouldRetry func(error) bool, onRetry func(attempt int, err error, delay time.Duration)) error {
	var lastErr error

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if shouldRetry != nil && !shouldRetry(lastErr) {
			return lastErr
		}

		if attempt == cfg.MaxRetries {
			break
		}

		delay := time.Duration(float64(cfg.InitialDelay) * math.Pow(cfg.BackoffFactor, float64(attempt)))
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}

		if onRetry != nil {
			onRetry(attempt+1, lastErr, delay)
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("%w (last error: %v)", ctx.Err(), lastErr)
		case <-time.After(delay):
		}
	}

	return fmt.Errorf("after %d attempts: %w", cfg.MaxRetries+1, lastErr)
}
