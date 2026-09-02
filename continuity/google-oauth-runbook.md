# Google OAuth Runbook — running file, updated every iteration

**How this works** (founder real-time, 2026-09-02): this file always names the ONE next concrete
step. Read the "NEXT STEP" section below, try it, then either:
- **it worked** → tell me, I mark it done here and write the next one, or
- **you got stuck / it didn't work / something's unclear** → tell me what happened, I update this
  file with a corrected/more specific step, and we go again.

An Apple gets filed every time this file changes, so there's a real audit trail of the back-and-
forth alongside the file itself.

---

## Why this is blocked on you specifically

Creating/configuring OAuth credentials happens in the Google Cloud Console, a web UI — there's no
API or CLI path a coding agent can drive for the consent-screen/credential-creation steps
themselves (this has been the real, standing blocker on Google sign-in all session:
`internal/http/handlers/portal.go`'s own header comment already says "Google sign-in stays
blocked on a human-only GCP Console step"). Everything on the CODE side (the actual sign-in
button, the backend token verification, the session cookie) is already real and built — see
"What's already built" below.

## What's already built, real, and just waiting on real credentials

- **`GoogleAuthHandler`** (`POST /api/v1/auth/google`, `internal/http/handlers/auth.go`) — verifies
  a real Google ID token (Google Identity Services client-side "Sign in with Google" button flow,
  not a server-side redirect), issues an IDUNA JWT + session cookie. Used by both `/admin/login`'s
  own eventual Google option and the developer portal's own `/portal/login`.
- **`WebCeremonyHandler`** (`internal/http/handlers/web_ceremony.go`) — a SEPARATE, real
  server-side "authorization code" OAuth flow (needs `GOOGLE_CLIENT_SECRET`, does a real token
  exchange), used for a different ceremony/onboarding flow. This one needs a real, exact
  **Authorized redirect URI** registered in Google Cloud Console (see Step 1 below) — the GIS
  button flow above does not.
- Both already read real env vars once set: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`,
  `GOOGLE_REDIRECT_URI` / `CEREMONY_OAUTH_REDIRECT_URI` (see `IDUNA/CLAUDE.md`'s own Key Env Vars
  table). None of these are currently set anywhere I can see from here.

## Real, already-existing lead worth knowing before Step 1

A GCP project already exists and is in real use for something else: `GOOGLE_CLOUD_PROJECT` is
already set to `project-d24a71e9-2daf-4b2d-917` in this environment (Vertex AI / Gemini related,
per the other `GOOGLE_CLOUD_*`/`GOOGLE_GENAI_*` vars alongside it). **You can very likely reuse
this same project** for the OAuth client instead of creating a new one — one less real decision to
make in Step 1, but confirm it's the right project when you're in the Console (a genAI-only
project may deliberately be kept separate from public-facing OAuth, only you know that context).

The real, live public domain for IDUNA, per an earlier session's own confirmed-via-curl note
(`EMILY/continuity/2026-07-18-concrete-next-steps.md`): `iduna.farthq.com`. **This is my own best
current guess for what to register in the Console below — confirm or correct it in your feedback**
if IDUNA's real current public URL is different now.

---

## NEXT STEP (do this one)

1. Go to **https://console.cloud.google.com/apis/credentials** (Google Cloud Console →
   APIs & Services → Credentials). Select the `project-d24a71e9-2daf-4b2d-917` project (or
   whichever real project you want IDUNA's OAuth client to live in) from the project picker at
   the top, if it isn't already selected.

2. If you haven't already configured an **OAuth consent screen** for this project, you'll be
   prompted to do that first (Console → APIs & Services → OAuth consent screen):
   - User type: **External** (unless every real IDUNA user is inside a Google Workspace org you
     control, in which case **Internal** is simpler — your call, either works with the code as
     built).
   - App name: something recognizable, e.g. `IDUNA` or `EINHORN_INDUSTRIAL IDUNA`.
   - Support email / developer contact: your own email.
   - Scopes: the defaults (`.../auth/userinfo.email`, `.../auth/userinfo.profile`, `openid`) are
     enough — the code only reads `sub` (Google's own stable user ID) and `email` from the
     verified token.
   - If **External** + "Testing" publishing status: add your own Google account (and anyone else
     who needs to sign in before this goes fully public) under **Test users**.

3. Create a new **OAuth 2.0 Client ID**: Credentials → **+ Create Credentials** → **OAuth client
   ID** → Application type: **Web application**.
   - Name: e.g. `IDUNA web`.
   - **Authorized JavaScript origins** (needed for the Sign-in-with-Google button,
     `GoogleAuthHandler`): add
     - `https://iduna.farthq.com` (confirm/correct this is the real current domain)
     - `http://localhost:8080` (for local dev/testing)
   - **Authorized redirect URIs** (needed for `WebCeremonyHandler`'s own separate flow): add
     - `https://iduna.farthq.com/`
     - `http://localhost:8080/`
   - Click **Create**. Google will show you a real **Client ID** and **Client Secret** —
     copy both somewhere safe for the next step (don't paste them back to me in chat; treat them
     like any other credential).

4. **Tell me once you've done this** (you don't need to hand me the actual secret values) —
   the next step will be setting `GOOGLE_CLIENT_ID`/`GOOGLE_CLIENT_SECRET` as real env vars
   wherever IDUNA's real production process actually runs (I don't have visibility into that
   host from this sandbox, so I'll need you to tell me how IDUNA's env is currently configured
   there — a systemd unit's `Environment=` lines, a `.env` file it sources, etc. — so the next
   step in this file is exact rather than generic).

**If you get stuck on any of the above** (can't find a menu, a field asks for something not
listed here, the domain guess is wrong, etc.) — tell me exactly where, and I'll correct this file.

---

## Log

- 2026-09-02 — file created, Step 1 (consent screen + OAuth client) written. Not yet attempted.
