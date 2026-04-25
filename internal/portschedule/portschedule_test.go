package portschedule_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portschedule"
)

func makeTime(hour int) time.Time {
	return time.Date(2024, 1, 15, hour, 30, 0, 0, time.UTC)
}

func TestNew_NoWindows(t *testing.T) {
	s, err := portschedule.New(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Active(makeTime(3)) {
		t.Error("expected schedule with no windows to always be active")
	}
}

func TestNew_InvalidStartHour(t *testing.T) {
	_, err := portschedule.New([]portschedule.Window{{StartHour: -1, EndHour: 8}}, nil)
	if err == nil {
		t.Fatal("expected error for invalid start hour")
	}
}

func TestNew_InvalidEndHour(t *testing.T) {
	_, err := portschedule.New([]portschedule.Window{{StartHour: 8, EndHour: 24}}, nil)
	if err == nil {
		t.Fatal("expected error for invalid end hour")
	}
}

func TestNew_SameStartAndEnd(t *testing.T) {
	_, err := portschedule.New([]portschedule.Window{{StartHour: 9, EndHour: 9}}, nil)
	if err == nil {
		t.Fatal("expected error when start equals end")
	}
}

func TestActive_WithinWindow(t *testing.T) {
	s, _ := portschedule.New([]portschedule.Window{{StartHour: 8, EndHour: 18}}, nil)
	if !s.Active(makeTime(10)) {
		t.Error("expected hour 10 to be active in window 08-18")
	}
}

func TestActive_OutsideWindow(t *testing.T) {
	s, _ := portschedule.New([]portschedule.Window{{StartHour: 8, EndHour: 18}}, nil)
	if s.Active(makeTime(20)) {
		t.Error("expected hour 20 to be inactive in window 08-18")
	}
}

func TestActive_OvernightWindow(t *testing.T) {
	s, _ := portschedule.New([]portschedule.Window{{StartHour: 22, EndHour: 6}}, nil)
	if !s.Active(makeTime(23)) {
		t.Error("expected hour 23 active in overnight window 22-06")
	}
	if !s.Active(makeTime(2)) {
		t.Error("expected hour 2 active in overnight window 22-06")
	}
	if s.Active(makeTime(10)) {
		t.Error("expected hour 10 inactive in overnight window 22-06")
	}
}

func TestActive_MultipleWindows(t *testing.T) {
	s, _ := portschedule.New([]portschedule.Window{
		{StartHour: 6, EndHour: 9},
		{StartHour: 17, EndHour: 20},
	}, nil)
	if !s.Active(makeTime(7)) {
		t.Error("expected hour 7 active")
	}
	if !s.Active(makeTime(18)) {
		t.Error("expected hour 18 active")
	}
	if s.Active(makeTime(12)) {
		t.Error("expected hour 12 inactive")
	}
}

func TestNextActive_AlreadyActive(t *testing.T) {
	s, _ := portschedule.New([]portschedule.Window{{StartHour: 8, EndHour: 18}}, nil)
	now := makeTime(10)
	if got := s.NextActive(now); !got.Equal(now) {
		t.Errorf("expected NextActive to return input time when already active, got %v", got)
	}
}

func TestNextActive_AdvancesToWindow(t *testing.T) {
	s, _ := portschedule.New([]portschedule.Window{{StartHour: 8, EndHour: 18}}, nil)
	now := makeTime(20) // outside window
	next := s.NextActive(now)
	if next.Hour() != 8 {
		t.Errorf("expected next active hour to be 8, got %d", next.Hour())
	}
}

func TestNextActive_NoWindows(t *testing.T) {
	s, _ := portschedule.New(nil, nil)
	now := makeTime(3)
	if got := s.NextActive(now); !got.Equal(now) {
		t.Errorf("expected NextActive to return input when no windows set")
	}
}
