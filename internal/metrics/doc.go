// Package metrics provides lightweight, thread-safe runtime counters for
// portwatch monitoring sessions.
//
// Usage:
//
//	ctr := metrics.New()
//	ctr.IncScans()
//	ctr.IncAlerts()
//	ctr.SetOpenPorts(3)
//	ctr.WriteTo(os.Stdout)
//
// All counter operations are safe for concurrent use without additional
// locking. Counters are never reset during a session; consumers should
// compute deltas if rate information is needed.
package metrics
