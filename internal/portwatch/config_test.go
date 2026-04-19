package portwatch

import (
	"testing"
	"time"
)

func TestValidate_AppliesDefaults(t *testing.T) {
	cfg := Config{Host: "localhost"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Ports.From != defaultFrom {
		t.Errorf("From: want %d, got %d", defaultFrom, cfg.Ports.From)
	}
	if cfg.Ports.To != defaultTo {
		t.Errorf("To: want %d, got %d", defaultTo, cfg.Ports.To)
	}
	if cfg.Interval != minInterval {
		t.Errorf("Interval: want %v, got %v", minInterval, cfg.Interval)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	cfg := Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := Config{
		Host:     "10.0.0.1",
		Ports:    PortRange{From: 80, To: 443},
		Interval: 30 * time.Second,
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_InvertedRange(t *testing.T) {
	cfg := Config{
		Host:  "localhost",
		Ports: PortRange{From: 500, To: 100},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestValidate_IntervalBelowMin(t *testing.T) {
	cfg := Config{
		Host:     "localhost",
		Interval: 1 * time.Second,
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != minInterval {
		t.Errorf("expected interval clamped to %v, got %v", minInterval, cfg.Interval)
	}
}
