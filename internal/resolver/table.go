package resolver

// ServiceTable returns a copy of the built-in well-known port map.
// Callers may inspect or extend it before passing to New.
func ServiceTable() map[int]string {
	copy := make(map[int]string, len(wellKnown))
	for k, v := range wellKnown {
		copy[k] = v
	}
	return copy
}

// Merge combines base with overrides, returning a new map.
// Keys present in overrides take precedence over base.
func Merge(base, overrides map[int]string) map[int]string {
	out := make(map[int]string, len(base)+len(overrides))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overrides {
		out[k] = v
	}
	return out
}
