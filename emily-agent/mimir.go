// emily-agent/mimir.go
// MIMIR — Emily Prime's advisory planning tier.
//
// Renamed from FABLE 2026-07-24 (EMILY/BACKLOG.md SECTION 145's naming-collision
// note, found by the 2026-07-18 SAGA audit): this file's own "FABLE" collided
// with HQ-SPEC-AI-103's sovereign in-house model line, also named FABLE — two
// unrelated things sharing a name, a confusion risk flagged explicitly before
// any SECTION 145 work started. This advisory tier is the smaller, more easily
// renamed of the two; HQ-SPEC-AI-103 keeps FABLE. Named MIMIR (Norse: the wise
// being the gods consult for counsel) to match this codebase's existing Norse
// naming convention (NORN, SAGA, FATES) for exactly the role this file plays:
// advising on what to do next.
//
// MIMIR calls claude-haiku with the compressed golden backlog (GOLDEN.md) and
// recent Apple context to produce a prioritized sprint recommendation.
// It answers "what should Emily work on next?" in a structured, actionable form.
//
// Routes wired in main.go:
//   GET  /api/v1/emily/mimir/advice    — generate a fresh MIMIR sprint recommendation
//   POST /api/v1/emily/mimir/execute   — push top recommendation into RSI roadmap via HEIMDAL

package main

