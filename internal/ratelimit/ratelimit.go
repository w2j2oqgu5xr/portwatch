// Package ratelimit provides a simple token-bucket rate limiter used to
// throttle alert notifications so that a burst of port changes does not
// flood the operator with messages.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a thread-safe token-bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens added per second
	lastTick time.Time
	now      func() time.Time // injectable for testing
}

// New returns a Limiter that allows up to burst events and refills at the
// given rate (events per second).
func New(rate float64, burst int) *Limiter {
	if rate <= 0 {
		rate = 1
	}
	if burst <= 0 {
		burst = 1
	}
	return &Limiter{
		tokens:   float64(burst),
		max:      float64(burst),
		rate:     rate,
		lastTick: time.Now(),
		now:      time.Now,
	}
}

// Allow reports whether one event may proceed. It consumes one token when
// returning true.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Reset restores the bucket to its maximum capacity.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tokens = l.max
	l.lastTick = l.now()
}
