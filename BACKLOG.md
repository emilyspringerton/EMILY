# EMILY PRIME — CROSS-REPO GOLDEN BACKLOG
## Owner: Emily Prime | Machine-readable | Git-authoritative
### Last updated: 2026-06-07 | RSI master loop + Bloomberg TUI — Claude Code build

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

- [x] **IDUNA embedded SQLite backend** — Zero-MySQL startup. Bootstrap auto-detects.
  All agents can authenticate without human MySQL provisioning. Apple #1 (Tyler activation).
  Ref: IDUNA commit 8490bb5. DONE 2026-06-07.

- [x] **Tyler agent registration (Build 0016)** — .claude/tyler_agent.md, settings.json,
  emily.sh wired. Tyler authenticated on live IDUNA. Apple #1 filed. Blame path: OCCUPIED.
  DONE 2026-06-07.

- [x] **EMILY/BACKLOG.md created** — This file. Cross-repo golden backlog. Machine-readable.
  Apple #3 (Emily Prime decision) authorized. DONE 2026-06-07.

- [x] **Run emily.sh — first autonomous RSI iteration** — Tyler ran emily.sh. Build 0017 (MPT
  package) + Build 0018 (1952 Detroit annual composite + architecture nodes). 59 BACKLOG items
  complete, 3 remaining (MPT-blocked). Apple #16 (rsi_iteration) filed. DONE 2026-06-07.

- [x] **Wire EMILY RSI engine to IDUNA** — `EMILY/scripts/start-emily-agent.sh` sources
  IDUNA credentials and launches emily-agent with `IDUNA_BASE_URL`, `IDUNA_AGENT_NAME`,
  `IDUNA_AGENT_SECRET` set. cron.go already has Apple submission code (lines 244-257) — it
  fires when IdunaClient is non-nil. Restart running agent using start-emily-agent.sh to enable.
  DONE 2026-06-07.

- [x] **Process supervision: IDUNA systemd unit** — Unit file written to
  `IDUNA/scripts/iduna.service`. Embedded SQLite, no MYSQL_DSN. WorkingDirectory=/home/fatbaby/IDUNA,
  EnvironmentFile=~/.config/iduna/env (optional). Restart=on-failure, RestartSec=10.
  Deploy: `cp iduna.service ~/.config/systemd/user/ && systemctl --user enable --now iduna.service`.
  Binary: `go build -o ~/.local/bin/iduna .` in IDUNA repo. Apple #50. DONE 2026-06-07.
  (Human must build binary and run systemctl commands to activate.)

- [x] **Process supervision: emily.sh cron or supervisor** — `TYLER/scripts/cron-emily.sh`
  written. Lock-guarded, IDUNA-check before run, sources Tyler credentials, git pull + 1
  iteration + push. Crontab: `0 */4 * * * /home/fatbaby/TYLER/scripts/cron-emily.sh`.
  Apple #18 filed. DONE 2026-06-07. (Human must add crontab entry.)

---

## SECTION 2: MONEYPRINTERTURBO (video pipeline)

- [x] **MPT config.toml written** — LLM: litellm → anthropic/claude-sonnet-4-6.
  Pexels key placeholder. Build 0017. DONE 2026-06-05.

- [ ] **MPT Pexels API key** — Set YOUR_PEXELS_API_KEY_HERE in MoneyPrinterTurbo/config.toml.
  Human action required (key must come from the user). BLOCKED: waiting on key.

- [ ] **S01E01 cold open compiled clip** — Once MPT is running, run the cold open compilation.
  See TYLER BACKLOG.md for full spec. Dependency: Pexels key ✓, MPT service running.

- [ ] **MPT → TYLER RSI trigger** — When emily.sh produces a new episode script, auto-invoke
  MPT. Spec exists in engine/moneyprinter_pipeline.md. Dependency: cold open clip ✓.

---

## SECTION 3: SYSTEM OBSERVABILITY (Emily Prime sees everything)

