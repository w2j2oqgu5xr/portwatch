// Package portevict tracks ports that have been evicted (closed after being
// open for a sustained period) and provides a simple eviction log.
package portevict

import (
	"sync"
	"time"
)

// Entry records a single port eviction event.
type Entry struct {
	Port      int
	Host      string
	OpenedAt  time.Time
	EvictedAt time.Time
	Duration  time.Duration
}

// Log holds an in-memory eviction history.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	max     int
}

// New creates a Log that retains at most max entries.
// If max is zero it defaults to 256.
func New(max int) *Log {
	if max <= 0 {
		max = 256
	}
	return &Log{max: max}
}

// Record appends an eviction entry derived from the port, host, and the time
// the port was first observed open.
func (l *Log) Record(port int, host string, openedAt time.Time) Entry {
	now := time.Now()
	e := Entry{
		Port:      port,
		Host:      host,
		OpenedAt:  openedAt,
		EvictedAt: now,
		Duration:  now.Sub(openedAt),
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, e)
	if len(l.entries) > l.max {
		l.entries = l.entries[len(l.entries)-l.max:]
	}
	return e
}

// All returns a shallow copy of all recorded entries.
func (l *Log) All() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Len returns the number of entries currently held.
func (l *Log) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}

// Clear removes all entries.
func (l *Log) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = l.entries[:0]
}
