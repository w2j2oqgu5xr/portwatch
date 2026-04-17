// Package portstate provides a thread-safe tracker for monitoring the
// open/closed state of network ports across successive scans.
//
// Use New to create a Tracker, then call Update with the latest list of
// open ports after each scan. Update returns a slice of Change values
// describing any ports that have been opened or closed since the previous
// call. The Tracker retains history for all ports it has ever seen so that
// closed ports remain visible in snapshots.
package portstate
