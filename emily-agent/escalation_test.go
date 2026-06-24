package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestImprovementTaskEscalationFields(t *testing.T) {
	now := time.Now().UTC()
	task := ImprovementTask{
		ID:          "S126-07",
		Description: "test task",
		Status:      "running",
		StuckCycles: 3,
		EscalatedAt: &now,
	}
	if task.StuckCycles != 3 {
		t.Errorf("want StuckCycles=3, got %d", task.StuckCycles)
	}
	if task.EscalatedAt == nil {
		t.Error("EscalatedAt should not be nil")
	}
}

func TestImprovementTaskEscalationJSONRoundtrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	orig := ImprovementTask{
		ID:          "task-001",
		Status:      "running",
		StuckCycles: 4,
		EscalatedAt: &now,
	}
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatal(err)
	}
	var rt ImprovementTask
	if err := json.Unmarshal(data, &rt); err != nil {
		t.Fatal(err)
	}
	if rt.StuckCycles != 4 {
		t.Errorf("StuckCycles: want 4, got %d", rt.StuckCycles)
	}
	if rt.EscalatedAt == nil {
		t.Error("EscalatedAt not round-tripped")
	}
}

func TestEscalateTaskWithSonnetNoKey(t *testing.T) {
	task := &ImprovementTask{ID: "t1", Description: "do something", StuckCycles: 3}
	_, err := escalateTaskWithSonnet(t.Context(), task, "")
	if err == nil {
		t.Error("expected error when API key is empty")
	}
}

func TestEscalateTaskWithSonnetSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"content": []map[string]any{
				{"type": "text", "text": "Revised: focus only on the Go struct change."},
			},
		})
	}))
	defer srv.Close()

	// Patch the API URL via a round-tripper swap isn't easy; instead we test via
	// a mock server by directly calling the function with a patched transport.
	// Since escalateTaskWithSonnet uses http.DefaultClient, we verify the mock
	// server response structure instead.
	var reqBody map[string]any
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&reqBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"content": []map[string]any{
				{"type": "text", "text": "Revised approach: narrow scope to single file."},
			},
		})
	}))
	defer srv2.Close()

	// We can't easily redirect http.DefaultClient without affecting other tests.
	// Verify the function builds correct request payload structure by inspecting
	// what it would send. Instead, test the field placement and no-error path
	// with a stubbed transport.
	origTransport := http.DefaultTransport
	http.DefaultTransport = &mockTransport{
		handler: func(r *http.Request) (*http.Response, error) {
			// Check that the request targets anthropic.
			if !strings.Contains(r.URL.String(), "anthropic") {
				return nil, nil
			}
			rec := httptest.NewRecorder()
			json.NewEncoder(rec).Encode(map[string]any{
				"content": []map[string]any{
					{"type": "text", "text": "Revised: single-step approach."},
				},
			})
			return rec.Result(), nil
		},
	}
	defer func() { http.DefaultTransport = origTransport }()

	task := &ImprovementTask{ID: "t2", Description: "do something hard", StuckCycles: 4}
	revised, err := escalateTaskWithSonnet(t.Context(), task, "fake-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if revised == "" {
		t.Error("expected non-empty revised description")
	}
}

// mockTransport is an http.RoundTripper backed by a handler function.
type mockTransport struct {
	handler func(*http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return m.handler(r)
}

func TestEscalationThresholdConstant(t *testing.T) {
	if escalationThreshold != 3 {
		t.Errorf("escalationThreshold should be 3, got %d", escalationThreshold)
	}
}
