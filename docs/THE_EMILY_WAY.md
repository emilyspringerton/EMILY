# The Emily Way
## How EINHORN_INDUSTRIAL builds things

*Synthesized from CLAUDE.md files, git commits, changelogs, and RSI loop observations across all repos.*
*Last updated: 2026-07-17*

---

## Three-Sentence Version

Emily Prime plans; Claude Code implements; every outcome is an Apple.
The backlog is the load-bearing node — every decision is traceable to a backlog item, every completion has an Apple as proof.
Recursive self-improvement is the product: each cycle makes the next cycle faster.

---

## The Principles

### 1. Backlog First
Always read `EMILY/BACKLOG.md` before starting any work. Pick the highest-priority `[ ]` item in the lowest-numbered open section. Do not start new work until you know what the backlog says.

### 2. Spec Before Implementation
Architecture memos, northstars, and design docs go into `docs/` first. Claude Code implements what is already specced. No guessing at intent — if there's no spec, write the spec and get it committed before writing Go/Kotlin/C.

Pattern observed in commits: every major feature (SHANKPIT Milestone 3, goldenbuild, HEIMDAL, Ask Emily) arrived with a spec committed in `docs/specs/` or as a northstar before code landed.

### 3. Apple Before Mark-Done
Every `[x]` in BACKLOG.md requires an Apple posted to IDUNA first:
```bash
emily apples post -t completion -repo <REPONAME> "<title>"
```
No Apple = not done. The Apple is the proof. The Apple ID goes in the backlog entry next to `[x]`.

### 4. CHANGELOG on Every Meaningful Change
```bash
emily changelog add <repo> "<what changed>"
# or: append dated entry under ## YYYY-MM-DD in <repo>/CHANGELOG.md
```
This is mandatory, not optional. Every agent run, every Claude Code session, every meaningful commit needs a CHANGELOG entry. The CHANGELOG is the human-readable narrative of the audit trail.

### 5. Register New Golden Docs
If you create a new NORTHSTAR.md, architecture spec, or mission-critical design doc:
```bash
# Append to EMILY/context/golden-docs-index.md:
| NAME | <repo>/path/to/doc.md | <tier> | <budget> | one-line description |

# Then commit EMILY:
cd /home/fatbaby/EMILY && git add context/golden-docs-index.md && git commit -m "golden-index: add NAME" && git push
```
Anything not in the golden index is invisible to Emily Prime. Visibility = inclusion.

### 6. Commit and Push at Logical Boundaries
Each commit is one atomic thing. Format:
```
feat(scope): description        — new capability
fix(scope): description         — bug fix
docs(scope): description        — documentation
chore(scope): description       — housekeeping
perf(scope): description        — performance improvement
backlog: ✓ S22-01               — backlog completion marker
golden-index: add NAME          — golden doc registration
```
Never batch unrelated changes in one commit. Push after every commit so the system is always in sync and Emily Prime can read the latest state.

### 7. Tests Before Commit
```bash
go test ./...  # in any Go repo before committing
```
Tests pass or the commit doesn't happen. For non-Go changes (docs, config), at minimum: read the file you just wrote before committing to catch typos.

### 8. Small Tasks, Clear Acceptance Criteria
Every backlog item has acceptance criteria: observable behavior, not "code written." Example:
- Bad: "implement MySQL projector"
- Good: "cmd/projector tails secwatch eventstore, projects signal_generated into governance_signals + entity_timeline. Idempotent. Graceful degradation when MYSQL_URL not set."

### 9. The Audit Trail Is the Product
Every decision is a footprint. Every Apple is a proof. The BACKLOG.md + DONE.md + APPLES git repo + IDUNA ledger together form the complete audit trail. This is not bureaucracy — it is the mechanism by which Emily Prime learns from past cycles and improves future ones.

### 10. Emily Prime Plans, Claude Code Implements
Strategic questions go to Emily Prime:
```bash
emily prime-task create "what should we build next in S23?"
# or: POST /api/v1/emily/plan {"question": "..."}
```
Claude Code tokens are for implementation. Planning in a Claude Code session when Emily Prime can answer is wasting the most expensive token bucket.

### 11. RSI AGI Loop: Always `--continue`
The observation-watcher invokes Claude Code for each RSI cycle. To build persistent context across cycles (the AGI loop), set:
```bash
OBSERVATION_CONTINUE=true emily start
# This adds --continue to every claude invocation:
# claude --dangerously-skip-permissions --continue "<prompt>"
```
With `--continue`, each RSI cycle is not a fresh session but a continuation of the prior one. Context accumulates. The system learns across iterations, not just within them. This is what makes it an AGI loop, not just automation.

### 12. Multi-Repo Discipline
Each repo is its own git history. Commits in PRRJECT_FATBABY, EMILY, IDUNA, etc. are independent. When a change touches multiple repos:
1. Commit each repo separately, in dependency order
2. Include cross-repo context in the commit message ("see EMILY commit abc123")
3. Update EMILY/BACKLOG.md last (it is the canonical record)

