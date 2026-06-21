package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContentHash_Deterministic(t *testing.T) {
	input := []byte("hello world")
	h1 := contentHash(input, 0)
	h2 := contentHash(input, 0)
	if h1 != h2 {
		t.Errorf("contentHash not deterministic: %q != %q", h1, h2)
	}
	if len(h1) != 64 {
		t.Errorf("expected 64-char hex, got len=%d", len(h1))
	}
}

func TestContentHash_DifferentInputs(t *testing.T) {
	a := contentHash([]byte("aaa"), 0)
	b := contentHash([]byte("bbb"), 0)
	if a == b {
		t.Error("different inputs produced same hash")
	}
}

func TestContentHash_BudgetTruncates(t *testing.T) {
	full := []byte("abcdefgh")
	// budget=4 should hash only first 4 bytes
	h1 := contentHash(full, 4)
	h2 := contentHash([]byte("abcd"), 0)
	if h1 != h2 {
		t.Errorf("budget=4 should hash same as first 4 bytes; h1=%q h2=%q", h1, h2)
	}
}

func TestContentHash_BudgetZeroMeansAll(t *testing.T) {
	data := []byte("short")
	h1 := contentHash(data, 0)
	h2 := contentHash(data, 1000) // budget larger than data — no truncation
	if h1 != h2 {
		t.Errorf("budget=0 and budget>len should hash same bytes")
	}
}

func TestLoadGoldenIndex_Fallback(t *testing.T) {
	// Passing a non-existent index should return hardcoded fallback (non-empty).
	sources := loadGoldenIndex("/does/not/exist/index.md", "/base")
	if len(sources) == 0 {
		t.Error("expected hardcoded fallback sources, got empty slice")
	}
}

func TestLoadGoldenIndex_ParsesTable(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "golden-docs-index.md")
	content := `| name | path | tier | budget | description |
|------|------|------|--------|-------------|
| NORTHSTAR | docs/NORTHSTAR.md | 1 | 4000 | northstar doc |
| IGNORED | docs/ignored.md | 2 | 0 | tier 2 — should skip |
`
	if err := os.WriteFile(indexPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	sources := loadGoldenIndex(indexPath, "/base")
	if len(sources) != 1 {
		t.Fatalf("want 1 Tier-1 source, got %d", len(sources))
	}
	if sources[0].Name != "NORTHSTAR" {
		t.Errorf("source[0].Name = %q, want NORTHSTAR", sources[0].Name)
	}
	if sources[0].Budget != 4000 {
		t.Errorf("source[0].Budget = %d, want 4000", sources[0].Budget)
	}
}

func TestLoadGoldenIndex_EmptyTier1FallsBack(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "golden-docs-index.md")
	// Only Tier-2 rows — should fall back to hardcoded
	content := `| name | path | tier | budget | description |
|------|------|------|--------|-------------|
| IGNORED | docs/ignored.md | 2 | 0 | tier 2 only |
`
	if err := os.WriteFile(indexPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	sources := loadGoldenIndex(indexPath, "/base")
	if len(sources) == 0 {
		t.Error("expected hardcoded fallback when all rows are Tier 2, got empty")
	}
}
