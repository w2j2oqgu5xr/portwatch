package notify_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestBuild_ConsoleTarget(t *testing.T) {
	var buf bytes.Buffer
	n, err := notify.Build(notify.Config{Targets: "console", Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(baseEvent); err != nil {
		t.Fatalf("notify error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output in buffer")
	}
}

func TestBuild_FileTarget(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portwatch.log")

	n, err := notify.Build(notify.Config{Targets: "file", LogFile: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(baseEvent); err != nil {
		t.Fatalf("notify error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "8080") {
		t.Errorf("expected port in log file, got: %s", data)
	}
}

func TestBuild_FileTarget_MissingPath(t *testing.T) {
	_, err := notify.Build(notify.Config{Targets: "file"})
	if err == nil {
		t.Fatal("expected error when LogFile is empty")
	}
}

func TestBuild_WebhookTarget_MissingURL(t *testing.T) {
	_, err := notify.Build(notify.Config{Targets: "webhook"})
	if err == nil {
		t.Fatal("expected error when WebhookURL is empty")
	}
}

func TestBuild_EmptyTargetsFallsBackToConsole(t *testing.T) {
	var buf bytes.Buffer
	n, err := notify.Build(notify.Config{Targets: "", Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(baseEvent); err != nil {
		t.Fatalf("notify error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected fallback console output")
	}
}

func TestBuild_MultipleTargets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "multi.log")
	var buf bytes.Buffer

	n, err := notify.Build(notify.Config{
		Targets: "console,file",
		LogFile: path,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(baseEvent); err != nil {
		t.Fatalf("notify error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected console output")
	}
	data, _ := os.ReadFile(path)
	if len(data) == 0 {
		t.Error("expected file output")
	}
}
