package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeFakeObs(t *testing.T, dir string, name string, rec ObservationRecord) {
	t.Helper()
	raw, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), raw, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBuildObsDigestEmpty(t *testing.T) {
	obsDir := t.TempDir()
	memDir := t.TempDir()
	d, err := buildObservationDigest(obsDir, memDir)
	if err != nil {
		t.Fatal(err)
	}
	if d.TotalCount != 0 {
		t.Errorf("want 0, got %d", d.TotalCount)
	}
	// Digest file should be written.
	if _, err := os.Stat(filepath.Join(memDir, "observations-digest.json")); err != nil {
		t.Error("digest file not written")
	}
}

func TestBuildObsDigestCounts(t *testing.T) {
	obsDir := t.TempDir()
	memDir := t.TempDir()
	writeFakeObs(t, obsDir, "a.json", ObservationRecord{Severity: "warn", Source: "fatbaby", Timestamp: "2026-01-01T00:00:00Z", Summary: "obs A"})
	writeFakeObs(t, obsDir, "b.json", ObservationRecord{Severity: "error", Source: "emily", Timestamp: "2026-01-02T00:00:00Z", Summary: "obs B"})
	writeFakeObs(t, obsDir, "c.json", ObservationRecord{Severity: "warn", Source: "fatbaby", Timestamp: "2026-01-03T00:00:00Z", Summary: "obs C"})

	d, err := buildObservationDigest(obsDir, memDir)
	if err != nil {
		t.Fatal(err)
	}
	if d.TotalCount != 3 {
		t.Errorf("want 3, got %d", d.TotalCount)
	}
	if d.BySeverity["warn"] != 2 {
		t.Errorf("want 2 warns, got %d", d.BySeverity["warn"])
	}
	if d.BySeverity["error"] != 1 {
		t.Errorf("want 1 error, got %d", d.BySeverity["error"])
	}
}

func TestBuildObsDigestTopEntities(t *testing.T) {
	obsDir := t.TempDir()
	memDir := t.TempDir()
	writeFakeObs(t, obsDir, "a.json", ObservationRecord{Source: "fatbaby", Timestamp: "2026-01-01T00:00:00Z"})
	writeFakeObs(t, obsDir, "b.json", ObservationRecord{Source: "fatbaby", Timestamp: "2026-01-02T00:00:00Z"})
	writeFakeObs(t, obsDir, "c.json", ObservationRecord{Source: "emily", Timestamp: "2026-01-03T00:00:00Z"})
	writeFakeObs(t, obsDir, "d.json", ObservationRecord{Entity: "AAPL", Source: "sec", Timestamp: "2026-01-04T00:00:00Z"})

	d, err := buildObservationDigest(obsDir, memDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.TopEntities) == 0 {
		t.Error("expected at least one top entity")
	}
	// fatbaby appears 2× — should be first.
	if d.TopEntities[0] != "fatbaby" {
		t.Errorf("want fatbaby first, got %s", d.TopEntities[0])
	}
}

func TestBuildObsDigestRecentSummaries(t *testing.T) {
	obsDir := t.TempDir()
	memDir := t.TempDir()
	for i := 0; i < 5; i++ {
		ts := "2026-01-0"
		ts += string(rune('1' + i))
		ts += "T00:00:00Z"
		writeFakeObs(t, obsDir, "obs_"+string(rune('a'+i))+".json", ObservationRecord{
			Summary: "summary " + string(rune('A'+i)), Timestamp: ts,
		})
	}
	d, err := buildObservationDigest(obsDir, memDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.RecentSummaries) != 3 {
		t.Errorf("want 3 recent summaries, got %d", len(d.RecentSummaries))
	}
	// Most recent is obs_e (timestamp 2026-01-05).
	if d.RecentSummaries[0] != "summary E" {
		t.Errorf("want 'summary E' first, got %q", d.RecentSummaries[0])
	}
}

func TestBuildObsDigestLastSeen(t *testing.T) {
	obsDir := t.TempDir()
	memDir := t.TempDir()
	writeFakeObs(t, obsDir, "a.json", ObservationRecord{Timestamp: "2026-01-01T00:00:00Z"})
	writeFakeObs(t, obsDir, "b.json", ObservationRecord{Timestamp: "2026-06-15T00:00:00Z"})
	d, err := buildObservationDigest(obsDir, memDir)
	if err != nil {
		t.Fatal(err)
	}
	if d.LastSeenAt != "2026-06-15T00:00:00Z" {
		t.Errorf("want last seen 2026-06-15, got %s", d.LastSeenAt)
	}
}

func TestBuildObsDigestWriteFile(t *testing.T) {
	obsDir := t.TempDir()
	memDir := t.TempDir()
	writeFakeObs(t, obsDir, "x.json", ObservationRecord{Severity: "info", Source: "test", Timestamp: "2026-01-01T00:00:00Z"})
	_, err := buildObservationDigest(obsDir, memDir)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(memDir, "observations-digest.json"))
	if err != nil {
		t.Fatal("digest file missing")
	}
	var d ObservationDigest
	if err := json.Unmarshal(raw, &d); err != nil {
		t.Fatalf("invalid JSON in digest: %v", err)
	}
	if d.TotalCount != 1 {
		t.Errorf("want 1, got %d", d.TotalCount)
	}
}

func TestBuildObsDigestSkipsNonJSON(t *testing.T) {
	obsDir := t.TempDir()
	memDir := t.TempDir()
	// Write a non-JSON file that should be ignored.
	_ = os.WriteFile(filepath.Join(obsDir, "readme.txt"), []byte("ignore me"), 0o644)
	writeFakeObs(t, obsDir, "real.json", ObservationRecord{Severity: "info", Timestamp: "2026-01-01T00:00:00Z"})
	d, err := buildObservationDigest(obsDir, memDir)
	if err != nil {
		t.Fatal(err)
	}
	if d.TotalCount != 1 {
		t.Errorf("want 1, got %d", d.TotalCount)
	}
}
