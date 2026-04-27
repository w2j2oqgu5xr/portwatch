package portsampler_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portsampler"
)

func TestNew_ClampsWindowToOne(t *testing.T) {
	s := portsampler.New(0)
	s.Record([]int{80})
	s.Record([]int{443})
	// window=1 means only the latest is kept
	all := s.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(all))
	}
	if all[0].Ports[0] != 443 {
		t.Fatalf("expected port 443, got %d", all[0].Ports[0])
	}
}

func TestRecord_RetainsWindow(t *testing.T) {
	s := portsampler.New(3)
	for i := 0; i < 5; i++ {
		s.Record([]int{i})
	}
	if got := len(s.All()); got != 3 {
		t.Fatalf("expected 3 samples, got %d", got)
	}
}

func TestLatest_EmptyReturnsZeroValue(t *testing.T) {
	s := portsampler.New(5)
	latest := s.Latest()
	if !latest.Timestamp.IsZero() {
		t.Fatal("expected zero timestamp for empty sampler")
	}
	if len(latest.Ports) != 0 {
		t.Fatal("expected empty ports for empty sampler")
	}
}

func TestLatest_ReturnsMostRecent(t *testing.T) {
	s := portsampler.New(5)
	s.Record([]int{22})
	time.Sleep(time.Millisecond)
	s.Record([]int{80, 443})

	latest := s.Latest()
	if len(latest.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(latest.Ports))
	}
}

func TestAverageCount_NoSamples(t *testing.T) {
	s := portsampler.New(5)
	if avg := s.AverageCount(); avg != 0 {
		t.Fatalf("expected 0, got %f", avg)
	}
}

func TestAverageCount_Calculated(t *testing.T) {
	s := portsampler.New(10)
	s.Record([]int{80})           // 1
	s.Record([]int{80, 443})      // 2
	s.Record([]int{80, 443, 22})  // 3
	// average = (1+2+3)/3 = 2.0
	if avg := s.AverageCount(); avg != 2.0 {
		t.Fatalf("expected 2.0, got %f", avg)
	}
}

func TestReset_ClearsAllSamples(t *testing.T) {
	s := portsampler.New(5)
	s.Record([]int{80})
	s.Record([]int{443})
	s.Reset()
	if got := len(s.All()); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestRecord_IsolatesCopy(t *testing.T) {
	s := portsampler.New(5)
	ports := []int{80, 443}
	s.Record(ports)
	ports[0] = 9999 // mutate original
	if s.Latest().Ports[0] != 80 {
		t.Fatal("sampler should store a copy, not a reference")
	}
}
