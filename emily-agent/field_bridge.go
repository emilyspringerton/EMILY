package main

// field_bridge.go: wires the pkg/archetypes Field into the RSI AutonomousCycle.
// The archetypes package defines its own LLMClient interface so it can be imported
// without circular dependency. This file bridges main.AnthropicClient → archetypes.LLMClient.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"emily-agent/pkg/archetypes"
)

// archetypesBridge adapts main.AnthropicClient to archetypes.LLMClient.
// It re-implements the raw HTTP call using the same endpoint and auth as AnthropicClient.
type archetypesBridge struct {
	apiKey  string
	baseURL string
}

func newArchetypesBridge(apiKey string) *archetypesBridge {
	return &archetypesBridge{apiKey: apiKey, baseURL: "https://api.anthropic.com/v1"}
}

func (b *archetypesBridge) Complete(ctx context.Context, req archetypes.LLMRequest) (archetypes.LLMResponse, error) {
	body := map[string]any{
		"model":       req.Model,
		"max_tokens":  req.MaxTokens,
		"temperature": req.Temperature,
		"system":      req.SystemProm,
		"messages": []map[string]string{
			{"role": "user", "content": req.UserPrompt},
		},
	}
	raw, _ := json.Marshal(body)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", b.baseURL+"/messages", bytes.NewReader(raw))
	if err != nil {
		return archetypes.LLMResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", b.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return archetypes.LLMResponse{}, err
	}
	defer resp.Body.Close()

	var out struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return archetypes.LLMResponse{}, fmt.Errorf("anthropic decode: %w", err)
	}
	if out.Error != nil {
		return archetypes.LLMResponse{}, fmt.Errorf("anthropic: %s", out.Error.Message)
	}
	if len(out.Content) == 0 {
		return archetypes.LLMResponse{}, fmt.Errorf("anthropic: empty response")
	}
	return archetypes.LLMResponse{
		Content:      out.Content[0].Text,
		InputTokens:  out.Usage.InputTokens,
		OutputTokens: out.Usage.OutputTokens,
	}, nil
}

// NewFieldFromEnv builds an archetypes.Field using ANTHROPIC_API_KEY.
// Returns nil if the key is not set; caller must nil-check before use.
func NewFieldFromEnv() *archetypes.Field {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		return nil
	}
	bridge := newArchetypesBridge(key)
	return &archetypes.Field{
		E1Client: bridge,
		E2Client: bridge,
		E1Model:  "claude-sonnet-4-6",
		E2Model:  "claude-haiku-4-5-20251001",
	}
}

// AugmentTaskWithField runs a Field invocation for the task's intent and prepends
// the blended resonance output to the task description. Called in the RSI DECIDE
// phase after pickTask selects a new task. Operates with a short timeout so it
// never blocks the cycle.
func AugmentTaskWithField(ctx context.Context, field *archetypes.Field, task *ImprovementTask) {
	if field == nil || task == nil {
		return
	}
	// Only augment newly-promoted tasks, not ones already running.
	if task.Status != "pending" {
		return
	}

	// Derive intent from the task description (first 120 chars).
	intent := task.Description
	if len(intent) > 120 {
		intent = intent[:120]
	}

	// 30-second budget per augmentation — must not block the full 5-min cycle.
	augCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := field.Invoke(augCtx, intent, task.Description, false)
	if err != nil {
		log.Printf("[field] augment task %s: %v", task.ID, err)
		return
	}

	// Prepend the resonance header to the task description.
	header := fmt.Sprintf("[FIELD Ψ(t) | %s | Δφ=%.0f° | E1=%.0f%% E2=%.0f%% | spirits=%s]\n%s\n\n---\n\n",
		res.ResonanceState,
		res.PhaseDeltaDeg,
		res.E1Weight*100, res.E2Weight*100,
		spiritStack(res.ActiveSpirits),
		res.Output,
	)
	task.Description = header + task.Description

	log.Printf("[field] task %s augmented: state=%s Δφ=%.0f° spirits=%s",
		task.ID, res.ResonanceState, res.PhaseDeltaDeg, spiritStack(res.ActiveSpirits))
}

func spiritStack(ss []archetypes.Spirit) string {
	names := make([]string, len(ss))
	for i, s := range ss {
		names[i] = s.Name
	}
	return strings.Join(names, "+")
}
