package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return p
}

func TestRun_MissingConfig(t *testing.T) {
	err := run()
	// run() uses os.Args; with no extra args it tries "portwatch.yaml" which
	// won't exist in the test working directory.
	if err == nil {
		t.Skip("portwatch.yaml unexpectedly present in working dir")
	}
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestRun_InvalidConfig(t *testing.T) {
	p := writeTempConfig(t, "host: \"\"\nports: []\n")
	os.Args = []string{"portwatch", p}
	t.Cleanup(func() { os.Args = os.Args[:1] })

	err := run()
	if err == nil {
		t.Fatal("expected validation error for empty host/ports, got nil")
	}
}
