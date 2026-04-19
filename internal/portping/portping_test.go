package portping_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portping"
)

func startTCP(t *testing.T) (port int, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port, func() { _ = ln.Close() }
}

func TestPing_AlivePort(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()
	p := portping.New(time.Second)
	res := p.Ping(context.Background(), "127.0.0.1", port)
	if !res.Alive() {
		t.Fatalf("expected alive, got err: %v", res.Err)
	}
	if res.Latency <= 0 {
		t.Fatalf("expected positive latency")
	}
}

func TestPing_ClosedPort(t *testing.T) {
	p := portping.New(200 * time.Millisecond)
	res := p.Ping(context.Background(), "127.0.0.1", 1)
	if res.Alive() {
		t.Fatal("expected not alive")
	}
}

func TestPingAll_ReturnsAllResults(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()
	p := portping.New(time.Second)
	results := p.PingAll(context.Background(), "127.0.0.1", []int{port, 1})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Alive() {
		t.Errorf("first result should be alive")
	}
	if results[1].Alive() {
		t.Errorf("second result should not be alive")
	}
}

func TestPing_DefaultTimeout(t *testing.T) {
	p := portping.New(0)
	if p == nil {
		t.Fatal("expected non-nil pinger")
	}
}

func TestResult_String_Alive(t *testing.T) {
	r := portping.Result{Host: "localhost", Port: 80, Latency: 500 * time.Microsecond}
	s := r.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
