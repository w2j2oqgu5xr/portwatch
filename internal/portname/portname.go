// Package portname provides human-readable labels for well-known port numbers.
package portname

import "fmt"

// Lookup returns a display string for the given port, e.g. "80 (http)".
func Lookup(port int) string {
	if name, ok := wellKnown[port]; ok {
		return fmt.Sprintf("%d (%s)", port, name)
	}
	return fmt.Sprintf("%d", port)
}

// Name returns only the service name for a port, or "unknown" if not found.
func Name(port int) string {
	if name, ok := wellKnown[port]; ok {
		return name
	}
	return "unknown"
}

// wellKnown is a curated map of common TCP port numbers to service names.
var wellKnown = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}
