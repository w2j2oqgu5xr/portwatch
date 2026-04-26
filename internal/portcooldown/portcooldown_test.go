package portcooldown

import (
	"testing"
	"time"
)

func newFixed(t time.Time) *Cooldown {
	c := New(2 * time.Second)
	c.nowFunc = func() time.Time { return t }
	return c
}

func TestAllow_FirstTransitionAlwaysPermitted(t *testing.T) {
	c := newFixed(time.Now())
	if !c.Allow(80, "open") {
		t.Fatal("expected first transition to be allowed")
	}
}

func TestAllow_SameStateNotPermitted(t *testing.T) {
	c := newFixed(time.Now())
	c.Allow(80, "open")
	if c.Allow(80, "open") {
		t.Fatal("expected same-state transition to be suppressed")
	}
}

func TestAllow_SuppressedWithinCooldown(t *testing.T) {
	now := time.Now()
	c := newFixed(now)
	c.Allow(80, "open")
	// advance by less than the 2s period
	c.nowFunc = func() time.Time { return now.Add(1 * time.Second) }
	if c.Allow(80, "closed") {
		t.Fatal("expected transition to be suppressed within cooldown")
	}
}

func TestAllow_PermittedAfterCooldown(t *testing.T) {
	now := time.Now()
	c := newFixed(now)
	c.Allow(80, "open")
	c.nowFunc = func() time.Time { return now.Add(3 * time.Second) }
	if !c.Allow(80, "closed") {
		t.Fatal("expected transition to be allowed after cooldown")
	}
}

func TestAllow_IndependentPorts(t *testing.T) {
	now := time.Now()
	c := newFixed(now)
	c.Allow(80, "open")
	c.Allow(443, "open")
	// advance within cooldown for port 80
	c.nowFunc = func() time.Time { return now.Add(1 * time.Second) }
	if c.Allow(80, "closed") {
		t.Fatal("port 80 should still be suppressed")
	}
	if c.Allow(443, "closed") {
		t.Fatal("port 443 should still be suppressed")
	}
}

func TestReset_AllowsImmediateTransition(t *testing.T) {
	now := time.Now()
	c := newFixed(now)
	c.Allow(80, "open")
	c.nowFunc = func() time.Time { return now.Add(500 * time.Millisecond) }
	c.Reset(80)
	if !c.Allow(80, "closed") {
		t.Fatal("expected transition after reset to be allowed")
	}
}

func TestLen_TracksEntries(t *testing.T) {
	c := newFixed(time.Now())
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
	c.Allow(80, "open")
	c.Allow(443, "open")
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
	c.Reset(80)
	if c.Len() != 1 {
		t.Fatalf("expected 1 after reset, got %d", c.Len())
	}
}

func TestNew_DefaultCooldownApplied(t *testing.T) {
	c := New(0)
	if c.period != DefaultCooldown {
		t.Fatalf("expected default cooldown %v, got %v", DefaultCooldown, c.period)
	}
}
