// Package monitor provides continuous port-state monitoring for portwatch.
//
// It wraps the scanner package to periodically scan a set of ports on a given
// host and emits Change events whenever a port transitions between open and
// closed states.
//
// Basic usage:
//
//	 m := monitor.New("localhost", []int{80, 443, 8080}, 5*time.Second)
//	 go m.Start()
//
//	 for change := range m.Changes {
//	     fmt.Printf("port %d is now %s\n", change.Port, change.Status)
//	 }
//
// Call close(m.Stop) to shut down the monitor gracefully; the Changes channel
// will be closed once the background goroutine exits.
package monitor
