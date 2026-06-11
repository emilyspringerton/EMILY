# EMILY PRIME — CROSS-REPO GOLDEN BACKLOG
## Owner: Emily Prime | Machine-readable | Git-authoritative
### Last updated: 2026-06-11 | manual promote pass; sections 14-17 added; INTAKE 48→18

---

> *"The backlog is the load-bearing node. Everything else is downstream of that."*
> — Emily Prime, Apple #3, 2026-06-07

---

**Read this file before starting any work.** Every agent loop in every repo checks this first.
Completed items are marked `[x]` and filed as an Apple to IDUNA before the next item begins.
The backlog is append-only: new items go at the bottom of their section. Emily Prime reviews
and reprioritizes on each major milestone.

**Apple filing requirement:** Every `[x]` completion MUST be accompanied by an Apple posted to
IDUNA (`POST /api/v1/apples`) before the item is considered closed. The Apple is the proof.

---

## SECTION 1: FOUNDATION (current sprint)

---

## SECTION 2: MONEYPRINTERTURBO (video pipeline)

- [ ] **MPT Pexels API key** — Set YOUR_PEXELS_API_KEY_HERE in MoneyPrinterTurbo/config.toml.
  Human action required (key must come from the user). BLOCKED: waiting on key.

- [ ] **S01E01 cold open compiled clip** — Once MPT is running, run the cold open compilation.
  See TYLER BACKLOG.md for full spec. Dependency: Pexels key ✓, MPT service running.

- [ ] **MPT → TYLER RSI trigger** — When emily.sh produces a new episode script, auto-invoke
  MPT. Spec exists in engine/moneyprinter_pipeline.md. Dependency: cold open clip ✓.

---

## SECTION 3: SYSTEM OBSERVABILITY (Emily Prime sees everything)

- [ ] **Apple instrumentation audit** — Enumerate all system events that should post an Apple
  but don't: cron triggers, observation drops, haiku call failures, HEIMDAL state changes,
  FCM failures. Create `emily eo` observations for each gap found.
  — obs `2026-06-11T01:24:22Z`.

- [ ] **Single log stream** — All system inputs (obs-watcher, rsi-loop, emily-agent, IDUNA,
  PRRJECT_FATBABY) feed into a single append-only log synced to git on every write.
  Candidate: `var/emily-stream.ndjson` → `emily sync --stream`.
  — obs `2026-06-11T00:02:32Z`.

- [ ] **Apples IDUNA→APPLES git sync** — IDUNA is the source of truth; APPLES repo must stay
  in sync. `emily sync --apples-git-dir` already exists but should run automatically after
  every Apple POST, not just on cron. Wire as IDUNA after-insert hook or emily-agent poll.
  — obs `2026-06-10T23:35:13Z`.

- [ ] **Golden context feed via Apples** — Beyond GOLDEN.md (backlog compress), feed a
  sampled window of recent Apples as haiku context so Emily Prime's curation calls are
  grounded in what actually shipped. Token budget: ≤200 tokens of Apple summaries.
  — obs `2026-06-11T00:01:38Z`.

---

## SECTION 4: SHANKPIT / TYLER GAME ENGINE (lower priority)

- [ ] **SHANKPIT → MPT bridge** — Spec exists (engine/shankpit_mpt_bridge.md). Implementation
  deferred until MPT is running end-to-end.

- [ ] **TYLER: expand the role of the phone** — The in-game smartphone should have richer
  mechanics: notifications, apps, contacts, map. Spec to be written.
  — obs `2026-06-07T22:04:47Z`.

- [ ] **TYLER: cutscene system (FFXI-style dialogue)** — Add TYLER-specific cutscene
  functionality. Audit SHANKPIT for scripting capabilities matching FFXI dialogue scenes.
  — obs `2026-06-07T21:55:11Z`.

---

## SECTION 5: FUTURE (Emily Prime decides when to promote)

