package tui_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/tui"
)

func newRenderer(t *testing.T) (*tui.Renderer, *strings.Builder) {
	t.Helper()
	var buf strings.Builder
	m := metrics.New()
	return tui.New(&buf, m), &buf
}

func TestRender_ContainsHeader(t *testing.T) {
	r, buf := newRenderer(t)
	if err := r.Render(nil, time.Now()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "portwatch status") {
		t.Error("expected header in output")
	}
}

func TestRender_ShowsOpenPorts(t *testing.T) {
	r, buf := newRenderer(t)
	_ = r.Render([]int{8080, 443, 80}, time.Now())
	out := buf.String()
	for _, p := range []string{"80", "443", "8080"} {
		if !strings.Contains(out, p) {
			t.Errorf("expected port %s in output", p)
		}
	}
}

func TestRender_NoPorts(t *testing.T) {
	r, buf := newRenderer(t)
	_ = r.Render([]int{}, time.Now())
	if !strings.Contains(buf.String(), "none") {
		t.Error("expected 'none' when no open ports")
	}
}

func TestRender_MetricsDisplayed(t *testing.T) {
	var buf strings.Builder
	m := metrics.New()
	m.IncScans()
	m.IncScans()
	m.IncAlerts()
	r := tui.New(&buf, m)
	_ = r.Render(nil, time.Now())
	out := buf.String()
	if !strings.Contains(out, "2") {
		t.Error("expected scan count 2 in output")
	}
}

func TestRender_LongPortListTruncated(t *testing.T) {
	r, buf := newRenderer(t)
	ports := []int{80, 443, 8080, 8443, 9000, 9090, 3000, 3001}
	_ = r.Render(ports, time.Now())
	if !strings.Contains(buf.String(), "...") {
		t.Error("expected truncation indicator for long port list")
	}
}
