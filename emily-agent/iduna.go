// emily-agent/iduna.go
// IDUNA IAM client — authenticates Emily as an M2M agent and submits Apples
// (golden documentation records) to the Apples ledger after each RSI run.
//
// Spec: HQ-SPEC-IAM-096. Requires env vars:
//   IDUNA_BASE_URL    — e.g. https://iam.farthq.internal
//   IDUNA_AGENT_NAME  — registered agent name (e.g. EMILY_PRIME)
//   IDUNA_AGENT_SECRET — M2M credential set via the IDUNA Back Office
//
// When any of these env vars are absent, the client is nil and all calls
// are no-ops. Apple submission failures are always non-fatal.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// IdunaClient authenticates as an agent and submits Apples to IDUNA.
type IdunaClient struct {
	baseURL     string
	agentName   string
	agentSecret string
	httpClient  *http.Client

	mu           sync.Mutex
	token        string
	tokenExp     time.Time
	applesGitDir string // if set, runs emily sync after each successful Apple POST
}

// NewIdunaClientFromEnv reads IDUNA_BASE_URL, IDUNA_AGENT_NAME, IDUNA_AGENT_SECRET
// and returns a configured client. Returns nil if any required var is absent.
// Optionally reads APPLES_GIT_DIR: when set, each successful Apple POST triggers
// an async `emily sync --apples-git-dir <dir>` to keep the git backup in sync.
func NewIdunaClientFromEnv() *IdunaClient {
	base := envOr("IDUNA_BASE_URL", "")
	name := envOr("IDUNA_AGENT_NAME", "")
	secret := envOr("IDUNA_AGENT_SECRET", "")
	if base == "" || name == "" || secret == "" {
		return nil
	}
	return &IdunaClient{
		baseURL:      base,
		agentName:    name,
		agentSecret:  secret,
		httpClient:   &http.Client{Timeout: 15 * time.Second},
		applesGitDir: envOr("APPLES_GIT_DIR", ""),
	}
}

// authenticate fetches a fresh JWT if none exists or it expires within 5 minutes.
func (c *IdunaClient) authenticate(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && time.Until(c.tokenExp) > 5*time.Minute {
		return nil
	}
	body, _ := json.Marshal(map[string]string{
		"agent_name":   c.agentName,
		"agent_secret": c.agentSecret,
	})
	req, err := http.NewRequestWithContext(ctx, "POST",
		c.baseURL+"/api/v1/auth/agent", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("iduna auth: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("iduna auth status %d: %s", resp.StatusCode, string(raw))
	}
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return fmt.Errorf("iduna auth parse: %w", err)
	}
	c.token = result.AccessToken
	c.tokenExp = time.Now().UTC().Add(time.Duration(result.ExpiresIn) * time.Second)
	return nil
}

// ApplePayload is the request body for POST /api/v1/apples.
type ApplePayload struct {
	SourceRepo string          `json:"source_repo"`
	RunID      string          `json:"run_id"`
	AppleType  string          `json:"apple_type"`
	Title      string          `json:"title"`
	Body       string          `json:"body"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
}

// PostApple authenticates (if needed) and submits a golden documentation record.
// Returns the assigned Apple ID. Errors are logged but never fatal to the caller.
func (c *IdunaClient) PostApple(ctx context.Context, payload ApplePayload) (int64, error) {
	if err := c.authenticate(ctx); err != nil {
		return 0, fmt.Errorf("iduna authenticate: %w", err)
	}

	c.mu.Lock()
	tok := c.token
	c.mu.Unlock()

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST",
		c.baseURL+"/api/v1/apples", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("iduna post apple: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 16*1024))
	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("iduna post apple status %d: %s", resp.StatusCode, string(raw))
	}
	var result struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return 0, fmt.Errorf("iduna post apple parse: %w", err)
	}

	if c.applesGitDir != "" && result.ID > 0 {
		dir := c.applesGitDir
		go func() {
			cmd := exec.Command("emily", "sync", "--apples-git-dir", dir)
			if out, err := cmd.CombinedOutput(); err != nil {
				log.Printf("apples sync: %v: %s", err, out)
			}
		}()
	}

	return result.ID, nil
}

// GetPushToken returns the current FCM device token for the given agent name.
// Returns ("", nil) if no token is registered. Used by Emily Prime before
// firing FCM push notifications to MJOLNIR.
func (c *IdunaClient) GetPushToken(ctx context.Context, agentName string) (string, error) {
	if err := c.authenticate(ctx); err != nil {
		return "", fmt.Errorf("iduna authenticate: %w", err)
	}
	c.mu.Lock()
	tok := c.token
	c.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, "GET",
		c.baseURL+"/api/v1/push-tokens/"+agentName, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("iduna get push token: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4*1024))
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("iduna get push token status %d: %s", resp.StatusCode, raw)
	}
	var result struct {
		FCMToken string `json:"fcm_token"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("iduna get push token parse: %w", err)
	}
	return result.FCMToken, nil
}

