// Package porttrend tracks how port open/close counts change over time,
// providing a rolling window of observations for trend analysis.
package porttrend

import (
	"sync"
	"time"
)

// Sample is a single observation recorded at a point in time.
type Sample struct {
	At       time.Time
	OpenCount int
}

// Trend summarises direction of change across the window.
type Trend int

const (
	TrendStable  Trend = 0
	TrendRising  Trend = 1
	TrendFalling Trend = -1
)

func (t Trend) String() string {
	switch t {
	case TrendRising:
		return "rising"
	case TrendFalling:
		return "falling"
	default:
		return "stable"
	}
}

// Tracker maintains a rolling window of port-count samples.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	samples []Sample
}

// New returns a Tracker that retains samples within the given window duration.
// A zero or negative window defaults to 5 minutes.
func New(window time.Duration) *Tracker {
	if window <= 0 {
		window = 5 * time.Minute
	}
	return &Tracker{window: window}
}

// Record adds a new observation. Samples older than the window are evicted.
func (t *Tracker) Record(openCount int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	t.samples = append(t.samples, Sample{At: now, OpenCount: openCount})
	t.evict(now)
}

// Samples returns a copy of the current window of samples.
func (t *Tracker) Samples() []Sample {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Sample, len(t.samples))
	copy(out, t.samples)
	return out
}

// Current returns the Trend direction by comparing the oldest and newest
// samples in the window. Returns TrendStable when fewer than two samples exist.
func (t *Tracker) Current() Trend {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.samples) < 2 {
		return TrendStable
	}
	first := t.samples[0].OpenCount
	last := t.samples[len(t.samples)-1].OpenCount
	switch {
	case last > first:
		return TrendRising
	case last < first:
		return TrendFalling
	default:
		return TrendStable
	}
}

// evict removes samples that fall outside the retention window. Caller must hold mu.
func (t *Tracker) evict(now time.Time) {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(t.samples) && t.samples[i].At.Before(cutoff) {
		i++
	}
	t.samples = t.samples[i:]
}
