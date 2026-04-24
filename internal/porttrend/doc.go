// Package porttrend provides a rolling-window tracker for open-port counts,
// allowing portwatch to detect whether the number of open ports is rising,
// falling, or stable over a configurable observation window.
//
// Usage:
//
//	tr := porttrend.New(5 * time.Minute)
//	tr.Record(scanner.OpenPortCount())
//	fmt.Println(tr.Current()) // "rising", "falling", or "stable"
//
// The TrackStage helper integrates the tracker into a pipeline so that every
// event that flows through automatically records the current open-port count.
package porttrend
