# Emily Prime ↔ MJOLNIR Integration Spec
## How Emily Prime pushes intelligence to Emily's Android phone

**Author**: Emily Prime (fulfilling directed task 3146266896637121286)  
**Date**: 2026-06-09  
**Status**: Implementation complete for server side; Android skeleton pending

---

## Overview

Emily Prime is the sender. MJOLNIR is the receiver. Push notifications are the
last mile of the autonomous loop — when something critical happens in the agent
network, Emily Springerton's phone rings within seconds.

This document covers:
1. Which Apple types trigger FCM pushes
2. How Emily Prime resolves device tokens
3. The morning briefing cron
4. The integration seam in the codebase

---

## 1. Apple Severity Thresholds for Push

### Always push (CEO-visible escalations)
These come from `runPrimeTriage` in `integration.go` when an observation scores
`RequiresCEOVisibility = true`:

```go
// integration.go — triage loop
if t.RequiresCEOVisibility {
    go push(title, body, data)   // goroutine, never blocks triage
}
```

Triggers include:
- Observation severity `"critical"` + strategic score ≥ 80
- `ObservationType == "escalation"` regardless of score
- Signal anomalies affecting ≥ 5 watchlist tickers simultaneously
- Auth failures on any IDUNA agent (security event)

### Push on major RSI milestone (future — not yet wired)
When `buildCycleApple` produces `appleType == "improvement"` with
`tokens_saved > 5000` or `signals_added > 10`, fire a HIGH-priority push.
Deferred until RSI cycle throughput is measurable.

### Never push
- `appleType == "status"` (routine health checks)
- `appleType == "observation"` without CEO visibility
- Apples from `source_repo == "MJOLNIR"` (would create push loops)

---

## 2. Device Token Resolution

Emily Prime resolves the FCM token lazily — on demand, when escalation fires:

```go
// cron.go — runPrimeTriageCycle
push = func(title, body string, data map[string]string) {
    deviceToken, err := iduna.GetPushToken(ctx, "mjolnir-emily")
    if err != nil || deviceToken == "" {
        log.Printf("fcm: no device token (push skipped)")
        return
    }
    sender.Send(ctx, deviceToken, fcm.Message{...})
}
```

**Agent name**: `mjolnir-emily` — the logical name registered by the MJOLNIR app
at first launch via `POST /api/v1/push-tokens`.

**IDUNA endpoint**: `GET /api/v1/push-tokens/mjolnir-emily`  
**Auth**: M2M bearer token (Emily Prime's `IDUNA_AGENT_SECRET`)  
**Required permission**: `push_tokens.read` (granted to emily-prime agent)

**Graceful degradation**: if IDUNA is unreachable or no token is registered,
push is skipped silently (logged). The triage cycle continues normally.

---

## 3. Morning Briefing (Daily at 09:00 UTC)

Implementation: `EMILY/emily-agent/briefing.go` (or wired into cron.go).

The briefing fires from a daily cron function that:
1. Queries IDUNA `GET /api/v1/apples?limit=50` for Apples filed in the last 24h
2. Groups by source_repo and apple_type
3. Builds a summary: "12 Apples overnight: 8 improvements, 2 observations, 2 RSI iterations"
4. Fires FCM push with priority `"normal"` (does not bypass DND)
5. Data payload includes `{"deep_link": "mjolnir://feed"}` to open the Apple feed

The briefing is skipped if:
- No Apples were filed in the last 24h
- MJOLNIR device token is not registered

Schedule: added to the daily 09:00 UTC cron entry alongside the RSI loop.

---

## 4. FCM Sender Config (Required Environment Variables)

Set these in the environment where `emily-agent` runs (e.g. `start-emily-agent.sh`):

```bash
FCM_PROJECT_ID=einhorn-mjolnir
FCM_SERVICE_ACCOUNT_JSON='{"type":"service_account","project_id":"einhorn-mjolnir",...}'
# OR:
FCM_SERVICE_ACCOUNT_FILE=/home/fatbaby/.config/emily/fcm-service-account.json
```

The service account JSON comes from the Firebase Console:
`Project Settings → Service accounts → Generate new private key`

Emily Prime logs at startup: `fcm: sender initialized (project=einhorn-mjolnir)`  
If not configured: `fcm: init failed (push notifications disabled)` — system continues.

---

## 5. Codebase Seams

| File | Role |
|------|------|
| `emily-agent/pkg/fcm/sender.go` | FCM HTTP v1 sender, OAuth2 token caching |
| `emily-agent/pkg/fcm/jwt.go` | RS256 service account JWT builder |
| `emily-agent/cron.go` | Sender init + `PushFunc` closure in `runPrimeTriageCycle` |
| `emily-agent/integration.go` | `runPrimeTriage(…, push PushFunc)` — fires on CEO escalations |
| `emily-agent/iduna.go` | `GetPushToken(agentName)` — resolves FCM token from IDUNA |
| `IDUNA/internal/store/sqlite.go` | `UpsertPushToken` + `GetPushToken` impl |
| `IDUNA/internal/http/handlers/push_tokens.go` | `POST /api/v1/push-tokens` + `GET /…/{name}` |
| `MJOLNIR/app/…/FcmTokenManager.kt` | Android: registers token with IDUNA on launch |
| `MJOLNIR/app/…/MjolnirMessagingService.kt` | Android: receives FCM, routes to notification channels |

---

## 6. MJOLNIR App Registration Flow

On first Android launch:
1. App completes Google Sign-In → IDUNA JWT stored in `EncryptedSharedPreferences`
2. `FirebaseMessaging.getInstance().token` → FCM token string
3. App calls `POST /api/v1/push-tokens` with `{agent_name: "mjolnir-emily", platform: "android", fcm_token: "...", fingerprint: "..."}`
4. IDUNA stores in `push_tokens` table
5. Emily Prime can now resolve the token and fire pushes

---

## 7. Android Notification Channels

| Channel ID | Name | Importance | Bypasses DND |
|-----------|------|-----------|-------------|
| `MJOLNIR_CRITICAL` | Critical Alerts | URGENT | Yes |
| `MJOLNIR_HIGH` | High Priority | HIGH | No |
| `MJOLNIR_NORMAL` | Activity | DEFAULT | No |

Emily Prime uses `"high"` priority for CEO-visibility escalations (→ `MJOLNIR_HIGH`).
Morning briefing uses `"normal"` (→ `MJOLNIR_NORMAL`).
Future: manually triggered critical alerts use `"critical"` (→ `MJOLNIR_CRITICAL`, bypasses DND).

---

*EMILY PRIME ↔ MJOLNIR INTEGRATION | 2026-06-09*
*The push channel is live on the server side. Android registration completes the circuit.*
