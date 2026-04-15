// Package pipeline wires together scanning, filtering, throttling,
// and notification into a single reusable processing step.
package pipeline

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/throttle"
)

// Event represents a port state change travelling through the pipeline.
type Event struct {
	Port      int
	Host      string
	Kind      alert.EventKind
	Timestamp time.Time
}

// Stage is a function that may transform or suppress an event.
// Returning (nil, nil) drops the event silently.
type Stage func(ctx context.Context, e *Event) (*Event, error)

// Pipeline executes an ordered sequence of stages and then notifies.
type Pipeline struct {
	stages   []Stage
	notifier notify.Notifier
}

// New constructs a Pipeline with the given stages and notifier.
func New(notifier notify.Notifier, stages ...Stage) *Pipeline {
	return &Pipeline{
		stages:   stages,
		notifier: notifier,
	}
}

// Process runs e through every stage in order, then dispatches a
// notification. It stops early if a stage drops the event or returns
// an error.
func (p *Pipeline) Process(ctx context.Context, e *Event) error {
	current := e
	for _, s := range p.stages {
		var err error
		current, err = s(ctx, current)
		if err != nil {
			return err
		}
		if current == nil {
			return nil // dropped
		}
	}
	return p.notifier.Notify(ctx, alert.Event{
		Port:      current.Port,
		Host:      current.Host,
		Kind:      current.Kind,
		Timestamp: current.Timestamp,
	})
}

// FilterStage returns a Stage that uses the given filter.Filter to
// allow or suppress events based on port number.
func FilterStage(f *filter.Filter) Stage {
	return func(_ context.Context, e *Event) (*Event, error) {
		if !f.Allow(e.Port) {
			return nil, nil
		}
		return e, nil
	}
}

// ThrottleStage returns a Stage that suppresses repeated events for
// the same port within the throttle cooldown window.
func ThrottleStage(t *throttle.Throttle) Stage {
	return func(_ context.Context, e *Event) (*Event, error) {
		if !t.Allow(e.Port) {
			return nil, nil
		}
		return e, nil
	}
}
