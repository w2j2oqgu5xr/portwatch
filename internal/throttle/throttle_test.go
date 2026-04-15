package throttle_test

import (
	"sync"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/throttle"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	th := throttle.New(time.Second)
	if !th.Allow(8080) {
		t.Fatal("expected first alert for port 8080 to be allowed")
	}
}

func TestAllow_SuppressedWithinCooldown(t *testing.T) {
	th := throttle.New(time.Minute)
	th.Allow(9090) // prime
	if th.Allow(9090) {
		t.Fatal("expected second alert within cooldown to be suppressed")
	}
}

func TestAllow_PermittedAfterCooldown(t *testing.T) {
	th := throttle.New(20 * time.Millisecond)
	th.Allow(443)
	time.Sleep(30 * time.Millisecond)
	if !th.Allow(443) {
		t.Fatal("expected alert to be allowed after cooldown expires")
	}
}

func TestAllow_IndependentPorts(t *testing.T) {
	th := throttle.New(time.Minute)
	th.Allow(80)
	if !th.Allow(443) {
		t.Fatal("expected different port to be allowed independently")
	}
}

func TestReset_ClearsPort(t *testing.T) {
	th := throttle.New(time.Minute)
	th.Allow(22)
	th.Reset(22)
	if !th.Allow(22) {
		t.Fatal("expected port to be allowed after explicit Reset")
	}
}

func TestResetAll_ClearsAllPorts(t *testing.T) {
	th := throttle.New(time.Minute)
	th.Allow(80)
	th.Allow(443)
	th.ResetAll()
	if !th.Allow(80) || !th.Allow(443) {
		t.Fatal("expected all ports to be allowed after ResetAll")
	}
}

func TestAllow_ZeroCooldownDefaultsTo30s(t *testing.T) {
	th := throttle.New(0)
	if th.Cooldown() != 30*time.Second {
		t.Fatalf("expected default cooldown 30s, got %v", th.Cooldown())
	}
}

func TestAllow_ConcurrentSafe(t *testing.T) {
	th := throttle.New(time.Millisecond)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			th.Allow(p % 10)
		}(i)
	}
	wg.Wait()
}
