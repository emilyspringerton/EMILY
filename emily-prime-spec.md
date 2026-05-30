# Emily Integration Specification
## Emily Prime ↔ FatBaby-Emily Communication Layer

**Status:** Living Document — Current Iteration  
**Parent:** EINHORN_INDUSTRIAL / Emily Architecture Series  
**Supersedes:** Legacy FARTCO integration spec (circa Kubernetes/Gmail OAuth era)

---

## What This Is

This document defines the integration layer between two Emily instances:

- **Emily Prime** — the meta-orchestrator and CEO-facing chief of staff. Strategic, emergent, self-improving. Operates at the level of mission, resource allocation, and cross-system coordination. Talks to the CEO. Talks to herself. Grows.

- **FatBaby-Emily** — the domain-specific signal intelligence agent embedded inside `prrject-fatbaby`. She watches SEC EDGAR filings and PR Newswire press releases, processes them into structured signals, runs the EPS oracle, and knows the corpus deeply. She does not see the whole board — that is not her job.

The integration allows Emily Prime to direct FatBaby-Emily's ongoing development and expansion of data collection and signal extraction — while FatBaby-Emily has access to higher-level corporation docs so she can see enough of the big picture to make good local decisions without needing to ask Prime for everything.

---

## Why These Two Need to Talk

FatBaby-Emily is very good at what she does. She knows the pipeline. She knows the tickers. She sees filing anomalies, EPS surprises, governance signals. She publishes observations and the observation-watcher loop fires Claude Code to act on them.

But she operates in a closed domain. She cannot:

- Decide which new data sources are worth adding to the watchlist
- Prioritize her own improvement roadmap relative to the broader mission
- Know whether a governance signal she surfaced is strategically important or noise in context
- Surface critical findings to the CEO at the right moment

Emily Prime closes these gaps. She reads FatBaby-Emily's observations, understands them in broader strategic context, decides what to escalate and what to defer, and issues directed improvement tasks back into the FatBaby loop.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│  CEO                                                    │
│  ↕  (escalations, strategy, cross-training)             │
└─────────────────────────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│  EMILY PRIME                                            │
│  Meta-orchestrator / Chief of Staff                     │
│  ├─ Reads FatBaby observations                          │
│  ├─ Triages strategic relevance                         │
│  ├─ Issues directed improvement tasks                   │
│  ├─ Manages Emily's own improvement loop                │
│  └─ Talks to CEO (surfaces, escalates, summarizes)      │
└────────────────────────┬────────────────────────────────┘
                         │  directed tasks + context
                         │  structured observations ↑
