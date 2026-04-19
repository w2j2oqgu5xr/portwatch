package portwatch_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/portwatch"
)

// captureNotifier records events sent through the pipeline.
type captureNotifier struct {
	events []pipeline.Event
}

func (c *captureNotifier) Notify(_ context.Context, e pipeline.Event) error {
	c.events = append(c.events, e)
	return nil
}

func startTCP(t *testing.T) (port int, stop func()) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return l.Addr().(*net.TCPAddr).Port, func() { l.Close() }
}

func TestWatcher_DetectsOpenedPort(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()

	cap := &captureNotifier{}
	p := pipeline.New(cap)

	cfg := portwatch.Config{
		Host:  "127.0.0.1",
		Ports: portwatch.PortRange{From: port, To: port},
		Interval: 5 * time.Second,
	}
	w, err := portwatch.New(cfg, p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	w.Run(ctx) //nolint:errcheck

	if len(cap.events) == 0 {
		t.Fatal("expected at least one event, got none")
	}
	if cap.events[0].Kind != portwatch.EventOpened {
		t.Fatalf("expected %q, got %q", portwatch.EventOpened, cap.events[0].Kind)
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := portwatch.New(portwatch.Config{}, nil)
	if err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestConfig_Validate_BadRange(t *testing.T) {
	cfg := portwatch.Config{
		Host:  "localhost",
		Ports: portwatch.PortRange{From: 100, To: 10},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for inverted port range")
	}
}
