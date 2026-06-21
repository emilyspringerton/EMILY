// api_push.go — POST /api/v1/emily/push/test (S32-02)
//
// Fires a test FCM notification to the registered MJOLNIR device token.
// Requires FCM to be configured (FCM_PROJECT_ID + FCM_SERVICE_ACCOUNT_JSON).
// Requires at least one push token registered in IDUNA for agent "mjolnir-emily".
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"emily-agent/pkg/fcm"
)

func (s *Server) handlePushTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.cycle == nil {
		http.Error(w, "cycle not initialized", http.StatusServiceUnavailable)
		return
	}

	ac := s.cycle

	if ac.fcmSender == nil {
		writeJSONError(w, http.StatusServiceUnavailable,
			"FCM not configured — set FCM_PROJECT_ID + FCM_SERVICE_ACCOUNT_JSON and restart")
		return
	}
	if ac.iduna == nil {
		writeJSONError(w, http.StatusServiceUnavailable,
			"IDUNA not configured — IDUNA_BASE_URL + IDUNA_AGENT_SECRET required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	deviceToken, err := ac.iduna.GetPushToken(ctx, "mjolnir-emily")
	if err != nil || deviceToken == "" {
		writeJSONError(w, http.StatusNotFound,
			"no push token registered for mjolnir-emily — complete S32-01 first")
		return
	}

	if err := ac.fcmSender.Send(ctx, deviceToken, fcm.Message{
		Title: "Emily Prime — test push",
		Body:  "S32-02: FCM push test fired from emily-agent at " + time.Now().Format(time.RFC3339),
		Data: map[string]string{
			"event": "push_test",
			"ts":    time.Now().Format(time.RFC3339),
		},
	}); err != nil {
		writeJSONError(w, http.StatusBadGateway, "FCM send failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":       "sent",
		"device_token": deviceToken[:min(16, len(deviceToken))] + "…",
		"sent_at":      time.Now().Format(time.RFC3339),
	})
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
