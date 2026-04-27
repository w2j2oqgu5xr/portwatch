package portclassify_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portclassify"
)

func defaultClassifier() *portclassify.Classifier {
	return portclassify.New(portclassify.DefaultPolicy())
}

func TestClassify_SafePort(t *testing.T) {
	c := defaultClassifier()
	r := c.Classify(443)
	if r.Tier != portclassify.TierSafe {
		t.Fatalf("expected safe, got %s", r.Tier)
	}
}

func TestClassify_CautionPort(t *testing.T) {
	c := defaultClassifier()
	r := c.Classify(22)
	if r.Tier != portclassify.TierCaution {
		t.Fatalf("expected caution, got %s", r.Tier)
	}
}

func TestClassify_CriticalPort(t *testing.T) {
	c := defaultClassifier()
	r := c.Classify(3389)
	if r.Tier != portclassify.TierCritical {
		t.Fatalf("expected critical, got %s", r.Tier)
	}
}

func TestClassify_UnknownPort(t *testing.T) {
	c := defaultClassifier()
	r := c.Classify(9999)
	if r.Tier != portclassify.TierUnknown {
		t.Fatalf("expected unknown, got %s", r.Tier)
	}
}

func TestClassifyAll_ReturnsAllResults(t *testing.T) {
	c := defaultClassifier()
	results := c.ClassifyAll([]int{80, 22, 23, 9999})
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	expected := []portclassify.Tier{
		portclassify.TierSafe,
		portclassify.TierCaution,
		portclassify.TierCritical,
		portclassify.TierUnknown,
	}
	for i, r := range results {
		if r.Tier != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], r.Tier)
		}
	}
}

func TestResult_String_ContainsPort(t *testing.T) {
	r := portclassify.Result{Port: 443, Tier: portclassify.TierSafe, Reason: "web"}
	s := r.String()
	if !strings.Contains(s, "443") {
		t.Errorf("expected port in string, got: %s", s)
	}
}

func TestStaticPolicy_CustomRule(t *testing.T) {
	p := portclassify.NewStaticPolicy()
	p.Add(12345, portclassify.TierCritical, "test service")
	tier, reason, ok := p.Lookup(12345)
	if !ok || tier != portclassify.TierCritical || reason != "test service" {
		t.Fatalf("unexpected lookup result: ok=%v tier=%s reason=%s", ok, tier, reason)
	}
}

func TestStaticPolicy_MissingPort(t *testing.T) {
	p := portclassify.NewStaticPolicy()
	_, _, ok := p.Lookup(1)
	if ok {
		t.Fatal("expected not ok for missing port")
	}
}

func TestTier_String_AllValues(t *testing.T) {
	cases := map[portclassify.Tier]string{
		portclassify.TierSafe:     "safe",
		portclassify.TierCaution:  "caution",
		portclassify.TierCritical: "critical",
		portclassify.TierUnknown:  "unknown",
	}
	for tier, want := range cases {
		if tier.String() != want {
			t.Errorf("Tier(%d).String() = %q, want %q", tier, tier.String(), want)
		}
	}
}
