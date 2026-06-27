// emily-agent/kgraph.go — S138-05 KnowledgeQuery tool
//
// Queries the EINHORN INDEX knowledge graph via IDUNA /api/v1/kgraph/query.
// Emily uses this before Research() for any factual lookup.
// On miss or error: returns empty result (caller falls back to Research).
//
// Apple audit: all KG queries file an Apple metadata record (type=audit, not completion).

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// KGraphResult mirrors kgraph.GraphResult (from PRRJECT_FATBABY).
type KGraphResult struct {
	SubjectName string    `json:"subject_name"`
	Predicate   string    `json:"predicate"`
	ObjectName  string    `json:"object_name"`
	Confidence  float64   `json:"confidence"`
	SourceURL   string    `json:"source_url"`
	ExtractedAt time.Time `json:"extracted_at"`
}

// registerKnowledgeQueryTool adds the "knowledge_query" tool to the dispatcher.
func registerKnowledgeQueryTool(d *ToolDispatcher, iduna *IdunaClient) {
	d.Register(ToolDef{
		Name:        "knowledge_query",
		Description: "Query the EINHORN INDEX knowledge graph for entities and relationships. Use this BEFORE Research() for any factual lookup. Returns empty on miss — caller should fall back to Research().",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"entity":    {Type: "string", Description: "Entity name to look up (company, product, person, etc.)"},
				"predicate": {Type: "string", Description: "Optional relationship filter: sells, supplies, priced_at, employs, etc. Leave empty for all relationships."},
				"limit":     {Type: "integer", Description: "Max results (default: 10, max: 50)"},
			},
			Required: []string{"entity"},
		},
	}, func(args map[string]any) (string, error) {
		entity, _ := args["entity"].(string)
		predicate, _ := args["predicate"].(string)
		limitF, _ := args["limit"].(float64)
		if entity == "" {
			return "", fmt.Errorf("entity required")
		}
		limit := 10
		if limitF > 0 && int(limitF) <= 50 {
			limit = int(limitF)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		results, err := kgraphQuery(ctx, iduna, entity, predicate, limit)
		if err != nil || len(results) == 0 {
			if iduna != nil && err != nil {
				_, _ = iduna.PostApple(ctx, ApplePayload{
					AppleType:  "audit",
					Title:      fmt.Sprintf("kgraph: query miss — %s", truncate(entity, 60)),
					Body:       fmt.Sprintf("entity=%q predicate=%q error=%v", entity, predicate, err),
					SourceRepo: "EMILY",
				})
			}
			return fmt.Sprintf("Knowledge graph: no results for entity=%q predicate=%q. Use Research() as fallback.", entity, predicate), nil
		}

		// File audit Apple for hit.
		if iduna != nil {
			_, _ = iduna.PostApple(ctx, ApplePayload{
				AppleType:  "audit",
				Title:      fmt.Sprintf("kgraph: query hit — %s (%d results)", truncate(entity, 50), len(results)),
				Body:       fmt.Sprintf("entity=%q predicate=%q results=%d", entity, predicate, len(results)),
				SourceRepo: "EMILY",
			})
		}

		// Build human-readable response.
		var sb strings.Builder
		fmt.Fprintf(&sb, "## Knowledge Graph: %s\n\n", entity)
		for _, r := range results {
			fmt.Fprintf(&sb, "- **%s** → %s → **%s** (conf=%.0f%%, src=%s)\n",
				r.SubjectName, r.Predicate, r.ObjectName, r.Confidence*100, r.SourceURL)
		}
		return sb.String(), nil
	})
}

// kgraphQuery calls IDUNA GET /api/v1/kgraph/query.
func kgraphQuery(ctx context.Context, iduna *IdunaClient, entity, predicate string, limit int) ([]KGraphResult, error) {
	if iduna == nil {
		return nil, nil
	}
	if err := iduna.authenticate(ctx); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("entity", entity)
	if predicate != "" {
		q.Set("predicate", predicate)
	}
	q.Set("limit", fmt.Sprintf("%d", limit))

	reqURL := strings.TrimRight(iduna.baseURL, "/") + "/api/v1/kgraph/query?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+iduna.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusServiceUnavailable {
		return nil, nil // kgraph service not up — graceful miss
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IDUNA kgraph: %d", resp.StatusCode)
	}
	var out struct {
		Results []KGraphResult `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Results, nil
}
