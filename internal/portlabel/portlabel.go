// Package portlabel assigns human-readable labels to ports based on
// their state, service name, and optional custom annotations.
package portlabel

import "fmt"

// Label holds display metadata for a single port.
type Label struct {
	Port        int
	Service     string
	Annotation  string
	Open        bool
}

// String returns a formatted label string.
func (l Label) String() string {
	status := "closed"
	if l.Open {
		status = "open"
	}
	if l.Annotation != "" {
		return fmt.Sprintf("%d/%s (%s) [%s]", l.Port, l.Service, l.Annotation, status)
	}
	return fmt.Sprintf("%d/%s [%s]", l.Port, l.Service, status)
}

// Labeler produces Labels for ports.
type Labeler struct {
	resolver    func(int) string
	annotations map[int]string
}

// New returns a Labeler. resolver maps port numbers to service names.
func New(resolver func(int) string, annotations map[int]string) *Labeler {
	if annotations == nil {
		annotations = map[int]string{}
	}
	return &Labeler{resolver: resolver, annotations: annotations}
}

// Label builds a Label for the given port and open state.
func (l *Labeler) Label(port int, open bool) Label {
	svc := l.resolver(port)
	if svc == "" {
		svc = "unknown"
	}
	return Label{
		Port:       port,
		Service:    svc,
		Annotation: l.annotations[port],
		Open:       open,
	}
}

// LabelAll returns labels for all provided ports.
func (l *Labeler) LabelAll(ports []int, open bool) []Label {
	out := make([]Label, len(ports))
	for i, p := range ports {
		out[i] = l.Label(p, open)
	}
	return out
}
