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

## Commit Protocol (standing instruction)

Always commit and push completed work immediately — don't wait to be asked. This is the default for every repo in this monorepo.

