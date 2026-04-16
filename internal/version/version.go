// Package version exposes build-time version information for portwatch.
package version

import "fmt"

// These variables are set at build time via -ldflags.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// Info holds structured version metadata.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
}

// Get returns the current version Info.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
	}
}

// String returns a human-readable version string.
func (i Info) String() string {
	return fmt.Sprintf("portwatch %s (commit=%s, built=%s)", i.Version, i.Commit, i.BuildDate)
}