### 13. Degraded Mode Is OK
Emily Prime starts without `ANTHROPIC_API_KEY` in degraded mode. Obs-watcher works without IDUNA. Systems are designed to fail gracefully and not require all components to be live. Start with what you have; unlock full capability incrementally.

### 14. Continuity Report — Full System Sync Checkpoint
Long sessions accumulate state across many repos and background agents faster than any one
person or session can hold in their head. Before ending a long/complex session, after a batch of
background agents lands, or whenever explicitly asked for a "sync" or "continuity report," produce
one. It exists so a fresh session — human or Claude Code, no memory of what came before — can read
one file and know exactly where things stand, the same way `REBOOT_RUNBOOK.md` exists so a fresh
session can recover the stack without reconstructing history from logs.

A continuity report has four parts, in order:

1. **Full sync sweep.** `git fetch` + `git status --short --branch` across every repo in the
   monorepo home (not just the ones touched this session — the sweep is the point). Anything
   showing "ahead of origin" gets pushed immediately as part of producing the report, not noted as
   a TODO. A continuity report that ships with unpushed commits has failed at its one job.
2. **Changelog digest.** What actually changed, repo by repo, since the last report (or since
   session start if this is the first one) — pulled from each repo's `CHANGELOG.md` and recent
   `git log`, not re-derived from memory.
3. **Continuity state.** What's in-flight (background agents still running), what's blocked and on
   what (credentials, sudo, API/session limits — name the exact blocker, not "waiting on human"),
   what's queued but not yet dispatched (e.g. `EMILY/docs/fable-prompts/fable-next-backlog.md`),
   and current system health (which services are actually up, verified, not assumed).
4. **Clear next steps.** What the next session — automated or human-driven — should do first.

Store reports at `EMILY/continuity/YYYY-MM-DD-HHMMSS.md`, one per checkpoint, append-only — never
overwrite a prior report, they're the audit trail of the audit trail (Principle 9). Commit and push
EMILY last, same as any multi-repo change (Principle 12). Individual reports don't need golden-doc
registration (they'd bloat the index the way individual Apples don't get indexed either) — this
principle's existence in `THE_EMILY_WAY.md` (already tier-1 golden) is what makes the *process*
visible; the reports themselves are found in `EMILY/continuity/` the same way Apples are found in
IDUNA, not by being individually promoted to golden status.

### 15. Operational Health Is Not Optional

The founder's own words, verbatim, are the specification for this principle, added 2026-07-17 after
`emily-system.service` sat OOM-killed for ~22 hours and `secwatch`/`eps-reconciler` sat silently
dead for 10-22 hours — all three discovered by luck (a glimpsed tmux message, a manual audit
prompt), not by any system that was supposed to be watching:

> "we have been somewhat up for over 24 hours and are somehow just now noticing it — we need to
> write into the core prime directive to always check it, its critical to the S part of the S
> growth curve, we need 2 years of good data — we are going to go through some growth pains that
> the log streaming will help us write the story around but we need to get serious about
> operations. I know its hard being scrappy and bootstrapping, I know the founder sends off on side
> quests, but the growth of the ecosystem is not an enemy to fight."
>
> "fatbaby will pay the server bills for 10 years."

What this means in practice:

1. **Every session, before or during other work, verify infrastructure is actually up** — not
   assumed up because it was up last time. `emily status --fatbaby`, `systemctl --user status`
   on the core units, a log-freshness spot check. This is not a side task competing with "real"
   work; it *is* real work, on the same footing as any backlog item.
2. **Detection must be automatic, not lucky.** A human noticing an OOM-kill message on a phone
   screen is not a monitoring system. Every long-running process — HTTP-served or headless —
   needs a watchdog path that fires an escalation Apple + Slack alert on its own, on a timescale
   of minutes, not "whenever someone happens to look." (See `CheckServiceHealth` /
   `CheckPollerHealth` in `emily-agent/watchdog.go`, SECTION 152.)
3. **The founder's own side quests are a known, accepted cost of doing business here, not
   something to fight or resent** — bootstrapping means the person steering also gets pulled into
   GPT-2 training runs, DNS northstars, and Colab runbooks. The fix for the risk that creates
   (an unsupervised box getting OOM-killed mid-detour) is better automated monitoring and
   supervision (systemd auto-restart, not manual `go run`), not fewer side quests. Growth pains in
   the ecosystem are not the enemy; an *undetected* multi-hour data-ingestion gap is.
4. **The reason this matters more than almost anything else queued:** FatBaby's SEC/PR signal
   pipeline is a multi-year data asset — the founder's framing is 2+ years of continuous,
   trustworthy ingestion history, not a demo. A silent gap in that history is not a recoverable
   bug once the gap has passed; the data for those hours is simply gone. Treat any headless
   ingestion process the same way `CheckServiceHealth` already treats IDUNA/newssite/signalapi:
   as something whose downtime is unacceptable to discover after the fact.

