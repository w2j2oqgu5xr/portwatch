package portecho_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portecho"
)

func TestSummarise_AllOpen(t *testing.T) {
	results := []portecho.Result{
		{Open: true, Latency: 1 * time.Millisecond},
		{Open: true, Latency: 3 * time.Millisecond},
		{Open: true, Latency: 2 * time.Millisecond},
	}
	s := portecho.Summarise(results)
	if s.Total != 3 {
		t.Errorf("Total: want 3, got %d", s.Total)
	}
	if s.Open != 3 {
		t.Errorf("Open: want 3, got %d", s.Open)
	}
	if s.Closed != 0 {
		t.Errorf("Closed: want 0, got %d", s.Closed)
	}
	if s.MinRTT != 1*time.Millisecond {
		t.Errorf("MinRTT: want 1ms, got %s", s.MinRTT)
	}
	if s.MaxRTT != 3*time.Millisecond {
		t.Errorf("MaxRTT: want 3ms, got %s", s.MaxRTT)
	}
	if s.MeanRTT != 2*time.Millisecond {
		t.Errorf("MeanRTT: want 2ms, got %s", s.MeanRTT)
	}
}

func TestSummarise_SomeClosed(t *testing.T) {
	results := []portecho.Result{
		{Open: true, Latency: 4 * time.Millisecond},
		{Open: false},
		{Open: false},
	}
	s := portecho.Summarise(results)
	if s.Open != 1 || s.Closed != 2 {
		t.Errorf("open/closed mismatch: %d/%d", s.Open, s.Closed)
	}
	if s.MeanRTT != 4*time.Millisecond {
		t.Errorf("MeanRTT: want 4ms, got %s", s.MeanRTT)
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := portecho.Summarise(nil)
	if s.Total != 0 {
		t.Errorf("expected zero Stats for empty input")
	}
}

func TestStats_String(t *testing.T) {
	s := portecho.Stats{
		Total:   5,
		Open:    3,
		Closed:  2,
		MinRTT:  500 * time.Microsecond,
		MeanRTT: 1 * time.Millisecond,
		MaxRTT:  2 * time.Millisecond,
	}
	out := s.String()
	if out == "" {
		t.Error("expected non-empty Stats.String()")
	}
}
