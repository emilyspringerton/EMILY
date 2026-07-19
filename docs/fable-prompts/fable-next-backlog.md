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
| 1 | IDUNA front door funnel for Agents and Unagents | `iduna-front-door-funnel.md` | **ready to dispatch** | SAGA VS-spec audit landed 2026-07-16 (IDUNA `529834f`, EMILY `508a9be`) — no remaining blocker |
| 7 | AGI RSI next steps: human-in-the-loop tiers + revenue via EDIS | `agi-rsi-hitl-revenue-edis.md` | **ready to dispatch** | none — corrected 2026-07-19: EDIS/WordPress is actually already live on `iduna.farthq.com` (HTTP only, no cert); real blocker is (a) HTTPS + (b) the founder's "merge with okemily.com" ask, which folds into entry #1's already-scoped root-`/` collision |
| 8 | OKEMILY blog post — "Clean Builds First" (fiction substrate ↔ AGI RSI loop) | `okemily-blog-clean-builds-first-DRAFT.md` | **ready to dispatch** — full pre-pass draft already written by Claude Code, Fable's job is final line-edit + fact re-verify (repo state moves fast) + publish via IDUNA `blog.write` | none |
| 9 | Augment/improve OKEMILY landing page + signup funnel, grounded in golden docs | `okemily-landing-page-augment.md` | **ready to dispatch** | none — mailing-list vault unlocked by founder 2026-07-19; deploy itself needs a human (`sudo rsync` + nginx reload per `OKEMILY/CLAUDE.md`) |
| 10 | SHANKPIT funnel page on okemily.com (Phase 1 — CTA placeholder to GitHub release tag) | `okemily-shankpit-funnel.md` | **ready to dispatch** | none — founder-approved phased plan 2026-07-19; farthq.com stays up unchanged; CTA target verified live (`github.com/emilyspringerton/SHANKPIT/releases/tag/460`) |
| 11 | SHANKPIT-460 maps: build v0 map (bases + mid rocket tower, CTF-shaped/DM-launch, Halo-ish "Super Rumble" physics feel) + write longer-term pipeline report (AI-gen, editor, bot-fuzz-eval via NORN) | `shankpit460-maps-report.md` | **ready to dispatch** | none — current map system is 142 lines, one hardcoded box layout compiled into the binary; v0 map has a concrete founder brief as of 2026-07-19 |
| 12 | OKEMILY blog post — "Recursion for LLMs" (HEIMDAL outage as an LLM stepping into a broken LLM-translation pipeline) | `okemily-blog-recursion-for-llms-DRAFT.md` | **ready to dispatch** — full pre-pass draft already written by Claude Code, same publish path as entry #8 | none |
| 4 | Act III synthesis — reconcile 141-150 + KIKORYU audit, apply SAGA/NORN thinking by hand, write tournaments-platform backlog | `act3-synthesis.md` | **ready to dispatch, paused** | none — failed twice on transient `529 Overloaded` (not session-limit, not content). Per founder instruction 2026-07-16, not being auto-retried; dispatch deliberately whenever picked back up. |

---

## Dispatched (awaiting completion)

| # | Title | Dispatched | Notes |
|---|---|---|---|
| 5 | DNS operations northstar ("fates dns" — EINHORN operating its own DNS) | 2026-07-17 | Full standalone prompt written (`dns-operations-northstar.md`) after 2 prior attempts failed on transient 529s with no saved prompt; dispatched with real grounding (farthq.com's current Cloudflare NS confirmed via `dig`, real product-surface inventory, NORN/SAGA cross-references). |

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
