package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"go.uber.org/zap"
)

// Client handles Slack webhook communications
type Client struct {
	WebhookURL string
	Channel    string
	Username   string
	IconEmoji  string
	HTTPClient *http.Client
	// MaxRetries controls how many times to retry transient failures (0 = no retries)
	MaxRetries int
	// Backoff is the base backoff duration used for exponential backoff between retries
	Backoff time.Duration
	// Logger receives diagnostic logs; if nil, the package default logger is used
	Logger Logger
	// Sleep is used to wait between retries. Inject for testing.
	Sleep func(d time.Duration)
}

// Logger is the minimal logging interface used by the Slack client.
// Logger is a minimal structured logger supporting fields and leveled messages.
type Logger interface {
	WithFields(map[string]interface{}) Logger
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

// Option configures the Slack client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.HTTPClient = hc }
}

// WithMaxRetries sets number of retries for transient failures.
func WithMaxRetries(n int) Option {
	return func(c *Client) { c.MaxRetries = n }
}

// WithBackoff sets base backoff duration.
func WithBackoff(d time.Duration) Option {
	return func(c *Client) { c.Backoff = d }
}

// WithSleeper sets a custom sleep function (useful for tests).
func WithSleeper(sleep func(time.Duration)) Option {
	return func(c *Client) { c.Sleep = sleep }
}

// WithLogger sets a structured logger (implements Printf).
func WithLogger(l Logger) Option {
	return func(c *Client) { c.Logger = l }
}

// WithZapLogger configures the client to use a zap.Logger.
func WithZapLogger(z *zap.Logger) Option {
	return func(c *Client) { c.Logger = zapAdapter{z} }
}

// WithLogrLogger configures the client to use a logr.Logger.
func WithLogrLogger(l logr.Logger) Option {
	return func(c *Client) { c.Logger = logrAdapter{l} }
}

// WithTimeout sets HTTP client timeout conveniently.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if c.HTTPClient == nil {
			c.HTTPClient = &http.Client{Timeout: d}
		} else {
			c.HTTPClient.Timeout = d
		}
	}
}

// stdLogger wraps the standard library logger.
type stdLogger struct{ l *log.Logger }

func (s stdLogger) WithFields(f map[string]interface{}) Logger {
	// naive implementation: include fields in prefix
	prefix := ""
	for k, v := range f {
		prefix += fmt.Sprintf("%s=%v ", k, v)
	}
	return stdLogger{l: log.New(s.l.Writer(), prefix, s.l.Flags())}
}

func (s stdLogger) Infof(format string, v ...interface{})  { s.l.Printf(format, v...) }
func (s stdLogger) Errorf(format string, v ...interface{}) { s.l.Printf("ERROR: "+format, v...) }

// zapAdapter adapts zap.Logger to our Logger interface.
type zapAdapter struct{ z *zap.Logger }

func (z zapAdapter) WithFields(f map[string]interface{}) Logger {
	fields := make([]zap.Field, 0, len(f))
	for k, v := range f {
		fields = append(fields, zap.Any(k, v))
	}
	return zapAdapter{z: z.z.With(fields...)}
}

func (z zapAdapter) Infof(format string, v ...interface{})  { z.z.Sugar().Infof(format, v...) }
func (z zapAdapter) Errorf(format string, v ...interface{}) { z.z.Sugar().Errorf(format, v...) }

// logrAdapter adapts logr.Logger.
type logrAdapter struct{ l logr.Logger }

func (r logrAdapter) WithFields(f map[string]interface{}) Logger {
	// convert to key-values
	kv := make([]interface{}, 0, len(f)*2)
	for k, v := range f {
		kv = append(kv, k, v)
	}
	return logrAdapter{l: r.l.WithValues(kv...)}
}

func (r logrAdapter) Infof(format string, v ...interface{}) { r.l.Info(fmt.Sprintf(format, v...)) }
func (r logrAdapter) Errorf(format string, v ...interface{}) {
	r.l.Error(fmt.Errorf(format, v...), fmt.Sprintf(format, v...))
}

