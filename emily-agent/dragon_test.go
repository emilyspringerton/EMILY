package main

import (
	"encoding/json"
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
