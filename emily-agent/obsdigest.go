// obsdigest.go — observation digest builder for Emily Prime's PLAN phase.
//
// Reads all JSON files in signals/observations/, counts by severity, extracts
// the top 3 entity names (from summary/source fields), and writes
// emily-memory/observations-digest.json.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ObservationRecord is the on-disk format of a single observation file.
type ObservationRecord struct {
	Summary         string `json:"summary"`
	Findings        string `json:"findings"`
	ObservationType string `json:"observation_type"`
	Severity        string `json:"severity"`
	Source          string `json:"source"`
	SuggestedFix    string `json:"suggested_fix,omitempty"`
	Timestamp       string `json:"timestamp"`
	Entity          string `json:"entity,omitempty"` // optional explicit entity name
}

// ObservationDigest is the structured summary written to emily-memory.
type ObservationDigest struct {
	GeneratedAt    time.Time      `json:"generated_at"`
	TotalCount     int            `json:"total_count"`
	BySeverity     map[string]int `json:"by_severity"`
	TopEntities    []string       `json:"top_entities"`
	LastSeenAt     string         `json:"last_seen_at"` // ISO-8601 of most recent obs
	RecentSummaries []string      `json:"recent_summaries"` // top 3 by recency
}

// buildObservationDigest scans obsDir and writes the digest to memDir/observations-digest.json.
// obsDir = signals/observations/, memDir = emily-memory/
func buildObservationDigest(obsDir, memDir string) (ObservationDigest, error) {
	entries, err := os.ReadDir(obsDir)
	if err != nil {
		return ObservationDigest{}, fmt.Errorf("read observations dir: %w", err)
	}

	var records []ObservationRecord
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(obsDir, e.Name()))
		if err != nil {
			continue
		}
		var rec ObservationRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			continue
		}
		records = append(records, rec)
	}

	digest := ObservationDigest{
		GeneratedAt: time.Now().UTC(),
		TotalCount:  len(records),
		BySeverity:  make(map[string]int),
	}
	if len(records) == 0 {
		return digest, writeDigest(memDir, digest)
	}

	// Count by severity.
	for _, r := range records {
		sev := strings.ToLower(r.Severity)
		if sev == "" {
			sev = "unknown"
		}
		digest.BySeverity[sev]++
	}

	// Sort records by timestamp descending (filename is already ISO-prefixed, so sort is lexicographic).
	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp > records[j].Timestamp
	})

	// Last-seen timestamp from most recent record.
	if records[0].Timestamp != "" {
		digest.LastSeenAt = records[0].Timestamp
	}

	// Top 3 recent summaries.
	for i, r := range records {
		if i >= 3 {
			break
		}
		if s := strings.TrimSpace(r.Summary); s != "" {
			digest.RecentSummaries = append(digest.RecentSummaries, s)
		}
	}

	// Top 3 entity names: prefer explicit entity field, else source.
	entityCount := make(map[string]int)
	for _, r := range records {
		name := strings.TrimSpace(r.Entity)
		if name == "" {
			name = strings.TrimSpace(r.Source)
		}
		if name != "" && name != "unknown" {
			entityCount[name]++
		}
	}
	type kv struct {
		key   string
		count int
	}
	var pairs []kv
	for k, v := range entityCount {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].key < pairs[j].key
	})
	for i, p := range pairs {
		if i >= 3 {
			break
		}
		digest.TopEntities = append(digest.TopEntities, p.key)
	}

	return digest, writeDigest(memDir, digest)
}

func writeDigest(memDir string, digest ObservationDigest) error {
	if err := os.MkdirAll(memDir, 0o755); err != nil {
		return fmt.Errorf("mkdir emily-memory: %w", err)
	}
	out, err := json.MarshalIndent(digest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(memDir, "observations-digest.json"), out, 0o644)
}
