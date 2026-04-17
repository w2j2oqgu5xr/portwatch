// Package portstate tracks the current open/closed state of monitored ports
// and provides change detection between successive scans.
package portstate

import (
	"sync"
	"time"
)

// State represents the observed state of a single port.
type State struct {
	Port      int
	Open      bool
	LastSeen  time.Time
	FirstSeen time.Time
}

// Change describes a transition for a port.
type Change struct {
	Port   int
	Opened bool // true = newly opened, false = newly closed
	At     time.Time
}

// Tracker holds the last known state for a set of ports.
type Tracker struct {
	mu     sync.RWMutex
	states map[int]State
}

// New returns an empty Tracker.
func New() *Tracker {
	return &Tracker{states: make(map[int]State)}
}

// Update applies a new scan result and returns any changes detected.
func (t *Tracker) Update(openPorts []int) []Change {
	now := time.Now()
	newSet := make(map[int]struct{}, len(openPorts))
	for _, p := range openPorts {
		newSet[p] = struct{}{}
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	var changes []Change

	// Detect newly opened ports.
	for _, p := range openPorts {
		if s, exists := t.states[p]; !exists || !s.Open {
			first := now
			if exists {
				first = s.FirstSeen
			}
			t.states[p] = State{Port: p, Open: true, LastSeen: now, FirstSeen: first}
			changes = append(changes, Change{Port: p, Opened: true, At: now})
		} else {
			s.LastSeen = now
			t.states[p] = s
		}
	}

	// Detect newly closed ports.
	for p, s := range t.states {
		if _, open := newSet[p]; !open && s.Open {
			s.Open = false
			s.LastSeen = now
			t.states[p] = s
			changes = append(changes, Change{Port: p, Opened: false, At: now})
		}
	}

	return changes
}

// Snapshot returns a copy of all tracked states.
func (t *Tracker) Snapshot() []State {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]State, 0, len(t.states))
	for _, s := range t.states {
		out = append(out, s)
	}
	return out
}

// OpenPorts returns the list of currently open port numbers.
func (t *Tracker) OpenPorts() []int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var ports []int
	for p, s := range t.states {
		if s.Open {
			ports = append(ports, p)
		}
	}
	return ports
}
