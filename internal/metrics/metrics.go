// Package metrics tracks runtime counters for portwatch: scans performed,
// alerts fired, and events suppressed by throttle or filter.
package metrics

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// Counters holds all runtime statistics for a monitoring session.
type Counters struct {
	ScansTotal     atomic.Int64
	AlertsTotal    atomic.Int64
	Suppressed     atomic.Int64
	PortsOpen      atomic.Int64
	StartedAt      time.Time
}

// New returns a Counters instance with StartedAt set to now.
func New() *Counters {
	return &Counters{StartedAt: time.Now()}
}

// IncScans increments the total scan counter by 1.
func (c *Counters) IncScans() { c.ScansTotal.Add(1) }

// IncAlerts increments the alert counter by 1.
func (c *Counters) IncAlerts() { c.AlertsTotal.Add(1) }

// IncSuppressed increments the suppressed-event counter by 1.
func (c *Counters) IncSuppressed() { c.Suppressed.Add(1) }

// SetOpenPorts overwrites the current open-port gauge.
func (c *Counters) SetOpenPorts(n int) { c.PortsOpen.Store(int64(n)) }

// Uptime returns how long the monitor has been running.
func (c *Counters) Uptime() time.Duration {
	return time.Since(c.StartedAt).Truncate(time.Second)
}

// WriteTo formats a human-readable summary and writes it to w.
func (c *Counters) WriteTo(w io.Writer) (int64, error) {
	s := fmt.Sprintf(
		"uptime=%s scans=%d alerts=%d suppressed=%d open_ports=%d\n",
		c.Uptime(),
		c.ScansTotal.Load(),
		c.AlertsTotal.Load(),
		c.Suppressed.Load(),
		c.PortsOpen.Load(),
	)
	n, err := fmt.Fprint(w, s)
	return int64(n), err
}
