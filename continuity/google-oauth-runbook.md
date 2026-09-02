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

~~The real, live public domain for IDUNA, per an earlier session's own confirmed-via-curl note
(`EMILY/continuity/2026-07-18-concrete-next-steps.md`): `iduna.farthq.com`.~~ **Corrected by the
founder**: the real domain actually used for `/portal/login` is **`okemily.com`** — the guess
above was wrong (stale from a month-old note, or IDUNA's own real public path moved since). See
the current NEXT STEP below: the OAuth client likely needs `https://okemily.com` added to its own
Authorized JavaScript origins / redirect URIs, since Step 1 only ever registered the wrong
`iduna.farthq.com` guess.

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

## Step 2 attempt — real error hit: "your Google account may not be registered, sign in failed"

This is the exact "Testing" publishing-status case named in Step 1's own predictions: the OAuth
consent screen is almost certainly still in **Testing** mode, which only allows sign-in from
Google accounts explicitly added as **Test users** — every other account (including yours, unless
already added) gets turned away with a message shaped like this one, before the button even gets
to a real permission-grant screen.

## Real domain correction — okemily.com, not iduna.farthq.com

Confirmed by the founder: `/portal/login` is really served from **`https://okemily.com`**. Step
1's own guess (`iduna.farthq.com`) was wrong, so the OAuth client currently has the WRONG origin/
redirect URI registered — this is very likely the real, direct cause of the `redirect_uri_mismatch`
error below, separate from (and in addition to) the Test-users issue.

## NEXT STEP (do this one) — two real fixes on the same OAuth client

**Confirmed so far**: `redirect_uri_mismatch` (this iteration) and, on an earlier attempt,
"your Google account may not be registered" (the Testing/Test-users case). Both are real, and
both need fixing — do both while you're in the Console rather than one round trip each.

1. Go to **https://console.cloud.google.com/apis/credentials** (project
   `project-d24a71e9-2daf-4b2d-917`), open the **OAuth client ID** created in Step 1, and add:
   - **Authorized JavaScript origins**: `https://okemily.com`
   - **Authorized redirect URIs**: `https://okemily.com/` **and** `https://okemily.com/portal/login`
     — I don't know which EXACT one Google's own error page is asking for without seeing it, so
     add both rather than guessing once and possibly missing again. If the error page you saw
     printed an exact `redirect_uri=...` value (Google's own real error pages usually show this
     in the fine print), please tell me that exact string too — it's the one, certain fix, no
     guessing needed.
   - You can leave the old, wrong `iduna.farthq.com` entries in place (harmless) or delete them —
     your call.
   - **Save.**
2. Same visit, also do the **Test users** fix from the earlier "may not be registered" error: OAuth
   consent screen → **Audience** tab → **Test users** → **+ Add users** → add the real Google
   account you're testing with → **Save**.
3. **Try the real sign-in again** at `https://okemily.com/portal/login`, same account.
   - **Works** → tell me, done.
   - **Different error** → tell me exactly what it says.
   - **Same error** → tell me the exact `redirect_uri=...` value from Google's own error page
     text if you can find it; that's the one real, unambiguous fix.

---

## Do we have logging for this failure? — real, honest answer: no, not this specific one

This particular failure (Google rejecting the sign-in on Google's OWN side — wrong redirect URI,
or account not in Test users) happens entirely between the browser and `accounts.google.com`,
**before** any request ever reaches IDUNA. `GoogleAuthHandler` (the code that emits
`iduna:auth.google.failure` into the unified log, S226-02) only runs once Google hands the
browser a real ID token to POST to `/api/v1/auth/google` — a rejection at this earlier stage never
gets that far, so there's genuinely nothing for IDUNA's own logging to see or record here. This
isn't a gap in the logging work itself, just the real limitation of a client-side OAuth flow: the
config-mismatch/consent-screen failures are only ever visible in the browser (and, for anyone with
real Console access, in the OAuth client's own "Recent activity" / Google Cloud's own audit
surfaces — not this app's own log). A real, separate, later idea if this keeps needing debugging:
have the login page's own JS call `POST /services/collector` directly on a GIS error callback, so
even these "never reached IDUNA" failures land in the same place — not built now, just named.

---

## Log

- 2026-09-02 — file created, Step 1 (consent screen + OAuth client) written. Not yet attempted.
- 2026-09-02 — founder created the OAuth client, gave me the real Client ID + Secret directly in
  chat. Found IDUNA's real running env file (`~/.config/iduna/env`, a live `iduna.service`
  systemd user unit), added both values there, rebuilt the binary, restarted the service, and
  verified live (health check green, the real Client ID confirmed reaching `/portal/login`'s own
  server-rendered HTML). Next: a real end-to-end sign-in attempt in an actual browser.
- 2026-09-02 — founder tried the real sign-in at the real domain, `okemily.com/portal/login` (not
  `iduna.farthq.com` — Step 1's own domain guess was wrong, corrected here). Two real errors seen
  across attempts: "your Google account may not be registered" (Testing/Test-users case) and
  `redirect_uri_mismatch` (the direct, likely-primary consequence of the wrong domain guess).
  Founder also asked whether this failure is captured in the new unified logging backend —
  answered: no, honestly, since a Google-side rejection never reaches IDUNA's own code at all.
  Next: add the real `okemily.com` origin/redirect URIs to the OAuth client, and add the real
  test account, in the same Console visit.
