package portrank_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/portrank"
)

func TestNew_DefaultBaseScore(t *testing.T) {
	r := portrank.New(0)
	scores := r.Rank([]int{80})
	if len(scores) != 1 {
		t.Fatalf("expected 1 score, got %d", len(scores))
	}
	if scores[0].Value <= 0 {
		t.Errorf("expected positive score, got %f", scores[0].Value)
	}
}

func TestRank_SortedByScore(t *testing.T) {
	r := portrank.New(1.0)
	r.SetWeight(443, 5.0)
	r.SetWeight(80, 2.0)
	scores := r.Rank([]int{80, 443, 22})
	if scores[0].Port != 443 {
		t.Errorf("expected 443 ranked first, got %d", scores[0].Port)
	}
	if scores[1].Port != 80 {
		t.Errorf("expected 80 ranked second, got %d", scores[1].Port)
	}
}

func TestRank_AssignsRankNumbers(t *testing.T) {
	r := portrank.New(1.0)
	scores := r.Rank([]int{22, 80, 443})
	for i, s := range scores {
		if s.Rank != i+1 {
			t.Errorf("rank %d: expected %d, got %d", i, i+1, s.Rank)
		}
	}
}

func TestRecordChange_IncreasesScore(t *testing.T) {
	r := portrank.New(1.0)
	r.SetWeight(8080, 1.0)
	before := r.Rank([]int{8080})[0].Value
	r.RecordChange(8080)
	after := r.Rank([]int{8080})[0].Value
	if after <= before {
		t.Errorf("expected score to increase after change: before=%f after=%f", before, after)
	}
}

func TestReset_ClearsChanges(t *testing.T) {
	r := portrank.New(1.0)
	r.SetWeight(9090, 1.0)
	r.RecordChange(9090)
	r.RecordChange(9090)
	r.Reset()
	scores := r.Rank([]int{9090})
	if scores[0].Value != 1.0 {
		t.Errorf("expected score 1.0 after reset, got %f", scores[0].Value)
	}
}

func TestSetWeight_NegativeClampedToZero(t *testing.T) {
	r := portrank.New(1.0)
	r.SetWeight(1234, -5.0)
	scores := r.Rank([]int{1234})
	if scores[0].Value < 0 {
		t.Errorf("expected non-negative score, got %f", scores[0].Value)
	}
}

func TestRank_TieBreakByPort(t *testing.T) {
	r := portrank.New(2.0)
	scores := r.Rank([]int{9000, 8000, 7000})
	// All have equal weight; tie-break should order by port number ascending.
	if scores[0].Port != 7000 {
		t.Errorf("expected 7000 first in tie-break, got %d", scores[0].Port)
	}
}
