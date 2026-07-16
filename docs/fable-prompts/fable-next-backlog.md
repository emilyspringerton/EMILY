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

---

## Dispatched (awaiting completion)

| # | Title | Dispatched | Notes |
|---|---|---|---|
| 2 | Emilyify the mag book (`QUEENSALLYONLINEBOOKOFMAGIFICATIONANDUNICOR` → Go port + design doc) | 2026-07-16 | Feeds `EMILY/BACKLOG.md` SECTION 147 (S147-01). No standalone prompt file written — dispatched directly in-session; see continuity report / session transcript for the full brief if it needs re-sending. |
| — | TYLER S00E-2 (prequel bridge to S00E01) | 2026-07-16 | **Failed**, not landed — hit a Fable session/usage limit before writing anything (verified clean tree). Resets ~5:50am UTC 2026-07-17. No prompt file written yet either — should get one before this session ends so a fresh session can retry without re-deriving the brief. |

---

## Done

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
