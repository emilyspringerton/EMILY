// emily-agent/heimdal.go
// HEIMDAL sprint planning cycle.
//
// Emily Prime polls IDUNA for pending sprint items submitted via the MJOLNIR
// sprint planning UI, translates each raw requirement into an ImprovementTask
// with AcceptanceCriteria (using Claude haiku), adds it to the RSI roadmap,
// files an Apple, and sends an FCM push to MJOLNIR.
//
// This bridges the gap between "human says what they want" and "Emily executes
// it the Emily way" — criteria-driven, iterative, Apple-on-completion.

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
	"strconv"
	"strings"
	"time"

	"emily-agent/pkg/fcm"
)

// sprintSummary mirrors the summary shape from GET /api/v1/heimdal/sprints.
type sprintSummary struct {
	ID                  int64  `json:"id"`
	AgentName           string `json:"agent_name"`
	RequirementPreview  string `json:"requirement_preview"`
	Requirement         string `json:"requirement"`
	RoadmapID           string `json:"roadmap_id"`
	Status              string `json:"status"`
	AppleID             int64  `json:"apple_id"`
}

// ListPendingHeimdalSprints returns all pending sprint items from IDUNA.
func (c *IdunaClient) ListPendingHeimdalSprints(ctx context.Context) ([]sprintSummary, error) {
	if err := c.authenticate(ctx); err != nil {
		return nil, fmt.Errorf("iduna authenticate: %w", err)
	}
	c.mu.Lock()
	tok := c.token
	c.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, "GET",
		c.baseURL+"/api/v1/heimdal/sprints?status=pending&limit=10", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("heimdal list: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("heimdal list status %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Sprints []sprintSummary `json:"sprints"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("heimdal list parse: %w", err)
	}
	return result.Sprints, nil
}

// PatchHeimdalSprint updates a sprint item after Emily Prime has processed it.
func (c *IdunaClient) PatchHeimdalSprint(ctx context.Context, id int64, criteriaJSON, roadmapID, status string, appleID int64) error {
	if err := c.authenticate(ctx); err != nil {
		return fmt.Errorf("iduna authenticate: %w", err)
	}
	c.mu.Lock()
	tok := c.token
	c.mu.Unlock()

	body, _ := json.Marshal(map[string]any{
		"criteria_json": criteriaJSON,
		"roadmap_id":    roadmapID,
		"status":        status,
		"apple_id":      appleID,
	})
	req, err := http.NewRequestWithContext(ctx, "PATCH",
		fmt.Sprintf("%s/api/v1/heimdal/sprints/%d", c.baseURL, id), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("heimdal patch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4*1024))
		return fmt.Errorf("heimdal patch status %d: %s", resp.StatusCode, raw)
	}
	return nil
}

// translateRequirement calls Claude haiku to translate a raw product requirement
// into a structured RoadmapItem with AcceptanceCriteria.
// appleContext is an optional compact summary of recent system activity injected
// as extra context; pass "" to skip.
func translateRequirement(ctx context.Context, apiKey, sprintID, requirement, appleContext string) (RoadmapItem, error) {
	if apiKey == "" {
		return RoadmapItem{}, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	systemPrompt := `You are Emily Prime's sprint planning module. You translate raw product requirements into structured RSI roadmap items.

Output ONLY valid JSON — no markdown, no explanation, just the JSON object.

Output schema:
{
  "id": "heimdal-<sprint_id>",
  "name": "concise name (max 60 chars)",
  "description": "fuller description of what needs to be built",
  "criteria": [
    {"name": "criterion_slug", "description": "what this criterion measures", "target": "what passing looks like — specific and measurable"},
    ...
  ],
  "max_iters": 5
}

Rules:
- id must be "heimdal-<sprint_id>" exactly
- criteria array must have 2-5 items
- Each criterion name is lowercase_snake_case
- Each target is specific and verifiable
- max_iters should be 3-8 depending on complexity`

	userPrompt := fmt.Sprintf("Sprint ID: %s\n\nRequirement:\n%s", sprintID, requirement)
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
		return RoadmapItem{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return RoadmapItem{}, fmt.Errorf("anthropic api: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return RoadmapItem{}, fmt.Errorf("anthropic api %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return RoadmapItem{}, fmt.Errorf("parse response: %w", err)
	}
	var text string
	for _, block := range result.Content {
		if block.Type == "text" {
			text = strings.TrimSpace(block.Text)
			break
		}
	}
	if text == "" {
		return RoadmapItem{}, fmt.Errorf("no text block in response")
	}

	// Strip markdown code fences if Claude wrapped in them
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

	var item RoadmapItem
	if err := json.Unmarshal([]byte(text), &item); err != nil {
		return RoadmapItem{}, fmt.Errorf("parse roadmap item JSON: %w (raw: %s)", err, text[:min(len(text), 200)])
	}
	if item.ID == "" {
		item.ID = "heimdal-" + sprintID
	}
	if item.MaxIters <= 0 {
		item.MaxIters = 5
	}
	item.Status = "queued"
	item.Priority = 10
	return item, nil
}

// runHeimdalCycle is called in the PLAN phase each cron cycle.
// It drains up to 5 pending HEIMDAL sprints, translating each into an RSI
// roadmap item and firing an FCM push to MJOLNIR.
func (ac *AutonomousCycle) runHeimdalCycle(ctx context.Context, push PushFunc) string {
	if ac.iduna == nil {
		return "heimdal: iduna client not configured"
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "heimdal: ANTHROPIC_API_KEY not set"
	}

	sprints, err := ac.iduna.ListPendingHeimdalSprints(ctx)
	if err != nil {
		log.Printf("heimdal: list pending sprints: %v", err)
		return fmt.Sprintf("heimdal: list error: %v", err)
	}
	if len(sprints) == 0 {
		return "heimdal: no pending sprints"
	}

	// Fetch recent Apple context once for the whole batch (≤200 tokens).
	appleCtxStr := ""
	if appleCtxCtx, appleCtxCancel := context.WithTimeout(ctx, 8*time.Second); true {
		appleCtxStr = ac.iduna.FetchAppleContext(appleCtxCtx, 10)
		appleCtxCancel()
	}

	processed := 0
	for i, sprint := range sprints {
		if i >= 5 {
			break
		}

		sprintIDStr := fmt.Sprintf("%d", sprint.ID)
		log.Printf("heimdal: processing sprint id=%d agent=%s", sprint.ID, sprint.AgentName)

		translateCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		item, err := translateRequirement(translateCtx, apiKey, sprintIDStr, sprint.Requirement, appleCtxStr)
		cancel()
		if err != nil {
			log.Printf("heimdal: translate sprint %d: %v", sprint.ID, err)
			continue
		}

		// Add to the RSI roadmap
		if addErr := ac.AddRoadmapItem(item); addErr != nil {
			log.Printf("heimdal: add roadmap item %s: %v", item.ID, addErr)
			continue
		}

		// Serialize criteria for IDUNA patch
		criteriaRaw, _ := json.Marshal(item.Criteria)

		// Post Apple
		req := sprint.Requirement
		if len(req) > 200 {
			req = req[:197] + "..."
		}
		appleTitle := fmt.Sprintf("HEIMDAL sprint queued: %s", item.Name)
		appleBody := fmt.Sprintf("## HEIMDAL Sprint\n\n**Requirement:** %s\n\n**Roadmap ID:** %s\n\n**Criteria:**\n\n%s",
			req, item.ID, formatCriteria(item.Criteria))
		meta, _ := json.Marshal(map[string]any{
			"sprint_id":  sprint.ID,
			"agent_name": sprint.AgentName,
			"roadmap_id": item.ID,
		})

		var appleID int64
		if appleCtx, appleCancel := context.WithTimeout(ctx, 10*time.Second); true {
			id, appleErr := ac.iduna.PostApple(appleCtx, ApplePayload{
				SourceRepo: "emily",
				RunID:      fmt.Sprintf("heimdal-sprint-%d-%s", sprint.ID, time.Now().UTC().Format("20060102T150405Z")),
				AppleType:  "improvement",
				Title:      appleTitle,
				Body:       appleBody,
				Metadata:   json.RawMessage(meta),
			})
			appleCancel()
			if appleErr != nil {
				log.Printf("heimdal: post apple for sprint %d: %v", sprint.ID, appleErr)
			}
			appleID = id
		}

		// Patch sprint status in IDUNA
		patchCtx, patchCancel := context.WithTimeout(ctx, 10*time.Second)
		patchErr := ac.iduna.PatchHeimdalSprint(patchCtx, sprint.ID,
			string(criteriaRaw), item.ID, "queued", appleID)
		patchCancel()
		if patchErr != nil {
			log.Printf("heimdal: patch sprint %d: %v", sprint.ID, patchErr)
		}

		// FCM push to MJOLNIR
		if push != nil {
			push("Sprint queued: "+item.Name,
				fmt.Sprintf("%d criteria · %d max iters", len(item.Criteria), item.MaxIters),
				map[string]string{"type": "heimdal_sprint", "roadmap_id": item.ID},
			)
		}

		processed++
		log.Printf("heimdal: sprint %d queued as %s (%d criteria)", sprint.ID, item.ID, len(item.Criteria))
	}

	return fmt.Sprintf("heimdal: processed %d/%d pending sprints", processed, len(sprints))
}

// notifyHeimdalStatus is called in a goroutine when a heimdal-* RSI task reaches
// a terminal state (complete or blocked). It patches the IDUNA sprint record,
// files a completion Apple, and fires an FCM push to MJOLNIR.
func (ac *AutonomousCycle) notifyHeimdalStatus(task *ImprovementTask, newStatus string) {
	if ac.iduna == nil || task == nil {
		return
	}
	sprintIDStr := strings.TrimPrefix(task.ID, "heimdal-")
	sprintID, err := strconv.ParseInt(sprintIDStr, 10, 64)
	if err != nil {
		log.Printf("heimdal notify: parse sprint id from %q: %v", task.ID, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	appleTitle := fmt.Sprintf("HEIMDAL sprint %s: %s", newStatus, task.ID)
	appleBody := fmt.Sprintf("## HEIMDAL Sprint %s\n\n**Task:** %s  \n**Status:** %s  \n**Iterations:** %d/%d",
		newStatus, task.ID, newStatus, len(task.Iterations), task.MaxIters)
	if len(task.Lessons) > 0 {
		appleBody += "\n\n**Lessons:**\n"
		for _, l := range task.Lessons {
			appleBody += "- " + l + "\n"
		}
	}
	meta, _ := json.Marshal(map[string]any{
		"sprint_id":  sprintID,
		"task_id":    task.ID,
		"status":     newStatus,
		"iterations": len(task.Iterations),
	})

	appleID, appleErr := ac.iduna.PostApple(ctx, ApplePayload{
		SourceRepo: "emily",
		RunID:      fmt.Sprintf("heimdal-done-%d-%s", sprintID, time.Now().UTC().Format("20060102T150405Z")),
		AppleType:  "improvement",
		Title:      appleTitle,
		Body:       appleBody,
		Metadata:   json.RawMessage(meta),
	})
	if appleErr != nil {
		log.Printf("heimdal notify: post apple for sprint %d: %v", sprintID, appleErr)
	}

	patchCtx, patchCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer patchCancel()
	if patchErr := ac.iduna.PatchHeimdalSprint(patchCtx, sprintID, "", task.ID, newStatus, appleID); patchErr != nil {
		log.Printf("heimdal notify: patch sprint %d to %s: %v", sprintID, newStatus, patchErr)
	} else {
		log.Printf("heimdal notify: sprint %d patched to %s (apple=%d)", sprintID, newStatus, appleID)
	}

	if ac.fcmSender == nil {
		return
	}
	pushCtx, pushCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pushCancel()
	deviceToken, err := ac.iduna.GetPushToken(pushCtx, "mjolnir-emily")
	if err != nil || deviceToken == "" {
		return
	}
	emoji := "✅"
	if newStatus == "blocked" {
		emoji = "⚠️"
	}
	_ = ac.fcmSender.Send(pushCtx, deviceToken, fcm.Message{
		Title:    fmt.Sprintf("%s Sprint %s: %s", emoji, newStatus, task.ID),
		Body:     fmt.Sprintf("%d iterations · %d criteria", len(task.Iterations), len(task.Criteria)),
		Priority: "high",
		Data:     map[string]string{"type": "heimdal_done", "task_id": task.ID, "status": newStatus},
	})
}

func formatCriteria(criteria []AcceptanceCriterion) string {
	var sb strings.Builder
	for i, c := range criteria {
		sb.WriteString(fmt.Sprintf("%d. **%s** — %s  \n   Target: %s\n", i+1, c.Name, c.Description, c.Target))
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
