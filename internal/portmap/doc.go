// Package portmap provides a thread-safe registry that associates open port
// numbers with runtime metadata including protocol, owning process name, and
// PID. It is used to enrich scan results and alert messages with contextual
// information beyond the bare port number.
//
// Usage:
//
//	m := portmap.New()
//	m.Set(80, portmap.Entry{Port: 80, Protocol: "tcp", PID: 1234, Process: "nginx"})
//	if e, ok := m.Get(80); ok {
//		fmt.Println(e)
//	}
package portmap
