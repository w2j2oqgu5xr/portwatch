// Package baseline stores and compares a trusted set of open ports
// to distinguish expected from unexpected port changes.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// ErrNoBaseline is returned when no baseline file exists.
var ErrNoBaseline = errors.New("baseline: no baseline found")

// Baseline holds a trusted snapshot of open ports.
type Baseline struct {
	Ports     []int     `json:"ports"`
	RecordedAt time.Time `json:"recorded_at"`
}

// Save writes the baseline to the given path as JSON.
func Save(path string, ports []int) error {
	b := Baseline{
		Ports:      ports,
		RecordedAt: time.Now().UTC(),
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(b)
}

// Load reads a baseline from the given path.
func Load(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoBaseline
		}
		return nil, err
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Unexpected returns ports that are open but not in the baseline.
func (b *Baseline) Unexpected(open []int) []int {
	trusted := make(map[int]struct{}, len(b.Ports))
	for _, p := range b.Ports {
		trusted[p] = struct{}{}
	}
	var out []int
	for _, p := range open {
		if _, ok := trusted[p]; !ok {
			out = append(out, p)
		}
	}
	return out
}

// Missing returns ports in the baseline that are no longer open.
func (b *Baseline) Missing(open []int) []int {
	current := make(map[int]struct{}, len(open))
	for _, p := range open {
		current[p] = struct{}{}
	}
	var out []int
	for _, p := range b.Ports {
		if _, ok := current[p]; !ok {
			out = append(out, p)
		}
	}
	return out
}
