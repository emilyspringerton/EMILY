// emily-agent/goldenbuild.go
// GoldenDocCompiler reads Tier 1 golden docs from all repos and compresses
// them into EMILY/context/full-system-context.md via claude-haiku bilingual
// Chinese/English compression (~150 tokens per source).
//
// Called each cron cycle via MaybeRebuild (only rebuilds if a source changed).
// The compiled context is prepended to emilyStaticPrompt and fed into FABLE.

package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// goldenCacheEntry caches one compressed section keyed by source content hash.
type goldenCacheEntry struct {
	Hash       string `json:"hash"`       // sha256 of source content at last compress
	Compressed string `json:"compressed"` // haiku-compressed section
	UpdatedAt  string `json:"updated_at"` // RFC3339
}

// goldenCache persists per-source compressed sections so unchanged sources
// skip the haiku call on the next rebuild. Saves ~13 haiku calls per rebuild
// when only 1-2 sources change.
type goldenCache struct {
	Entries map[string]goldenCacheEntry `json:"entries"` // source path → entry
}

// GoldenSource is one documentation file to compress into the system context.
type GoldenSource struct {
	Name   string // label in compiled output, e.g. "FATBABY"
	Path   string // absolute path
	Budget int    // max chars to read; 0 = read all
}

// GoldenDocCompiler builds EMILY/context/full-system-context.md from Tier 1 golden docs.
type GoldenDocCompiler struct {
	apiKey    string
	outPath   string
	cachePath string // context/golden-cache.json
	sources   []GoldenSource
}

// NewGoldenDocCompiler creates a compiler whose source list is loaded from
// EMILY/context/golden-docs-index.md. If that file is absent, falls back to
// the hardcoded Tier 1 list so the compiler always has something to work with.
// emilyRoot should be the absolute path to the EMILY repo root.
func NewGoldenDocCompiler(emilyRoot, apiKey string) *GoldenDocCompiler {
	base := filepath.Dir(emilyRoot)
	contextDir := filepath.Join(emilyRoot, "context")
	sources := loadGoldenIndex(filepath.Join(contextDir, "golden-docs-index.md"), base)
	return &GoldenDocCompiler{
		apiKey:    apiKey,
		outPath:   filepath.Join(contextDir, "full-system-context.md"),
		cachePath: filepath.Join(contextDir, "golden-cache.json"),
		sources:   sources,
	}
}

// loadGoldenIndex parses EMILY/context/golden-docs-index.md and returns GoldenSources
// for all Tier 1 rows. Falls back to a minimal hardcoded set if the file can't be read.
func loadGoldenIndex(indexPath, base string) []GoldenSource {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		log.Printf("goldenbuild: index not found (%v) — using hardcoded fallback", err)
		return hardcodedSources(base)
	}

	var sources []GoldenSource
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "|---") || strings.HasPrefix(line, "| name") {
			continue
		}
		// Parse markdown table row: | name | path | tier | budget | description |
		fields := strings.Split(line, "|")
		if len(fields) < 6 {
			continue
		}
		name := strings.TrimSpace(fields[1])
		relPath := strings.TrimSpace(fields[2])
		tier := strings.TrimSpace(fields[3])
		budgetStr := strings.TrimSpace(fields[4])
		if name == "" || relPath == "" || tier != "1" {
			continue
		}
		budget := 0
		fmt.Sscanf(budgetStr, "%d", &budget)
		sources = append(sources, GoldenSource{
			Name:   name,
			Path:   filepath.Join(base, relPath),
			Budget: budget,
		})
	}

	if len(sources) == 0 {
		log.Printf("goldenbuild: index parsed 0 Tier 1 sources — using hardcoded fallback")
		return hardcodedSources(base)
	}
	log.Printf("goldenbuild: loaded %d Tier 1 sources from index", len(sources))
	return sources
}