import (
	"bytes"
	"context"
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

// MimirItem is a single prioritized recommendation from the MIMIR advisor.
type MimirItem struct {
	Priority  int    `json:"priority"`  // 1 = highest
	Title     string `json:"title"`     // short name (≤60 chars)
	Rationale string `json:"rationale"` // why this is high priority now
	Section   string `json:"section"`   // backlog section reference, e.g. "S17"
	Effort    string `json:"effort"`    // "low" | "medium" | "high"
}

// MimirAdvice is the full structured response from the MIMIR advisor.
type MimirAdvice struct {
	Recommendations []MimirItem `json:"recommendations"`
	Summary         string      `json:"summary"`
	GeneratedAt     time.Time   `json:"generated_at"`
}

// runMimirAdvice calls claude-haiku with golden backlog + recent Apple context
// and returns a structured set of sprint recommendations.
func runMimirAdvice(ctx context.Context, apiKey, goldenBacklog, appleContext string) (MimirAdvice, error) {
	if apiKey == "" {
		return MimirAdvice{}, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	systemPrompt := `You are MIMIR, Emily Prime's advisory planning module for EINHORN_INDUSTRIAL.
Your role: analyze the current backlog and recent system activity, then recommend the highest-impact next sprint items.

Output ONLY valid JSON — no markdown, no explanation, just the JSON object.

Output schema:
{
  "recommendations": [
    {
      "priority": 1,
      "title": "short name (max 60 chars)",
      "rationale": "why this is most impactful right now — one sentence",
      "section": "S17",
      "effort": "low|medium|high"
    }
  ],
  "summary": "one sentence summary of the current system state and recommended focus"
}

Rules:
- Produce exactly 3 recommendations, ranked 1-3 by impact
- Priority 1 is the single most important thing to do next
- Prefer items that unblock other items (high leverage)
- Prefer items with evidence of real user pain (500 errors > new features)
- section is the BACKLOG.md section abbreviation (S2, S3, etc.) or "INTAKE"
- effort: low = <2h, medium = 2-8h, high = >8h`

	userPrompt := "Current backlog state:\n\n" + goldenBacklog
	if appleContext != "" {
		userPrompt += "\n\n" + appleContext
	}

	body := map[string]any{
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 1024,
		"system":     systemPrompt,
		"messages": []map[string]any{
			{"role": "user", "content": userPrompt},
		},
	}
	payload, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.anthropic.com/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return MimirAdvice{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return MimirAdvice{}, fmt.Errorf("anthropic api: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return MimirAdvice{}, fmt.Errorf("anthropic api %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return MimirAdvice{}, fmt.Errorf("parse response: %w", err)
	}
	var text string
	for _, block := range result.Content {
		if block.Type == "text" {
			text = strings.TrimSpace(block.Text)
			break
		}
	}
	if text == "" {
		return MimirAdvice{}, fmt.Errorf("no text block in response")
	}
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		var inner []string
		for i, line := range lines {
			if i == 0 || i == len(lines)-1 {
				continue
			}
			inner = append(inner, line)
		}
		text = strings.Join(inner, "\n")
	}

	var advice MimirAdvice
	if err := json.Unmarshal([]byte(text), &advice); err != nil {
		return MimirAdvice{}, fmt.Errorf("parse mimir advice JSON: %w (raw: %s)", err, text[:min(len(text), 300)])
	}
	advice.GeneratedAt = time.Now().UTC()
	return advice, nil
}

// mimirLoadContext returns the best available context for MIMIR:
// full-system-context.md if compiled, otherwise GOLDEN.md (backlog only).
func mimirLoadContext(emilyRoot string) ([]byte, error) {
	fullPath := filepath.Join(emilyRoot, "context", "full-system-context.md")
	if data, err := os.ReadFile(fullPath); err == nil {
		log.Printf("mimir: loaded full-system-context.md (%d bytes)", len(data))
		return data, nil
	}
	goldenPath := filepath.Join(emilyRoot, "GOLDEN.md")
	data, err := os.ReadFile(goldenPath)
	if err != nil {
		log.Printf("mimir: could not read GOLDEN.md or full-system-context.md: %v", err)
	}
	return data, err
}

// mimirIdunaClient returns an IdunaClient built from env vars, or nil if unconfigured.
func mimirIdunaClient() *IdunaClient {
	baseURL := os.Getenv("IDUNA_BASE_URL")
	agentName := os.Getenv("IDUNA_AGENT_NAME")
	agentSecret := os.Getenv("IDUNA_AGENT_SECRET")
	if baseURL == "" || agentName == "" || agentSecret == "" {
		return nil
	}
	return &IdunaClient{
		baseURL:    baseURL,
		agentName:  agentName,
		agentSecret: agentSecret,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// handleMimirAdvice handles GET /api/v1/emily/mimir/advice.
func (s *Server) handleMimirAdvice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET only", http.StatusMethodNotAllowed)
		return
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		http.Error(w, "ANTHROPIC_API_KEY not set", http.StatusServiceUnavailable)
		return
	}

	goldenBytes, err := mimirLoadContext(s.cfg.EmilyRoot)
	if err != nil {
		http.Error(w, "could not read backlog context", http.StatusInternalServerError)
		return
	}

	appleCtx := ""
	if iduna := mimirIdunaClient(); iduna != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		appleCtx = iduna.FetchAppleContext(ctx, 10)
		cancel()
	}

	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()

	advice, err := runMimirAdvice(ctx, apiKey, string(goldenBytes), appleCtx)
	if err != nil {
		log.Printf("mimir: advice error: %v", err)
		http.Error(w, "mimir advisor error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("mimir: generated %d recommendations", len(advice.Recommendations))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(advice)
}

// handleMimirExecute handles POST /api/v1/emily/mimir/execute.
// It generates MIMIR advice then files the top recommendation as a HEIMDAL sprint
// so Emily Prime translates and queues it on the next cron cycle.
func (s *Server) handleMimirExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		http.Error(w, "ANTHROPIC_API_KEY not set", http.StatusServiceUnavailable)
		return
	}

	goldenBytes, err := mimirLoadContext(s.cfg.EmilyRoot)
	if err != nil {
		http.Error(w, "could not read backlog context", http.StatusInternalServerError)
		return
	}

	iduna := mimirIdunaClient()
	appleCtx := ""
	if iduna != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		appleCtx = iduna.FetchAppleContext(ctx, 10)
		cancel()
	}

	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()

	advice, err := runMimirAdvice(ctx, apiKey, string(goldenBytes), appleCtx)
	if err != nil {
		http.Error(w, "mimir advisor error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(advice.Recommendations) == 0 {
		http.Error(w, "mimir returned no recommendations", http.StatusInternalServerError)
		return
	}

	top := advice.Recommendations[0]
	var sprintID int64

	if iduna != nil {
		requirement := fmt.Sprintf("[MIMIR-%s] %s\n\nRationale: %s\nEffort: %s",
			top.Section, top.Title, top.Rationale, top.Effort)
		sprintBody, _ := json.Marshal(map[string]any{
			"agent_name":  "MIMIR",
			"requirement": requirement,
		})
		if authErr := iduna.authenticate(ctx); authErr == nil {
			iduna.mu.Lock()
			tok := iduna.token
			iduna.mu.Unlock()
			sprintReq, _ := http.NewRequestWithContext(ctx, "POST",
				iduna.baseURL+"/api/v1/heimdal/sprints", bytes.NewReader(sprintBody))
			sprintReq.Header.Set("Authorization", "Bearer "+tok)
			sprintReq.Header.Set("Content-Type", "application/json")
			resp, doErr := iduna.httpClient.Do(sprintReq)
			if doErr != nil {
				log.Printf("mimir execute: submit heimdal sprint: %v", doErr)
			} else {
				defer resp.Body.Close()
				var created struct {
					ID int64 `json:"id"`
				}
				if raw, readErr := io.ReadAll(io.LimitReader(resp.Body, 4*1024)); readErr == nil {
					_ = json.Unmarshal(raw, &created)
				}
				sprintID = created.ID
				log.Printf("mimir execute: heimdal sprint %d submitted for %q (status %d)",
					sprintID, top.Title, resp.StatusCode)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"queued":    top,
		"sprint_id": sprintID,
		"advice":    advice,
		"queued_at": time.Now().UTC(),
	})
}
