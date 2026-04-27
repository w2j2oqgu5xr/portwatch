package portrank_test

import (
	"context"
	"testing"

	"github.com/yourorg/portwatch/internal/pipeline"
	"github.com/yourorg/portwatch/internal/portrank"
)

func nopHandler(_ context.Context, _ pipeline.Event) error { return nil }

func makeEvent(t pipeline.EventType, port int) pipeline.Event {
	return pipeline.Event{Type: t, Port: port}
}

func TestTrackStage_RecordsOpenedPort(t *testing.T) {
	r := portrank.New(1.0)
	r.SetWeight(80, 1.0)
	stage := portrank.TrackStage(r)

	before := r.Rank([]int{80})[0].Value
	_ = stage(context.Background(), makeEvent(pipeline.EventPortOpened, 80), nopHandler)
	after := r.Rank([]int{80})[0].Value

	if after <= before {
		t.Errorf("expected score to increase: before=%f after=%f", before, after)
	}
}

func TestTrackStage_RecordsClosedPort(t *testing.T) {
	r := portrank.New(1.0)
	r.SetWeight(443, 1.0)
	stage := portrank.TrackStage(r)

	before := r.Rank([]int{443})[0].Value
	_ = stage(context.Background(), makeEvent(pipeline.EventPortClosed, 443), nopHandler)
	after := r.Rank([]int{443})[0].Value

	if after <= before {
		t.Errorf("expected score to increase on close: before=%f after=%f", before, after)
	}
}

func TestTrackStage_ForwardsEvent(t *testing.T) {
	r := portrank.New(1.0)
	stage := portrank.TrackStage(r)
	forwarded := false
	handler := func(_ context.Context, _ pipeline.Event) error {
		forwarded = true
		return nil
	}
	_ = stage(context.Background(), makeEvent(pipeline.EventPortOpened, 22), handler)
	if !forwarded {
		t.Error("expected event to be forwarded to next handler")
	}
}

func TestTrackStage_IgnoresUnknownType(t *testing.T) {
	r := portrank.New(1.0)
	r.SetWeight(22, 1.0)
	stage := portrank.TrackStage(r)

	before := r.Rank([]int{22})[0].Value
	_ = stage(context.Background(), makeEvent(pipeline.EventType("unknown"), 22), nopHandler)
	after := r.Rank([]int{22})[0].Value

	if after != before {
		t.Errorf("expected score unchanged for unknown event type: before=%f after=%f", before, after)
	}
}
