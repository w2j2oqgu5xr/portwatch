package portevict_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portevict"
)

func TestRecord_StoresEntry(t *testing.T) {
	l := portevict.New(10)
	opened := time.Now().Add(-5 * time.Second)
	e := l.Record(8080, "localhost", opened)

	if e.Port != 8080 {
		t.Fatalf("expected port 8080, got %d", e.Port)
	}
	if e.Host != "localhost" {
		t.Fatalf("expected host localhost, got %s", e.Host)
	}
	if e.Duration < 4*time.Second {
		t.Fatalf("expected duration >= 4s, got %v", e.Duration)
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	l := portevict.New(10)
	opened := time.Now().Add(-time.Second)
	l.Record(80, "host", opened)
	l.Record(443, "host", opened)

	if got := len(l.All()); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}

func TestLog_CapEnforcedOnOverflow(t *testing.T) {
	l := portevict.New(3)
	opened := time.Now()
	for i := 0; i < 5; i++ {
		l.Record(1000+i, "h", opened)
	}
	if l.Len() != 3 {
		t.Fatalf("expected cap 3, got %d", l.Len())
	}
	entries := l.All()
	if entries[0].Port != 1002 {
		t.Fatalf("expected oldest evicted entry port 1002, got %d", entries[0].Port)
	}
}

func TestClear_RemovesAllEntries(t *testing.T) {
	l := portevict.New(10)
	opened := time.Now()
	l.Record(22, "h", opened)
	l.Clear()
	if l.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", l.Len())
	}
}

func TestNew_DefaultMax(t *testing.T) {
	l := portevict.New(0)
	opened := time.Now()
	for i := 0; i < 300; i++ {
		l.Record(i, "h", opened)
	}
	if l.Len() != 256 {
		t.Fatalf("expected default cap 256, got %d", l.Len())
	}
}
