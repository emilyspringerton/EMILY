// emily-agent/briefing.go
// Daily morning briefing push notification to MJOLNIR.
//
// At 09:00 UTC each day (±30 min window) Emily Prime queries IDUNA for Apples
// from the last 24 h, summarises them by type, and fires an FCM push to the
// mjolnir-emily device with a mjolnir://feed deep link. A sentinel file in the
// state directory prevents duplicate pushes within the same calendar day.
//
// When FATBABY_SIGNAL_API_URL and FATBABY_SIGNAL_API_KEY are set, the briefing
// also includes the top 5 high-confidence signals from PRRJECT_FATBABY.

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
func runMorningBriefing(ctx context.Context, iduna *IdunaClient, push PushFunc, stateDir, earningsCalDir string) {
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

	// Fetch top FatBaby signals if signalapi is configured.
	signals := fetchTopSignals(ctx)

	title, body := buildBriefingMessage(recent)
	if len(signals) > 0 {
		body += "\n\n── SIGNAL INTELLIGENCE ──"
		for _, s := range signals {
			body += fmt.Sprintf("\n  %-6s  %s  [%s]", s.Ticker, s.Summary, s.SignalType)
		}
	}

	// Earnings calendar section: upcoming reports in the next 7 days.
	if earningsCalDir != "" {
		if upcoming := loadUpcomingEarnings(earningsCalDir, 7); len(upcoming) > 0 {
			body += "\n\n── EARNINGS THIS WEEK ──"
			for _, e := range upcoming {
				timing := ""
				if e.BeforeMarket != nil {
					if *e.BeforeMarket {
						timing = " BMO"
					} else {
						timing = " AMC"
					}
				}
				body += fmt.Sprintf("\n  %-6s  %s%s  [%s]", e.Ticker, e.ReportDate, timing, e.Status)
			}
		}
	}

	log.Printf("briefing: firing morning push — %d apples in 24h, %d signals", len(recent), len(signals))

	go push(title, body, map[string]string{
		"deep_link":     "mjolnir://feed",
		"briefing_date": time.Now().UTC().Format("2006-01-02"),
		"apple_count":   fmt.Sprintf("%d", len(recent)),
	})

	markBriefingSent(stateDir)

	// File a status Apple so the briefing is visible in the audit trail.
	if iduna != nil {
		go func() {
			appleCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, _ = iduna.PostApple(appleCtx, ApplePayload{
				AppleType:  "status",
				SourceRepo: "EMILY",
				Title:      fmt.Sprintf("morning briefing sent — %d apples in 24h", len(recent)),
				Body:       body,
			})
		}()
	}
}

// FatBabySignalEntry is a minimal view of a signal from PRRJECT_FATBABY signalapi.
type FatBabySignalEntry struct {
	Ticker         string  `json:"ticker"`
	SignalType     string  `json:"signal_type"`
	Importance     int     `json:"importance"`
	Sentiment      float64 `json:"sentiment"`
	Summary        string  `json:"summary"`
	ImpactAnalysis string  `json:"impact_analysis"`
}

// fetchTopSignals fetches the top signals (importance ≥ 7) from the
// PRRJECT_FATBABY signalapi. Returns at most 5 results.
// FATBABY_SIGNAL_API_URL must be set; FATBABY_SIGNAL_API_KEY is optional.
// Non-fatal: returns nil on any error.
func fetchTopSignals(ctx context.Context) []FatBabySignalEntry {
	base := strings.TrimRight(os.Getenv("FATBABY_SIGNAL_API_URL"), "/")
	if base == "" {
		return nil
	}
	apiKey := os.Getenv("FATBABY_SIGNAL_API_KEY")
	url := base + "/v1/data-quality"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	defer resp.Body.Close()

	var quality struct {
		ByTicker []struct {
			Ticker       string  `json:"ticker"`
			AvgConf      float64 `json:"avg_confidence"`
			SignalCount  int     `json:"signal_count"`
			AvgImp       float64 `json:"avg_importance"`
		} `json:"by_ticker"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&quality); err != nil {
		return nil
	}

	// Keep tickers with avg_confidence ≥ 0.75 and ≥ 1 signal, limit to 5.
	var out []FatBabySignalEntry
	for _, tq := range quality.ByTicker {
		if tq.AvgConf < 0.75 || tq.SignalCount == 0 {
			continue
		}
		out = append(out, FatBabySignalEntry{
			Ticker:     tq.Ticker,
			SignalType: fmt.Sprintf("avg_importance=%.1f", tq.AvgImp),
			Importance: tq.SignalCount,
			Summary:    fmt.Sprintf("confidence=%.2f (%d signals)", tq.AvgConf, tq.SignalCount),
		})
		if len(out) >= 5 {
			break
		}
	}
	return out
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

// briefingEarningsDate is the minimal shape we need from dates.ndjson.
type briefingEarningsDate struct {
	Ticker       string  `json:"ticker"`
	ReportDate   string  `json:"report_date"`
	Status       string  `json:"status"`
	BeforeMarket *bool   `json:"before_market"`
}

// loadUpcomingEarnings reads var/earnings-calendar/dates.ndjson and returns
// records whose report_date falls in [today, today+days), sorted soonest first.
// Returns nil on any error (non-fatal: briefing still sends without this section).
func loadUpcomingEarnings(calDir string, days int) []briefingEarningsDate {
	today := time.Now().UTC().Format("2006-01-02")
	cutoff := time.Now().UTC().AddDate(0, 0, days).Format("2006-01-02")

	f, err := os.Open(filepath.Join(calDir, "dates.ndjson"))
	if err != nil {
		return nil
	}
	defer f.Close()

	// Deduplicate by ticker+date; prefer higher status (confirmed>announced>backfilled).
	statusPriority := map[string]int{"confirmed": 3, "announced": 2, "backfilled": 1}
	type key struct{ ticker, date string }
	best := map[key]briefingEarningsDate{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var d briefingEarningsDate
		if err := json.Unmarshal([]byte(line), &d); err != nil {
			continue
		}
		if d.Ticker == "" || d.ReportDate < today || d.ReportDate >= cutoff {
			continue
		}
		k := key{d.Ticker, d.ReportDate}
		if cur, ok := best[k]; !ok || statusPriority[d.Status] > statusPriority[cur.Status] {
			best[k] = d
		}
	}

	out := make([]briefingEarningsDate, 0, len(best))
	for _, d := range best {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ReportDate == out[j].ReportDate {
			return out[i].Ticker < out[j].Ticker
		}
		return out[i].ReportDate < out[j].ReportDate
	})
	return out
}
