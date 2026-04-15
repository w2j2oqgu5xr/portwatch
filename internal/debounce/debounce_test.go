package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

const tick = 20 * time.Millisecond

func TestSubmit_CallsAfterWait(t *testing.T) {
	d := debounce.New(tick)
	var called int32

	d.Submit("8080", func() { atomic.AddInt32(&called, 1) })

	time.Sleep(tick * 3)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected callback to be called once, got %d", called)
	}
}

func TestSubmit_ResetsWindowOnRepeat(t *testing.T) {
	d := debounce.New(tick * 2)
	var called int32

	// Fire three times in quick succession — only the last should trigger.
	for i := 0; i < 3; i++ {
		d.Submit("9090", func() { atomic.AddInt32(&called, 1) })
		time.Sleep(tick / 2)
	}

	time.Sleep(tick * 5)
	if n := atomic.LoadInt32(&called); n != 1 {
		t.Fatalf("expected exactly 1 call after debounce, got %d", n)
	}
}

func TestSubmit_IndependentKeys(t *testing.T) {
	d := debounce.New(tick)
	var a, b int32

	d.Submit("portA", func() { atomic.AddInt32(&a, 1) })
	d.Submit("portB", func() { atomic.AddInt32(&b, 1) })

	time.Sleep(tick * 3)
	if atomic.LoadInt32(&a) != 1 || atomic.LoadInt32(&b) != 1 {
		t.Fatalf("expected both callbacks called once, got a=%d b=%d", a, b)
	}
}

func TestCancel_PreventsCallback(t *testing.T) {
	d := debounce.New(tick * 2)
	var called int32

	d.Submit("443", func() { atomic.AddInt32(&called, 1) })
	d.Cancel("443")

	time.Sleep(tick * 4)
	if atomic.LoadInt32(&called) != 0 {
		t.Fatal("expected callback to be suppressed after Cancel")
	}
}

func TestPending_TracksActiveTimers(t *testing.T) {
	d := debounce.New(tick * 3)

	if d.Pending() != 0 {
		t.Fatal("expected 0 pending timers initially")
	}

	d.Submit("22", func() {})
	d.Submit("80", func() {})

	if d.Pending() != 2 {
		t.Fatalf("expected 2 pending timers, got %d", d.Pending())
	}

	time.Sleep(tick * 6)
	if d.Pending() != 0 {
		t.Fatalf("expected 0 pending timers after expiry, got %d", d.Pending())
	}
}

func TestNew_ZeroWaitDefaulted(t *testing.T) {
	d := debounce.New(0)
	var called int32
	d.Submit("53", func() { atomic.AddInt32(&called, 1) })
	time.Sleep(300 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("expected callback to fire with defaulted wait duration")
	}
}
