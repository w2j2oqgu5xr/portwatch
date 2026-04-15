package notify_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
)

var baseEvent = notify.Event{
	Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	Host:      "localhost",
	Port:      8080,
	Kind:      "opened",
	Message:   "new port detected",
}

func TestWriterNotifier_FormatsEvent(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewWriterNotifier(&buf)
	if err := n.Notify(baseEvent); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "opened") {
		t.Errorf("expected 'opened' in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected host in output, got: %s", out)
	}
}

func TestMultiNotifier_FansOut(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	m := notify.NewMulti(
		notify.NewWriterNotifier(&buf1),
		notify.NewWriterNotifier(&buf2),
	)
	if err := m.Notify(baseEvent); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf1.Len() == 0 || buf2.Len() == 0 {
		t.Error("expected both notifiers to receive the event")
	}
}

func TestMultiNotifier_CollectsErrors(t *testing.T) {
	var buf bytes.Buffer
	bad := &errorNotifier{err: errors.New("fail")}
	m := notify.NewMulti(notify.NewWriterNotifier(&buf), bad)
	err := m.Notify(baseEvent)
	if err == nil {
		t.Fatal("expected error from multi notifier")
	}
	if !strings.Contains(err.Error(), "fail") {
		t.Errorf("expected 'fail' in error, got: %v", err)
	}
}

func TestWebhookNotifier_PostsEvent(t *testing.T) {
	var received string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		buf.ReadFrom(r.Body)
		received = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wh := notify.NewWebhookNotifier(ts.URL)
	if err := wh.Notify(baseEvent); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(received, "8080") {
		t.Errorf("expected port in webhook body, got: %s", received)
	}
}

func TestWebhookNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	wh := notify.NewWebhookNotifier(ts.URL)
	if err := wh.Notify(baseEvent); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

// errorNotifier is a test helper that always returns an error.
type errorNotifier struct{ err error }

func (e *errorNotifier) Notify(_ notify.Event) error { return e.err }
