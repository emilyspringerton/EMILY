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

	"emily-agent/pkg/archetypes"
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
// decides what to escalate, fires events to the TRAPX server, routes decisions
// through the FIELD archetype engine, and posts Dragon Apples to IDUNA.
func dragonACT(ctx context.Context, iduna *IdunaClient, cycleNum int, streamPath string) {
	dragonACTWithField(ctx, iduna, nil, cycleNum, streamPath)
}

// dragonACTWithField is the full Dragon ACT implementation with optional FIELD routing.
func dragonACTWithField(ctx context.Context, iduna *IdunaClient, field *archetypes.Field, cycleNum int, streamPath string) {
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

		// Archetype routing (S121-05): augment the Dragon decision with FIELD resonance.
		archetypeTag := DragonArchetypeAugment(ctx, field, d, cs)
		streamMsg := fmt.Sprintf("fired %s district=%s reason=%s", d.EventType, d.DistrictID, d.Reason)
		if archetypeTag != "" {
			streamMsg += " " + archetypeTag
		}
		appendStream(streamPath, StreamEntry{
			Type:  "dragon_act",
			Cycle: cycleNum,
			Msg:   streamMsg,
		})

		// Post Dragon Apple with archetype resonance metadata.
		if iduna != nil {
			appleCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
			defer cancel()
			body := fmt.Sprintf("cycle=%d event=%s district=%s reason=%s",
				cycleNum, d.EventType, d.DistrictID, d.Reason)
			if archetypeTag != "" {
				body += "\n" + archetypeTag
			}
			_, _ = iduna.PostApple(appleCtx, ApplePayload{
				AppleType:  "observation",
				SourceRepo: "TRAPX",
				Title:      fmt.Sprintf("Dragon: %s in %s", d.EventType, d.DistrictID),
				Body:       body,
			})
		}
	}
	if len(decisions) == 0 {
		log.Printf("[dragon act] city stable — no events fired (cycle %d)", cycleNum)
	}
}

// DragonArchetypeAugment runs the FIELD archetype engine for a Dragon city-event
// decision and returns an archetype resonance string to attach to Dragon Apples.
// Default spirit stack: Raum#40 (city/information) + Amon#7 (insight) + Gaap#33 (foresight).
// Returns "" if field is nil or invocation fails — Dragon fires regardless.
func DragonArchetypeAugment(ctx context.Context, field *archetypes.Field, decision dragonDecision, cs DragonCityState) string {
	if field == nil {
		return ""
	}
	intent := fmt.Sprintf("city escalation: %s in %s", decision.EventType, decision.DistrictID)
	ctxText := fmt.Sprintf("tech_pressure=%.0f rogue_count=%d reason=%s", cs.TechPressure, cs.RogueCount, decision.Reason)

	augCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	res, err := field.Invoke(augCtx, intent, ctxText, false)
	if err != nil {
		log.Printf("[dragon field] invoke failed: %v", err)
		return ""
	}
	return fmt.Sprintf("archetype_corridor=%s spirit_stack=%s Δφ=%.0f°",
		res.ResonanceState, spiritStack(res.ActiveSpirits), res.PhaseDeltaDeg)
}

// ── EMILY ↔ TRAPX event bridge ─────────────────────────────────────────────

// trapxMUDAPIURL returns the MUD world-event API URL, or "" if unset.
func trapxMUDAPIURL() string { return os.Getenv("TRAPX_MUD_API_URL") }

// DragonBroadcastWorldEvent POSTs an event to the MUD world-event API
// (TRAPX_MUD_API_URL/api/world-events). Non-fatal: logs on error.
func DragonBroadcastWorldEvent(ctx context.Context, eventType, message, district string) {
	base := trapxMUDAPIURL()
	if base == "" {
		return
	}
	payload := map[string]string{
		"type":     eventType,
		"message":  message,
		"district": district,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST",
		strings.TrimRight(base, "/")+"/api/world-events",
		strings.NewReader(string(body)))
	if err != nil {
		log.Printf("[dragon-bridge] build request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[dragon-bridge] POST world-event: %v", err)
		return
	}
	resp.Body.Close()
	log.Printf("[dragon-bridge] world-event posted: type=%s district=%s status=%d", eventType, district, resp.StatusCode)
}

// TrapXAppleWatcher polls IDUNA /api/v1/apples for new GoblinFoxDragon completion
// Apples and forwards them to the MUD world-event queue. Runs as a goroutine.
// Exits when ctx is cancelled.
func TrapXAppleWatcher(ctx context.Context, idunaBaseURL string, intervalSec int) {
	if trapxMUDAPIURL() == "" {
		return
	}
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()
	var lastSeen string

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			lastSeen = trapxPollApples(ctx, idunaBaseURL, lastSeen)
		}
	}
}

// trapxPollApples fetches completion Apples from IDUNA and forwards new ones to the MUD.
// Returns the ID of the most recent Apple seen (for pagination).
func trapxPollApples(ctx context.Context, idunaBaseURL, afterID string) string {
	url := strings.TrimRight(idunaBaseURL, "/") + "/api/v1/apples?source_repo=GoblinFoxDragon&type=completion&limit=10"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("[trapx-watcher] build request: %v", err)
		return afterID
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[trapx-watcher] poll: %v", err)
		return afterID
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return afterID
	}

	var result struct {
		Apples []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"apples"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return afterID
	}

	newLatest := afterID
	for _, a := range result.Apples {
		if a.ID == afterID {
			break // already seen
		}
		if newLatest == afterID || a.ID > newLatest {
			newLatest = a.ID
		}
		pollCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		DragonBroadcastWorldEvent(pollCtx,
			"emily_cast",
			"Emily dispatch received. Stay alert. ("+a.Title+")",
			"")
		cancel()
	}
	return newLatest
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
