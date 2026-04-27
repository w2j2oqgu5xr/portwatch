package portweight_test

import (
	"testing"

	"github.com/user/portwatch/internal/portweight"
)

func TestNew_DefaultBaseScore(t *testing.T) {
	s := portweight.New(nil)
	w := s.Compute(80)
	if w.Score != 0 {
		t.Fatalf("expected 0, got %f", w.Score)
	}
	if w.Reason != "base" {
		t.Fatalf("expected reason 'base', got %q", w.Reason)
	}
}

func TestObserveOpen_IncreasesScore(t *testing.T) {
	s := portweight.New(nil)
	s.ObserveOpen(443)
	s.ObserveOpen(443)
	w := s.Compute(443)
	// 0.5 * 2 = 1.0
	if w.Score != 1.0 {
		t.Fatalf("expected 1.0, got %f", w.Score)
	}
	if w.Reason != "open-frequency" {
		t.Fatalf("unexpected reason %q", w.Reason)
	}
}

func TestRecordChange_IncreasesScore(t *testing.T) {
	s := portweight.New(nil)
	s.RecordChange(22)
	w := s.Compute(22)
	// 2.0 * 1 = 2.0
	if w.Score != 2.0 {
		t.Fatalf("expected 2.0, got %f", w.Score)
	}
	if w.Reason != "change-activity" {
		t.Fatalf("unexpected reason %q", w.Reason)
	}
}

func TestBaseScore_CombinesWithObserved(t *testing.T) {
	s := portweight.New(map[int]float64{80: 5.0})
	s.ObserveOpen(80)
	w := s.Compute(80)
	// 5.0 + 0.5*1 = 5.5
	if w.Score != 5.5 {
		t.Fatalf("expected 5.5, got %f", w.Score)
	}
}

func TestComputeAll_SortedDescending(t *testing.T) {
	s := portweight.New(map[int]float64{443: 1.0, 22: 3.0, 80: 2.0})
	weights := s.ComputeAll()
	if len(weights) != 3 {
		t.Fatalf("expected 3 weights, got %d", len(weights))
	}
	for i := 1; i < len(weights); i++ {
		if weights[i].Score > weights[i-1].Score {
			t.Fatalf("weights not sorted: index %d (%f) > index %d (%f)",
				i, weights[i].Score, i-1, weights[i-1].Score)
		}
	}
	if weights[0].Port != 22 {
		t.Fatalf("expected port 22 first, got %d", weights[0].Port)
	}
}

func TestComputeAll_EmptyScorer(t *testing.T) {
	s := portweight.New(nil)
	if ws := s.ComputeAll(); len(ws) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(ws))
	}
}
