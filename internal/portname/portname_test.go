package portname_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portname"
)

func TestLookup_WellKnown(t *testing.T) {
	result := portname.Lookup(80)
	if result != "80 (http)" {
		t.Fatalf("expected '80 (http)', got %q", result)
	}
}

func TestLookup_Unknown(t *testing.T) {
	result := portname.Lookup(9999)
	if result != "9999" {
		t.Fatalf("expected '9999', got %q", result)
	}
}

func TestLookup_ContainsPortNumber(t *testing.T) {
	for _, port := range []int{22, 443, 3306} {
		result := portname.Lookup(port)
		if !strings.Contains(result, fmt.Sprintf("%d", port)) {
			t.Errorf("Lookup(%d) = %q: missing port number", port, result)
		}
	}
}

func TestName_WellKnown(t *testing.T) {
	if got := portname.Name(22); got != "ssh" {
		t.Fatalf("expected 'ssh', got %q", got)
	}
}

func TestName_Unknown(t *testing.T) {
	if got := portname.Name(1); got != "unknown" {
		t.Fatalf("expected 'unknown', got %q", got)
	}
}

func TestName_Postgres(t *testing.T) {
	if got := portname.Name(5432); got != "postgres" {
		t.Fatalf("expected 'postgres', got %q", got)
	}
}

import "fmt"
