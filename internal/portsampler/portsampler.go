// Package portsampler provides periodic sampling of open ports with
// statistical aggregation over a configurable window.
package portsampler

import (
	"sync"
	"time"
)

// Sample holds a single observation of open ports at a point in time.
type Sample struct {
	Timestamp time.Time
	Ports     []int
}

// Sampler collects port samples and exposes aggregate statistics.
type Sampler struct {
	mu      sync.Mutex
	window  int
	samples []Sample
}

// New creates a Sampler that retains at most window samples.
// If window is less than 1 it is clamped to 1.
func New(window int) *Sampler {
	if window < 1 {
		window = 1
	}
	return &Sampler{window: window}
}

// Record appends a new sample, evicting the oldest when the window is full.
func (s *Sampler) Record(ports []int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cp := make([]int, len(ports))
	copy(cp, ports)

	s.samples = append(s.samples, Sample{Timestamp: time.Now(), Ports: cp})
	if len(s.samples) > s.window {
		s.samples = s.samples[len(s.samples)-s.window:]
	}
}

// Latest returns the most recent sample, or an empty Sample if none exist.
func (s *Sampler) Latest() Sample {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) == 0 {
		return Sample{}
	}
	return s.samples[len(s.samples)-1]
}

// All returns a copy of all retained samples in chronological order.
func (s *Sampler) All() []Sample {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Sample, len(s.samples))
	copy(out, s.samples)
	return out
}

// AverageCount returns the mean number of open ports across all retained samples.
// Returns 0 if no samples have been recorded.
func (s *Sampler) AverageCount() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) == 0 {
		return 0
	}
	var total int
	for _, sm := range s.samples {
		total += len(sm.Ports)
	}
	return float64(total) / float64(len(s.samples))
}

// Reset discards all recorded samples.
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.samples = nil
}
