package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func makeTestCityState(rogueCount int, techPressure float64) DragonCityState {
	return DragonCityState{
		At:           time.Date(2026, 6, 24, 12, 0, 0, 0, time.UTC),
		TechPressure: techPressure,
		TechTier:     "LeashFrays",
		CrownFired:   false,
		RogueCount:   rogueCount,
		Districts: []DragonDistrictState{
			{DistrictID: "district-residential", ControlIntegrity: 0.9, IsRogue: false, ScarCount: 0, Attention: 120},
			{DistrictID: "district-commercial", ControlIntegrity: 0.2, IsRogue: rogueCount > 0, ScarCount: 2, Attention: 750, IsUnderAudit: true},
		},
		RecentReceipts: []string{"[12:00:00] fo-001 — CLAIMED by player1"},
	}
}

func TestDragonSummaryContainsHeader(t *testing.T) {
	cs := makeTestCityState(0, 0)
	s := dragonSummary(cs)
	if !strings.Contains(s, "DRAGON CITY STATE") {
		t.Error("summary should contain DRAGON CITY STATE header")
	}
}

func TestDragonSummaryContainsDistricts(t *testing.T) {
	cs := makeTestCityState(0, 0)
	s := dragonSummary(cs)
	if !strings.Contains(s, "district-residential") {
		t.Error("summary should list district-residential")
	}
}

func TestDragonSummaryCriticalStatus(t *testing.T) {
	cs := makeTestCityState(1, 300)
	s := dragonSummary(cs)
	if !strings.Contains(s, "CRITICAL") {
		t.Error("district with CI=0.2 should be flagged CRITICAL")
	}
}

func TestDragonSummaryRogueFlag(t *testing.T) {
	cs := makeTestCityState(1, 0)
	s := dragonSummary(cs)
	if !strings.Contains(s, "ROGUE SWARM") {
		t.Error("rogue district should be flagged [ROGUE SWARM]")
	}
}

func TestDragonSummaryAuditFlag(t *testing.T) {
	cs := makeTestCityState(0, 0)
	cs.Districts[1].IsUnderAudit = true
	s := dragonSummary(cs)
	if !strings.Contains(s, "AUDIT") {
		t.Error("district under audit should be flagged [AUDIT]")
	}
}

func TestDragonSummaryTechPressure(t *testing.T) {
	cs := makeTestCityState(0, 350)
	s := dragonSummary(cs)
	if !strings.Contains(s, "350") {
		t.Error("summary should include tech pressure value")
	}
}

func TestDragonSummaryRecentReceipts(t *testing.T) {
	cs := makeTestCityState(0, 0)
	s := dragonSummary(cs)
	if !strings.Contains(s, "CLAIMED") {
		t.Error("summary should include recent receipts")
	}
}

func TestDragonObserveNoURL(t *testing.T) {
	t.Setenv("TRAPX_SERVER_URL", "")
	result := DragonObserve()
	if result != "" {
		t.Error("DragonObserve should return empty string when TRAPX_SERVER_URL unset")
	}
}

func TestDragonObserveWithServer(t *testing.T) {
	cs := makeTestCityState(0, 100)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/trapx/city-state" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cs)
	}))
	defer srv.Close()

	t.Setenv("TRAPX_SERVER_URL", srv.URL)
	result := DragonObserve()
	if result == "" {
		t.Error("DragonObserve should return summary when server responds")
	}
	if !strings.Contains(result, "DRAGON CITY STATE") {
		t.Errorf("result should contain city state header, got: %q", result)
	}
}

func TestDragonFireEventNoURL(t *testing.T) {
	t.Setenv("TRAPX_SERVER_URL", "")
	err := DragonFireEvent("rogue_swarm", "district-residential", "", "test")
	if err == nil {
		t.Error("should return error when TRAPX_SERVER_URL not set")
	}
}

func TestDragonFireEventSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "message": "done"})
	}))
	defer srv.Close()

	t.Setenv("TRAPX_SERVER_URL", srv.URL)
	if err := DragonFireEvent("rogue_swarm", "district-residential", "test", "emily-dragon"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInjectDragonObservationNoURL(t *testing.T) {
	t.Setenv("TRAPX_SERVER_URL", "")
	obs := []string{"existing observation"}
	result := InjectDragonObservation(obs)
	if len(result) != 1 {
		t.Errorf("should not add to observations when no URL set, got len=%d", len(result))
	}
}

func TestInjectDragonObservationWithURL(t *testing.T) {
	cs := makeTestCityState(0, 0)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cs)
	}))
	defer srv.Close()

	t.Setenv("TRAPX_SERVER_URL", srv.URL)
	obs := []string{"existing observation"}
	result := InjectDragonObservation(obs)
	if len(result) != 2 {
		t.Errorf("should add Dragon observation, got len=%d", len(result))
	}
}

// ── dragonDecide ──────────────────────────────────────────────────────────────

