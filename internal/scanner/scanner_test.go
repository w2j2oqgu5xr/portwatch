package scanner

import (
	"net"
	"strconv"
	"testing"
	"time"
)

// startTCPServer opens a local TCP listener and returns its port and a stop func.
func startTCPServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	return port, func() { _ = ln.Close() }
}

func TestScan_OpenPort(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	results, err := Scan(ScanOptions{
		Host:    "127.0.0.1",
		Ports:   []int{port},
		Timeout: time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	results, err := Scan(ScanOptions{
		Host:    "127.0.0.1",
		Ports:   []int{1},
		Timeout: 200 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestScan_EmptyHost(t *testing.T) {
	_, err := Scan(ScanOptions{Host: "", Ports: []int{80}})
	if err == nil {
		t.Error("expected error for empty host")
	}
}

func TestScan_InvalidPort(t *testing.T) {
	_, err := Scan(ScanOptions{Host: "127.0.0.1", Ports: []int{0}})
	if err == nil {
		t.Error("expected error for invalid port 0")
	}
}

func TestOpenPorts(t *testing.T) {
	states := []PortState{
		{Port: 80, Open: true},
		{Port: 81, Open: false},
		{Port: 443, Open: true},
	}
	open := OpenPorts(states)
	if len(open) != 2 {
		t.Errorf("expected 2 open ports, got %d", len(open))
	}
}
