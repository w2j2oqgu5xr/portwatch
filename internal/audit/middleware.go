package audit

import (
	"context"

	"github.com/user/portwatch/internal/pipeline"
)

// Stage returns a pipeline.Stage that records every event that passes
// through to the provided Logger. Events are never dropped by this stage.
func Stage(l *Logger) pipeline.Stage {
	return func(ctx context.Context, e pipeline.Event, next func(pipeline.Event)) {
		severity := "info"
		if e.Type == "opened" {
			severity = "warn"
		}
		_ = l.Record(Entry{
			Event:    e.Type,
			Port:     e.Port,
			Host:     e.Host,
			Severity: severity,
		})
		next(e)
	}
}
