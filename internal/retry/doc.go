// Package retry provides configurable retry logic with exponential backoff
// for use when performing transient operations such as network scans or
// webhook deliveries.
//
// Usage:
//
//	err := retry.Do(ctx, retry.DefaultPolicy(), func() error {
//		return doSomething()
//	})
package retry
