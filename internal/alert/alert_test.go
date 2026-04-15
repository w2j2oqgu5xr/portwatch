package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func baseEvent(prevState, newState bool) alert.Event {
	return alert.Event{
		Host:      "localhost",
		Port:      8080,
		PrevState: prevState,
		NewState:  newState,
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}
}

func TestConsoleNotifier_PortOpened(t *testing.T) {
	var buf bytes.Buffer
	n := &alert.ConsoleNotifier{Out: &buf}

	evt := baseEvent(false, true)
	if err := n.Notify(alert.LevelAlert, evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", output)
	}
	if !strings.Contains(output, "open") {
		t.Errorf("expected 'open' in output, got: %s", output)
	}
	if !strings.Contains(output, "localhost:8080") {
		t.Errorf("expected host:port in output, got: %s", output)
	}
}

func TestConsoleNotifier_PortClosed(t *testing.T) {
	var buf bytes.Buffer
	n := &alert.ConsoleNotifier{Out: &buf}

	evt := baseEvent(true, false)
	if err := n.Notify(alert.LevelWarn, evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "WARN") {
		t.Errorf("expected WARN in output, got: %s", output)
	}
	if !strings.Contains(output, "closed") {
		t.Errorf("expected 'closed' in output, got: %s", output)
	}
}

func TestConsoleNotifier_TimestampInOutput(t *testing.T) {
	var buf bytes.Buffer
	n := &alert.ConsoleNotifier{Out: &buf}

	evt := baseEvent(false, true)
	if err := n.Notify(alert.LevelAlert, evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "2024-01-15") {
		t.Errorf("expected timestamp in output, got: %s", output)
	}
}

func TestLevelForEvent_PortOpened(t *testing.T) {
	evt := baseEvent(false, true)
	if got := alert.LevelForEvent(evt); got != alert.LevelAlert {
		t.Errorf("expected ALERT, got %s", got)
	}
}

func TestLevelForEvent_PortClosed(t *testing.T) {
	evt := baseEvent(true, false)
	if got := alert.LevelForEvent(evt); got != alert.LevelWarn {
		t.Errorf("expected WARN, got %s", got)
	}
}

func TestLevelForEvent_NoChange(t *testing.T) {
	evt := baseEvent(true, true)
	if got := alert.LevelForEvent(evt); got != alert.LevelInfo {
		t.Errorf("expected INFO, got %s", got)
	}
}
