package portping

import (
	"fmt"
	"time"
)

// Stats summarises a series of Results for a single port.
type Stats struct {
	Host    string
	Port    int
	Sent    int
	Received int
	Min     time.Duration
	Max     time.Duration
	Avg     time.Duration
}

// PacketLoss returns the percentage of lost pings (0–100).
func (s Stats) PacketLoss() float64 {
	if s.Sent == 0 {
		return 0
	}
	return float64(s.Sent-s.Received) / float64(s.Sent) * 100
}

// String returns a compact summary line.
func (s Stats) String() string {
	return fmt.Sprintf("%s:%d sent=%d recv=%d loss=%.0f%% min=%s avg=%s max=%s",
		s.Host, s.Port, s.Sent, s.Received, s.PacketLoss(),
		s.Min.Round(time.Microsecond), s.Avg.Round(time.Microsecond), s.Max.Round(time.Microsecond))
}

// Summarise aggregates a slice of Results into Stats.
func Summarise(host string, port int, results []Result) Stats {
	s := Stats{Host: host, Port: port, Sent: len(results)}
	if len(results) == 0 {
		return s
	}
	var total time.Duration
	first := true
	for _, r := range results {
		if !r.Alive() {
			continue
		}
		s.Received++
		total += r.Latency
		if first || r.Latency < s.Min {
			s.Min = r.Latency
		}
		if r.Latency > s.Max {
			s.Max = r.Latency
		}
		first = false
	}
	if s.Received > 0 {
		s.Avg = total / time.Duration(s.Received)
	}
	return s
}
