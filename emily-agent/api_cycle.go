// api_cycle.go — stable /api/v1/emily/ endpoints for external RSI orchestration.
//
// POST /api/v1/emily/cycle       — trigger one autonomous RSI cycle async
// GET  /api/v1/emily/state       — return current CycleState as JSON
// POST /api/v1/emily/roadmap/items — add a roadmap item to the RSI queue

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// handleCycleTrigger fires one autonomous RSI cycle in a goroutine.
// Returns immediately; use GET /api/v1/emily/state to observe progress.
func (s *Server) handleCycleTrigger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.cycleMu.TryLock() {
		http.Error(w, `{"error":"cycle already running"}`, http.StatusConflict)
		return
	}

	go func() {
		defer s.cycleMu.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), s.cycle.cfg.CycleDuration+30*time.Second)
		defer cancel()
		_ = ctx
		if err := s.cycle.RunOnce(); err != nil {
			// RunOnce logs its own errors; nothing more to do here.
			return
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"cycle_triggered"}`)
}

// handleCycleState returns the persisted CycleState as JSON.
func (s *Server) handleCycleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state, err := s.cycle.loadState()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(state)
}

// handleRoadmapItems supports POST to add a new RoadmapItem to the RSI queue.
// Body: { "id":"...", "name":"...", "description":"...", "priority":5, "criteria":[] }
func (s *Server) handleRoadmapItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var item RoadmapItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
		return
	}
	if item.ID == "" || item.Name == "" {
		http.Error(w, `{"error":"id and name are required"}`, http.StatusBadRequest)
		return
	}
	item.Status = "queued"
	if item.Priority == 0 {
		item.Priority = 50
	}

	if err := s.cycle.AddRoadmapItem(item); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"status":"queued","id":%q,"name":%q}`+"\n", item.ID, item.Name)
}
