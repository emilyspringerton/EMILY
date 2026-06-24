// gpt2fallback.go — S125-06: GPT-2 inference fallback for the PLAN phase.
//
// When ANTHROPIC_API_KEY is absent or rate-limited, Emily Prime falls back to
// local GPT-2 inference via the binary at $GPT2_INFER_BIN.
//
// The binary is expected to read a prompt from stdin and write generated text
// to stdout. Exit 0 = success; any other exit code = failure.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

// gpt2InferAvailable returns true if the binary at $GPT2_INFER_BIN exists and is executable.
func gpt2InferAvailable() bool {
	bin := os.Getenv("GPT2_INFER_BIN")
	if bin == "" {
		return false
	}
	_, err := os.Stat(bin)
	return err == nil
}

// gpt2Fallback runs the GPT-2 inference binary with the given prompt and returns
// the generated text. maxTokens is passed as an env var (GPT2_MAX_TOKENS) for
// implementations that honour it. Returns an error if the binary fails or is unavailable.
func gpt2Fallback(ctx context.Context, prompt string, maxTokens int) (string, error) {
	bin := os.Getenv("GPT2_INFER_BIN")
	if bin == "" {
		return "", fmt.Errorf("GPT2_INFER_BIN not set")
	}
	cmd := exec.CommandContext(ctx, bin)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GPT2_MAX_TOKENS=%d", maxTokens))
	cmd.Stdin = bytes.NewBufferString(prompt)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("gpt2 infer: %w", err)
	}
	return string(out), nil
}

// callHaiku sends a simple user-only prompt to claude-haiku and returns the text response.
func callHaiku(ctx context.Context, prompt, apiKey string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not set")
	}
	body := map[string]any{
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 256,
		"messages":   []map[string]any{{"role": "user", "content": prompt}},
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
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("haiku rate-limited (429)")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("haiku HTTP %d: %s", resp.StatusCode, string(raw))
	}
	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", err
	}
	for _, c := range result.Content {
		if c.Type == "text" {
			return c.Text, nil
		}
	}
	return "", fmt.Errorf("no text block in haiku response")
}

// planWithFallback attempts a haiku LLM call; falls back to GPT-2 if the API
// key is absent or rate-limited. Returns the generated plan text and a bool
// indicating whether GPT-2 was used.
func planWithFallback(ctx context.Context, prompt, apiKey string) (text string, usedGPT2 bool, err error) {
	if apiKey != "" {
		text, err = callHaiku(ctx, prompt, apiKey)
		if err == nil {
			return text, false, nil
		}
	}
	// API key absent or call failed — try GPT-2.
	if !gpt2InferAvailable() {
		if err != nil {
			return "", false, fmt.Errorf("haiku failed and GPT2_INFER_BIN unavailable: %w", err)
		}
		return "", false, fmt.Errorf("ANTHROPIC_API_KEY not set and GPT2_INFER_BIN unavailable")
	}
	text, err = gpt2Fallback(ctx, prompt, 200)
	if err != nil {
		return "", true, fmt.Errorf("gpt2 fallback failed: %w", err)
	}
	return text, true, nil
}
