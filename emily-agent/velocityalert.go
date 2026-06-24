package main

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

	"emily-agent/pkg/fcm"
)

// velocityAlertEntry is one alert from the signalapi /v1/velocity-alerts endpoint.
type velocityAlertEntry struct {
	Ticker      string    `json:"ticker"`
	Count       int       `json:"count"`
	WindowStart time.Time `json:"window_start"`
	WindowEnd   time.Time `json:"window_end"`
}

// fetchVelocityAlerts polls FATBABY_SIGNAL_API_URL/v1/velocity-alerts.
// Returns nil if the env var is not set or the request fails (non-fatal).
func fetchVelocityAlerts(ctx context.Context) []velocityAlertEntry {
	base := strings.TrimRight(os.Getenv("FATBABY_SIGNAL_API_URL"), "/")
	if base == "" {
		return nil
	}
	apiKey := os.Getenv("FATBABY_SIGNAL_API_KEY")
	url := base + "/v1/velocity-alerts"
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
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	var result struct {
		Alerts []velocityAlertEntry `json:"alerts"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil
	}
	return result.Alerts
}

// observeVelocityAlerts fetches velocity alerts and dispatches Apple + FCM for each.
// Called from the OBSERVE phase. All errors are logged, never fatal.
func (ac *AutonomousCycle) observeVelocityAlerts(ctx context.Context, cycleNum int) {
	alerts := fetchVelocityAlerts(ctx)
	if len(alerts) == 0 {
		return
	}
	for _, a := range alerts {
		msg := fmt.Sprintf("Signal velocity spike: %s produced %d high-confidence signals (%s–%s UTC)",
			a.Ticker, a.Count,
			a.WindowStart.UTC().Format("15:04:05"),
			a.WindowEnd.UTC().Format("15:04:05"),
		)
		log.Printf("[cycle %d] velocity alert: %s", cycleNum, msg)

		if ac.iduna != nil {
			alertCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			if _, err := ac.iduna.PostApple(alertCtx, ApplePayload{
				SourceRepo: "PRRJECT_FATBABY",
				AppleType:  "escalation",
				Title:      fmt.Sprintf("velocity alert: %s ×%d", a.Ticker, a.Count),
				Body:       msg,
			}); err != nil {
				log.Printf("[cycle %d] velocity apple warn: %v", cycleNum, err)
			}
			cancel()
		}

		if ac.fcmSender != nil && ac.iduna != nil {
			pushCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			deviceToken, err := ac.iduna.GetPushToken(pushCtx, "mjolnir-emily")
			cancel()
			if err == nil && deviceToken != "" {
				_ = ac.fcmSender.Send(ctx, deviceToken, fcm.Message{
					Title:    fmt.Sprintf("⚡ %s velocity spike", a.Ticker),
					Body:     fmt.Sprintf("%d signals in <1h (importance>60)", a.Count),
					Priority: "high",
					Data:     map[string]string{"ticker": a.Ticker, "count": fmt.Sprintf("%d", a.Count)},
				})
			}
		}
	}
}
