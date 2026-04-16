package baseline

import (
	"fmt"
	"io"
)

// Violation describes a deviation from the baseline.
type Violation struct {
	Port int
	Kind string // "unexpected" or "missing"
}

func (v Violation) String() string {
	return fmt.Sprintf("port %d is %s", v.Port, v.Kind)
}

// Check compares open ports against the baseline and returns all violations.
func Check(b *Baseline, open []int) []Violation {
	var violations []Violation
	for _, p := range b.Unexpected(open) {
		violations = append(violations, Violation{Port: p, Kind: "unexpected"})
	}
	for _, p := range b.Missing(open) {
		violations = append(violations, Violation{Port: p, Kind: "missing"})
	}
	return violations
}

// Report writes a human-readable violations report to w.
func Report(w io.Writer, violations []Violation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "baseline: no violations detected")
		return
	}
	fmt.Fprintf(w, "baseline: %d violation(s) detected\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(w, "  - %s\n", v)
	}
}
