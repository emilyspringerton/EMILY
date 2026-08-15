package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

// TestRunBacklogCuration_LongMultibyteSummaryStaysValidUTF8 is the regression
// test for a real bug found live 2026-08-15: the previous truncation
// (summary[:119], a byte-slice cut) could land mid-multibyte-character on
// non-ASCII text and emit an invalid UTF-8 byte into BACKLOG.md. That
// silently broke plain `grep` against the whole file (grep treats a file
// with invalid UTF-8 as binary and skips text search), which is how this
// was originally discovered -- searches that should have matched came back
// empty. This test uses a real Chinese summary long enough to require
// truncation at the exact boundary that broke in production.
func TestRunBacklogCuration_LongMultibyteSummaryStaysValidUTF8(t *testing.T) {
	gitRoot := t.TempDir()
	fatBabyRoot := t.TempDir()

	obsDir := filepath.Join(fatBabyRoot, "var", "emily-observations")
	if err := os.MkdirAll(obsDir, 0o755); err != nil {
		t.Fatalf("mkdir obs dir: %v", err)
	}

	// The real production summary that triggered this bug (179 runes long,
	// truncation must land inside a run of Chinese/English mixed text).
	longSummary := "自我修正:剛才Apple #13608跟REDGARDEN commit都誤稱這是'首次'真正end-to-end --autocurriculum訓練——實際上2026-08-11已經跑過一次(500K timesteps,75% win rate,REDGARDEN CHANGELOG已有記錄,sess-20260810-0505-a53abca2)。這次(2026-08-14)是第二次,只跑了200K timesteps"
	if utf8.RuneCountInString(longSummary) <= 120 {
		t.Fatalf("test fixture summary must be >120 runes to exercise truncation, got %d", utf8.RuneCountInString(longSummary))
	}

	obs := struct {
		Timestamp string `json:"timestamp"`
		Summary   string `json:"summary"`
	}{Timestamp: "2026-08-15T00-00-00Z", Summary: longSummary}
	data, err := json.Marshal(obs)
	if err != nil {
		t.Fatalf("marshal obs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(obsDir, "2026-08-15T00-00-00Z.json"), data, 0o644); err != nil {
		t.Fatalf("write obs: %v", err)
	}

	backlogPath := filepath.Join(gitRoot, "BACKLOG.md")
	initial := "# BACKLOG\n\n## BACKLOG PROTOCOL\n\nrules here\n"
	if err := os.WriteFile(backlogPath, []byte(initial), 0o644); err != nil {
		t.Fatalf("write initial BACKLOG.md: %v", err)
	}

	s := &IntegrationStore{gitRoot: gitRoot, fatBabyRoot: fatBabyRoot}
	n, err := s.runBacklogCuration()
	if err != nil {
		t.Fatalf("runBacklogCuration: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 item curated, got %d", n)
	}

	out, err := os.ReadFile(backlogPath)
	if err != nil {
		t.Fatalf("read BACKLOG.md: %v", err)
	}
	if !utf8.Valid(out) {
		t.Error("BACKLOG.md contains invalid UTF-8 after curation -- the exact regression this test guards against")
	}
	if _, err := os.ReadFile(backlogPath); err != nil {
		t.Fatalf("re-read BACKLOG.md: %v", err)
	}
	// Sanity: the truncated summary should still be present and end with the
	// ellipsis marker, not be silently dropped or empty.
	if !strings.Contains(string(out), "…") {
		t.Error("expected the truncated summary to end with an ellipsis marker")
	}
}
