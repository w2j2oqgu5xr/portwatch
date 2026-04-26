package portquota_test

import (
	"testing"

	"github.com/user/portwatch/internal/portquota"
)

func TestNew_DefaultLimitClamped(t *testing.T) {
	e := portquota.New(0)
	if e.Limit() != 1 {
		t.Fatalf("expected limit 1, got %d", e.Limit())
	}
}

func TestNew_ValidLimit(t *testing.T) {
	e := portquota.New(10)
	if e.Limit() != 10 {
		t.Fatalf("expected limit 10, got %d", e.Limit())
	}
}

func TestSet_BelowLimit_ReturnsNil(t *testing.T) {
	e := portquota.New(5)
	if v := e.Set(3); v != nil {
		t.Fatalf("expected no violation, got %v", v)
	}
}

func TestSet_AtLimit_ReturnsNil(t *testing.T) {
	e := portquota.New(5)
	if v := e.Set(5); v != nil {
		t.Fatalf("expected no violation at limit, got %v", v)
	}
}

func TestSet_ExceedsLimit_ReturnsViolation(t *testing.T) {
	e := portquota.New(5)
	v := e.Set(7)
	if v == nil {
		t.Fatal("expected violation, got nil")
	}
	if v.Limit != 5 {
		t.Errorf("violation limit: want 5, got %d", v.Limit)
	}
	if v.Actual != 7 {
		t.Errorf("violation actual: want 7, got %d", v.Actual)
	}
}

func TestSet_NegativeCountClamped(t *testing.T) {
	e := portquota.New(3)
	if v := e.Set(-5); v != nil {
		t.Fatalf("expected no violation for negative count, got %v", v)
	}
	if e.Count() != 0 {
		t.Errorf("expected count 0, got %d", e.Count())
	}
}

func TestCount_ReflectsLastSet(t *testing.T) {
	e := portquota.New(10)
	e.Set(4)
	if e.Count() != 4 {
		t.Errorf("expected count 4, got %d", e.Count())
	}
}

func TestSetLimit_UpdatesLimit(t *testing.T) {
	e := portquota.New(5)
	e.SetLimit(20)
	if e.Limit() != 20 {
		t.Errorf("expected limit 20, got %d", e.Limit())
	}
}

func TestViolation_String_ContainsValues(t *testing.T) {
	e := portquota.New(3)
	v := e.Set(9)
	if v == nil {
		t.Fatal("expected violation")
	}
	s := v.String()
	for _, sub := range []string{"9", "3", "quota exceeded"} {
		if len(s) == 0 {
			t.Fatalf("empty violation string")
		}
		// simple substring check without importing strings
		found := false
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("violation string %q missing %q", s, sub)
		}
	}
}
