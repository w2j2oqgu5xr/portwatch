// Package schedule provides a context-aware, interval-based Ticker for
// driving periodic port scans inside the portwatch monitor loop.
//
// Usage:
//
//	ticker, err := schedule.New(ctx, 5*time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer ticker.Stop()
//
//	for range ticker.C {
//		// perform scan
//	}
//
// The minimum allowed interval is schedule.MinInterval (500 ms) to
// prevent accidental resource exhaustion during misconfiguration.
package schedule
