// emily-agent/alerting.go
// Check-in monitor alerting worker.
//
// Polls IDUNA /api/v1/monitors/overdue every checkInterval seconds.
// For each overdue monitor: fires Slack alert (+ email when configured),
// then marks it alerted via POST /api/v1/monitors/:id/alerted so we don't
// re-fire until the monitor recovers (next check-in clears alerted_at).

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	alertingPollInterval = 5 * time.Minute
	alertingHTTPTimeout  = 15 * time.Second
)

// MonitorRecord mirrors the JSON returned by IDUNA /api/v1/monitors/overdue.
type MonitorRecord struct {
	ID                int64      `json:"id"`
	Name              string     `json:"name"`
	Slug              string     `json:"slug"`
	TimeoutSeconds    int        `json:"timeout_seconds"`
	Owner             string     `json:"owner"`
	LastCheckinAt     *time.Time `json:"last_checkin_at,omitempty"`
	AlertSlackChannel string     `json:"alert_slack_channel,omitempty"`
	AlertEmail        string     `json:"alert_email,omitempty"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
}

// CheckinAlertWorker polls IDUNA for overdue monitors and fires notifications.
type CheckinAlertWorker struct {
	Iduna  *IdunaClient
	Slack  *SlackNotifier
	Gmail  *GmailClient
	AlertTo string // fallback email recipient; used when monitor has no alert_email set
}

// NewCheckinAlertWorker builds a worker from the current environment.
func NewCheckinAlertWorker(iduna *IdunaClient, slack *SlackNotifier, gmail *GmailClient) *CheckinAlertWorker {
	return &CheckinAlertWorker{
		Iduna:   iduna,
		Slack:   slack,
		Gmail:   gmail,
		AlertTo: envOr("ALERT_EMAIL", ""),
	}
}

// Run starts the polling loop. Blocks until ctx is cancelled.
func (w *CheckinAlertWorker) Run(ctx context.Context) {
	if w.Iduna == nil {
		log.Println("[alerting] IDUNA client not configured — check-in alerting disabled")
		return
	}
	log.Printf("[alerting] check-in alert worker started (poll interval %s)", alertingPollInterval)
	ticker := time.NewTicker(alertingPollInterval)
	defer ticker.Stop()
	w.poll(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.poll(ctx)
		}
	}
}

func (w *CheckinAlertWorker) poll(ctx context.Context) {
	monitors, err := w.fetchOverdue(ctx)
	if err != nil {
		log.Printf("[alerting] fetch overdue: %v", err)
		return
	}
	if len(monitors) == 0 {
		return
	}
	log.Printf("[alerting] %d overdue monitor(s)", len(monitors))
	for _, m := range monitors {
		w.fireAlert(ctx, m)
	}
}

func (w *CheckinAlertWorker) fetchOverdue(ctx context.Context) ([]MonitorRecord, error) {
	if err := w.Iduna.authenticate(ctx); err != nil {
		return nil, fmt.Errorf("auth: %w", err)
	}
	url := strings.TrimRight(w.Iduna.baseURL, "/") + "/api/v1/monitors/overdue"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	w.Iduna.mu.Lock()
	token := w.Iduna.token
	w.Iduna.mu.Unlock()
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: alertingHTTPTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IDUNA returned %d for monitors/overdue", resp.StatusCode)
	}
	var result struct {
		Monitors []MonitorRecord `json:"monitors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Monitors, nil
}

func (w *CheckinAlertWorker) fireAlert(ctx context.Context, m MonitorRecord) {
	log.Printf("[alerting] firing alert for monitor %q (id=%d, timeout=%ds)", m.Name, m.ID, m.TimeoutSeconds)

	// Slack alert.
	if w.Slack != nil && w.Slack.Enabled() {
		if err := w.Slack.SendCheckinMiss(m.Name, m.Slug, m.TimeoutSeconds, m.LastCheckinAt); err != nil {
			log.Printf("[alerting] slack alert for %q: %v", m.Name, err)
		}
	}

	// Email alert.
	to := m.AlertEmail
	if to == "" {
		to = w.AlertTo
	}
	if to != "" && w.Gmail != nil {
		var lastStr string
		if m.LastCheckinAt == nil {
			lastStr = "never"
		} else {
			lastStr = fmt.Sprintf("%s ago", time.Since(*m.LastCheckinAt).Round(time.Second))
		}
		alert := AlertEmail{
			To:      to,
			Subject: fmt.Sprintf("[EMILY ALERT] Check-in missed: %s", m.Name),
			Body: fmt.Sprintf(
				"Monitor: %s\nOwner: %s\nTimeout: %ds\nLast check-in: %s\nStatus: %s\n\n"+
					"Check-in URL: %s/api/v1/monitors/checkin/%s\n\n"+
					"This alert was fired automatically by Emily Prime.",
				m.Name, m.Owner, m.TimeoutSeconds, lastStr, m.Status,
				strings.TrimRight(w.Iduna.baseURL, "/"), m.Slug,
			),
		}
		if err := w.Gmail.SendAlert(ctx, alert); err != nil {
			log.Printf("[alerting] email alert for %q: %v", m.Name, err)
		}
	}

	// Mark alerted in IDUNA so we don't re-fire until the monitor recovers.
	if err := w.markAlerted(ctx, m.ID); err != nil {
		log.Printf("[alerting] mark alerted for %q: %v", m.Name, err)
	}
}

func (w *CheckinAlertWorker) markAlerted(ctx context.Context, id int64) error {
	url := fmt.Sprintf("%s/api/v1/monitors/%d/alerted",
		strings.TrimRight(w.Iduna.baseURL, "/"), id)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader([]byte("{}")))
	if err != nil {
		return err
	}
	w.Iduna.mu.Lock()
	token := w.Iduna.token
	w.Iduna.mu.Unlock()
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: alertingHTTPTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("IDUNA returned %d for mark-alerted", resp.StatusCode)
	}
	return nil
}
