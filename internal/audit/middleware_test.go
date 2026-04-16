package audit_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/pipeline"
)

func TestStage_RecordsEvent(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	stage := audit.Stage(l)

	called := false
	stage(context.Background(), pipeline.Event{Type: "opened", Port: 9090, Host: "localhost"}, func(e pipeline.Event) {
		called = true
	})

	if !called {
		t.Error("expected next to be called")
	}
	if !strings.Contains(buf.String(), `"port":9090`) {
		t.Errorf("expected port in audit log, got: %s", buf.String())
	}
}

func TestStage_SeverityWarnOnOpened(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	stage := audit.Stage(l)
	stage(context.Background(), pipeline.Event{Type: "opened", Port: 80, Host: "localhost"}, func(e pipeline.Event) {})
	if !strings.Contains(buf.String(), `"severity":"warn"`) {
		t.Errorf("expected warn severity for opened event")
	}
}

func TestStage_SeverityInfoOnClosed(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	stage := audit.Stage(l)
	stage(context.Background(), pipeline.Event{Type: "closed", Port: 80, Host: "localhost"}, func(e pipeline.Event) {})
	if !strings.Contains(buf.String(), `"severity":"info"`) {
		t.Errorf("expected info severity for closed event")
	}
}
