package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/yourusername/ops-tool/pkg/metrics"
)

// Dashboard serves the ops-tool web dashboard
type Dashboard struct {
	addr         string
	metricsStore *metrics.MetricsStore
	mu           sync.RWMutex
	isRunning    bool
	stopChan     chan bool
	refreshRate  time.Duration
}

// NewDashboard creates a new dashboard
func NewDashboard(addr string, store *metrics.MetricsStore) *Dashboard {
	return &Dashboard{
		addr:         addr,
		metricsStore: store,
		refreshRate:  5 * time.Second,
		stopChan:     make(chan bool),
	}
}

// Start starts the dashboard server
func (d *Dashboard) Start() error {
	d.mu.Lock()
	if d.isRunning {
		d.mu.Unlock()
		return fmt.Errorf("dashboard already running")
	}
	d.isRunning = true
	d.mu.Unlock()

	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/metrics", d.handleMetrics)
	mux.HandleFunc("/api/events", d.handleEvents)
	mux.HandleFunc("/api/alerts", d.handleAlerts)
	mux.HandleFunc("/api/stats", d.handleStats)
	mux.HandleFunc("/api/health", d.handleHealth)

	// Web UI
	mux.HandleFunc("/", d.handleDashboard)
	mux.HandleFunc("/css/dashboard.css", d.handleCSS)
	mux.HandleFunc("/js/dashboard.js", d.handleJS)

	server := &http.Server{
		Addr:    d.addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Dashboard server error: %v\n", err)
		}
	}()

	fmt.Printf("✅ Dashboard started at http://%s\n", d.addr)

	// Wait for stop signal
	<-d.stopChan
	server.Close()

	d.mu.Lock()
	d.isRunning = false
	d.mu.Unlock()

	return nil
}

// Stop stops the dashboard server
func (d *Dashboard) Stop() {
	d.mu.Lock()
	if d.isRunning {
		d.stopChan <- true
	}
	d.mu.Unlock()
}

// handleDashboard serves the main dashboard HTML
func (d *Dashboard) handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, dashboardHTML)
}

// handleCSS serves dashboard CSS
func (d *Dashboard) handleCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, dashboardCSS)
}

// handleJS serves dashboard JavaScript
func (d *Dashboard) handleJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, dashboardJS)
}

// handleMetrics returns metrics as JSON
func (d *Dashboard) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	metrics := d.metricsStore.GetMetrics()
	json.NewEncoder(w).Encode(metrics)
}

// handleEvents returns events as JSON
func (d *Dashboard) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	events := d.metricsStore.GetEvents()
	json.NewEncoder(w).Encode(events)
}

// handleAlerts returns alerts as JSON
func (d *Dashboard) handleAlerts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	alerts := d.metricsStore.GetAlerts()
	json.NewEncoder(w).Encode(alerts)
}

// handleStats returns summary statistics
func (d *Dashboard) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stats := d.metricsStore.GetStats()
	json.NewEncoder(w).Encode(stats)
}

// handleHealth returns dashboard health status
func (d *Dashboard) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"uptime":    time.Since(time.Now().Add(-1 * time.Hour)), // Example uptime
	}
	json.NewEncoder(w).Encode(response)
}

const dashboardHTML = `
<!DOCTYPE html>
<html>
<head>
	<title>ops-tool Dashboard</title>
	<link rel="stylesheet" href="/css/dashboard.css">
	<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
	<div class="container">
		<header>
			<h1>🚀 ops-tool Dashboard</h1>
			<p>Real-time DevOps Operations Monitoring</p>
		</header>

		<div class="stats-grid">
			<div class="stat-card">
				<h3>Active Alerts</h3>
				<p class="stat-value" id="activeAlerts">-</p>
			</div>
			<div class="stat-card">
				<h3>Total Events</h3>
				<p class="stat-value" id="totalEvents">-</p>
			</div>
			<div class="stat-card">
				<h3>Successful Ops</h3>
				<p class="stat-value success" id="successOps">-</p>
			</div>
			<div class="stat-card">
				<h3>Failed Ops</h3>
				<p class="stat-value error" id="failedOps">-</p>
			</div>
		</div>

		<div class="content-grid">
			<section class="widget alerts-widget">
				<h2>⚠️ Active Alerts</h2>
				<div id="alertsList" class="list">
					<p class="loading">Loading...</p>
				</div>
			</section>

			<section class="widget events-widget">
				<h2>📊 Recent Events</h2>
				<div id="eventsList" class="list">
					<p class="loading">Loading...</p>
				</div>
			</section>

			<section class="widget metrics-widget">
				<h2>📈 Metrics</h2>
				<div id="metricsList" class="list">
					<p class="loading">Loading...</p>
				</div>
			</section>
		</div>

		<footer>
			<p>ops-tool v1.0 | Last updated: <span id="lastUpdate">-</span></p>
		</footer>
	</div>

	<script src="/js/dashboard.js"></script>
</body>
</html>
`

