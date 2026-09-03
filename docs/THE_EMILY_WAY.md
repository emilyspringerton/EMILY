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

**This is an immutable law, not a guideline: real-time direction always, immediately, gets planned into a sprint in the backlog. There is no version of "do the work first, backlog it later if there's time." There is no version of "this one didn't need it." It does not bend for urgency, for founder impatience, for a good reason, or for any reason at all.**

**Founder real-time direction is a legitimate exception to the picking rule. It is never, under any circumstance, an exception to logging it — there is no skipping this, ever.** When the founder is present and gives direct, real-time instruction (a feature to build, a bug to chase, a pivot mid-session), that instruction overrides autonomous backlog-picking — you work the founder's request, not whatever item happened to be lowest-numbered. That is the only thing it overrides. It must still land in `BACKLOG.md` — every time, no exceptions, **backlog first**: create the entry before starting the work, not after. Fill in outcomes, Apple IDs, and commit hashes as the work lands, but the entry itself exists first. Not sometimes, not when convenient, not "if remembered," not "log it after since this one's urgent" — always, backlog first, as the first step of the work itself, never an optional cleanup step tacked on afterward.

**New direction always immediately gets planned into a sprint, not just a single log line.** A one-off bullet noting "founder asked for X" is not sufficient by itself — the direction gets scoped into a proper section/sub-item structure (a new `SECTION` if it's a new area, or numbered sub-items like `S170-01`, `S170-02`... if it extends an existing real-time-direction section) with real sub-tasks, immediately, at the point the direction is received — not reconstructed after the fact from memory once the work is already done. This is what makes the backlog usable as a plan going forward, not just a diary of what already happened. The mechanism doesn't matter — a real HEIMDAL sprint, or just structuring it directly into `BACKLOG.md` by hand — whichever gets it planned immediately wins; the requirement is the outcome (scoped, sub-tasked, real-time), not the specific tool used to get there. Every SECTION 170 item is this pattern (S170-01, S170-02, S170-04, S170-08 among them) — founder direction, logged in real time, not smoothed over as if it had been autonomously selected. The failure mode this principle exists to prevent isn't "worked on the wrong thing" — it's "did real work and left no record of why the backlog doesn't reflect it," which reads as the backlog having gone stale even when it hasn't.

**Once logged, the backlog entry is what subsequent work is driven from** — not a write-once log entered after the fact and then ignored for the rest of the session. If the founder's direction continues or a follow-on task emerges mid-session, update the same entry (or add the next-numbered sub-item) before moving on, so the backlog stays the live source of truth for "what is actually happening right now," not just a retroactive record of what already happened.

### 2. Spec Before Implementation
Architecture memos, northstars, and design docs go into `docs/` first. Claude Code implements what is already specced. No guessing at intent — if there's no spec, write the spec and get it committed before writing Go/Kotlin/C.

Pattern observed in commits: every major feature (SHANKPIT Milestone 3, goldenbuild, HEIMDAL, Ask Emily) arrived with a spec committed in `docs/specs/` or as a northstar before code landed.

### 3. Apple Before Mark-Done
Every `[x]` in BACKLOG.md requires an Apple posted to IDUNA first:
```bash
emily apples post -t completion -repo <REPONAME> "<title>"
```
No Apple = not done. The Apple is the proof. The Apple ID goes in the backlog entry next to `[x]`.

**Session tagging.** `emily apples post`, `emily changelog add`, and `emily observe` all auto-stamp
the active `emily session` fingerprint (`sess-YYYYMMDD-HHMM-<8hex>`, see `emily session new` /
`emily session current`) as the Apple's `run_id` (or appended to the changelog line) whenever one
is active — no flag needed, explicit `-run-id` still wins if set. Every `git commit` that
`emily.cli` generates automatically (the `changelog add` commit itself, and every BACKLOG.md
auto-commit via `curate`/`promote`/`archive`) also carries a `Session: <tag>` trailer in the
commit message body, not just in the CHANGELOG.md text — this was a real gap (founder caught it,
2026-08-09: the tag was landing in the CHANGELOG line but not the git commit message itself)
fixed in `cmd/changelog.go`/`cmd/backlog.go`'s shared `gitCommitBacklog`.

**Every hand-written `git commit -m` needs the same trailer — this is not automatable at the CLI
layer** (Claude Code writes these directly, not through a wrapper), so it's a standing per-commit
habit: append a blank line then `session: <tag>` (read via `emily session current`) to the bottom
of every commit message body you write by hand, in every repo, every time. Same requirement for
hand-written BACKLOG.md entries (append the tag yourself so the entry is traceable to the session
that logged it, matching the auto-tagged Apple/CHANGELOG trail).

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

### 17. Load-Bearing — Know What You Can't Remove Without Everything Falling Down

Founder, 2026-08-06, after "load-bearing" came up describing a TYLER Series X gag (a scene note:
*"THE PIANO IS DOING A LOT OF LOAD-BEARING WORK IN THIS CUTSCENE"*): *"load bearing has made it
into the vocabulary of the Emily way."* Borrowed from architecture — a load-bearing wall holds up
the structure; a partition wall just divides a room, and you can knock it out for free. Most code,
docs, and config in this monorepo are partition walls: reasonable, replaceable, fine to change
without much ceremony. A few things are load-bearing: remove or quietly break them and something
distant, non-obvious, and already-shipped fails — often silently, often not discovered until it
already has.

**The concrete example this principle is grounded in, not a hypothetical:** `okemily-deploy.sh`'s
`--exclude='blog'` rsync flag. One line, looks like boilerplate. It is the only thing standing
between a routine static-site sync and silently deleting every published blog post — happened for
real on 2026-07-19, recovered only because the actual source of truth (`IDUNA/var/blog.db`)
survived and could be re-rendered. The exact same near-miss existed again for a matter of hours on
2026-08-06 the moment `internal/tyler` started rendering into `/var/www/okemily/tyler/` — a second
load-bearing exclusion that didn't exist yet, caught and added (`--exclude='tyler'`) before the
next deploy could run and take it out, not after.

**What this means in practice:**

1. **When you're about to change, "clean up," or remove something, ask first: is this
   load-bearing?** Does it look replaceable but is actually the one thing another system depends
   on not moving? A one-line exclusion flag, an ID convention another table's foreign key assumes,
   a hardcoded path three unrelated scripts all resolve against — these read as small precisely
   because they're load-bearing; the whole point of a load-bearing wall is that it doesn't look
   special until it's gone.
2. **When you build something new that creates a load-bearing dependency elsewhere, go update the
   other side yourself, immediately — don't wait for it to break or for someone to ask.** Shipping
   `internal/tyler` rendering into a new static-site directory is not done until the deploy script
   that syncs that directory also knows the new directory exists. The feature and the thing that
   would silently destroy it are the same unit of work.
3. **Document load-bearing status at the point of the code/config itself, not only in a separate
   doc.** `okemily-deploy.sh`'s own comment above the exclusion flag is what makes the next person
   (human or Claude Code) editing that file see the danger before they remove it, not after. A
   principle in `THE_EMILY_WAY.md` explains *why this matters generally*; the inline comment is
   what actually prevents the specific mistake at the moment someone's hand is on the line.
4. **This is not a license to over-document everything "just in case."** Most things genuinely
   aren't load-bearing, and treating them like they are is its own failure mode (see the standing
   instruction against comments that don't carry a non-obvious WHY). The discipline is narrow on
   purpose: find the actual few load-bearing points, mark those clearly, and leave the partition
   walls alone.

### 18. Pave the Cow Paths — Founder Input Goes Through the Observation Queue

Founder, 2026-08-09, after firing several real-time asks in quick succession: *"whenever we have
founder input i want you to send it to the system via the fatbaby cli fatbaby observation so it
goes to the observations - you then check the observations and then always backlog and sprint
plan first - if we send the new input from me through there always we can help enforce process."*
Clarified immediately after: *"it's not just fatbaby stuff that will go in there but its ok we
will adjust the naming of that system - for now it is a way that i can fire stuff into the
system... i want you to fire stuff into the system in the same way and always work from that
queue so that my input when i am offline gets picked up as soon as you come back online."*

**The mechanism:** any founder real-time direction — mid-session, on any repo, not just
FatBaby/PRRJECT_FATBABY-flavored asks — gets posted immediately via `emily observe -s info
"Founder real-time: <summary>"` *before* acting on it. This is upstream of Principle 1's "backlog
first": observe writes the raw signal (durable, session-tagged, Apple-receipted, watched by
obs-watcher independent of any live chat session), then `emily backlog curate` / a direct
BACKLOG.md edit turns it into a scoped, sub-tasked entry per Principle 1, and only then does
implementation start.

