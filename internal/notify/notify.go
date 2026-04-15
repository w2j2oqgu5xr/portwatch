// Package notify provides notification channel abstractions for portwatch.
// It supports multiple output targets such as console, file, and webhook.
package notify

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Event represents a port change notification payload.
type Event struct {
	Timestamp time.Time
	Host      string
	Port      int
	Kind      string // "opened" or "closed"
	Message   string
}

// Notifier is the interface implemented by all notification targets.
type Notifier interface {
	Notify(e Event) error
}

// Multi fans out a single event to multiple notifiers.
type Multi struct {
	notifiers []Notifier
}

// NewMulti returns a Multi notifier that delegates to all provided notifiers.
func NewMulti(nn ...Notifier) *Multi {
	return &Multi{notifiers: nn}
}

// Notify sends the event to every registered notifier, collecting errors.
func (m *Multi) Notify(e Event) error {
	var errs []string
	for _, n := range m.notifiers {
		if err := n.Notify(e); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("notify errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

// WriterNotifier writes formatted events to any io.Writer.
type WriterNotifier struct {
	w io.Writer
}

// NewWriterNotifier returns a WriterNotifier backed by w.
func NewWriterNotifier(w io.Writer) *WriterNotifier {
	return &WriterNotifier{w: w}
}

// Notify formats the event and writes it to the underlying writer.
func (wn *WriterNotifier) Notify(e Event) error {
	_, err := fmt.Fprintf(wn.w, "[%s] %s port %d on %s: %s\n",
		e.Timestamp.Format(time.RFC3339), e.Kind, e.Port, e.Host, e.Message)
	return err
}

// WebhookNotifier posts events as plain-text HTTP POST requests.
type WebhookNotifier struct {
	URL    string
	client *http.Client
}

// NewWebhookNotifier returns a WebhookNotifier that posts to the given URL.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Notify sends the event as a plain-text POST to the configured webhook URL.
func (wh *WebhookNotifier) Notify(e Event) error {
	body := fmt.Sprintf("[%s] %s port %d on %s: %s",
		e.Timestamp.Format(time.RFC3339), e.Kind, e.Port, e.Host, e.Message)
	resp, err := wh.client.Post(wh.URL, "text/plain", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}
