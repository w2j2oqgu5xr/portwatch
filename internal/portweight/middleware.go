package portweight

import (
	"context"

	"github.com/user/portwatch/internal/pipeline"
)

// TrackStage returns a pipeline.Stage that feeds port open/close
// events into the Scorer and then forwards the event unchanged.
func TrackStage(s *Scorer) pipeline.Stage {
	return func(ctx context.Context, event pipeline.Event, next pipeline.Handler) error {
		switch event.Type {
		case pipeline.EventPortOpened:
			s.ObserveOpen(event.Port)
			s.RecordChange(event.Port)
		case pipeline.EventPortClosed:
			s.RecordChange(event.Port)
		}
		return next(ctx, event)
	}
}
