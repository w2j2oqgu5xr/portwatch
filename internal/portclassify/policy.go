package portclassify

// Policy resolves a port number to a Tier and a human-readable reason.
// Implementations may combine static tables with dynamic rules.
type Policy interface {
	Lookup(port int) (tier Tier, reason string, ok bool)
}

// StaticPolicy is a simple map-backed Policy.
type StaticPolicy struct {
	rules map[int]entry
}

type entry struct {
	tier   Tier
	reason string
}

// NewStaticPolicy returns an empty StaticPolicy.
func NewStaticPolicy() *StaticPolicy {
	return &StaticPolicy{rules: make(map[int]entry)}
}

// Add registers a tier and reason for a port.
func (p *StaticPolicy) Add(port int, tier Tier, reason string) {
	p.rules[port] = entry{tier: tier, reason: reason}
}

// Lookup implements Policy.
func (p *StaticPolicy) Lookup(port int) (Tier, string, bool) {
	if e, ok := p.rules[port]; ok {
		return e.tier, e.reason, true
	}
	return TierUnknown, "", false
}

// DefaultPolicy returns a StaticPolicy pre-populated with common port tiers.
func DefaultPolicy() *StaticPolicy {
	p := NewStaticPolicy()
	// safe
	for _, port := range []int{80, 443, 8080, 8443} {
		p.Add(port, TierSafe, "common web service")
	}
	// caution
	for _, port := range []int{22, 21, 25, 3306, 5432} {
		p.Add(port, TierCaution, "sensitive service")
	}
	// critical
	for _, port := range []int{23, 445, 3389, 4444, 6666} {
		p.Add(port, TierCritical, "high-risk service")
	}
	return p
}
