package resolver_test

import (
	"testing"

	"github.com/user/portwatch/internal/resolver"
)

func TestServiceTable_ContainsHTTP(t *testing.T) {
	tbl := resolver.ServiceTable()
	if tbl[80] != "http" {
		t.Fatalf("expected http for 80, got %s", tbl[80])
	}
}

func TestServiceTable_IsCopy(t *testing.T) {
	tbl := resolver.ServiceTable()
	tbl[80] = "modified"

	r := resolver.New(nil)
	if r.Name(80) != "http" {
		t.Fatal("mutating ServiceTable result should not affect internal table")
	}
}

func TestMerge_OverrideWins(t *testing.T) {
	base := map[int]string{80: "http", 22: "ssh"}
	over := map[int]string{80: "custom"}
	out := resolver.Merge(base, over)
	if out[80] != "custom" {
		t.Fatalf("expected custom, got %s", out[80])
	}
	if out[22] != "ssh" {
		t.Fatalf("expected ssh, got %s", out[22])
	}
}

func TestMerge_BasePreserved(t *testing.T) {
	base := map[int]string{443: "https"}
	out := resolver.Merge(base, nil)
	if out[443] != "https" {
		t.Fatalf("expected https, got %s", out[443])
	}
}
