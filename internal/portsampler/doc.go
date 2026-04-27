// Package portsampler provides a sliding-window sampler for open port
// observations. It is designed to be embedded in monitoring pipelines to
// accumulate periodic snapshots of which ports are open, enabling downstream
// consumers to compute statistics such as average port count or detect
// anomalous spikes in activity.
//
// # Usage
//
//	s := portsampler.New(60) // retain last 60 samples
//	s.Record(openPorts)      // call after each scan
//	avg := s.AverageCount() // mean open-port count over the window
//
// # Pipeline integration
//
// TrackStage wraps a Sampler as a pipeline.Stage so it can be inserted into
// an existing processing chain without additional glue code.
package portsampler
