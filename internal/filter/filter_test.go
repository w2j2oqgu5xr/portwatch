package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
)

func TestAllow_EmptyLists(t *testing.T) {
	f := filter.New(nil, nil)
	for _, port := range []int{22, 80, 443, 8080} {
		if !f.Allow(port) {
			t.Errorf("expected port %d to be allowed with empty filter", port)
		}
	}
}

func TestAllow_AllowList(t *testing.T) {
	f := filter.New([]int{80, 443}, nil)
	if !f.Allow(80) {
		t.Error("expected port 80 to be allowed")
	}
	if !f.Allow(443) {
		t.Error("expected port 443 to be allowed")
	}
	if f.Allow(22) {
		t.Error("expected port 22 to be blocked by allow list")
	}
}

func TestAllow_DenyList(t *testing.T) {
	f := filter.New(nil, []int{22, 23})
	if f.Allow(22) {
		t.Error("expected port 22 to be denied")
	}
	if f.Allow(23) {
		t.Error("expected port 23 to be denied")
	}
	if !f.Allow(80) {
		t.Error("expected port 80 to be allowed")
	}
}

func TestAllow_DenyOverridesAllow(t *testing.T) {
	f := filter.New([]int{80, 443}, []int{443})
	if !f.Allow(80) {
		t.Error("expected port 80 to be allowed")
	}
	if f.Allow(443) {
		t.Error("expected port 443 to be denied even though in allow list")
	}
}

func TestApply_FiltersSlice(t *testing.T) {
	f := filter.New([]int{80, 443}, nil)
	input := []int{22, 80, 443, 8080}
	got := f.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
	if got[0] != 80 || got[1] != 443 {
		t.Errorf("unexpected ports: %v", got)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	f := filter.New([]int{80}, nil)
	got := f.Apply([]int{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}

func TestValidate_ValidPorts(t *testing.T) {
	if err := filter.Validate([]int{1, 80, 443, 65535}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	if err := filter.Validate([]int{0}); err == nil {
		t.Error("expected error for port 0")
	}
	if err := filter.Validate([]int{65536}); err == nil {
		t.Error("expected error for port 65536")
	}
}
