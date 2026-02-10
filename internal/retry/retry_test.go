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
	"errors"
	"testing"
	"time"
)

func TestDo_SuccessOnFirstTry(t *testing.T) {
	calls := 0
	err := Do(context.Background(), Config{
		MaxRetries:    3,
		InitialDelay:  time.Millisecond,
		MaxDelay:      10 * time.Millisecond,
		BackoffFactor: 2.0,
	}, func() error {
		calls++
		return nil
	}, nil, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_SuccessAfterRetries(t *testing.T) {
	calls := 0
	err := Do(context.Background(), Config{
		MaxRetries:    5,
		InitialDelay:  time.Millisecond,
		MaxDelay:      10 * time.Millisecond,
		BackoffFactor: 2.0,
	}, func() error {
		calls++
		if calls < 3 {
			return errors.New("transient error")
		}
		return nil
	}, nil, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_AllRetriesExhausted(t *testing.T) {
	calls := 0
	err := Do(context.Background(), Config{
		MaxRetries:    3,
		InitialDelay:  time.Millisecond,
		MaxDelay:      10 * time.Millisecond,
		BackoffFactor: 2.0,
	}, func() error {
		calls++
		return errors.New("persistent error")
	}, nil, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// 1 initial + 3 retries = 4
	if calls != 4 {
		t.Fatalf("expected 4 calls, got %d", calls)
	}
}

func TestDo_NonRetryableError(t *testing.T) {
	authErr := errors.New("401 unauthorized")
	calls := 0

	err := Do(context.Background(), Config{
		MaxRetries:    5,
		InitialDelay:  time.Millisecond,
		MaxDelay:      10 * time.Millisecond,
		BackoffFactor: 2.0,
	}, func() error {
		calls++
		return authErr
	}, func(err error) bool {
		return false // nothing is retryable
	}, nil)

	if !errors.Is(err, authErr) {
		t.Fatalf("expected auth error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call (no retry), got %d", calls)
	}
}

func TestDo_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	err := Do(ctx, Config{
		MaxRetries:    100,
		InitialDelay:  50 * time.Millisecond,
		MaxDelay:      50 * time.Millisecond,
		BackoffFactor: 1.0,
	}, func() error {
		calls++
		return errors.New("keep failing")
	}, nil, nil)

	if err == nil {
		t.Fatal("expected error from context cancellation")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDo_OnRetryCallback(t *testing.T) {
	var retryAttempts []int
	var retryErrors []error

	targetErr := errors.New("fail")

	err := Do(context.Background(), Config{
		MaxRetries:    3,
		InitialDelay:  time.Millisecond,
		MaxDelay:      10 * time.Millisecond,
		BackoffFactor: 2.0,
	}, func() error {
		return targetErr
	}, nil, func(attempt int, err error, delay time.Duration) {
		retryAttempts = append(retryAttempts, attempt)
		retryErrors = append(retryErrors, err)
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if len(retryAttempts) != 3 {
		t.Fatalf("expected 3 onRetry calls, got %d", len(retryAttempts))
	}
	for i, a := range retryAttempts {
		if a != i+1 {
			t.Fatalf("expected attempt %d, got %d", i+1, a)
		}
	}
	for _, e := range retryErrors {
		if !errors.Is(e, targetErr) {
			t.Fatalf("expected target error in callback, got %v", e)
		}
	}
}

func TestDo_BackoffCappedAtMaxDelay(t *testing.T) {
	var delays []time.Duration

	_ = Do(context.Background(), Config{
		MaxRetries:    3,
		InitialDelay:  10 * time.Millisecond,
		MaxDelay:      15 * time.Millisecond,
		BackoffFactor: 10.0,
	}, func() error {
		return errors.New("fail")
	}, nil, func(attempt int, err error, delay time.Duration) {
		delays = append(delays, delay)
	})

	for i, d := range delays {
		if d > 15*time.Millisecond {
			t.Fatalf("delay %d exceeded max: %v", i, d)
		}
	}
}
