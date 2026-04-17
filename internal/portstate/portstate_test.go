package portstate_test

import (
	"testing"

	"github.com/user/portwatch/internal/portstate"
)

func TestUpdate_DetectsOpenedPort(t *testing.T) {
	tr := portstate.New()
	changes := tr.Update([]int{80, 443})
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	for _, c := range changes {
		if !c.Opened {
			t.Errorf("expected Opened=true for port %d", c.Port)
		}
	}
}

func TestUpdate_DetectsClosedPort(t *testing.T) {
	tr := portstate.New()
	tr.Update([]int{80, 443})
	changes := tr.Update([]int{80})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Port != 443 || changes[0].Opened {
		t.Errorf("expected port 443 closed, got %+v", changes[0])
	}
}

func TestUpdate_NoChanges(t *testing.T) {
	tr := portstate.New()
	tr.Update([]int{80})
	changes := tr.Update([]int{80})
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestUpdate_ReopenedPort(t *testing.T) {
	tr := portstate.New()
	tr.Update([]int{8080})
	tr.Update([]int{})
	changes := tr.Update([]int{8080})
	if len(changes) != 1 || !changes[0].Opened {
		t.Errorf("expected port 8080 to be reopened")
	}
}

func TestOpenPorts_ReflectsCurrentState(t *testing.T) {
	tr := portstate.New()
	tr.Update([]int{22, 80, 443})
	tr.Update([]int{22, 80})
	ports := tr.OpenPorts()
	if len(ports) != 2 {
		t.Errorf("expected 2 open ports, got %d", len(ports))
	}
}

func TestSnapshot_ContainsAllTracked(t *testing.T) {
	tr := portstate.New()
	tr.Update([]int{22, 80})
	snap := tr.Snapshot()
	if len(snap) != 2 {
		t.Errorf("expected 2 states in snapshot, got %d", len(snap))
	}
}

func TestSnapshot_ClosedPortRetained(t *testing.T) {
	tr := portstate.New()
	tr.Update([]int{9000})
	tr.Update([]int{})
	snap := tr.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 state, got %d", len(snap))
	}
	if snap[0].Open {
		t.Errorf("expected port 9000 to be marked closed")
	}
}
