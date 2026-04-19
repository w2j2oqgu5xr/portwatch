// Package portwatch provides a high-level Watcher that combines port
// scanning, snapshot diffing, and event pipeline processing into a single
// reusable component.
//
// Basic usage:
//
//	cfg := portwatch.Config{
//		Host:     "localhost",
//		Ports:    portwatch.PortRange{From: 1, To: 1024},
//		Interval: 30 * time.Second,
//	}
//	w, err := portwatch.New(cfg, myPipeline)
//	if err != nil { ... }
//	w.Run(ctx)
package portwatch
