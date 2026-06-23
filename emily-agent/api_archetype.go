package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

// archetypeEngineURL returns the configured archetype-engine address.
func archetypeEngineURL() string {
	if u := os.Getenv("ARCHETYPE_ENGINE_URL"); u != "" {
		return u
	}
	return "http://localhost:8090"
}

// handleArchetypeStatus proxies GET /api/v1/emily/archetype/status to the archetype engine.
// Returns the engine status JSON with an added `reachable` field.
func (s *Server) handleArchetypeStatus(w http.ResponseWriter, r *http.Request) {
	engineURL := archetypeEngineURL()
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(engineURL + "/status")

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"reachable":   false,
			"engine_url":  engineURL,
			"error":       err.Error(),
			"queried_at":  time.Now().UTC().Format(time.RFC3339),
		})
		return
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
	if err != nil || resp.StatusCode != http.StatusOK {
		json.NewEncoder(w).Encode(map[string]any{
			"reachable":  false,
			"engine_url": engineURL,
			"status":     resp.StatusCode,
		})
		return
	}

	// Merge engine status with reachable=true.
	var status map[string]any
	if err := json.Unmarshal(raw, &status); err != nil {
		w.Write(raw) // pass through verbatim if unmarshal fails
		return
	}
	status["reachable"] = true
	status["engine_url"] = engineURL
	json.NewEncoder(w).Encode(status)
}

// handleArchetypeSpirits proxies GET /api/v1/emily/archetype/spirits to the archetype engine.
func (s *Server) handleArchetypeSpirits(w http.ResponseWriter, r *http.Request) {
	engineURL := archetypeEngineURL()
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(engineURL + "/spirits")
	if err != nil {
		http.Error(w, "archetype engine unreachable: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
