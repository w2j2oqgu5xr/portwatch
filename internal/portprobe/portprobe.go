// Package portprobe provides active probing of individual ports with
// configurable timeout and retry semantics, returning a structured Result.
package portprobe

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single port probe.
type Result struct {
	Host    string
	Port    int
	Open    bool
	Latency time.Duration
	Err     error
}

// String returns a human-readable representation of the result.
func (r Result) String() string {
	state := "closed"
	if r.Open {
		state = "open"
	}
	return fmt.Sprintf("%s:%d %s (%s)", r.Host, r.Port, state, r.Latency.Round(time.Microsecond))
}

// Options configures probe behaviour.
type Options struct {
	Timeout time.Duration
	Retries int
}

func (o *Options) applyDefaults() {
	if o.Timeout <= 0 {
		o.Timeout = 2 * time.Second
	}
	if o.Retries < 0 {
		o.Retries = 0
	}
}

// Probe checks whether a single TCP port is open on the given host.
// It retries up to Options.Retries additional times on failure.
func Probe(ctx context.Context, host string, port int, opts Options) Result {
	opts.applyDefaults()

	addr := fmt.Sprintf("%s:%d", host, port)
	atttempts := 1 + opts.Retries

	var last error
	for i := 0; i < atttempts; i++ {
		start := time.Now()
		conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
		latency := time.Since(start)
		if err == nil {
			conn.Close()
			return Result{Host: host, Port: port, Open: true, Latency: latency}
		}
		last = err
		// respect context cancellation between retries
		if ctx.Err() != nil {
			break
		}
	}
	return Result{Host: host, Port: port, Open: false, Err: last}
}

// ProbeAll probes each port in the provided slice concurrently and returns
// results in the same order as the input.
func ProbeAll(ctx context.Context, host string, ports []int, opts Options) []Result {
	results := make([]Result, len(ports))
	type indexed struct {
		i int
		r Result
	}
	ch := make(chan indexed, len(ports))
	for i, p := range ports {
		go func(idx, port int) {
			ch <- indexed{idx, Probe(ctx, host, port, opts)}
		}(i, p)
	}
	for range ports {
		item := <-ch
		results[item.i] = item.r
	}
	return results
}
