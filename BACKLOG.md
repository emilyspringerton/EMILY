# EMILY PRIME — CROSS-REPO GOLDEN BACKLOG
## Owner: Emily Prime | Machine-readable | Git-authoritative
### Last updated: 2026-06-09 | MJOLNIR northstar + Selenium web audit added

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

---

## SECTION 4: SHANKPIT / TYLER GAME ENGINE (lower priority)

- [ ] **SHANKPIT → MPT bridge** — Spec exists (engine/shankpit_mpt_bridge.md). Implementation
  deferred until MPT is running end-to-end.

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

## SECTION 13: OBSERVATION → BACKLOG CURATION PIPELINE (added Emily Prime 2026-06-11)

- [ ] **Backlog intake: EmilyOS mobility edition — bare-metal exokernel, ISO2424242** —
  FatBaby observation `2026-06-11T00:54:39Z`: "emily os mobility edition bare metal exokernel
  ISO2424242". New product initiative: EmilyOS as a bare-metal exokernel OS targeting
  mobile/embedded hardware. ISO identifier: ISO2424242. Northstar to be drafted; new repo
  to be created (candidate name: EmilyOS/). Related pending observations: EmilyOS Arch vs Debian
  + BAZEL-equivalent repo (2026-06-11T00:00:18Z); PITVIPER standalone SDL2 terminal
  (2026-06-11T00:52:56Z). Status: INTAKE PENDING curation pipeline above.
  ADDED: 2026-06-11 (first item produced by curation session focus).

---

## INTAKE QUEUE (curated by emily backlog curate)

Items below have been auto-curated from FatBaby observations. Emily Prime reviews
and promotes them into the appropriate section when she plans the next sprint.

- [ ] **find all the gaps in the instrumentation where we should be publishing apples but are not and create new observations f…** — obs `2026-06-11T01:24:22Z`. CURATED: 2026-06-11.
- [ ] **emily os mobility edition bare metal exokernel ISO2424242** — obs `2026-06-11T00:54:39Z`. CURATED: 2026-06-11.
- [ ] **the fatbaby custom terminal v0 northstar the built the emily way the repo is named PITVIPER** — obs `2026-06-11T00:52:56Z`. CURATED: 2026-06-11.
- [ ] **the news site should have stock charts - do it the emily way - lets build our own charting library so we can easily add…** — obs `2026-06-11T00:51:58Z`. CURATED: 2026-06-11.
- [ ] **all system inputs need to feed into a single log stream that is continuously synced to git** — obs `2026-06-11T00:02:32Z`. CURATED: 2026-06-11.
- [ ] **HACK feed a certain amount of apples as context to save tokens for claude having to read certain files he will glean co…** — obs `2026-06-11T00:01:38Z`. CURATED: 2026-06-11.
- [ ] **EmilyOS Debian first or Arch? we want to have our own BAZEL or BAZEL equivalent soc2 software repository** — obs `2026-06-11T00:00:18Z`. CURATED: 2026-06-11.
- [ ] **implement basic haiku -> FABLE advisor tool for emily prime** — obs `2026-06-10T23:59:02Z`. CURATED: 2026-06-11.
- [ ] **implement emily prime haiku/sonnet -> advisor FABLE for planning sprints indo HEIMDAL sprint planning API - we need to …** — obs `2026-06-10T23:56:50Z`. CURATED: 2026-06-11.
- [ ] **we need to have an emily prime api so we can use her agent to help direct the rsi loops so we can split the token usage…** — obs `2026-06-10T23:54:59Z`. CURATED: 2026-06-11.

---

