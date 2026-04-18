package portlabel_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portlabel"
)

func resolver(port int) string {
	switch port {
	case 80:
		return "http"
	case 443:
		return "https"
	case 22:
		return "ssh"
	}
	return ""
}

func TestLabel_KnownService(t *testing.T) {
	l := portlabel.New(resolver, nil)
	lbl := l.Label(80, true)
	if lbl.Service != "http" {
		t.Fatalf("expected http, got %s", lbl.Service)
	}
	if !lbl.Open {
		t.Fatal("expected open")
	}
}

func TestLabel_UnknownService(t *testing.T) {
	l := portlabel.New(resolver, nil)
	lbl := l.Label(9999, false)
	if lbl.Service != "unknown" {
		t.Fatalf("expected unknown, got %s", lbl.Service)
	}
}

func TestLabel_WithAnnotation(t *testing.T) {
	annotations := map[int]string{8080: "dev-proxy"}
	l := portlabel.New(resolver, annotations)
	lbl := l.Label(8080, true)
	if lbl.Annotation != "dev-proxy" {
		t.Fatalf("expected dev-proxy, got %s", lbl.Annotation)
	}
	if !strings.Contains(lbl.String(), "dev-proxy") {
		t.Fatalf("String() missing annotation: %s", lbl.String())
	}
}

func TestLabel_String_ClosedPort(t *testing.T) {
	l := portlabel.New(resolver, nil)
	s := l.Label(22, false).String()
	if !strings.Contains(s, "closed") {
		t.Fatalf("expected closed in %s", s)
	}
}

func TestLabelAll_ReturnsCorrectCount(t *testing.T) {
	l := portlabel.New(resolver, nil)
	ports := []int{22, 80, 443}
	labels := l.LabelAll(ports, true)
	if len(labels) != 3 {
		t.Fatalf("expected 3 labels, got %d", len(labels))
	}
	for _, lbl := range labels {
		if !lbl.Open {
			t.Errorf("expected all open, port %d was closed", lbl.Port)
		}
	}
}

func TestLabel_String_ContainsPort(t *testing.T) {
	l := portlabel.New(resolver, nil)
	s := l.Label(443, true).String()
	if !strings.Contains(s, "443") {
		t.Fatalf("expected port number in %s", s)
	}
}
