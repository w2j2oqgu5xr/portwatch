// Package baseline provides functionality for recording and comparing
// a trusted set of open ports (the "baseline") against the current
// observed state.
//
// A baseline is saved to disk as JSON and can be loaded on subsequent
// runs. Unexpected returns ports that are open but absent from the
// baseline; Missing returns baseline ports no longer observed open.
package baseline
