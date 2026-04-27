package portsampler

import (
	"context"

	"github.com/user/portwatch/internal/pipeline"
)

// TrackStage returns a pipeline.Stage that records each passing event's port
// into the provided Sampler. The stage is transparent — it always forwards
// the event to the next handler regardless of whether recording succeeds.
//
// Only events whose Port field is greater than zero are recorded; control
// events with Port == 0 are forwarded without affecting the sampler.
func TrackStage(s *Sampler) pipeline.Stage {
	return func(ctx context.Context, ev pipeline.Event, next pipeline.Handler) error {
		if ev.Port > 0 {
			// Collect all currently tracked ports from the latest sample and
			// append this event's port so the sampler reflects the running set.
			latest := s.Latest()
			portSet := make(map[int]struct{}, len(latest.Ports)+1)
			for _, p := range latest.Ports {
				portSet[p] = struct{}{}
			}
			if ev.Type == pipeline.EventOpened {
				portSet[ev.Port] = struct{}{}
			} else if ev.Type == pipeline.EventClosed {
				delete(portSet, ev.Port)
			}
			ports := make([]int, 0, len(portSet))
			for p := range portSet {
				ports = append(ports, p)
			}
			s.Record(ports)
		}
		return next(ctx, ev)
	}
}
