package daemon_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
)

func minimalConfig() *config.Config {
	return &config.Config{
		Host:            "127.0.0.1",
		IntervalSeconds: 1,
		Ports:           []int{12340, 12341},
		Notify:          config.NotifyConfig{Targets: nil},
	}
}

func TestNew_ValidConfig(t *testing.T) {
	cfg := minimalConfig()
	d, err := daemon.New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("New() returned nil daemon")
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	cfg := minimalConfig()
	d, err := daemon.New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	err = d.Run(ctx)
	if err == nil {
		t.Fatal("Run() expected context error, got nil")
	}
	// context.DeadlineExceeded or context.Canceled are both acceptable.
	if ctx.Err() == nil {
		t.Errorf("Run() returned %v but context shows no error", err)
	}
}

func TestMetrics_IncrementsAfterRun(t *testing.T) {
	cfg := minimalConfig()
	cfg.IntervalSeconds = 1

	d, err := daemon.New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	defer cancel()

	_ = d.Run(ctx) // runs for ~1 s, should complete at least one scan tick

	if d.Metrics().Scans() == 0 {
		t.Error("expected at least one scan to be recorded in metrics")
	}
}
