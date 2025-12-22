package dashboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/yourusername/ops-tool/pkg/audit"
	authpkg "github.com/yourusername/ops-tool/pkg/auth"
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

// auth resources for HTTP handlers (use same backing files)
var (
	httpTokenStore *authpkg.TokenStore
	httpUserStore  = authpkg.NewUserStore()
)

// role hierarchy for simple checks
var roleLevel = map[authpkg.Role]int{
	authpkg.RoleViewer:   10,
	authpkg.RoleOperator: 20,
	authpkg.RoleAdmin:    30,
}

// authMiddleware wraps handlers and enforces a minimum role. It accepts either
// an Authorization: Bearer <token> header or X-Actor: <username> custom header.
func authMiddleware(minRole authpkg.Role, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor := ""
		// check Authorization header first
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// expect `Bearer <token>`
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
				if token != "" {
					tok, err := httpTokenStore.Validate(token)
					if err != nil {
						if rerr := audit.Record("auth.check", "", r.URL.Path, map[string]any{"allowed": false, "reason": "invalid token"}); rerr != nil {
							fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
						}
						http.Error(w, "invalid token", http.StatusUnauthorized)
						return
					}
					actor = tok.User
				}
			}
		}
		// fallback to X-Actor header
		if actor == "" {
			actor = r.Header.Get("X-Actor")
		}
		if actor == "" {
			if err := audit.Record("auth.check", "", r.URL.Path, map[string]any{"allowed": false, "reason": "no actor"}); err != nil {
				fmt.Fprintf(os.Stderr, "audit record failed: %v\n", err)
			}
			http.Error(w, "missing actor or token", http.StatusUnauthorized)
			return
		}
		u, err := httpUserStore.GetUser(actor)
		if err != nil {
			if rerr := audit.Record("auth.check", actor, r.URL.Path, map[string]any{"allowed": false, "reason": "actor not found"}); rerr != nil {
				fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
			}
			http.Error(w, "actor not found", http.StatusUnauthorized)
			return
		}
		if roleLevel[u.Role] < roleLevel[minRole] {
			if rerr := audit.Record("auth.check", actor, r.URL.Path, map[string]any{"allowed": false, "required": minRole, "have": u.Role}); rerr != nil {
				fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
			}
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		// allowed
		if err := audit.Record("auth.check", actor, r.URL.Path, map[string]any{"allowed": true, "have": u.Role}); err != nil {
			fmt.Fprintf(os.Stderr, "audit record failed: %v\n", err)
		}
		h(w, r)
	}
}

// authMethodMiddleware enforces different minimum roles depending on HTTP method.
// GET/HEAD/OPTIONS use `viewRole`, while POST/PUT/DELETE use `postRole`.
func authMethodMiddleware(viewRole, postRole authpkg.Role, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// pick required role by method
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodDelete:
			authMiddleware(postRole, h)(w, r)
			return
		default:
			authMiddleware(viewRole, h)(w, r)
			return
		}
	}
}

