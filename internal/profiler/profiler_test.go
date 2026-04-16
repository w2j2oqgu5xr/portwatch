package profiler_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/user/portwatch/internal/profiler"
)

func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freeAddr: %v", err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

func TestStart_ServesDebugEndpoint(t *testing.T) {
	addr := freeAddr(t)
	s := profiler.New(addr)
	if err := s.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Shutdown(context.Background())

	url := fmt.Sprintf("http://%s/debug/pprof/", addr)
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestAddr_ReturnsConfiguredAddr(t *testing.T) {
	s := profiler.New("localhost:6060")
	if got := s.Addr(); got != "localhost:6060" {
		t.Errorf("Addr() = %q, want %q", got, "localhost:6060")
	}
}

func TestShutdown_StopsServer(t *testing.T) {
	addr := freeAddr(t)
	s := profiler.New(addr)
	if err := s.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown: %v", err)
	}
}
