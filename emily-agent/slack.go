// emily-agent/slack.go
// Slack notification integration for Emily Prime.
//
// Reads SLACK_WEBHOOK_URL from env. If unset, all Send calls are no-ops.
// Supports plain text messages and structured alert blocks.
//
// For full Slack app OAuth (bot token + channel posting), set SLACK_BOT_TOKEN
// and SLACK_DEFAULT_CHANNEL instead; the webhook path is simpler for most uses.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// SlackNotifier sends messages to a Slack webhook URL.
// Zero value is safe but sends nothing (Enabled() returns false).
type SlackNotifier struct {
	WebhookURL     string
	DefaultChannel string // optional; only used when sending via bot token (future)
	client         *http.Client
}

// NewSlackNotifier reads SLACK_WEBHOOK_URL from the environment.
func NewSlackNotifier() *SlackNotifier {
	return &SlackNotifier{
		WebhookURL:     os.Getenv("SLACK_WEBHOOK_URL"),
		DefaultChannel: os.Getenv("SLACK_DEFAULT_CHANNEL"),
		client:         &http.Client{Timeout: 10 * time.Second},
	}
}

// Enabled returns true when a webhook URL is configured.
func (s *SlackNotifier) Enabled() bool {
	return s.WebhookURL != ""
}

// Send posts a plain-text message to the configured webhook.
func (s *SlackNotifier) Send(text string) error {
	if !s.Enabled() {
		return nil
	}
	payload := map[string]any{"text": text}
	return s.post(payload)
}

// SendAlert posts a structured alert with level, title, and body.
// Level should be one of: info, warning, error, critical.
func (s *SlackNotifier) SendAlert(level, title, body string) error {
	if !s.Enabled() {
		return nil
	}

	emoji := map[string]string{
		"info":     ":information_source:",
		"warning":  ":warning:",
		"error":    ":x:",
		"critical": ":rotating_light:",
	}
	icon := emoji[level]
	if icon == "" {
		icon = ":bell:"
	}

	text := fmt.Sprintf("%s *[%s]* %s\n%s", icon, level, title, body)
	return s.Send(text)
}

// SendDown fires a site-down alert.
func (s *SlackNotifier) SendDown(serviceName, detail string, downSince time.Time) error {
	dur := time.Since(downSince).Round(time.Second)
	msg := fmt.Sprintf(":rotating_light: *SITE DOWN* — `%s`\nDown for: %s\n%s",
		serviceName, dur, detail)
	return s.Send(msg)
}

// SendUp fires a service-recovered notification.
func (s *SlackNotifier) SendUp(serviceName string, downDuration time.Duration) error {
	msg := fmt.Sprintf(":white_check_mark: *RECOVERED* — `%s` (was down %s)",
		serviceName, downDuration.Round(time.Second))
	return s.Send(msg)
}

// SendCheckinMiss fires an alert when a monitored service misses its check-in window.
func (s *SlackNotifier) SendCheckinMiss(monitorName, slug string, timeout int, lastCheckin *time.Time) error {
	var lastStr string
	if lastCheckin == nil {
		lastStr = "never"
	} else {
		lastStr = fmt.Sprintf("%s ago", time.Since(*lastCheckin).Round(time.Second))
	}
	msg := fmt.Sprintf(":skull: *CHECK-IN MISSED* — `%s`\nTimeout: %ds | Last seen: %s\nCheck-in URL: /api/v1/monitors/checkin/%s",
		monitorName, timeout, lastStr, slug)
	return s.Send(msg)
}

func (s *SlackNotifier) post(payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal: %w", err)
	}
	resp, err := s.client.Post(s.WebhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// slackNotifyOrLog sends a Slack message, logging errors instead of returning them.
func slackNotifyOrLog(s *SlackNotifier, text string) {
	if s == nil || !s.Enabled() {
		return
	}
	if err := s.Send(text); err != nil {
		log.Printf("[slack] send error: %v", err)
	}
}
