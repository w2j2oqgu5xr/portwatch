// Package notify defines the Notifier interface and concrete implementations
// used by portwatch to dispatch port-change events to one or more output
// channels.
//
// Available notifiers:
//
//   - WriterNotifier — writes formatted events to any io.Writer (stdout, file).
//   - WebhookNotifier — HTTP POSTs event payloads to a configured URL.
//   - Multi — fans a single event out to a slice of Notifier implementations.
//
// Example usage:
//
//	n := notify.NewMulti(
//		notify.NewWriterNotifier(os.Stdout),
//		notify.NewWebhookNotifier("https://hooks.example.com/portwatch"),
//	)
//	n.Notify(notify.Event{Host: "localhost", Port: 443, Kind: "opened"})
package notify
