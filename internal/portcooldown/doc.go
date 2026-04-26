// Package portcooldown enforces a quiet period between successive port
// state transitions.
//
// Ports that flap between open and closed states within the configured
// window are suppressed until the state has been stable for at least
// the cooldown duration. This prevents downstream pipeline stages and
// notifiers from being flooded by transient connectivity changes.
//
// Basic usage:
//
//	cd := portcooldown.New(5 * time.Second)
//
//	if cd.Allow(port, "open") {
//	    // forward the event
//	}
package portcooldown
