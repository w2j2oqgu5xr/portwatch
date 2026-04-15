// Package daemon wires together the scanner, monitor, pipeline, and
// notifier into a long-running portwatch process.
package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/schedule"
)

// Daemon orchestrates periodic port scanning and event delivery.
type Daemon struct {
	cfg      *config.Config
	mon      *monitor.Monitor
	pipe     *pipeline.Pipeline
	sched    *schedule.Scheduler
	met      *metrics.Metrics
	notifier notify.Notifier
}

// New constructs a Daemon from the supplied configuration.
func New(cfg *config.Config) (*Daemon, error) {
	met := metrics.New()

	notifier, err := notify.Build(cfg.Notify)
	if err != nil {
		return nil, err
	}

	pipe := pipeline.New(notifier)

	mon, err := monitor.New(cfg)
	if err != nil {
		return nil, err
	}

	interval := time.Duration(cfg.IntervalSeconds) * time.Second
	sched, err := schedule.New(interval)
	if err != nil {
		return nil, err
	}

	return &Daemon{
		cfg:      cfg,
		mon:      mon,
		pipe:     pipe,
		sched:    sched,
		met:      met,
		notifier: notifier,
	}, nil
}

// Run starts the daemon and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch daemon starting (interval=%ds)\n", d.cfg.IntervalSeconds)

	ticker := d.sched.Ticker(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch daemon stopped")
			return ctx.Err()
		case <-ticker:
			events, err := d.mon.Scan(ctx)
			if err != nil {
				log.Printf("scan error: %v", err)
				continue
			}
			d.met.IncScans()
			d.met.SetOpenPorts(d.mon.OpenCount())
			for _, ev := range events {
				if err := d.pipe.Process(ctx, ev); err != nil {
					log.Printf("pipeline error: %v", err)
				}
			}
		}
	}
}

// Metrics returns the live metrics snapshot.
func (d *Daemon) Metrics() *metrics.Metrics { return d.met }
