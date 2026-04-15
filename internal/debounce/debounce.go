// Package debounce provides a mechanism to suppress rapid repeated events
// for the same port, emitting only after a quiet period has elapsed.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays forwarding of events until no new event for the same key
// has arrived within the configured wait duration.
type Debouncer struct {
	wait    time.Duration
	mu      sync.Mutex
	timers  map[string]*time.Timer
}

// New creates a Debouncer that waits for the given duration of silence
// before invoking the callback. wait must be positive; if not, it defaults
// to 100 ms.
func New(wait time.Duration) *Debouncer {
	if wait <= 0 {
		wait = 100 * time.Millisecond
	}
	return &Debouncer{
		wait:   wait,
		timers: make(map[string]*time.Timer),
	}
}

// Submit schedules fn to be called after the debounce window for key expires.
// If Submit is called again for the same key before the window elapses, the
// previous scheduled call is cancelled and the window resets.
func (d *Debouncer) Submit(key string, fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
	}

	d.timers[key] = time.AfterFunc(d.wait, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		fn()
	})
}

// Cancel discards any pending callback for key without invoking it.
func (d *Debouncer) Cancel(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
		delete(d.timers, key)
	}
}

// Pending returns the number of keys that currently have a scheduled callback.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.timers)
}
