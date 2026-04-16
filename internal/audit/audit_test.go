package audit_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/audit"
)

func baseEntry() audit.Entry {
	return audit.Entry{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Event:     "opened",
		Port:      8080,
		Host:      "localhost",
		Severity:  "warn",
	}
}

func TestRecord_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	if err := l.Record(baseEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"port":8080`) {
		t.Errorf("expected port in output, got: %s", buf.String())
	}
}

func TestRecord_AutoFillsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	e := baseEntry()
	e.Timestamp = time.Time{}
	if err := l.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"timestamp"`) {
		t.Errorf("expected timestamp in output")
	}
}

func TestReadAll_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	for i := 0; i < 3; i++ {
		e := baseEntry()
		e.Port = 8080 + i
		l.Record(e)
	}
	entries, err := audit.ReadAll(&buf)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestReadAll_EmptyReader(t *testing.T) {
	entries, err := audit.ReadAll(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries")
	}
}

func TestReadAll_CorruptLine(t *testing.T) {
	_, err := audit.ReadAll(strings.NewReader("{bad json\n"))
	if err == nil {
		t.Error("expected error on corrupt input")
	}
}

func TestNewFileLogger_BadPath(t *testing.T) {
	_, err := audit.NewFileLogger("/nonexistent/dir/audit.log")
	if err == nil {
		t.Error("expected error for bad path")
	}
}
