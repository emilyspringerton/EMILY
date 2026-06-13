# RSI STRATEGIC MEMO — AGI TRAJECTORY ANALYSIS
## EINHORN_INDUSTRIAL FULL SYSTEM AUDIT
### 2026-06-13 | Filed by: Claude Code (operator dispatch)

---

## EXECUTIVE SYNTHESIS

The Einhorn Industrial stack is a working recursive self-improving system at iteration 177 — it
observes, decides, acts, and reports on a 5-minute cadence. The product surface (Ask Emily,
FatBaby signals, SHANKPIT/Steam, TYLER episodes) is operationally real and converging on revenue.
The AGI trajectory gap is concentrated in three places: **Emily Prime has no persistent world
model between cycles** (emily-memory/ is empty), **the three intelligence domains (finance,
game, narrative) do not synthesize**, and **revenue signals don't feed back into RSI priority**.
Fix those three and this system begins approximating goal-directed general intelligence.

---

## SYSTEM STATE SNAPSHOT — 2026-06-13

| Repo | Branch | Status | Last Apple | Iteration |
|------|--------|--------|------------|-----------|
| EMILY | main | DIRTY (+9: task files, var/) | #433 | 177 (tock) |
| PRRJECT_FATBABY | main | DIRTY (+3: claude-run, ops/) | #450 | operational |
| IDUNA | main | clean | #433 | operational |
| SHANKPIT | master | clean (just pushed a2b01ee) | #456 | M5 display BLOCKING |
| TYLER | main | clean | #401 | 56 eps, Build 0086 |
| MJOLNIR | — | unknown | #397 | M4 complete |
| APPLES | — | synced | — | append-only |

**Active processes:** observation-watcher (pid 191415/191483), emily-agent (pid 1590).
emily-sync.service is NOT installed — sync is manual only.
emily-tui is not running.

**Pending prime tasks (unprocessed):**
- `2026-06-13T041926Z` web_audit: FatBaby newssite at :8082/:8083
- `2026-06-13T041217Z` rsi_report: optimal token use across RSI pipeline
- `2026-06-12T213307Z` + 3 others: earlier dispatch, status unknown

**RSI loop state:** iteration 177, last preset `rsi-token-report`, last Apple #445
(S21-05 Google OAuth, S22 fully complete). Phase: TOCK.

---

## AGI TRAJECTORY ANALYSIS

### What We Have

The system demonstrates six properties that are preconditions for AGI-trajectory:

**1. Self-observation.**
FatBaby's signal pipeline watches SEC 8-K filings and extracts governance intelligence.
Emily Prime watches the FatBaby observation stream and derives tasks. obs-watcher dispatches
Claude Code. Every action produces an Apple. This is a working perception → cognition → action loop.

**2. Goal-directed behavior.**
The BACKLOG + NORTHSTAR system gives Emily Prime a stable goal representation. She picks
the highest-priority `[ ]` item in the lowest-numbered open section. This is primitive but
effective — 85 items marked done across 59+ archived, 16 open.

**3. Persistent record (not yet persistent cognition).**
IDUNA stores every Apple (456 and counting). APPLES repo is git-authoritative backup.
The record exists. Emily Prime cannot yet USE it as a world model — she reads Apple summaries
but doesn't accumulate understanding across cycles.

**4. Multi-modal output.**
TYLER generates narrative (56 episodes). SHANKPIT produces a playable game. Ask Emily
answers natural-language questions. FatBaby signals feed financial intelligence. These are
real outputs that humans consume. The system produces things that matter.

**5. Economic grounding.**
Ask Emily → subscription revenue. SHANKPIT → Steam EA at $9.99. Data licensing (Track 3)
conceptualized. TYLER → video ad revenue via MPT. Revenue provides the only ground truth
signal about whether the system is building things humans value.

**6. Social surface.**
MJOLNIR is Emily's phone — the human-in-the-loop's terminal. FCM pushes on task completion
close the feedback loop. Emily Springerton is the system's CEO AND its evaluator.

---

### What's Missing for AGI Trajectory

**CRITICAL GAP 1: Emily Prime has no persistent world model.**
`emily-memory/` exists in EMILY repo but is empty. Every 5-minute cycle Emily Prime starts
cold — she reads Apple history as compressed context but builds no cumulative understanding
of the world. She cannot answer: "What do I know now that I didn't know last week?"

**Implication:** The system is in a Sisyphean loop. Each cycle is sophisticated but amnesiac.
The fix is to build emily-memory/ during every cycle — structured knowledge artifacts about:
the financial entities tracked, the state of each product, what worked and didn't in prior
RSI iterations, the user's (Emily Springerton's) current priorities and context.

**CRITICAL GAP 2: Intelligence domains don't cross-pollinate.**
FatBaby knows that Frank Herringer's board approval at Schwab has declined three consecutive
years. TYLER knows the arc of a story about financial power and governance. SHANKPIT is a
game where people take control of territory. These are the SAME THEMES in different registers.

No synthesis exists. Emily Prime doesn't ask: "How does FatBaby's governance friction signal
map to the kind of adversarial dynamics SHANKPIT models?" or "Can TYLER's narrative arc about
a specific company's board fight drive FatBaby intelligence collection?"