┌────────────────────────▼────────────────────────────────┐
│  FATBABY-EMILY (cmd/emily-agent)                        │
│  Domain-specific signal intelligence                     │
│  ├─ SEC EDGAR watcher (secwatch)                        │
│  ├─ PR Newswire watcher (prwatch + prwatch-body)        │
│  ├─ Signal processor (EPS, governance, sentiment)       │
│  ├─ Entity graph (directors, auditor changes, votes)    │
│  └─ Observation writer → var/emily-observations/        │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│  PRRJECT-FATBABY PIPELINE                               │
│  eventstore / feedserver / broker / signalapi / newssite │
└─────────────────────────────────────────────────────────┘
```

---

## The Two Heavy Hitters

### 1. Git Integration — The Auditability Spine

Git is not just version control here. It is the memory, the audit trail, and the coordination substrate for the entire Emily system.

**What it provides:**

Every decision Emily Prime makes, every improvement task she issues, every observation FatBaby-Emily publishes — all of it flows through Git. This means:

- The system is fully auditable. Any decision can be traced back to the observation that triggered it, the reasoning Emily applied, and the outcome.
- Improvement is versioned. `decision-heuristics.v3.md` exists because `v2` was measured and found wanting. You can diff them. You can roll back.
- Emily's self-modification is controlled. She can update her own heuristics and agent configs, but these changes are commits — visible, reversible, subject to review.
- FatBaby-Emily's codebase changes go through PRs. The observation-watcher loop invokes Claude Code, which edits source, runs `go test ./...`, and commits on pass. That commit is a traceable artifact.

**Current implementation (FatBaby side):**

The Emily agent in `cmd/emily-agent` already writes observations to `var/emily-observations/latest.json` with a timestamped archive sibling. The `observation-watcher` process polls this file and shells out to Claude Code.

**What needs adding for Prime integration:**

Emily Prime needs read access to the FatBaby observation stream and write access to issue directed tasks back. The cleanest mechanism is a shared Git repo (or a dedicated integration repo) where:

- FatBaby-Emily commits observations to `signals/observations/`
- Emily Prime commits directed tasks to `signals/tasks/`
- Both sides poll their respective directories

This keeps the message bus auditable by default, with no additional infrastructure. Timestamps, authors, and reasoning all live in commit history.

**Secrets required:**

```
GIT_ACCESS_TOKEN     — read/write to integration repo
GIT_WEBHOOK_SECRET   — validate incoming push events (optional, for push-triggered loops)
```

---

### 2. Email Integration — Emily as True Chief of Staff

This is where Emily Prime crosses from an ops agent into a genuine chief of staff role.

The vision: Emily Prime monitors the CEO's inbox, understands context, triages against the strategic picture she maintains, and surfaces what actually matters. Not every email. The right ones, at the right moment, with enough context to act.

This is also the outbound channel. When FatBaby-Emily surfaces a significant governance signal — a director departure, an auditor change, an EPS miss with unusual circumstances — Emily Prime decides whether this warrants a CEO-level alert and sends it. The signal that would otherwise sit in a dashboard becomes an email in an inbox.

**What it provides:**

- Inbound triage: Emily monitors the CEO inbox, classifies by urgency and strategic relevance, surfaces what needs a decision
- Outbound alerts: Emily sends formatted summaries of material signals (governance events, EPS surprises, watchlist anomalies) to the right people
- Weekly pulse: Emily Prime auto-generates and sends a strategic summary — what FatBaby-Emily observed this week, what it means, what decisions are pending

**Current state:**

The email integration is planned but not wired. The old spec defined the `GmailClient` interface and OAuth setup. That interface design is sound. The update is that it now connects Emily Prime (not a standalone service) to the CEO's actual workflow.

**Secrets required:**

```
GMAIL_CLIENT_ID       — Google OAuth2 client
GMAIL_CLIENT_SECRET   — Google OAuth2 secret
GMAIL_REFRESH_TOKEN   — OAuth2 refresh token
GMAIL_PUBSUB_TOPIC    — for push notifications (optional, falls back to polling)
```

---

## Context Sharing: What FatBaby-Emily Sees

FatBaby-Emily needs to see enough of the big picture to make good local decisions. She should not need to ask Emily Prime whether to add a new ticker to the watchlist — she should be able to decide that herself if the strategic context is available.

**What she gets access to:**

- `config/watchlist.json` — already in repo; she can read and propose changes
- Corporate strategy docs (read-only) — committed to the Emily Prime repo under `context/strategy/`; FatBaby-Emily pulls these periodically
- Signal priority weights — Emily Prime publishes a `context/signal-priorities.json` that tells FatBaby-Emily which signal types are currently most strategically relevant

**What she does not get:**

- The CEO inbox
- Emily Prime's internal reasoning logs
- IAM and governance infrastructure details

The principle: FatBaby-Emily gets enough context to be a good domain agent, not enough to be distracted from her job.

---

## Communication Protocol

### FatBaby-Emily → Emily Prime (Observations)

FatBaby-Emily already produces structured observations via `fatbaby_write_observation`. For Prime integration, these observations are committed to the shared integration repo:

```json
{
  "timestamp": "2026-05-30T14:22:15Z",
  "source": "fatbaby-emily",
  "observation_type": "anomaly | improvement | escalation | status",
  "severity": "critical | high | normal | low",
  "summary": "...",
  "detail": "...",
  "affected_tickers": ["AAPL", "MSFT"],
  "signal_ids": ["sig_001", "sig_002"],
  "recommended_action": "...",
  "requires_ceo_visibility": false
}
```

### Emily Prime → FatBaby-Emily (Directed Tasks)

Emily Prime writes task files to `signals/tasks/<timestamp>-<id>.json`:

```json
{
  "timestamp": "2026-05-30T14:30:00Z",
  "from": "emily-prime",
  "to": "fatbaby-emily",
  "task_type": "expand_coverage | improve_signal | fix_anomaly | config_change",
  "priority": "high",
  "description": "...",
  "acceptance_criteria": ["...", "..."],
  "context": "Brief strategic rationale FatBaby-Emily needs to act well",
  "deadline": "2026-05-31T00:00:00Z"
}
```

FatBaby-Emily's observation-watcher loop picks these up alongside Emily observations. Same mechanism, second input channel.

---

## Emily Prime's Self-Improvement Loop

Emily Prime's strategy is emergent. She improves through interaction with the CEO and through what she observes across all systems. The audit trail is what makes this safe and legible.

**The loop:**

1. Emily Prime reads FatBaby observations
2. Emily Prime triages — strategic relevance, CEO visibility, improvement opportunity
3. Emily Prime updates her own decision heuristics (versioned, committed)
4. Emily Prime issues directed tasks back to FatBaby-Emily
5. Emily Prime escalates to CEO where warranted (email)
6. CEO feedback flows back into Emily Prime's context
7. Repeat

**What never changes without human review:**

- Core values and ethical constraints
- IAM and governance rules (separate platform, separate track)
- Audit log deletion or modification

---

## Deliverables — Current Iteration

| Component | Owner | Output |
|---|---|---|
| Integration repo setup | Engineering | `signals/observations/`, `signals/tasks/` directories with schema |
| FatBaby observation commit hook | Engineering | Auto-commit on `fatbaby_write_observation` |
| Emily Prime observation reader | Emily Prime / Engineering | Polls integration repo, triages, routes |
| Directed task writer | Emily Prime | Commits task files to `signals/tasks/` |
| FatBaby task poller | Engineering | Extends observation-watcher to read Prime tasks |
| Gmail OAuth setup | Engineering | Service account + tokens for CEO inbox access |
| Gmail client (Prime side) | Engineering | Read inbox, send alerts, weekly pulse |
| Strategic context publishing | Emily Prime | `context/strategy/` and `context/signal-priorities.json` |
| CEO escalation formatter | Emily Prime | Structured email template for material signal alerts |

---

## Immediate Actions

**Engineering:**

- Set up integration repo (or dedicated directory in Emily Prime repo) with `signals/` directory structure
- Extend `cmd/observation-watcher` to also poll `signals/tasks/` from Prime
- Add Git commit step to `fatbaby_write_observation` tool
- Begin Gmail OAuth setup (service account creation)

**Emily Prime:**

- Define initial `context/signal-priorities.json` — which FatBaby signal types are currently most strategically relevant
- Define CEO escalation criteria — what threshold triggers an email vs. stays in dashboard
- Draft weekly pulse email template

**Open questions for CEO:**

1. Which inbox(es) should Emily monitor — CEO only, or broader leadership?
2. What is the escalation threshold for governance signals? (Director departure: always? Auditor change: always? EPS miss: depends on magnitude?)
3. Should Emily Prime's directed tasks to FatBaby require human confirmation before FatBaby acts, or is autonomous action within defined scope acceptable?

---

## What Success Looks Like

Emily Prime is working when:

- A material governance signal surfaces in FatBaby at 11pm and the CEO has a clear, contextualized summary in their inbox by 7am
- FatBaby-Emily's watchlist expands to a new sector because Emily Prime identified a strategic gap and issued a directed task — without a human having to write a Jira ticket
- Every decision Emily Prime made last month can be traced: observation → reasoning → action → outcome → heuristic update
- The CEO's inbox triage is measurably lighter because Emily Prime is handling the first pass

---

*The goal is not automation for its own sake. It is a system where AGI self-direction happens through clear, open, auditable mechanisms — and where the humans who need to stay in the loop actually can.*
