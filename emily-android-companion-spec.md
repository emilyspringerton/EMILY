# Emily Android Companion — Technical Specification
### Push notifications for agent escalation, with authenticated human response

**Status:** Draft v0.1
**Scope:** A mobile client whose single job is to get a human in the loop fast when Emily escalates, show enough context to decide, and send an authenticated decision back. Nothing more — this is a control surface, not a dashboard.

---

## 1. Purpose

Emily escalates a fraction of decisions to humans (legal/compliance, budget, strategic direction, values/ethics, novel situations, emergencies). Today those escalations sit in `emily/state/escalations/` until someone checks. This app removes the polling: when an escalation is created, the right person's phone buzzes within seconds, they see the context and options, and they respond. Their response unblocks (or halts) Emily.

The design priority, in order: **reliable delivery → trustworthy authorization → fast comprehension.** Because an approval here can authorize an autonomous agent to act, the authorization path is treated as security-critical, not as a convenience feature.

---

## 2. System context

```
┌──────────────┐   escalation     ┌─────────────────────┐
│   Emily      │  event created   │  Escalation Service │
│  (backend)   ├─────────────────►│  (new component)    │
└──────▲───────┘                  └─────────┬───────────┘
       │                                    │ push (FCM)
       │ decision applied                   ▼
       │                          ┌─────────────────────┐
       │                          │  Firebase Cloud      │
       │                          │  Messaging           │
       │                          └─────────┬───────────┘
       │                                    │
       │   POST /decision (authenticated)   ▼
       └──────────────────────────┌─────────────────────┐
                                   │  Android app        │
                                   │  (Kotlin/Compose)   │
                                   └─────────────────────┘
```

The **Escalation Service** is the only new backend piece. It owns device registration, push fan-out, decision intake, auth verification, and the audit log. Emily writes escalations to it and reads decisions from it; Emily never talks to FCM or the phone directly.

---

## 3. Escalation data model

The contract between Emily and the app. Emily already produces most of these fields per the delegation framework (context, options, recommendation, timeline impact); this just formalizes them.

```json
{
  "escalation_id": "esc_2026_0531_a17f",
  "created_at": "2026-05-31T14:03:22Z",
  "severity": "emergency | high | normal | info",
  "category": "legal | budget | strategy | values | novel | emergency | ops",
  "title": "Production deploy blocked: license ambiguity on new data source",
  "summary": "Emily wants to add source X to the training corpus but cannot confirm the license permits commercial use. Collection is paused pending a decision.",
  "context": {
    "what_happened": "…",
    "why_escalated": "Crosses the legal/compliance boundary; outside Emily's authority.",
    "what_is_blocked": "Phase-2 corpus expansion (source X only). Other collection continues.",
    "relevant_links": ["…"]
  },
  "options": [
    {
      "id": "opt_approve",
      "label": "Approve inclusion",
      "pros": ["…"], "cons": ["…"],
      "requires_strong_auth": true,
      "is_destructive": false
    },
    {
      "id": "opt_exclude",
      "label": "Exclude source X, continue without it",
      "pros": ["…"], "cons": ["…"],
      "requires_strong_auth": false,
      "is_destructive": false
    }
  ],
  "recommended_option_id": "opt_exclude",
  "timeline_impact": "Each hour blocked delays Phase 2 by ~1h. No hard deadline.",
  "deadline": null,
  "free_text_allowed": true,
  "state": "open | acknowledged | decided | expired | rescinded"
}
```

Notes:
- `requires_strong_auth` and `is_destructive` are set by Emily per option, so the app knows when to demand biometric confirmation and when to show an extra "are you sure" step.
- `deadline` drives auto-escalation (Section 8). `null` means no timeout, only reminders.
- `state` is authoritative on the server; the app reflects it and must handle a `rescinded` escalation arriving for one already shown (Emily resolved it itself).

---

## 4. Notification design

**Channels** (Android notification channels, set once at first launch):

| Channel | Severity | Sound/behavior | Bypasses DND |
|---|---|---|---|
| `emergency` | emergency | Full-screen intent, loud, repeating | Yes (registered as a high-priority/critical alert) |
| `high` | high | Heads-up, distinct sound | No |
| `normal` | normal | Standard heads-up | No |
| `info` | info | Silent, tray only | No |

- **Emergency** (e.g. Emily halted itself) uses an FCM high-priority message and a full-screen intent so it surfaces even on the lock screen. Treat this like a phone call, not a chat ping.
- Push payload is **data-only** (not a `notification` payload) so the app constructs the notification itself and controls channel routing even in the background.
- The push body carries the `escalation_id` and minimal metadata only — **never the full context**. The app fetches detail over the authenticated API. This keeps sensitive escalation content off the push transport and out of notification logs.

**Payload (FCM data message):**
```json
{
  "data": {
    "escalation_id": "esc_2026_0531_a17f",
    "severity": "high",
    "category": "legal",
    "title_preview": "Production deploy blocked: license ambiguity",
    "deadline": "null",
    "ts": "2026-05-31T14:03:22Z",
    "sig": "<HMAC over the above fields>"
  },
  "android": { "priority": "high" }
}
```
The `sig` lets the app reject spoofed pushes before it even calls the API (Section 7).

---

## 5. App screens

Three screens, deliberately minimal.

1. **Escalation list** — open items, sorted by severity then age. Each row: severity dot, category, title, age, deadline countdown if present. Pull-to-refresh; also refreshes on FCM receipt.
2. **Escalation detail** — title, summary, expandable context, the options as cards (showing pros/cons and the recommended badge), timeline impact, deadline countdown. Primary action area pinned to bottom.
3. **Decision confirmation** — appears when an option is tapped. For `requires_strong_auth` options, this is gated by biometric/device-credential. For `is_destructive`, an extra explicit confirm. Optional free-text note field if `free_text_allowed`.

