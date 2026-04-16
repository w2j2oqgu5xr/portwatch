// Package portname maps TCP port numbers to human-readable service names.
//
// It provides two functions:
//
//   - Lookup returns a formatted string combining the port number and service
//     name, e.g. "443 (https)". For unknown ports only the number is returned.
//
//   - Name returns only the bare service name string, or "unknown" when the
//     port is not in the built-in table.
//
// The built-in table covers the most common well-known ports. It is not
// intended to be exhaustive; use the resolver package for configurable
// overrides and extended lookups.
package portname
