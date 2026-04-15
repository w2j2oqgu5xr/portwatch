// Package pipeline provides a composable processing pipeline for port
// change events.
//
// A Pipeline is built from an ordered list of Stage functions and a
// Notifier. Each Stage may transform or silently drop an event before
// it reaches the notifier. Two built-in stage constructors are
// provided:
//
//   - FilterStage wraps a filter.Filter to allow only ports that match
//     the configured allow/deny lists.
//
//   - ThrottleStage wraps a throttle.Throttle to suppress duplicate
//     alerts within a cooldown window.
//
// Example usage:
//
//	f := filter.New(filter.Config{AllowPorts: []int{80, 443}})
//	t := throttle.New(30 * time.Second)
//	p := pipeline.New(notifier, pipeline.FilterStage(f), pipeline.ThrottleStage(t))
//	p.Process(ctx, &pipeline.Event{Port: 80, Kind: alert.KindOpened})
package pipeline
