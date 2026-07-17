// emily-agent/enrichworker.go
// GPT-2 tower Apple enrichment worker (S147-02).
//
// Polls IDUNA for recent Apples missing a GPT-2 tower fingerprint, generates
// text via serve.py, runs it through towerprint.Compute() in-process (no
// Python for the transform itself), and PATCHes the result back. Per
// gpt2-alpine-c/docs/TOWERPRINT.md §5's decision: async, caller-side — the
// Apple POST never blocks on :8088, and this worker's failure mode is
// "field missing, retried next poll," never a lost/slow Apple. Same shape
// as CheckinAlertWorker (alerting.go).

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gpt2-alpine-c/pkg/towerprint"
)

const (
	enrichPollInterval = 5 * time.Minute
	enrichBatchLimit   = 20 // recent Apples inspected per poll
)

// enrichGenerateClient has no Timeout of its own — the 60s deadline in
// enrichOne's genCtx governs each generate call instead. gpt2HealthClient
// (gpt2tools.go) has a 5s timeout meant for /health pings, too short for a
// real 100-token generation.
var enrichGenerateClient = &http.Client{}

// ApplEnrichWorker polls IDUNA and enriches Apples missing gpt2_fingerprint.
type ApplEnrichWorker struct {
	Iduna *IdunaClient
}

// NewApplEnrichWorker builds a worker from an existing IDUNA client.
func NewApplEnrichWorker(iduna *IdunaClient) *ApplEnrichWorker {
	return &ApplEnrichWorker{Iduna: iduna}
}

// Run starts the polling loop. Blocks until ctx is cancelled.
func (w *ApplEnrichWorker) Run(ctx context.Context) {
	if w.Iduna == nil {
		log.Println("[enrich] IDUNA client not configured — GPT-2 tower enrichment disabled")
		return
	}
	log.Printf("[enrich] GPT-2 tower enrichment worker started (poll interval %s)", enrichPollInterval)
	ticker := time.NewTicker(enrichPollInterval)
	defer ticker.Stop()
	w.poll(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.poll(ctx)
		}
	}
}

func (w *ApplEnrichWorker) poll(ctx context.Context) {
	apples, err := w.Iduna.ListApples(ctx, "", "", enrichBatchLimit)
	if err != nil {
		log.Printf("[enrich] list apples: %v", err)
		return
	}
	var candidates []AppleListItem
	for _, a := range apples {
		if !a.HasGpt2Fingerprint {
			candidates = append(candidates, a)
		}
	}
	if len(candidates) == 0 {
		return
	}
	log.Printf("[enrich] %d apple(s) missing gpt2_fingerprint", len(candidates))
	for _, a := range candidates {
		w.enrichOne(ctx, a)
	}
}

func (w *ApplEnrichWorker) enrichOne(ctx context.Context, a AppleListItem) {
	genCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	text, model, err := generateForFingerprint(genCtx, a.Title)
	if err != nil {
		// Expected and non-alarming when serve.py isn't running — this
		// Apple stays a candidate and is retried on the next poll, exactly
		// the "field missing, retried later" degradation TOWERPRINT.md §5
		// calls for.
		log.Printf("[enrich] apple %d: generate failed (will retry): %v", a.ID, err)
		return
	}

	fp, err := towerprint.Compute(text)
	if err != nil {
		// A generation with no letters at all is unusual but not an error
		// worth logging loudly every poll — quiet skip, will retry.
		log.Printf("[enrich] apple %d: towerprint compute failed: %v", a.ID, err)
		return
	}

	fpJSON, err := json.Marshal(map[string]any{
		"generated":   text,
		"squished":    fp.Squished,
		"tower":       fp.Tower,
		"magic_tower": fp.Magic,
		"codze":       fp.Codze,
		"seed":        towerprint.Seed(time.Now()),
	})
	if err != nil {
		log.Printf("[enrich] apple %d: marshal fingerprint: %v", a.ID, err)
		return
	}

	updates := map[string]json.RawMessage{"gpt2_fingerprint": fpJSON}
	if model != "" {
		if modelJSON, err := json.Marshal(model); err == nil {
			updates["model_fingerprint"] = modelJSON
		}
	}

	patchCtx, patchCancel := context.WithTimeout(ctx, 15*time.Second)
	defer patchCancel()
	if err := w.Iduna.PatchApple(patchCtx, a.ID, updates); err != nil {
		log.Printf("[enrich] apple %d: patch failed (will retry): %v", a.ID, err)
		return
	}
	log.Printf("[enrich] apple %d enriched (model=%q)", a.ID, model)
}

// generateForFingerprint calls serve.py's /generate with prompt (the
// Apple's title, per TOWERPRINT.md §5's "suggested prompt") and returns the
// generated continuation plus the model tag serve.py reports. Separate from
// gpt2tools.go's gpt2_generate tool callback — that one is shaped for the
// LLM tool-dispatcher's args-map calling convention; this is a plain,
// context-aware Go call for a background worker.
func generateForFingerprint(ctx context.Context, prompt string) (text, model string, err error) {
	reqBody, err := json.Marshal(map[string]any{
		"prompt":      prompt,
		"max_tokens":  100,
		"temperature": 0.8,
	})
	if err != nil {
		return "", "", fmt.Errorf("marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, gpt2ServeURL+"/generate", bytes.NewReader(reqBody))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := enrichGenerateClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("gpt2 server unreachable: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("gpt2 server returned %d: %s", resp.StatusCode, raw)
	}

	var out struct {
		Text  string `json:"text"`
		Model string `json:"model"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", "", fmt.Errorf("parse response: %w", err)
	}
	if out.Error != "" {
		return "", "", fmt.Errorf("gpt2 generate error: %s", out.Error)
	}
	return out.Text, out.Model, nil
}
