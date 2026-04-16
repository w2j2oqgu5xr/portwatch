// Package retry provides simple retry logic with backoff for transient errors.
package retry

import (
	"context"
	"time"
)

// Policy defines retry behaviour.
type Policy struct {
	MaxAttempts int
	InitialDelay time.Duration
	Multiplier float64
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		Multiplier:   2.0,
	}
}

// Do executes fn up to p.MaxAttempts times, backing off between attempts.
// It returns the last error if all attempts fail.
func Do(ctx context.Context, p Policy, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	if p.Multiplier <= 0 {
		p.Multiplier = 1
	}

	delay := p.InitialDelay
	var err error

	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err = fn()
		if err == nil {
			return nil
		}
		if attempt < p.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * p.Multiplier)
		}
	}
	return err
}
