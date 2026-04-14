package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
}

// ScanOptions configures the port scanner.
type ScanOptions struct {
	Host    string
	Ports   []int
	Timeout time.Duration
}

// Scan checks whether each port in opts.Ports is open on opts.Host.
// It returns a slice of PortState results.
func Scan(opts ScanOptions) ([]PortState, error) {
	if opts.Host == "" {
		return nil, fmt.Errorf("scanner: host must not be empty")
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 500 * time.Millisecond
	}

	results := make([]PortState, 0, len(opts.Ports))

	for _, port := range opts.Ports {
		if port < 1 || port > 65535 {
			return nil, fmt.Errorf("scanner: invalid port number %d", port)
		}

		address := fmt.Sprintf("%s:%d", opts.Host, port)
		conn, err := net.DialTimeout("tcp", address, opts.Timeout)

		state := PortState{
			Port:     port,
			Protocol: "tcp",
			Open:     err == nil,
		}

		if conn != nil {
			_ = conn.Close()
		}

		results = append(results, state)
	}

	return results, nil
}

// OpenPorts filters a slice of PortState and returns only the open ones.
func OpenPorts(states []PortState) []PortState {
	open := make([]PortState, 0)
	for _, s := range states {
		if s.Open {
			open = append(open, s)
		}
	}
	return open
}
