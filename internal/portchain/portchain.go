// Package portchain provides a composable middleware chain for processing
// port scan results through a series of transformation and filtering stages
// before they reach a final handler.
//
// A Chain is built once and reused across scan cycles. Each stage receives
// a port number and may pass it downstream, drop it, or enrich the context.
package portchain

import (
	"context"
	"fmt"
	"sync"
)

// Port represents a single TCP port number.
type Port = int

// Handler is the terminal function that receives a port after all stages
// have processed it.
type Handler func(ctx context.Context, port Port)

// Stage is a middleware function. It receives the current port and the next
// handler in the chain. A stage may call next to continue processing, or
// return early to drop the port.
type Stage func(ctx context.Context, port Port, next Handler)

// Chain is an ordered sequence of stages that terminates at a Handler.
type Chain struct {
	mu      sync.RWMutex
	stages  []Stage
	handler Handler
}

// New creates a Chain with the given terminal handler. Stages can be appended
// via Use before the chain is executed.
func New(handler Handler) *Chain {
	if handler == nil {
		handler = func(_ context.Context, _ Port) {}
	}
	return &Chain{handler: handler}
}

// Use appends one or more stages to the chain. Stages are applied in the
// order they are added: the first stage added is the outermost wrapper.
// Use is safe to call before the chain processes any ports, but must not
// be called concurrently with Process.
func (c *Chain) Use(stages ...Stage) *Chain {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stages = append(c.stages, stages...)
	return c
}

// Process runs port through every stage in order, then calls the terminal
// handler if no stage drops the port. The provided context is forwarded
// unchanged to each stage and the final handler.
func (c *Chain) Process(ctx context.Context, port Port) {
	c.mu.RLock()
	stages := make([]Stage, len(c.stages))
	copy(stages, c.stages)
	handler := c.handler
	c.mu.RUnlock()

	build(stages, handler)(ctx, port)
}

// build recursively wraps handler with each stage, right to left, so that
// the first stage in the slice is the outermost call.
func build(stages []Stage, final Handler) Handler {
	if len(stages) == 0 {
		return final
	}
	next := build(stages[1:], final)
	current := stages[0]
	return func(ctx context.Context, port Port) {
		current(ctx, port, next)
	}
}

// LogStage returns a Stage that prints a debug line for each port that passes
// through. It is intended for development and testing only.
func LogStage(logf func(format string, args ...any)) Stage {
	if logf == nil {
		logf = func(format string, args ...any) { fmt.Printf(format+"\n", args...) }
	}
	return func(ctx context.Context, port Port, next Handler) {
		logf("portchain: processing port %d", port)
		next(ctx, port)
	}
}

// AllowStage returns a Stage that only forwards ports for which allow returns
// true. Ports that do not pass the predicate are silently dropped.
func AllowStage(allow func(port Port) bool) Stage {
	return func(ctx context.Context, port Port, next Handler) {
		if allow(port) {
			next(ctx, port)
		}
	}
}

// TransformStage returns a Stage that maps a port to a different port number
// before forwarding it. This is useful for normalisation in tests.
func TransformStage(fn func(port Port) Port) Stage {
	return func(ctx context.Context, port Port, next Handler) {
		next(ctx, fn(port))
	}
}
