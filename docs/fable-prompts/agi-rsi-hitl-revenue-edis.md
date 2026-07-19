# Queued Fable prompt — AGI RSI next steps: human-in-the-loop tiers + revenue via EDIS

**Status:** queued, not yet dispatched. Written 2026-07-19 (post-reboot session).
**Dependency:** none — dispatch whenever picked up. Grounded in a live state check done the same
day (see "Reality check" below); re-verify anything time-sensitive if dispatching later than a
few days out.

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## Reality check (2026-07-19, do not skip re-verifying this)

The founder asked for "a full Fable planning step for next steps: AGI RSI, human-in-the-loop,
revenue, EDIS — let's go." Before writing the prompt body, a live check turned up a genuine
blocker that has to be the spine of this plan, not a footnote:

- **Correction (2026-07-19, later same session):** the first pass of this reality check was wrong
  — EDIS/WordPress **is** actually deployed and live. `/var/www/edis/wp-config.php` exists,
  `wp-content/plugins/` has edis-core, edis-ask-emily, edis-dis, edis-earnings, edis-signals, and
  akismet all installed, MySQL is running, and `curl http://iduna.farthq.com/` returns a real
  200 with `<title>IDUNA Intelligence Platform</title>`. `EMILY/BACKLOG.md` `S23-01` being
  unchecked is a stale checkbox, not reality — verify the real state directly (`curl`, `ls
  /var/www/edis/`) rather than trusting that line item. **What's actually still missing: HTTPS.**
  `https://iduna.farthq.com/` connection-refuses (`HTTP 000`) because certbot was never run for
  that domain — `/etc/nginx/sites-available/edis` is explicitly still the "HTTP only until certbot
  runs" bootstrap config. That's the real, narrow blocker: get a cert issued for iduna.farthq.com
  (or resolve the domain question below first, since which domain gets the cert is exactly the
  open question), not a from-scratch deploy.
- **Separately, and now the more load-bearing open question:** EDIS is live on
  **`iduna.farthq.com`**, a different domain from **`okemily.com`** (the polished, HTTPS-working,
  actively-used landing page + blog + mailing list). The founder asked, later in this same session,
  to "launch EDIS merge with okemily.com" — i.e. connect or consolidate these two. That's a real
  architecture decision (which domain serves WordPress at root, what happens to okemily.com's
  existing static landing page and IDUNA's own planned static frontend at bare `/`, how WP_HOME/
  WP_SITEURL migration is handled for a live WordPress install, cert consolidation) with a
  documented outage precedent for exactly this class of nginx surgery (see `OKEMILY/CLAUDE.md`'s
  2026-07-18 incident writeup). This is squarely the same collision entry #1
  (`iduna-front-door-funnel.md`) already exists to resolve — fold the "merge with okemily.com"
  question into that design pass rather than treating it as a separate ad hoc decision. Whoever
  picks this up should read entry #1 first and decide whether it still needs its own dispatch or
  whether this doc's revenue plan supersedes/absorbs it.
- **New, separate ask from the founder (2026-07-19):** EDIS/WordPress logins should go through
  IDUNA for SSO rather than native `wp-login.php`. This is real, scoped engineering (an IDUNA-auth
  WordPress plugin/must-use plugin doing an OAuth-style handoff) independent of the domain-merge
  question above — it can be built and committed without touching live production, and activated
  whenever the domain question is settled. Check `IDUNA/internal/http/handlers/local_auth.go` and
  the existing Google-OAuth consumer flow (`internal/http/handlers/auth.go`) as the starting point
  for what IDUNA already exposes versus what a WordPress-facing SSO handoff would need net-new.
- Everything *behind* that front door is comparatively mature: Ask Emily's chat endpoint, rate
  limiting, IDUNA OAuth, and MySQL/MongoDB read models are all done (SECTION 20/21, all `[x]`).
  The revenue mechanism is built; the door to it is closed.
- The human-in-the-loop primitive for the AGI RSI loop already exists in spec form —
  `EMILY/docs/hq-specs/HQ-SPEC-PRIME-101-norn-loop-kernel.md` (NORN) — as **gate approval tiers**:
  `autonomous` (no human), `prime_ack` (Emily Prime acknowledges before promotion), `biometric`
  (companion-app human approval, reserved for artifacts feeding a physical/financial execution
  path). This is the mechanism to build "human-in-the-loop" *on top of*, not a new concept to
  invent from scratch.

This prompt asks Fable to turn those three facts into a concrete, sequenced plan — not to
re-discover them.

---

## Prompt body (copy from here down)

You're doing a planning pass, not an implementation pass. Deliverable: a written plan (new doc,
see "Where this lands" below), not code changes — though it's fine to note exact files/functions
future implementation work will touch.

### The question

EINHORN_INDUSTRIAL's AGI RSI loop (Emily Prime's 5-minute OBSERVE→DECIDE→ACT→PLAN cron cycle,
`EMILY/emily-agent`) has been running unattended for weeks. The founder wants the next phase to
add **deliberate human-in-the-loop checkpoints** to that loop, and to connect it to **revenue**
via **EDIS** (the WordPress product surface). Plan the concrete next steps that get from where the
system is today to that state — sequenced, not a wishlist.