// Start starts the dashboard server
func (d *Dashboard) Start() error {
	// write a debug marker so we can trace startup failures from logs
	func() {
		if f, err := os.OpenFile("dashboard.start.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
			if _, werr := fmt.Fprintf(f, "%s: Start() called for %s\n", time.Now().Format(time.RFC3339), d.addr); werr != nil {
				log.Printf("failed to write dashboard.start.log: %v", werr)
			}
			if cerr := f.Close(); cerr != nil {
				log.Printf("failed to close dashboard.start.log: %v", cerr)
			}
		}
	}()

	d.mu.Lock()
	if d.isRunning {
		d.mu.Unlock()
		return fmt.Errorf("dashboard already running")
	}
	d.isRunning = true
	d.mu.Unlock()

	// ensure token store is initialized and watcher started for the dashboard
	if httpTokenStore == nil {
		httpTokenStore = authpkg.NewTokenStore("tokens.json")
		httpTokenStore.StartWatcher()
	}

	mux := http.NewServeMux()

	// API endpoints (viewer for reads, operator for creates/changes)
	mux.HandleFunc("/api/metrics", authMethodMiddleware(authpkg.RoleViewer, authpkg.RoleOperator, d.handleMetrics))
	mux.HandleFunc("/api/events", authMethodMiddleware(authpkg.RoleViewer, authpkg.RoleOperator, d.handleEvents))
	mux.HandleFunc("/api/alerts", authMethodMiddleware(authpkg.RoleViewer, authpkg.RoleOperator, d.handleAlerts))
	// audit is sensitive: admin only
	mux.HandleFunc("/api/audit", authMiddleware(authpkg.RoleAdmin, d.handleAudit))
	mux.HandleFunc("/api/stats", authMethodMiddleware(authpkg.RoleViewer, authpkg.RoleOperator, d.handleStats))
	mux.HandleFunc("/api/health", authMiddleware(authpkg.RoleViewer, d.handleHealth))

	// Web UI (require viewer)
	mux.HandleFunc("/", authMiddleware(authpkg.RoleViewer, d.handleDashboard))
	mux.HandleFunc("/css/dashboard.css", authMiddleware(authpkg.RoleViewer, d.handleCSS))
	mux.HandleFunc("/js/dashboard.js", authMiddleware(authpkg.RoleViewer, d.handleJS))

	server := &http.Server{
		Addr:    d.addr,
		Handler: mux,
	}

	// bind first to surface address/port errors immediately
	ln, err := net.Listen("tcp", d.addr)
	if err != nil {
		// persist the error to a local file for debugging
		f, ferr := os.OpenFile("dashboard.error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if ferr == nil {
			if _, werr := fmt.Fprintf(f, "%s: bind error: %v\n", time.Now().Format(time.RFC3339), err); werr != nil {
				log.Printf("failed to write dashboard.error.log: %v", werr)
			}
			if cerr := f.Close(); cerr != nil {
				log.Printf("failed to close dashboard.error.log: %v", cerr)
			}
		}
		if rerr := audit.Record("dashboard.error", "", d.addr, map[string]any{"error": err.Error()}); rerr != nil {
			log.Printf("audit record failed: %v", rerr)
		}
		return err
	}

	// record successful bind to help debugging
	if f, ferr := os.OpenFile("dashboard.start.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); ferr == nil {
		if _, werr := fmt.Fprintf(f, "%s: bind OK on %s\n", time.Now().Format(time.RFC3339), d.addr); werr != nil {
			log.Printf("failed to write dashboard.start.log: %v", werr)
		}
		if cerr := f.Close(); cerr != nil {
			log.Printf("failed to close dashboard.start.log: %v", cerr)
		}
	}

	// Run server.Serve on the bound listener so bind errors are already checked.
	errCh := make(chan error, 1)
	if f, ferr := os.OpenFile("dashboard.start.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); ferr == nil {
		if _, werr := fmt.Fprintf(f, "%s: launching serve goroutine\n", time.Now().Format(time.RFC3339)); werr != nil {
			log.Printf("failed to write dashboard.start.log: %v", werr)
		}
		if cerr := f.Close(); cerr != nil {
			log.Printf("failed to close dashboard.start.log: %v", cerr)
		}
	}
	go func() {
		errCh <- server.Serve(ln)
	}()

	fmt.Printf("✅ Dashboard started at http://%s\n", d.addr)

	select {
	case err := <-errCh:
		// persist the received err (can be nil or http.ErrServerClosed)
		if f, ferr := os.OpenFile("dashboard.start.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); ferr == nil {
			if _, werr := fmt.Fprintf(f, "%s: errCh returned: %v\n", time.Now().Format(time.RFC3339), err); werr != nil {
				log.Printf("failed to write dashboard.start.log: %v", werr)
			}
			if cerr := f.Close(); cerr != nil {
				log.Printf("failed to close dashboard.start.log: %v", cerr)
			}
		}
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("Dashboard server error: %v\n", err)
			// persist the error to a local file for debugging
			f, ferr := os.OpenFile("dashboard.error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
			if ferr == nil {
				if _, werr := fmt.Fprintf(f, "%s: Dashboard server error: %v\n", time.Now().Format(time.RFC3339), err); werr != nil {
					log.Printf("failed to write dashboard.error.log: %v", werr)
				}
				if cerr := f.Close(); cerr != nil {
					log.Printf("failed to close dashboard.error.log: %v", cerr)
				}
			}
			if rerr := audit.Record("dashboard.error", "", d.addr, map[string]any{"error": err.Error()}); rerr != nil {
				log.Printf("audit record failed: %v", rerr)
			}
			return err
		}
	case <-d.stopChan:
		if cerr := server.Close(); cerr != nil {
			log.Printf("server close failed: %v", cerr)
		}
	}

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
	if _, err := fmt.Fprint(w, dashboardHTML); err != nil {
		fmt.Fprintf(os.Stderr, "write response failed: %v\n", err)
	}
}

// handleCSS serves dashboard CSS
func (d *Dashboard) handleCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, dashboardCSS); err != nil {
		fmt.Fprintf(os.Stderr, "write response failed: %v\n", err)
	}
}

// handleJS serves dashboard JavaScript
func (d *Dashboard) handleJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, dashboardJS); err != nil {
		fmt.Fprintf(os.Stderr, "write response failed: %v\n", err)
	}
}

// handleMetrics returns metrics as JSON
func (d *Dashboard) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	metrics := d.metricsStore.GetMetrics()
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, "failed to encode metrics", http.StatusInternalServerError)
		if rerr := audit.Record("dashboard.error", "", d.addr, map[string]any{"error": err.Error()}); rerr != nil {
			fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
		}
		return
	}
}

// handleEvents returns events as JSON
func (d *Dashboard) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	events := d.metricsStore.GetEvents()
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "failed to encode events", http.StatusInternalServerError)
		if rerr := audit.Record("dashboard.error", "", d.addr, map[string]any{"error": err.Error()}); rerr != nil {
			fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
		}
		return
	}
}

// handleAlerts returns alerts as JSON
func (d *Dashboard) handleAlerts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	alerts := d.metricsStore.GetAlerts()
	if err := json.NewEncoder(w).Encode(alerts); err != nil {
		http.Error(w, "failed to encode alerts", http.StatusInternalServerError)
		if rerr := audit.Record("dashboard.error", "", d.addr, map[string]any{"error": err.Error()}); rerr != nil {
			fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
		}
		return
	}
}

// handleStats returns summary statistics
func (d *Dashboard) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stats := d.metricsStore.GetStats()
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "failed to encode stats", http.StatusInternalServerError)
		if rerr := audit.Record("dashboard.error", "", d.addr, map[string]any{"error": err.Error()}); rerr != nil {
			fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
		}
		return
	}
}

// handleHealth returns dashboard health status
func (d *Dashboard) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"uptime":    time.Since(time.Now().Add(-1 * time.Hour)), // Example uptime
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode health", http.StatusInternalServerError)
		if rerr := audit.Record("dashboard.error", "", d.addr, map[string]any{"error": err.Error()}); rerr != nil {
			fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
		}
		return
	}
}

// handleAudit returns audit log entries as JSON
func (d *Dashboard) handleAudit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	entries, err := audit.ReadEntries("")
	if err != nil {
		http.Error(w, "failed to read audit", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, "failed to encode audit entries", http.StatusInternalServerError)
		if rerr := audit.Record("dashboard.error", "", d.addr, map[string]any{"error": err.Error()}); rerr != nil {
			fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
		}
		return
	}
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
