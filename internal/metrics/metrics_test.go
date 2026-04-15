package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestNew_StartsAtZero(t *testing.T) {
	c := metrics.New()
	if c.ScansTotal.Load() != 0 {
		t.Fatalf("expected ScansTotal=0, got %d", c.ScansTotal.Load())
	}
	if c.AlertsTotal.Load() != 0 {
		t.Fatalf("expected AlertsTotal=0, got %d", c.AlertsTotal.Load())
	}
}

func TestIncScans_Increments(t *testing.T) {
	c := metrics.New()
	c.IncScans()
	c.IncScans()
	if got := c.ScansTotal.Load(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestIncAlerts_Increments(t *testing.T) {
	c := metrics.New()
	c.IncAlerts()
	if got := c.AlertsTotal.Load(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestIncSuppressed_Increments(t *testing.T) {
	c := metrics.New()
	c.IncSuppressed()
	c.IncSuppressed()
	c.IncSuppressed()
	if got := c.Suppressed.Load(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestSetOpenPorts_StoresValue(t *testing.T) {
	c := metrics.New()
	c.SetOpenPorts(7)
	if got := c.PortsOpen.Load(); got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
	c.SetOpenPorts(2)
	if got := c.PortsOpen.Load(); got != 2 {
		t.Fatalf("expected 2 after update, got %d", got)
	}
}

func TestUptime_NonNegative(t *testing.T) {
	c := metrics.New()
	time.Sleep(10 * time.Millisecond)
	if c.Uptime() < 0 {
		t.Fatal("uptime should not be negative")
	}
}

func TestWriteTo_ContainsFields(t *testing.T) {
	c := metrics.New()
	c.IncScans()
	c.IncAlerts()
	c.IncSuppressed()
	c.SetOpenPorts(4)

	var buf bytes.Buffer
	_, err := c.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"scans=1", "alerts=1", "suppressed=1", "open_ports=4", "uptime="} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}
