// Package watchdog implements a heartbeat-based stall detector for the
// portwatch scan loop.
//
// Usage:
//
//	beat := make(chan struct{}, 1)
//	wd := watchdog.New(beat, 10*time.Second, func() {
//		// restart or alert
//	})
//	go wd.Run(ctx)
//
//	// inside scan loop:
//	beat <- struct{}{}
package watchdog
