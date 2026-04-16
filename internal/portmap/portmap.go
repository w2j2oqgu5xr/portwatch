// Package portmap provides a thread-safe mapping of ports to metadata
// such as process name, PID, and protocol for enriching scan results.
package portmap

import (
	"fmt"
	"sync"
)

// Entry holds metadata associated with an open port.
type Entry struct {
	Port     int
	Protocol string
	PID      int
	Process  string
}

// String returns a human-readable representation of the entry.
func (e Entry) String() string {
	if e.Process != "" {
		return fmt.Sprintf("%d/%s (pid=%d, %s)", e.Port, e.Protocol, e.PID, e.Process)
	}
	return fmt.Sprintf("%d/%s", e.Port, e.Protocol)
}

// Map is a thread-safe store of port entries.
type Map struct {
	mu      sync.RWMutex
	entries map[int]Entry
}

// New returns an empty Map.
func New() *Map {
	return &Map{entries: make(map[int]Entry)}
}

// Set stores an entry for the given port.
func (m *Map) Set(port int, e Entry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[port] = e
}

// Get retrieves the entry for a port. Returns false if not found.
func (m *Map) Get(port int) (Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[port]
	return e, ok
}

// Delete removes the entry for a port.
func (m *Map) Delete(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, port)
}

// All returns a snapshot of all entries.
func (m *Map) All() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Entry, 0, len(m.entries))
	for _, e := range m.entries {
		out = append(out, e)
	}
	return out
}

// Len returns the number of entries.
func (m *Map) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.entries)
}
