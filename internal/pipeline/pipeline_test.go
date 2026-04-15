package pipeline_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/throttle"
)

// recordingNotifier captures every event it receives.
type recordingNotifier struct {
	mu     sync.Mutex
	events []alert.Event
}

func (r *recordingNotifier) Notify(_ context.Context, e alert.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, e)
	return nil
}

func (r *recordingNotifier) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.events)
}

func newEvent(port int) *pipeline.Event {
	return &pipeline.Event{
		Port:      port,
		Host:      "localhost",
		Kind:      alert.KindOpened,
		Timestamp: time.Now(),
	}
}

func TestProcess_ReachesNotifier(t *testing.T) {
	n := &recordingNotifier{}
	p := pipeline.New(n)
	if err := p.Process(context.Background(), newEvent(80)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.count() != 1 {
		t.Fatalf("expected 1 notification, got %d", n.count())
	}
}

func TestProcess_StageDropsEvent(t *testing.T) {
	n := &recordingNotifier{}
	drop := func(_ context.Context, _ *pipeline.Event) (*pipeline.Event, error) {
		return nil, nil
	}
	p := pipeline.New(n, drop)
	if err := p.Process(context.Background(), newEvent(80)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.count() != 0 {
		t.Fatalf("expected 0 notifications, got %d", n.count())
	}
}

func TestProcess_StageReturnsError(t *testing.T) {
	n := &recordingNotifier{}
	boom := errors.New("stage failure")
	fail := func(_ context.Context, _ *pipeline.Event) (*pipeline.Event, error) {
		return nil, boom
	}
	p := pipeline.New(n, fail)
	if err := p.Process(context.Background(), newEvent(80)); !errors.Is(err, boom) {
		t.Fatalf("expected boom error, got %v", err)
	}
}

func TestFilterStage_AllowsMatchingPort(t *testing.T) {
	n := &recordingNotifier{}
	f, _ := filter.New(filter.Config{AllowPorts: []int{443}})
	p := pipeline.New(n, pipeline.FilterStage(f))
	_ = p.Process(context.Background(), newEvent(443))
	if n.count() != 1 {
		t.Fatalf("expected 1 notification, got %d", n.count())
	}
}

func TestFilterStage_DropsNonMatchingPort(t *testing.T) {
	n := &recordingNotifier{}
	f, _ := filter.New(filter.Config{AllowPorts: []int{443}})
	p := pipeline.New(n, pipeline.FilterStage(f))
	_ = p.Process(context.Background(), newEvent(80))
	if n.count() != 0 {
		t.Fatalf("expected 0 notifications, got %d", n.count())
	}
}

func TestThrottleStage_SuppressesRepeat(t *testing.T) {
	n := &recordingNotifier{}
	th := throttle.New(10 * time.Second)
	p := pipeline.New(n, pipeline.ThrottleStage(th))
	_ = p.Process(context.Background(), newEvent(80))
	_ = p.Process(context.Background(), newEvent(80))
	if n.count() != 1 {
		t.Fatalf("expected 1 notification, got %d", n.count())
	}
}
