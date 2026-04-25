// Package portschedule provides time-window based scanning schedules,
// allowing scans to be restricted to specific hours of the day.
package portschedule

import (
	"errors"
	"fmt"
	"time"
)

// Window represents a time range within a single day (24-hour clock).
type Window struct {
	StartHour int
	EndHour   int
}

// Schedule holds a set of active windows and a location for evaluation.
type Schedule struct {
	windows  []Window
	location *time.Location
}

// New creates a Schedule with the provided windows and timezone location.
// If loc is nil, UTC is used. Returns an error if any window is invalid.
func New(windows []Window, loc *time.Location) (*Schedule, error) {
	if loc == nil {
		loc = time.UTC
	}
	for _, w := range windows {
		if err := validateWindow(w); err != nil {
			return nil, err
		}
	}
	return &Schedule{windows: windows, location: loc}, nil
}

// Active reports whether the given time falls within any configured window.
// If no windows are configured, Active always returns true.
func (s *Schedule) Active(t time.Time) bool {
	if len(s.windows) == 0 {
		return true
	}
	local := t.In(s.location)
	h := local.Hour()
	for _, w := range s.windows {
		if w.StartHour <= w.EndHour {
			if h >= w.StartHour && h < w.EndHour {
				return true
			}
		} else {
			// overnight window e.g. 22-06
			if h >= w.StartHour || h < w.EndHour {
				return true
			}
		}
	}
	return false
}

// NextActive returns the next time at or after t when the schedule becomes active.
// If the schedule is already active at t, t is returned unchanged.
// If no windows are configured, t is returned immediately.
func (s *Schedule) NextActive(t time.Time) time.Time {
	if len(s.windows) == 0 || s.Active(t) {
		return t
	}
	// Step forward hour-by-hour (max 24 iterations).
	candidate := t.Truncate(time.Hour).Add(time.Hour)
	for i := 0; i < 25; i++ {
		if s.Active(candidate) {
			return candidate
		}
		candidate = candidate.Add(time.Hour)
	}
	return t
}

func validateWindow(w Window) error {
	if w.StartHour < 0 || w.StartHour > 23 {
		return fmt.Errorf("portschedule: start hour %d out of range [0,23]", w.StartHour)
	}
	if w.EndHour < 0 || w.EndHour > 23 {
		return fmt.Errorf("portschedule: end hour %d out of range [0,23]", w.EndHour)
	}
	if w.StartHour == w.EndHour {
		return errors.New("portschedule: start and end hour must differ")
	}
	return nil
}
