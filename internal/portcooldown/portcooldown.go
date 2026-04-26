// Package portcooldown enforces a minimum quiet period before a port
// transition is considered stable and forwarded downstream.
//
// A port that flaps open→closed→open within the cooldown window is
// suppressed until it settles, reducing alert noise caused by transient
// connectivity blips.
package portcooldown

import (
	"sync"
	"time"
)

// DefaultCooldown is the quiet period applied when none is specified.
const DefaultCooldown = 5 * time.Second

// entry tracks the last state change for a single port.
type entry struct {
	state     string
	changedAt time.Time
}

// Cooldown suppresses repeated state changes for a port within a
// configurable quiet period.
type Cooldown struct {
	mu       sync.Mutex
	period   time.Duration
	entries  map[int]*entry
	nowFunc  func() time.Time
}

// New returns a Cooldown with the given quiet period.
// If period is zero or negative, DefaultCooldown is used.
func New(period time.Duration) *Cooldown {
	if period <= 0 {
		period = DefaultCooldown
	}
	return &Cooldown{
		period:  period,
		entries: make(map[int]*entry),
		nowFunc: time.Now,
	}
}

// Allow reports whether the state change for port should be forwarded.
// The first transition for a port is always allowed. Subsequent
// transitions are suppressed until the cooldown period has elapsed
// since the previous allowed transition.
func (c *Cooldown) Allow(port int, state string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.nowFunc()
	e, ok := c.entries[port]
	if !ok {
		c.entries[port] = &entry{state: state, changedAt: now}
		return true
	}
	if e.state == state {
		return false
	}
	if now.Sub(e.changedAt) < c.period {
		return false
	}
	e.state = state
	e.changedAt = now
	return true
}

// Reset removes the cooldown record for port, allowing its next
// transition to pass through unconditionally.
func (c *Cooldown) Reset(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, port)
}

// Len returns the number of ports currently tracked.
func (c *Cooldown) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
