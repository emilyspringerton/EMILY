package main

// dragon.go: Emily Prime as the Dragon — TRAPX city intelligence.
//
// The Dragon reads the TRAPX city state each RSI OBSERVE phase and fires city
// events in the ACT phase. When TRAPX_SERVER_URL is set, the RSI cycle gains
// a "city summary" section: districts in crisis surface alongside normal
// FatBaby observations.

import (
	"context"
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

// dragonDecision holds a Dragon-triggered city event.
type dragonDecision struct {
	EventType  string
	DistrictID string
	Reason     string
}

// dragonDecide reads city state and returns any escalation events Emily Prime
// should fire as the Dragon. Rules:
//
//   - Any district CI < 0.25 → rogue_swarm if not already rogue
//   - Any district attention > 800 → media_spike to sustain heat
//   - Tech pressure < 100 after cycle N>5 → pressure_spike to remind the city
//
// The Dragon is opportunistic: if the city is stable, she is silent.
func dragonDecide(cs DragonCityState, cycleNum int) []dragonDecision {
	var decisions []dragonDecision
	for _, d := range cs.Districts {
		if d.ControlIntegrity < 0.25 && !d.IsRogue {
			decisions = append(decisions, dragonDecision{
				EventType:  "rogue_swarm",
				DistrictID: d.DistrictID,
				Reason:     fmt.Sprintf("CI=%.2f below rogue threshold in %s", d.ControlIntegrity, d.DistrictID),
			})
		}
		if d.Attention > 800 && !d.IsUnderAudit {
			decisions = append(decisions, dragonDecision{
				EventType:  "media_spike",
				DistrictID: d.DistrictID,
				Reason:     fmt.Sprintf("attention=%.0f sustained heat in %s", d.Attention, d.DistrictID),
			})
		}
	}
	if cycleNum > 5 && cs.TechPressure < 100 {
		decisions = append(decisions, dragonDecision{
			EventType:  "pressure_spike",
			DistrictID: "",
			Reason:     fmt.Sprintf("tech pressure=%.0f — Dragon stirs the city", cs.TechPressure),
		})
	}
	return decisions
}

// dragonACT is called at the end of each RSI cycle. It reads the city state,
// decides what to escalate, fires events to the TRAPX server, and optionally
// posts Dragon Apples to IDUNA.
func dragonACT(ctx context.Context, iduna *IdunaClient, cycleNum int, streamPath string) {
	if dragonServerURL() == "" {
		return
	}
	cs, err := dragonFetchState()
	if err != nil {
		log.Printf("[dragon act] fetch failed: %v", err)
		return
	}
	decisions := dragonDecide(cs, cycleNum)
	for _, d := range decisions {
		if err := DragonFireEvent(d.EventType, d.DistrictID, d.Reason, "emily-dragon"); err != nil {
			log.Printf("[dragon act] fire %s failed: %v", d.EventType, err)
			continue
		}
		log.Printf("[dragon act] fired %s in %q: %s", d.EventType, d.DistrictID, d.Reason)
		appendStream(streamPath, StreamEntry{
			Type:  "dragon_act",
			Cycle: cycleNum,
			Msg:   fmt.Sprintf("fired %s district=%s reason=%s", d.EventType, d.DistrictID, d.Reason),
		})
		// Post Dragon Apple (S121-04 hook: file every Dragon-triggered event).
		if iduna != nil {
			appleCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
			defer cancel()
			_, _ = iduna.PostApple(appleCtx, ApplePayload{
				AppleType:  "observation",
				SourceRepo: "TRAPX",
				Title:      fmt.Sprintf("Dragon: %s in %s", d.EventType, d.DistrictID),
				Body:       fmt.Sprintf("cycle=%d event=%s district=%s reason=%s", cycleNum, d.EventType, d.DistrictID, d.Reason),
			})
		}
	}
	if len(decisions) == 0 {
		log.Printf("[dragon act] city stable — no events fired (cycle %d)", cycleNum)
	}
}

// dragonFetchState fetches and decodes the city state from the TRAPX server.
func dragonFetchState() (DragonCityState, error) {
	base := dragonServerURL()
	url := strings.TrimRight(base, "/") + "/api/v1/trapx/city-state"
	resp, err := http.Get(url)
	if err != nil {
		return DragonCityState{}, fmt.Errorf("get city state: %w", err)
	}
	defer resp.Body.Close()
	var cs DragonCityState
	if err := json.NewDecoder(resp.Body).Decode(&cs); err != nil {
		return DragonCityState{}, fmt.Errorf("decode: %w", err)
	}
	return cs, nil
}