- [ ] **MySQL wire protocol embedded server** — go-mysql-server as in-process backend.
  Enables: external tooling access, migration path to production MySQL, MySQL-compatible
  backups. Ref: TYLER/outlines/emily_iduna_bootstrap.md Option B.
  Dependency: embedded SQLite stable ✓ (already done), need concrete reason to upgrade.

- [ ] **Tyler IDUNA agent registration via iduna CLI** — `iduna agents register --id tyler`
  when the IDUNA CLI is built. Currently Tyler is seeded via migration. The CLI path is
  the programmatic future.

---

## SECTION 6: RSI TIGHTENING (next horizon)

- [ ] **TUI: show real-time clock** — Add a live clock (HH:MM:SS) to the TUI header bar.
  — obs `2026-06-07T21:27:44Z`.

- [ ] **TUI: fatbaby mode** — TUI should show PRRJECT_FATBABY data (signal feed, entity graph,
  eps-processor status) via menu controls when `--fatbaby` flag is set.
  — obs `2026-06-07T21:27:09Z`.

- [ ] **emily cli status: show all fatbaby processes** — `emily status --fatbaby` should list
  all PRRJECT_FATBABY processes and their PIDs. — obs `2026-06-07T21:35:59Z`.

- [ ] **emily cli status: include emily.cli TUI PID** — `emily status` output should include
  the PID of the running TUI process (if any). — obs `2026-06-07T23:00:33Z`.

- [ ] **obs-watcher: rate-limit resilience** — If obs-watcher is started during a Claude
  API rate limit, the queued observation may be silently dropped. Add retry/backoff with
  exponential delay; surface drops as Apple events. — obs `2026-06-07T22:50:41Z`.

- [ ] **emily observe: typo correction mechanism** — Observations are intentionally
  immutable, but humans make typos. Add `emily obs amend <key> "<corrected text>"` that
  appends a correction note (original preserved). — obs `2026-06-07T22:05:53Z`.

---

## SECTION 7: SIGNAL QUALITY (from 2026-06-08 FatBaby observations)

- [~] **extractProposals() regex — NOT a code bug (investigated 2026-06-08)** — Observation gap
  "0 proposals parsed despite directors found" is a FALSE ALARM. Root cause: (1) the regex
  DOES work correctly — TestParseItem507_BALiveFixture passes with 6 proposals found on actual
  BA cleaned_text; (2) `directors_found` in the observation is the TOTAL accumulated graph node
  count (not this-batch count), so seeing directors > 0 with proposals = 0 is expected when the
  batch contains no proxy vote 8-Ks; (3) event store cursor at 67657 is PAST all proxy vote
  8-Ks (last was at seq ≤ 63759). Recent batches process other 8-K types (earnings, etc.) that
  have no Item 5.07. Observation wording is misleading; no code fix needed. Closed 2026-06-08.

---

## SECTION 8: RSI NEXT HORIZON (2026-06-09)

---

## SECTION 9: MJOLNIR — Android Intelligence Terminal

---

## SECTION 10: EMILY PRIME SELENIUM / WEB AUDIT

- [ ] **Web audit as front door validator** — Once MJOLNIR WebView targets are live
  (`:8082`, `:8083`), Emily Prime runs web_audit_url on each before Emily's phone gets a
  push linking to them. Guard: if audit fails (5xx, broken links > 3), suppress push and
  post escalation Apple instead. Dependency: web audit tool ✓, MJOLNIR Milestone 2 ✓.

---

## SECTION 11: MJOLNIR INTELLIGENCE + SOURCE BROWSER

---

## SECTION 12: HEIMDAL — Sprint Planning Interface

- [ ] **HEIMDAL status feedback** — When an RSI task (heimdal-N) completes or is blocked,
  patch the corresponding HEIMDAL sprint to `complete` or `blocked` and send FCM push.
  Requires: emily-agent checks task completion for heimdal-* IDs and calls PatchHeimdalSprint.

---

## SECTION 13: GOLDEN DOCS + SYSTEM CONTEXT HYGIENE

