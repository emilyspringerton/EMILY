package main

// dragon.go: Emily Prime as the Dragon — TRAPX city intelligence.
//
// The Dragon reads the TRAPX city state each RSI OBSERVE phase and fires city
// events in the ACT phase. When TRAPX_SERVER_URL is set, the RSI cycle gains
// a "city summary" section: districts in crisis surface alongside normal
// FatBaby observations.

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// DragonCityState mirrors trapxapi.CityState for JSON decode.
type DragonCityState struct {
	At             time.Time             `json:"at"`
	TechPressure   float64               `json:"tech_pressure"`
	TechTier       string                `json:"tech_tier"`
	CrownFired     bool                  `json:"crown_fired"`
	RogueCount     int                   `json:"rogue_count"`
	Districts      []DragonDistrictState `json:"districts"`
	RecentReceipts []string              `json:"recent_receipts"`
}

// DragonDistrictState mirrors trapxapi.DistrictState.
type DragonDistrictState struct {
	DistrictID       string  `json:"district_id"`
	ControlIntegrity float64 `json:"control_integrity"`
	IsRogue          bool    `json:"is_rogue"`
	ScarCount        int     `json:"scar_count"`
	Attention        float64 `json:"attention"`
	IsUnderAudit     bool    `json:"is_under_audit"`
	FOCount          int     `json:"fo_count"`
	TotalFlow        float64 `json:"total_flow"`
}

// dragonServerURL returns the TRAPX server URL from env, or "" if unset.
func dragonServerURL() string {
	return os.Getenv("TRAPX_SERVER_URL")
}

// DragonObserve fetches the TRAPX city state and returns a formatted summary
// string for injection into the RSI OBSERVE context window. Returns "" if
// TRAPX_SERVER_URL is not set or the fetch fails (non-fatal).
func DragonObserve() string {
	base := dragonServerURL()
	if base == "" {
		return ""
	}
	url := strings.TrimRight(base, "/") + "/api/v1/trapx/city-state"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[dragon] city-state fetch failed: %v", err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("[dragon] city-state: unexpected status %d", resp.StatusCode)
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[dragon] city-state read: %v", err)
		return ""
	}
	var cs DragonCityState
	if err := json.Unmarshal(body, &cs); err != nil {
		log.Printf("[dragon] city-state decode: %v", err)
		return ""
	}
	return dragonSummary(cs)
}

// dragonSummary formats a compact city state summary for the RSI context window.
func dragonSummary(cs DragonCityState) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== DRAGON CITY STATE (as of %s) ===\n", cs.At.Format("15:04:05")))
	sb.WriteString(fmt.Sprintf("TechPressure=%.0f (%s)  RogueDistricts=%d  CrownFired=%v\n",
		cs.TechPressure, cs.TechTier, cs.RogueCount, cs.CrownFired))

	sb.WriteString("Districts:\n")
	for _, d := range cs.Districts {
		flags := ""
		if d.IsRogue {
			flags += " [ROGUE SWARM]"
		}
		if d.IsUnderAudit {
			flags += " [AUDIT]"
		}
		crisis := "stable"
		if d.ControlIntegrity < 0.3 {
			crisis = "CRITICAL"
		} else if d.ControlIntegrity < 0.6 {
			crisis = "stressed"
		}
		sb.WriteString(fmt.Sprintf("  %-25s CI=%.2f attn=%.0f scars=%d status=%s%s\n",
			d.DistrictID, d.ControlIntegrity, d.Attention, d.ScarCount, crisis, flags))
	}

	if len(cs.RecentReceipts) > 0 {
		sb.WriteString("Recent city receipts:\n")
		for _, r := range cs.RecentReceipts {
			sb.WriteString("  " + r + "\n")
		}
	}
	sb.WriteString("=== END DRAGON CITY STATE ===\n")
	return sb.String()
}

// DragonFireEvent sends a city event to the TRAPX server. Returns "" on success
// or an error string. Used by the Emily Prime ACT phase.
func DragonFireEvent(eventType, districtID, params, actor string) error {
	base := dragonServerURL()
	if base == "" {
		return fmt.Errorf("TRAPX_SERVER_URL not set")
	}
	url := strings.TrimRight(base, "/") + "/api/v1/trapx/events"

	payload := map[string]string{
		"event_type":  eventType,
		"district_id": districtID,
		"params":      params,
		"actor":       actor,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("dragon fire event: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		rb, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("dragon fire event: status %d: %s", resp.StatusCode, string(rb))
	}
	log.Printf("[dragon] fired event type=%s district=%s", eventType, districtID)
	return nil
}

// InjectDragonObservation appends the Dragon city summary to a string slice of
// observations. The slice is the same format used by the RSI OBSERVE phase for
// FatBaby observation files. Returns the extended slice.
func InjectDragonObservation(observations []string) []string {
	summary := DragonObserve()
	if summary == "" {
		return observations
	}
	return append(observations, summary)
}
