// emily-agent/watchdog.go
// Service health watchdog for Emily Prime.
//
// CheckServiceHealth is called at the start of each cron cycle (RunOnce).
// It pings critical services and fires escalation Apples when:
//   - A service is down for >= 2 consecutive minutes
//   - A log file in EMILY/var/logs/ exceeds 500 MB
//
// Downtime tracking persists across cycles in a JSON file (watchdog-state.json).

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	watchdogStateFile    = "watchdog-state.json"
	serviceDownThreshold = 2 * time.Minute
	logSizeAlertBytes    = 500 * 1024 * 1024 // 500 MB
	watchdogHTTPTimeout  = 5 * time.Second
)

// ServiceConfig describes one monitored service.
type ServiceConfig struct {
	Name        string
	HealthURL   string // URL to GET; 2xx = up
	Description string
}

// defaultServices returns the standard EINHORN_INDUSTRIAL service set.
func defaultServices() []ServiceConfig {
	fatbabyRoot := envOr("FATBABY_ROOT", "/home/fatbaby/PRRJECT_FATBABY")
	_ = fatbabyRoot
	shankpitAdmin := envOr("SHANKPIT_ADMIN_URL", "http://localhost:6970")
	return []ServiceConfig{
		{"IDUNA", "http://localhost:8080/health", "IAM + trust authority"},
		{"newssite", "http://localhost:8082/healthz", "FatBaby news site"},
		{"signalapi", "http://localhost:9091/healthz", "Signal query API"},
		{"emily-agent", "http://localhost:8086/api/v1/emily/state", "Emily Prime agent"},
		{"shankpit", shankpitAdmin + "/healthz", "SHANKPIT game server"},
	}
}

// WatchdogState tracks per-service downtime.
type WatchdogState struct {
	Services  map[string]ServiceDownRecord `json:"services"`
	UpdatedAt time.Time                    `json:"updated_at"`
}

// ServiceDownRecord tracks when a service first went down.
type ServiceDownRecord struct {
	Name      string    `json:"name"`
	DownSince time.Time `json:"down_since"`
	Alerted   bool      `json:"alerted"` // true once we've fired an Apple for this outage
}

// loadWatchdogState reads persisted watchdog state from the state dir.
func loadWatchdogState(stateDir string) *WatchdogState {
	path := filepath.Join(stateDir, watchdogStateFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return &WatchdogState{Services: make(map[string]ServiceDownRecord)}
	}
	var ws WatchdogState
	if err := json.Unmarshal(data, &ws); err != nil {
		return &WatchdogState{Services: make(map[string]ServiceDownRecord)}
	}
	if ws.Services == nil {
		ws.Services = make(map[string]ServiceDownRecord)
	}
	return &ws
}

func saveWatchdogState(stateDir string, ws *WatchdogState) {
	ws.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return
	}
	path := filepath.Join(stateDir, watchdogStateFile)
	_ = os.WriteFile(path, data, 0o644)
}

