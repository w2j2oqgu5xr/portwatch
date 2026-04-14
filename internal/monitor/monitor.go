package monitor

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Change represents a detected port state change.
type Change struct {
	Port   int
	Status string // "opened" or "closed"
}

// Monitor continuously scans ports and reports changes via a channel.
type Monitor struct {
	Host     string
	Ports    []int
	Interval time.Duration
	Changes  chan Change
	Stop     chan struct{}
}

// New creates a new Monitor instance.
func New(host string, ports []int, interval time.Duration) *Monitor {
	return &Monitor{
		Host:     host,
		Ports:    ports,
		Interval: interval,
		Changes:  make(chan Change, 16),
		Stop:     make(chan struct{}),
	}
}

// Start begins monitoring and sends Change events when port states differ.
func (m *Monitor) Start() {
	previous := make(map[int]bool)

	// Populate initial state silently.
	initial, _ := scanner.Scan(m.Host, m.Ports)
	for _, p := range initial {
		previous[p] = true
	}

	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.Stop:
			close(m.Changes)
			return
		case <-ticker.C:
			current, err := scanner.Scan(m.Host, m.Ports)
			if err != nil {
				fmt.Printf("[portwatch] scan error: %v\n", err)
				continue
			}

			currentSet := make(map[int]bool, len(current))
			for _, p := range current {
				currentSet[p] = true
			}

			// Detect newly opened ports.
			for _, p := range m.Ports {
				if currentSet[p] && !previous[p] {
					m.Changes <- Change{Port: p, Status: "opened"}
				} else if !currentSet[p] && previous[p] {
					m.Changes <- Change{Port: p, Status: "closed"}
				}
			}

			previous = currentSet
		}
	}
}
