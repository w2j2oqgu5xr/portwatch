package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/retry"
)

var errTemp = errors.New("temporary error")

func TestDo_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultPolicy(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnError(t *testing.T) {
	calls := 0
	p := retry.Policy{MaxAttempts: 3, InitialDelay: time.Millisecond, Multiplier: 1}
	err := retry.Do(context.Background(), p, func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ReturnsLastError(t *testing.T) {
	p := retry.Policy{MaxAttempts: 2, InitialDelay: time.Millisecond, Multiplier: 1}
	err := retry.Do(context.Background(), p, func() error {
		return errTemp
	})
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected errTemp, got %v", err)
	}
}

func TestDo_RespectsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := retry.Policy{MaxAttempts: 5, InitialDelay: time.Millisecond, Multiplier: 1}
	err := retry.Do(ctx, p, func() error { return errTemp })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDo_ZeroMaxAttemptsRunsOnce(t *testing.T) {
	calls := 0
	p := retry.Policy{MaxAttempts: 0, InitialDelay: time.Millisecond, Multiplier: 1}
	retry.Do(context.Background(), p, func() error { calls++; return errTemp }) //nolint
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}
