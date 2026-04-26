package portlock

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestLock_And_IsLocked(t *testing.T) {
	r := New()
	if err := r.Lock(8080, "maintenance", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.IsLocked(8080) {
		t.Fatal("expected port 8080 to be locked")
	}
}

func TestUnlock_RemovesLock(t *testing.T) {
	r := New()
	_ = r.Lock(443, "test", 0)
	r.Unlock(443)
	if r.IsLocked(443) {
		t.Fatal("expected port 443 to be unlocked")
	}
}

func TestIsLocked_UnknownPort(t *testing.T) {
	r := New()
	if r.IsLocked(9999) {
		t.Fatal("expected unknown port to not be locked")
	}
}

func TestLock_InvalidPort(t *testing.T) {
	r := New()
	if err := r.Lock(0, "bad", 0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Lock(70000, "bad", 0); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestLock_ExpiresAfterTTL(t *testing.T) {
	base := time.Now()
	r := New()
	r.now = fixedNow(base)
	_ = r.Lock(22, "ttl-test", 5*time.Minute)

	// before expiry
	if !r.IsLocked(22) {
		t.Fatal("expected port to be locked before TTL")
	}

	// after expiry
	r.now = fixedNow(base.Add(10 * time.Minute))
	if r.IsLocked(22) {
		t.Fatal("expected port to be unlocked after TTL")
	}
}

func TestAll_ExcludesExpired(t *testing.T) {
	base := time.Now()
	r := New()
	r.now = fixedNow(base)
	_ = r.Lock(80, "permanent", 0)
	_ = r.Lock(81, "short", time.Minute)

	r.now = fixedNow(base.Add(2 * time.Minute))
	all := r.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 active lock, got %d", len(all))
	}
	if all[0].Port != 80 {
		t.Fatalf("expected port 80, got %d", all[0].Port)
	}
}

func TestPurge_RemovesExpired(t *testing.T) {
	base := time.Now()
	r := New()
	r.now = fixedNow(base)
	_ = r.Lock(100, "a", time.Minute)
	_ = r.Lock(101, "b", time.Minute)
	_ = r.Lock(102, "c", 0)

	r.now = fixedNow(base.Add(2 * time.Minute))
	n := r.Purge()
	if n != 2 {
		t.Fatalf("expected 2 purged, got %d", n)
	}
	if len(r.All()) != 1 {
		t.Fatalf("expected 1 remaining lock after purge")
	}
}

func TestLock_String_NoExpiry(t *testing.T) {
	l := Lock{Port: 22, Reason: "ssh", LockedAt: time.Now()}
	s := l.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
