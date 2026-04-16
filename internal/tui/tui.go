// Package tui provides a simple terminal status renderer for portwatch.
package tui

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

// Renderer writes a compact status table to an io.Writer.
type Renderer struct {
	out     io.Writer
	metrics *metrics.Metrics
}

// New returns a Renderer that writes to out.
func New(out io.Writer, m *metrics.Metrics) *Renderer {
	return &Renderer{out: out, metrics: m}
}

// Render writes the current status snapshot to the writer.
func (r *Renderer) Render(openPorts []int, lastScan time.Time) error {
	var sb strings.Builder

	sb.WriteString("╔══════════════════════════════╗\n")
	sb.WriteString("║       portwatch status       ║\n")
	sb.WriteString("╠══════════════════════════════╣\n")
	fmt.Fprintf(&sb, "║ Last scan : %-17s║\n", lastScan.Format("15:04:05"))
	fmt.Fprintf(&sb, "║ Scans     : %-17d║\n", r.metrics.Scans())
	fmt.Fprintf(&sb, "║ Alerts    : %-17d║\n", r.metrics.Alerts())
	fmt.Fprintf(&sb, "║ Suppressed: %-17d║\n", r.metrics.Suppressed())
	sb.WriteString("╠══════════════════════════════╣\n")

	sorted := make([]int, len(openPorts))
	copy(sorted, openPorts)
	sort.Ints(sorted)

	if len(sorted) == 0 {
		sb.WriteString("║ Open ports: none             ║\n")
	} else {
		parts := make([]string, len(sorted))
		for i, p := range sorted {
			parts[i] = fmt.Sprintf("%d", p)
		}
		line := strings.Join(parts, ", ")
		if len(line) > 17 {
			line = line[:14] + "..."
		}
		fmt.Fprintf(&sb, "║ Open ports: %-17s║\n", line)
	}

	sb.WriteString("╚══════════════════════════════╝\n")

	_, err := fmt.Fprint(r.out, sb.String())
	return err
}
