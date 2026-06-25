package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// heartbeatPayload matches POST /api/v1/agents/heartbeat on IDUNA.
type heartbeatPayload struct {
	AgentID      string   `json:"agent_id"`
	ClusterID    string   `json:"cluster_id"`
	Capabilities []string `json:"capabilities"`
	LoadScore    float64  `json:"load_score"`
}

// SendHeartbeat POSTs a cluster heartbeat to IDUNA. Non-fatal: logs on error.
// clusterID identifies this installation (defaults to hostname).
// capabilities are auto-detected from environment.
func (c *IdunaClient) SendHeartbeat(ctx context.Context, clusterID string, loadScore float64) error {
	if c == nil {
		return nil
	}
	if err := c.authenticate(ctx); err != nil {
		return fmt.Errorf("heartbeat auth: %w", err)
	}

	agentID := c.agentName // IDUNA identifies agents by name in the heartbeat table
	caps := detectCapabilities()
	payload := heartbeatPayload{
		AgentID:      agentID,
		ClusterID:    clusterID,
		Capabilities: caps,
		LoadScore:    loadScore,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/api/v1/agents/heartbeat", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("heartbeat POST: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// detectCapabilities returns a list of capability tags based on environment and hardware.
func detectCapabilities() []string {
	caps := []string{"emily_prime"}

	// GPU detection: check for CUDA device or NVIDIA env hints.
	if os.Getenv("CUDA_VISIBLE_DEVICES") != "" || os.Getenv("NVIDIA_VISIBLE_DEVICES") != "" {
		caps = append(caps, "gpu")
	}
	// Check if secwatch store is configured.
	if os.Getenv("FATBABY_ROOT") != "" || os.Getenv("FATBABY_SIGNAL_API_URL") != "" {
		caps = append(caps, "secwatch", "prwatch", "earnings_calendar")
	}
	// MJOLNIR push configured.
	if os.Getenv("FCM_PROJECT_ID") != "" {
		caps = append(caps, "fcm_push")
	}
	// Arch tag.
	caps = append(caps, runtime.GOARCH)

	return caps
}

// clusterID returns a stable identifier for this cluster.
// Prefers EMILY_CLUSTER_ID env; falls back to hostname.
func clusterID() string {
	if id := os.Getenv("EMILY_CLUSTER_ID"); id != "" {
		return id
	}
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "unknown"
	}
	return strings.ToLower(h)
}

// startHeartbeatLoop fires an initial heartbeat and then one every 60 s.
// Runs as a goroutine; stops when ctx is cancelled.
func startHeartbeatLoop(ctx context.Context, iduna *IdunaClient, loadFn func() float64) {
	if iduna == nil {
		return
	}
	id := clusterID()
	send := func() {
		hbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := iduna.SendHeartbeat(hbCtx, id, loadFn()); err != nil {
			log.Printf("cluster heartbeat: %v", err)
		} else {
			log.Printf("cluster heartbeat: sent cluster_id=%s", id)
		}
	}
	send() // immediate on startup
	go func() {
		t := time.NewTicker(60 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				send()
			}
		}
	}()
}
