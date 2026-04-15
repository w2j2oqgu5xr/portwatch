package schedule_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/schedule"
)

func TestNew_BelowMinInterval(t *testing.T) {
	_, err := schedule.New(context.Background(), 100*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for interval below minimum, got nil")
	}
}

func TestNew_ValidInterval(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker, err := schedule.New(ctx, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ticker.Stop()

	if ticker.Interval() != time.Second {
		t.Errorf("expected interval 1s, got %v", ticker.Interval())
	}
}

func TestTicker_EmitsTicks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker, err := schedule.New(ctx, schedule.MinInterval)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ticker.Stop()

	select {
	case _, ok := <-ticker.C:
		if !ok {
			t.Fatal("channel closed unexpectedly")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for first tick")
	}
}

func TestTicker_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ticker, err := schedule.New(ctx, schedule.MinInterval)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cancel()

	time.Sleep(100 * time.Millisecond)

	// drain any buffered tick then confirm channel closes
	for range ticker.C {
	}
	// reaching here means channel was closed — pass
}

func TestTicker_StopMethod(t *testing.T) {
	ctx := context.Background()

	ticker, err := schedule.New(ctx, schedule.MinInterval)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ticker.Stop()

	time.Sleep(100 * time.Millisecond)

	for range ticker.C {
	}
}
