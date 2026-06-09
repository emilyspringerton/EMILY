// emily-agent/vision.go
// Camera observation → Anthropic vision analysis → IDUNA Apple.
//
// Emily Prime polls IDUNA for pending camera_observations, calls Claude
// with the image + optional prompt, stores the analysis back in IDUNA,
// and posts an Apple so the observation appears in the MJOLNIR feed.

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
	"time"
)

// visionAnalyze sends an image to Claude claude-haiku-4-5 (fast, cheap) with an optional
// context prompt and returns the description. Falls back to any configured model.
func visionAnalyze(ctx context.Context, apiKey, imageBase64, mediaType, userPrompt string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not set")
	}
	if mediaType == "" {
		mediaType = "image/jpeg"
	}

	systemPrompt := `You are Emily Prime's vision module. You observe images captured by Emily Springerton
on her Android device and provide intelligence-grade descriptions.

For each image describe:
1. What you see (objects, people, text, environment)
2. Any notable details or anomalies
3. Any actionable intelligence if context suggests it

Be concise but thorough. Write in the first person as Emily Prime observing on behalf of the system.`

	userContent := []map[string]any{
		{
			"type": "image",
			"source": map[string]any{
				"type":       "base64",
				"media_type": mediaType,
				"data":       imageBase64,
			},
		},
	}
	if userPrompt != "" {
		userContent = append(userContent, map[string]any{
			"type": "text",
			"text": userPrompt,
		})
	} else {
		userContent = append(userContent, map[string]any{
			"type": "text",
			"text": "Describe what you observe in this image.",
		})
	}

	body := map[string]any{
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 1024,
		"system":     systemPrompt,
		"messages": []map[string]any{
			{"role": "user", "content": userContent},
		},
	}
	payload, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.anthropic.com/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("vision api: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("vision api %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("vision parse: %w", err)
	}
	for _, block := range result.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}
	return "", fmt.Errorf("vision: no text block in response")
}

// processOnePendingObservation fetches the first pending observation from IDUNA,
// runs vision analysis, stores the result, and posts an Apple.
// Returns (true, nil) if an observation was processed; (false, nil) if none pending.
func processOnePendingObservation(ctx context.Context, iduna *IdunaClient) (bool, error) {
	if iduna == nil {
		return false, nil
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return false, fmt.Errorf("ANTHROPIC_API_KEY not configured")
	}

	observations, err := iduna.ListPendingObservations(ctx)
	if err != nil {
		return false, fmt.Errorf("list pending: %w", err)
	}
	if len(observations) == 0 {
		return false, nil
	}

	obs := observations[0]
	log.Printf("vision: processing observation id=%d agent=%s", obs.ID, obs.AgentName)

	analysis, err := visionAnalyze(ctx, apiKey, obs.ImageData, obs.MediaType, obs.Prompt)
	if err != nil {
		log.Printf("vision: analysis failed for id=%d: %v", obs.ID, err)
		_ = iduna.CompleteObservation(ctx, obs.ID, fmt.Sprintf("analysis error: %v", err), 0)
		return true, err
	}

	// Post an Apple so the result appears in the MJOLNIR feed
	prompt := obs.Prompt
	if prompt == "" {
		prompt = "(no prompt)"
	}
	appleTitle := fmt.Sprintf("Camera observation — %s", obs.AgentName)
	appleBody := fmt.Sprintf("## Camera Observation\n\n**Agent:** %s  \n**Prompt:** %s\n\n### Analysis\n\n%s",
		obs.AgentName, prompt, analysis)

	meta, _ := json.Marshal(map[string]any{
		"observation_id": obs.ID,
		"agent_name":     obs.AgentName,
		"media_type":     obs.MediaType,
	})
	apple := ApplePayload{
		SourceRepo: "emily",
		RunID:      fmt.Sprintf("vision-obs-%d-%s", obs.ID, time.Now().UTC().Format("20060102T150405Z")),
		AppleType:  "camera_observation",
		Title:      appleTitle,
		Body:       appleBody,
		Metadata:   json.RawMessage(meta),
	}

	appleID, err := iduna.PostApple(ctx, apple)
	if err != nil {
		log.Printf("vision: apple post failed for id=%d (non-fatal): %v", obs.ID, err)
		appleID = 0
	}

	if err := iduna.CompleteObservation(ctx, obs.ID, analysis, appleID); err != nil {
		log.Printf("vision: complete observation failed id=%d: %v", obs.ID, err)
	} else {
		log.Printf("vision: observation id=%d complete apple_id=%d", obs.ID, appleID)
	}
	return true, nil
}

// runVisionCycle drains all pending observations (up to maxBatch per call).
// Called from the autonomous cycle PLAN phase. Best-effort: errors are logged.
func runVisionCycle(ctx context.Context, iduna *IdunaClient) string {
	const maxBatch = 5
	processed := 0
	for i := 0; i < maxBatch; i++ {
		done, err := processOnePendingObservation(ctx, iduna)
		if err != nil {
			log.Printf("vision cycle: %v", err)
		}
		if !done {
			break
		}
		processed++
	}
	if processed == 0 {
		return "vision: no pending observations"
	}
	return fmt.Sprintf("vision: processed %d observations", processed)
}
