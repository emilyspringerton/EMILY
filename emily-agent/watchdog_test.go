package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func touchLog(t *testing.T, path string, mtime time.Time) {
	t.Helper()
	if err := os.WriteFile(path, []byte("log line\n"), 0o644); err != nil {
		t.Fatalf("write log: %v", err)
	}
	if err := os.Chtimes(path, mtime, mtime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
}

// TestCheckPollerHealth_FreshLogNoAlert covers the common case: a poller
// whose log was just written to must not be flagged.
func TestCheckPollerHealth_FreshLogNoAlert(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "secwatch.log")
	touchLog(t, logPath, time.Now())

	pollers := []PollerConfig{{"secwatch", logPath, 30 * time.Minute, "test"}}
	alerts := CheckPollerHealth(dir, pollers)
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts for a fresh log, got: %v", alerts)
	}
}

// TestCheckPollerHealth_StaleLogEventuallyAlerts covers the real bug this
// was written for: a poller whose log has gone stale (process silently
// killed) must be detected and, after the debounce window, alerted.
func TestCheckPollerHealth_StaleLogEventuallyAlerts(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "secwatch.log")
	// Last write was well beyond MaxStaleness ago.
	touchLog(t, logPath, time.Now().Add(-2*time.Hour))

	pollers := []PollerConfig{{"secwatch", logPath, 30 * time.Minute, "test"}}

	// First check: detects staleness, starts the downtime clock, but
	// doesn't alert yet (matches CheckServiceHealth's own debounce shape).
	alerts := CheckPollerHealth(dir, pollers)
	if len(alerts) != 0 {
		t.Fatalf("expected no alert on first detection, got: %v", alerts)
	}

	// Simulate serviceDownThreshold having already elapsed by directly
	// backdating the persisted DownSince, the same way a real second cron
	// cycle several minutes later would observe it.
	ws := loadWatchdogState(dir)
	rec := ws.Services["poller:secwatch"]
	rec.DownSince = time.Now().UTC().Add(-serviceDownThreshold - time.Second)
	ws.Services["poller:secwatch"] = rec
	saveWatchdogState(dir, ws)

	alerts = CheckPollerHealth(dir, pollers)
	if len(alerts) != 1 {
		t.Fatalf("expected exactly 1 alert once past the debounce threshold, got: %v", alerts)
	}
}

// TestCheckPollerHealth_MissingLogTreatedAsDown covers a poller that has
// never written a log at all (e.g. fresh box, log rotated away) — must be
// treated the same as stale, not silently ignored.
func TestCheckPollerHealth_MissingLogTreatedAsDown(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "does-not-exist.log")

	pollers := []PollerConfig{{"ghost-poller", logPath, 30 * time.Minute, "test"}}
	CheckPollerHealth(dir, pollers) // first detection

	ws := loadWatchdogState(dir)
	if _, tracked := ws.Services["poller:ghost-poller"]; !tracked {
		t.Fatal("expected a missing log file to be tracked as down, not ignored")
	}
}

// TestCheckPollerHealth_RecoveryClearsState confirms a poller that resumes
// writing to its log clears its downtime record instead of alerting forever.
func TestCheckPollerHealth_RecoveryClearsState(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "secwatch.log")
	touchLog(t, logPath, time.Now().Add(-2*time.Hour))

	pollers := []PollerConfig{{"secwatch", logPath, 30 * time.Minute, "test"}}
	CheckPollerHealth(dir, pollers) // marks it down

	// The process comes back and writes a fresh line.
	touchLog(t, logPath, time.Now())
	CheckPollerHealth(dir, pollers)

	ws := loadWatchdogState(dir)
	if _, stillTracked := ws.Services["poller:secwatch"]; stillTracked {
		t.Fatal("expected downtime record to clear once the log is fresh again")
	}
}
