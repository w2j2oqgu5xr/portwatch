package snapshot_test

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	ports := []int{80, 443, 8080}
	if err := snapshot.Save(path, ports); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	s, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(s.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(s.Ports))
	}
	if s.RecordedAt.IsZero() {
		t.Error("expected RecordedAt to be set")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if !errors.Is(err, snapshot.ErrNoSnapshot) {
		t.Fatalf("expected ErrNoSnapshot, got %v", err)
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o600)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for corrupt file, got nil")
	}
}

func TestDiff_OpenedAndClosed(t *testing.T) {
	prev := []int{80, 443, 22}
	current := []int{80, 8080}

	opened, closed := snapshot.Diff(prev, current)

	sort.Ints(opened)
	sort.Ints(closed)

	if len(opened) != 1 || opened[0] != 8080 {
		t.Errorf("expected opened=[8080], got %v", opened)
	}
	if len(closed) != 2 || closed[0] != 22 || closed[1] != 443 {
		t.Errorf("expected closed=[22 443], got %v", closed)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	ports := []int{80, 443}
	opened, closed := snapshot.Diff(ports, ports)

	if len(opened) != 0 {
		t.Errorf("expected no opened ports, got %v", opened)
	}
	if len(closed) != 0 {
		t.Errorf("expected no closed ports, got %v", closed)
	}
}

func TestDiff_EmptyPrev(t *testing.T) {
	current := []int{22, 80}
	opened, closed := snapshot.Diff(nil, current)

	sort.Ints(opened)
	if len(opened) != 2 || opened[0] != 22 || opened[1] != 80 {
		t.Errorf("expected opened=[22 80], got %v", opened)
	}
	if len(closed) != 0 {
		t.Errorf("expected no closed ports, got %v", closed)
	}
}
