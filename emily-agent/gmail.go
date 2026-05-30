// emily-agent/gmail.go
// Gmail integration for Emily Prime — CEO inbox triage and outbound signal alerts.
//
// Emily Prime monitors the CEO inbox, classifies by urgency and strategic relevance,
// surfaces what needs a decision. Outbound: formatted signal alerts and weekly pulse.
//
// Auth: OAuth2 refresh token flow (no browser required after initial setup).
// Credentials are loaded from environment variables — no files on disk.
//
// Required env vars:
//   GMAIL_CLIENT_ID       — Google OAuth2 client ID
//   GMAIL_CLIENT_SECRET   — Google OAuth2 client secret
//   GMAIL_REFRESH_TOKEN   — OAuth2 refresh token (from one-time browser auth)
//   GMAIL_CEO_ADDRESS     — CEO email address to monitor and send to
//
// Optional:
//   GMAIL_SEND_FROM       — address Emily sends from (defaults to CEO address)
//   GMAIL_MAX_TRIAGE      — max emails to triage per cycle (default 20)

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// OAuth2 token management
// ---------------------------------------------------------------------------

type oauth2Token struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"-"`
}

type GmailCredentials struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	CEOAddress   string
	SendFrom     string
}

func loadGmailCredentials() (GmailCredentials, bool) {
	creds := GmailCredentials{
		ClientID:     os.Getenv("GMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("GMAIL_CLIENT_SECRET"),
		RefreshToken: os.Getenv("GMAIL_REFRESH_TOKEN"),
		CEOAddress:   os.Getenv("GMAIL_CEO_ADDRESS"),
		SendFrom:     os.Getenv("GMAIL_SEND_FROM"),
	}
	if creds.ClientID == "" || creds.ClientSecret == "" || creds.RefreshToken == "" || creds.CEOAddress == "" {
		return creds, false
	}
	if creds.SendFrom == "" {
		creds.SendFrom = creds.CEOAddress
	}
	return creds, true
}

// ---------------------------------------------------------------------------
// GmailClient
// ---------------------------------------------------------------------------

type GmailClient struct {
	creds  GmailCredentials
	client *http.Client
	token  *oauth2Token
}

func NewGmailClient(creds GmailCredentials) *GmailClient {
	return &GmailClient{
		creds:  creds,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *GmailClient) accessToken(ctx context.Context) (string, error) {
	if g.token != nil && time.Now().Before(g.token.ExpiresAt.Add(-60*time.Second)) {
		return g.token.AccessToken, nil
	}
	data := url.Values{
		"client_id":     {g.creds.ClientID},
		"client_secret": {g.creds.ClientSecret},
		"refresh_token": {g.creds.RefreshToken},
		"grant_type":    {"refresh_token"},
	}
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://oauth2.googleapis.com/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gmail token refresh: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("gmail token refresh status=%d: %s", resp.StatusCode, raw)
	}
	var tok struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(raw, &tok); err != nil {
		return "", fmt.Errorf("gmail token parse: %w", err)
	}
	g.token = &oauth2Token{
		AccessToken: tok.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second),
	}
	return g.token.AccessToken, nil
}

func (g *GmailClient) gmailGET(ctx context.Context, path string) ([]byte, error) {
	tok, err := g.accessToken(ctx)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://gmail.googleapis.com/gmail/v1"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("gmail GET %s status=%d: %s", path, resp.StatusCode, raw)
	}
	return raw, nil
}

func (g *GmailClient) gmailPOST(ctx context.Context, path string, body any) ([]byte, error) {
	tok, err := g.accessToken(ctx)
	if err != nil {
		return nil, err
	}
	payload, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://gmail.googleapis.com/gmail/v1"+path,
		bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("gmail POST %s status=%d: %s", path, resp.StatusCode, raw)
	}
	return raw, nil
}

// ---------------------------------------------------------------------------
// InboxMessage — lightweight representation of an email
// ---------------------------------------------------------------------------

type InboxMessage struct {
	ID        string    `json:"id"`
	ThreadID  string    `json:"thread_id"`
	From      string    `json:"from"`
	Subject   string    `json:"subject"`
	Snippet   string    `json:"snippet"`
	Date      time.Time `json:"date"`
	Labels    []string  `json:"labels"`
	BodyText  string    `json:"body_text,omitempty"`
	Unread    bool      `json:"unread"`
}

// TriagedMessage is an InboxMessage with Emily Prime's relevance assessment.
type TriagedMessage struct {
	Message           InboxMessage `json:"message"`
	Urgency           string       `json:"urgency"`           // critical|high|normal|low
	Category          string       `json:"category"`          // strategic|operational|noise|unknown
	StrategicRelevance float64     `json:"strategic_relevance"` // 0.0–1.0
	SuggestedAction   string       `json:"suggested_action"`
	NeedsReply        bool         `json:"needs_reply"`
}

// ---------------------------------------------------------------------------
// ReadInbox — list recent unread messages in the CEO inbox
// ---------------------------------------------------------------------------

func (g *GmailClient) ReadInbox(ctx context.Context, maxResults int) ([]InboxMessage, error) {
	if maxResults <= 0 {
		maxResults = 20
	}
	path := fmt.Sprintf("/users/me/messages?maxResults=%d&labelIds=INBOX&q=is:unread", maxResults)
	raw, err := g.gmailGET(ctx, path)
	if err != nil {
		return nil, err
	}
	var list struct {
		Messages []struct {
			ID       string `json:"id"`
			ThreadID string `json:"threadId"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}

	var msgs []InboxMessage
	for _, m := range list.Messages {
		msg, err := g.fetchMessage(ctx, m.ID)
		if err != nil {
			log.Printf("gmail fetch_message id=%s err=%v", m.ID, err)
			continue
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (g *GmailClient) fetchMessage(ctx context.Context, id string) (InboxMessage, error) {
	raw, err := g.gmailGET(ctx, "/users/me/messages/"+id+"?format=metadata&metadataHeaders=From&metadataHeaders=Subject&metadataHeaders=Date")
	if err != nil {
		return InboxMessage{}, err
	}
	var msg struct {
		ID       string   `json:"id"`
		ThreadID string   `json:"threadId"`
		Snippet  string   `json:"snippet"`
		LabelIDs []string `json:"labelIds"`
		Payload  struct {
			Headers []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"headers"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(raw, &msg); err != nil {
		return InboxMessage{}, err
	}
	out := InboxMessage{
		ID:       msg.ID,
		ThreadID: msg.ThreadID,
		Snippet:  msg.Snippet,
		Labels:   msg.LabelIDs,
	}
	for _, h := range msg.Payload.Headers {
		switch strings.ToLower(h.Name) {
		case "from":
			out.From = h.Value
		case "subject":
			out.Subject = h.Value
		case "date":
			if t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", h.Value); err == nil {
				out.Date = t
			}
		}
	}
	for _, l := range msg.LabelIDs {
		if l == "UNREAD" {
			out.Unread = true
		}
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// SendAlert — send a formatted signal alert to the CEO
// ---------------------------------------------------------------------------

type AlertEmail struct {
	To      string
	Subject string
	Body    string // plain text
}

func (g *GmailClient) SendAlert(ctx context.Context, alert AlertEmail) error {
	raw := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		g.creds.SendFrom, alert.To, alert.Subject, alert.Body)
	encoded := base64.URLEncoding.EncodeToString([]byte(raw))
	_, err := g.gmailPOST(ctx, "/users/me/messages/send", map[string]string{"raw": encoded})
	return err
}

// ---------------------------------------------------------------------------
// FormatSignalAlert — formats a FatBaby observation as a CEO-ready email
// ---------------------------------------------------------------------------

func FormatSignalAlert(obs Observation, triage TriageResult) AlertEmail {
	var sb strings.Builder
	sb.WriteString("EMILY PRIME — SIGNAL ALERT\n")
	sb.WriteString(strings.Repeat("─", 50) + "\n\n")
	fmt.Fprintf(&sb, "Severity:   %s\n", strings.ToUpper(obs.Severity))
	fmt.Fprintf(&sb, "Relevance:  %.0f%%\n", triage.StrategicRelevance*100)
	if len(obs.AffectedTickers) > 0 {
		fmt.Fprintf(&sb, "Tickers:    %s\n", strings.Join(obs.AffectedTickers, ", "))
	}
	if len(triage.TriggeredSignals) > 0 {
		fmt.Fprintf(&sb, "Signals:    %s\n", strings.Join(triage.TriggeredSignals, ", "))
	}
	sb.WriteString("\n")

	if obs.Detail != "" {
		fmt.Fprintf(&sb, "Detail:\n%s\n\n", obs.Detail)
	} else if obs.Findings != "" {
		fmt.Fprintf(&sb, "Findings:\n%s\n\n", obs.Findings)
	}

	if obs.RecommendedAction != "" {
		fmt.Fprintf(&sb, "Recommended action:\n%s\n\n", obs.RecommendedAction)
	}

	fmt.Fprintf(&sb, "Rationale: %s\n\n", triage.Rationale)
	fmt.Fprintf(&sb, "Timestamp: %s\n", obs.Timestamp)
	sb.WriteString("\n— Emily Prime\n")

	urgencyPrefix := ""
	switch {
	case triage.StrategicRelevance >= 0.9:
		urgencyPrefix = "[CRITICAL] "
	case triage.StrategicRelevance >= 0.75:
		urgencyPrefix = "[HIGH] "
	}

	subject := urgencyPrefix + "Signal Alert: " + obs.Summary
	if len(obs.AffectedTickers) > 0 {
		subject += " (" + strings.Join(obs.AffectedTickers, ", ") + ")"
	}

	return AlertEmail{
		To:      "",
		Subject: subject,
		Body:    sb.String(),
	}
}

// ---------------------------------------------------------------------------
// WeeklyPulse — auto-generated strategic summary
// ---------------------------------------------------------------------------

type WeeklyPulse struct {
	PeriodStart    time.Time     `json:"period_start"`
	PeriodEnd      time.Time     `json:"period_end"`
	Observations   []Observation `json:"observations"`
	TriageResults  []TriageResult `json:"triage_results"`
	CEOEscalations int           `json:"ceo_escalations"`
	TasksIssued    int           `json:"tasks_issued"`
}

func FormatWeeklyPulse(pulse WeeklyPulse) AlertEmail {
	var sb strings.Builder
	sb.WriteString("EMILY PRIME — WEEKLY PULSE\n")
	sb.WriteString(strings.Repeat("─", 50) + "\n\n")
	fmt.Fprintf(&sb, "Period: %s → %s\n\n",
		pulse.PeriodStart.Format("2006-01-02"),
		pulse.PeriodEnd.Format("2006-01-02"))
	fmt.Fprintf(&sb, "Observations received:  %d\n", len(pulse.Observations))
	fmt.Fprintf(&sb, "CEO escalations:        %d\n", pulse.CEOEscalations)
	fmt.Fprintf(&sb, "Directed tasks issued:  %d\n\n", pulse.TasksIssued)

	var highRelevance []TriageResult
	for _, t := range pulse.TriageResults {
		if t.StrategicRelevance >= 0.75 {
			highRelevance = append(highRelevance, t)
		}
	}
	if len(highRelevance) > 0 {
		sb.WriteString("HIGH-RELEVANCE SIGNALS THIS WEEK:\n")
		for i, t := range highRelevance {
			fmt.Fprintf(&sb, "  %d. %s (relevance=%.0f%%) — %s\n",
				i+1, t.ObservationTimestamp, t.StrategicRelevance*100, t.Rationale)
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("No high-relevance signals this week.\n\n")
	}

	sb.WriteString("— Emily Prime\n")

	return AlertEmail{
		Subject: fmt.Sprintf("Emily Prime Weekly Pulse — %s", pulse.PeriodEnd.Format("2006-01-02")),
		Body:    sb.String(),
	}
}

// ---------------------------------------------------------------------------
// Register Gmail tools on Emily Prime's dispatcher
// ---------------------------------------------------------------------------

func registerGmailTools(d *ToolDispatcher, gmail *GmailClient) {
	d.Register(ToolDef{
		Name:        "gmail_read_inbox",
		Description: "Read recent unread messages from the CEO inbox. Returns subject, sender, snippet, date for each message.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"max_results": {Type: "number", Description: "Max messages to return (default 20)"},
			},
		},
	}, func(args map[string]any) (string, error) {
		max := 20
		if v, ok := args["max_results"].(float64); ok && v > 0 {
			max = int(v)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		msgs, err := gmail.ReadInbox(ctx, max)
		if err != nil {
			return "", err
		}
		if len(msgs) == 0 {
			return "inbox clear — no unread messages", nil
		}
		data, _ := json.MarshalIndent(msgs, "", "  ")
		return string(data), nil
	})

	d.Register(ToolDef{
		Name:        "gmail_send_alert",
		Description: "Send a signal alert email to the CEO. Use for material signals that exceed the CEO visibility threshold.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"subject": {Type: "string", Description: "Email subject line"},
				"body":    {Type: "string", Description: "Plain text email body"},
				"to":      {Type: "string", Description: "Recipient (default: CEO address from config)"},
			},
			Required: []string{"subject", "body"},
		},
	}, func(args map[string]any) (string, error) {
		to := stringArg(args, "to")
		if to == "" {
			to = gmail.creds.CEOAddress
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := gmail.SendAlert(ctx, AlertEmail{
			To:      to,
			Subject: stringArg(args, "subject"),
			Body:    stringArg(args, "body"),
		})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("alert sent to %s", to), nil
	})

	d.Register(ToolDef{
		Name:        "gmail_send_weekly_pulse",
		Description: "Generate and send the weekly pulse email summarising what FatBaby-Emily observed and what Emily Prime decided.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"to": {Type: "string", Description: "Recipient (default: CEO address)"},
			},
		},
	}, func(args map[string]any) (string, error) {
		to := stringArg(args, "to")
		if to == "" {
			to = gmail.creds.CEOAddress
		}
		// Build a minimal pulse from whatever observations are available.
		pulse := WeeklyPulse{
			PeriodStart: time.Now().AddDate(0, 0, -7),
			PeriodEnd:   time.Now(),
		}
		alert := FormatWeeklyPulse(pulse)
		alert.To = to
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := gmail.SendAlert(ctx, alert); err != nil {
			return "", err
		}
		return fmt.Sprintf("weekly pulse sent to %s", to), nil
	})
}

// buildGmailClient returns a configured GmailClient or nil if credentials aren't set.
func buildGmailClient() *GmailClient {
	creds, ok := loadGmailCredentials()
	if !ok {
		return nil
	}
	return NewGmailClient(creds)
}
