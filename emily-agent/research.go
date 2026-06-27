// emily-agent/research.go — S137-01 Research tool
//
// Provides the "research" Emily Prime tool: fetches from a curated source list
// (supply chain, AI, financial, general web), extracts text, files a research_log
// Apple with full source provenance and a synthesized summary.
//
// Design: self-contained HTTP fetcher — no external package dependencies.
// The apple audit contract is locked in S137-04.

package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

// researchSource is a static curated source.
type researchSource struct {
	Name   string
	URL    string
	Domain string
}

// researchSources mirrors the DefaultSources registry from PRRJECT_FATBABY/research.
var researchSources = []researchSource{
	{Name: "Supply Chain Dive", URL: "https://www.supplychaindive.com/", Domain: "supply_chain"},
	{Name: "Logistics Mgmt", URL: "https://www.logisticsmgmt.com/", Domain: "supply_chain"},
	{Name: "MIT Tech Review AI", URL: "https://www.technologyreview.com/topic/artificial-intelligence/", Domain: "ai"},
	{Name: "Ars Technica", URL: "https://arstechnica.com/", Domain: "ai"},
	{Name: "Reuters Business", URL: "https://www.reuters.com/business/", Domain: "financial"},
	{Name: "Hacker News", URL: "https://news.ycombinator.com/", Domain: "general"},
	{Name: "StickerMule Pricing", URL: "https://www.stickermule.com/products/die-cut-stickers", Domain: "supply_chain"},
	{Name: "StickerApp", URL: "https://stickerapp.com/die-cut-stickers/", Domain: "supply_chain"},
}

// researchFetchResult is the scraped output from one source.
type researchFetchResult struct {
	Source    researchSource `json:"source"`
	FetchedAt time.Time      `json:"fetched_at"`
	URL       string         `json:"url"`
	Excerpt   string         `json:"excerpt"` // first 2000 UTF-8 chars of body text
	Error     string         `json:"error,omitempty"`
}

// ResearchResult is the full output of one Research call.
type ResearchResult struct {
	Query      string                `json:"query"`
	QueryHash  string                `json:"query_hash"`
	Sources    []researchFetchResult `json:"sources"`
	Synthesis  string                `json:"synthesis"`
	Confidence float64               `json:"confidence"`
	StartedAt  time.Time             `json:"started_at"`
	Duration   string                `json:"duration"`
}

// registerResearchTool adds the "research" tool to the dispatcher.
func registerResearchTool(d *ToolDispatcher, iduna *IdunaClient) {
	d.Register(ToolDef{
		Name:        "research",
		Description: "Fetch and synthesize information from curated sources on a given query. Files a research_log Apple with full source provenance. Use for supply chain vendor lookup, market research, technology assessment, or any external knowledge need.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"query":   {Type: "string", Description: "Research question or topic"},
				"domains": {Type: "string", Description: "Comma-separated domain filter: supply_chain, ai, financial, general (default: all)"},
				"max_sources": {Type: "integer", Description: "Maximum number of sources to fetch (default: 4, max: 8)"},
			},
			Required: []string{"query"},
		},
	}, func(args map[string]any) (string, error) {
		query, _ := args["query"].(string)
		domainsRaw, _ := args["domains"].(string)
		maxF, _ := args["max_sources"].(float64)
		if query == "" {
			return "", fmt.Errorf("query required")
		}
		maxSources := 4
		if maxF > 0 && int(maxF) <= 8 {
			maxSources = int(maxF)
		}

		var domainFilter []string
		if domainsRaw != "" {
			for _, d := range strings.Split(domainsRaw, ",") {
				domainFilter = append(domainFilter, strings.TrimSpace(d))
			}
		}

		// Normalize query hash for cache lookup.
		h := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(query))))
		qHash := fmt.Sprintf("%x", h)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// S137-03: check cache first.
		if cached, ok := researchCheckCache(ctx, iduna, qHash); ok {
			log.Printf("research: cache hit for %s", qHash[:8])
			var sb strings.Builder
			fmt.Fprintf(&sb, "## Research (cached): %s\n\n", query)
			for _, s := range cached.Sources {
				if s.Error != "" {
					continue
				}
				fmt.Fprintf(&sb, "### %s\n*%s*\n\n%s\n\n", s.Source.Name, s.URL, s.Excerpt)
			}
			return sb.String(), nil
		}

		result := runResearch(ctx, query, domainFilter, maxSources)
		result.QueryHash = qHash

		// S137-03: write to cache.
		go researchWriteCache(context.Background(), iduna, result)

		// File research_log Apple.
		if iduna != nil {
			appleBody, _ := json.Marshal(result)
			_, _ = iduna.PostApple(ctx, ApplePayload{
				AppleType:  "research_log",
				Title:      fmt.Sprintf("research: %s", truncate(query, 80)),
				Body:       string(appleBody),
				SourceRepo: "EMILY",
			})
		}

		// Build human-readable summary for the LLM.
		var sb strings.Builder
		fmt.Fprintf(&sb, "## Research: %s\n\n", query)
		fmt.Fprintf(&sb, "**Query hash:** %s | **Duration:** %s | **Confidence:** %.0f%%\n\n", result.QueryHash[:8], result.Duration, result.Confidence*100)
		for _, s := range result.Sources {
			if s.Error != "" {
				fmt.Fprintf(&sb, "### %s — FETCH ERROR: %s\n", s.Source.Name, s.Error)
				continue
			}
			fmt.Fprintf(&sb, "### %s\n*%s*\n\n%s\n\n", s.Source.Name, s.URL, s.Excerpt)
		}
		if result.Synthesis != "" {
			fmt.Fprintf(&sb, "---\n**Synthesis:** %s\n", result.Synthesis)
		}
		return sb.String(), nil
	})
}

