// Package config handles loading and validating portwatch configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Defaults applied when values are absent in the config file.
const (
	DefaultInterval   = 5 * time.Second
	DefaultHost       = "localhost"
	DefaultSnapshotDB = "portwatch.snapshot"
)

// Config holds the full portwatch runtime configuration.
type Config struct {
	Host        string        `yaml:"host"`
	Ports       []int         `yaml:"ports"`
	AllowPorts  []int         `yaml:"allow_ports"`
	DenyPorts   []int         `yaml:"deny_ports"`
	Interval    time.Duration `yaml:"interval"`
	SnapshotDB  string        `yaml:"snapshot_db"`
	HistoryFile string        `yaml:"history_file"`
	Verbose     bool          `yaml:"verbose"`
}

// Load reads a YAML config file from path and applies defaults.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}
	applyDefaults(&cfg)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func applyDefaults(c *Config) {
	if c.Host == "" {
		c.Host = DefaultHost
	}
	if c.Interval == 0 {
		c.Interval = DefaultInterval
	}
	if c.SnapshotDB == "" {
		c.SnapshotDB = DefaultSnapshotDB
	}
}

// Validate checks that the configuration is semantically valid.
func (c *Config) Validate() error {
	if len(c.Ports) == 0 && len(c.AllowPorts) == 0 {
		return errors.New("config: at least one port must be specified in 'ports' or 'allow_ports'")
	}
	for _, p := range append(c.Ports, append(c.AllowPorts, c.DenyPorts...)...) {
		if p < 1 || p > 65535 {
			return fmt.Errorf("config: invalid port %d", p)
		}
	}
	if c.Interval < time.Second {
		return fmt.Errorf("config: interval %v is too short (minimum 1s)", c.Interval)
	}
	return nil
}