- [ ] **emily prime should have all of the api capabilities of the emily.cli command** — obs `2026-06-10T23:38:51Z`. CURATED: 2026-06-11.
- [ ] **how are these prime-task s getting fired off? is that claude via the observation watcher? or is that claude larping as …** — obs `2026-06-10T23:36:46Z`. CURATED: 2026-06-11.
- [ ] **the apples need to all be synced into the git repo APPLES (IDUNA needs to do it because she is the APPLES source of tru…** — obs `2026-06-10T23:35:13Z`. CURATED: 2026-06-11.
- [ ] **i think we may need to review all of the golden documentations in all the repos and update the emily prime system propt…** — obs `2026-06-10T23:33:35Z`. CURATED: 2026-06-11.
- [ ] **IDUNA LOGIN PAGE and ADMIN DASHBOARD - is it built? if not build it - if it is built add to the readme how to get to th…** — obs `2026-06-10T22:01:42Z`. CURATED: 2026-06-11.
- [ ] **RSI cycle 6: conditional rules inlining + dedup process list** — obs `2026-06-10T20:34:42Z`. CURATED: 2026-06-11.
- [ ] **what is emily sync? is it pulling observations from fatbaby and pushing them to emily prime to curate the golden backlo…** — obs `2026-06-11T01:26:02Z`. CURATED: 2026-06-11.
- [ ] ** go run . apples list
error: iduna auth: Post "http://localhost:8080/api/v1/auth/agent": dial tcp [::1]:8080: connect: …** — obs `2026-06-10T00:51:05Z`. CURATED: 2026-06-11.
- [ ] **gpt-2 c fork git@github.com:emilyspringerton/gpt2-alpine-c.git parity with og gpt-2 repo (we may heve to build tensorfl…** — obs `2026-06-10T00:50:04Z`. CURATED: 2026-06-11.
- [ ] **gpt-2 as an entropy source** — obs `2026-06-10T00:48:37Z`. CURATED: 2026-06-11.
- [ ] **i think there is a tui bug with the process handling it no longer works on a clean reboot it used to** — obs `2026-06-10T00:25:48Z`. CURATED: 2026-06-11.
- [ ] **add emily.cli to emily.cli status including the pid of the tui** — obs `2026-06-07T23:00:33Z`. CURATED: 2026-06-11.
- [ ] **fosho the bug with the emily.cli tui hang still exists even if we hit q to exit gracefully (unsure if it always happens…** — obs `2026-06-07T22:53:25Z`. CURATED: 2026-06-11.
- [ ] **if fatbaby observation watcher is started during a claude rate limit a prime task may be skipped** — obs `2026-06-07T22:50:41Z`. CURATED: 2026-06-11.
- [ ] **filing observations are intentionally uneditable but what happens when a human screws up and types the wrong thing? do …** — obs `2026-06-07T22:05:53Z`. CURATED: 2026-06-11.
- [ ] **TYLER expand the role of the phoen** — obs `2026-06-07T22:04:47Z`. CURATED: 2026-06-11.
- [ ] **add TYLER specific cutscene functionality - audit the capabilities of SHANKPIT for scripting cutscenes like in ffxi (di…** — obs `2026-06-07T21:55:11Z`. CURATED: 2026-06-11.
- [ ] **news site 500s - iduna related? too big of an index to query? dependent processes? (i started the news site manualy)** — obs `2026-06-07T21:43:40Z`. CURATED: 2026-06-11.
- [ ] **emily cli status should show all of the fatbaby processes if we pass a fatbaby flag** — obs `2026-06-07T21:35:59Z`. CURATED: 2026-06-11.
- [ ] **the tui should show the real time if possible** — obs `2026-06-07T21:27:44Z`. CURATED: 2026-06-11.
- [ ] **TUI should show fatbaby data via menu controls (fatbaby mode)** — obs `2026-06-07T21:27:09Z`. CURATED: 2026-06-11.
- [ ] **TUI hang ensure fixed (gracfully exit the session when we kill it with ctrl c** — obs `2026-06-07T21:26:18Z`. CURATED: 2026-06-11.
- [ ] **fatbaby@localhost:~/emily.cli$ systemctl --user status daemon-reload
Failed to connect to bus: Permission denied
fatbab…** — obs `2026-06-07T21:23:23Z`. CURATED: 2026-06-11.
- [ ] **fatbaby news site does not start via emily cli start (all)** — obs `2026-06-07T21:18:53Z`. CURATED: 2026-06-11.
- [ ] **observation-watcher must inject full reporting and git sync requirements into the Claude Code prompt** — obs `2026-05-31T20:56:45Z`. CURATED: 2026-06-11.
- [ ] **All required and optional environment variables must be documented at the top of the README** — obs `2026-05-31T20:43:15Z`. CURATED: 2026-06-11.
- [ ] **Feature: create a GitHub issue automatically whenever Emily writes an observation** — obs `2026-05-31T20:41:56Z`. CURATED: 2026-06-11.
- [ ] **UX: ticker search should auto-navigate on click and on Enter key — remove redundant Go button step** — obs `2026-05-31T20:39:51Z`. CURATED: 2026-06-11.
- [ ] **Bug: newssite articles use publication date not filing date — historical filings appear as breaking news** — obs `2026-05-30T22:14:19Z`. CURATED: 2026-06-11.
- [ ] **Feature request: self-improving model training pipeline — user data flywheel for Emily fine-tuning** — obs `2026-05-30T22:10:22Z`. CURATED: 2026-06-11.
- [ ] **Feature request: full GTM funnel — Ask Emily free tier, Emily+ subscription, community editorial engine, Merkle query…** — obs `2026-05-30T22:02:43Z`. CURATED: 2026-06-11.
- [ ] **Feature request: newssite ingest endpoint for Emily-authored governance commentary articles** — obs `2026-05-30T21:53:49Z`. CURATED: 2026-06-11.
- [ ] **EDGAR submissions endpoint returning truncated JSON for all 5 major bank tickers — BAC, C, GS, JPM, MS** — obs `2026-05-30T21:45:32Z`. CURATED: 2026-06-11.
- [ ] **entity-graph cannot detect 8-K documents from persisted store — form/source_type field mismatch, 846 docs unprocessab…** — obs `2026-05-30T21:32:12Z`. CURATED: 2026-06-11.
- [ ] **entity-graph parsing all 8-K subtypes — Item 5.07 not found in non-proxy filings, producing 100% parse failure rate** — obs `2026-05-30T21:29:24Z`. CURATED: 2026-06-11.
- [ ] **entity-graph reads 0 filings despite 846 source documents in var/secwatch** — obs `2026-05-30T09:52:12Z`. CURATED: 2026-06-11.
- [ ] **eps-processor ticker map has only 2 entries — all press releases are being dropped silently** — obs `2026-05-30T09:46:52Z`. CURATED: 2026-06-11.
- [ ] **eps-processor ticker map has only 2 entries — all press releases are being dropped silently** — obs `2026-05-30T09:42:55Z`. CURATED: 2026-06-11.
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
