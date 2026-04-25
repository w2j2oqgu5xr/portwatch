// Package portage tracks how long each port has been in its current state
// and exposes helpers to classify ports as "new", "stable", or "stale".
package portage

import (
	"sync"
	"time"
)

// Status classifies a port based on how long it has been open.
type Status string

const (
	StatusNew    Status = "new"    // open for less than NewThreshold
	StatusStable Status = "stable" // open for between New and Stale thresholds
	StatusStale  Status = "stale"  // open for longer than StaleThreshold
)

// DefaultNewThreshold is the duration below which a port is considered new.
const DefaultNewThreshold = 5 * time.Minute

// DefaultStaleThreshold is the duration above which a port is considered stale.
const DefaultStaleThreshold = 24 * time.Hour

// Tracker records the first-seen timestamp for each port and classifies it.
type Tracker struct {
	mu             sync.RWMutex
	firstSeen      map[int]time.Time
	newThreshold   time.Duration
	staleThreshold time.Duration
}

// New creates a Tracker with the supplied thresholds.
// Zero values fall back to package-level defaults.
func New(newThreshold, staleThreshold time.Duration) *Tracker {
	if newThreshold <= 0 {
		newThreshold = DefaultNewThreshold
	}
	if staleThreshold <= 0 {
		staleThreshold = DefaultStaleThreshold
	}
	return &Tracker{
		firstSeen:      make(map[int]time.Time),
		newThreshold:   newThreshold,
		staleThreshold: staleThreshold,
	}
}

// Observe records the first time a port is seen open.
// Subsequent calls for the same port are no-ops.
func (t *Tracker) Observe(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.firstSeen[port]; !ok {
		t.firstSeen[port] = time.Now()
	}
}

// Forget removes a port from the tracker (e.g. when it closes).
func (t *Tracker) Forget(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.firstSeen, port)
}

// Age returns how long a port has been tracked and true, or zero and false if
// the port is not known.
func (t *Tracker) Age(port int) (time.Duration, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	ts, ok := t.firstSeen[port]
	if !ok {
		return 0, false
	}
	return time.Since(ts), true
}

// Classify returns the Status of a port. If the port is not tracked it returns
// StatusNew as a safe default.
func (t *Tracker) Classify(port int) Status {
	age, ok := t.Age(port)
	if !ok {
		return StatusNew
	}
	switch {
	case age >= t.staleThreshold:
		return StatusStale
	case age >= t.newThreshold:
		return StatusStable
	default:
		return StatusNew
	}
}
