// Package report provides functionality for generating human-readable
// summaries of port scan results and change history.
package report

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/snapshot"
)

// Summary holds aggregated data for a single report.
type Summary struct {
	GeneratedAt time.Time
	OpenPorts   []int
	Opened      []int
	Closed      []int
	Events      []history.Event
}

// Writer renders a Summary to an io.Writer.
type Writer struct {
	out io.Writer
}

// NewWriter returns a Writer that writes to out.
func NewWriter(out io.Writer) *Writer {
	return &Writer{out: out}
}

// Write renders the summary as plain text.
func (w *Writer) Write(s Summary) error {
	fmt.Fprintf(w.out, "portwatch report — %s\n", s.GeneratedAt.Format(time.RFC1123))
	fmt.Fprintln(w.out, strings.Repeat("-", 40))

	fmt.Fprintf(w.out, "Open ports (%d): ", len(s.OpenPorts))
	if len(s.OpenPorts) == 0 {
		fmt.Fprintln(w.out, "none")
	} else {
		parts := make([]string, len(s.OpenPorts))
		for i, p := range s.OpenPorts {
			parts[i] = fmt.Sprintf("%d", p)
		}
		fmt.Fprintln(w.out, strings.Join(parts, ", "))
	}

	fmt.Fprintf(w.out, "Newly opened : %s\n", formatPorts(s.Opened))
	fmt.Fprintf(w.out, "Newly closed : %s\n", formatPorts(s.Closed))

	if len(s.Events) > 0 {
		fmt.Fprintln(w.out, "\nRecent events:")
		for _, e := range s.Events {
			fmt.Fprintf(w.out, "  [%s] port %d %s\n",
				e.Timestamp.Format("15:04:05"), e.Port, e.Kind)
		}
	}
	return nil
}

// BuildSummary constructs a Summary from current ports, a diff and event log.
func BuildSummary(open []int, diff snapshot.Diff, events []history.Event) Summary {
	return Summary{
		GeneratedAt: time.Now(),
		OpenPorts:   open,
		Opened:      diff.Opened,
		Closed:      diff.Closed,
		Events:      events,
	}
}

func formatPorts(ports []int) string {
	if len(ports) == 0 {
		return "none"
	}
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(parts, ", ")
}
