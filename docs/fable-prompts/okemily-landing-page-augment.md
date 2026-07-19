# Queued Fable prompt — augment/improve the OKEMILY landing page + signup funnel

**Status:** queued, not yet dispatched. Written 2026-07-19 (post-reboot session).
**Dependency:** none to start, but see the mailing-list note below — verify the signup funnel is
actually accepting subscribers before or as part of this work, not after.

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## Why this is queued, not just handed to a general edit

The founder's ask was specifically to use the **golden docs** as grounding and improve the page
from there — not a blind copywriting pass. The risk with landing-page copy on a company this fast-
moving is claims going stale or overclaiming (saying something is live that isn't — see the EDIS
deploy status found earlier this session, entry #7 in this same index). Fable reading the golden
doc set first is the guard against that; it's the reason this is a dispatch rather than a five-
minute inline edit.

## Reality check (2026-07-19, re-verify if dispatching later)

- The current live page (`OKEMILY/index.html`) is deliberately minimal — see `OKEMILY/CLAUDE.md`:
  "**Not the product funnel.** EDIS (WordPress, Ask Emily chat, signal widgets) is the eventual
  full product front-end and is a separate, later effort." Built narrow on purpose, for Google
  Cloud for Startups credibility, not as the real conversion funnel.
- The only actual conversion action on the page today is the mailing-list signup form
  (`POST /api/v1/mailing-list/subscribe`, same-origin proxied to IDUNA).
- **The mailing-list vault was found locked earlier in this session** (blocking every signup
  attempt with a silent failure since IDUNA's last restart) — the founder was given the unlock
  command to run interactively. Confirm it's actually unlocked (`curl -s
  http://127.0.0.1:8080/api/v1/mailing-list/subscribe -X POST -H 'Content-Type: application/json'
  -d '{"email":"test@example.com","consent":true}'` should not error with a locked-vault message)
  before treating the funnel as working — there's no point improving copy that feeds a broken pipe.
- **EDIS production has never actually deployed** (`iduna.farthq.com` connection-refuses as of
  2026-07-19, S23-01 still unchecked in `EMILY/BACKLOG.md`). Do not write copy that implies Ask
  Emily, signal widgets, or any EDIS-hosted product is live — it isn't yet.
- A separate, already-queued Fable prompt (`iduna-front-door-funnel.md`, entry #1) covers IDUNA's
  own identity/onboarding funnel (Google OAuth → honor code → gamertag). That's a different,
  deeper piece of work with its own security framing (dedicated clean origin, no WordPress in the
  path). This prompt is scoped to the *marketing* landing page and its signup funnel only — don't
  let scope bleed into that design question; reference it, don't redo it.
- Live deploy of any change here requires a human running `sudo rsync ... && sudo systemctl reload
  nginx` per `OKEMILY/CLAUDE.md` — and that doc has a hard-won warning about certbot rewriting the
  live nginx config in place (a blind overwrite caused a real HTTPS outage on 2026-07-18). Your
  deliverable is the improved repo files, committed and pushed. **Do not attempt the sudo deploy
  steps yourself** — say clearly in your final summary that a human needs to run the deploy, and
  give the exact command from `OKEMILY/CLAUDE.md`.

## Prompt body (copy from here down)

You're improving EINHORN_INDUSTRIAL's public landing page (`OKEMILY/index.html`, plain static
HTML/CSS/JS, no build step, no framework — keep it that way) and its signup funnel. The founder
loves the page as it stands and wants it **augmented**, not replaced — preserve its voice, its
dark/light theme system, the easter eggs (the triple-click year reveal, the TYLER footer quote,
"Enter the void 🕳️"). This is refinement, not a redesign.

### Required reading first — the golden docs, for real

Read `EMILY/context/golden-docs-index.md` and pull the current version of every doc listed there
at tier 1, at minimum. These specifically matter for landing-page accuracy:

- `EMILY/docs/NORTHSTAR.md` (EMILY-NORTH) and `EMILY/emily-memory/world-state.md` (EMILY-MEMORY)
  — current, real status of every product line. The page's "three pillars" section must not claim
  more maturity than what these say is actually true today.
- `PRRJECT_FATBABY/docs/northstar/executive_summary.md` (FATBABY-EXEC) — for the capital-markets
  pillar's copy.
- `SHANKPIT/docs2/NORTHSTAR.md` (SHANKPIT) — for the game-worlds pillar, including real Steam
  Early Access track status.
- `PRRJECT_FATBABY/docs/GTM_FUNNEL.md` (GTM) — the actual Ask Emily funnel design (free tier →
  subscription). This is the closest existing doc to "what should our funnel eventually look
  like," even though Ask Emily itself lives on EDIS, not here — use it to inform what a *waitlist*
  page should be promising and collecting interest for.
- `EDIS/NORTHSTAR.md` (EDIS) — cross-check against the reality check above; the doc will describe
  EDIS's intended end state, which is not the same as what's live today. Don't conflate the two in
  copy.
- `EMILY/docs/THE_EMILY_WAY.md` (THE-EMILY-WAY) — the values section on the page ("build in the
  open," "real data over demos," etc.) should stay true to this, not drift into generic startup
  copy.
- `TYLER/README.md` (TYLER-BIBLE) — for the tone of the existing fiction easter egg; if you extend
  or add to that thread, match its register (see the OKEMILY blog's "Activation #114" post,
  `/var/www/okemily/blog/activation-114/index.html`, for how the company already blends its real
  and fictional layers in public-facing writing).

### What "augment and improve" means here, concretely

Your call on specifics, but the plan should cover:

1. **Accuracy pass.** Walk every factual claim on the current page against the golden docs above.
   Fix anything stale or aspirational-sounding-as-fact.
2. **Funnel strength.** Right now there's exactly one CTA (email signup) plus a footer link to
   "the void." Is that the right single funnel, or does the page need a second, softer CTA (e.g.
   "read the blog," "browse the code") for visitors not ready to hand over an email? Look at what
   analogous credibility/waitlist pages do well, but don't cargo-cult generic SaaS landing-page
   patterns onto a page whose whole identity is "not a generic SaaS landing page."
3. **Proof points.** The page currently asserts things ("real production pipeline," "not a demo")
   without much concrete evidence. Consider surfacing something real and specific — a live stat,
   a recent Apple, a link to a specific blog post — that backs the claim rather than just stating
   it. Ground any number you add in something you actually queried (IDUNA `/api/v1/apples`,
   `emily status`), not a guess.
4. **Structure/copy only** — do not touch the signup form's JS/endpoint wiring, the privacy
   policy link, or the consent checkbox; that's a compliance surface, leave it alone unless you
   find it's actually broken (see the vault-lock check above), in which case flag it rather than
   silently fixing IDUNA-side code from this dispatch.

### Where this lands

Edit `OKEMILY/index.html` in place. Update `OKEMILY/CHANGELOG.md` per its existing convention.
Commit and push OKEMILY. In your final summary, state plainly: (a) what changed and why, citing
which golden doc grounded each factual claim you touched, (b) the exact deploy command from
`OKEMILY/CLAUDE.md` for the founder to run, and (c) whether the mailing-list vault check above
passed or failed.

---

## Prompt body ends here
