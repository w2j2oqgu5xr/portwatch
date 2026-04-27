// Package portrank scores and ranks monitored ports by risk.
//
// A Ranker combines a static classification weight (set via SetWeight)
// with a dynamic change-frequency counter (incremented via RecordChange)
// to produce a floating-point risk score for each port.  Ports with
// higher weights or more frequent state transitions receive a higher
// score and therefore a lower (more urgent) rank number.
//
// Typical usage:
//
//	r := portrank.New(1.0)
//	r.SetWeight(443, 5.0)  // HTTPS is high-value
//	r.SetWeight(22, 3.0)   // SSH is medium-value
//	scores := r.Rank(openPorts)
//
// The TrackStage helper integrates the Ranker into a pipeline so that
// scores are updated automatically as port-opened / port-closed events
// flow through the processing chain.
package portrank
