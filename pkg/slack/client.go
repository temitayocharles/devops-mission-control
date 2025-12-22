package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client handles Slack webhook communications
type Client struct {
	WebhookURL string
	Channel    string
	Username   string
	IconEmoji  string
}

// NewClient creates a new Slack client
func NewClient(webhookURL, channel string) *Client {
	return &Client{
		WebhookURL: webhookURL,
		Channel:    channel,
		Username:   "ops-tool",
		IconEmoji:  ":robot_face:",
	}
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
	if severity == "warning" {
		color = "#ff9900" // orange
	} else if severity == "critical" {
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
	if status == "failed" {
		color = "#ff0000"
	} else if status == "pending" {
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
	if status == "failed" {
		emoji = "❌"
	} else if status == "pending" {
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
	if status == "failed" {
		color = "#ff0000"
	} else if status == "pending" {
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
	if status == "failed" {
		emoji = "❌"
	} else if status == "pending" {
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
	if status == "unhealthy" || status == "error" {
		color = "#ff0000"
	} else if status == "degraded" || status == "warning" {
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

	resp, err := http.Post(
		c.WebhookURL,
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return fmt.Errorf("failed to send slack message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack webhook returned status %d", resp.StatusCode)
	}

	return nil
}