// hardcodedSources is the bootstrap fallback used when golden-docs-index.md is absent.
func hardcodedSources(base string) []GoldenSource {
	return []GoldenSource{
		{Name: "FATBABY", Path: filepath.Join(base, "PRRJECT_FATBABY/docs/northstar/northstar.md"), Budget: 8000},
		{Name: "FATBABY-EXEC", Path: filepath.Join(base, "PRRJECT_FATBABY/docs/northstar/executive_summary.md"), Budget: 4000},
		{Name: "EMILY-SPEC", Path: filepath.Join(base, "EMILY/emily-prime-spec.md"), Budget: 4000},
		{Name: "EMILY-NORTH", Path: filepath.Join(base, "EMILY/docs/NORTHSTAR.md"), Budget: 3000},
		{Name: "EMILY-BACKLOG", Path: filepath.Join(base, "EMILY/GOLDEN.md"), Budget: 0},
		{Name: "EMIREE", Path: filepath.Join(base, "EMILY/emiree-emily-fatbaby.md"), Budget: 4000},
		{Name: "RSI-ENGINE", Path: filepath.Join(base, "EMILY/docs/emily-rsi-engine-spec.md"), Budget: 0},
		{Name: "EMIREE-SPEC", Path: filepath.Join(base, "EMILY/docs/emiree-over-agent-spec.md"), Budget: 0},
		{Name: "IDUNA", Path: filepath.Join(base, "IDUNA/golden.md"), Budget: 3000},
		{Name: "APPLES", Path: filepath.Join(base, "APPLES/docs/NORTHSTAR.md"), Budget: 0},
		{Name: "MJOLNIR", Path: filepath.Join(base, "MJOLNIR/docs/NORTHSTAR.md"), Budget: 0},
		{Name: "SHANKPIT", Path: filepath.Join(base, "SHANKPIT/docs2/NORTHSTAR.md"), Budget: 4000},
		{Name: "EMILYOS", Path: filepath.Join(base, "EmilyOS/docs/NORTHSTAR.md"), Budget: 0},
		{Name: "PITVIPER", Path: filepath.Join(base, "PITVIPER/docs/NORTHSTAR.md"), Budget: 0},
		{Name: "EMILY-CLI", Path: filepath.Join(base, "emily.cli/docs/NORTHSTAR.md"), Budget: 0},
		{Name: "GTM", Path: filepath.Join(base, "PRRJECT_FATBABY/docs/GTM_FUNNEL.md"), Budget: 3000},
		{Name: "TRAINING", Path: filepath.Join(base, "EMILY/docs/TRAINING_PIPELINE.md"), Budget: 2000},
		{Name: "EDIS", Path: filepath.Join(base, "EDIS/NORTHSTAR.md"), Budget: 3000},
	}
}

// MaybeRebuild rebuilds the full context if any source content changed since the last
// compile. Uses per-source SHA-256 hashes stored in golden-cache.json so that only
// changed sources are recompressed — unchanged sources reuse their cached section,
// avoiding ~13 haiku calls per rebuild when only 1-2 sources change.
func (g *GoldenDocCompiler) MaybeRebuild(ctx context.Context) error {
	cache := g.loadCache()
	anyChanged := false
	for _, src := range g.sources {
		raw, err := os.ReadFile(src.Path)
		if err != nil {
			continue
		}
		h := contentHash(raw, src.Budget)
		if entry, ok := cache.Entries[src.Path]; !ok || entry.Hash != h {
			anyChanged = true
			break
		}
	}
	if !anyChanged {
		return nil
	}
	return g.Build(ctx)
}

