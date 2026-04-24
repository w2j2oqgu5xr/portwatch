package portdiff_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portdiff"
	"github.com/user/portwatch/internal/resolver"
)

func newComputer() *portdiff.Computer {
	res := resolver.New(nil)
	return portdiff.New(res)
}

func TestCompute_OpenedPorts(t *testing.T) {
	c := newComputer()
	d := c.Compute([]int{80}, []int{80, 443})

	if len(d.Opened) != 1 {
		t.Fatalf("expected 1 opened, got %d", len(d.Opened))
	}
	if d.Opened[0].Port != 443 {
		t.Errorf("expected port 443, got %d", d.Opened[0].Port)
	}
	if len(d.Closed) != 0 {
		t.Errorf("expected no closed ports, got %d", len(d.Closed))
	}
}

func TestCompute_ClosedPorts(t *testing.T) {
	c := newComputer()
	d := c.Compute([]int{80, 8080}, []int{80})

	if len(d.Closed) != 1 {
		t.Fatalf("expected 1 closed, got %d", len(d.Closed))
	}
	if d.Closed[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", d.Closed[0].Port)
	}
}

func TestCompute_NoChanges(t *testing.T) {
	c := newComputer()
	d := c.Compute([]int{22, 80}, []int{22, 80})

	if !d.IsEmpty() {
		t.Errorf("expected empty diff, got opened=%d closed=%d", len(d.Opened), len(d.Closed))
	}
}

func TestCompute_SortedOutput(t *testing.T) {
	c := newComputer()
	d := c.Compute([]int{}, []int{9000, 443, 80})

	ports := make([]int, len(d.Opened))
	for i, ch := range d.Opened {
		ports[i] = ch.Port
	}
	for i := 1; i < len(ports); i++ {
		if ports[i] < ports[i-1] {
			t.Errorf("opened ports not sorted: %v", ports)
		}
	}
}

func TestChange_String(t *testing.T) {
	ch := portdiff.Change{Port: 80, Service: "http", Direction: "opened"}
	got := ch.String()
	if !strings.Contains(got, "80") || !strings.Contains(got, "http") {
		t.Errorf("unexpected Change.String(): %q", got)
	}
}

func TestDiff_Summary_WithChanges(t *testing.T) {
	c := newComputer()
	d := c.Compute([]int{22}, []int{22, 80, 443})
	s := d.Summary()
	if !strings.Contains(s, "opened") {
		t.Errorf("summary missing 'opened': %q", s)
	}
}

func TestDiff_Summary_NoChanges(t *testing.T) {
	d := portdiff.Diff{}
	if got := d.Summary(); got != "no changes" {
		t.Errorf("expected 'no changes', got %q", got)
	}
}

func TestCompute_ServiceAnnotated(t *testing.T) {
	c := newComputer()
	d := c.Compute([]int{}, []int{22})

	if len(d.Opened) != 1 {
		t.Fatalf("expected 1 opened port")
	}
	if d.Opened[0].Service == "" {
		t.Error("expected non-empty service name for port 22")
	}
}
