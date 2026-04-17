package resolver_test

import (
	"testing"

	"github.com/user/portwatch/internal/resolver"
)

func TestName_WellKnown(t *testing.T) {
	r := resolver.New(nil)
	if got := r.Name(22); got != "ssh" {
		t.Fatalf("expected ssh, got %s", got)
	}
}

func TestName_Unknown(t *testing.T) {
	r := resolver.New(nil)
	if got := r.Name(9999); got != "9999" {
		t.Fatalf("expected 9999, got %s", got)
	}
}

func TestName_ExtraOverridesDefault(t *testing.T) {
	r := resolver.New(map[int]string{80: "myapp"})
	if got := r.Name(80); got != "myapp" {
		t.Fatalf("expected myapp, got %s", got)
	}
}

func TestName_ExtraCustomPort(t *testing.T) {
	r := resolver.New(map[int]string{12345: "internal-api"})
	if got := r.Name(12345); got != "internal-api" {
		t.Fatalf("expected internal-api, got %s", got)
	}
}

func TestLabel_Format(t *testing.T) {
	r := resolver.New(nil)
	if got := r.Label(443); got != "443/https" {
		t.Fatalf("expected 443/https, got %s", got)
	}
}

func TestLabel_UnknownPort(t *testing.T) {
	r := resolver.New(nil)
	if got := r.Label(55555); got != "55555/55555" {
		t.Fatalf("expected 55555/55555, got %s", got)
	}
}

func TestNew_NilExtraDoesNotPanic(t *testing.T) {
	r := resolver.New(nil)
	_ = r.Name(80)
}

func TestName_WellKnownPorts(t *testing.T) {
	r := resolver.New(nil)
	cases := []struct {
		port int
		want string
	}{
		{21, "ftp"},
		{25, "smtp"},
		{53, "dns"},
		{80, "http"},
		{443, "https"},
	}
	for _, tc := range cases {
		if got := r.Name(tc.port); got != tc.want {
			t.Errorf("Name(%d): expected %s, got %s", tc.port, tc.want, got)
		}
	}
}
