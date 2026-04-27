package portrank

import "sort"

// Score represents the computed risk score for a port.
type Score struct {
	Port  int
	Value float64
	Rank  int
}

// Ranker scores and ranks ports by risk based on frequency of state
// changes and classification weight.
type Ranker struct {
	weights   map[int]float64
	changes   map[int]int
	baseScore float64
}

// New returns a Ranker with the given base score applied to unknown ports.
func New(baseScore float64) *Ranker {
	if baseScore <= 0 {
		baseScore = 1.0
	}
	return &Ranker{
		weights:   make(map[int]float64),
		changes:   make(map[int]int),
		baseScore: baseScore,
	}
}

// SetWeight assigns a classification weight to a port. Higher weights
// increase the port's risk score.
func (r *Ranker) SetWeight(port int, weight float64) {
	if weight < 0 {
		weight = 0
	}
	r.weights[port] = weight
}

// RecordChange increments the state-change counter for a port.
func (r *Ranker) RecordChange(port int) {
	r.changes[port]++
}

// Rank computes risk scores for the supplied ports and returns them
// sorted from highest to lowest score.
func (r *Ranker) Rank(ports []int) []Score {
	scores := make([]Score, 0, len(ports))
	for _, p := range ports {
		w, ok := r.weights[p]
		if !ok {
			w = r.baseScore
		}
		v := w * float64(1+r.changes[p])
		scores = append(scores, Score{Port: p, Value: v})
	}
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].Value != scores[j].Value {
			return scores[i].Value > scores[j].Value
		}
		return scores[i].Port < scores[j].Port
	})
	for i := range scores {
		scores[i].Rank = i + 1
	}
	return scores
}

// Reset clears all recorded state changes while preserving weights.
func (r *Ranker) Reset() {
	r.changes = make(map[int]int)
}