### Required reading first

1. `EMILY/docs/hq-specs/HQ-SPEC-PRIME-101-norn-loop-kernel.md` — NORN, the loop kernel. Section 4
   ("Gate Policy & Approval Tiers") is the actual human-in-the-loop mechanism: `autonomous` /
   `prime_ack` / `biometric` tiers, tier relaxation as itself a gated promotion. Section 6 (the
   instantiation table) shows what's already wired to which tier.
2. `EMILY/docs/THE_EMILY_WAY.md` — Principle 11 (AGI loop / `--continue` rationale) and whatever
   principle numbers now cover continuity reporting and memory-constraint doctrine (check current
   numbering, it has grown past what's cited in older docs).
3. `EMILY/BACKLOG.md` SECTION 19–24 (the three revenue tracks: SHANKPIT→Steam, Ask Emily, data
   licensing) and SECTION 35 (EDIS DIS production hardening) — read what's `[x]` vs `[ ]` for
   real, don't trust any summary of it that predates this dispatch.
4. `EDIS/NORTHSTAR.md`, `EDIS/docs/architecture.md`, `EDIS/CLAUDE.md` — what EDIS actually is
   (WordPress front end over signalapi + Emily Prime, three plugins: edis-core, edis-signals,
   edis-ask-emily) and where WooCommerce/Ask Emily paid-tier logic already lives.
5. `EDIS/ops/sprint-deploy.sh` — read it, don't just cite that it exists. What does it actually do,
   what does it assume about the target host, and what's the smallest safe first real run look
   like (dry-run flag? staging target? or is production the only target that exists)?
6. `EMILY/docs/fable-prompts/iduna-front-door-funnel.md` — a related, already-queued Fable prompt
   about IDUNA's front-door funnel UI colliding with EDIS's WordPress root path. Don't duplicate
   this work; your plan should name where it depends on or should sequence around entry #1's
   outcome (dispatch order matters if both eventually run).
7. `git -C /home/fatbaby/EMILY log --oneline -20` and the same for IDUNA, PRRJECT_FATBABY, EDIS —
   confirm nothing above has moved since 2026-07-19 before you plan against it as current state.

### What the plan needs to cover

**A. Unblock revenue (EDIS is live, but incomplete and on the wrong-for-branding domain).**
EDIS/WordPress already runs at `iduna.farthq.com` over HTTP — this is not a from-scratch deploy,
it's finishing one. Two real "someone runs sudo on a production host" actions remain, and your job
is to sequence them, not perform them: (1) get HTTPS issued for whichever domain wins the domain
question below, (2) resolve the `okemily.com` merge the founder asked for, which is the same
root-`/` collision entry #1 already scoped (WordPress wants `/`, IDUNA's own planned static
frontend wants `/`, and now okemily.com's existing static landing page is a third claimant on the
same question). Say plainly which of the two actions blocks the other — cert issuance depends on
knowing the final domain, so the domain decision almost certainly has to come first. Don't let this
plan re-litigate entry #1 from scratch; read it, decide whether it still needs its own dispatch or
whether your plan settles it directly, and say which.

**B. Human-in-the-loop, concretely.** NORN's gate tiers already exist as a spec. What RSI-loop
decisions in `emily-agent` today run fully `autonomous` that should move to `prime_ack` or
`biometric`? Name specific decision points in the cron cycle (task promotion, backlog curation,
CEO escalation triggers, HEIMDAL sprint translation) and propose which tier each belongs at and
why — don't propose "add human review everywhere," that's not a plan, that's a slogan. Consider
MJOLNIR (the Android companion app) as the natural `biometric`-tier approval surface — check
whether it already has any UI hook for this or whether that's new work.

**C. Sequencing.** Order the above two tracks plus whatever quick wins fall out of the required
reading, given real dependencies (e.g., does human-in-the-loop tooling need to exist before EDIS
goes live, or can they land independently and in parallel?). Flag anything that's a decision only
the founder can make (the sudo deploy action itself; any spend; anything touching real payment
flows given WooCommerce is already scaffolded per S135-03) versus what's safe for an agent to just
execute.

**D. Known gaps to fold in, don't relitigate:** `ANTHROPIC_API_KEY` credit balance hit zero during
this session (2026-07-19), pausing Claude invocations in the RSI loop for periods at a time —
note this as an operational risk for any human-in-the-loop design that assumes the loop is always
live to *ask* for approval, but don't spend the plan solving billing.

### Where this lands

Write the plan to `EMILY/docs/planning/agi-rsi-hitl-revenue-edis-plan.md` (new file — create the
`planning/` directory if it doesn't exist). Register it in
`EMILY/context/golden-docs-index.md` per the standard golden-doc protocol if you judge it
load-bearing enough for Emily Prime's own cron cycle to read (your call — argue it in one line if
you skip registration). Add concrete `[ ]` backlog items to `EMILY/BACKLOG.md` under a new
section for whatever this plan concludes should happen first, don't leave the plan un-actioned in
prose only. Commit EDIS/EMILY changes per the normal Apple + CHANGELOG + commit protocol.

---

## Prompt body ends here
