package porttrend

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
)

func makeEvent(port int, kind string) pipeline.Event {
	return pipeline.Event{Port: port, Type: kind}
}

func nopHandler(_ context.Context, _ pipeline.Event) error { return nil }

func TestTrackStage_RecordsSample(t *testing.T) {
	tr := New(time.Minute)
	counter := 0
	openCount := func() int {
		counter++
		return 4
	}
	stage := TrackStage(tr, openCount)
	err := stage(context.Background(), makeEvent(80, "opened"), nopHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if counter != 1 {
		t.Errorf("expected openCount called once, got %d", counter)
	}
	samples := tr.Samples()
	if len(samples) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(samples))
	}
	if samples[0].OpenCount != 4 {
		t.Errorf("expected OpenCount 4, got %d", samples[0].OpenCount)
	}
}

func TestTrackStage_ForwardsEvent(t *testing.T) {
	tr := New(time.Minute)
	var received pipeline.Event
	handler := func(_ context.Context, e pipeline.Event) error {
		received = e
		return nil
	}
	stage := TrackStage(tr, func() int { return 1 })
	evt := makeEvent(443, "closed")
	_ = stage(context.Background(), evt, handler)
	if received.Port != 443 {
		t.Errorf("expected forwarded port 443, got %d", received.Port)
	}
}

func TestTrackStage_AccumulatesTrend(t *testing.T) {
	tr := New(time.Minute)
	counts := []int{2, 5, 8}
	i := 0
	openCount := func() int {
		v := counts[i]
		i++
		return v
	}
	stage := TrackStage(tr, openCount)
	for range counts {
		_ = stage(context.Background(), makeEvent(22, "opened"), nopHandler)
	}
	if got := tr.Current(); got != TrendRising {
		t.Errorf("expected rising trend, got %v", got)
	}
}