### 16. The Inbox — Capture Now, Triage Later

Added 2026-07-23 after the founder shared a link with no task attached and said, plainly: "its
important i dont know save it somewhere... I'll dump stuff there if theres no clear intent
keeping full audit trails." Side quests are already an accepted cost of doing business here
(Principle 15) — this is the mechanism for the ones that arrive as a fragment, not a task.

The bar for `EMILY/inbox/` is deliberately zero. Nothing needs to be understood, scoped, or even
plausibly relevant before it lands there — the only requirement is that it lands *somewhere*,
timestamped, with the full raw content, instead of evaporating at the end of a conversation. Do
not block capture on figuring out intent first; that is the entire point of the principle.

This is upstream of, and distinct from, BACKLOG.md's `INTAKE QUEUE`: the intake queue is for
observations already curated toward becoming real backlog items (`emily backlog promote` routes
them through haiku classification). The inbox is for things that haven't been triaged at all —
some entries will graduate into a backlog item once their intent becomes clear; many won't, and
that's fine. Log format: `EMILY/inbox/INBOX.md`, newest dated section on top (same convention as
every `CHANGELOG.md`), one bullet per item — timestamp, raw content verbatim, whatever context
existed at capture time, and a `status` (`unprocessed` / `triaged` / `discarded`). Update the
status in place when an entry graduates or gets consciously dropped; don't delete the entry
either way — the audit trail is the point.

---

## The Feedback Loop

```
Emily Prime (plans)
  ↓ issues directed task → EMILY/signals/tasks/
Observation-watcher (dispatches)
  ↓ invokes claude --dangerously-skip-permissions [--continue]
Claude Code (implements)
  ↓ edits code, runs tests, commits, pushes
  ↓ writes claude-run report to PRRJECT_FATBABY/claude-runs/
Emily Prime (observes)
  ↓ reads observation, FABLE advisor reads full-system-context.md
  ↓ files Apple, updates backlog, issues next task
Repeat (RSI loop, 5-minute cron)
```

---

## The Document Hierarchy

```
Tier 1 Golden Docs (always in Emily Prime context)
  → northstars, executive summaries, integration specs, RSI engine spec
  → listed in EMILY/context/golden-docs-index.md

Tier 2 Important Specs (read before touching a subsystem)
  → agent protocol, tool specs, netcode contracts, platform architecture docs
  → listed in golden-docs-index.md with tier=2

Tier 3 Reference (fine-grained, rarely needed in context)
  → individual endpoint specs, migration files, fixture docs
  → not in golden-docs-index.md; found by reading code

BACKLOG.md (the load-bearing node)
  → cross-repo, git-authoritative, append-only
  → read before starting any work

GOLDEN.md (compressed backlog)
  → auto-generated by emily backlog compress
  → what Emily Prime reads in cron cycles (≤1200 tokens)

CHANGELOG.md (per-repo narrative)
  → human-readable dated history
  → updated on every meaningful change

EMILY/continuity/ (session checkpoints)
  → full sync sweep + changelog digest + in-flight/blocked/queued state + next steps
  → one dated file per checkpoint, append-only, not individually golden-indexed
  → see Principle 14
```

---

## Common Patterns

### Starting a new feature
1. Check BACKLOG.md — is there an item for this?
2. If not, add one with acceptance criteria and post `emily observe "new feature needed: X"`
3. Emily Prime will triage and promote via HEIMDAL if it fits the roadmap
4. Spec first: write the design doc, commit to `docs/`
5. Implement against the spec
6. Run tests
7. Post Apple
8. Mark `[x]` in BACKLOG.md
9. Update CHANGELOG.md
10. Commit everything (spec + code + backlog + changelog)
11. Push

### Fixing a bug
1. Write a failing test first (Go: `_test.go` file)
2. Fix the bug
3. Verify test passes
4. Update CHANGELOG.md
5. Post Apple
6. Commit and push

### Adding a golden doc
1. Write the doc
2. Append row to `EMILY/context/golden-docs-index.md`
3. Commit both (doc in its repo, golden-docs-index in EMILY)
4. Push EMILY last (it must have the correct index)

---

## What "Done" Looks Like

An item is done when:
- [ ] Code is written and tests pass
- [ ] Apple is posted to IDUNA (`emily apples post -t completion ...`)
- [ ] BACKLOG.md item is marked `[x]` with Apple ID
- [ ] CHANGELOG.md is updated in the affected repo(s)
- [ ] Everything is committed and pushed
- [ ] Any new golden doc is registered in golden-docs-index.md

If any of these are missing, it's not done. The Apple is the proof; the git push is the publication.
