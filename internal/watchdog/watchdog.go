// Package watchdog provides a self-healing mechanism that restarts
// the scan loop if it stalls or panics unexpectedly.
package watchdog

import (
	"context"
	"log"
	"time"
)

// DefaultTimeout is the maximum time allowed between heartbeats.
const DefaultTimeout = 30 * time.Second

// Watchdog monitors a heartbeat channel and calls the restart function
// if no beat is received within the timeout window.
type Watchdog struct {
	timeout time.Duration
	heartbeat <-chan struct{}
	onStall func()
}

// New creates a Watchdog. heartbeat should be sent to regularly by the
// monitored goroutine. onStall is invoked when the timeout elapses.
func New(heartbeat <-chan struct{}, timeout time.Duration, onStall func()) *Watchdog {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &Watchdog{
		timeout:  timeout,
		heartbeat: heartbeat,
		onStall:  onStall,
	}
}

// Run starts the watchdog loop. It blocks until ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context) {
	timer := time.NewTimer(w.timeout)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.heartbeat:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(w.timeout)
		case <-timer.C:
			log.Println("[watchdog] stall detected — invoking recovery")
			w.onStall()
			timer.Reset(w.timeout)
		}
	}
}
