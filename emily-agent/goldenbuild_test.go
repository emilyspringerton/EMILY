package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

// buildTestCompiler sets up a GoldenDocCompiler with one real source file.
// pureCompress is deterministic and needs no network/API key, so Build()
// always produces the real thing here — there is no more "degraded" path
// to simulate.
func buildTestCompiler(t *testing.T, sourceContent string) (*GoldenDocCompiler, string) {
	t.Helper()
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "source.md")
	if err := os.WriteFile(srcPath, []byte(sourceContent), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	cachePath := filepath.Join(dir, "golden-cache.json")
	g := &GoldenDocCompiler{
		outPath:   filepath.Join(dir, "full-system-context.md"),
		cachePath: cachePath,
		sources:   []GoldenSource{{Name: "TESTDOC", Path: srcPath, Budget: 0}},
	}
	return g, cachePath
}

func TestBuild_ProducesRealCompressedContentNoNetworkNeeded(t *testing.T) {
	g, _ := buildTestCompiler(t, "# Heading\nSome real body text.\n")
	if err := g.Build(context.Background()); err != nil {
		t.Fatalf("Build: %v", err)
	}
	cache := g.loadCache()
	entry, ok := cache.Entries[g.sources[0].Path]
	if !ok {
		t.Fatal("expected a cache entry after Build")
	}
	if !strings.Contains(entry.Compressed, "## TESTDOC") {
		t.Errorf("expected compressed section to carry the source label, got %q", entry.Compressed)
	}
	if !strings.Contains(entry.Compressed, "Heading") || !strings.Contains(entry.Compressed, "Some real body text.") {
		t.Errorf("expected the real header+lead-line extracted content, got %q", entry.Compressed)
	}
}

func TestBuild_CachedEntryIsReusedWhenContentUnchanged(t *testing.T) {
	g, cachePath := buildTestCompiler(t, "content")
	raw, _ := os.ReadFile(g.sources[0].Path)
	h := contentHash(raw, g.sources[0].Budget)
	seed := goldenCache{Entries: map[string]goldenCacheEntry{
		g.sources[0].Path: {Hash: h, Compressed: "## TESTDOC\nseeded content\n", UpdatedAt: "2020-01-01T00:00:00Z"},
	}}
	b, _ := json.Marshal(seed)
	if err := os.WriteFile(cachePath, b, 0o644); err != nil {
		t.Fatalf("seed cache: %v", err)
	}

	if err := g.Build(context.Background()); err != nil {
		t.Fatalf("Build: %v", err)
	}
	cache := g.loadCache()
	entry := cache.Entries[g.sources[0].Path]
	if entry.UpdatedAt != "2020-01-01T00:00:00Z" {
		t.Errorf("expected the cache entry to be reused untouched (hash matched), got UpdatedAt=%q", entry.UpdatedAt)
	}
	if entry.Compressed != "## TESTDOC\nseeded content\n" {
		t.Errorf("expected the seeded compressed content to be preserved, got %q", entry.Compressed)
	}
}

func TestMaybeRebuild_SkipsWhenNoChange(t *testing.T) {
	g, cachePath := buildTestCompiler(t, "unchanged content")
	raw, _ := os.ReadFile(g.sources[0].Path)
	h := contentHash(raw, g.sources[0].Budget)
	seed := goldenCache{Entries: map[string]goldenCacheEntry{
		g.sources[0].Path: {Hash: h, Compressed: "## TESTDOC\nreal output\n", UpdatedAt: "2020-01-01T00:00:00Z"},
	}}
	b, _ := json.Marshal(seed)
	if err := os.WriteFile(cachePath, b, 0o644); err != nil {
		t.Fatalf("seed cache: %v", err)
	}

	if err := g.MaybeRebuild(context.Background()); err != nil {
		t.Fatalf("MaybeRebuild: %v", err)
	}
	if _, err := os.Stat(g.outPath); !os.IsNotExist(err) {
		t.Error("expected MaybeRebuild to skip Build() entirely (no output file written) when nothing changed")
	}
}

func TestMaybeRebuild_TriggersWhenContentChanges(t *testing.T) {
	g, cachePath := buildTestCompiler(t, "new content, differs from the seeded hash below")
	seed := goldenCache{Entries: map[string]goldenCacheEntry{
		g.sources[0].Path: {Hash: "stale-hash-that-wont-match", Compressed: "## TESTDOC\nold\n", UpdatedAt: "2020-01-01T00:00:00Z"},
	}}
	b, _ := json.Marshal(seed)
	if err := os.WriteFile(cachePath, b, 0o644); err != nil {
		t.Fatalf("seed cache: %v", err)
	}

	if err := g.MaybeRebuild(context.Background()); err != nil {
		t.Fatalf("MaybeRebuild: %v", err)
	}
	out, err := os.ReadFile(g.outPath)
	if err != nil {
		t.Fatalf("MaybeRebuild should have written full-system-context.md (proving Build() ran), but: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected non-empty full-system-context.md after a triggered rebuild")
	}
}

func TestPureCompress_ExtractsHeadersAndLeadLines(t *testing.T) {
	text := "# Title\nFirst line under title.\nIgnored second line.\n\n## Sub\nLead line under sub.\n"
	out := pureCompress("DOC", text, 900)
	for _, want := range []string{"## DOC", "Title", "First line under title.", "Sub", "Lead line under sub."} {
		if !strings.Contains(out, want) {
			t.Errorf("pureCompress output missing %q; got %q", want, out)
		}
	}
	if strings.Contains(out, "Ignored second line.") {
		t.Errorf("pureCompress should only keep the FIRST body line after a header, got %q", out)
	}
}

func TestPureCompress_FallsBackToPlainTruncationWithNoHeaders(t *testing.T) {
	text := "just plain lines\nno markdown structure at all\nthird line\n"
	out := pureCompress("PLAIN", text, 900)
	if !strings.Contains(out, "## PLAIN") {
		t.Errorf("expected the section label, got %q", out)
	}
	if !strings.Contains(out, "just plain lines") || !strings.Contains(out, "third line") {
		t.Errorf("expected plain-text fallback to keep real content, got %q", out)
	}
}

func TestPureCompress_RespectsMaxChars(t *testing.T) {
	text := strings.Repeat("# H\nbody line\n", 200)
	out := pureCompress("BIG", text, 200)
	if len(out) > 210 { // small slack for the "…\n" suffix
		t.Errorf("expected output capped near maxChars=200, got %d chars", len(out))
	}
}

func TestPureCompress_Deterministic(t *testing.T) {
	text := "# H\nsome body\n"
	a := pureCompress("D", text, 900)
	b := pureCompress("D", text, 900)
	if a != b {
		t.Error("pureCompress should be deterministic for identical input")
	}
}
