// emily-agent/plan.go
// POST /api/v1/emily/plan — answer a planning question via claude-haiku.
//
// Routes wired in main.go:
//   POST /api/v1/emily/plan — accept {question, context?}, return {sprints, summary}
//
// Offloads sprint-batch planning from Claude Code (sonnet tokens) to haiku.
// Reads full-system-context.md (goldenbuild output) or falls back to GOLDEN.md.

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
	"strings"
	"time"
)

// PlanRequest is the body for POST /api/v1/emily/plan.
type PlanRequest struct {
	Question string `json:"question"`
	Context  string `json:"context,omitempty"`
}

// PlanSprint is one sprint item in the plan response.
type PlanSprint struct {
	Title     string `json:"title"`
	Rationale string `json:"rationale"`
	Section   string `json:"section"`
	Effort    string `json:"effort"`
}

// PlanResponse is returned by POST /api/v1/emily/plan.
type PlanResponse struct {
	Sprints     []PlanSprint `json:"sprints"`
	Summary     string       `json:"summary"`
	GeneratedAt time.Time    `json:"generated_at"`
}

const planSystemPrompt = `You are Emily Prime's planning module for EINHORN_INDUSTRIAL.
You are given a planning question and the current system backlog context.
Produce a structured, actionable sprint batch answering the question.

Output ONLY valid JSON — no markdown, no explanation, no prose outside the JSON.

Output schema:
{
  "sprints": [
    {
      "title": "short actionable title (max 60 chars)",
      "rationale": "why this sprint answers the question — one sentence",
      "section": "backlog section (e.g. S22, INTAKE) or empty string",
      "effort": "low|medium|high"
    }
  ],
  "summary": "one sentence describing the recommended focus and why"
}

Rules:
- Produce 2–5 sprints, ordered by impact (highest first)
- effort: low = <2h, medium = 2–8h, high = >8h
- Prefer high-leverage items that unblock others
- Prefer items with concrete acceptance criteria over vague explorations
- section must match existing BACKLOG sections if possible, else empty
- Frame-break lens (REDGARDEN/NORTHSTAR.md §28): before picking sprints, consider what general
  structural pattern the planning question is one instance of — this can sharpen which sprints
  you pick and their rationale, but never changes the output contract above; still ONLY the JSON
  schema, no visible reframing/commentary in the output itself`

// runPlan calls claude-haiku with backlog context + the user's question
// and returns a structured sprint batch.
func runPlan(ctx context.Context, apiKey, systemContext, question, extraContext string) (PlanResponse, error) {
	if apiKey == "" {
		return PlanResponse{}, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	userPrompt := "System backlog context:\n\n" + systemContext
	userPrompt += "\n\nPlanning question: " + question
	if extraContext != "" {
		userPrompt += "\n\nAdditional context:\n" + extraContext
	}

	body := map[string]any{
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 1024,
		"system":     planSystemPrompt,
		"messages": []map[string]any{
			{"role": "user", "content": userPrompt},
		},
	}
	payload, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.anthropic.com/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return PlanResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return PlanResponse{}, fmt.Errorf("anthropic api: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return PlanResponse{}, fmt.Errorf("anthropic api %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return PlanResponse{}, fmt.Errorf("parse response: %w", err)
	}
	var text string
	for _, block := range result.Content {
		if block.Type == "text" {
			text = strings.TrimSpace(block.Text)
			break
		}
	}
	if text == "" {
		return PlanResponse{}, fmt.Errorf("no text block in response")
	}
	// Strip markdown code fences if haiku wrapped the JSON.
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

	var plan PlanResponse
	if err := json.Unmarshal([]byte(text), &plan); err != nil {
		return PlanResponse{}, fmt.Errorf("parse plan JSON: %w (raw: %s)", err, text[:min(len(text), 300)])
	}
	plan.GeneratedAt = time.Now().UTC()
	return plan, nil
}

// handlePlan handles POST /api/v1/emily/plan.
func (s *Server) handlePlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		http.Error(w, "ANTHROPIC_API_KEY not set", http.StatusServiceUnavailable)
		return
	}

	var req PlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Question) == "" {
		http.Error(w, `"question" is required`, http.StatusBadRequest)
		return
	}

	sysCtxBytes, err := mimirLoadContext(s.cfg.EmilyRoot)
	if err != nil {
		http.Error(w, "could not read backlog context", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()

	plan, err := runPlan(ctx, apiKey, string(sysCtxBytes), req.Question, req.Context)
	if err != nil {
		log.Printf("plan: error: %v", err)
		http.Error(w, "plan error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("plan: generated %d sprints for question: %q", len(plan.Sprints), req.Question)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}