func TestDragonDecideStableCity(t *testing.T) {
	cs := DragonCityState{
		TechPressure: 150,
		Districts: []DragonDistrictState{
			{DistrictID: "district-residential", ControlIntegrity: 0.9, Attention: 200},
		},
	}
	decisions := dragonDecide(cs, 3)
	if len(decisions) != 0 {
		t.Errorf("stable city should produce no decisions, got %d", len(decisions))
	}
}

func TestDragonDecideRogueSwarmTriggered(t *testing.T) {
	cs := DragonCityState{
		Districts: []DragonDistrictState{
			{DistrictID: "district-commercial", ControlIntegrity: 0.15, IsRogue: false},
		},
	}
	decisions := dragonDecide(cs, 10)
	if len(decisions) == 0 {
		t.Fatal("expected rogue_swarm decision for CI=0.15")
	}
	if decisions[0].EventType != "rogue_swarm" {
		t.Errorf("expected rogue_swarm, got %q", decisions[0].EventType)
	}
	if decisions[0].DistrictID != "district-commercial" {
		t.Errorf("expected district-commercial, got %q", decisions[0].DistrictID)
	}
}

func TestDragonDecideSkipsAlreadyRogue(t *testing.T) {
	cs := DragonCityState{
		Districts: []DragonDistrictState{
			{DistrictID: "district-commercial", ControlIntegrity: 0.1, IsRogue: true},
		},
	}
	decisions := dragonDecide(cs, 10)
	for _, d := range decisions {
		if d.EventType == "rogue_swarm" {
			t.Error("should not fire rogue_swarm on already-rogue district")
		}
	}
}

func TestDragonDecideMediaSpike(t *testing.T) {
	cs := DragonCityState{
		Districts: []DragonDistrictState{
			{DistrictID: "district-residential", ControlIntegrity: 0.8, Attention: 850, IsUnderAudit: false},
		},
	}
	decisions := dragonDecide(cs, 10)
	var found bool
	for _, d := range decisions {
		if d.EventType == "media_spike" {
			found = true
		}
	}
	if !found {
		t.Error("expected media_spike for attention=850")
	}
}

func TestDragonDecidePressureSpike(t *testing.T) {
	cs := DragonCityState{
		TechPressure: 50,
		Districts:    []DragonDistrictState{},
	}
	decisions := dragonDecide(cs, 10)
	var found bool
	for _, d := range decisions {
		if d.EventType == "pressure_spike" {
			found = true
		}
	}
	if !found {
		t.Error("expected pressure_spike for low tech pressure after cycle 5")
	}
}

func TestDragonDecideNoPressureSpikeEarlyCycles(t *testing.T) {
	cs := DragonCityState{
		TechPressure: 50,
		Districts:    []DragonDistrictState{},
	}
	decisions := dragonDecide(cs, 3) // cycle < 5, no pressure spike
	for _, d := range decisions {
		if d.EventType == "pressure_spike" {
			t.Error("should not pressure_spike in early cycles")
		}
	}
}

// ── dragonACT integration ──────────────────────────────────────────────────────

func TestDragonACTNoURL(t *testing.T) {
	t.Setenv("TRAPX_SERVER_URL", "")
	// Should be a no-op; no panic.
	dragonACT(t.Context(), nil, 1, "/tmp/dragon-test-stream.ndjson")
}

// ── TRAPX event bridge ─────────────────────────────────────────────────────────

func TestDragonBroadcastWorldEventNoURL(t *testing.T) {
	t.Setenv("TRAPX_MUD_API_URL", "")
	// Should be a no-op; no panic.
	DragonBroadcastWorldEvent(t.Context(), "emily_cast", "test", "")
}

func TestDragonBroadcastWorldEventPosts(t *testing.T) {
	var received []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/world-events" || r.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()
	t.Setenv("TRAPX_MUD_API_URL", srv.URL)

	DragonBroadcastWorldEvent(t.Context(), "faction_war_start", "war begins", "district-residential")

	if received == nil {
		t.Fatal("no request received by mock MUD server")
	}
	var body map[string]string
	if err := json.Unmarshal(received, &body); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if body["type"] != "faction_war_start" {
		t.Errorf("type: want faction_war_start, got %s", body["type"])
	}
	if body["district"] != "district-residential" {
		t.Errorf("district: want district-residential, got %s", body["district"])
	}
}

func TestTrapXPollApplesNoNew(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"apples": []map[string]any{
				{"id": "apple-001", "title": "S126-01 done"},
			},
		})
	}))
	defer srv.Close()
	t.Setenv("TRAPX_MUD_API_URL", "http://127.0.0.1:19999") // not listening — but poll should not reach broadcast for already-seen

	latest := trapxPollApples(t.Context(), srv.URL, "apple-001")
	// already-seen apple-001 → no new forward; latest stays the same
	if latest != "apple-001" {
		t.Errorf("want apple-001, got %s", latest)
	}
}