const dashboardCSS = `
* {
	margin: 0;
	padding: 0;
	box-sizing: border-box;
}

body {
	font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: #333;
	min-height: 100vh;
}

.container {
	max-width: 1400px;
	margin: 0 auto;
	padding: 20px;
}

header {
	color: white;
	margin-bottom: 40px;
	text-align: center;
}

header h1 {
	font-size: 2.5em;
	margin-bottom: 10px;
}

header p {
	font-size: 1.1em;
	opacity: 0.9;
}

.stats-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
	gap: 20px;
	margin-bottom: 40px;
}

.stat-card {
	background: white;
	padding: 20px;
	border-radius: 8px;
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
	text-align: center;
}

.stat-card h3 {
	font-size: 0.9em;
	color: #666;
	margin-bottom: 10px;
	text-transform: uppercase;
}

.stat-value {
	font-size: 2.5em;
	font-weight: bold;
	color: #667eea;
}

.stat-value.success {
	color: #10b981;
}

.stat-value.error {
	color: #ef4444;
}

.stat-value.warning {
	color: #f59e0b;
}

.content-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
	gap: 20px;
	margin-bottom: 40px;
}

.widget {
	background: white;
	padding: 20px;
	border-radius: 8px;
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.widget h2 {
	font-size: 1.3em;
	margin-bottom: 20px;
	color: #333;
}

.list {
	max-height: 400px;
	overflow-y: auto;
}

.list-item {
	padding: 12px;
	border-left: 4px solid #667eea;
	margin-bottom: 10px;
	background: #f9fafb;
	border-radius: 4px;
}

.list-item.success {
	border-left-color: #10b981;
}

.list-item.error {
	border-left-color: #ef4444;
}

.list-item.warning {
	border-left-color: #f59e0b;
}

.list-item.critical {
	border-left-color: #ef4444;
}

.list-item-title {
	font-weight: bold;
	margin-bottom: 5px;
}

.list-item-desc {
	font-size: 0.9em;
	color: #666;
}

.list-item-time {
	font-size: 0.8em;
	color: #999;
	margin-top: 5px;
}

.loading {
	text-align: center;
	color: #999;
	padding: 20px;
}

footer {
	text-align: center;
	color: white;
	padding: 20px;
	border-top: 1px solid rgba(255, 255, 255, 0.2);
}

@media (max-width: 768px) {
	.stats-grid {
		grid-template-columns: repeat(2, 1fr);
	}
	
	.content-grid {
		grid-template-columns: 1fr;
	}
	
	header h1 {
		font-size: 1.8em;
	}
}
`

const dashboardJS = `
const API_BASE = '/api';
const REFRESH_INTERVAL = 5000;

async function fetchData(endpoint) {
	try {
		const response = await fetch(API_BASE + endpoint);
		return await response.json();
	} catch (error) {
		console.error('Failed to fetch ' + endpoint + ':', error);
		return null;
	}
}

async function updateDashboard() {
	const stats = await fetchData('/stats');
	const alerts = await fetchData('/alerts');
	const events = await fetchData('/events');
	const metrics = await fetchData('/metrics');

	if (stats) {
		document.getElementById('activeAlerts').textContent = stats.active_alerts || 0;
		document.getElementById('totalEvents').textContent = stats.total_events || 0;
		document.getElementById('successOps').textContent = stats.successful_ops || 0;
		document.getElementById('failedOps').textContent = stats.failed_ops || 0;
	}

	if (alerts) {
		const activeAlerts = alerts.filter(a => !a.resolved_at);
		const alertsHTML = activeAlerts.length > 0
			? activeAlerts.slice(-10).map(a => '<div class="list-item ' + a.severity + '"><div class="list-item-title">' + a.name + '</div><div class="list-item-desc">' + a.message + '</div><div class="list-item-time">' + new Date(a.timestamp).toLocaleString() + '</div></div>').join('')
			: '<p class="loading">No active alerts</p>';
		document.getElementById('alertsList').innerHTML = alertsHTML;
	}

	if (events) {
		const eventsHTML = events.length > 0
			? events.slice(-10).reverse().map(e => '<div class="list-item ' + e.status + '"><div class="list-item-title">' + e.message + '</div><div class="list-item-desc">Duration: ' + (e.duration / 1e9).toFixed(2) + 's</div><div class="list-item-time">' + new Date(e.timestamp).toLocaleString() + '</div></div>').join('')
			: '<p class="loading">No events</p>';
		document.getElementById('eventsList').innerHTML = eventsHTML;
	}

	if (metrics) {
		const metricsHTML = metrics.length > 0
			? metrics.slice(-10).map(m => '<div class="list-item"><div class="list-item-title">' + m.name + '</div><div class="list-item-desc">Value: ' + m.value.toFixed(2) + ' ' + m.unit + '</div><div class="list-item-time">' + new Date(m.timestamp).toLocaleString() + '</div></div>').join('')
			: '<p class="loading">No metrics</p>';
		document.getElementById('metricsList').innerHTML = metricsHTML;
	}

	document.getElementById('lastUpdate').textContent = new Date().toLocaleTimeString();
}

// Initial load and setup refresh interval
updateDashboard();
setInterval(updateDashboard, REFRESH_INTERVAL);
`
