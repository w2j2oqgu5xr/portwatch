package baseline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func TestCheck_DetectsUnexpectedAndMissing(t *testing.T) {
	b := &baseline.Baseline{Ports: []int{22, 80}}
	violations := baseline.Check(b, []int{80, 8080})
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d: %v", len(violations), violations)
	}
	kinds := map[string]bool{}
	for _, v := range violations {
		kinds[v.Kind] = true
	}
	if !kinds["unexpected"] || !kinds["missing"] {
		t.Errorf("expected both kinds, got %v", kinds)
	}
}

func TestCheck_NoViolations(t *testing.T) {
	b := &baseline.Baseline{Ports: []int{22, 80}}
	violations := baseline.Check(b, []int{22, 80})
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestReport_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	baseline.Report(&buf, nil)
	if !strings.Contains(buf.String(), "no violations") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestReport_ListsViolations(t *testing.T) {
	var buf bytes.Buffer
	violations := []baseline.Violation{
		{Port: 8080, Kind: "unexpected"},
		{Port: 22, Kind: "missing"},
	}
	baseline.Report(&buf, violations)
	out := buf.String()
	if !strings.Contains(out, "2 violation") {
		t.Errorf("expected violation count in output: %s", out)
	}
	if !strings.Contains(out, "8080") || !strings.Contains(out, "22") {
		t.Errorf("expected port numbers in output: %s", out)
	}
}

func TestViolation_String(t *testing.T) {
	v := baseline.Violation{Port: 443, Kind: "unexpected"}
	if !strings.Contains(v.String(), "443") {
		t.Errorf("String() missing port: %s", v.String())
	}
}
