// Package portttl tracks how long ports have been in their current state.
package portttl

import (
	"sync"
	"time"
)

// Entry holds the state transition time for a single port.
type Entry struct {
	Port      int
	Since     time.Time
	OpenState bool
}

// Tracker records when each port last changed state.
type Tracker struct {
	mu      sync.RWMutex
	entries map[int]Entry
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{entries: make(map[int]Entry)}
}

// Record marks a port as having changed to the given state at the current time.
func (t *Tracker) Record(port int, open bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[port]
	if ok && e.OpenState == open {
		return // state unchanged, keep original timestamp
	}
	t.entries[port] = Entry{Port: port, Since: time.Now(), OpenState: open}
}

// Age returns how long the port has been in its current state.
// Returns zero and false if the port has never been recorded.
func (t *Tracker) Age(port int) (time.Duration, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[port]
	if !ok {
		return 0, false
	}
	return time.Since(e.Since), true
}

// Get returns the Entry for a port, and whether it exists.
func (t *Tracker) Get(port int) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[port]
	return e, ok
}

// Delete removes the tracking entry for a port.
func (t *Tracker) Delete(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// All returns a snapshot of all tracked entries.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}
