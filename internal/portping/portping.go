// Package portping measures round-trip latency to open TCP ports.
package portping

import (
	"context"
	"fmt"
	"net"
	"time"
)

// DefaultTimeout is used when no timeout is specified.
const DefaultTimeout = 2 * time.Second

// Result holds the outcome of a single ping attempt.
type Result struct {
	Host    string
	Port    int
	Latency time.Duration
	Err     error
}

// Alive reports whether the ping succeeded.
func (r Result) Alive() bool { return r.Err == nil }

// String returns a human-readable summary.
func (r Result) String() string {
	if !r.Alive() {
		return fmt.Sprintf("%s:%d unreachable (%v)", r.Host, r.Port, r.Err)
	}
	return fmt.Sprintf("%s:%d alive latency=%s", r.Host, r.Port, r.Latency.Round(time.Microsecond))
}

// Pinger sends TCP pings to a host:port.
type Pinger struct {
	timeout time.Duration
}

// New creates a Pinger with the given timeout. If timeout <= 0 DefaultTimeout is used.
func New(timeout time.Duration) *Pinger {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &Pinger{timeout: timeout}
}

// Ping dials host:port once and returns the result.
func (p *Pinger) Ping(ctx context.Context, host string, port int) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	d := net.Dialer{Timeout: p.timeout}
	conn, err := d.DialContext(ctx, "tcp", addr)
	latency := time.Since(start)
	if err != nil {
		return Result{Host: host, Port: port, Err: err}
	}
	_ = conn.Close()
	return Result{Host: host, Port: port, Latency: latency}
}

// PingAll pings each port in ports and returns all results.
func (p *Pinger) PingAll(ctx context.Context, host string, ports []int) []Result {
	results := make([]Result, 0, len(ports))
	for _, port := range ports {
		results = append(results, p.Ping(ctx, host, port))
	}
	return results
}
