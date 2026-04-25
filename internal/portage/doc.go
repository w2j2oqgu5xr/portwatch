// Package portage classifies open ports by how long they have been observed.
//
// A Tracker records the first time each port is seen open and assigns it one
// of three statuses:
//
//   - new    – seen for less than the configured new threshold (default 5 min)
//   - stable – seen for between the new and stale thresholds
//   - stale  – seen for longer than the stale threshold (default 24 h)
//
// Typical usage:
//
//	tr := portage.New(0, 0) // use package defaults
//	tr.Observe(8080)
//	status := tr.Classify(8080) // portage.StatusNew
//
// When a port closes, call Forget so the tracker does not accumulate stale
// entries indefinitely.
package portage