// NewClient creates a new Slack client
// NewClient creates a new Slack client. Use functional options to customize behavior.
func NewClient(webhookURL, channel string, opts ...Option) *Client {
	c := &Client{
		WebhookURL: webhookURL,
		Channel:    channel,
		Username:   "ops-tool",
		IconEmoji:  ":robot_face:",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		MaxRetries: 2,
		Backoff:    500 * time.Millisecond,
		Logger:     stdLogger{l: log.Default()},
	}

	for _, o := range opts {
		o(c)
	}

	// ensure sensible defaults after options
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	}
	if c.Backoff <= 0 {
		c.Backoff = 500 * time.Millisecond
	}
	if c.Logger == nil {
		c.Logger = stdLogger{l: log.Default()}
	}

	return c
}

// Message represents a Slack message
type Message struct {
	Channel     string       `json:"channel,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents a Slack message attachment
type Attachment struct {
	Color      string  `json:"color,omitempty"`
	Title      string  `json:"title,omitempty"`
	TitleLink  string  `json:"title_link,omitempty"`
	Text       string  `json:"text,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
	ImageURL   string  `json:"image_url,omitempty"`
	ThumbURL   string  `json:"thumb_url,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	FooterIcon string  `json:"footer_icon,omitempty"`
	Timestamp  int64   `json:"ts,omitempty"`
}

// Field represents a field in a Slack attachment
type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

// SendMessage sends a text message to Slack
func (c *Client) SendMessage(text string) error {
	msg := Message{
		Channel:   c.Channel,
		Username:  c.Username,
		IconEmoji: c.IconEmoji,
		Text:      text,
	}

	return c.sendWebhook(msg)
}

// SendAlert sends an alert notification to Slack
func (c *Client) SendAlert(alertName, severity, message string, metadata map[string]string) error {
	color := "#36a64f" // green
	switch severity {
	case "warning":
		color = "#ff9900" // orange
	case "critical":
		color = "#ff0000" // red
	}

	fields := []Field{
		{
			Title: "Severity",
			Value: severity,
			Short: true,
		},
		{
			Title: "Alert",
			Value: alertName,
			Short: true,
		},
	}

	for key, value := range metadata {
		fields = append(fields, Field{
			Title: key,
			Value: value,
			Short: true,
		})
	}

	attachment := Attachment{
		Color:  color,
		Title:  "⚠️ Alert Triggered",
		Text:   message,
		Fields: fields,
	}

	msg := Message{
		Channel:     c.Channel,
		Username:    c.Username,
		IconEmoji:   c.IconEmoji,
		Attachments: []Attachment{attachment},
	}

	return c.sendWebhook(msg)
}

// SendDeployment sends a deployment notification to Slack
func (c *Client) SendDeployment(app, status, version string, metadata map[string]string) error {
	color := "#36a64f"
	switch status {
	case "failed":
		color = "#ff0000"
	case "pending":
		color = "#0099ff"
	}

	fields := []Field{
		{
			Title: "App",
			Value: app,
			Short: true,
		},
		{
			Title: "Status",
			Value: status,
			Short: true,
		},
		{
			Title: "Version",
			Value: version,
			Short: true,
		},
	}

	for key, value := range metadata {
		fields = append(fields, Field{
			Title: key,
			Value: value,
			Short: true,
		})
	}

	emoji := "✅"
	switch status {
	case "failed":
		emoji = "❌"
	case "pending":
		emoji = "⏳"
	}

	attachment := Attachment{
		Color:  color,
		Title:  emoji + " Deployment Notification",
		Text:   fmt.Sprintf("%s deployment of %s", status, app),
		Fields: fields,
	}

	msg := Message{
		Channel:     c.Channel,
		Username:    c.Username,
		IconEmoji:   c.IconEmoji,
		Attachments: []Attachment{attachment},
	}

	return c.sendWebhook(msg)
}

// SendOperationStatus sends an operation status update to Slack
func (c *Client) SendOperationStatus(operation, status string, duration float64, metadata map[string]string) error {
	color := "#36a64f"
	switch status {
	case "failed":
		color = "#ff0000"
	case "pending":
		color = "#0099ff"
	}

	fields := []Field{
		{
			Title: "Operation",
			Value: operation,
			Short: true,
		},
		{
			Title: "Status",
			Value: status,
			Short: true,
		},
		{
			Title: "Duration",
			Value: fmt.Sprintf("%.2f seconds", duration),
			Short: true,
		},
	}

	for key, value := range metadata {
		fields = append(fields, Field{
			Title: key,
			Value: value,
			Short: true,
		})
	}

	emoji := "✅"
	switch status {
	case "failed":
		emoji = "❌"
	case "pending":
		emoji = "⏳"
	}

	attachment := Attachment{
		Color:  color,
		Title:  emoji + " Operation Status",
		Text:   fmt.Sprintf("%s: %s", operation, status),
		Fields: fields,
	}

	msg := Message{
		Channel:     c.Channel,
		Username:    c.Username,
		IconEmoji:   c.IconEmoji,
		Attachments: []Attachment{attachment},
	}

	return c.sendWebhook(msg)
}

// SendResourceStatus sends resource status update to Slack
func (c *Client) SendResourceStatus(resource, resourceType, status string, stats map[string]string) error {
	color := "#36a64f"
	switch status {
	case "unhealthy", "error":
		color = "#ff0000"
	case "degraded", "warning":
		color = "#ff9900"
	}

	fields := []Field{
		{
			Title: "Resource",
			Value: resource,
			Short: true,
		},
		{
			Title: "Type",
			Value: resourceType,
			Short: true,
		},
		{
			Title: "Status",
			Value: status,
			Short: true,
		},
	}

	for key, value := range stats {
		fields = append(fields, Field{
			Title: key,
			Value: value,
			Short: true,
		})
	}

	attachment := Attachment{
		Color:  color,
		Title:  "📊 Resource Status",
		Text:   fmt.Sprintf("%s is %s", resource, status),
		Fields: fields,
	}

	msg := Message{
		Channel:     c.Channel,
		Username:    c.Username,
		IconEmoji:   c.IconEmoji,
		Attachments: []Attachment{attachment},
	}

	return c.sendWebhook(msg)
}

// SendMetricsSummary sends a metrics summary to Slack
func (c *Client) SendMetricsSummary(title string, metrics map[string]interface{}) error {
	fields := []Field{}

	for key, value := range metrics {
		fields = append(fields, Field{
			Title: key,
			Value: fmt.Sprintf("%v", value),
			Short: true,
		})
	}

	attachment := Attachment{
		Color:  "#0099ff",
		Title:  "📈 Metrics Summary",
		Text:   title,
		Fields: fields,
	}

	msg := Message{
		Channel:     c.Channel,
		Username:    c.Username,
		IconEmoji:   c.IconEmoji,
		Attachments: []Attachment{attachment},
	}

	return c.sendWebhook(msg)
}

// sendWebhook sends the message via Slack webhook
func (c *Client) sendWebhook(msg Message) error {
	if c.WebhookURL == "" {
		return fmt.Errorf("slack webhook URL not configured")
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Ensure sensible defaults
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	}
	if c.Logger == nil {
		c.Logger = stdLogger{l: log.Default()}
	}
	if c.Backoff <= 0 {
		c.Backoff = 500 * time.Millisecond
	}

	var lastErr error
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		// create a fresh reader for each attempt
		resp, err := c.HTTPClient.Post(
			c.WebhookURL,
			"application/json",
			bytes.NewReader(payload),
		)
		if err != nil {
			lastErr = fmt.Errorf("failed to send slack message: %w", err)
			c.Logger.Errorf("slack send attempt %d/%d failed: %v", attempt+1, c.MaxRetries+1, err)
			if attempt < c.MaxRetries {
				sleep := c.Backoff * (1 << attempt)
				if c.Sleep != nil {
					c.Sleep(sleep)
				} else {
					time.Sleep(sleep)
				}
				continue
			}
			return lastErr
		}

		body, readErr := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		if readErr != nil {
			c.Logger.Errorf("failed reading response body on attempt %d: %v", attempt+1, readErr)
		}
		if closeErr != nil {
			c.Logger.Errorf("failed closing response body on attempt %d: %v", attempt+1, closeErr)
		}

		c.Logger.Infof("slack webhook response attempt %d status=%d body=%s", attempt+1, resp.StatusCode, string(body))

		// 200 OK -> done
		if resp.StatusCode == http.StatusOK {
			return nil
		}

		// Retry on 5xx server errors
		if resp.StatusCode >= 500 && resp.StatusCode < 600 && attempt < c.MaxRetries {
			c.Logger.Infof("retrying slack send due to server error %d", resp.StatusCode)
			if c.Sleep != nil {
				c.Sleep(c.Backoff * (1 << attempt))
			} else {
				time.Sleep(c.Backoff * (1 << attempt))
			}
			continue
		}

		// Non-retryable error
		return fmt.Errorf("slack webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return lastErr
}