A resolved/rescinded escalation opened from a stale notification shows a clear "already resolved" state rather than letting the user act on it.

---

## 6. Backend API (Escalation Service)

All endpoints require auth (Section 7). JSON over HTTPS, TLS 1.3.

```
POST /v1/devices/register
  body: { device_token, platform: "android", app_version }
  → binds the FCM token to the authenticated user; replaces old token.

DELETE /v1/devices/{device_token}
  → on logout / token rotation.

GET /v1/escalations?state=open
  → list view; returns lightweight summaries.

GET /v1/escalations/{id}
  → full escalation object (Section 3).

POST /v1/escalations/{id}/ack
  → marks acknowledged; stops reminder cadence. Idempotent.

POST /v1/escalations/{id}/decision
  body: {
    option_id,
    free_text?,
    client_decision_id,        // UUID for idempotency / replay protection
    auth_assertion             // strong-auth proof when required (Section 7)
  }
  → 200 {state:"decided"} | 409 if already decided/expired/rescinded.
```

Decision intake is **idempotent on `client_decision_id`** and the server enforces single-winner: the first valid decision wins; a later one returns `409` with the recorded outcome. This handles double-taps, retries, and two approvers acting at once.

Emily reads decided escalations from the service and applies them; the loop from Section 2 closes here.

---

## 7. Security & authorization

This is the part that matters most, because the "Approve" button can authorize an autonomous, self-directing agent to act.

- **User auth:** OIDC against your IdP; short-lived access token + refresh token in Android Keystore-backed encrypted storage. No long-lived bearer tokens on device.
- **Strong-auth for sensitive approvals:** options flagged `requires_strong_auth` require a fresh biometric / device-credential prompt (`BiometricPrompt`) at decision time. The result is bound into an `auth_assertion` the server verifies — not just a local boolean the app could lie about. Practically: the device signs the decision with a Keystore key that is gated by biometric (`setUserAuthenticationRequired(true)`), and the server verifies the signature against the device's registered public key.
- **Push authenticity:** the `sig` HMAC on the push payload is verified before the app trusts or displays it, so a forged push can't fabricate an escalation. Detail is always re-fetched from the API regardless.
- **Replay protection:** `client_decision_id` + server-side single-winner enforcement + short assertion validity window.
- **Transport:** TLS 1.3, plus certificate pinning to the Escalation Service. Consider mTLS for the Emily↔Service hop.
- **Least context on the wire:** sensitive escalation content never rides the push channel; only IDs and previews do.
- **Audit:** the service records, immutably, every escalation, every ack, every decision, who decided, the auth method used, and timestamps. This log is part of Emily's existing audit trail and must not be editable from the app.
- **Authorization scope:** approving here authorizes *only the specific option presented*. The app cannot issue free-form commands to Emily; it can only select among options Emily defined or send a note. This is a hard boundary — the phone is not a remote shell.

A blunt design rule: if an approval would let Emily do something irreversible (delete data, spend money, change production), the spec requires both `requires_strong_auth` and `is_destructive` to be set, forcing biometric + explicit confirm. Convenience never overrides that.

---

## 8. Reliability & delivery guarantees

Escalations — emergencies especially — can't quietly fail to arrive.

- **High-priority FCM** for `high`/`emergency` to survive Doze/background limits.
- **Acknowledgement loop:** the service expects an `/ack` within a per-severity window. No ack → re-push with backoff.
- **Multi-channel fallback for emergencies:** if no ack within N minutes on an `emergency`, fall back to SMS and/or an automated phone call via a third-party gateway, then to the next person in an on-call rotation.
- **Escalation routing / on-call:** the service maps category+severity to a responder (or rotation). Unacked criticals walk down the list rather than dead-ending on one silent phone.
- **Deadline handling:** if `deadline` passes with no decision, Emily applies its predefined safe default (usually: stay halted / take the conservative option) and the escalation is marked `expired` — never "guess and proceed."
- **Offline app:** queued decisions are held with their `client_decision_id` and synced on reconnect; the server's single-winner rule makes a late sync safe (it'll cleanly lose to any decision already recorded).

---

## 9. Tech stack

- **Language/UI:** Kotlin, Jetpack Compose.
- **Push:** Firebase Cloud Messaging (data messages).
- **Auth/crypto:** AndroidX Security, `BiometricPrompt`, Android Keystore (StrongBox where available).
- **Networking:** Retrofit/OkHttp with certificate pinning.
- **Local store:** encrypted DataStore/Room for the cached escalation list and queued decisions (so context survives backgrounding, but is wiped on logout).
- **Min SDK:** 26+ (notification channels); target current.
- **Backend (Escalation Service):** your choice; needs FCM Admin SDK, your IdP integration, an append-only audit store, and a small REST surface (Section 6).

---

## 10. Out of scope (v0.1)

No live agent chat, no log browsing, no metrics dashboards, no ability to *start* tasks — only respond to escalations Emily raises. No iOS (separate spec; the data model and API are platform-neutral and reusable). Keeping the surface small is itself a security property: the fewer things the phone can authorize, the less a lost/compromised phone can do.

---

## 11. Open questions for you

- Who's on the responder list per category, and is there a rotation, or single owners?
- For emergencies, do you want the SMS/phone-call fallback in v0.1, or is high-priority push enough to start?
- Should multiple people be able to see an escalation but only designated ones decide (view vs. act roles)?
- Do you have an existing IdP (Okta/Auth0/Google Workspace/etc.) to wire OIDC into, or does the service issue its own identities?

Answer those four and the spec is implementable as written.
