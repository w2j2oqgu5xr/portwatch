package portwatch

import (
	"errors"
	"time"
)

// PortRange defines the inclusive port range to scan.
type PortRange struct {
	From int
	To   int
}

// Config holds watcher configuration.
type Config struct {
	Host     string
	Ports    PortRange
	Interval time.Duration
}

const (
	minInterval  = 5 * time.Second
	defaultFrom  = 1
	defaultTo    = 1024
)

// Validate returns an error if the Config is invalid and applies defaults
// where values are missing.
func (c *Config) Validate() error {
	if c.Host == "" {
		return errors.New("portwatch: host must not be empty")
	}
	if c.Ports.From <= 0 {
		c.Ports.From = defaultFrom
	}
	if c.Ports.To <= 0 {
		c.Ports.To = defaultTo
	}
	if c.Ports.From > c.Ports.To {
		return errors.New("portwatch: port range From must be <= To")
	}
	if c.Interval < minInterval {
		c.Interval = minInterval
	}
	return nil
}
