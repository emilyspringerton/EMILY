// pkg/fcm/sender.go — Firebase Cloud Messaging HTTP v1 API sender.
//
// Emily Prime uses this to push notifications to MJOLNIR (Emily's Android phone)
// when a critical Apple is filed or CEO visibility is required.
//
// Auth: Google OAuth2 with a service account JSON credential.
// The service account JSON is stored in IDUNA secrets (or as env var FCM_SERVICE_ACCOUNT_JSON).
// FCM project ID comes from env var FCM_PROJECT_ID.
//
// Usage:
//
//	sender, err := fcm.NewFromEnv()
//	if err != nil { ... }
//	err = sender.Send(ctx, deviceToken, fcm.Message{
//	    Title:    "Critical: PRRJECT_FATBABY",
//	    Body:     "governance_health_index score dropped to 0 for AMZN",
//	    Priority: "high",
//	    Data:     map[string]string{"apple_id": "42", "source_repo": "PRRJECT_FATBABY"},
//	})

package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const fcmSendURL = "https://fcm.googleapis.com/v1/projects/%s/messages:send"

// Message is a push notification to send.
type Message struct {
	Title    string            // notification title (shown in tray)
	Body     string            // notification body (≤ 140 chars recommended)
	Priority string            // "high" (wakes device) or "normal"
	Data     map[string]string // arbitrary key-value data payload for the app
}

// Sender dispatches FCM messages via the HTTP v1 API.
type Sender struct {
	projectID      string
	serviceAccount []byte // raw service account JSON

	mu          sync.Mutex
	cachedToken string
	tokenExpiry time.Time
}

// NewFromEnv creates a Sender from environment variables:
//
//	FCM_PROJECT_ID              — Google Cloud project ID (required)
//	FCM_SERVICE_ACCOUNT_JSON    — service account JSON inline (preferred)
//	FCM_SERVICE_ACCOUNT_FILE    — path to service account JSON file (fallback)
func NewFromEnv() (*Sender, error) {
	projectID := os.Getenv("FCM_PROJECT_ID")
	if projectID == "" {
		return nil, errors.New("FCM_PROJECT_ID not set")
	}

	var saJSON []byte
	if v := os.Getenv("FCM_SERVICE_ACCOUNT_JSON"); v != "" {
		saJSON = []byte(v)
	} else if f := os.Getenv("FCM_SERVICE_ACCOUNT_FILE"); f != "" {
		var err error
		saJSON, err = os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("read service account file: %w", err)
		}
	} else {
		return nil, errors.New("FCM_SERVICE_ACCOUNT_JSON or FCM_SERVICE_ACCOUNT_FILE not set")
	}

	return &Sender{projectID: projectID, serviceAccount: saJSON}, nil
}

// IsConfigured returns true if FCM env vars are present (use to gate the feature at startup).
func IsConfigured() bool {
	return os.Getenv("FCM_PROJECT_ID") != "" &&
		(os.Getenv("FCM_SERVICE_ACCOUNT_JSON") != "" || os.Getenv("FCM_SERVICE_ACCOUNT_FILE") != "")
}

// Send dispatches a push notification to the given FCM device token.
func (s *Sender) Send(ctx context.Context, deviceToken string, msg Message) error {
	if deviceToken == "" {
		return errors.New("device token is empty")
	}

	accessToken, err := s.getAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("fcm access token: %w", err)
	}

	priority := msg.Priority
	if priority == "" {
		priority = "high"
	}
	channelID := channelForPriority(priority)

	payload := map[string]any{
		"message": map[string]any{
			"token": deviceToken,
			"notification": map[string]string{
				"title": msg.Title,
				"body":  truncate(msg.Body, 200),
			},
			"data": msg.Data,
			"android": map[string]any{
				"priority": priority,
				"notification": map[string]string{
					"channel_id":   channelID,
					"click_action": "OPEN_APPLE_DETAIL",
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	url := fmt.Sprintf(fcmSendURL, s.projectID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fcm http: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fcm %d: %s", resp.StatusCode, respBody)
	}
	return nil
}

// getAccessToken returns a valid OAuth2 bearer token for the FCM scope,
// using the service account JSON. Caches the token until 60s before expiry.
func (s *Sender) getAccessToken(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cachedToken != "" && time.Until(s.tokenExpiry) > 60*time.Second {
		return s.cachedToken, nil
	}

	// Parse service account JSON to get key info.
	var sa struct {
		Type        string `json:"type"`
		ProjectID   string `json:"project_id"`
		PrivateKey  string `json:"private_key"`
		ClientEmail string `json:"client_email"`
		TokenURI    string `json:"token_uri"`
	}
	if err := json.Unmarshal(s.serviceAccount, &sa); err != nil {
		return "", fmt.Errorf("parse service account: %w", err)
	}
	if sa.TokenURI == "" {
		sa.TokenURI = "https://oauth2.googleapis.com/token"
	}

	// Build a JWT grant for the OAuth2 token endpoint.
	// Uses RS256 signing from the service account private key.
	jwt, err := buildServiceAccountJWT(sa.ClientEmail, sa.PrivateKey, sa.TokenURI,
		"https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return "", fmt.Errorf("build jwt: %w", err)
	}

	// Exchange JWT for access token.
	form := "grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Ajwt-bearer&assertion=" + jwt
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sa.TokenURI,
		strings.NewReader(form))
	if err != nil {
		return "", fmt.Errorf("build token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token exchange: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return "", errors.New("empty access token from token endpoint")
	}

	s.cachedToken = tokenResp.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return s.cachedToken, nil
}

func channelForPriority(priority string) string {
	switch priority {
	case "high":
		return "MJOLNIR_HIGH"
	case "critical":
		return "MJOLNIR_CRITICAL"
	default:
		return "MJOLNIR_NORMAL"
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