// AppleListItem is a summary record from GET /api/v1/apples.
type AppleListItem struct {
	ID                 int64     `json:"id"`
	AgentID            string    `json:"agent_id"`
	SourceRepo         string    `json:"source_repo"`
	RunID              string    `json:"run_id"`
	AppleType          string    `json:"apple_type"`
	Title              string    `json:"title"`
	RecordedAt         time.Time `json:"recorded_at"`
	HasGpt2Fingerprint bool      `json:"has_gpt2_fingerprint"`
}

// ListApples returns up to limit Apples, optionally filtered by apple_type and source_repo.
// Empty strings skip the filter. Apples are returned newest-first.
func (c *IdunaClient) ListApples(ctx context.Context, appleType, sourceRepo string, limit int) ([]AppleListItem, error) {
	if err := c.authenticate(ctx); err != nil {
		return nil, fmt.Errorf("iduna authenticate: %w", err)
	}
	c.mu.Lock()
	tok := c.token
	c.mu.Unlock()

	u := c.baseURL + fmt.Sprintf("/api/v1/apples?limit=%d", limit)
	if appleType != "" {
		u += "&apple_type=" + appleType
	}
	if sourceRepo != "" {
		u += "&source_repo=" + sourceRepo
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("iduna list apples: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iduna list apples status %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Apples []struct {
			ID                 int64  `json:"id"`
			AgentID            string `json:"agent_id"`
			SourceRepo         string `json:"source_repo"`
			RunID              string `json:"run_id"`
			AppleType          string `json:"apple_type"`
			Title              string `json:"title"`
			RecordedAt         string `json:"recorded_at"`
			HasGpt2Fingerprint bool   `json:"has_gpt2_fingerprint"`
		} `json:"apples"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("iduna list apples parse: %w", err)
	}

	items := make([]AppleListItem, 0, len(result.Apples))
	for _, a := range result.Apples {
		t, _ := time.Parse(time.RFC3339Nano, a.RecordedAt)
		items = append(items, AppleListItem{
			ID:                 a.ID,
			AgentID:            a.AgentID,
			SourceRepo:         a.SourceRepo,
			RunID:              a.RunID,
			AppleType:          a.AppleType,
			Title:              a.Title,
			RecordedAt:         t,
			HasGpt2Fingerprint: a.HasGpt2Fingerprint,
		})
	}
	return items, nil
}

// PatchApple merges the given enrichment fields into an Apple's metadata via
// PATCH /api/v1/apples/{id}. Used by the async enrichment worker (S147-02)
// — never called from the hot POST path, per TOWERPRINT.md §5's decision
// that enrichment must never block or risk an Apple filing.
func (c *IdunaClient) PatchApple(ctx context.Context, id int64, updates map[string]json.RawMessage) error {
	if err := c.authenticate(ctx); err != nil {
		return fmt.Errorf("iduna authenticate: %w", err)
	}
	c.mu.Lock()
	tok := c.token
	c.mu.Unlock()

	body, err := json.Marshal(updates)
	if err != nil {
		return fmt.Errorf("iduna patch apple: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.baseURL+fmt.Sprintf("/api/v1/apples/%d", id), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("iduna patch apple: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return fmt.Errorf("iduna patch apple status %d: %s", resp.StatusCode, raw)
	}
	return nil
}

// CameraObservation is a pending intelligence request from MJOLNIR.
type CameraObservation struct {
	ID        int64  `json:"id"`
	AgentName string `json:"agent_name"`
	ImageData string `json:"image_data"` // base64
	MediaType string `json:"media_type"`
	Prompt    string `json:"prompt"`
	Status    string `json:"status"`
}

// ListPendingObservations returns camera observations with status=pending.
func (c *IdunaClient) ListPendingObservations(ctx context.Context) ([]CameraObservation, error) {
	if err := c.authenticate(ctx); err != nil {
		return nil, fmt.Errorf("iduna authenticate: %w", err)
	}
	c.mu.Lock()
	tok := c.token
	c.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, "GET",
		c.baseURL+"/api/v1/intelligence/observations?status=pending&limit=10", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("iduna list observations: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iduna list observations status %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Observations []CameraObservation `json:"observations"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("iduna list observations parse: %w", err)
	}
	return result.Observations, nil
}

