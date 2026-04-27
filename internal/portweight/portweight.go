// Package portweight assigns a numeric weight to ports based on
// their observed open frequency, service criticality, and recent
// change activity. Weights can be used by ranking and alerting
// pipelines to prioritise high-value ports.
package portweight

import "sync"

// Weight holds the computed weight for a single port.
type Weight struct {
	Port   int
	Score  float64
	Reason string
}

// Scorer computes weights for ports.
type Scorer struct {
	mu       sync.RWMutex
	base     map[int]float64 // caller-supplied base scores
	freq     map[int]int     // open-observation count
	changes  map[int]int     // state-change count
}

// New returns a Scorer with optional base scores.
// Base scores are merged with observed data at compute time.
func New(base map[int]float64) *Scorer {
	b := make(map[int]float64, len(base))
	for k, v := range base {
		b[k] = v
	}
	return &Scorer{
		base:    b,
		freq:    make(map[int]int),
		changes: make(map[int]int),
	}
}

// ObserveOpen records that port was seen open during a scan cycle.
func (s *Scorer) ObserveOpen(port int) {
	s.mu.Lock()
	s.freq[port]++
	s.mu.Unlock()
}

// RecordChange increments the state-change counter for port.
func (s *Scorer) RecordChange(port int) {
	s.mu.Lock()
	s.changes[port]++
	s.mu.Unlock()
}

// Compute returns a Weight for the given port.
// Score = base + 0.5*freq + 2.0*changes
func (s *Scorer) Compute(port int) Weight {
	s.mu.RLock()
	defer s.mu.RUnlock()

	score := s.base[port]
	score += 0.5 * float64(s.freq[port])
	score += 2.0 * float64(s.changes[port])

	reason := "base"
	if s.changes[port] > 0 {
		reason = "change-activity"
	} else if s.freq[port] > 0 {
		reason = "open-frequency"
	}
	return Weight{Port: port, Score: score, Reason: reason}
}

// ComputeAll returns weights for every port that has been observed
// or has a base score, sorted by descending score.
func (s *Scorer) ComputeAll() []Weight {
	s.mu.RLock()
	seen := make(map[int]struct{})
	for p := range s.base {
		seen[p] = struct{}{}
	}
	for p := range s.freq {
		seen[p] = struct{}{}
	}
	for p := range s.changes {
		seen[p] = struct{}{}
	}
	s.mu.RUnlock()

	out := make([]Weight, 0, len(seen))
	for p := range seen {
		out = append(out, s.Compute(p))
	}
	sortWeights(out)
	return out
}

func sortWeights(ws []Weight) {
	for i := 1; i < len(ws); i++ {
		for j := i; j > 0 && ws[j].Score > ws[j-1].Score; j-- {
			ws[j], ws[j-1] = ws[j-1], ws[j]
		}
	}
}
