// Package daemon provides the top-level runtime loop for portwatch.
//
// A Daemon is constructed from a [config.Config], wires all internal
// subsystems together, and exposes a single Run method that blocks until
// the supplied context is cancelled.
//
// Typical usage:
//
//	cfg, err := config.Load(path)
//	if err != nil { ... }
//
//	d, err := daemon.New(cfg)
//	if err != nil { ... }
//
//	if err := d.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//	    log.Fatal(err)
//	}
package daemon