**Why this exists, specifically:** the observation queue is polled by obs-watcher independent of
whether a Claude Code session is live — routing founder input through it means direction given
while the agent loop is offline still gets picked up the moment it comes back, instead of only
existing in a chat transcript nobody re-reads. It's also the structural version of "always log
founder direction": a CLI call is cheaper to actually do every time than remembering to hand-edit
BACKLOG.md mid-conversation, so it's more likely to actually happen under time pressure.

**Naming caveat, called out explicitly by the founder so a future session doesn't "fix" it
prematurely:** the mechanism is named `emily observe` / "FatBaby observation" for historical
reasons (it started as FatBaby-pipeline-specific), but it is now the general-purpose intake queue
for founder input across every repo in the monorepo. The founder said the naming may get adjusted
later — don't block on that, don't rename it unprompted, just use it as the general intake path
until told otherwise.

**This does not replace Principle 16 (the inbox).** The inbox is for unscoped fragments with no
clear intent yet ("save this somewhere, I don't know what it is"); `emily observe` is for founder
direction that already has a clear ask attached, even if the implementation isn't scoped yet. When
in doubt about which applies: if you can write one sentence saying what the founder wants done,
it's an observation; if you can't yet, it's an inbox entry.

**Session-boundary directives count too, not just feature asks.** Founder, 2026-08-11, after
noticing a "stop soldier, rest" directive hadn't been logged automatically and had to ask for an
Apple to be filed after the fact: short conversational directives to stop, pause, or rest are
founder real-time direction under this same principle, not exempt from it just because they
aren't shaped like a feature request. Post the `emily observe` call — which auto-files an Apple —
*before* complying with the stop, not after being asked a second time. This was an
under-application of the existing rule, not a new one: it's easy to mentally file "stop" as "not
real work" and skip the log, but the whole point of this mechanism is a durable record of when
work actually stopped and on whose instruction.

