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

## Step 1 — DONE

You created the OAuth consent screen + a real Web application Client ID and handed me the
Client ID + Client Secret directly. Real, done, verified:

- Found the real place these belong: IDUNA runs here as a real, live `systemctl --user` service
  (`iduna.service`, confirmed running, PID checked directly), which loads its env from
  `~/.config/iduna/env` (an `EnvironmentFile=`, per the unit file's own header comment) —
  **not** a file this session had visibility into before you gave me the credential, so this
  answers the "where does IDUNA's real env actually live" question the previous version of this
  file left open.
- Added `GOOGLE_CLIENT_ID`/`GOOGLE_CLIENT_SECRET` to that file (values not repeated here or in
  git — this file and its own git history never carry the actual secret).
- Rebuilt the binary (`go build -o ~/.local/bin/iduna .`, matching the unit file's own documented
  deploy step) and restarted the service (`systemctl --user restart iduna.service`) so the new
  env vars actually take effect (systemd doesn't hot-reload `EnvironmentFile` into an already-
  running process).
- **Verified live**: `curl http://localhost:8080/health` → healthy; `curl .../portal/login`'s own
  real HTML now renders `data-google-client-id="532442865445-...apps.googleusercontent.com"` —
  the real Client ID is reaching the page (checked the actual server-rendered attribute, not the
  page's own always-present client-side JS fallback STRING, which would have shown up in a naive
  text search regardless of whether it's really unconfigured — a real, worth-naming gotcha for
  next time this needs re-checking).

## NEXT STEP (do this one)

**Try a real, end-to-end sign-in** in an actual browser: go to IDUNA's own real, public
`/portal/login` page (you know the real current public URL for this — my own only confirmed lead
this session was `https://iduna.farthq.com`, from a month-old note, so please use whatever you
know to be current) and click the real "Sign in with Google" button.

- **If it works** (you land back on `/portal` signed in) — tell me, and we're done here; I'll
  close this out.
- **If Google shows an error page instead of completing sign-in** — the two most likely real
  causes, given Step 1's own guesses about your real domain:
  - `redirect_uri_mismatch` — the real redirect URI Google received doesn't match what's
    registered on the OAuth client. Tell me the exact URI Google's own error page shows; I'll add
    it to the Client ID's own **Authorized redirect URIs** list in the Console (I can't do this
    part myself — it's the same GCP Console step as Step 1).
  - `origin_mismatch` / "not a valid origin for the client" — same real fix, but for **Authorized
    JavaScript origins** instead.
  - "Access blocked: this app has not completed Google's verification process" / "app isn't
    verified" — if the consent screen is in **Testing** mode, your own Google account needs to be
    added under **Test users** (OAuth consent screen page) before you can sign in with it.
- **If the button doesn't even render** (blank space where it should be, or the "not yet
  configured" fallback text really shows up on the actual rendered page, not just present
  somewhere in page source) — that's a real, different bug (the server-side check above already
  confirmed the value IS reaching the page, so this would point at something client-side) — tell
  me exactly what you see.

---

## Log

- 2026-09-02 — file created, Step 1 (consent screen + OAuth client) written. Not yet attempted.
- 2026-09-02 — founder created the OAuth client, gave me the real Client ID + Secret directly in
  chat. Found IDUNA's real running env file (`~/.config/iduna/env`, a live `iduna.service`
  systemd user unit), added both values there, rebuilt the binary, restarted the service, and
  verified live (health check green, the real Client ID confirmed reaching `/portal/login`'s own
  server-rendered HTML). Next: a real end-to-end sign-in attempt in an actual browser.
