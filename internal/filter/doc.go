// Package filter provides port filtering for portwatch.
//
// A Filter is constructed with an optional allow list and an optional
// deny list. When an allow list is present, only listed ports are
// considered for monitoring. Ports in the deny list are always
// excluded, regardless of the allow list.
//
// Example usage:
//
//	f := filter.New([]int{80, 443}, []int{8080})
//	open := f.Apply(scannedPorts)
package filter
