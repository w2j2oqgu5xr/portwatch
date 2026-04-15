package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, "host: 127.0.0.1\nports: [80, 443]\ninterval: 10s\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", cfg.Host)
	}
	if len(cfg.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(cfg.Ports))
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("expected 10s interval, got %v", cfg.Interval)
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	path := writeTemp(t, "ports: [22]\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Host != config.DefaultHost {
		t.Errorf("expected default host, got %s", cfg.Host)
	}
	if cfg.Interval != config.DefaultInterval {
		t.Errorf("expected default interval, got %v", cfg.Interval)
	}
	if cfg.SnapshotDB != config.DefaultSnapshotDB {
		t.Errorf("expected default snapshot db, got %s", cfg.SnapshotDB)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidate_NoPorts(t *testing.T) {
	path := writeTemp(t, "host: localhost\ninterval: 5s\n")
	_, err := config.Load(path)
	if err == nil {
		t.Error("expected validation error when no ports specified")
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	path := writeTemp(t, "ports: [0]\n")
	_, err := config.Load(path)
	if err == nil {
		t.Error("expected validation error for port 0")
	}
}

func TestValidate_IntervalTooShort(t *testing.T) {
	path := writeTemp(t, "ports: [80]\ninterval: 500ms\n")
	_, err := config.Load(path)
	if err == nil {
		t.Error("expected validation error for interval < 1s")
	}
}
