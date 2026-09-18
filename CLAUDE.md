# EMILY — Emily Prime Agent

Emily Prime is the meta-orchestrator and chief of staff for EINHORN_INDUSTRIAL. She runs as a Go
HTTP service on port 8086, executes RSI (Recursive Self-Improvement) cron cycles every 15 minutes,
triages FatBaby observations, issues directed tasks to the obs-watcher loop, files Apples to IDUNA
after each cycle, and sends FCM push notifications to MJOLNIR on critical events.

## What Lives Here

| Path | What it is |
|---|---|
| `emily-agent/` | Core agent binary: cron, RSI loop, HEIMDAL, vision, briefing, prime triage |
| `scripts/rsi-loop.sh` | TIC→TOCK→ENTROPY→ANALYZE tic-toc shell loop |
| `signals/observations/` | FatBaby observation files Emily Prime reads each cycle |
| `signals/tasks/` | Directed task files obs-watcher picks up and dispatches to Claude Code |
| `inbox/` | No-context dump — anything shared with no clear task attached lands here first; see THE_EMILY_WAY.md Principle 16 |
| `var/rsi-loop-state.json` | Live RSI loop state (read by TUI) |
| `BACKLOG.md` | The canonical cross-repo golden backlog — read before starting any work |
| `GOLDEN.md` | Compressed backlog context for haiku (≤1200 tokens) |
| `DONE.md` | Archived completed items |
| `emily-memory/` | Emily Prime's persistent memory across cycles |

## Key Env Vars

```
ANTHROPIC_API_KEY   — required for all LLM calls (haiku + sonnet)
IDUNA_BASE_URL      — e.g. http://localhost:8080
IDUNA_AGENT_NAME    — EMILY_PRIME
IDUNA_AGENT_SECRET  — M2M credential
APPLES_GIT_DIR      — /home/fatbaby/APPLES (triggers auto-sync after each Apple POST)
FCM_PROJECT_ID      — Firebase project for MJOLNIR push notifications
FCM_SERVICE_ACCOUNT_JSON — path to service account JSON
GMAIL_CLIENT_ID     — for CEO escalation emails (optional)
GMAIL_CLIENT_SECRET
GMAIL_REFRESH_TOKEN
EMILY_STATE_DIR     — default: ./emily-state
```

## Emily Prime Cron Cycle (RunOnce)

Each 15-minute cycle has four phases:

1. **OBSERVE** — load state, count roadmap items
2. **DECIDE** — pick highest-priority queued RSI task (Emiree gear-aware)
3. **ACT** — run one RSI iteration (generate → evaluate → update task)
4. **PLAN** — triage FatBaby observations, drain HEIMDAL sprints, vision cycle, morning briefing, file Apple

After each cycle: Emiree witch-engine updates gear (ACTIVE/COAST/REST) based on outcomes.

## Apple Filing Protocol

Every meaningful cycle outcome is filed as an Apple via `POST /api/v1/apples`. Apple types:
- `improvement` — task succeeded or iterated
- `observation` — triage found new directed tasks
- `audit` — idle/monitoring cycle
- `escalation` — CEO-visible critical signal

## HEIMDAL Integration

MJOLNIR sends product requirements → IDUNA `heimdal_sprints` → Emily Prime translates via
claude-haiku → RSI roadmap item + Apple + FCM push → Claude Code executes → Emily Prime
patches sprint to `complete`/`blocked` + FCM push.

## Backlog Protocol

1. Read BACKLOG.md before starting any work.
2. Pick the highest-priority `[ ]` item in the lowest-numbered open section — unless founder
   real-time direction is present, in which case route it through `emily observe` first (see
   `docs/THE_EMILY_WAY.md` Principle 18), then log it into BACKLOG.md before working it.
3. Do the work.
4. Post Apple to IDUNA (`emily apples post` auto-tags with the active `emily session`).
5. Mark `[x]` with Apple ID; hand-written BACKLOG.md entries should also carry the session tag
   (`emily session current`) for traceability.
6. `git add BACKLOG.md && git commit && git push`.

## Related Repos

- `PRRJECT_FATBABY` — signal pipeline; publishes observations Emily Prime reads
- `IDUNA` — IAM + Apples store (`:8080`)
- `emily.cli` — CLI for Emily (emily observe, emily status, emily tui)
- `MJOLNIR` — Android app; Emily's phone
- `APPLES` — git-authoritative Apple backup (`emily sync --apples-git-dir`)

