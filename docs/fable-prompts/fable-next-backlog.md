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
| 6 | **TOP PRIORITY** — entity-graph/signalapi/newssite full-in-memory-replay fragility fix | `replay-fragility-northstar.md` | **ready to dispatch, top priority** | none — also pinned at the top of `BACKLOG.md` SECTION 1 so it's picked up first regardless of this queue's order |
| 1 | IDUNA front door funnel for Agents and Unagents | `iduna-front-door-funnel.md` | **ready to dispatch** | SAGA VS-spec audit landed 2026-07-16 (IDUNA `529834f`, EMILY `508a9be`) — no remaining blocker |
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
| 2 | Emilyify the mag book → `gpt2-alpine-c/pkg/towerprint` + `docs/TOWERPRINT.md` | 2026-07-16 | gpt2-alpine-c `819c400`, EMILY `c0efd2f` |

*(none yet)*

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
