// Package alert provides alerting mechanisms for port state changes.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a port state change event.
type Event struct {
	Host      string
	Port      int
	PrevState bool
	NewState  bool
	Timestamp time.Time
}

// Notifier is the interface for sending alerts.
type Notifier interface {
	Notify(level Level, event Event) error
}

// ConsoleNotifier writes alerts to an io.Writer (defaults to os.Stdout).
type ConsoleNotifier struct {
	Out io.Writer
}

// NewConsoleNotifier creates a ConsoleNotifier writing to stdout.
func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{Out: os.Stdout}
}

// Notify formats and writes the event to the configured writer.
func (c *ConsoleNotifier) Notify(level Level, event Event) error {
	state := "closed"
	if event.NewState {
		state = "open"
	}
	_, err := fmt.Fprintf(
		c.Out,
		"[%s] [%s] %s:%d is now %s\n",
		event.Timestamp.Format(time.RFC3339),
		level,
		event.Host,
		event.Port,
		state,
	)
	return err
}

// LevelForEvent determines the appropriate alert level for a port change.
func LevelForEvent(event Event) Level {
	if event.NewState && !event.PrevState {
		// Port unexpectedly opened
		return LevelAlert
	}
	if !event.NewState && event.PrevState {
		// Port closed
		return LevelWarn
	}
	return LevelInfo
}
