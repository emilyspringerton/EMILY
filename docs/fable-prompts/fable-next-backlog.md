# Fable Next — queued dispatch tracker

*Index of prompts written for future Fable dispatch. Each entry links to its full prompt file in
this directory. Status values: **queued** (written, not sent), **blocked** (queued but has an
explicit dependency that hasn't landed yet), **dispatched** (sent, awaiting completion),
**done** (landed — move the entry to the Done section with the commit that closed it).*

*Convention: write the full prompt as its own file in this directory (see
`iduna-front-door-funnel.md` for the shape — status/dependency header, then the dispatch-ready
prompt body), then add a one-line entry here. Update the entry's status as it moves. Don't let
prompt content live only in this index — this file is a table of contents, not the prompts
themselves.*

---

## Queued

| # | Title | File | Status | Depends on |
|---|---|---|---|---|
| 3 | TYLER S00E-2 (prequel bridge to S00E01) | `tyler-s00e-2-retry.md` | **ready to dispatch** | none — still hold until #14 (Season 12, failed 2026-07-19, needs redispatch below) actually lands, to avoid two concurrent agents both writing TYLER canon |
| 15 | Evolve the OKEMILY TYLER easter egg into an ARG using Field Activation (YOLO, one-shot) | `tyler-arg-field-activation.md` | **ready to dispatch** | none — same hold as #3 |
| 4 | Act III synthesis — reconcile 141-150 + KIKORYU audit, apply SAGA/NORN thinking by hand, write tournaments-platform backlog | `act3-synthesis.md` | **ready to dispatch, paused** | none — failed twice on transient `529 Overloaded` (not session-limit, not content). Per founder instruction 2026-07-16, not being auto-retried; dispatch deliberately whenever picked back up. |

---

## Dispatched (awaiting completion)

| # | Title | Dispatched | Notes |
|---|---|---|---|
| 5 | DNS operations northstar ("fates dns" — EINHORN operating its own DNS) | 2026-07-17 | Full standalone prompt written (`dns-operations-northstar.md`) after 2 prior attempts failed on transient 529s with no saved prompt; dispatched with real grounding (farthq.com's current Cloudflare NS confirmed via `dig`, real product-surface inventory, NORN/SAGA cross-references). |
| 8 | OKEMILY blog post — "Clean Builds First" | published directly 2026-07-19 | https://okemily.com/blog/clean-builds-first/ — no Fable dispatch needed, Claude Code published it directly. |
| 12 | OKEMILY blog post — "Recursion for LLMs" | published directly 2026-07-19 | https://okemily.com/blog/recursion-for-llms/ |

---

## Failed — needs redispatch (2026-07-19, session-limit batch)

*All 7 dispatched together, 2026-07-19 ~03:00 UTC — mistake: fired in one parallel batch, which
burned through the whole session's Fable usage limit immediately. Every one failed within seconds
on "You've hit your session limit · resets 5:40am (UTC)," before any of them produced real work
(most were still in required-reading). No file changes were made by any of them (confirmed via
`git worktree list` across every affected repo — all clean), so nothing needs cleanup, only
redispatch. Do not batch-redispatch these — one or two at a time, after 5:40am UTC, so a real
session-limit signal is distinguishable from a repeat of this mistake.*

| # | Title | File |
|---|---|---|
| 1 | IDUNA front door funnel for Agents and Unagents | `iduna-front-door-funnel.md` |
| 7 | AGI RSI next steps: human-in-the-loop tiers + revenue via EDIS | `agi-rsi-hitl-revenue-edis.md` |
| 9 | Augment/improve OKEMILY landing page + signup funnel | `okemily-landing-page-augment.md` |
| 10 | SHANKPIT funnel page on okemily.com (Phase 1) | `okemily-shankpit-funnel.md` |
| 13 | "The Gauntlet" northstar | `gauntlet-press-release-publishing.md` |
| 14 | TYLER Season 12 | `tyler-season-12.md` |

---

## Failed — needs redispatch

| # | Title | File | Notes |
|---|---|---|---|
| 3 | TYLER S00E-2 (prequel bridge to S00E01) | `tyler-s00e-2-retry.md` | Hit a Fable session/usage limit before writing anything (verified clean tree). Resets ~5:50am UTC 2026-07-17. Prompt file is written and ready — just dispatch it after the reset. |

---

## Done

| # | Title | Landed | Commits |
|---|---|---|---|
| 6 | **TOP PRIORITY** — entity-graph/signalapi/newssite full-in-memory-replay fragility fix → `PRRJECT_FATBABY/docs/northstar/replay-fragility.md` (decision: streaming eventstore `Scan` API + per-process SQLite snapshot-plus-tail checkpoints; root cause was O(n²) full-file re-reads in `FileStore.ReadFrom`) | 2026-07-18 | PRRJECT_FATBABY `a628947` (EMILY golden-index + backlog commit follows this file). Apple #9986 |
| 2 | Emilyify the mag book → `gpt2-alpine-c/pkg/towerprint` + `docs/TOWERPRINT.md` | 2026-07-16 | gpt2-alpine-c `819c400`, EMILY `c0efd2f` |

---

## Notes / related non-Fable prep work

- `IDUNA/ops/nginx-front-door-snippet.conf` — ready-to-apply nginx location blocks for IDUNA's
  API/auth/admin/jwks paths on `iduna.farthq.com` (sudo-gated, not yet applied). Written as prep
  for entry #1, not itself a Fable task — pure ops mechanics.
- Entry #1 was updated 2026-07-16 with the real deployment conflict it needs to resolve
  (`iduna.farthq.com` is currently 100% EDIS WordPress; IDUNA's own static frontend collides with
  WordPress's root path) and the founder's leading resolution idea (host the funnel UI as a new
  `edis-iduna` WordPress plugin inside EDIS — Emily *Distributed Intelligence System* — rather than
  serving IDUNA's own competing frontend).