// CheckServiceHealth pings all services and returns a list of alert messages.
// Each alert message is suitable for an escalation Apple body.
// Also checks log file sizes and alerts if any exceeds 500 MB.
func CheckServiceHealth(ctx context.Context, stateDir, logDir string, services []ServiceConfig) []string {
	if services == nil {
		services = defaultServices()
	}

	ws := loadWatchdogState(stateDir)
	var alerts []string
	client := &http.Client{Timeout: watchdogHTTPTimeout}
	now := time.Now().UTC()

	for _, svc := range services {
		up := pingService(ctx, client, svc.HealthURL)

		if up {
			// Service recovered — clear downtime record.
			if rec, was_down := ws.Services[svc.Name]; was_down {
				downDuration := now.Sub(rec.DownSince).Round(time.Second)
				log.Printf("[watchdog] %s recovered (was down %s)", svc.Name, downDuration)
			}
			delete(ws.Services, svc.Name)
			continue
		}

		// Service is down.
		rec, tracked := ws.Services[svc.Name]
		if !tracked {
			// First time we see it down — record the time.
			ws.Services[svc.Name] = ServiceDownRecord{
				Name:      svc.Name,
				DownSince: now,
			}
			log.Printf("[watchdog] %s DOWN (first detection)", svc.Name)
			continue
		}

		downDuration := now.Sub(rec.DownSince)
		log.Printf("[watchdog] %s still DOWN (%.0fs, threshold=%s, alerted=%v)",
			svc.Name, downDuration.Seconds(), serviceDownThreshold, rec.Alerted)

		if downDuration >= serviceDownThreshold && !rec.Alerted {
			msg := fmt.Sprintf("%s has been DOWN for %s. Health URL: %s. Description: %s.",
				svc.Name, downDuration.Round(time.Second), svc.HealthURL, svc.Description)
			alerts = append(alerts, msg)
			rec.Alerted = true
			ws.Services[svc.Name] = rec
		}
	}

	// Check log file sizes.
	if logDir != "" {
		logAlerts := checkLogSizes(logDir)
		alerts = append(alerts, logAlerts...)
	}

	saveWatchdogState(stateDir, ws)
	return alerts
}

// pingService returns true if the URL responds with a 2xx status.
func pingService(ctx context.Context, client *http.Client, url string) bool {
	if url == "" {
		return true // unconfigured services are not monitored
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// PollerConfig describes one monitored headless poller process — one with
// no HTTP endpoint to ping (secwatch, prwatch, eps-reconciler, ...),
// monitored instead by log-file freshness: is the process still writing to
// its own log within the window its poll interval implies?
//
// Added 2026-07-17 after secwatch and eps-reconciler were both silently
// killed (OOM) and stayed down for hours — secwatch's data gap ran from
// ~13:00 until a human happened to run a manual "check on all the fatbaby
// data" audit at ~23:00 — with zero automated detection, because
// CheckServiceHealth's HTTP-only monitoring has no visibility into a
// process with no health endpoint. The founder's own framing: FatBaby's
// data is the multi-year asset this whole system exists to build (data
// resale readiness, years of continuous operation) — a silent day-long
// ingestion gap discovered by luck, not alerting, is exactly the failure
// mode that framing can't tolerate.
type PollerConfig struct {
	Name         string
	LogPath      string
	MaxStaleness time.Duration // alert if the log hasn't been written to within this window
	Description  string
}

// defaultPollers returns the core FatBaby ingestion processes that have no
// HTTP health endpoint. MaxStaleness is set generously above each process's
// own poll interval (see each process's -poll-interval default / ops/
// systemd unit) to avoid false alarms from a merely-slow cycle.
func defaultPollers() []PollerConfig {
	fatbabyRoot := envOr("FATBABY_ROOT", "/home/fatbaby/PRRJECT_FATBABY")
	logDir := filepath.Join(fatbabyRoot, "var", "logs")
	return []PollerConfig{
		{"secwatch", filepath.Join(logDir, "secwatch.log"), 30 * time.Minute, "SEC EDGAR filing poller (5m interval) — the primary data-ingestion process"},
		{"processor", filepath.Join(logDir, "processor.log"), 15 * time.Minute, "Filing→signal processor (15s interval)"},
		{"prwatch", filepath.Join(logDir, "prwatch.log"), 15 * time.Minute, "PR Newswire discovery poller (30s interval)"},
		{"prwatch-body", filepath.Join(logDir, "prwatch-body.log"), 15 * time.Minute, "PR Newswire body fetcher (15s interval)"},
		{"eps-reconciler", filepath.Join(logDir, "eps-reconciler.log"), 8 * time.Hour, "EPS oracle reconciliation (6h interval)"},
	}
}

// CheckPollerHealth checks each poller's log-file freshness and returns
// alert messages using the same downtime-tracking/debounce state
// (WatchdogState) and 2-minute confirmation threshold as CheckServiceHealth,
// just keyed by "poller:<name>" so the two check types never collide.
func CheckPollerHealth(stateDir string, pollers []PollerConfig) []string {
	if pollers == nil {
		pollers = defaultPollers()
	}

	ws := loadWatchdogState(stateDir)
	var alerts []string
	now := time.Now().UTC()

	for _, p := range pollers {
		key := "poller:" + p.Name
		info, err := os.Stat(p.LogPath)
		fresh := err == nil && now.Sub(info.ModTime()) < p.MaxStaleness

		if fresh {
			if rec, wasDown := ws.Services[key]; wasDown {
				downDuration := now.Sub(rec.DownSince).Round(time.Second)
				log.Printf("[watchdog] poller %s recovered (was stale %s)", p.Name, downDuration)
			}
			delete(ws.Services, key)
			continue
		}

		rec, tracked := ws.Services[key]
		if !tracked {
			ws.Services[key] = ServiceDownRecord{Name: key, DownSince: now}
			log.Printf("[watchdog] poller %s STALE (first detection — log not updated within %s)", p.Name, p.MaxStaleness)
			continue
		}

		downDuration := now.Sub(rec.DownSince)
		log.Printf("[watchdog] poller %s still STALE (%.0fs since first detection, alerted=%v)",
			p.Name, downDuration.Seconds(), rec.Alerted)

		if downDuration >= serviceDownThreshold && !rec.Alerted {
			msg := fmt.Sprintf(
				"%s log has not been updated in over %s (expected activity within %s) — likely dead or hung. Log: %s. Description: %s.",
				p.Name, (p.MaxStaleness + downDuration).Round(time.Minute), p.MaxStaleness, p.LogPath, p.Description)
			alerts = append(alerts, msg)
			rec.Alerted = true
			ws.Services[key] = rec
		}
	}

	saveWatchdogState(stateDir, ws)
	return alerts
}

// checkLogSizes scans logDir for files exceeding logSizeAlertBytes.
// Returns alert strings for each oversized file.
func checkLogSizes(logDir string) []string {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return nil
	}

	var alerts []string
	var oversized []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".log") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.Size() > logSizeAlertBytes {
			oversized = append(oversized, fmt.Sprintf("%s (%.0f MB)",
				name, float64(info.Size())/1024/1024))
		}
	}

	if len(oversized) > 0 {
		sort.Strings(oversized)
		alerts = append(alerts, fmt.Sprintf(
			"Log files exceeding 500 MB in %s: %s. Run logrotate or truncate.",
			logDir, strings.Join(oversized, ", ")))
	}
	return alerts
}

