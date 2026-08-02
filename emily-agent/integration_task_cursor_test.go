package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestRunPrimeTriage_DoesNotReissueTasksAcrossCycles is a regression test for a real,
// long-running bug (found live 2026-08-02, founder: "the rsi loop keeps putting the same
// stale 7ish tasks into the backlog"): runPrimeTriage re-read the same ReadObservations(10)
// window every cycle with no cursor tracking which observations had already been triaged for
// tasking, so a single observation that stayed near the top of that window (i.e. no newer
// observations arrived to push it out) generated a fresh duplicate task file every cycle
// forever. WriteTask's own recentDuplicateExists only rate-limited this to once per ~4h12m
// (its 4h rolling window plus 15-min poll granularity), not stopped it -- 260 duplicate copies
// of the same 9 tasks had accumulated in signals/tasks/ over roughly 45 days before this was
// caught. This test proves the fix: two full runPrimeTriage cycles over the same unchanged
// observation set produce exactly one task file, not two.
func TestRunPrimeTriage_DoesNotReissueTasksAcrossCycles(t *testing.T) {
	tmp := t.TempDir()
	signalsDir := filepath.Join(tmp, "signals")
	gitRoot := filepath.Join(tmp, "gitroot")

	store, err := NewIntegrationStore(signalsDir, gitRoot)
	if err != nil {
		t.Fatalf("NewIntegrationStore: %v", err)
	}
	store.fatBabyRoot = "" // disable backlog curation side effects for this test

	// Minimal valid signal-priorities.json — LoadSignalPriorities requires this file to exist.
	if err := os.MkdirAll(filepath.Join(gitRoot, "context"), 0o755); err != nil {
		t.Fatal(err)
	}
	priorities := `{
		"version": 1,
		"updated_at": "2026-08-02T00:00:00Z",
		"signal_weights": {},
		"ceo_visibility_threshold": 0.95,
		"escalation_criteria": {}
	}`
	if err := os.WriteFile(filepath.Join(gitRoot, "context", "signal-priorities.json"), []byte(priorities), 0o644); err != nil {
		t.Fatal(err)
	}

	// A single observation guaranteed to triage above the 0.5 task-write threshold:
	// Severity "critical" alone sets maxWeight=0.9 in Triage, independent of signal_weights.
	obs := Observation{
		Timestamp: "2026-08-02T00:00:00Z",
		Source:    "fatbaby-emily",
		Severity:  "critical",
		Summary:   "test observation that should only ever produce one task",
	}
	if err := store.WriteObservation(obs); err != nil {
		t.Fatalf("WriteObservation: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := runPrimeTriage(ctx, store, nil, nil, nil); err != nil {
		t.Fatalf("runPrimeTriage (cycle 1): %v", err)
	}
	tasksAfterFirst := countTaskFiles(t, signalsDir)
	if tasksAfterFirst != 1 {
		t.Fatalf("after cycle 1: want 1 task file, got %d", tasksAfterFirst)
	}

	// Backdate the task file written by cycle 1 past WriteTask's own 4h recentDuplicateExists
	// window (renaming it to an older timestamp-prefixed filename, matching how that dedup
	// actually reads recency). This reproduces the real bug's real trigger condition -- without
	// this, cycle 2 would land inside the 4h window and get suppressed by that OTHER, adjacent
	// dedup layer regardless of whether the cursor fix under test here does anything at all,
	// which would make this test pass for the wrong reason.
	backdateTaskFiles(t, signalsDir, -5*time.Hour)

	// Second cycle: same observation still present, no new observations arrived, and the
	// existing task is now old enough that WriteTask's own 4h dedup would no longer catch a
	// duplicate on its own. Before this test's fix (the task cursor), this cycle would write a
	// second, duplicate task file (fresh timestamp/task_id, identical description) -- exactly
	// the real, live bug.
	if _, err := runPrimeTriage(ctx, store, nil, nil, nil); err != nil {
		t.Fatalf("runPrimeTriage (cycle 2): %v", err)
	}
	tasksAfterSecond := countTaskFiles(t, signalsDir)
	if tasksAfterSecond != 1 {
		t.Fatalf("after cycle 2 (same observation, no new arrivals, old task past the 4h dedup window): want still 1 task file, got %d -- the task cursor did not prevent re-issuing", tasksAfterSecond)
	}

	// The cursor file itself should now reflect the observation's own timestamp.
	cursorData, err := os.ReadFile(filepath.Join(signalsDir, "observations", ".task-cursor"))
	if err != nil {
		t.Fatalf("read .task-cursor: %v", err)
	}
	if got := string(cursorData); got != obs.Timestamp {
		t.Errorf(".task-cursor = %q, want %q", got, obs.Timestamp)
	}
}

// TestRunPrimeTriage_NewObservationStillGetsTasked confirms the fix doesn't over-suppress:
// a genuinely new observation (newer timestamp than the cursor) still gets a task written.
func TestRunPrimeTriage_NewObservationStillGetsTasked(t *testing.T) {
	tmp := t.TempDir()
	signalsDir := filepath.Join(tmp, "signals")
	gitRoot := filepath.Join(tmp, "gitroot")

	store, err := NewIntegrationStore(signalsDir, gitRoot)
	if err != nil {
		t.Fatalf("NewIntegrationStore: %v", err)
	}
	store.fatBabyRoot = ""

	if err := os.MkdirAll(filepath.Join(gitRoot, "context"), 0o755); err != nil {
		t.Fatal(err)
	}
	priorities := `{
		"version": 1,
		"updated_at": "2026-08-02T00:00:00Z",
		"signal_weights": {},
		"ceo_visibility_threshold": 0.95,
		"escalation_criteria": {}
	}`
	if err := os.WriteFile(filepath.Join(gitRoot, "context", "signal-priorities.json"), []byte(priorities), 0o644); err != nil {
		t.Fatal(err)
	}

	obs1 := Observation{Timestamp: "2026-08-02T00:00:00Z", Source: "fatbaby-emily", Severity: "critical", Summary: "first observation"}
	if err := store.WriteObservation(obs1); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := runPrimeTriage(ctx, store, nil, nil, nil); err != nil {
		t.Fatalf("runPrimeTriage (cycle 1): %v", err)
	}
	if n := countTaskFiles(t, signalsDir); n != 1 {
		t.Fatalf("after cycle 1: want 1 task file, got %d", n)
	}

	obs2 := Observation{Timestamp: "2026-08-02T01:00:00Z", Source: "fatbaby-emily", Severity: "critical", Summary: "second, genuinely new observation"}
	if err := store.WriteObservation(obs2); err != nil {
		t.Fatal(err)
	}
	if _, err := runPrimeTriage(ctx, store, nil, nil, nil); err != nil {
		t.Fatalf("runPrimeTriage (cycle 2): %v", err)
	}
	if n := countTaskFiles(t, signalsDir); n != 2 {
		t.Fatalf("after cycle 2 (one genuinely new observation arrived): want 2 task files, got %d", n)
	}
}

// backdateTaskFiles renames every task file to a filename timestamp `age` in the past.
// recentDuplicateExists reads recency purely from the filename's RFC3339-with-colons-stripped
// prefix (see its own implementation), so renaming is sufficient to simulate the passage of
// time without needing to actually wait or fake the clock.
func backdateTaskFiles(t *testing.T, signalsDir string, age time.Duration) {
	t.Helper()
	dir := filepath.Join(signalsDir, "tasks")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read tasks dir: %v", err)
	}
	oldPrefix := strings.ReplaceAll(time.Now().UTC().Add(age).Format(time.RFC3339), ":", "")
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		// Filenames are "<timestamp>-<task_id>.json"; keep the task_id suffix, replace the
		// timestamp prefix only.
		parts := strings.SplitN(e.Name(), "-task-", 2)
		if len(parts) != 2 {
			t.Fatalf("unexpected task filename shape: %s", e.Name())
		}
		newName := oldPrefix + "-task-" + parts[1]
		if err := os.Rename(filepath.Join(dir, e.Name()), filepath.Join(dir, newName)); err != nil {
			t.Fatalf("backdate rename: %v", err)
		}
	}
}

func countTaskFiles(t *testing.T, signalsDir string) int {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(signalsDir, "tasks"))
	if err != nil {
		t.Fatalf("read tasks dir: %v", err)
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			n++
		}
	}
	return n
}
