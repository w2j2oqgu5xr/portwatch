package portping_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portping"
)

func TestSummarise_AllAlive(t *testing.T) {
	results := []portping.Result{
		{Host: "h", Port: 80, Latency: 1 * time.Millisecond},
		{Host: "h", Port: 80, Latency: 3 * time.Millisecond},
		{Host: "h", Port: 80, Latency: 2 * time.Millisecond},
	}
	s := portping.Summarise("h", 80, results)
	if s.Sent != 3 || s.Received != 3 {
		t.Fatalf("sent=%d recv=%d", s.Sent, s.Received)
	}
	if s.Min != 1*time.Millisecond {
		t.Errorf("min: got %s want 1ms", s.Min)
	}
	if s.Max != 3*time.Millisecond {
		t.Errorf("max: got %s want 3ms", s.Max)
	}
	if s.Avg != 2*time.Millisecond {
		t.Errorf("avg: got %s want 2ms", s.Avg)
	}
	if s.PacketLoss() != 0 {
		t.Errorf("expected 0 loss")
	}
}

func TestSummarise_SomeLost(t *testing.T) {
	results := []portping.Result{
		{Host: "h", Port: 80, Latency: 1 * time.Millisecond},
		{Host: "h", Port: 80, Err: fmt.Errorf("refused")},
	}
	s := portping.Summarise("h", 80, results)
	if s.PacketLoss() != 50 {
		t.Errorf("expected 50%% loss, got %.0f", s.PacketLoss())
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := portping.Summarise("h", 80, nil)
	if s.Sent != 0 || s.Received != 0 {
		t.Errorf("unexpected counts")
	}
}

func TestStats_String(t *testing.T) {
	s := portping.Stats{Host: "localhost", Port: 443, Sent: 4, Received: 4,
		Min: time.Millisecond, Max: 3 * time.Millisecond, Avg: 2 * time.Millisecond}
	if s.String() == "" {
		t.Fatal("expected non-empty string")
	}
}
