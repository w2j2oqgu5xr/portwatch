// Package profiler exposes an optional pprof HTTP endpoint for runtime profiling.
package profiler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// Server wraps an HTTP server exposing pprof endpoints.
type Server struct {
	addr string
	srv  *http.Server
}

// New creates a new profiler Server bound to the given address (e.g. "localhost:6060").
func New(addr string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.DefaultServeMux)
	return &Server{
		addr: addr,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

// Start begins listening in a background goroutine.
// It returns an error if the server fails to start within 200 ms.
func (s *Server) Start() error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		return fmt.Errorf("profiler: %w", err)
	case <-time.After(200 * time.Millisecond):
		return nil
	}
}

// Shutdown gracefully stops the profiler server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// Addr returns the address the profiler is bound to.
func (s *Server) Addr() string {
	return s.addr
}
