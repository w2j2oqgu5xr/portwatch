package banner_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/banner"
)

func startBannerServer(t *testing.T, msg string) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.String())
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port64, _ := strconv.ParseInt(portStr, 10, 32)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.Write([]byte(msg))
	}()
	t.Cleanup(func() { ln.Close() })
	return int(port64)
}

func TestGrab_ReturnsBanner(t *testing.T) {
	port := startBannerServer(t, "SSH-2.0-OpenSSH")
	g := banner.New(time.Second)
	r := g.Grab("127.0.0.1", port)
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if r.Banner != "SSH-2.0-OpenSSH" {
		t.Errorf("got %q, want SSH-2.0-OpenSSH", r.Banner)
	}
}

func TestGrab_ClosedPort_ReturnsError(t *testing.T) {
	g := banner.New(200 * time.Millisecond)
	r := g.Grab("127.0.0.1", 1)
	if r.Err == nil {
		t.Error("expected error for closed port")
	}
}

func TestGrabAll_ReturnsAllResults(t *testing.T) {
	p1 := startBannerServer(t, "HTTP/1.1 200 OK")
	p2 := startBannerServer(t, "220 SMTP")
	g := banner.New(time.Second)
	results := g.GrabAll("127.0.0.1", []int{p1, p2})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestNew_DefaultTimeout(t *testing.T) {
	g := banner.New(0)
	if g.Timeout != banner.DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", banner.DefaultTimeout, g.Timeout)
	}
}
