package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.ndjson")
}

func TestRecord_CreatesFile(t *testing.T) {
	p := tempPath(t)
	rec := history.NewRecorder(p)

	err := rec.Record(history.Entry{Event: "opened", Host: "localhost", Port: 8080})
	if err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestRecord_TimestampAutoFilled(t *testing.T) {
	p := tempPath(t)
	rec := history.NewRecorder(p)

	before := time.Now().UTC().Add(-time.Second)
	_ = rec.Record(history.Entry{Event: "closed", Host: "127.0.0.1", Port: 22})
	after := time.Now().UTC().Add(time.Second)

	entries, err := history.ReadAll(p)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v outside expected range [%v, %v]", ts, before, after)
	}
}

func TestReadAll_MultipleEntries(t *testing.T) {
	p := tempPath(t)
	rec := history.NewRecorder(p)

	events := []history.Entry{
		{Event: "opened", Host: "localhost", Port: 3000},
		{Event: "opened", Host: "localhost", Port: 3001},
		{Event: "closed", Host: "localhost", Port: 3000},
	}
	for _, e := range events {
		if err := rec.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	entries, err := history.ReadAll(p)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != len(events) {
		t.Fatalf("expected %d entries, got %d", len(events), len(entries))
	}
	for i, e := range entries {
		if e.Event != events[i].Event || e.Port != events[i].Port {
			t.Errorf("entry %d mismatch: got %+v, want %+v", i, e, events[i])
		}
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	entries, err := history.ReadAll("/nonexistent/path/history.ndjson")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(entries))
	}
}

func TestRecord_AppendsBetweenCalls(t *testing.T) {
	p := tempPath(t)

	rec1 := history.NewRecorder(p)
	_ = rec1.Record(history.Entry{Event: "opened", Host: "h", Port: 1})

	// Simulate a second process / recorder instance.
	rec2 := history.NewRecorder(p)
	_ = rec2.Record(history.Entry{Event: "closed", Host: "h", Port: 1})

	entries, err := history.ReadAll(p)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}
