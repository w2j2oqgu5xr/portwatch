package portsampler_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/portsampler"
)

func nopHandler(_ context.Context, _ pipeline.Event) error { return nil }

func makeEvent(typ pipeline.EventType, port int) pipeline.Event {
	return pipeline.Event{Type: typ, Port: port}
}

func TestTrackStage_RecordsOpenedPort(t *testing.T) {
	s := portsampler.New(10)
	stage := portsampler.TrackStage(s)

	ev := makeEvent(pipeline.EventOpened, 80)
	if err := stage(context.Background(), ev, nopHandler); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	latest := s.Latest()
	if len(latest.Ports) != 1 || latest.Ports[0] != 80 {
		t.Fatalf("expected [80], got %v", latest.Ports)
	}
}

func TestTrackStage_RemovesClosedPort(t *testing.T) {
	s := portsampler.New(10)
	s.Record([]int{80, 443})
	stage := portsampler.TrackStage(s)

	ev := makeEvent(pipeline.EventClosed, 80)
	if err := stage(context.Background(), ev, nopHandler); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	latest := s.Latest()
	for _, p := range latest.Ports {
		if p == 80 {
			t.Fatal("port 80 should have been removed after close event")
		}
	}
}

func TestTrackStage_ForwardsEvent(t *testing.T) {
	s := portsampler.New(5)
	stage := portsampler.TrackStage(s)

	var received pipeline.Event
	handler := func(_ context.Context, ev pipeline.Event) error {
		received = ev
		return nil
	}

	ev := makeEvent(pipeline.EventOpened, 443)
	_ = stage(context.Background(), ev, handler)

	if received.Port != 443 {
		t.Fatalf("expected forwarded port 443, got %d", received.Port)
	}
}

func TestTrackStage_SkipsZeroPort(t *testing.T) {
	s := portsampler.New(5)
	stage := portsampler.TrackStage(s)

	ev := makeEvent(pipeline.EventOpened, 0)
	_ = stage(context.Background(), ev, nopHandler)

	if len(s.All()) != 0 {
		t.Fatal("zero-port event should not produce a sample")
	}
}
