// Package portecho provides a round-trip latency prober that connects to a
// port, writes a small payload, and measures the time until the connection is
// accepted. It is used to distinguish slow-but-open ports from truly closed
// ones and to surface response-time regressions over successive scans.
package portecho

import (
	"context"
	"fmt"
	"net"
	"time"
)

// DefaultTimeout is used when Config.Timeout is zero.
const DefaultTimeout = 2 * time.Second

// Result holds the outcome of a single echo probe.
type Result struct {
	Host    string
	Port    int
	Open    bool
	Latency time.Duration
	Err     error
}

// String returns a human-readable summary of the result.
func (r Result) String() string {
	if !r.Open {
		return fmt.Sprintf("%s:%d closed (%v)", r.Host, r.Port, r.Err)
	}
	return fmt.Sprintf("%s:%d open latency=%s", r.Host, r.Port, r.Latency.Round(time.Microsecond))
}

// Prober sends echo probes to one or more ports on a host.
type Prober struct {
	host    string
	timeout time.Duration
}

// New returns a Prober targeting host. If timeout is zero, DefaultTimeout is used.
func New(host string, timeout time.Duration) *Prober {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &Prober{host: host, timeout: timeout}
}

// Probe connects to a single port and returns a Result.
func (p *Prober) Probe(ctx context.Context, port int) Result {
	addr := fmt.Sprintf("%s:%d", p.host, port)
	start := time.Now()

	dialer := &net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	latency := time.Since(start)

	if err != nil {
		return Result{Host: p.host, Port: port, Open: false, Err: err}
	}
	_ = conn.Close()
	return Result{Host: p.host, Port: port, Open: true, Latency: latency}
}

// ProbeAll probes each port in ports concurrently and returns all results.
func (p *Prober) ProbeAll(ctx context.Context, ports []int) []Result {
	type indexed struct {
		i int
		r Result
	}
	ch := make(chan indexed, len(ports))
	for i, port := range ports {
		go func(i, port int) {
			ch <- indexed{i, p.Probe(ctx, port)}
		}(i, port)
	}
	results := make([]Result, len(ports))
	for range ports {
		ix := <-ch
		results[ix.i] = ix.r
	}
	return results
}
