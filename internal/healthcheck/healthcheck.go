// Package healthcheck provides a simple HTTP endpoint that exposes
// runtime metrics and daemon status for external monitoring systems.
package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Status represents the current health of the portwatch daemon.
type Status struct {
	OK        bool      `json:"ok"`
	Uptime    string    `json:"uptime"`
	Scans     uint64    `json:"scans_total"`
	Alerts    uint64    `json:"alerts_total"`
	OpenPorts int       `json:"open_ports"`
	CheckedAt time.Time `json:"checked_at"`
}

// Provider supplies live metric values to the health handler.
type Provider interface {
	Scans() uint64
	Alerts() uint64
	OpenPorts() int
}

// Server is a lightweight HTTP server exposing a /healthz endpoint.
type Server struct {
	addr     string
	start    time.Time
	provider Provider
	ready    atomic.Bool
	server   *http.Server
}

// New creates a Server that listens on addr and reads metrics from p.
func New(addr string, p Provider) *Server {
	s := &Server{
		addr:     addr,
		start:    time.Now(),
		provider: p,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	s.server = &http.Server{Addr: addr, Handler: mux}
	return s
}

// SetReady marks the daemon as fully initialised.
func (s *Server) SetReady(ready bool) { s.ready.Store(ready) }

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error { return s.server.ListenAndServe() }

// Shutdown gracefully stops the server.
func (s *Server) Shutdown() error {
	return s.server.Close()
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	status := Status{
		OK:        s.ready.Load(),
		Uptime:    fmt.Sprintf("%.0fs", time.Since(s.start).Seconds()),
		Scans:     s.provider.Scans(),
		Alerts:    s.provider.Alerts(),
		OpenPorts: s.provider.OpenPorts(),
		CheckedAt: time.Now().UTC(),
	}
	code := http.StatusOK
	if !status.OK {
		code = http.StatusServiceUnavailable
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(status)
}
