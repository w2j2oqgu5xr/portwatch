package portmap_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/portmap"
)

func TestSet_And_Get(t *testing.T) {
	m := portmap.New()
	e := portmap.Entry{Port: 80, Protocol: "tcp", PID: 123, Process: "nginx"}
	m.Set(80, e)

	got, ok := m.Get(80)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Process != "nginx" {
		t.Errorf("expected nginx, got %s", got.Process)
	}
}

func TestGet_Missing(t *testing.T) {
	m := portmap.New()
	_, ok := m.Get(9999)
	if ok {
		t.Error("expected missing entry")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	m := portmap.New()
	m.Set(443, portmap.Entry{Port: 443, Protocol: "tcp"})
	m.Delete(443)
	_, ok := m.Get(443)
	if ok {
		t.Error("expected entry to be deleted")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	m := portmap.New()
	m.Set(22, portmap.Entry{Port: 22, Protocol: "tcp", Process: "sshd"})
	m.Set(80, portmap.Entry{Port: 80, Protocol: "tcp", Process: "nginx"})

	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestLen_ReflectsCount(t *testing.T) {
	m := portmap.New()
	if m.Len() != 0 {
		t.Error("expected empty map")
	}
	m.Set(8080, portmap.Entry{Port: 8080})
	if m.Len() != 1 {
		t.Errorf("expected 1, got %d", m.Len())
	}
}

func TestEntry_String_WithProcess(t *testing.T) {
	e := portmap.Entry{Port: 80, Protocol: "tcp", PID: 42, Process: "nginx"}
	s := e.String()
	if s != "80/tcp (pid=42, nginx)" {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestEntry_String_NoProcess(t *testing.T) {
	e := portmap.Entry{Port: 443, Protocol: "tcp"}
	if e.String() != "443/tcp" {
		t.Errorf("unexpected string: %s", e.String())
	}
}

func TestConcurrentAccess(t *testing.T) {
	m := portmap.New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			m.Set(p, portmap.Entry{Port: p, Protocol: "tcp"})
			m.Get(p)
		}(i)
	}
	wg.Wait()
	if m.Len() != 50 {
		t.Errorf("expected 50 entries, got %d", m.Len())
	}
}
