// Package resolver maps port numbers to well-known service names.
package resolver

import (
	"fmt"
	"strconv"
)

// wellKnown maps common port numbers to service names.
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
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Resolver resolves port numbers to human-readable service names.
type Resolver struct {
	extra map[int]string
}

// New returns a Resolver with optional extra mappings merged over defaults.
func New(extra map[int]string) *Resolver {
	r := &Resolver{extra: make(map[int]string)}
	for k, v := range extra {
		r.extra[k] = v
	}
	return r
}

// Name returns the service name for port, or a numeric string if unknown.
func (r *Resolver) Name(port int) string {
	if v, ok := r.extra[port]; ok {
		return v
	}
	if v, ok := wellKnown[port]; ok {
		return v
	}
	return strconv.Itoa(port)
}

// Label returns a formatted label like "80/http" or "9999/9999".
func (r *Resolver) Label(port int) string {
	return fmt.Sprintf("%d/%s", port, r.Name(port))
}
