package watchdog_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestWatchdog_CallsOnStallWhenNoHeartbeat(t *testing.T) {
	beat := make(chan struct{}, 1)
	var called atomic.Int32
	wd := watchdog.New(beat, 50*time.Millisecond, func() {
		called.Add(1)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	wd.Run(ctx)
	if called.Load() == 0 {
		t.Fatal("expected onStall to be called")
	}
}

func TestWatchdog_NoStallWithRegularHeartbeats(t *testing.T) {
	beat := make(chan struct{}, 1)
	var called atomic.Int32
	wd := watchdog.New(beat, 100*time.Millisecond, func() {
		called.Add(1)
	})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		wd.Run(ctx)
		close(done)
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case beat <- struct{}{}:
				time)
			.Sleep(250 * time.Millisecondcel()
	 called.Load() != 0 {
		t.Fatalf("expected no stall, got %d", called.Load())
	}
}

func TestWatchdog_DefaultTimeoutUsedWhenZero(t *testing.T) {
	beat := make(chan struct{}, 1)
	wd := watchdog.New(beat, 0, func() {})
	if wd == nil {
		t.Fatal("expected non-nil watchdog")
	}
}

func TestWatchdog_StopsOnContextCancel(t *testing.T) {
	beat := make(chan struct{}, 1)
	wd := watchdog.New(beat, 50*time.Millisecond, func() {})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		wd.Run(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("watchdog did not stop after context cancel")
	}
}
