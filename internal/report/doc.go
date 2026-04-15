// Package report generates human-readable summaries of portwatch scan
// results and change history.
//
// A Summary bundles the current open-port list, the diff from the previous
// snapshot, and any recent history events into a single value.  A Writer
// renders that value as plain text suitable for terminal output or log files.
//
// Typical usage:
//
//	events, _ := history.ReadAll(historyPath)
//	diff       := snapshot.Diff(previous, current)
//	summary    := report.BuildSummary(current, diff, events)
//	w          := report.NewWriter(os.Stdout)
//	w.Write(summary)
package report
