package portstate

import (
	"context"

	"github.com/user/portwatch/internal/pipeline"
)

// TrackStage returns a pipeline.Stage that records each event's port
// transition into the provided Tracker. The event is always forwarded
// unchanged so downstream stages continue to receive it.
func TrackStage(tr *Tracker) pipeline.Stage {
	return func(ctx context.Context, ev pipeline.Event, next func(pipeline.Event)) {
		switch ev.Type {
		case pipeline.EventOpened:
			tr.mu.Lock()
			s := tr.states[ev.Port]
			s.Port = ev.Port
			s.Open = true
			s.LastSeen = ev.At
			if s.FirstSeen.IsZero() {
				s.FirstSeen = ev.At
			}
			tr.states[ev.Port] = s
			tr.mu.Unlock()
		case pipeline.EventClosed:
			tr.mu.Lock()
			s := tr.states[ev.Port]
			s.Open = false
			s.LastSeen = ev.At
			tr.states[ev.Port] = s
			tr.mu.Unlock()
		}
		next(ev)
	}
}