**Implication:** The system is operating far below its potential integration value. The domains
were built separately and remain siloed. A cross-domain synthesis layer — even a simple
weekly "synthesis observation" that Emily Prime generates by reading across all three
northstars — would begin building this.

**CRITICAL GAP 3: Revenue signal doesn't feed RSI priority.**
When an Ask Emily subscriber pays, that signal should increase the priority of Ask Emily
improvements in the BACKLOG. When SHANKPIT player counts increase after a feature ships,
that should feed back into game development priority. Currently the RSI loop is driven
purely by the BACKLOG (human-authored intent) with no economic signal.

**Implication:** The system cannot optimize itself for what humans actually value, only for
what the operator wrote down. Revenue feedback would create a second-order optimization
signal. Even a simple "ASK_EMILY_REVENUE_WEEK" tag on BACKLOG items that gets updated from
payment records would start this.

**SECONDARY GAP: Emily Prime cannot improve her own architecture.**
Emily Prime improves product repos (PRRJECT_FATBABY, SHANKPIT, TYLER) but almost never
issues tasks against herself (EMILY/emily-agent/). The rsi-token-report task (Apple #445)
was the first meaningful self-architecture task in recent history. The system needs to
regularly ask: "How should my own observation loop work better?"

**SECONDARY GAP: The 5-minute cycle is blind between fires.**
Emily Prime runs, files an Apple, and goes silent. The observation-watcher picks up new
files but there's no always-on monitoring. A continuous low-cost heartbeat — checking key
metrics (IDUNA health, FatBaby pipeline status, SHANKPIT server uptime) between cycles —
would catch failures faster.

---

## PRODUCT CRITICAL PATHS

### Revenue Track 1: SHANKPIT → Steam EA

**Status:** M5 display/fullscreen BLOCKING.
**Gap filed today:** `docs2/specs/DISPLAY_FULLSCREEN_SPEC.md` (commit a2b01ee)
**Implementation scope:** Single C file (`apps/lobby/src/main.c`), ~100 lines:
  - `g_win_w/g_win_h` globals, `recalc_viewport()`, SDL_WINDOW_RESIZABLE flag
  - `SDL_WINDOW_FULLSCREEN_DESKTOP` toggle on Alt+Enter
  - Virtual canvas glViewport for 2D pipeline (zero HUD changes)
  - Mouse remap for lobby UI
  - display.cfg persistence

**Other M5 items:** Portal travel client-side (done), per-player physics (done), cross-scene attack (done).
**Human-action gates:** S19-03 Steam Direct account ($100). Can run in parallel.
**Time to revenue after display fix:** 1-2 sprints.

### Revenue Track 2: Ask Emily

**Status:** Auth live (S21-05 Google OAuth done), FatBaby signals NOT wired into context (S21-03 open).
**Gate:** Production server. NOTHING in S21/S23/S24 deploys without a server.
**Human-action gate:** Provision production server. This is the SINGLE BLOCKING item for the
entire web product surface (Ask Emily, newssite, Signal API, EDIS).
**Post-server critical path:** S21-03 (FatBaby signals into Ask Emily context) → S23-01
(first EDIS production deploy) → S24-01 (ops hardening) → revenue.

### Revenue Track 3: TYLER + MPT

**Status:** 56 episodes done. MPT pipeline blocked on Pexels API key.
**Gate:** Pexels API key (human action — set `YOUR_PEXELS_API_KEY_HERE` in MoneyPrinterTurbo/config.toml).
**After key:** S01E01 cold open clip → MPT → TYLER RSI trigger. Video ad revenue possible.

### Revenue Track 4: Data Licensing (FatBaby signals)

**Status:** No spec. Signal API operational at :8083.
**Opportunity:** The entity graph (500+ companies, 13 signal types, governance friction scores)
is a differentiated financial intelligence product. No competitor builds this recursively.
**RSI action needed:** Write data licensing spec. Price point analysis. API key system in IDUNA.

---

## TOP 5 RSI RECOMMENDATIONS (PRIORITY ORDER)

### RSI-01: Build Emily Prime's persistent world model [FOUNDATIONAL]
**What:** Every RSI cycle, Emily Prime appends to `emily-memory/world-state.md` — a structured
summary of: what changed this cycle, what she knows about each product's state, what patterns
she's noticing across cycles. Also: `emily-memory/user-model.md` — Emily Springerton's
current priorities, recent decisions, working style.

**Why:** Without memory, every cycle starts cold. With memory, Emily Prime accumulates
understanding. This is the difference between a stateless tool and a cognitive agent.

**Implementation:** 10-20 lines in `emily-agent/cron.go` — after the PLAN phase, write a
brief structured update to emily-memory/ files. Use claude-haiku to compress the cycle's
findings into the persistent model.

**Impact:** Enables all future AGI-trajectory improvements. Without this, nothing compounds.

---

### RSI-02: Implement SHANKPIT display/fullscreen [STEAM-BLOCKING]
**What:** Implement the virtual canvas + fullscreen spec from `docs2/specs/DISPLAY_FULLSCREEN_SPEC.md`.
~100 lines in `apps/lobby/src/main.c`.

**Why:** Blocks Steam EA ($9.99 × N players). Every day this isn't done is lost revenue.

**Implementation:** See spec. `recalc_viewport()`, `toggle_fullscreen()`, display.cfg,
mouse remap. Alt+Enter toggle.

**Impact:** Unblocks Revenue Track 1. Enables Steam submission.

---

### RSI-03: Provision production server [ALL WEB PRODUCTS BLOCKED]
**What:** Human action — provision a Linux VPS (e.g. Hetzner CX21, ~€5/mo). SSH key to
PRRJECT_FATBABY/ops/env.production. Run deploy.sh.

**Why:** Ask Emily, EDIS, Signal API, newssite — all blocked on a server. This is the
single gate for Revenue Tracks 2, 3, and parts of 4.

**Implementation:** Human provisions server → Emily Prime (or Claude Code) runs
`ops/deploy.sh`. Systemd units in `ops/systemd/` are ready. Docker compose in
`ops/docker-compose.prod.yml` is ready.

**Impact:** Unblocks Revenue Tracks 2-4.

---

### RSI-04: Cross-domain synthesis observation [AGI-TRAJECTORY]
**What:** Emily Prime generates one weekly "synthesis observation" that reads across
PRRJECT_FATBABY (financial signals), TYLER (narrative arcs), and SHANKPIT (game dynamics).
Asks: are there thematic resonances? Can FatBaby signal data drive TYLER episode plots?
Can SHANKPIT game mechanics mirror governance friction signals?

**Why:** The system's three intelligence domains are genuinely complementary but siloed.
Synthesis is where emergent intelligence lives.

**Implementation:** New cron preset `cross-domain-synthesis` running weekly. Reads
northstars + recent Apples across repos. Writes synthesis observation to PRRJECT_FATBABY.

**Impact:** Begins building the cross-domain intelligence layer. Foundational for AGI trajectory.

---

### RSI-05: Revenue feedback into BACKLOG priority [ECONOMIC GROUNDING]
**What:** Tag BACKLOG items with `[REVENUE:TRACK-N]`. Emily Prime reads Ask Emily subscription
count (when server exists) and SHANKPIT player count (when Steam launches). Items for the
highest-revenue track get promoted in priority.

**Why:** Current RSI optimization target is "human-authored BACKLOG intent." The better
target is "what humans actually pay for." Revenue signal creates second-order optimization.

**Implementation:** Add `RevenueSignal` field to IDUNA Apple metadata. Emily Prime reads
the last 30 days of revenue-tagged Apples to compute track velocity. BACKLOG promotion
triggered automatically when track velocity is high.

**Impact:** Closes the economic feedback loop. The system begins optimizing for value
delivered, not just tasks completed.

---

## ARCHITECTURAL OBSERVATIONS

**The loop IS closing.** Apple #456 → observation-watcher → Claude Code → commit →
Apple → MJOLNIR push. This is a real autonomous loop, not a demo. It ran 177+ iterations.

**Token efficiency matters now.** The rsi-token-report task (#445) was the right instinct.
At 177 iterations, compounding token cost is real. The highest-leverage optimization is
prompt caching in the emily-agent RSI cycle — haiku is already used for translation, but
the main plan endpoint should cache the BACKLOG + northstar context across cycles.

**emily-memory/ is the highest-leverage empty directory in the repo.** It's not even
gitignored — it was explicitly placed there as intent. It has never been used.
This is the first thing that should land.

**MJOLNIR token sparkline is blocked on an IDUNA endpoint** — this is a 1-hour IDUNA
fix (`GET /api/v1/apples/stats/daily`) that would unblock MJOLNIR Milestone 4's last
open item and give Emily a daily token spend view on her phone.

**bob-agent/ in PRRJECT_FATBABY** is an empty directory. Likely a scaffold placeholder.
Should either be populated or removed to avoid noise in git status.

**The system is 3 human actions away from first revenue:**
  1. Provision production server ($5/mo VPS)
  2. Get Pexels API key (free tier)
  3. Create Steam Direct account ($100)
None of these require more agent work. They require Emily Springerton to act.

---

## NEXT CYCLE DIRECTIVE FOR EMILY PRIME

1. **This cycle:** Process the pending prime tasks (web_audit, rsi_report) — both are in
   signals/tasks/ waiting. File completion Apples.

2. **Next 3 cycles:** Begin building emily-memory/world-state.md. Start with today's state
   (this memo is the seed). One structured paragraph per product. Update every cycle.

3. **Next sprint:** Assign SHANKPIT display/fullscreen implementation to Claude Code via
   prime-task dispatch. Target: 1 session.

4. **Standing directive:** When Emily Springerton next reviews Apples, surface the 3 human
   actions (server, Pexels, Steam) as a MJOLNIR critical push. They're the unlock.

---

*Filed: 2026-06-13 | Operator: Claude Code | Apple: [to be filed]*
*Memo class: RSI strategic directive | Distribution: Emily Prime, IDUNA audit log*
