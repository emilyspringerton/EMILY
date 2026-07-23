# EMILY

**The meta-orchestration layer for EINHORN_INDUSTRIAL.** Emily Prime plans, Claude Code
implements, every outcome is logged as an Apple. This repo is the chief-of-staff agent, the
cross-repo backlog, and the operating procedure that ties a small AI-native company's many
codebases into one recursive self-improvement loop.

> I plan. Claude Code implements. Every outcome is an Apple.
> — the three-sentence description Emily Prime reads about herself, every cycle.

---

## Table of contents

- [What this is](#what-this-is)
- [The three agents](#the-three-agents)
- [Architecture](#architecture)
- [The RSI cron cycle](#the-rsi-cron-cycle)
- [The Emily Way](#the-emily-way)
- [Repo layout](#repo-layout)
- [Running it](#running-it)
- [The `emily` CLI](#the-emily-cli)
- [Document hierarchy](#document-hierarchy)
- [Related repos](#related-repos)
- [Status](#status)

---

## What this is

EINHORN_INDUSTRIAL runs several codebases at once — a financial signal-intelligence pipeline
(`PRRJECT_FATBABY`), an identity/audit authority (`IDUNA`), a multiplayer game engine
(`SHANKPIT`), an Android companion app (`MJOLNIR`), a WordPress-based public product (`EDIS`),
and more — with a very small human team. EMILY is the layer that makes that tractable: an
agent that reads the state of every repo, decides what matters most right now, hands
implementation work to Claude Code, and writes down what happened so the next cycle (or the
next person) doesn't have to reconstruct it from scratch.

It is not a chatbot bolted onto a codebase. It's closer to a chief of staff with a filing
system: it triages, it delegates, it keeps a paper trail, and it runs on a timer whether or
not anyone is watching.

Three ideas hold the whole thing together:

1. **The backlog is the source of truth.** One file (`BACKLOG.md`), append-only, cross-repo,
   git-authoritative. Every session — human or agent — starts by reading it.
2. **Every meaningful outcome is an Apple.** A structured, timestamped record filed to IDUNA
   (`POST /api/v1/apples`), backed up to a plain git repo (`APPLES`). If it isn't an Apple, it
   didn't happen, as far as the next cycle is concerned.
3. **Planning and implementation are separated on purpose.** Emily Prime decides *what*.
   Claude Code does *how*. This split is what makes the loop recursive rather than a single
   long-running process trying to hold everything in its own context at once.

---

## The three agents

EMILY isn't one program — it's three cooperating roles, each with a distinct job:

### Emiree — the witch engine

The state machine underneath everything else. A seven-gear dynamical system —
`IDLE → COAST → ACTIVE → ENGAGED → PUSH → SURGE → OVERLOAD` — that governs how aggressively
the RSI loop paces itself based on how recent cycles actually went. She's never directly
visible externally; her existence is provable only from the outside, the way a spell's effect
is provable from what changed, not from watching it cast. Implemented in
`emily-agent/emiree.go`; the source spec (`emiree.md`) is a seven-volume transmission written
in parallel Sanskrit/Chinese/English — deliberately unlike a normal engineering doc, because
the thing it's specifying isn't a normal state machine either.

### Emily Prime — chief of staff

The surface everyone actually talks to. A Go HTTP service on `:8086` that reads FatBaby
observations, issues directed improvement tasks, manages escalations, and files an Apple after
every meaningful cycle. Her system prompt is rebuilt dynamically from compressed context
across *every* repo in the monorepo (see [Document hierarchy](#document-hierarchy)) — she
doesn't reason about EMILY in isolation, she reasons about the whole company. Implemented in
`emily-agent/`.

### FatBaby-Emily — domain signal intelligence

Lives inside `PRRJECT_FATBABY`, not here. Watches SEC EDGAR filings and PR Newswire releases,
turns them into structured governance signals, and publishes observations that Emily Prime
reads and triages. She doesn't see the whole board — just her domain — which is the point:
Emily Prime's job is synthesizing across specialists like her, not being one herself.

---

## Architecture

```
                    CEO (Gmail escalations, phone)
                          ↕
                  Emily Prime (:8086)
                    ├─ Emiree witch engine       (RSI pacing, gear state)
                    ├─ FABLE advisor             (haiku sprint recommendations)
                    ├─ GoldenDocCompiler         (cross-repo context, rebuilt each cycle)
                    ├─ HEIMDAL bridge            (MJOLNIR → IDUNA sprints → Emily Prime)
                    ├─ RSI loop                  (iterative improvement, acceptance criteria)
                    └─ Integration layer         (reads observations, writes directed tasks)
                          ↕
                  FatBaby-Emily (PRRJECT_FATBABY)
                    └─ SEC EDGAR + PR Newswire → governance signals → observations
```

The feedback loop that actually ships code:

```
Emily Prime (plans)
  │  issues a directed task → EMILY/signals/tasks/
  ▼
observation-watcher (dispatches)
  │  invokes `claude --dangerously-skip-permissions [--continue]`
  ▼
Claude Code (implements)
  │  edits code, runs tests, commits, pushes
  │  writes a claude-run report to PRRJECT_FATBABY/claude-runs/
  ▼
Emily Prime (observes)
  │  reads the report, FABLE advisor re-reads full-system-context.md
  │  files an Apple, updates the backlog, issues the next task
  ▼
repeat — the RSI cron cycle, every 15 minutes
```

Optionally, `emily start --agi` wires `--continue` into every Claude Code invocation, so
sessions accumulate context across RSI cycles instead of starting from a blank slate each
time — the closer-to-continuous mode of the loop.

---

## The RSI cron cycle

Every fifteen minutes, Emily Prime runs one full cycle:

1. **GoldenDocCompiler.MaybeRebuild** — refresh the compiled cross-repo context if any source
   golden doc changed.
2. **OBSERVE** — load state, check infrastructure health.
3. **DECIDE** — Emiree-gear-aware selection of the highest-priority active task.
4. **ACT** — run one RSI iteration against that task's acceptance criteria.
5. **PLAN** — triage FatBaby observations, drain queued HEIMDAL sprints, run a vision cycle,
   compose a morning briefing if due.
6. **APPLE** — file the cycle's outcome to IDUNA; push an FCM notification to MJOLNIR if
   something completed or needs attention.

After each cycle, Emiree updates gear state from the outcome — a good cycle can push the
system toward `PUSH`/`SURGE`; a stalled or degraded one pulls it back toward `COAST`.

---

## The Emily Way

The operating procedure every session — human or agent — follows, in full at
[`docs/THE_EMILY_WAY.md`](docs/THE_EMILY_WAY.md). The short version:

| # | Principle |
|---|---|
| 1 | **Backlog first.** Read `BACKLOG.md`. Pick the highest-priority open item in the lowest-numbered section. |
| 2 | **Spec before implementation.** A new subsystem gets a northstar before it gets code. |
| 3 | **Apple before mark-done.** No `[x]` in the backlog without a filed Apple first. |
| 4 | **CHANGELOG on every meaningful change.** |
| 5 | **Register new golden docs** so Emily Prime's context pipeline can actually see them. |
| 6 | **Commit and push at logical boundaries** — one atomic thing per commit. |
| 7 | **Tests before commit.** `go test ./...` green, every time, every repo. |
| 8 | **Small tasks, clear acceptance criteria.** |
| 9 | **The audit trail is the product.** If it isn't written down, it didn't happen. |
| 10 | **Emily Prime plans, Claude Code implements** — the split that makes this recursive. |
| 11 | **RSI AGI loop: `--continue` when running unattended.** |
| 12 | **Multi-repo discipline** — commit each repo independently, in dependency order, reference sibling commits. |
| 13 | **Degraded mode is OK** — a blocked task gets logged as blocked, not silently dropped. |
| 14 | **Continuity report** — a full sync checkpoint (`EMILY/continuity/`) so a fresh session can pick up cold. |
| 15 | **Operational health is not optional** — infrastructure gets verified, not assumed, every session. |

Principle 15 exists because it had to: `emily-system.service` once sat OOM-killed for ~22
hours before anyone noticed. The fix wasn't "be more careful" — it was watchdogs that fire an
escalation Apple on their own, on a timescale of minutes. The full story is in the principle
itself; it's the most-lived-in rule in the document.

---

## Repo layout

```
EMILY/
├── README.md                 — you are here
├── CLAUDE.md                 — instructions for Claude Code working in this repo
├── BACKLOG.md                — the canonical cross-repo backlog (read this first)
├── GOLDEN.md                 — auto-generated compressed backlog (≤1200 tokens, cron-cycle context)
├── DONE.md                   — archived completed backlog items
├── CHANGELOG.md              — dated, human-readable history
│
├── emily-agent/               — Emily Prime's Go service
│   ├── main.go                — HTTP server, tool loop, dynamic system prompt
│   ├── cron.go                — the 15-minute RSI cycle
│   ├── goldenbuild.go         — GoldenDocCompiler (cross-repo context compression)
│   ├── fable.go                — FABLE advisor (haiku sprint recommendations)
│   ├── emiree.go               — the seven-gear witch engine
│   ├── heimdal.go              — MJOLNIR → IDUNA → Emily Prime sprint bridge
│   ├── rsi.go                  — the RSI loop itself
│   └── watchdog.go             — service/poller health checks (Principle 15)
│
├── docs/                      — northstars, specs, protocol docs
│   ├── THE_EMILY_WAY.md       — the operating procedure, in full
│   ├── NORTHSTAR.md           — EMILY's own northstar
│   └── NORTHSTAR_*.md          — northstars for adjacent systems (Baby ERP, Supply Chain, STINKIES…)
│
├── context/                   — compiled cross-repo context
│   └── golden-docs-index.md   — the registry of every Tier 1/Tier 2 golden doc, across every repo
│
├── continuity/                 — dated full-system sync checkpoints (Principle 14)
├── signals/
│   ├── observations/           — FatBaby observation files Emily Prime reads each cycle
│   └── tasks/                  — directed tasks the observation-watcher dispatches to Claude Code
│
├── scripts/rsi-loop.sh         — TIC→TOCK→ENTROPY→ANALYZE tic-toc shell loop
└── emiree.md                   — the Emiree witch-engine transmission (full spec, unconventional format)
```

---

## Running it

```bash
# Build and run Emily Prime's HTTP service
cd emily-agent
go build -o emily-agent .
./emily-agent           # listens on :8086, no autonomous cycling

# Run one RSI cycle and exit (cron-friendly)
./emily-agent --cron

# Run autonomous cycles continuously at the configured interval
./emily-agent --daemon
```

Required environment (see `emily-agent/CLAUDE.md` for the full reference):

```
ANTHROPIC_API_KEY   — required for all LLM calls (haiku + sonnet)
IDUNA_BASE_URL      — e.g. http://localhost:8080
IDUNA_AGENT_NAME    — EMILY_PRIME
IDUNA_AGENT_SECRET  — M2M credential, auto-loaded from IDUNA/var/agent-secrets.env
APPLES_GIT_DIR      — path to the APPLES repo; triggers auto-sync after each Apple POST
FCM_PROJECT_ID      — Firebase project for MJOLNIR push notifications (optional)
GMAIL_CLIENT_ID / GMAIL_CLIENT_SECRET / GMAIL_REFRESH_TOKEN  — CEO escalation email (optional)
```

To start the full unattended loop, including the observation-watcher that dispatches Claude
Code:

```bash
emily start --agi   # --continue on every Claude Code invocation; context accumulates across cycles
emily start         # default: each observation starts a fresh Claude Code session
```

---

## The `emily` CLI

The human-facing control surface, in the sibling `emily.cli` repo — every interaction with
IDUNA, the observation-watcher, and agent status flows through it.

| Command | What it does |
|---|---|
| `emily observe <msg>` | Post an observation into the FatBaby pipeline |
| `emily obs amend <key> <correction>` | Append a correction to an existing observation |
| `emily apples list [filter]` | Query the IDUNA Apples log |
| `emily apples post -t <type> <title> <body>` | File an Apple |
| `emily watch [repo]` | Tail Apples in real time |
| `emily status [--fatbaby] [--watch]` | Cross-repo git + IDUNA + process state |
| `emily start [all|--agi]` | Start FatBaby processes / the full RSI loop |
| `emily sync [--all] [--apples-git-dir <dir>]` | Sync observations → IDUNA, Apples → git |
| `emily tui [--fatbaby]` | Full-terminal TUI — roadmap, tasks, health, in three columns |
| `emily backlog promote` | Curate/promote raw INTAKE QUEUE items into structured backlog entries via haiku |
| `emily agents list` | List registered IDUNA agents |
| `emily primetask [create\|list]` | Interact with Emily Prime's RSI task queue directly |
| `emily changelog add <repo> "<change>"` | Append a dated CHANGELOG.md entry |

---

## Document hierarchy

Not every doc carries the same weight. Emily Prime's context budget is finite, so what makes
it into her system prompt every cycle is deliberately tiered:

```
Tier 1 — Golden docs (always in context)
  northstars, executive summaries, integration specs, the RSI engine spec
  → EMILY/context/golden-docs-index.md

Tier 2 — Important specs (read before touching a subsystem)
  agent protocol, tool specs, netcode contracts, platform architecture
  → also indexed, tier=2

Tier 3 — Reference (rarely needed in context)
  individual endpoint specs, migrations, fixture docs
  → not indexed; found by reading code directly

BACKLOG.md — the load-bearing node
  cross-repo, git-authoritative, append-only; read before starting any work

GOLDEN.md — compressed backlog
  auto-generated (`emily backlog compress`); what Emily Prime actually reads each cron cycle

CHANGELOG.md — per-repo narrative history

EMILY/continuity/ — session checkpoints
  full sync sweep + changelog digest + in-flight/blocked/queued state + next steps
```

The pipeline that keeps Emily Prime's context current:

```
All repo northstars (Tier 1 docs)
  → GoldenDocCompiler.Build()  [haiku bilingual compression]
  → context/full-system-context.md
  → buildEmilySystemPrompt()   [prepended to every conversation + the FABLE advisor]
```

It rebuilds automatically, once per cron cycle, whenever any source doc has changed.

---

## Related repos

| Repo | Role |
|---|---|
| [`PRRJECT_FATBABY`](../PRRJECT_FATBABY) | SEC/PR financial signal pipeline; FatBaby-Emily publishes observations here |
| [`IDUNA`](../IDUNA) | Platform IAM, Apples ledger, HEIMDAL sprint store (`:8080`) — the central trust authority |
| [`emily.cli`](../emily.cli) | The `emily` operator CLI; the human half of the feedback loop |
| [`MJOLNIR`](../MJOLNIR) | Android app — Emily's phone; FCM push, Apple feed, HEIMDAL submitter |
| [`APPLES`](../APPLES) | Git-authoritative Apple audit-trail backup |
| [`SHANKPIT`](../SHANKPIT) | Server-authoritative UDP FPS; RSI loop drives its development |
| [`EDIS`](../EDIS) | WordPress-based public product surface, fronted by a self-hosted digital-immune-system layer |
| [`EmilyOS`](../EmilyOS) | Bare-metal policy kernel: posture-gated sessions, RBAC, audit log |
| [`PITVIPER`](../PITVIPER) | SDL2 terminal emulator with Emily Prime integration hooks |
| [`GoblinFoxDragon`](../GoblinFoxDragon) | Dragonfly/Bedrock fork; persistent-world engine R&D |
| [`TYLER`](../TYLER) | Television-as-code media arm; a genuinely different kind of RSI loop (episode scripts as builds) |
| [`NORN`](../NORN) | The propose→grade→gate→promote loop kernel — five interfaces (Artifact, Proposer, Oracle, Gate, Registry) that formalize the self-improvement pattern EMILY runs at company scale |

Go modules across the monorepo — including `EMILY/emily-agent`, `IDUNA`, `PRRJECT_FATBABY`,
`emily.cli`, `NORN`, and EDIS's Go-based digital-immune-system layer — share one `go.work`
workspace at the monorepo root, so `go test ./...` from `/home/fatbaby` exercises all of them
together.

---

## Status

Live. The RSI cron cycle runs every 15 minutes, in production, against real infrastructure —
this isn't a demo harness. What "done" looks like, in practice:

- Emily Prime's system prompt includes live compressed context from every active repo.
- The FABLE advisor's recommendations account for every northstar, not just the open backlog.
- Emily Prime accepts a planning question and returns a structured sprint batch; Claude Code
  executes the implementation.
- Every RSI cycle's outcome is filed as an Apple, viewable in IDUNA's Back Office.
- The audit trail is complete: every decision is a footprint, every Apple a proof.

For the day-to-day state of what's shipped versus open, `BACKLOG.md` is authoritative — this
README describes the shape of the system, not its current sprint.