- [x] **emily CLI (emily.cli repo)** — Full Go CLI v0.5.0. 69 tests, 5 packages.
  Commands: observe (stdin+auto-Apple), apples (list/post/get), watch (IDUNA tail -f),
  status (--watch live dashboard), sync (--watch daemon), install (--cron/--systemd),
  prime-task (directed tasks→EMILY/signals/tasks), agents (activity dashboard), help.
  Apples #37–#48. README.md complete. DONE 2026-06-07.

- [x] **Apples dashboard query** — `EMILY/scripts/apples.sh` written. Authenticates as
  EMILY-PRIME via M2M, queries `/api/v1/apples`, prints tabular view. Flags: `--full`,
  `--limit=N`, filter by source_repo. Apple #4 filed. DONE 2026-06-07.

- [x] **Cross-repo status script** — `EMILY/scripts/status.sh` written. Git status for all
  5 repos (branch, dirty count, last commit, backlog stats) + last Apple per agent from IDUNA.
  DONE 2026-06-07.

- [x] **PRRJECT_FATBABY → IDUNA Apple filing** — `EMILY/scripts/sync-fatbaby.sh` syncs
  observation files to IDUNA as `signal_observation` Apples. State-tracked (EMILY/var/
  fatbaby-synced.txt) to avoid duplicates. 9 historical observations backfilled, Apples #6–14.
  Run on cron or after each FatBaby cycle. DONE 2026-06-07.

---

## SECTION 4: SHANKPIT / TYLER GAME ENGINE (lower priority)

- [x] **Tyler mode SHANKPIT spec** — `TYLER/engine/tyler_hum_mechanic.md` written. First
  Tyler↔SHANKPIT implementation bridge: HUM_FIELD_ACTIVE server event (6-unit radius, 45s,
  Jiangshi NPC hold + documentation -80%). Go interface sketched. RSI receipts documented.
  Apple #18 filed. DONE 2026-06-07.

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

- [x] **PRRJECT_FATBABY Emily Prime loop (CLI layer)** — `emily prime-task` command writes
  directed task JSON to EMILY/signals/tasks/. Observation-watcher picks up within 10s and
  invokes Claude Code on FatBaby. Loop: operator→CLI→tasks/→watcher→claude. Apple #44.
  DONE 2026-06-07. (Full autonomous loop requires obs-watcher running with --prime-tasks flag.)

- [x] **Obs-watcher with --prime-tasks** — `emily start` launches observation-watcher with
  `--prime-tasks EMILY/signals/tasks` automatically. Already wired in cmd/start.go.
  Auto-detects EMILY sibling dir if not specified. DONE 2026-06-07.

- [x] **RSI master loop (rsi-loop.sh)** — `EMILY/scripts/rsi-loop.sh` orchestrates the
  full tic-toc RSI cycle: TIC (emily prime-task --preset rsi-token-report) → TOCK (wait for
  Claude Code completion via claude-runs/ polling) → ENTROPY (TYLER/emily.sh N iterations) →
  ANALYZE (emily observe posts observation → triggers next cycle). Configurable via env vars:
  MAX_ITERS, TYLER_ITERS, SLEEP_BETWEEN, DRY_RUN, SKIP_TYLER. Writes per-iteration state to
  EMILY/var/rsi-loop-state.json. DONE 2026-06-07.

- [x] **emily tui — Bloomberg terminal** — `emily tui` command (v0.6.0). Three-column live
  dashboard: repos+tasks+token budget | Apple feed (auto-refresh) | process health+RSI loop+
  actions. Hotkeys: F1=fire RSI task, F2=run Tyler, F3=start system, F4=tail logs, r=refresh,
  q=quit. Uses tview (first external dep). 15s auto-refresh. DONE 2026-06-07.

---

## SECTION 6: RSI TIGHTENING (next horizon)

- [ ] **rsi-loop.sh: task completion detection via Apple** — Instead of polling claude-runs/
  file count, poll IDUNA for a new Apple with run_id matching the task_id. This makes
  completion detection precise regardless of claude-runs/ layout.
  Depends: IDUNA running + emily CLI auth.

