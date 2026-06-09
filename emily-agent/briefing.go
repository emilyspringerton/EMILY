// emily-agent/briefing.go
// Daily morning briefing push notification to MJOLNIR.
//
// At 09:00 UTC each day (±30 min window) Emily Prime queries IDUNA for Apples
// from the last 24 h, summarises them by type, and fires an FCM push to the
// mjolnir-emily device with a mjolnir://feed deep link. A sentinel file in the
// state directory prevents duplicate pushes within the same calendar day.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const briefingSentinelFile = "last-briefing-date.txt"

// briefingDue returns true if the current UTC time is between 08:30 and 09:30
// and today's briefing has not yet been sent.
func briefingDue(stateDir string) bool {
	now := time.Now().UTC()
	h, m := now.Hour(), now.Minute()
	windowStart := h*60+m >= 8*60+30
	windowEnd := h*60+m < 9*60+30
	if !windowStart || !windowEnd {
		return false
	}
	today := now.Format("2006-01-02")
	path := filepath.Join(stateDir, briefingSentinelFile)
	data, err := os.ReadFile(path)
	if err == nil && strings.TrimSpace(string(data)) == today {
		return false // already sent today
	}
	return true
}

// markBriefingSent records that today's briefing has been sent.
func markBriefingSent(stateDir string) {
	today := time.Now().UTC().Format("2006-01-02")
	path := filepath.Join(stateDir, briefingSentinelFile)
	_ = os.WriteFile(path, []byte(today), 0o644)
}

// runMorningBriefing fetches the last 24 h of Apples from IDUNA, summarises
// them by type, and fires an FCM push to the MJOLNIR app. It is called from the
// autonomous cycle when briefingDue() is true. Non-fatal: errors are logged.
func runMorningBriefing(ctx context.Context, iduna *IdunaClient, push PushFunc, stateDir string) {
	if iduna == nil || push == nil {
		return
	}

	apples, err := iduna.ListApples(ctx, "", "", 200)
	if err != nil {
		log.Printf("briefing: list apples failed: %v", err)
		return
	}

	cutoff := time.Now().UTC().Add(-24 * time.Hour)
	var recent []AppleListItem
	for _, a := range apples {
		if a.RecordedAt.After(cutoff) {
			recent = append(recent, a)
		}
	}

	title, body := buildBriefingMessage(recent)
	log.Printf("briefing: firing morning push — %d apples in 24h", len(recent))

	go push(title, body, map[string]string{
		"deep_link":      "mjolnir://feed",
		"briefing_date":  time.Now().UTC().Format("2006-01-02"),
		"apple_count":    fmt.Sprintf("%d", len(recent)),
	})

	markBriefingSent(stateDir)
}

// buildBriefingMessage turns a list of recent Apples into push title + body.
func buildBriefingMessage(apples []AppleListItem) (title, body string) {
	if len(apples) == 0 {
		return "Good morning — Emily Prime",
			"No new activity in the last 24 hours. All systems quiet."
	}

	// Group by type
	byType := map[string][]AppleListItem{}
	for _, a := range apples {
		byType[a.AppleType] = append(byType[a.AppleType], a)
	}

	// Sort types by count desc for headline
	type typeCount struct {
		t string
		n int
	}
	var counts []typeCount
	for t, items := range byType {
		counts = append(counts, typeCount{t, len(items)})
	}
	sort.Slice(counts, func(i, j int) bool { return counts[i].n > counts[j].n })

	// Title: concise summary
	parts := make([]string, 0, len(counts))
	for _, tc := range counts {
		parts = append(parts, fmt.Sprintf("%d %s", tc.n, tc.t))
	}
	title = fmt.Sprintf("Morning briefing — %d apples: %s", len(apples), strings.Join(parts, ", "))
	if len(title) > 90 {
		title = fmt.Sprintf("Morning briefing — %d apples (%d types)", len(apples), len(counts))
	}

	// Body: one line per type, capped at top 5 recent escalations/completions
	var sb strings.Builder
	for _, tc := range counts {
		label := typeLabel(tc.t)
		sb.WriteString(fmt.Sprintf("%s: %d  ", label, tc.n))
	}
	sb.WriteString("\n")

	// Show up to 3 most recent escalations or completions
	var notable []AppleListItem
	for _, a := range apples {
		if a.AppleType == "escalation" || a.AppleType == "completion" || a.AppleType == "improvement" {
			notable = append(notable, a)
		}
	}
	sort.Slice(notable, func(i, j int) bool {
		return notable[i].RecordedAt.After(notable[j].RecordedAt)
	})
	if len(notable) > 3 {
		notable = notable[:3]
	}
	for _, a := range notable {
		title := a.Title
		if len(title) > 60 {
			title = title[:57] + "..."
		}
		sb.WriteString(fmt.Sprintf("• %s\n", title))
	}

	return title, strings.TrimSpace(sb.String())
}

func typeLabel(t string) string {
	switch t {
	case "completion":
		return "✓"
	case "escalation":
		return "!"
	case "improvement":
		return "↑"
	case "signal_observation":
		return "~"
	case "rsi_iteration":
		return "⟳"
	case "status":
		return "·"
	default:
		return t
	}
}
