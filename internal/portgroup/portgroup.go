// Package portgroup provides named groupings of ports for organized monitoring.
package portgroup

import "fmt"

// Group represents a named collection of ports.
type Group struct {
	Name  string
	Ports []int
}

// Registry holds named port groups.
type Registry struct {
	groups map[string]Group
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{groups: make(map[string]Group)}
}

// Add registers a named group. Returns an error if the name is already taken.
func (r *Registry) Add(name string, ports []int) error {
	if _, exists := r.groups[name]; exists {
		return fmt.Errorf("portgroup: group %q already registered", name)
	}
	if len(ports) == 0 {
		return fmt.Errorf("portgroup: group %q must contain at least one port", name)
	}
	copy := make([]int, len(ports))
	for i, p := range ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("portgroup: invalid port %d in group %q", p, name)
		}
		copy[i] = p
	}
	r.groups[name] = Group{Name: name, Ports: copy}
	return nil
}

// Get returns the Group for the given name, or false if not found.
func (r *Registry) Get(name string) (Group, bool) {
	g, ok := r.groups[name]
	return g, ok
}

// All returns all registered groups.
func (r *Registry) All() []Group {
	out := make([]Group, 0, len(r.groups))
	for _, g := range r.groups {
		out = append(out, g)
	}
	return out
}

// Resolve expands a list of group names into a deduplicated port slice.
func (r *Registry) Resolve(names []string) ([]int, error) {
	seen := make(map[int]struct{})
	var ports []int
	for _, name := range names {
		g, ok := r.groups[name]
		if !ok {
			return nil, fmt.Errorf("portgroup: unknown group %q", name)
		}
		for _, p := range g.Ports {
			if _, dup := seen[p]; !dup {
				seen[p] = struct{}{}
				ports = append(ports, p)
			}
		}
	}
	return ports, nil
}
