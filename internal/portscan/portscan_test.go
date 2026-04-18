package portscan_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portscan"
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
			conn.Write([]byte("hello"))
			conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port, func() { ln.Close() }
}

func collectResults(ch <-chan portscan.Result) []portscan.Result {
	var out []portscan.Result
	for r := range ch {
		out = append(out, r)
	}
	return out
}

func TestScan_OpenPort(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()

	ctx := context.Background()
	ch := portscan.Scan(ctx, "127.0.0.1", []int{port}, portscan.Options{Timeout: time.Second})
	results := collectResults(ch)

	if len(results) != 1 || !results[0].Open {
		t.Fatalf("expected open port, got %+v", results)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	ctx := context.Background()
	ch := portscan.Scan(ctx, "127.0.0.1", []int{1}, portscan.Options{Timeout: 200 * time.Millisecond})
	results := collectResults(ch)

	if len(results) != 1 || results[0].Open {
		t.Fatalf("expected closed port, got %+v", results)
	}
}

func TestScan_BannerGrab(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()

	ctx := context.Background()
	ch := portscan.Scan(ctx, "127.0.0.1", []int{port}, portscan.Options{
		Timeout:    time.Second,
		GrabBanner: true,
	})
	results := collectResults(ch)

	if len(results) != 1 || results[0].Banner != "hello" {
		t.Fatalf("expected banner 'hello', got %+v", results)
	}
}

func TestScan_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ports := make([]int, 50)
	for i := range ports {
		ports[i] = 10000 + i
	}
	ch := portscan.Scan(ctx, "127.0.0.1", ports, portscan.Options{Timeout: time.Second})
	// should drain without blocking
	for range ch {
	}
}
