# Concrete Next Steps — 2026-07-18 (end of session)

*A punch list, not a narrative — what's actually done, what's blocked on you, and what's
genuinely next, for every deliverable touched tonight.*

---

## 1. okemily.com — done, stable, needs nothing right now

Landing page, privacy policy, encrypted mailing-list signup (Mailchimp live-verified), blog
(8 real posts), API playground, status page, HTTPS+HSTS. Nothing pending here unless you want
more content or features.

## 2. FatBaby pipeline health — the single most important open item

**This is the real next action, already queued and ready:**

- **Tomorrow, first thing**: dispatch `EMILY/docs/fable-prompts/replay-fragility-northstar.md` to
  Fable. It's pinned at the top of `BACKLOG.md` SECTION 1 specifically so saying "ok"/"go"/
  "continue" picks it up automatically. This fixes the real root cause behind tonight's
  `signalapi` incident and the ongoing `newssite` crash-loop: full event-history replay on every
  process restart, with no persisted/cached index. Recommended shape (already in the prompt):
  snapshot-plus-tail via SQLite, matching the pattern this session already used for
  `IDUNA/internal/blog` and `IDUNA/internal/mailinglist` — not Mongo, not primarily "a bigger box."
- **`signalapi` stays stopped/disabled** until that fix lands — it will not come back on its own,
  and re-enabling it before the fix would likely just repeat tonight's thrash.
- **`entity-graph`** was deliberately left as an unsupervised `go run` process (not migrated to
  systemd) given what happened to `signalapi` — same reasoning applies, wait for the real fix.
- **`newssite`** is still OOM-crash-looping periodically. I fixed one real, verified bug tonight
  (article deep-links were doing a full-history scan on every request — 60s+ → 1.85s for
  already-indexed docs), but that doesn't stop the crash loop itself, just makes each surviving
  window more useful. Same underlying fix (above) resolves this too.

## 3. EDIS (the full product front-end) — one command away

WordPress core is already live at `iduna.farthq.com` (confirmed via curl this session), just
untitled/unbranded and HTTP-only. `EDIS/ops/sprint-deploy.sh` is idempotent and, on inspection,
appears to have no missing secrets or hard blockers — it just needs to be run:

```
sudo bash /home/fatbaby/EDIS/ops/sprint-deploy.sh
```

After that: fix the WordPress site title (currently "IDUNA Intelligence Platform" — should be
something FatBaby/EDIS-branded), then decide whether to link EDIS from okemily.com's footer
(low-risk) or eventually move it to an okemily.com subdomain (bigger, deliberate move, not tonight).

## 4. SHANKPIT-460 — forked, running, needs a real design pass

Game server is live and systemd-supervised (`shankpit460-server.service`, UDP `:6969`,
auto-restart verified). **What's genuinely next**: an actual NORTHSTAR.md for this repo — what
specifically gets stripped down for the low-spec/esports target, what stays from the parent
SHANKPIT codebase. Nothing was improvised here tonight on purpose; this is real design work still
waiting on you or a Fable dispatch, whichever you'd rather do.

## 5. Two new northstars written tonight — neither implemented yet, both ready to build from

- `PRRJECT_FATBABY/docs/headlines/live-feed-northstar.md` — MTWire-style combined headline feed,
  builds on the existing (not-yet-running) `feedserver`.
- `PRRJECT_FATBABY/docs/northstar/tina-engine.md` — TINA (Trading Idea, Not Advice), the
  compliance-forward layer under `jon-agent`'s existing persona.

Both are design-only by intent (Emily Way: spec before implementation). Pick whichever matters
more when there's bandwidth to actually build one.

## 6. Smaller, real, still-open items

- IDUNA's uncommitted `202606180001_local_users.sql` diff — pre-existing, not from tonight, still
  needs a decision (commit as a proper new migration, or discard).
- `EMILY/BACKLOG.md` S153-06 — board of directors / non-profit ownership structure. Explicitly
  parked by you, no urgency, existing docs worth a real pass whenever it's time
  (`emily-complete-vision.md`, `emily-press-package.md`, `THE_FIELD.md`, and others listed in
  BACKLOG).
- S153-09 — IDUNA's OpenAPI spec is known-stale (missing the blog/mailing-list endpoints added
  tonight, plus a second unreconciled `openapi.yaml`). Deferred by your own explicit instruction.
- Kubernetes/k3s migration northstar — still owed from earlier this session, still not written.
  Worth doing before the Mercury-gated infra window opens (2026-07-23 onward) so there's a plan
  ready rather than starting cold.
- Swap-space expansion (`fallocate` a swapfile) — proposed earlier, never run, still needs your
  sudo password if you want the extra buffer. Lower priority now that the replay-fragility fix is
  queued as the real solution to tonight's actual memory pressure.

## 7. What genuinely doesn't need anything from you right now

Everything in §1 (okemily.com), the `secwatch`/`processor`/`prwatch`/`prwatch-body`/
`eps-reconciler`/`newssite`/`shankpit460-server` systemd units (all auto-restart-verified), and
both new northstar docs (they're complete as design documents; building them is the open item, not
finishing writing them).
