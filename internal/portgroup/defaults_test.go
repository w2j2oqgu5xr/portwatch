package portgroup_test

import (
	"testing"

	"github.com/user/portwatch/internal/portgroup"
)

func TestLoadDefaults_PopulatesRegistry(t *testing.T) {
	r := portgroup.New()
	portgroup.LoadDefaults(r)
	groups := r.All()
	if len(groups) == 0 {
		t.Fatal("expected default groups to be loaded")
	}
}

func TestLoadDefaults_WebGroupExists(t *testing.T) {
	r := portgroup.New()
	portgroup.LoadDefaults(r)
	g, ok := r.Get("web")
	if !ok {
		t.Fatal("expected 'web' group")
	}
	if len(g.Ports) == 0 {
		t.Error("expected web group to have ports")
	}
}

func TestLoadDefaults_SkipsDuplicates(t *testing.T) {
	r := portgroup.New()
	r.Add("web", []int{80})
	// Should not panic or error even though "web" already exists.
	portgroup.LoadDefaults(r)
	g, _ := r.Get("web")
	// Original single-port entry should be preserved.
	if len(g.Ports) != 1 {
		t.Errorf("expected original web group unchanged, got %d ports", len(g.Ports))
	}
}

func TestLoadDefaults_DatabaseGroupResolvable(t *testing.T) {
	r := portgroup.New()
	portgroup.LoadDefaults(r)
	ports, err := r.Resolve([]string{"database"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) == 0 {
		t.Error("expected database ports")
	}
}
