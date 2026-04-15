// Package cli handles command-line argument parsing for portwatch.
//
// It exposes a thin Parse function that accepts a slice of argument strings
// (typically os.Args[1:]) and an io.Writer for usage output, returning a
// populated Options struct.
//
// Example:
//
//	opts, err := cli.Parse(os.Args[1:], os.Stderr)
//	if err != nil {
//		os.Exit(2)
//	}
//	if opts.Version {
//		cli.PrintVersion(os.Stdout)
//		os.Exit(0)
//	}
package cli
