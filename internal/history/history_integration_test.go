package history_test

import (
	"testing"

	"github.com/user/portwatch/internal/history"
)

// TestRecorder_ConcurrentWrites verifies that concurrent Record calls do not
// corrupt the log file and that all entries are recoverable.
func TestRecorder_ConcurrentWrites(t *testing.T) {
	p := tempPath(t)
	rec := history.NewRecorder(p)

	const workers = 10
	const perWorker = 5

	done := make(chan error, workers)
	for w := 0; w < workers; w++ {
		go func(port int) {
			var lastErr error
			for i := 0; i < perWorker; i++ {
				if err := rec.Record(history.Entry{
					Event: "opened",
					Host:  "localhost",
					Port:  port,
				}); err != nil {
					lastErr = err
				}
			}
			done <- lastErr
		}(8000 + w)
	}

	for i := 0; i < workers; i++ {
		if err := <-done; err != nil {
			t.Errorf("worker error: %v", err)
		}
	}

	entries, err := history.ReadAll(p)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	want := workers * perWorker
	if len(entries) != want {
		t.Errorf("expected %d entries, got %d", want, len(entries))
	}
}
