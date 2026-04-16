// Package portgroup provides named groupings of ports to simplify
// monitor configuration. Users can define logical groups such as
// "web", "database", or "mail" and reference them by name rather
// than repeating port lists throughout configuration files.
//
// A Registry holds all named groups and can resolve one or more
// group names into a deduplicated list of port numbers ready for
// use with the scanner.
//
// Default well-known groups can be loaded via LoadDefaults.
package portgroup
