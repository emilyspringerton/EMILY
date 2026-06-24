// api_health.go — S125-07: GET /health
//
// Simple health endpoint for gfdapi (PitViper webmaster mode) and other
// downstream services that poll Emily Prime to determine gear state.
// Returns:
//   {"ok":true,"service":"emily-agent","gear":"ACTIVE"}
//
// "gear" is the Emiree witch-engine gear name (idle/slow/cruise/active/drive/high-power/overload).
// Callers may map any non-"idle" gear to "ACTIVE" for simplified pass/fail checks.
package main

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gear := "UNKNOWN"
	if s.emiree != nil {
		gear = s.emiree.State.Gear.String()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"ok":      true,
		"service": "emily-agent",
		"gear":    gear,
	})
}
