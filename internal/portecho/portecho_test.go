package portecho_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portecho"
)

func startTCP(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

func TestProbe_OpenPort(t *testing.T) {
	port := startTCP(t)
	p := portecho.New("127.0.0.1", 0)
	r := p.Probe(context.Background(), port)
	if !r.Open {
		t.Fatalf("expected open, got closed: %v", r.Err)
	}
	if r.Latency <= 0 {
		t.Errorf("expected positive latency, got %s", r.Latency)
	}
}

func TestProbe_ClosedPort(t *testing.T) {
	p := portecho.New("127.0.0.1", 200*time.Millisecond)
	r := p.Probe(context.Background(), 1)
	if r.Open {
		t.Fatal("expected closed port")
	}
	if r.Err == nil {
		t.Error("expected non-nil error for closed port")
	}
}

func TestProbe_DefaultTimeout(t *testing.T) {
	p := portecho.New("127.0.0.1", 0)
	// We can only verify it doesn't panic; timeout is internal.
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
}

func TestProbeAll_ReturnsAllResults(t *testing.T) {
	port1 := startTCP(t)
	port2 := startTCP(t)

	p := portecho.New("127.0.0.1", 0)
	results := p.ProbeAll(context.Background(), []int{port1, port2, 1})

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	openCount := 0
	for _, r := range results {
		if r.Open {
			openCount++
		}
	}
	if openCount != 2 {
		t.Errorf("expected 2 open, got %d", openCount)
	}
}

func TestResult_String_Open(t *testing.T) {
	r := portecho.Result{Host: "127.0.0.1", Port: 8080, Open: true, Latency: 500 * time.Microsecond}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestResult_String_Closed(t *testing.T) {
	r := portecho.Result{Host: "127.0.0.1", Port: 9999, Open: false}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
