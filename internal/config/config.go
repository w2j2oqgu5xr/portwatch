// Package config handles loading and validating portwatch configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	// Host is the target host to scan (default: "localhost").
	Host string `json:"host"`

	// Ports lists the port numbers to monitor.
	Ports []int `json:"ports"`

	// Interval is how often to re-scan the ports.
	Interval Duration `json:"interval"`
}

// Duration is a wrapper around time.Duration that supports JSON unmarshalling
// from a human-readable string such as "5s" or "1m".
type Duration struct{ time.Duration }

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("config: invalid duration %q: %w", s, err)
	}
	d.Duration = parsed
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Validate returns an error if the configuration contains invalid values.
func (c *Config) Validate() error {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if len(c.Ports) == 0 {
		return fmt.Errorf("config: at least one port must be specified")
	}
	for _, p := range c.Ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("config: port %d is out of valid range (1-65535)", p)
		}
	}
	if c.Interval.Duration <= 0 {
		c.Interval.Duration = 5 * time.Second
	}
	return nil
}
