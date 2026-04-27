// Package portclassify categorises ports into risk tiers based on
// well-known service associations and configurable policy rules.
package portclassify

import "fmt"

// Tier represents the risk classification of a port.
type Tier int

const (
	TierUnknown  Tier = iota
	TierSafe           // expected, low-risk services
	TierCaution        // elevated attention warranted
	TierCritical       // high-risk or unexpected exposure
)

func (t Tier) String() string {
	switch t {
	case TierSafe:
		return "safe"
	case TierCaution:
		return "caution"
	case TierCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// Result holds the classification outcome for a single port.
type Result struct {
	Port   int
	Tier   Tier
	Reason string
}

func (r Result) String() string {
	return fmt.Sprintf("port %d: %s (%s)", r.Port, r.Tier, r.Reason)
}

// Classifier assigns risk tiers to ports.
type Classifier struct {
	policy Policy
}

// New returns a Classifier using the provided Policy.
func New(p Policy) *Classifier {
	return &Classifier{policy: p}
}

// Classify returns the risk tier for the given port.
func (c *Classifier) Classify(port int) Result {
	if tier, reason, ok := c.policy.Lookup(port); ok {
		return Result{Port: port, Tier: tier, Reason: reason}
	}
	return Result{Port: port, Tier: TierUnknown, Reason: "no policy match"}
}

// ClassifyAll classifies a slice of ports and returns all results.
func (c *Classifier) ClassifyAll(ports []int) []Result {
	out := make([]Result, len(ports))
	for i, p := range ports {
		out[i] = c.Classify(p)
	}
	return out
}
