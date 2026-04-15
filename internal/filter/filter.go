// Package filter provides port filtering utilities for portwatch.
// It allows users to include or exclude specific ports or port ranges
// from monitoring based on configurable rules.
package filter

import "fmt"

// Rule defines a single filter rule applied to a port number.
type Rule struct {
	Allow bool
	Ports []int
}

// Filter holds a set of rules and evaluates whether a port should
// be monitored.
type Filter struct {
	allowList map[int]struct{}
	denyList  map[int]struct{}
}

// New constructs a Filter from allow and deny port lists.
// If allowList is non-empty, only those ports are monitored.
// Ports in denyList are always excluded.
func New(allow, deny []int) *Filter {
	f := &Filter{
		allowList: make(map[int]struct{}, len(allow)),
		denyList:  make(map[int]struct{}, len(deny)),
	}
	for _, p := range allow {
		f.allowList[p] = struct{}{}
	}
	for _, p := range deny {
		f.denyList[p] = struct{}{}
	}
	return f
}

// Allow reports whether the given port passes the filter.
func (f *Filter) Allow(port int) bool {
	if _, blocked := f.denyList[port]; blocked {
		return false
	}
	if len(f.allowList) == 0 {
		return true
	}
	_, ok := f.allowList[port]
	return ok
}

// Apply returns only the ports from the input slice that pass the filter.
func (f *Filter) Apply(ports []int) []int {
	out := make([]int, 0, len(ports))
	for _, p := range ports {
		if f.Allow(p) {
			out = append(out, p)
		}
	}
	return out
}

// Validate checks that port numbers are within the valid TCP/UDP range.
func Validate(ports []int) error {
	for _, p := range ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("invalid port number: %d (must be 1–65535)", p)
		}
	}
	return nil
}
