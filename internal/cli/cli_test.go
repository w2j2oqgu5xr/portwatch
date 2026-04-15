package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestParse_Defaults(t *testing.T) {
	opts, err := Parse([]string{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.ConfigPath != "portwatch.yaml" {
		t.Errorf("ConfigPath = %q, want %q", opts.ConfigPath, "portwatch.yaml")
	}
	if opts.Verbose {
		t.Error("Verbose should default to false")
	}
	if opts.Version {
		t.Error("Version should default to false")
	}
}

func TestParse_CustomConfig(t *testing.T) {
	opts, err := Parse([]string{"-config", "/etc/pw.yaml"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.ConfigPath != "/etc/pw.yaml" {
		t.Errorf("ConfigPath = %q, want /etc/pw.yaml", opts.ConfigPath)
	}
}

func TestParse_VerboseAndVersion(t *testing.T) {
	opts, err := Parse([]string{"-verbose", "-version"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Verbose {
		t.Error("expected Verbose=true")
	}
	if !opts.Version {
		t.Error("expected Version=true")
	}
}

func TestParse_UnknownFlag(t *testing.T) {
	var buf bytes.Buffer
	_, err := Parse([]string{"-unknown"}, &buf)
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestPrintVersion(t *testing.T) {
	var buf bytes.Buffer
	PrintVersion(&buf)
	if !strings.Contains(buf.String(), "portwatch") {
		t.Errorf("version output missing 'portwatch': %q", buf.String())
	}
	if !strings.Contains(buf.String(), version) {
		t.Errorf("version output missing version string %q: %q", version, buf.String())
	}
}
