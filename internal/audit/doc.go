// Package audit provides structured, append-only audit logging for
// portwatch. Each port change event is recorded as a newline-delimited
// JSON entry containing the event type, affected port, host, severity,
// and a UTC timestamp.
//
// Usage:
//
//	logger, err := audit.NewFileLogger("/var/log/portwatch/audit.log")
//	if err != nil { ... }
//	logger.Record(audit.Entry{
//		Event:    "opened",
//		Port:     8080,
//		Host:     "localhost",
//		Severity: "warn",
//	})
package audit
