// Package ipinfo fetches metadata for IP addresses via the ipinfo.io API.
//
// Usage:
//
//	lookup := ipinfo.New("optional-api-token")
//	info, err := lookup.Get("8.8.8.8")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(info) // 8.8.8.8 (Mountain View, US) AS15169 Google LLC
//
// The token may be empty for anonymous access, subject to rate limiting.
package ipinfo
