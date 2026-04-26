// Package portlock provides a mechanism for locking (pinning) specific ports
// so that alerts are suppressed for those ports during maintenance windows
// or known-good states.
package portlock

import (
	"fmt"
	"sync"
	"time"
)

// Lock represents a pinned port entry.
type Lock struct {
	Port      int
	Reason    string
	LockedAt  time.Time
	ExpiresAt time.Time // zero means no expiry
}

// Expired reports whether the lock has passed its expiry time.
func (l Lock) Expired(now time.Time) bool {
	return !l.ExpiresAt.IsZero() && now.After(l.ExpiresAt)
}

// String returns a human-readable description of the lock.
func (l Lock) String() string {
	if l.ExpiresAt.IsZero() {
		return fmt.Sprintf("port %d locked: %s (no expiry)", l.Port, l.Reason)
	}
	return fmt.Sprintf("port %d locked: %s (expires %s)", l.Port, l.Reason, l.ExpiresAt.Format(time.RFC3339))
}

// Registry holds the set of currently locked ports.
type Registry struct {
	mu    sync.RWMutex
	locks map[int]Lock
	now   func() time.Time
}

// New creates a new Registry.
func New() *Registry {
	return &Registry{
		locks: make(map[int]Lock),
		now:   time.Now,
	}
}

// Lock pins a port with an optional TTL (zero means no expiry).
func (r *Registry) Lock(port int, reason string, ttl time.Duration) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portlock: invalid port %d", port)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.now()
	var exp time.Time
	if ttl > 0 {
		exp = now.Add(ttl)
	}
	r.locks[port] = Lock{Port: port, Reason: reason, LockedAt: now, ExpiresAt: exp}
	return nil
}

// Unlock removes the lock for a port.
func (r *Registry) Unlock(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.locks, port)
}

// IsLocked reports whether a port is currently locked (and not expired).
func (r *Registry) IsLocked(port int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.locks[port]
	if !ok {
		return false
	}
	return !l.Expired(r.now())
}

// All returns a snapshot of all active (non-expired) locks.
func (r *Registry) All() []Lock {
	r.mu.RLock()
	defer r.mu.RUnlock()
	now := r.now()
	out := make([]Lock, 0, len(r.locks))
	for _, l := range r.locks {
		if !l.Expired(now) {
			out = append(out, l)
		}
	}
	return out
}

// Purge removes all expired locks.
func (r *Registry) Purge() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.now()
	removed := 0
	for port, l := range r.locks {
		if l.Expired(now) {
			delete(r.locks, port)
			removed++
		}
	}
	return removed
}
