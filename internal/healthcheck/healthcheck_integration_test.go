package healthcheck_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

type staticProvider struct {
	scans     uint64
	alerts    uint64
	openPorts int
}

func (s *staticProvider) Scans() uint64    { return s.scans }
func (s *staticProvider) Alerts() uint64   { return s.alerts }
func (s *staticProvider) OpenPorts() int   { return s.openPorts }

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

func TestIntegration_ListenAndServe(t *testing.T) {
	addr := freePort(t)
	p := &staticProvider{scans: 7, alerts: 3, openPorts: 2}

	srv := healthcheck.New(addr, p)
	srv.SetReady(true)

	go func() { _ = srv.ListenAndServe() }()
	defer srv.Shutdown()

	// Wait briefly for the server to start.
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var status healthcheck.Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatal(err)
	}
	if status.Scans != 7 {
		t.Errorf("expected scans=7, got %d", status.Scans)
	}
	if status.Alerts != 3 {
		t.Errorf("expected alerts=3, got %d", status.Alerts)
	}
	if status.OpenPorts != 2 {
		t.Errorf("expected open_ports=2, got %d", status.OpenPorts)
	}
	if status.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
}

func TestIntegration_ShutdownStopsServer(t *testing.T) {
	addr := freePort(t)
	p := &staticProvider{}
	srv := healthcheck.New(addr, p)
	srv.SetReady(true)

	go func() { _ = srv.ListenAndServe() }()
	time.Sleep(30 * time.Millisecond)

	if err := srv.Shutdown(); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}

	time.Sleep(20 * time.Millisecond)
	_, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	if err == nil {
		t.Error("expected connection refused after shutdown")
	}
}
