// Package snapshot provides functionality for persisting and loading
// port scan state to disk, enabling portwatch to detect changes across
// process restarts.
package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// ErrNoSnapshot is returned when no snapshot file exists at the given path.
var ErrNoSnapshot = errors.New("snapshot: no snapshot found")

// Snapshot holds a recorded set of open ports at a point in time.
type Snapshot struct {
	Ports     []int     `json:"ports"`
	RecordedAt time.Time `json:"recorded_at"`
}

// Save writes the snapshot to the given file path as JSON.
func Save(path string, ports []int) error {
	s := Snapshot{
		Ports:      ports,
		RecordedAt: time.Now().UTC(),
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Load reads a snapshot from the given file path.
// Returns ErrNoSnapshot if the file does not exist.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoSnapshot
		}
		return nil, err
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Diff compares a previous snapshot's ports against a current set of open
// ports and returns the ports that were opened and ports that were closed.
func Diff(prev []int, current []int) (opened []int, closed []int) {
	prevSet := toSet(prev)
	currSet := toSet(current)

	for p := range currSet {
		if !prevSet[p] {
			opened = append(opened, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			closed = append(closed, p)
		}
	}
	return opened, closed
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
