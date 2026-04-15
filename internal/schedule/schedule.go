// Package schedule provides interval-based ticker utilities for
// periodically triggering port scans within the portwatch monitor loop.
package schedule

import (
	"context"
	"errors"
	"time"
)

// MinInterval is the smallest allowed scan interval to prevent
// accidental resource exhaustion.
const MinInterval = 500 * time.Millisecond

// Ticker emits ticks on a fixed interval until the context is cancelled.
type Ticker struct {
	C       <-chan time.Time
	interval time.Duration
	stop    func()
}

// New creates a Ticker that fires every interval. Returns an error if
// interval is below MinInterval.
func New(ctx context.Context, interval time.Duration) (*Ticker, error) {
	if interval < MinInterval {
		return nil, errors.New("schedule: interval below minimum allowed value")
	}

	t := time.NewTicker(interval)
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan time.Time, 1)

	go func() {
		defer close(ch)
		for {
			select {
			case tick := <-t.C:
				select {
				case ch <- tick:
				default:
					// drop tick if consumer is behind
				}
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}()

	return &Ticker{
		C:        ch,
		interval: interval,
		stop:     cancel,
	}, nil
}

// Interval returns the configured tick interval.
func (t *Ticker) Interval() time.Duration {
	return t.interval
}

// Stop cancels the underlying context and halts the ticker.
func (t *Ticker) Stop() {
	t.stop()
}
