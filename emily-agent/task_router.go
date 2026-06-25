package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// clusterHeartbeatView mirrors the JSON from GET /api/v1/agents?type=emily_cluster&active=true.
type clusterHeartbeatView struct {
	AgentID      string   `json:"agent_id"`
	ClusterID    string   `json:"cluster_id"`
	Capabilities []string `json:"capabilities"`
	LoadScore    float64  `json:"load_score"`
	LastSeen     string   `json:"last_seen"`
}

type clusterListResp struct {
	Clusters []clusterHeartbeatView `json:"clusters"`
}

// ListActiveClusters queries IDUNA for active emily_cluster agents.
// Returns nil on any error (non-fatal; caller falls back to local execution).
func (c *IdunaClient) ListActiveClusters(ctx context.Context) ([]clusterHeartbeatView, error) {
	if c == nil {
		return nil, nil
	}
	if err := c.authenticate(ctx); err != nil {
		return nil, fmt.Errorf("list clusters auth: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/api/v1/agents?type=emily_cluster&active=true", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list clusters: status %d", resp.StatusCode)
	}
	var result clusterListResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Clusters, nil
}

// hasCapability reports whether a cluster advertises cap.
func hasCapability(c clusterHeartbeatView, cap string) bool {
	for _, v := range c.Capabilities {
		if v == cap {
			return true
		}
	}
	return false
}

// routeTask decides whether to route a directed task to a remote cluster.
// Returns "" when local execution is preferred.
//
// Routing rule: if any remote cluster (clusterID != localID) has:
//   - the required capability, AND
//   - load_score < localLoad
//
// pick the lowest-load qualifying remote cluster.
// The caller then stamps routed_to_cluster in the task JSON and skips local execution.
func routeTask(clusters []clusterHeartbeatView, localID string, localLoad float64, requiredCap string) string {
	best := ""
	bestLoad := localLoad
	for _, c := range clusters {
		if c.ClusterID == localID {
			continue // skip self
		}
		if requiredCap != "" && !hasCapability(c, requiredCap) {
			continue
		}
		if c.LoadScore < bestLoad {
			bestLoad = c.LoadScore
			best = c.ClusterID
		}
	}
	return best
}

// MaybeRouteTask is called before dispatching a directed task to the local obs-watcher.
// If a better cluster exists it stamps the task file with routed_to_cluster,
// logs the routing decision, and returns true (caller must skip local dispatch).
// Returns false when the task should run locally.
func MaybeRouteTask(ctx context.Context, iduna *IdunaClient, taskBytes []byte, requiredCap string) ([]byte, bool) {
	if iduna == nil {
		return taskBytes, false
	}
	clusters, err := iduna.ListActiveClusters(ctx)
	if err != nil {
		log.Printf("task-router: list clusters: %v (routing locally)", err)
		return taskBytes, false
	}
	if len(clusters) <= 1 {
		return taskBytes, false // only self; no routing possible
	}

	localID := clusterID()
	target := routeTask(clusters, localID, 0.0, requiredCap)
	if target == "" {
		return taskBytes, false
	}

	// Stamp routed_to_cluster into the raw JSON.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(taskBytes, &raw); err != nil {
		return taskBytes, false
	}
	targetJSON, _ := json.Marshal(target)
	raw["routed_to_cluster"] = targetJSON
	raw["routed_from_cluster"], _ = json.Marshal(localID)

	out, err := marshalIndentCompact(raw)
	if err != nil {
		return taskBytes, false
	}
	log.Printf("task-router: routing task to cluster %s (has cap=%q, lower load)", target, requiredCap)
	return out, true
}

// marshalIndentCompact marshals map with sorted keys to pretty JSON.
func marshalIndentCompact(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}
