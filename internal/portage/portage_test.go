package portage

import (
	"testing"
	"time"
)

func TestNew_DefaultThresholds(t *testing.T) {
	tr := New(0, 0)
	if tr.newThreshold != DefaultNewThreshold {
		t.Fatalf("expected default new threshold %v, got %v", DefaultNewThreshold, tr.newThreshold)
	}
	if tr.staleThreshold != DefaultStaleThreshold {
		t.Fatalf("expected default stale threshold %v, got %v", DefaultStaleThreshold, tr.staleThreshold)
	}
}

func TestObserve_RecordsFirstSeen(t *testing.T) {
	tr := New(time.Minute, time.Hour)
	tr.Observe(8080)
	age, ok := tr.Age(8080)
	if !ok {
		t.Fatal("expected port to be tracked")
	}
	if age < 0 {
		t.Fatalf("unexpected negative age: %v", age)
	}
}

func TestObserve_IdempotentForSamePort(t *testing.T) {
	tr := New(time.Minute, time.Hour)
	tr.Observe(443)
	age1, _ := tr.Age(443)
	time.Sleep(5 * time.Millisecond)
	tr.Observe(443) // second call must not update timestamp
	age2, _ := tr.Age(443)
	if age2 < age1 {
		t.Fatal("second Observe should not reset first-seen timestamp")
	}
}

func TestForget_RemovesPort(t *testing.T) {
	tr := New(time.Minute, time.Hour)
	tr.Observe(22)
	tr.Forget(22)
	_, ok := tr.Age(22)
	if ok {
		t.Fatal("port should have been removed after Forget")
	}
}

func TestClassify_New(t *testing.T) {
	tr := New(time.Hour, 2*time.Hour)
	tr.Observe(9000)
	// age is ~0, well below newThreshold
	if got := tr.Classify(9000); got != StatusNew {
		t.Fatalf("expected StatusNew, got %s", got)
	}
}

func TestClassify_Stable(t *testing.T) {
	tr := New(time.Millisecond, time.Hour)
	tr.Observe(9001)
	time.Sleep(5 * time.Millisecond)
	if got := tr.Classify(9001); got != StatusStable {
		t.Fatalf("expected StatusStable, got %s", got)
	}
}

func TestClassify_Stale(t *testing.T) {
	tr := New(time.Millisecond, 2*time.Millisecond)
	tr.Observe(9002)
	time.Sleep(10 * time.Millisecond)
	if got := tr.Classify(9002); got != StatusStale {
		t.Fatalf("expected StatusStale, got %s", got)
	}
}

func TestClassify_UnknownPortDefaultsToNew(t *testing.T) {
	tr := New(time.Minute, time.Hour)
	if got := tr.Classify(1234); got != StatusNew {
		t.Fatalf("unknown port should default to StatusNew, got %s", got)
	}
}

func TestAge_MissingPort(t *testing.T) {
	tr := New(time.Minute, time.Hour)
	_, ok := tr.Age(5555)
	if ok {
		t.Fatal("expected ok=false for untracked port")
	}
}