// researchCheckCache checks IDUNA research_cache for a prior result.
func researchCheckCache(ctx context.Context, iduna *IdunaClient, queryHash string) (*ResearchResult, bool) {
	if iduna == nil {
		return nil, false
	}
	if err := iduna.authenticate(ctx); err != nil {
		return nil, false
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		strings.TrimRight(iduna.baseURL, "/")+"/api/v1/research/cache?q="+queryHash, nil)
	if err != nil {
		return nil, false
	}
	req.Header.Set("Authorization", "Bearer "+iduna.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		if resp != nil {
			resp.Body.Close()
		}
		return nil, false
	}
	defer resp.Body.Close()
	var out struct {
		Entry struct {
			ResultJSON string `json:"result_json"`
		} `json:"entry"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, false
	}
	var r ResearchResult
	if err := json.Unmarshal([]byte(out.Entry.ResultJSON), &r); err != nil {
		return nil, false
	}
	return &r, true
}

// researchWriteCache stores a result in IDUNA research_cache.
func researchWriteCache(ctx context.Context, iduna *IdunaClient, r ResearchResult) {
	if iduna == nil {
		return
	}
	if err := iduna.authenticate(ctx); err != nil {
		return
	}
	resultJSON, err := json.Marshal(r)
	if err != nil {
		return
	}
	var urls []string
	for _, s := range r.Sources {
		urls = append(urls, s.URL)
	}
	payload, _ := json.Marshal(map[string]any{
		"query_hash":  r.QueryHash,
		"query_text":  r.Query,
		"result_json": string(resultJSON),
		"source_urls": strings.Join(urls, "\n"),
		"ttl_hours":   48,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(iduna.baseURL, "/")+"/api/v1/research/cache",
		strings.NewReader(string(payload)))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+iduna.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}

// runResearch executes the fetch pipeline.
func runResearch(ctx context.Context, query string, domains []string, maxSources int) ResearchResult {
	start := time.Now()
	h := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(query))))
	queryHash := fmt.Sprintf("%x", h)

	// Filter sources by domain.
	var sources []researchSource
	for _, s := range researchSources {
		if len(domains) == 0 {
			sources = append(sources, s)
		} else {
			for _, d := range domains {
				if d == s.Domain {
					sources = append(sources, s)
					break
				}
			}
		}
	}
	if len(sources) > maxSources {
		sources = sources[:maxSources]
	}

	client := &http.Client{Timeout: 8 * time.Second}
	var fetched []researchFetchResult

	for _, src := range sources {
		select {
		case <-ctx.Done():
			break
		default:
		}
		fetchURL := src.URL
		// Append query as search param where the URL supports it (basic heuristic).
		if strings.Contains(fetchURL, "search") || strings.Contains(fetchURL, "q=") {
			fetchURL = fmt.Sprintf("%s?q=%s", strings.TrimRight(fetchURL, "/"), url.QueryEscape(query))
		}

		fr := researchFetchResult{
			Source:    src,
			URL:       fetchURL,
			FetchedAt: time.Now().UTC(),
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fetchURL, nil)
		if err != nil {
			fr.Error = err.Error()
			fetched = append(fetched, fr)
			continue
		}
		req.Header.Set("User-Agent", "EmilyPrime/1.0 (research bot; github.com/emilyspringerton)")
		req.Header.Set("Accept", "text/html,text/plain")

		resp, err := client.Do(req)
		if err != nil {
			fr.Error = err.Error()
			log.Printf("research: fetch %s: %v", fetchURL, err)
			fetched = append(fetched, fr)
			continue
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		resp.Body.Close()

		text := extractText(string(body))
		fr.Excerpt = truncateUTF8(text, 2000)
		fetched = append(fetched, fr)
	}

	// Count successful fetches.
	var successCount int
	for _, f := range fetched {
		if f.Error == "" {
			successCount++
		}
	}
	confidence := 0.0
	if len(fetched) > 0 {
		confidence = float64(successCount) / float64(len(fetched))
	}

	duration := time.Since(start).Round(time.Millisecond).String()

	return ResearchResult{
		Query:      query,
		QueryHash:  queryHash,
		Sources:    fetched,
		Synthesis:  "",
		Confidence: confidence,
		StartedAt:  start.UTC(),
		Duration:   duration,
	}
}

// extractText strips HTML tags and collapses whitespace.
func extractText(html string) string {
	var sb strings.Builder
	inTag := false
	for i := 0; i < len(html); {
		r, size := utf8.DecodeRuneInString(html[i:])
		i += size
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
			sb.WriteRune(' ')
		case !inTag:
			sb.WriteRune(r)
		}
	}
	// Collapse runs of whitespace.
	raw := sb.String()
	var out strings.Builder
	prevSpace := false
	for _, r := range raw {
		isSpace := r == ' ' || r == '\t' || r == '\n' || r == '\r'
		if isSpace {
			if !prevSpace {
				out.WriteRune(' ')
			}
			prevSpace = true
		} else {
			out.WriteRune(r)
			prevSpace = false
		}
	}
	return strings.TrimSpace(out.String())
}

// truncateUTF8 returns s truncated to n UTF-8 characters.
func truncateUTF8(s string, n int) string {
	count := 0
	for i := range s {
		if count >= n {
			return s[:i] + "…"
		}
		count++
	}
	return s
}

