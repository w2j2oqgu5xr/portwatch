package portttl_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portttl"
)

func TestRecord_StoresEntry(t *testing.T) {
	tr := portttl.New()
	tr.Record(8080, true)
	e, ok := tr.Get(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if !e.OpenState {
		t.Error("expected open state")
	}
	if e.Port != 8080 {
		t.Errorf("expected port 8080, got %d", e.Port)
	}
}

func TestRecord_SameStateDoesNotResetTimestamp(t *testing.T) {
	tr := portttl.New()
	tr.Record(443, true)
	e1, _ := tr.Get(443)
	time.Sleep(10 * time.Millisecond)
	tr.Record(443, true)
	e2, _ := tr.Get(443)
	if !e1.Since.Equal(e2.Since) {
		t.Error("timestamp should not change when state is the same")
	}
}

func TestRecord_StateChangeResetsTimestamp(t *testing.T) {
	tr := portttl.New()
	tr.Record(22, true)
	e1, _ := tr.Get(22)
	time.Sleep(10 * time.Millisecond)
	tr.Record(22, false)
	e2, _ := tr.Get(22)
	if !e2.Since.After(e1.Since) {
		t.Error("timestamp should update on state change")
	}
}

func TestAge_ReturnsPositiveDuration(t *testing.T) {
	tr := portttl.New()
	tr.Record(80, true)
	time.Sleep(5 * time.Millisecond)
	age, ok := tr.Age(80)
	if !ok {
		t.Fatal("expected age to be available")
	}
	if age <= 0 {
		t.Errorf("expected positive age, got %v", age)
	}
}

func TestAge_MissingPort(t *testing.T) {
	tr := portttl.New()
	_, ok := tr.Age(9999)
	if ok {
		t.Error("expected false for untracked port")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	tr := portttl.New()
	tr.Record(3000, true)
	tr.Delete(3000)
	_, ok := tr.Get(3000)
	if ok {
		t.Error("expected entry to be deleted")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	tr := portttl.New()
	tr.Record(1, true)
	tr.Record(2, false)
	tr.Record(3, true)
	all := tr.All()
	if len(all) != 3 {
		t.Errorf("expected 3 entries, got %d", len(all))
	}
}
