package monitor_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

func startTCPServer(t *testing.T, port int) (stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		t.Fatalf("failed to start TCP server on port %d: %v", port, err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return func() { ln.Close() }
}

func TestMonitor_DetectsOpenedPort(t *testing.T) {
	const port = 19201

	m := monitor.New("127.0.0.1", []int{port}, 50*time.Millisecond)
	go m.Start()
	defer close(m.Stop)

	// Give monitor one tick to establish baseline (port closed).
	time.Sleep(80 * time.Millisecond)

	stop := startTCPServer(t, port)
	defer stop()

	select {
	case change := <-m.Changes:
		if change.Port != port || change.Status != "opened" {
			t.Errorf("unexpected change: %+v", change)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for 'opened' change")
	}
}

func TestMonitor_DetectsClosedPort(t *testing.T) {
	const port = 19202

	stop := startTCPServer(t, port)

	m := monitor.New("127.0.0.1", []int{port}, 50*time.Millisecond)
	go m.Start()
	defer close(m.Stop)

	// Allow baseline to capture port as open.
	time.Sleep(80 * time.Millisecond)
	stop()

	select {
	case change := <-m.Changes:
		if change.Port != port || change.Status != "closed" {
			t.Errorf("unexpected change: %+v", change)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for 'closed' change")
	}
}

func TestMonitor_NoFalsePositives(t *testing.T) {
	const port = 19203
	stop := startTCPServer(t, port)
	defer stop()

	m := monitor.New("127.0.0.1", []int{port}, 50*time.Millisecond)
	go m.Start()
	defer close(m.Stop)

	// Port stays open — no changes expected.
	time.Sleep(200 * time.Millisecond)

	select {
	case change := <-m.Changes:
		t.Errorf("unexpected change for stable port: %+v", change)
	default:
		// pass
	}
}
