// Package alert defines the alerting interface and built-in notifiers
// used by portwatch to report port state changes.
//
// # Notifiers
//
// A Notifier sends an alert for a given Event at a specified Level.
// The package ships with a ConsoleNotifier that writes human-readable
// messages to any io.Writer.
//
// # Levels
//
// Three alert levels are defined:
//
//   - INFO  – no meaningful state change
//   - WARN  – a previously open port has closed
//   - ALERT – a new port has unexpectedly opened
//
// Use LevelForEvent to automatically derive the appropriate level from
// an Event's previous and new state.
package alert
