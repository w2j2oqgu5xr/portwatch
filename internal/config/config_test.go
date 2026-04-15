package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-cfg-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `{"host":"127.0.0.1","ports":[80,443],"interval":"10s"}`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("host: got %q, want %q", cfg.Host, "127.0.0.1")
	}
	if len(cfg.Ports) != 2 {
		t.Errorf("ports length: got %d, want 2", len(cfg.Ports))
	}
	if cfg.Interval.Duration != 10*time.Second {
		t.Errorf("interval: got %v, want 10s", cfg.Interval.Duration)
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	path := writeTemp(t, `{"ports":[8080]}`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Host != "localhost" {
		t.Errorf("default host: got %q, want %q", cfg.Host, "localhost")
	}
	if cfg.Interval.Duration != 5*time.Second {
		t.Errorf("default interval: got %v, want 5s", cfg.Interval.Duration)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestValidate_NoPorts(t *testing.T) {
	cfg := &config.Config{Host: "localhost"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty ports, got nil")
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	cfg := &config.Config{Ports: []int{0}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port 0, got nil")
	}
}

func TestDuration_RoundTrip(t *testing.T) {
	d := config.Duration{Duration: 30 * time.Second}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var d2 config.Duration
	if err := json.Unmarshal(b, &d2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if d.Duration != d2.Duration {
		t.Errorf("round-trip: got %v, want %v", d2.Duration, d.Duration)
	}
}
