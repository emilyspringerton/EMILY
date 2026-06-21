// api_posture.go — GET /api/v1/emily/posture (S49-01)
//
// Returns current EmilyOS posture so downstream services can gate on it
// without reading posture.json directly.
package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) handlePostureGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	p := readPosture()
	llmBlocked := p == PostureSiege || p == PostureExited

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"posture":     p,
		"llm_blocked": llmBlocked,
		"queried_at":  time.Now().UTC().Format(time.RFC3339),
	})
}
