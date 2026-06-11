// emily-agent/fable.go
// FABLE — Emily Prime's advisory planning tier.
//
// FABLE calls claude-haiku with the compressed golden backlog (GOLDEN.md) and
// recent Apple context to produce a prioritized sprint recommendation.
// It answers "what should Emily work on next?" in a structured, actionable form.
//
// Routes wired in main.go:
//   GET  /api/v1/emily/fable/advice    — generate a fresh FABLE sprint recommendation
//   POST /api/v1/emily/fable/execute   — push top recommendation into RSI roadmap via HEIMDAL

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

// FableItem is a single prioritized recommendation from the FABLE advisor.
type FableItem struct {
	Priority  int    `json:"priority"`  // 1 = highest
	Title     string `json:"title"`     // short name (≤60 chars)
	Rationale string `json:"rationale"` // why this is high priority now
	Section   string `json:"section"`   // backlog section reference, e.g. "S17"
	Effort    string `json:"effort"`    // "low" | "medium" | "high"
}

// FableAdvice is the full structured response from the FABLE advisor.
type FableAdvice struct {
	Recommendations []FableItem `json:"recommendations"`
	Summary         string      `json:"summary"`
	GeneratedAt     time.Time   `json:"generated_at"`
}

// runFableAdvice calls claude-haiku with golden backlog + recent Apple context
// and returns a structured set of sprint recommendations.
func runFableAdvice(ctx context.Context, apiKey, goldenBacklog, appleContext string) (FableAdvice, error) {
	if apiKey == "" {
		return FableAdvice{}, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	systemPrompt := `You are FABLE, Emily Prime's advisory planning module for EINHORN_INDUSTRIAL.
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
		return FableAdvice{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return FableAdvice{}, fmt.Errorf("anthropic api: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return FableAdvice{}, fmt.Errorf("anthropic api %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return FableAdvice{}, fmt.Errorf("parse response: %w", err)
	}
	var text string
	for _, block := range result.Content {
		if block.Type == "text" {
			text = strings.TrimSpace(block.Text)
			break
		}
	}
	if text == "" {
		return FableAdvice{}, fmt.Errorf("no text block in response")
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

	var advice FableAdvice
	if err := json.Unmarshal([]byte(text), &advice); err != nil {
		return FableAdvice{}, fmt.Errorf("parse fable advice JSON: %w (raw: %s)", err, text[:min(len(text), 300)])
	}
	advice.GeneratedAt = time.Now().UTC()
	return advice, nil
}

// fableIdunaClient returns an IdunaClient built from env vars, or nil if unconfigured.
func fableIdunaClient() *IdunaClient {
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

// handleFableAdvice handles GET /api/v1/emily/fable/advice.
func (s *Server) handleFableAdvice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET only", http.StatusMethodNotAllowed)
		return
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		http.Error(w, "ANTHROPIC_API_KEY not set", http.StatusServiceUnavailable)
		return
	}

	goldenPath := filepath.Join(s.cfg.EmilyRoot, "GOLDEN.md")
	goldenBytes, err := os.ReadFile(goldenPath)
	if err != nil {
		log.Printf("fable: read GOLDEN.md: %v", err)
		http.Error(w, "could not read backlog context", http.StatusInternalServerError)
		return
	}

	appleCtx := ""
	if iduna := fableIdunaClient(); iduna != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		appleCtx = iduna.FetchAppleContext(ctx, 10)
		cancel()
	}

	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()

	advice, err := runFableAdvice(ctx, apiKey, string(goldenBytes), appleCtx)
	if err != nil {
		log.Printf("fable: advice error: %v", err)
		http.Error(w, "fable advisor error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("fable: generated %d recommendations", len(advice.Recommendations))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(advice)
}

// handleFableExecute handles POST /api/v1/emily/fable/execute.
// It generates FABLE advice then files the top recommendation as a HEIMDAL sprint
// so Emily Prime translates and queues it on the next cron cycle.
func (s *Server) handleFableExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		http.Error(w, "ANTHROPIC_API_KEY not set", http.StatusServiceUnavailable)
		return
	}

	goldenPath := filepath.Join(s.cfg.EmilyRoot, "GOLDEN.md")
	goldenBytes, err := os.ReadFile(goldenPath)
	if err != nil {
		http.Error(w, "could not read backlog context", http.StatusInternalServerError)
		return
	}

	iduna := fableIdunaClient()
	appleCtx := ""
	if iduna != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		appleCtx = iduna.FetchAppleContext(ctx, 10)
		cancel()
	}

	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()

	advice, err := runFableAdvice(ctx, apiKey, string(goldenBytes), appleCtx)
	if err != nil {
		http.Error(w, "fable advisor error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(advice.Recommendations) == 0 {
		http.Error(w, "fable returned no recommendations", http.StatusInternalServerError)
		return
	}

	top := advice.Recommendations[0]
	var sprintID int64

	if iduna != nil {
		requirement := fmt.Sprintf("[FABLE-%s] %s\n\nRationale: %s\nEffort: %s",
			top.Section, top.Title, top.Rationale, top.Effort)
		sprintBody, _ := json.Marshal(map[string]any{
			"agent_name":  "FABLE",
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
				log.Printf("fable execute: submit heimdal sprint: %v", doErr)
			} else {
				defer resp.Body.Close()
				var created struct {
					ID int64 `json:"id"`
				}
				if raw, readErr := io.ReadAll(io.LimitReader(resp.Body, 4*1024)); readErr == nil {
					_ = json.Unmarshal(raw, &created)
				}
				sprintID = created.ID
				log.Printf("fable execute: heimdal sprint %d submitted for %q (status %d)",
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
