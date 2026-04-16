package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tempPath(t)
	ports := []int{22, 80, 443}
	if err := baseline.Save(path, ports); err != nil {
		t.Fatalf("Save: %v", err)
	}
	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(b.Ports) != len(ports) {
		t.Errorf("got %d ports, want %d", len(b.Ports), len(ports))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err != baseline.ErrNoBaseline {
		t.Errorf("expected ErrNoBaseline, got %v", err)
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	path := tempPath(t)
	os.WriteFile(path, []byte("not json{"), 0644)
	_, err := baseline.Load(path)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}

func TestUnexpected_ReturnsNewPorts(t *testing.T) {
	b := &baseline.Baseline{Ports: []int{22, 80}}
	got := b.Unexpected([]int{22, 80, 8080})
	if len(got) != 1 || got[0] != 8080 {
		t.Errorf("unexpected ports: %v", got)
	}
}

func TestUnexpected_EmptyWhenAllTrusted(t *testing.T) {
	b := &baseline.Baseline{Ports: []int{22, 80}}
	got := b.Unexpected([]int{22, 80})
	if len(got) != 0 {
		t.Errorf("expected no unexpected ports, got %v", got)
	}
}

func TestMissing_ReturnsMissingPorts(t *testing.T) {
	b := &baseline.Baseline{Ports: []int{22, 80, 443}}
	got := b.Missing([]int{22})
	if len(got) != 2 {
		t.Errorf("expected 2 missing ports, got %v", got)
	}
}

func TestMissing_NoneWhenAllPresent(t *testing.T) {
	b := &baseline.Baseline{Ports: []int{22, 80}}
	got := b.Missing([]int{22, 80, 443})
	if len(got) != 0 {
		t.Errorf("expected none missing, got %v", got)
	}
}
