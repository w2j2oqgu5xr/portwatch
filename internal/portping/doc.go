// Package portping provides TCP-level latency measurement for individual ports.
//
// Use New to create a Pinger, then call Ping or PingAll to measure
// round-trip time to one or more ports on a host. Summarise aggregates
// a series of results into packet-loss and latency statistics.
//
// Example:
//
//	p := portping.New(2 * time.Second)
//	res := p.Ping(ctx, "localhost", 80)
//	fmt.Println(res)
package portping