- [ ] **emily tui: token spend from IDUNA Apples** — Parse token_used field from
  rsi_iteration Apple bodies for actual spend (vs current rough estimate of 8.2k/run).
  Adds real cost visibility to the Bloomberg terminal.

- [ ] **emily tui: keyboard command input bar** — Add a bottom input field (tview.InputField)
  so operators can type `emily prime-task "..."` directly from the TUI without leaving.

- [x] **rsi-loop.sh: FatBaby + EMILY combined tick** — After TOCK phase, POSTs to
  `$FATBABY_AGENT_URL/tick` (default http://localhost:8080/tick). HTTP 200 = sweep accepted;
  non-200 is a warning (agent offline) and does not abort the loop. Skip with SKIP_FATBABY_TICK=1.
  DONE 2026-06-08.

- [x] **rsi-loop.sh preset rotation** — Cycles through PRESET_LIST env var (default:
  rsi-token-report entity-graph-refinement eps-coverage-review). Each iteration uses
  next preset in the list (modulo). Prevents loop from only optimizing one surface.
  DONE 2026-06-08.

---

## SECTION 7: SIGNAL QUALITY (from 2026-06-08 FatBaby observations)

- [x] **director_long_tenure spurious entities from BA filings** — Proposal-topic nodes
  ("Rights Code", "Special Meetings", "Written Consent", "Political Activity", etc.) were
  persisted as NodeDirector in var/entity-graph/nodes.ndjson from 2011–2013 Boeing 8-Ks
  where reProposalSplitter missed the boundary. Fixed via: (1) 13 new nonNameWords in parser.go
  blocking future ingestion; (2) isSpuriousName() guard in ScoreLongTenure() skipping existing
  graph store nodes. Both in PRRJECT_FATBABY. go test ./... passes. DONE 2026-06-08.

- [ ] **extractProposals() regex — BA-style filing boundary miss** — Observation gap 2026-06-08:
  "0 proposals parsed despite directors found — proposal-splitter regex likely did not match
  filing text format". reProposalSplitter misses some BA 8-K proposal boundaries, causing
  dirBody to span the full body. This explains 0 proposals in most runs today (only 2 of 15+
  runs parsed any proposals). Fix: add reProposalSplitter alternation for BA-format boundaries
  (e.g., "Shareholder Proposal Relating to..." as a prose splitter), or extend reProseSplitter.
  Requires actual failing filing text as test fixture.
  Ref: PRRJECT_FATBABY/internal/entitygraph/parser.go extractProposals().

- [ ] **governance_health_index always score=0** — Every ticker's health index is scoring 0
  (critical) in every run on 2026-06-08. Either the scoring calculation is broken or the
  calibration thresholds produce 0 for all current data. Investigate
  internal/entitygraph/signals.go health score calculation.

- [ ] **Signal accuracy feedback loop: precision always 0** — All accuracy_scores have
  precision=0 with 100% pending predictions across all signal types. The confirmation/refutation
  mechanism exists (AccuracyRecord) but nothing is closing the loop. Need a mechanism to
  confirm or refute predictions from subsequent filings or price data.

- [ ] **director_link always 0** — No director_link signals in any of today's 15+ runs.
  Cross-board director linkage should fire when directors appear at multiple tickers. Either
  the edge-building logic isn't finding cross-ticker matches, or the scoring threshold is too
  high. Investigate BuildEdgesFromFiling + ScoreDirectorLinks.

- [x] **`eo` alias for `emily observe`** — Human observation 2026-06-08: "dedicated emily
  observe command eo". `emily eo` now routes to RunObserve (wired in main.go). DONE 2026-06-08.

- [ ] **APPLES dedicated git repo auto-sync** — Human observation 2026-06-08: "auto git sync
  all apples with dedicated APPLES repo in real-ish time". Spec: new `APPLES` repo (or
  subdirectory) that `emily sync --watch` commits to as a git-native Apples archive. Enables
  offline audit, branching on Apple history, and cross-machine Apple portability without IDUNA.

---

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
