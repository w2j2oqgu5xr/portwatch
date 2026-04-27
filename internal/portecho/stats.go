package portecho

import (
	"fmt"
	"time"
)

// Stats summarises a batch of echo probe results.
type Stats struct {
	Total   int
	Open    int
	Closed  int
	MinRTT  time.Duration
	MaxRTT  time.Duration
	MeanRTT time.Duration
}

// String returns a compact, human-readable summary.
func (s Stats) String() string {
	return fmt.Sprintf(
		"total=%d open=%d closed=%d min=%s mean=%s max=%s",
		s.Total, s.Open, s.Closed,
		s.MinRTT.Round(time.Microsecond),
		s.MeanRTT.Round(time.Microsecond),
		s.MaxRTT.Round(time.Microsecond),
	)
}

// Summarise computes aggregate statistics over a slice of Results.
// Latency values from closed ports are excluded from RTT calculations.
func Summarise(results []Result) Stats {
	if len(results) == 0 {
		return Stats{}
	}

	s := Stats{Total: len(results)}
	var total time.Duration
	var first = true

	for _, r := range results {
		if r.Open {
			s.Open++
			total += r.Latency
			if first || r.Latency < s.MinRTT {
				s.MinRTT = r.Latency
			}
			if r.Latency > s.MaxRTT {
				s.MaxRTT = r.Latency
			}
			first = false
		} else {
			s.Closed++
		}
	}

	if s.Open > 0 {
		s.MeanRTT = total / time.Duration(s.Open)
	}
	return s
}
