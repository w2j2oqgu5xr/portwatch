// Package throttle provides per-port alert suppression to prevent
// notification floods when a port repeatedly changes state.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks the last alert time for each port and suppresses
// duplicate alerts that arrive within the cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSeen map[int]time.Time
}

// New creates a Throttle with the given cooldown duration.
// If cooldown is zero or negative, a default of 30 seconds is used.
func New(cooldown time.Duration) *Throttle {
	if cooldown <= 0 {
		cooldown = 30 * time.Second
	}
	return &Throttle{
		cooldown: cooldown,
		lastSeen: make(map[int]time.Time),
	}
}

// Allow returns true if an alert for the given port should be emitted.
// It returns false when the port was alerted within the cooldown window.
func (t *Throttle) Allow(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if last, ok := t.lastSeen[port]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.lastSeen[port] = now
	return true
}

// Reset clears the suppression record for a specific port, allowing
// the next alert for that port to pass through immediately.
func (t *Throttle) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSeen, port)
}

// ResetAll clears suppression state for every port.
func (t *Throttle) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastSeen = make(map[int]time.Time)
}

// Cooldown returns the configured cooldown duration.
func (t *Throttle) Cooldown() time.Duration {
	return t.cooldown
}
