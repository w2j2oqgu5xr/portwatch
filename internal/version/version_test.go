package version_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/version"
)

func TestGet_ReturnsDefaults(t *testing.T) {
	info := version.Get()
	if info.Version == "" {
		t.Error("expected non-empty Version")
	}
	if info.Commit == "" {
		t.Error("expected non-empty Commit")
	}
	if info.BuildDate == "" {
		t.Error("expected non-empty BuildDate")
	}
}

func TestInfo_String_ContainsVersion(t *testing.T) {
	version.Version = "1.2.3"
	version.Commit = "abc1234"
	version.BuildDate = "2024-01-01"

	info := version.Get()
	s := info.String()

	for _, want := range []string{"portwatch", "1.2.3", "abc1234", "2024-01-01"} {
		if !strings.Contains(s, want) {
			t.Errorf("String() = %q, want it to contain %q", s, want)
		}
	}
}

func TestInfo_String_Format(t *testing.T) {
	version.Version = "0.1.0"
	version.Commit = "deadbeef"
	version.BuildDate = "2024-06-15"

	info := version.Get()
	s := info.String()

	if !strings.HasPrefix(s, "portwatch ") {
		t.Errorf("expected string to start with 'portwatch ', got: %q", s)
	}
}
