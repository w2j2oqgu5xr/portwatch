package porttrend

import (
	"context"

	"github.com/user/portwatch/internal/pipeline"
)

// TrackStage returns a pipeline.Stage that records the current open-port count
// into the given Tracker each time an event passes through.
// The event is always forwarded unchanged; tracking is a side-effect.
func TrackStage(tr *Tracker, openCount func() int) pipeline.Stage {
	return func(ctx context.Context, event pipeline.Event, next pipeline.Handler) error {
		tr.Record(openCount())
		return next(ctx, event)
	}
}