// CompleteObservation stores the analysis result for a camera observation.
func (c *IdunaClient) CompleteObservation(ctx context.Context, id int64, analysis string, appleID int64) error {
	if err := c.authenticate(ctx); err != nil {
		return fmt.Errorf("iduna authenticate: %w", err)
	}
	c.mu.Lock()
	tok := c.token
	c.mu.Unlock()

	body, _ := json.Marshal(map[string]any{
		"analysis": analysis,
		"apple_id": appleID,
		"status":   "done",
	})
	req, err := http.NewRequestWithContext(ctx, "PATCH",
		fmt.Sprintf("%s/api/v1/intelligence/observations/%d", c.baseURL, id),
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("iduna complete observation: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4*1024))
		return fmt.Errorf("iduna complete observation status %d: %s", resp.StatusCode, raw)
	}
	return nil
}

// FetchAppleContext returns a compact summary of the n most recent Apples,
// PostCheckin records a check-in for the monitor identified by slug.
// This is the heartbeat call that resets the monitor's overdue timer.
// Non-fatal: errors are logged but not returned.
func (c *IdunaClient) PostCheckin(ctx context.Context, slug string) {
	if c == nil {
		return
	}
	if err := c.authenticate(ctx); err != nil {
		return
	}
	url := strings.TrimRight(c.baseURL, "/") + "/api/v1/monitors/checkin/" + slug
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		raw, _ := io.ReadAll(resp.Body)
		log.Printf("checkin %s: IDUNA %d: %s", slug, resp.StatusCode, string(raw))
	}
}

// EnsureCronMonitor creates the emily-prime-cron monitor in IDUNA if it does not
// already exist. Timeout 1800s (30m) — cron fires every 15m so missing two ticks
// is the alert threshold. Safe to call on every startup.
func (c *IdunaClient) EnsureCronMonitor(ctx context.Context) {
	if c == nil {
		return
	}
	if err := c.authenticate(ctx); err != nil {
		return
	}
	body, _ := json.Marshal(map[string]any{
		"name":          "Emily Prime cron",
		"slug":          "emily-prime-cron",
		"kind":          "cron",
		"grace_seconds": 1800,
		"alert_email":   "emilyspringerton@gmail.com",
	})
	url := strings.TrimRight(c.baseURL, "/") + "/api/v1/monitors"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 201 = created; 409 = already exists (both OK).
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		raw, _ := io.ReadAll(resp.Body)
		log.Printf("EnsureCronMonitor: IDUNA %d: %s", resp.StatusCode, string(raw))
	}
}

// formatted for injection into LLM prompts as a ≤200-token context window.
// Returns an empty string on error so callers can degrade gracefully.
func (c *IdunaClient) FetchAppleContext(ctx context.Context, n int) string {
	items, err := c.ListApples(ctx, "", "", n)
	if err != nil || len(items) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("Recent system activity (latest Apples):\n")
	for _, a := range items {
		fmt.Fprintf(&sb, "- #%d [%s] %s (%s)\n",
			a.ID, a.AppleType, a.Title, a.RecordedAt.UTC().Format("01-02 15:04"))
	}
	return sb.String()
}
