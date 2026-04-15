package report_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/report"
	"github.com/user/portwatch/internal/snapshot"
)

func baseEvents() []history.Event {
	return []history.Event{
		{Port: 8080, Kind: "opened", Timestamp: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)},
		{Port: 22, Kind: "closed", Timestamp: time.Date(2024, 1, 1, 10, 1, 0, 0, time.UTC)},
	}
}

func TestWrite_ContainsOpenPorts(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf)
	s := report.BuildSummary(
		[]int{80, 443, 8080},
		snapshot.Diff{Opened: []int{8080}, Closed: []int{}},
		nil,
	)
	if err := w.Write(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"80", "443", "8080"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestWrite_ShowsOpenedAndClosed(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf)
	s := report.BuildSummary(
		[]int{80},
		snapshot.Diff{Opened: []int{80}, Closed: []int{22}},
		nil,
	)
	_ = w.Write(s)
	out := buf.String()
	if !strings.Contains(out, "Newly opened") {
		t.Error("expected 'Newly opened' section")
	}
	if !strings.Contains(out, "Newly closed") {
		t.Error("expected 'Newly closed' section")
	}
}

func TestWrite_NoOpenPorts(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf)
	s := report.BuildSummary([]int{}, snapshot.Diff{}, nil)
	_ = w.Write(s)
	if !strings.Contains(buf.String(), "none") {
		t.Error("expected 'none' when no ports are open")
	}
}

func TestWrite_IncludesRecentEvents(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf)
	s := report.BuildSummary([]int{8080}, snapshot.Diff{}, baseEvents())
	_ = w.Write(s)
	out := buf.String()
	if !strings.Contains(out, "Recent events") {
		t.Error("expected 'Recent events' section")
	}
	if !strings.Contains(out, "8080") {
		t.Error("expected port 8080 in events")
	}
}

func TestBuildSummary_SetsGeneratedAt(t *testing.T) {
	before := time.Now()
	s := report.BuildSummary(nil, snapshot.Diff{}, nil)
	after := time.Now()
	if s.GeneratedAt.Before(before) || s.GeneratedAt.After(after) {
		t.Errorf("GeneratedAt %v not within expected range", s.GeneratedAt)
	}
}