## CHANGELOG Protocol

After any meaningful change, update CHANGELOG.md:
```bash
emily changelog add EMILY "<what changed>"
# or manually: append a dated bullet under ## YYYY-MM-DD in EMILY/CHANGELOG.md
```

## Golden Doc Registration

If you create a new NORTHSTAR.md, architecture spec, or mission-critical design doc in this repo,
append a row to `EMILY/context/golden-docs-index.md` so Emily Prime picks it up on the next cycle:
```
| NAME | <repo>/path/to/doc.md | 1 | <budget-or-0> | one-line description |
```
Then commit and push EMILY:
```bash
cd /home/fatbaby/EMILY && git add context/golden-docs-index.md && git commit -m "golden-index: add NAME" && git push
```

## RSI AGI Loop

To start the full RSI loop with persistent context across cycles:
```bash
emily start --agi        # enables --continue: each claude invocation continues the prior session
emily start              # default: each observation starts a fresh claude session
```

The `--agi` flag wires `OBSERVATION_CONTINUE=true` into obs-watcher, which appends `--continue`
to every claude invocation. Sessions accumulate context across RSI cycles — each iteration starts
from where the last one ended rather than a blank slate.

See `EMILY/docs/THE_EMILY_WAY.md` principle 11 for the full AGI loop rationale.

## CONSTRUCT File Generation (standing instruction, monorepo Principle 21)

EMILY auto-generates CONSTRUCT files on each release via CI. A CONSTRUCT is a deterministic plaintext snapshot of all repo source files with file metadata — used for reproducible builds, audit trails, and offline source access. 

See the main `CLAUDE.md`'s "Principle 21: CONSTRUCT Files" section for the full rationale and shared implementation patterns. The generation is automatic in CI; no manual work needed.

## README Reality — SAGA reconciliation (standing instruction, monorepo-wide)

Founder real-time, 2026-09-18: if a change of yours **substantially changes the claim of this project's core README**,
then per SAGA protocols (`EMILY/docs/SAGA_SYSTEM_AUDIT_2026-07-18.md`, HQ-SPEC-DOC-102: intent ↔ claim ledger ↔ reality)
you **must update `README.md` in the same unit of work** so it reflects current reality. The README is the project's public
claim; it must not lag behind the code.

- **When it applies:** a capability is added or removed; status moves ("design only" → "working", "planned" → "shipped");
  the stack, build, run or install steps change; a claim in the README is now false or stale; or you add a **meaningful,
  genuinely interesting piece of kit** (a new tool, engine capability, protocol, pipeline, game system). For that last case
  especially: put it in the README — what it is, how to run it, and its honest status and limits.
- **When it does not:** ordinary fixes, refactors and small features that leave the README's claims true.
- **How:** re-read the README against what you just changed; fix or delete stale lines (including "not built yet" notes that
  are now built); verify any new claim by actually running it, and mark anything untested as untested; commit the README
  with (or immediately after) the change, and mention it in the CHANGELOG entry.

## Frame-Break Reframing

Founder-sourced prompting technique (REDGARDEN/NORTHSTAR.md §28, full origin in
REDGARDEN/docs2/MULTI_AGENT_RD_RESEARCH_NOTES.md §5): given a request, name the underlying
structural/systemic pattern it's one instance of — one level of abstraction up — as an added
lens during planning/triage/judgment calls. Use it to spot the general case behind a specific
ask. It augments judgment, it does not replace doing the work: direct, concrete execution of
the literal task asked for still happens every time.

## Commit Protocol (standing instruction)

Always commit and push completed work immediately — don't wait to be asked. This is the default for every repo in this monorepo.

Every commit — human-written or produced by automated code paths (git-commit helpers in emily-agent, emily.cli, IDUNA handlers, etc.) — must carry the active `emily session` fingerprint as a `session: <tag>` trailer (blank line, then the trailer). This was silently missing from several independently-implemented automated commit helpers across the monorepo until an audit on 2026-08-10 (founder, real-time: "where in the fuck is my llm session id anywhere"). If you add a new automated git-commit code path anywhere, wire in the session tag the same way — don't assume an existing helper already does it.
