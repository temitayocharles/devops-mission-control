package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test that the client retries on 5xx responses and succeeds when the server recovers.
func TestSlackClientRetries(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls <= 2 {
			// First two attempts: server error
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("err"))
			return
		}
		// Third attempt: success
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	// capture logs and timings
	logger := &capturingLogger{}
	// capturing sleeper to record sleeps instead of actually sleeping
	var totalSleep time.Duration
	sleeper := func(d time.Duration) { totalSleep += d }

	client := NewClient(srv.URL, "#test",
		WithMaxRetries(3),
		WithBackoff(10*time.Millisecond),
		WithLogger(logger),
		WithHTTPClient(&http.Client{Timeout: 2 * time.Second}),
		WithSleeper(sleeper),
	)

	if err := client.SendMessage("hello"); err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}

	if calls < 3 {
		t.Fatalf("expected at least 3 calls, got %d", calls)
	}

	// verify backoff durations were recorded via sleeper (10ms + 20ms = 30ms minimum)
	if totalSleep < 30*time.Millisecond {
		t.Fatalf("expected total sleep >= 30ms, got %v", totalSleep)
	}
}

// capturingLogger records timestamps of first and last logged events and messages.
type capturingLogger struct {
	start time.Time
	end   time.Time
}

func (c *capturingLogger) WithFields(map[string]interface{}) Logger { return c }
func (c *capturingLogger) Infof(format string, v ...interface{}) {
	now := time.Now()
	if c.start.IsZero() {
		c.start = now
	}
	c.end = now
}
func (c *capturingLogger) Errorf(format string, v ...interface{}) {
	now := time.Now()
	if c.start.IsZero() {
		c.start = now
	}
	c.end = now
}
