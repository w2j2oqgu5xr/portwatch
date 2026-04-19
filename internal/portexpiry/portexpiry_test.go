package portexpiry_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portexpiry"
)

func TestNew_DefaultGrace(t *testing.T) {
	tr := portexpiry.New(0)
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestMarkClosed_TracksPending(t *testing.T) {
	tr := portexpiry.New(time.Minute)
	tr.MarkClosed(8080)
	if got := tr.Pending(); got != 1 {
		t.Fatalf("expected 1 pending, got %d", got)
	}
}

func TestMarkOpen_RemovesPending(t *testing.T) {
	tr := portexpiry.New(time.Minute)
	tr.MarkClosed(8080)
	tr.MarkOpen(8080)
	if got := tr.Pending(); got != 0 {
		t.Fatalf("expected 0 pending, got %d", got)
	}
}

func TestExpired_EmptyBeforeGrace(t *testing.T) {
	tr := portexpiry.New(time.Hour)
	tr.MarkClosed(443)
	if entries := tr.Expired(); len(entries) != 0 {
		t.Fatalf("expected no expired entries, got %d", len(entries))
	}
}

func TestExpired_ReturnsAfterGrace(t *testing.T) {
	tr := portexpiry.New(time.Millisecond)
	tr.MarkClosed(9000)
	time.Sleep(5 * time.Millisecond)
	entries := tr.Expired()
	if len(entries) != 1 {
		t.Fatalf("expected 1 expired entry, got %d", len(entries))
	}
	if entries[0].Port != 9000 {
		t.Errorf("expected port 9000, got %d", entries[0].Port)
	}
}

func TestExpired_RemovesFromTracking(t *testing.T) {
	tr := portexpiry.New(time.Millisecond)
	tr.MarkClosed(7070)
	time.Sleep(5 * time.Millisecond)
	tr.Expired()
	if got := tr.Pending(); got != 0 {
		t.Fatalf("expected 0 pending after expiry, got %d", got)
	}
}

func TestMarkClosed_IdempotentTimestamp(t *testing.T) {
	tr := portexpiry.New(time.Hour)
	tr.MarkClosed(22)
	tr.MarkClosed(22) // second call should not reset
	if got := tr.Pending(); got != 1 {
		t.Fatalf("expected 1 pending, got %d", got)
	}
}

func TestExpired_MultiplePortsMixed(t *testing.T) {
	tr := portexpiry.New(time.Millisecond)
	tr.MarkClosed(1111)
	time.Sleep(5 * time.Millisecond)
	tr.MarkClosed(2222) // added after sleep, still within grace
	entries := tr.Expired()
	if len(entries) != 1 || entries[0].Port != 1111 {
		t.Errorf("expected only port 1111 expired, got %+v", entries)
	}
	if got := tr.Pending(); got != 1 {
		t.Fatalf("expected 1 still pending, got %d", got)
	}
}
