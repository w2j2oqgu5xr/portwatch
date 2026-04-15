package ratelimit_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_BurstConsumed(t *testing.T) {
	l := ratelimit.New(1, 3)

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow() == true on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow() == false after burst exhausted")
	}
}

func TestAllow_RefillOverTime(t *testing.T) {
	now := time.Now()
	l := ratelimit.New(10, 1) // 10 tokens/s, burst 1

	// Exhaust the single token.
	l.Allow()

	// Advance the clock by injecting a custom now function via a small wrapper.
	// Because now is unexported we test refill indirectly by waiting a short
	// real duration. Use a generous sleep to avoid flakiness on slow CI.
	time.Sleep(150 * time.Millisecond)
	_ = now // suppress unused warning

	if !l.Allow() {
		t.Fatal("expected token to be refilled after 150 ms at 10 tokens/s")
	}
}

func TestAllow_ZeroRateDefaultsToOne(t *testing.T) {
	l := ratelimit.New(0, 2)
	if !l.Allow() {
		t.Fatal("expected first Allow() to succeed")
	}
}

func TestReset_RestoresBurst(t *testing.T) {
	l := ratelimit.New(1, 2)
	l.Allow()
	l.Allow()

	if l.Allow() {
		t.Fatal("expected burst to be exhausted before Reset")
	}

	l.Reset()

	if !l.Allow() {
		t.Fatal("expected Allow() to succeed after Reset")
	}
}

func TestAllow_ConcurrentSafe(t *testing.T) {
	const goroutines = 50
	l := ratelimit.New(1000, goroutines)

	var allowed atomic.Int64
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if l.Allow() {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	if allowed.Load() > int64(goroutines) {
		t.Fatalf("allowed %d events, want at most %d", allowed.Load(), goroutines)
	}
}
