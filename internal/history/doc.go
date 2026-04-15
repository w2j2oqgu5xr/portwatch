// Package history provides persistent recording and retrieval of port-change
// events observed by portwatch.
//
// Events are stored as newline-delimited JSON (NDJSON) so the log file remains
// human-readable and easy to process with standard Unix tools such as jq.
//
// Basic usage:
//
//	rec := history.NewRecorder("/var/log/portwatch/history.ndjson")
//
//	err := rec.Record(history.Entry{
//		Event: "opened",
//		Host:  "localhost",
//		Port:  8080,
//	})
//
//	entries, err := history.ReadAll("/var/log/portwatch/history.ndjson")
package history
