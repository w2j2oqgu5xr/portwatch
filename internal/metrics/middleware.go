package metrics

import (
	"context"

	"github.com/user/portwatch/internal/pipeline"
)

// TrackAlertsStage returns a pipeline.Stage that increments c.AlertsTotal
// for every event that passes through, then forwards the event unchanged.
func TrackAlertsStage(c *Counters) pipeline.Stage {
	return func(ctx context.Context, ev pipeline.Event, next func(pipeline.Event)) {
		c.IncAlerts()
		next(ev)
	}
}

// TrackSuppressedStage returns a pipeline.Stage that increments c.Suppressed
// when a downstream stage drops the event (i.e. next is not called).
// It wraps the provided inner stage to observe its behaviour.
func TrackSuppressedStage(c *Counters, inner pipeline.Stage) pipeline.Stage {
	return func(ctx context.Context, ev pipeline.Event, next func(pipeline.Event)) {
		called := false
		inner(ctx, ev, func(forwarded pipeline.Event) {
			called = true
			next(forwarded)
		})
		if !called {
			c.IncSuppressed()
		}
	}
}
