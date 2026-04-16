// Package banner provides TCP service banner grabbing for open ports.
package banner

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// DefaultTimeout is the maximum time to wait for a banner response.
const DefaultTimeout = 2 * time.Second

// Result holds the grabbed banner for a single port.
type Result struct {
	Port    int
	Banner  string
	Err     error
}

// Grabber grabs banners from open TCP ports.
type Grabber struct {
	Timeout time.Duration
}

// New returns a Grabber with the given timeout. If timeout is zero,
// DefaultTimeout is used.
func New(timeout time.Duration) *Grabber {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &Grabber{Timeout: timeout}
}

// Grab attempts to read a banner from host:port.
func (g *Grabber) Grab(host string, port int) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, g.Timeout)
	if err != nil {
		return Result{Port: port, Err: err}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(g.Timeout))

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil && n == 0 {
		return Result{Port: port, Err: err}
	}

	banner := strings.TrimSpace(string(buf[:n]))
	return Result{Port: port, Banner: banner}
}

// GrabAll grabs banners for all provided ports concurrently.
func (g *Grabber) GrabAll(host string, ports []int) []Result {
	results := make([]Result, len(ports))
	type indexed struct {
		i int
		r Result
	}
	ch := make(chan indexed, len(ports))
	for i, p := range ports {
		go func(idx, port int) {
			ch <- indexed{idx, g.Grab(host, port)}
		}(i, p)
	}
	for range ports {
		v := <-ch
		results[v.i] = v.r
	}
	return results
}
