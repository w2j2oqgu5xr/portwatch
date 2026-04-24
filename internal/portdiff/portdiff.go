// Package portdiff computes human-readable diffs between two port sets,
// annotating each change with its resolved service name and direction.
package portdiff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/portwatch/internal/resolver"
)

// Change describes a single port transition.
type Change struct {
	Port      int
	Service   string
	Direction string // "opened" or "closed"
}

// String returns a compact, human-readable representation of the change.
func (c Change) String() string {
	return fmt.Sprintf("%s %d/%s", c.Direction, c.Port, c.Service)
}

// Diff holds the result of comparing two port snapshots.
type Diff struct {
	Opened []Change
	Closed []Change
}

// IsEmpty reports whether there are no changes.
func (d Diff) IsEmpty() bool {
	return len(d.Opened) == 0 && len(d.Closed) == 0
}

// Summary returns a single-line description suitable for log output.
func (d Diff) Summary() string {
	if d.IsEmpty() {
		return "no changes"
	}
	parts := make([]string, 0, 2)
	if n := len(d.Opened); n > 0 {
		parts = append(parts, fmt.Sprintf("%d opened", n))
	}
	if n := len(d.Closed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d closed", n))
	}
	return strings.Join(parts, ", ")
}

// Computer computes annotated diffs between port sets.
type Computer struct {
	res *resolver.Resolver
}

// New returns a Computer that annotates ports using res.
func New(res *resolver.Resolver) *Computer {
	return &Computer{res: res}
}

// Compute returns the annotated diff between prev and next port lists.
func (c *Computer) Compute(prev, next []int) Diff {
	prevSet := toSet(prev)
	nextSet := toSet(next)

	var d Diff
	for port := range nextSet {
		if !prevSet[port] {
			d.Opened = append(d.Opened, Change{
				Port:      port,
				Service:   c.res.Name(port),
				Direction: "opened",
			})
		}
	}
	for port := range prevSet {
		if !nextSet[port] {
			d.Closed = append(d.Closed, Change{
				Port:      port,
				Service:   c.res.Name(port),
				Direction: "closed",
			})
		}
	}
	sort.Slice(d.Opened, func(i, j int) bool { return d.Opened[i].Port < d.Opened[j].Port })
	sort.Slice(d.Closed, func(i, j int) bool { return d.Closed[i].Port < d.Closed[j].Port })
	return d
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
