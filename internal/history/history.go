// Package history records port change events to a persistent log file,
// allowing users to review past alerts after the fact.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents a single recorded port-change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"` // "opened" | "closed"
	Host      string    `json:"host"`
	Port      int       `json:"port"`
}

// Recorder appends history entries to a newline-delimited JSON file.
type Recorder struct {
	mu   sync.Mutex
	path string
}

// NewRecorder returns a Recorder that writes to the given file path.
func NewRecorder(path string) *Recorder {
	return &Recorder{path: path}
}

// Record appends a single entry to the history file.
func (r *Recorder) Record(e Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("history: open file: %w", err)
	}
	defer f.Close()

	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("history: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// ReadAll returns all entries stored in the history file.
// If the file does not exist an empty slice is returned without error.
func ReadAll(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: read file: %w", err)
	}

	var entries []Entry
	decoder := json.NewDecoder(
		// wrap raw bytes in a reader line-by-line via a simple approach
		newLineReader(data),
	)
	for decoder.More() {
		var e Entry
		if err := decoder.Decode(&e); err != nil {
			return nil, fmt.Errorf("history: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Filter returns only the entries that satisfy the provided predicate.
// It is a convenience helper for callers that need to search or narrow
// down results returned by ReadAll without reimplementing the loop.
func Filter(entries []Entry, fn func(Entry) bool) []Entry {
	var out []Entry
	for _, e := range entries {
		if fn(e) {
			out = append(out, e)
		}
	}
	return out
}
