// Package portscan provides a concurrent port scanner with timeout and result streaming.
package portscan

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// Result holds the outcome of scanning a single port.
type Result struct {
	Port   int
	Open   bool
	Banner string
	Err    error
}

// Options configures the concurrent scanner.
type Options struct {
	Concurrency int
	Timeout     time.Duration
	GrabBanner  bool
}

func defaultOptions(o Options) Options {
	if o.Concurrency <= 0 {
		o.Concurrency = 100
	}
	if o.Timeout <= 0 {
		o.Timeout = 2 * time.Second
	}
	return o
}

// Scan concurrently scans the given ports on host and streams results to the returned channel.
func Scan(ctx context.Context, host string, ports []int, opts Options) <-chan Result {
	opts = defaultOptions(opts)
	out := make(chan Result, len(ports))

	go func() {
		defer close(out)
		sem := make(chan struct{}, opts.Concurrency)
		var wg sync.WaitGroup

		for _, p := range ports {
			select {
			case <-ctx.Done():
				return
			case sem <- struct{}{}:
			}

			wg.Add(1)
			go func(port int) {
				defer wg.Done()
				defer func() { <-sem }()
				out <- probe(ctx, host, port, opts)
			}(p)
		}
		wg.Wait()
	}()

	return out
}

func probe(ctx context.Context, host string, port int, opts Options) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	d := net.Dialer{Timeout: opts.Timeout}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return Result{Port: port, Open: false, Err: err}
	}
	defer conn.Close()

	r := Result{Port: port, Open: true}
	if opts.GrabBanner {
		_ = conn.SetReadDeadline(time.Now().Add(opts.Timeout))
		buf := make([]byte, 256)
		n, _ := conn.Read(buf)
		if n > 0 {
			r.Banner = string(buf[:n])
		}
	}
	return r
}
