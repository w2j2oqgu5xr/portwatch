package portgroup_test

import (
	"testing"

	"github.com/user/portwatch/internal/portgroup"
)

func TestAdd_And_Get(t *testing.T) {
	r := portgroup.New()
	if err := r.Add("web", []int{80, 443}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, ok := r.Get("web")
	if !ok {
		t.Fatal("expected group to exist")
	}
	if len(g.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(g.Ports))
	}
}

func TestAdd_DuplicateName(t *testing.T) {
	r := portgroup.New()
	r.Add("web", []int{80})
	if err := r.Add("web", []int{443}); err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestAdd_EmptyPorts(t *testing.T) {
	r := portgroup.New()
	if err := r.Add("empty", []int{}); err == nil {
		t.Fatal("expected error for empty ports")
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	r := portgroup.New()
	if err := r.Add("bad", []int{0}); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Add("bad2", []int{65536}); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestResolve_DeduplicatesPorts(t *testing.T) {
	r := portgroup.New()
	r.Add("web", []int{80, 443})
	r.Add("api", []int{443, 8080})
	ports, err := r.Resolve([]string{"web", "api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 3 {
		t.Errorf("expected 3 unique ports, got %d", len(ports))
	}
}

func TestResolve_UnknownGroup(t *testing.T) {
	r := portgroup.New()
	_, err := r.Resolve([]string{"unknown"})
	if err == nil {
		t.Fatal("expected error for unknown group")
	}
}

func TestAll_ReturnsAllGroups(t *testing.T) {
	r := portgroup.New()
	r.Add("web", []int{80})
	r.Add("db", []int{5432})
	if len(r.All()) != 2 {
		t.Errorf("expected 2 groups, got %d", len(r.All()))
	}
}
