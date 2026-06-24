// memory_consolidate.go — S125-08: weekly emily-memory consolidation.
//
// Reads all *.json files in memDir, merges flat key→value maps,
// writes consolidated.json. Called from the PLAN phase at most once per 7 days.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// consolidateMemoryFragments merges all *.json fragments in memDir into
// consolidated.json. Non-destructive: source files are not modified.
func consolidateMemoryFragments(memDir string) error {
	entries, err := os.ReadDir(memDir)
	if err != nil {
		return fmt.Errorf("readdir %s: %w", memDir, err)
	}

	merged := map[string]json.RawMessage{}
	var sources []string
	count := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		if e.Name() == "consolidated.json" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(memDir, e.Name()))
		if err != nil {
			continue
		}
		var fragment map[string]json.RawMessage
		if err := json.Unmarshal(raw, &fragment); err != nil {
			continue
		}
		for k, v := range fragment {
			merged[k] = v
		}
		sources = append(sources, e.Name())
		count++
	}
	if count == 0 {
		return nil
	}

	now, _ := json.Marshal(time.Now().UTC().Format(time.RFC3339))
	srcJSON, _ := json.Marshal(sources)
	merged["_consolidated_at"] = now
	merged["_sources"] = srcJSON

	out, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	outPath := filepath.Join(memDir, "consolidated.json")
	return os.WriteFile(outPath, out, 0o644)
}
