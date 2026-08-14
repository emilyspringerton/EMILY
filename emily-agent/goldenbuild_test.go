package main

import (
	"context"
	"encoding/json"
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

// buildTestCompiler sets up a GoldenDocCompiler with one real source file and
// no API key, so compress() takes its "no ANTHROPIC_API_KEY" stub path
// (err == nil, degraded content) -- exercises the same IsFallback marking as
// a real haiku failure without needing to mock the Anthropic API.
func buildTestCompiler(t *testing.T, sourceContent string) (*GoldenDocCompiler, string) {
	t.Helper()
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "source.md")
	if err := os.WriteFile(srcPath, []byte(sourceContent), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	cachePath := filepath.Join(dir, "golden-cache.json")
	g := &GoldenDocCompiler{
		apiKey:    "", // forces compress()'s degraded no-key stub path
		outPath:   filepath.Join(dir, "full-system-context.md"),
		cachePath: cachePath,
		sources:   []GoldenSource{{Name: "TESTDOC", Path: srcPath, Budget: 0}},
	}
	return g, cachePath
}

func TestBuild_MarksDegradedNoAPIKeyResultAsFallback(t *testing.T) {
	g, cachePath := buildTestCompiler(t, "some real content for the test doc")
	if err := g.Build(context.Background()); err != nil {
		t.Fatalf("Build: %v", err)
	}
	cache := g.loadCache()
	entry, ok := cache.Entries[g.sources[0].Path]
	if !ok {
		t.Fatal("expected a cache entry after Build")
	}
	if !entry.IsFallback {
		t.Error("expected IsFallback=true for the no-API-key degraded result, got false")
	}
	_ = cachePath
}

func TestBuild_RetriesFallbackEntryOnNextBuildEvenWithMatchingHash(t *testing.T) {
	g, _ := buildTestCompiler(t, "unchanged content across both builds")

	if err := g.Build(context.Background()); err != nil {
		t.Fatalf("first Build: %v", err)
	}
	cache1 := g.loadCache()
	entry1 := cache1.Entries[g.sources[0].Path]
	if !entry1.IsFallback {
		t.Fatal("setup invariant broken: first build should have produced a fallback entry")
	}
	firstUpdatedAt := entry1.UpdatedAt

	// Second Build with byte-identical source content -- the pre-fix bug would
	// reuse the cached fallback forever (hash matches, Compressed is non-empty).
	// The fix must retry regardless, since the cached entry is IsFallback.
	if err := g.Build(context.Background()); err != nil {
		t.Fatalf("second Build: %v", err)
	}
	cache2 := g.loadCache()
	entry2 := cache2.Entries[g.sources[0].Path]
	if !entry2.IsFallback {
		t.Error("expected the retried entry to still be marked IsFallback (no API key still set)")
	}
	if entry2.UpdatedAt == "" || entry2.UpdatedAt < firstUpdatedAt {
		t.Errorf("expected UpdatedAt to advance on retry (proves compress() was actually called again), got %q (was %q)", entry2.UpdatedAt, firstUpdatedAt)
	}
}

func TestBuild_RealCompressionIsNotMarkedFallbackAndIsReused(t *testing.T) {
	g, cachePath := buildTestCompiler(t, "content")
	// Simulate a prior real (non-fallback) compression already in the cache,
	// as if haiku had succeeded on an earlier run.
	raw, _ := os.ReadFile(g.sources[0].Path)
	h := contentHash(raw, g.sources[0].Budget)
	seed := goldenCache{Entries: map[string]goldenCacheEntry{
		g.sources[0].Path: {Hash: h, Compressed: "## TESTDOC\nreal haiku output\n", UpdatedAt: "2020-01-01T00:00:00Z", IsFallback: false},
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
		t.Errorf("expected a real (non-fallback) cache entry to be reused untouched, got UpdatedAt=%q", entry.UpdatedAt)
	}
	if entry.Compressed != "## TESTDOC\nreal haiku output\n" {
		t.Errorf("expected the real compressed content to be preserved, got %q", entry.Compressed)
	}
}

func TestMaybeRebuild_TriggersOnFallbackEntryEvenWithMatchingHash(t *testing.T) {
	g, cachePath := buildTestCompiler(t, "unchanged content")
	raw, _ := os.ReadFile(g.sources[0].Path)
	h := contentHash(raw, g.sources[0].Budget)
	// Seed a cache entry whose hash matches the current source exactly, but
	// which is flagged as a fallback -- the pre-fix bug's own hash-only check
	// would see "hash matches" and skip Build() entirely, leaving this source
	// stuck on its degraded content forever regardless of Build()'s own fix.
	seed := goldenCache{Entries: map[string]goldenCacheEntry{
		g.sources[0].Path: {Hash: h, Compressed: "## TESTDOC\nold degraded stub\n", UpdatedAt: "2020-01-01T00:00:00Z", IsFallback: true},
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
	cache := g.loadCache()
	if cache.Entries[g.sources[0].Path].UpdatedAt == "2020-01-01T00:00:00Z" {
		t.Error("expected the fallback entry's UpdatedAt to advance (proves Build() actually retried it, not just a no-op skip)")
	}
}

func TestMaybeRebuild_SkipsWhenNoChangeAndNoFallback(t *testing.T) {
	g, cachePath := buildTestCompiler(t, "unchanged content")
	raw, _ := os.ReadFile(g.sources[0].Path)
	h := contentHash(raw, g.sources[0].Budget)
	seed := goldenCache{Entries: map[string]goldenCacheEntry{
		g.sources[0].Path: {Hash: h, Compressed: "## TESTDOC\nreal output\n", UpdatedAt: "2020-01-01T00:00:00Z", IsFallback: false},
	}}
	b, _ := json.Marshal(seed)
	if err := os.WriteFile(cachePath, b, 0o644); err != nil {
		t.Fatalf("seed cache: %v", err)
	}

	if err := g.MaybeRebuild(context.Background()); err != nil {
		t.Fatalf("MaybeRebuild: %v", err)
	}
	if _, err := os.Stat(g.outPath); !os.IsNotExist(err) {
		t.Error("expected MaybeRebuild to skip Build() entirely (no output file written) when nothing changed and nothing is fallback-flagged")
	}
}
