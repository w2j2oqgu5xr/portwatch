package porttrend

import (
	"testing"
	"time"
)

func TestNew_DefaultWindow(t *testing.T) {
	tr := New(0)
	if tr.window != 5*time.Minute {
		t.Fatalf("expected 5m window, got %v", tr.window)
	}
}

func TestRecord_StoresSample(t *testing.T) {
	tr := New(time.Minute)
	tr.Record(3)
	samples := tr.Samples()
	if len(samples) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(samples))
	}
	if samples[0].OpenCount != 3 {
		t.Errorf("expected OpenCount 3, got %d", samples[0].OpenCount)
	}
}

func TestCurrent_StableWhenFewSamples(t *testing.T) {
	tr := New(time.Minute)
	if tr.Current() != TrendStable {
		t.Error("expected stable with no samples")
	}
	tr.Record(5)
	if tr.Current() != TrendStable {
		t.Error("expected stable with single sample")
	}
}

func TestCurrent_Rising(t *testing.T) {
	tr := New(time.Minute)
	tr.Record(2)
	tr.Record(5)
	if got := tr.Current(); got != TrendRising {
		t.Errorf("expected rising, got %v", got)
	}
}

func TestCurrent_Falling(t *testing.T) {
	tr := New(time.Minute)
	tr.Record(10)
	tr.Record(4)
	if got := tr.Current(); got != TrendFalling {
		t.Errorf("expected falling, got %v", got)
	}
}

func TestCurrent_Stable(t *testing.T) {
	tr := New(time.Minute)
	tr.Record(7)
	tr.Record(7)
	if got := tr.Current(); got != TrendStable {
		t.Errorf("expected stable, got %v", got)
	}
}

func TestEvict_RemovesOldSamples(t *testing.T) {
	tr := New(50 * time.Millisecond)
	tr.Record(1)
	time.Sleep(60 * time.Millisecond)
	tr.Record(2) // this call triggers eviction
	samples := tr.Samples()
	if len(samples) != 1 {
		t.Fatalf("expected 1 sample after eviction, got %d", len(samples))
	}
	if samples[0].OpenCount != 2 {
		t.Errorf("expected surviving sample OpenCount=2, got %d", samples[0].OpenCount)
	}
}

func TestTrend_String(t *testing.T) {
	cases := []struct {
		trend Trend
		want  string
	}{
		{TrendStable, "stable"},
		{TrendRising, "rising"},
		{TrendFalling, "falling"},
	}
	for _, c := range cases {
		if got := c.trend.String(); got != c.want {
			t.Errorf("Trend(%d).String() = %q, want %q", c.trend, got, c.want)
		}
	}
}
