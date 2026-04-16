// Package audit provides structured audit logging for port change events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Port      int       `json:"port"`
	Host      string    `json:"host"`
	Severity  string    `json:"severity"`
}

// Logger writes audit entries as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// NewLogger returns a Logger writing to w.
func NewLogger(w io.Writer) *Logger {
	return &Logger{w: w}
}

// NewFileLogger opens or creates the file at path and returns a Logger.
func NewFileLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	return NewLogger(f), nil
}

// Record writes an audit entry. Timestamp is set to now if zero.
func (l *Logger) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}

// ReadAll decodes all entries from r.
func ReadAll(r io.Reader) ([]Entry, error) {
	var entries []Entry
	dec := json.NewDecoder(r)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
