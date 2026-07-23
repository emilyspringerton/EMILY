package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
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

// writeCheckpointDB creates a minimal SQLite file with a meta.snapshot_at row
// set to snapAt, mirroring the schema PRRJECT_FATBABY's indexcheckpoint and
// entity-graph's filingindex/accuracyindex packages actually write.
func writeCheckpointDB(t *testing.T, path string, snapAt time.Time) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open checkpoint db: %v", err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE meta (key TEXT PRIMARY KEY, value TEXT)`); err != nil {
		t.Fatalf("create meta table: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO meta(key, value) VALUES ('snapshot_at', ?)`, snapAt.UTC().Format(time.RFC3339)); err != nil {
		t.Fatalf("insert snapshot_at: %v", err)
	}
}

// TestCheckCheckpointHealth_FreshSnapshotNoAlert covers the common case: a
// checkpoint whose snapshot_at was just written must not be flagged.
func TestCheckCheckpointHealth_FreshSnapshotNoAlert(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "signalapi-index.db")
	writeCheckpointDB(t, dbPath, time.Now())

	checkpoints := []CheckpointConfig{{"signalapi", dbPath, 5 * time.Minute, "test"}}
	alerts := CheckCheckpointHealth(dir, checkpoints)
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts for a fresh checkpoint, got: %v", alerts)
	}
}

// TestCheckCheckpointHealth_StaleSnapshotEventuallyAlerts covers the real gap
// Phase 3 closes: a checkpoint whose snapshot_at has stopped advancing (write
// path wedged, even though the owning process might still look alive to
// CheckServiceHealth/CheckPollerHealth) must be detected and, after the
// debounce window, alerted.
func TestCheckCheckpointHealth_StaleSnapshotEventuallyAlerts(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "entity-graph-filings.db")
	writeCheckpointDB(t, dbPath, time.Now().Add(-2*time.Hour))

	checkpoints := []CheckpointConfig{{"entity-graph-filings", dbPath, 5 * time.Minute, "test"}}

	alerts := CheckCheckpointHealth(dir, checkpoints)
	if len(alerts) != 0 {
		t.Fatalf("expected no alert on first detection, got: %v", alerts)
	}

	ws := loadWatchdogState(dir)
	rec := ws.Services["checkpoint:entity-graph-filings"]
	rec.DownSince = time.Now().UTC().Add(-serviceDownThreshold - time.Second)
	ws.Services["checkpoint:entity-graph-filings"] = rec
	saveWatchdogState(dir, ws)

	alerts = CheckCheckpointHealth(dir, checkpoints)
	if len(alerts) != 1 {
		t.Fatalf("expected exactly 1 alert once past the debounce threshold, got: %v", alerts)
	}
}

// TestCheckCheckpointHealth_MissingFileTreatedAsStale covers a checkpoint
// that was deleted (an explicitly supported operator move -- checkpoints are
// a disposable cache) or never created yet -- must be treated as stale, not
// silently ignored, same as CheckPollerHealth's missing-log case.
func TestCheckCheckpointHealth_MissingFileTreatedAsStale(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "does-not-exist.db")

	checkpoints := []CheckpointConfig{{"ghost-checkpoint", dbPath, 5 * time.Minute, "test"}}
	CheckCheckpointHealth(dir, checkpoints) // first detection

	ws := loadWatchdogState(dir)
	if _, tracked := ws.Services["checkpoint:ghost-checkpoint"]; !tracked {
		t.Fatal("expected a missing checkpoint file to be tracked as stale, not ignored")
	}
}

// TestCheckCheckpointHealth_RecoveryClearsState confirms a checkpoint whose
// snapshot_at advances again (e.g. after a delete-and-rebuild, or the process
// resuming its poll loop) clears its downtime record instead of alerting
// forever.
func TestCheckCheckpointHealth_RecoveryClearsState(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "signalapi-index.db")
	writeCheckpointDB(t, dbPath, time.Now().Add(-2*time.Hour))

	checkpoints := []CheckpointConfig{{"signalapi", dbPath, 5 * time.Minute, "test"}}
	CheckCheckpointHealth(dir, checkpoints) // marks it down

	// The process (or a rebuild after an operator delete) writes a fresh
	// snapshot.
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("remove stale checkpoint: %v", err)
	}
	writeCheckpointDB(t, dbPath, time.Now())
	CheckCheckpointHealth(dir, checkpoints)

	ws := loadWatchdogState(dir)
	if _, stillTracked := ws.Services["checkpoint:signalapi"]; stillTracked {
		t.Fatal("expected downtime record to clear once snapshot_at is fresh again")
	}
}
