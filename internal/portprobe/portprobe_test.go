package portprobe_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portprobe"
)

func startTCP(t *testing.T) (port int, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	p, _ := strconv.Atoi(portStr)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return p, func() { ln.Close() }
}

func TestProbe_OpenPort(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()

	r := portprobe.Probe(context.Background(), "127.0.0.1", port, portprobe.Options{})
	if !r.Open {
		t.Fatalf("expected port %d to be open, err: %v", port, r.Err)
	}
	if r.Latency <= 0 {
		t.Errorf("expected positive latency, got %s", r.Latency)
	}
}

func TestProbe_ClosedPort(t *testing.T) {
	// bind then immediately close to get a free port number
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	ln.Close()

	r := portprobe.Probe(context.Background(), "127.0.0.1", port, portprobe.Options{Timeout: 200 * time.Millisecond})
	if r.Open {
		t.Fatalf("expected port %d to be closed", port)
	}
	if r.Err == nil {
		t.Error("expected non-nil error for closed port")
	}
}

func TestProbe_DefaultTimeout(t *testing.T) {
	// defaults should be applied without panic
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	ln.Close()

	// zero-value options should not panic
	r := portprobe.Probe(context.Background(), "127.0.0.1", port, portprobe.Options{})
	if r.Open {
		t.Fatal("expected closed")
	}
}

func TestProbeAll_ReturnsAllResults(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()

	ports := []int{port, port}
	results := portprobe.ProbeAll(context.Background(), "127.0.0.1", ports, portprobe.Options{})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Open {
			t.Errorf("expected port %d open", r.Port)
		}
	}
}

func TestResult_String(t *testing.T) {
	r := portprobe.Result{Host: "127.0.0.1", Port: 80, Open: true, Latency: 500 * time.Microsecond}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
