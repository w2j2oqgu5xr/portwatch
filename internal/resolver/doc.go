// Package resolver provides port-to-service name resolution for portwatch.
//
// It maintains a built-in table of well-known port assignments (ssh, http,
// postgres, etc.) and supports caller-supplied overrides for custom services.
//
// Usage:
//
//	r := resolver.New(map[int]string{8888: "my-service"})
//	fmt.Println(r.Label(22))   // "22/ssh"
//	fmt.Println(r.Label(8888)) // "8888/my-service"
package resolver