// CheckpointConfig describes one SQLite index checkpoint (see PRRJECT_FATBABY
// internal/indexcheckpoint and cmd/entity-graph/{filingindex,accuracyindex}.go)
// whose meta.snapshot_at value should advance roughly every poll interval.
//
// Added for docs/northstar/replay-fragility.md §4c/§9 Phase 3: these
// checkpoints are a disposable cache of the event store (the invariant is
// "deleting one is always a safe operator move" -- see the ops runbook), but
// if a checkpoint stops advancing while its owning process still looks alive
// by CheckServiceHealth/CheckPollerHealth's own signals, that is a distinct
// failure mode (write path broken, batch wedged mid-checkpoint-write) worth
// catching on its own, per Principle 15 -- in minutes, not at the next
// incident.
type CheckpointConfig struct {
	Name         string
	DBPath       string
	MaxStaleness time.Duration // alert if snapshot_at hasn't advanced within this window
	Description  string
}

// defaultCheckpoints returns the FatBaby SQLite index checkpoints. MaxStaleness
// is set generously above each owning process's poll interval (30s for all
// four today) to avoid false alarms from a merely-slow cycle.
func defaultCheckpoints() []CheckpointConfig {
	fatbabyRoot := envOr("FATBABY_ROOT", "/home/fatbaby/PRRJECT_FATBABY")
	varDir := filepath.Join(fatbabyRoot, "var")
	return []CheckpointConfig{
		{"signalapi", filepath.Join(varDir, "signalapi-index.db"), 5 * time.Minute, "signalapi SQLite index checkpoint (signals+docs)"},
		{"newssite", filepath.Join(varDir, "newssite-index.db"), 5 * time.Minute, "newssite SQLite index checkpoint (signals+docs)"},
		{"entity-graph-filings", filepath.Join(varDir, "entity-graph", "filings-index.db"), 5 * time.Minute, "entity-graph incremental filing-date/form index"},
		{"entity-graph-accuracy", filepath.Join(varDir, "entity-graph", "accuracy-index.db"), 5 * time.Minute, "entity-graph deduplicated accuracy-verdict index"},
	}
}

