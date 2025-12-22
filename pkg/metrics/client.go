package metrics

import (
	"fmt"
	"sync"
	"time"
)

// Metric represents a single metric data point
type Metric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Unit      string            `json:"unit"`
	Timestamp time.Time         `json:"timestamp"`
	Tags      map[string]string `json:"tags,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// Event represents an operational event (command execution, deployment, etc.)
type Event struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"` // "command", "deployment", "alert", "error"
	Message   string            `json:"message"`
	Status    string            `json:"status"` // "success", "failure", "pending"
	Timestamp time.Time         `json:"timestamp"`
	Duration  time.Duration     `json:"duration,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	User      string            `json:"user,omitempty"`
	Resource  string            `json:"resource,omitempty"`
}

// Alert represents a triggered alert
type Alert struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Severity   string            `json:"severity"` // "info", "warning", "critical"
	Message    string            `json:"message"`
	Timestamp  time.Time         `json:"timestamp"`
	Resolved   bool              `json:"resolved"`
	ResolvedAt *time.Time        `json:"resolved_at,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// MetricsStore holds all collected metrics and events
type MetricsStore struct {
	mu      sync.RWMutex
	Metrics []Metric
	Events  []Event
	Alerts  []Alert
	MaxSize int // Maximum number of items to store
}

// NewMetricsStore creates a new metrics store
func NewMetricsStore(maxSize int) *MetricsStore {
	if maxSize <= 0 {
		maxSize = 10000
	}
	return &MetricsStore{
		Metrics: make([]Metric, 0, maxSize),
		Events:  make([]Event, 0, maxSize),
		Alerts:  make([]Alert, 0, maxSize),
		MaxSize: maxSize,
	}
}

// RecordMetric adds a metric to the store
func (s *MetricsStore) RecordMetric(name string, value float64, unit string, tags map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metric := Metric{
		Name:      name,
		Value:     value,
		Unit:      unit,
		Timestamp: time.Now(),
		Tags:      tags,
	}

	s.Metrics = append(s.Metrics, metric)
	if len(s.Metrics) > s.MaxSize {
		s.Metrics = s.Metrics[len(s.Metrics)-s.MaxSize:]
	}
}

// RecordEvent adds an event to the store
func (s *MetricsStore) RecordEvent(eventType, message, status string, duration time.Duration, metadata map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event := Event{
		ID:        fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Type:      eventType,
		Message:   message,
		Status:    status,
		Timestamp: time.Now(),
		Duration:  duration,
		Metadata:  metadata,
	}

	s.Events = append(s.Events, event)
	if len(s.Events) > s.MaxSize {
		s.Events = s.Events[len(s.Events)-s.MaxSize:]
	}
}

// CreateAlert creates a new alert
func (s *MetricsStore) CreateAlert(name, severity, message string, metadata map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	alert := Alert{
		ID:        fmt.Sprintf("alert_%d", time.Now().UnixNano()),
		Name:      name,
		Severity:  severity,
		Message:   message,
		Timestamp: time.Now(),
		Resolved:  false,
		Metadata:  metadata,
	}

	s.Alerts = append(s.Alerts, alert)
	if len(s.Alerts) > s.MaxSize {
		s.Alerts = s.Alerts[len(s.Alerts)-s.MaxSize:]
	}
}

// ResolveAlert marks an alert as resolved
func (s *MetricsStore) ResolveAlert(alertID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Alerts {
		if s.Alerts[i].ID == alertID {
			s.Alerts[i].Resolved = true
			now := time.Now()
			s.Alerts[i].ResolvedAt = &now
			return nil
		}
	}
	return fmt.Errorf("alert %s not found", alertID)
}

// GetMetrics returns all metrics
func (s *MetricsStore) GetMetrics() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Metric{}, s.Metrics...)
}

// GetEvents returns all events
func (s *MetricsStore) GetEvents() []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Event{}, s.Events...)
}

// GetAlerts returns all alerts
func (s *MetricsStore) GetAlerts() []Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Alert{}, s.Alerts...)
}

// GetActiveAlerts returns only unresolved alerts
func (s *MetricsStore) GetActiveAlerts() []Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var active []Alert
	for _, alert := range s.Alerts {
		if !alert.Resolved {
			active = append(active, alert)
		}
	}
	return active
}

// GetMetricsByName returns metrics with a specific name
func (s *MetricsStore) GetMetricsByName(name string) []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Metric
	for _, m := range s.Metrics {
		if m.Name == name {
			result = append(result, m)
		}
	}
	return result
}

// GetEventsByType returns events of a specific type
func (s *MetricsStore) GetEventsByType(eventType string) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Event
	for _, e := range s.Events {
		if e.Type == eventType {
			result = append(result, e)
		}
	}
	return result
}

// Clear removes all data from the store
func (s *MetricsStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Metrics = []Metric{}
	s.Events = []Event{}
	s.Alerts = []Alert{}
}

// GetStats returns summary statistics
func (s *MetricsStore) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	successCount := 0
	failureCount := 0
	for _, e := range s.Events {
		if e.Status == "success" {
			successCount++
		} else if e.Status == "failure" {
			failureCount++
		}
	}

	activeAlerts := 0
	for _, a := range s.Alerts {
		if !a.Resolved {
			activeAlerts++
		}
	}

	return map[string]interface{}{
		"total_metrics":  len(s.Metrics),
		"total_events":   len(s.Events),
		"total_alerts":   len(s.Alerts),
		"active_alerts":  activeAlerts,
		"successful_ops": successCount,
		"failed_ops":     failureCount,
		"timestamp":      time.Now(),
	}
}
