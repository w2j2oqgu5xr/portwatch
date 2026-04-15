// Package throttle implements per-port alert cooldown logic for portwatch.
//
// When a port flaps or a scan loop fires rapidly, the same open/close event
// can trigger many notifications in quick succession. Throttle suppresses
// duplicate alerts for a configurable cooldown window, ensuring operators
// receive actionable signals rather than notification floods.
//
// Basic usage:
//
//	th := throttle.New(30 * time.Second)
//	if th.Allow(port) {
//		notifier.Notify(event)
//	}
//
// Throttle is safe for concurrent use by multiple goroutines.
package throttle
