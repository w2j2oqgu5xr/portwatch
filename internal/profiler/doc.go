// Package profiler provides an optional pprof HTTP server for runtime
// profiling of portwatch. When enabled via configuration, it binds to a
// local address and exposes the standard /debug/pprof/* endpoints.
//
// Usage:
//
//	s := profiler.New("localhost:6060")
//	if err := s.Start(); err != nil {
//		log.Fatal(err)
//	}
//	defer s.Shutdown(ctx)
package profiler
