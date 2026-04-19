// Package portexpiry tracks ports that have been closed and emits expiry
// events once a configurable grace period has elapsed without reopening.
package portexpiry

import (
	"sync"
	"time"
)

// Entry records when a port was marked closed.
type Entry struct {
	Port      int
	ClosedAt  time.Time
	ExpiredAt time.Time
}

// Tracker monitors closed ports and reports those past their grace period.
type Tracker struct {
	mu     sync.Mutex
	grace  time.Duration
	closed map[int]time.Time
}

// New creates a Tracker with the given grace duration.
// If grace is zero or negative, it defaults to 5 minutes.
func New(grace time.Duration) *Tracker {
	if grace <= 0 {
		grace = 5 * time.Minute
	}
	return &Tracker{
		grace:  grace,
		closed: make(map[int]time.Time),
	}
}

// MarkClosed records the time a port transitioned to closed.
// If the port is already tracked, the timestamp is not updated.
func (t *Tracker) MarkClosed(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.closed[port]; !exists {
		t.closed[port] = time.Now()
	}
}

// MarkOpen removes a port from expiry tracking (it reopened).
func (t *Tracker) MarkOpen(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.closed, port)
}

// Expired returns all ports whose grace period has elapsed.
// Expired ports are removed from tracking.
func (t *Tracker) Expired() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	var out []Entry
	for port, closedAt := range t.closed {
		if now.Sub(closedAt) >= t.grace {
			out = append(out, Entry{
				Port:      port,
				ClosedAt:  closedAt,
				ExpiredAt: now,
			})
			delete(t.closed, port)
		}
	}
	return out
}

// Pending returns the number of ports currently in the grace period.
func (t *Tracker) Pending() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.closed)
}
