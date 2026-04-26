// Package portquota enforces a maximum number of simultaneously open ports
// and emits a violation when the observed count exceeds the configured limit.
package portquota

import (
	"fmt"
	"sync"
)

// Violation describes a quota breach at a point in time.
type Violation struct {
	Limit int
	Actual int
}

func (v Violation) String() string {
	return fmt.Sprintf("quota exceeded: %d open ports (limit %d)", v.Actual, v.Limit)
}

// Enforcer tracks open-port counts and reports quota violations.
type Enforcer struct {
	mu    sync.Mutex
	limit int
	count int
}

// New creates an Enforcer with the given maximum open-port limit.
// If limit is less than 1 it is clamped to 1.
func New(limit int) *Enforcer {
	if limit < 1 {
		limit = 1
	}
	return &Enforcer{limit: limit}
}

// Set updates the current open-port count and returns a non-nil Violation
// when the count exceeds the configured limit.
func (e *Enforcer) Set(count int) *Violation {
	e.mu.Lock()
	defer e.mu.Unlock()

	if count < 0 {
		count = 0
	}
	e.count = count

	if e.count > e.limit {
		v := &Violation{Limit: e.limit, Actual: e.count}
		return v
	}
	return nil
}

// Count returns the most recently recorded open-port count.
func (e *Enforcer) Count() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.count
}

// Limit returns the configured maximum.
func (e *Enforcer) Limit() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.limit
}

// SetLimit updates the quota limit at runtime.
// If limit is less than 1 it is clamped to 1.
func (e *Enforcer) SetLimit(limit int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if limit < 1 {
		limit = 1
	}
	e.limit = limit
}
