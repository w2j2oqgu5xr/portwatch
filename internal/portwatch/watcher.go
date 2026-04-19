// Package portwatch ties together scanning, diffing, and alerting into a
// single high-level watch loop that external callers can embed.
package portwatch

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Event kinds emitted by the watcher.
const (
	EventOpened = "port.opened"
	EventClosed = "port.closed"
)

// Event describes a single port-state change detected by the watcher.
type Event struct {
	Kind      string
	Host      string
	Port      int
	Timestamp time.Time
}

// Watcher continuously scans a host and emits events through a pipeline.
type Watcher struct {
	cfg Config
	pipe *pipeline.Pipeline
}

// New creates a Watcher with the supplied Config and Pipeline.
func New(cfg Config, p *pipeline.Pipeline) (*Watcher, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Watcher{cfg: cfg, pipe: p}, nil
}

// Run starts the watch loop and blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	var prev []int

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			current, err := scanner.Scan(w.cfg.Host, scanner.PortRange{
				From: w.cfg.Ports.From,
				To:   w.cfg.Ports.To,
			})
			if err != nil {
				continue
			}
			diff := snapshot.Diff(prev, current)
			for _, p := range diff.Opened {
				w.pipe.Process(ctx, toEvent(EventOpened, w.cfg.Host, p))
			}
			for _, p := range diff.Closed {
				w.pipe.Process(ctx, toEvent(EventClosed, w.cfg.Host, p))
			}
			prev = current
		}
	}
}

func toEvent(kind, host string, port int) pipeline.Event {
	return pipeline.Event{
		Kind:      kind,
		Host:      host,
		Port:      port,
		Timestamp: time.Now(),
	}
}
