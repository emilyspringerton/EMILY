# MJOLNIR — Emily Prime Context

**For future Emily Prime sessions — read this before touching FCM, push tokens, or Android.**

MJOLNIR is Emily Springerton's Android app. It receives push notifications from Emily Prime
and shows the Apple feed. The server-side implementation is complete as of 2026-06-09.

**What's live**: IDUNA push_tokens table + API (`/api/v1/push-tokens`). FCM sender package
(`emily-agent/pkg/fcm/`). Triage escalation path fires FCM push when `RequiresCEOVisibility == true`.
IDUNA `GetPushToken("mjolnir-emily")` resolves the device token.

**What's pending**: Android project skeleton (Kotlin/Compose). FCM project creation in Firebase
Console. `google-services.json` in the MJOLNIR repo. Emily registering her device for the first time.

**Environment variables needed** (not yet set):  
`FCM_PROJECT_ID`, `FCM_SERVICE_ACCOUNT_JSON` (or `FCM_SERVICE_ACCOUNT_FILE`).  
Without these, FCM sender is disabled gracefully — triage continues, Gmail still works.

**Agent name for MJOLNIR device**: `"mjolnir-emily"` — used in all push token lookups.

**Full spec**: `EMILY/docs/MJOLNIR_INTEGRATION.md`  
**Android spec**: `MJOLNIR/docs/NORTHSTAR.md`, `MJOLNIR/docs/SPEC.md`
