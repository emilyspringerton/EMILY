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
	"net/http"
	"sync"
	"time"
)

// IdunaClient authenticates as an agent and submits Apples to IDUNA.
type IdunaClient struct {
	baseURL     string
	agentName   string
	agentSecret string
	httpClient  *http.Client

	mu        sync.Mutex
	token     string
	tokenExp  time.Time
}

// NewIdunaClientFromEnv reads IDUNA_BASE_URL, IDUNA_AGENT_NAME, IDUNA_AGENT_SECRET
// and returns a configured client. Returns nil if any required var is absent.
func NewIdunaClientFromEnv() *IdunaClient {
	base := envOr("IDUNA_BASE_URL", "")
	name := envOr("IDUNA_AGENT_NAME", "")
	secret := envOr("IDUNA_AGENT_SECRET", "")
	if base == "" || name == "" || secret == "" {
		return nil
	}
	return &IdunaClient{
		baseURL:     base,
		agentName:   name,
		agentSecret: secret,
		httpClient:  &http.Client{Timeout: 15 * time.Second},
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