// checkpointSnapshotAt opens dbPath read-only and returns its meta.snapshot_at
// value. A missing file, missing row, or unparsable timestamp all return an
// error -- callers treat any error the same as "stale" (matches
// CheckPollerHealth treating a missing log file as stale, not as healthy).
func checkpointSnapshotAt(dbPath string) (time.Time, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return time.Time{}, err
	}
	// mode=ro: this process only ever reads -- never compete with the
	// owning process (the sole writer) for the write lock.
	db, err := sql.Open("sqlite", dbPath+"?mode=ro&_busy_timeout=2000")
	if err != nil {
		return time.Time{}, err
	}
	defer db.Close()
	var raw string
	if err := db.QueryRow(`SELECT value FROM meta WHERE key='snapshot_at'`).Scan(&raw); err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, raw)
}

// CheckCheckpointHealth checks each checkpoint's meta.snapshot_at freshness
// and returns alert messages, using the same downtime-tracking/debounce
// state (WatchdogState) and 2-minute confirmation threshold as
// CheckServiceHealth/CheckPollerHealth, keyed by "checkpoint:<name>" so none
// of the three check types collide.
func CheckCheckpointHealth(stateDir string, checkpoints []CheckpointConfig) []string {
	if checkpoints == nil {
		checkpoints = defaultCheckpoints()
	}

	ws := loadWatchdogState(stateDir)
	var alerts []string
	now := time.Now().UTC()

	for _, c := range checkpoints {
		key := "checkpoint:" + c.Name
		snapAt, err := checkpointSnapshotAt(c.DBPath)
		fresh := err == nil && now.Sub(snapAt) < c.MaxStaleness

		if fresh {
			if rec, wasDown := ws.Services[key]; wasDown {
				downDuration := now.Sub(rec.DownSince).Round(time.Second)
				log.Printf("[watchdog] checkpoint %s recovered (was stale %s)", c.Name, downDuration)
			}
			delete(ws.Services, key)
			continue
		}

		rec, tracked := ws.Services[key]
		if !tracked {
			ws.Services[key] = ServiceDownRecord{Name: key, DownSince: now}
			log.Printf("[watchdog] checkpoint %s STALE (first detection -- snapshot_at not advancing within %s)", c.Name, c.MaxStaleness)
			continue
		}

		downDuration := now.Sub(rec.DownSince)
		log.Printf("[watchdog] checkpoint %s still STALE (%.0fs since first detection, alerted=%v)",
			c.Name, downDuration.Seconds(), rec.Alerted)

		if downDuration >= serviceDownThreshold && !rec.Alerted {
			msg := fmt.Sprintf(
				"Checkpoint %s (%s) has not advanced in over %s (expected within %s) -- its owning process may be alive but its checkpoint write path is wedged. Description: %s. Checkpoints are a disposable cache of the event store -- deleting %s is always a safe operator move; the owning process rebuilds it on next start.",
				c.Name, c.DBPath, (c.MaxStaleness + downDuration).Round(time.Minute), c.MaxStaleness, c.Description, c.DBPath)
			alerts = append(alerts, msg)
			rec.Alerted = true
			ws.Services[key] = rec
		}
	}

	saveWatchdogState(stateDir, ws)
	return alerts
}
