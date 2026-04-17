package portstate_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/portstate"
)

func makeEvent(port int, evType string) pipeline.Event {
	return pipeline.Event{Port: port, Type: evType, At: time.Now()}
}

func TestTrackStage_RecordsOpenedPort(t *testing.T) {
	tr := portstate.New()
	stage := portstate.TrackStage(tr)

	var forwarded []pipeline.Event
	stage(context.Background(), makeEvent(80, pipeline.EventOpened), func(e pipeline.Event) {
		forwarded = append(forwarded, e)
	})

	ports := tr.OpenPorts()
	if len(ports) != 1 || ports[0] != 80 {
		t.Errorf("expected port 80 to be tracked as open")
	}
	if len(forwarded) != 1 {
		t.Errorf("expected event to be forwarded")
	}
}

func TestTrackStage_RecordsClosedPort(t *testing.T) {
	tr := portstate.New()
	tr.Update([]int{443})
	stage := portstate.TrackStage(tr)

	stage(context.Background(), makeEvent(443, pipeline.EventClosed), func(e pipeline.Event) {})

	ports := tr.OpenPorts()
	for _, p := range ports {
		if p == 443 {
			t.Errorf("port 443 should not be open after closed event")
		}
	}
}

func TestTrackStage_ForwardsUnknownEventType(t *testing.T) {
	tr := portstate.New()
	stage := portstate.TrackStage(tr)

	var forwarded int
	stage(context.Background(), makeEvent(22, "unknown"), func(e pipeline.Event) {
		forwarded++
	})
	if forwarded != 1 {
		t.Errorf("expected unknown event to be forwarded")
	}
}
