package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Config holds the configuration required to build a Notifier pipeline.
type Config struct {
	// Targets is a comma-separated list of notifier types: "console", "file", "webhook".
	Targets string
	// LogFile is the path used when the "file" target is enabled.
	LogFile string
	// WebhookURL is the endpoint used when the "webhook" target is enabled.
	WebhookURL string
	// Output overrides the writer used for "console"; defaults to os.Stdout.
	Output io.Writer
}

// Build constructs a Notifier from the provided Config.
// Unknown target names are silently ignored so the tool remains forward-compatible.
func Build(cfg Config) (Notifier, error) {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	var nn []Notifier
	for _, raw := range strings.Split(cfg.Targets, ",") {
		target := strings.TrimSpace(strings.ToLower(raw))
		switch target {
		case "console":
			nn = append(nn, NewWriterNotifier(cfg.Output))
		case "file":
			if cfg.LogFile == "" {
				return nil, fmt.Errorf("notify: 'file' target requires LogFile to be set")
			}
			f, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				return nil, fmt.Errorf("notify: open log file: %w", err)
			}
			nn = append(nn, NewWriterNotifier(f))
		case "webhook":
			if cfg.WebhookURL == "" {
				return nil, fmt.Errorf("notify: 'webhook' target requires WebhookURL to be set")
			}
			nn = append(nn, NewWebhookNotifier(cfg.WebhookURL))
		}
	}

	if len(nn) == 0 {
		// Fallback: always emit to stdout so events are never silently dropped.
		nn = append(nn, NewWriterNotifier(cfg.Output))
	}

	return NewMulti(nn...), nil
}
