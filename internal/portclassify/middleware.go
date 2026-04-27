package portclassify

import (
	"context"
	"log/slog"

	"github.com/user/portwatch/internal/pipeline"
)

// Stage returns a pipeline.Stage that annotates log output with the risk tier
// of the affected port for every event that passes through.
func Stage(c *Classifier) pipeline.Stage {
	return func(ctx context.Context, event pipeline.Event, next pipeline.Handler) error {
		result := c.Classify(event.Port)
		slog.Debug("port classified",
			"port", event.Port,
			"tier", result.Tier.String(),
			"reason", result.Reason,
		)
		// Attach tier to event metadata if supported, then forward.
		if event.Meta == nil {
			event.Meta = make(map[string]string)
		}
		event.Meta["risk_tier"] = result.Tier.String()
		return next(ctx, event)
	}
}