// Build reads all sources, compresses changed ones via haiku (reusing cached sections
// for unchanged sources), and writes full-system-context.md + golden-cache.json.
func (g *GoldenDocCompiler) Build(ctx context.Context) error {
	cache := g.loadCache()
	var sections []string
	built, reused := 0, 0
	for _, src := range g.sources {
		raw, err := os.ReadFile(src.Path)
		if err != nil {
			log.Printf("goldenbuild: skip %s: %v", src.Name, err)
			continue
		}
		h := contentHash(raw, src.Budget)
		if entry, ok := cache.Entries[src.Path]; ok && entry.Hash == h && entry.Compressed != "" {
			sections = append(sections, entry.Compressed)
			reused++
			built++
			continue
		}

		text := string(raw)
		if src.Budget > 0 && len(text) > src.Budget {
			text = text[:src.Budget]
		}
		compressed, err := g.compress(ctx, src.Name, text)
		if err != nil {
			log.Printf("goldenbuild: compress %s failed (%v), using truncated fallback", src.Name, err)
			if len(text) > 500 {
				text = text[:500]
			}
			compressed = fmt.Sprintf("## %s\n%s\n", src.Name, text)
		}
		cache.Entries[src.Path] = goldenCacheEntry{
			Hash:       h,
			Compressed: compressed,
			UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
		}
		sections = append(sections, compressed)
		built++
	}
	log.Printf("goldenbuild: %d sources built (%d reused from cache, %d new haiku calls)", built, reused, built-reused)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("# 全系統上下文 FULL SYSTEM CONTEXT — %s\n", time.Now().UTC().Format("2006-01-02")))
	b.WriteString("*Auto-generated by GoldenDocCompiler. Dense bilingual context for Emily Prime.*\n")
	b.WriteString(fmt.Sprintf("*Sources compiled: %d/%d.*\n\n", built, len(g.sources)))
	for _, s := range sections {
		b.WriteString(s)
		b.WriteString("\n")
	}

	if err := os.MkdirAll(filepath.Dir(g.outPath), 0o755); err != nil {
		return fmt.Errorf("mkdir context/: %w", err)
	}
	if err := os.WriteFile(g.outPath, []byte(b.String()), 0o644); err != nil {
		return fmt.Errorf("write full-system-context.md: %w", err)
	}
	log.Printf("goldenbuild: wrote %s (%d bytes, %d/%d sources)", g.outPath, b.Len(), built, len(g.sources))
	g.saveCache(cache)
	return nil
}

// loadCache reads golden-cache.json; returns an empty cache on any error.
func (g *GoldenDocCompiler) loadCache() goldenCache {
	data, err := os.ReadFile(g.cachePath)
	if err != nil {
		return goldenCache{Entries: make(map[string]goldenCacheEntry)}
	}
	var c goldenCache
	if err := json.Unmarshal(data, &c); err != nil || c.Entries == nil {
		return goldenCache{Entries: make(map[string]goldenCacheEntry)}
	}
	return c
}

// saveCache writes golden-cache.json; errors are logged but non-fatal.
func (g *GoldenDocCompiler) saveCache(c goldenCache) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Printf("goldenbuild: cache marshal: %v", err)
		return
	}
	if err := os.WriteFile(g.cachePath, data, 0o644); err != nil {
		log.Printf("goldenbuild: cache write: %v", err)
	}
}

// contentHash returns a hex SHA-256 of raw, truncated to budget chars if budget>0.
func contentHash(raw []byte, budget int) string {
	content := raw
	if budget > 0 && len(content) > budget {
		content = content[:budget]
	}
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

// compress calls claude-haiku to compress doc into a dense bilingual ≤180 token section.
func (g *GoldenDocCompiler) compress(ctx context.Context, name, text string) (string, error) {
	if g.apiKey == "" {
		return fmt.Sprintf("## %s\n(no ANTHROPIC_API_KEY — raw context unavailable)\n", name), nil
	}

	system := `You are Emily Prime's context compressor for EINHORN_INDUSTRIAL.
Compress documentation into a dense ≤180 token bilingual section.
Rules:
- Chinese for concepts/state (系統/北極星/完成/待辦/架構/待實現), English for technical names/APIs/file paths
- Every token carries signal — no filler, no repetition, no preamble
- Preserve: purpose, current status, open blockers, key interfaces/endpoints/paths
- Start with ## NAME (use the exact label provided)
- Output ONLY the compressed section`

	user := fmt.Sprintf("Label: %s\n\nDoc:\n%s", name, text)

	body := map[string]any{
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 400,
		"system":     system,
		"messages":   []map[string]any{{"role": "user", "content": user}},
	}
	payload, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.anthropic.com/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", g.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("api %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	for _, block := range result.Content {
		if block.Type == "text" {
			return strings.TrimSpace(block.Text) + "\n", nil
		}
	}
	return "", fmt.Errorf("no text block in haiku response")
}