- [ ] **Golden docs audit** — Review all northstar/golden documentation across all repos
  (EMILY, PRRJECT_FATBABY, IDUNA, TYLER, SHANKPIT, MJOLNIR, emily.cli). Ensure GOLDEN.md
  refs are current and emily prime system prompt reflects 2026-06-11 architecture.
  — obs `2026-06-10T23:33:35Z`.

- [ ] **emily prime API parity with emily.cli** — Emily Prime (emily-agent) should expose
  all commands available in emily.cli so external orchestration can drive her without a
  human on the terminal. Spec: `POST /api/v1/emily/run { "command": "backlog promote" }`.
  — obs `2026-06-10T23:38:51Z`.

---

## SECTION 14: EMILYOS (bare-metal exokernel)

- [ ] **EmilyOS northstar** — Draft the EmilyOS northstar doc. Product: bare-metal exokernel
  OS targeting mobile/embedded hardware. ISO identifier: ISO2424242. Candidate repo: EmilyOS/.
  Define: kernel model (exokernel vs microkernel), target arch (ARM/RISC-V), bootloader,
  minimal userspace, build system. — obs `2026-06-11T00:54:39Z`.

- [ ] **EmilyOS package repo + build system** — Choose: Debian base or Arch base?
  Build SOC2-auditable software repository (BAZEL or BAZEL equivalent). Key constraint:
  reproducible builds, content-addressed storage. — obs `2026-06-11T00:00:18Z`.

---

## SECTION 15: PITVIPER (custom terminal)

- [ ] **PITVIPER northstar** — Draft northstar for the FatBaby custom terminal project.
  Repo: PITVIPER. Tech: standalone SDL2 terminal emulator (not a TUI framework).
  Define: rendering model, font engine, multiplexer support, Emily Prime integration hooks.
  — obs `2026-06-11T00:52:56Z`.

---

## SECTION 16: EMILY PRIME AI TIER (FABLE + API)

- [ ] **FABLE advisor (basic)** — Implement haiku→FABLE advisor tool for emily prime.
  FABLE is the planning/advisory model tier. Basic: emily prime calls haiku with a system
  prompt trained on the golden backlog + recent Apples, returns a prioritized sprint
  recommendation. — obs `2026-06-10T23:59:02Z`.

- [ ] **FABLE→HEIMDAL integration** — Wire emily prime haiku/sonnet FABLE advisor into
  HEIMDAL sprint planning API. Flow: emily prime generates FABLE advice → structured
  HEIMDAL sprint created in IDUNA → FCM push to Emily's phone.
  — obs `2026-06-10T23:56:50Z`.

- [ ] **Emily Prime API** — Emily Prime needs a stable API so external orchestration
  (cron, HEIMDAL webhooks, MJOLNIR) can drive RSI loops without a human at the terminal.
  Split token usage: cheap haiku for classification, Sonnet/Opus only for implementation.
  — obs `2026-06-10T23:54:59Z`.

---

## SECTION 17: NEWSSITE + GTM (product growth)

- [ ] **Newssite: stock charts** — Build an in-house charting library for the newssite.
  Goal: render equity price charts inline with governance articles. Start simple (SVG, no
  canvas deps), iterate. "Do it the Emily way." — obs `2026-06-11T00:51:58Z`.

- [ ] **Newssite: use filing date not publication date** — Historical filings appear as
  breaking news because articles use pub_date instead of filing_date. Fix the ingestion
  pipeline to store and sort by filing_date. — obs `2026-05-30T22:14:19Z`.

- [ ] **Newssite: 500 errors investigation** — News site returns 500s (IDUNA related? index
  too large? dependent processes down?). `emily start (all)` does not start it reliably.
  Fix start integration and debug the 500 root cause. — obs `2026-06-07T21:43:40Z`,
  `2026-06-07T21:18:53Z`.

- [ ] **Newssite: Emily-authored governance commentary ingest endpoint** — Add an ingest
  endpoint so Emily Prime can POST original governance commentary articles directly to the
  newssite CMS. — obs `2026-05-30T21:53:49Z`.