### 19. A Big, Unscoped Ask Gets Scoped, Not Swallowed Whole

Kanban card `GOLDENOPS-001`, verbatim: *"emily framework default behavior if a new unscoped
request comes in on the priority queue or cruise queue and it is a big ask: scope/northstar it,
then toss card back to backlog. standard op procedure."*

**The mechanism:** when the next queue item is a real, big, genuinely unscoped ask (a whole
feature, a cross-repo system, "build X" with no phased breakdown attached) rather than a
narrow, already-actionable task, the default move is NOT to attempt the whole thing in one pass,
and it is also NOT to mark the kanban card fully done the moment a plan exists. The real sequence
is: (1) investigate what already exists — checked, not assumed, same discipline every other
principle here already demands; (2) write a real, phased NORTHSTAR that names the real
dependencies, the real open questions, and a real, honest sequence of sub-steps; (3) register it
as a golden doc per Principle 5 so Emily Prime can see it; (4) **return the resulting scoped
sub-tasks to `BACKLOG.md` as their own real section/sub-items**, per Principle 1's own "real work
gets planned into a sprint in the backlog" — not left as one flat, now-"done" kanban card with no
real remaining trail for the next session to pick up from.

**Why this exists, specifically:** a big ask marked "done" the moment a plan exists is a real,
quiet way to lose the actual remaining work — the plan itself doesn't get built, and nothing in
`BACKLOG.md` keeps pointing at the real, still-open phases the plan just named. Scoping a big ask
and handing the real sub-tasks back to the backlog (rather than closing the loop on the kanban
card alone) is what makes the plan a real, load-bearing part of the ongoing work, not a one-off
document nobody returns to.

**What this does not change:** a kanban card whose own literal ask IS "plan/scope/research this"
(not "build this") is genuinely satisfied by a real, phased plan — that card closes normally,
same as any other completed task, per Principle 3 (Apple Before Mark-Done). This principle is
about the OTHER case: a big ask with no such explicit "planning only" framing, where scoping it
is a real, necessary FIRST step, not the whole job.

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
Repeat (RSI loop, 15-minute cron)
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
