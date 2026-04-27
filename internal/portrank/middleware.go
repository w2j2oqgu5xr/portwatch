package portrank

import (
	"context"

	"github.com/yourorg/portwatch/internal/pipeline"
)

// TrackStage returns a pipeline.Stage that records port state changes
// in the Ranker so that risk scores stay up to date as events flow
// through the processing pipeline.
func TrackStage(r *Ranker) pipeline.Stage {
	return func(ctx context.Context, event pipeline.Event, next pipeline.Handler) error {
		switch event.Type {
		case pipeline.EventPortOpened, pipeline.EventPortClosed:
			r.RecordChange(event.Port)
		}
		return next(ctx, event)
	}
}