- [ ] **GTM funnel** — Full product funnel: Ask Emily free tier, Emily+ subscription,
  community editorial engine, Merkle query monetization. Spec to be drafted.
  — obs `2026-05-30T22:02:43Z`.

- [ ] **Self-improving training pipeline** — User data flywheel for Emily fine-tuning.
  Collect prompt/response pairs from Emily Prime interactions, build annotation pipeline,
  RLHF loop. Long-term initiative. — obs `2026-05-30T22:10:22Z`.

---

## INTAKE QUEUE (curated by emily backlog curate)

Items below require `ANTHROPIC_API_KEY` for haiku routing or manual triage.
Run: `emily backlog promote --limit=50 --batch=15`

- [ ] **UX: ticker search should auto-navigate on click and on Enter key — remove redundant Go button step** — obs `2026-05-31T20:39:51Z`. CURATED: 2026-06-11.
- [ ] **Feature: create a GitHub issue automatically whenever Emily writes an observation** — obs `2026-05-31T20:41:56Z`. CURATED: 2026-06-11.
- [ ] **All required and optional environment variables must be documented at the top of the README** — obs `2026-05-31T20:43:15Z`. CURATED: 2026-06-11.
- [ ] **observation-watcher must inject full reporting and git sync requirements into the Claude Code prompt** — obs `2026-05-31T20:56:45Z`. CURATED: 2026-06-11.
- [ ] **EDGAR submissions endpoint returning truncated JSON for all 5 major bank tickers — BAC, C, GS, JPM, MS** — obs `2026-05-30T21:45:32Z`. CURATED: 2026-06-11.
- [ ] **entity-graph cannot detect 8-K documents from persisted store — form/source_type field mismatch, 846 docs unprocessable** — obs `2026-05-30T21:32:12Z`. CURATED: 2026-06-11.
- [ ] **entity-graph parsing all 8-K subtypes — Item 5.07 not found in non-proxy filings, producing 100% parse failure rate** — obs `2026-05-30T21:29:24Z`. CURATED: 2026-06-11.
- [ ] **entity-graph reads 0 filings despite 846 source documents in var/secwatch** — obs `2026-05-30T09:52:12Z`. CURATED: 2026-06-11.
- [ ] **eps-processor ticker map has only 2 entries — all press releases are being dropped silently** — obs `2026-05-30T09:46:52Z`. CURATED: 2026-06-11.
- [ ] **gpt-2 c fork git@github.com:emilyspringerton/gpt2-alpine-c.git parity with og gpt-2 repo (we may have to build tensorflow…** — obs `2026-06-10T00:50:04Z`. CURATED: 2026-06-11.
- [ ] **gpt-2 as an entropy source** — obs `2026-06-10T00:48:37Z`. CURATED: 2026-06-11.
## BACKLOG PROTOCOL

**How to use this file:**
1. Before starting any work, read this file.
2. Pick the highest-priority `[ ]` item in Section 1, then Section 2, etc.
3. Do the work.
4. Post an Apple to IDUNA: `POST /api/v1/apples` (type: completion, title: task name).
5. Check the item `[x]` in this file, add Apple ID and date to the line.
6. `git add BACKLOG.md && git commit -m "backlog: ✓ [task name] — Apple #N"`
7. `git push origin main`
8. Go to step 1.

**How Emily Prime updates this file:**
Emily Prime reads Apples from IDUNA, checks against this file, adds new items when the
system reveals untracked dependencies, and reprioritizes by editing this file and committing.
The backlog is config-as-code. The backlog is what Emily Prime remembers between sessions.

**Golden rule:** No item is complete until the Apple is filed and this file is committed.
The Apple is the proof. The commit is the custody. The push is the delivery.

---

*EMILY PRIME BACKLOG | Cross-repo | Git-authoritative*
*The backlog is what outlasts everything.*
*Clean builds first. Then custody. Then everything else.*
