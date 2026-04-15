// Package cli provides helpers for building the portwatch command-line
// interface, including flag parsing and usage formatting.
package cli

import (
	"flag"
	"fmt"
	"io"
)

// Options holds the parsed CLI flags.
type Options struct {
	ConfigPath string
	Verbose    bool
	Version    bool
}

const version = "0.1.0"

// Parse parses os.Args[1:] and returns Options.
// Output (usage / errors) is written to out.
func Parse(args []string, out io.Writer) (*Options, error) {
	fs := flag.NewFlagSet("portwatch", flag.ContinueOnError)
	fs.SetOutput(out)

	opts := &Options{}
	fs.StringVar(&opts.ConfigPath, "config", "portwatch.yaml", "path to config file")
	fs.BoolVar(&opts.Verbose, "verbose", false, "enable verbose logging")
	fs.BoolVar(&opts.Version, "version", false, "print version and exit")

	fs.Usage = func() {
		fmt.Fprintf(out, "Usage: portwatch [flags]\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return opts, nil
}

// PrintVersion writes the version string to out.
func PrintVersion(out io.Writer) {
	fmt.Fprintf(out, "portwatch %s\n", version)
}
