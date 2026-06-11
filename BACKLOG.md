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

- [x] **TYLER S06E03: "Four Hours"** — Build 0063. S06E03 written. A-269 (East Ferry Ave,
  Detroit). Found construction with anomalous radius reveals new category: AMPLIFICATION —
  Tyler's sustained presence at a found construction (4h40m in 1952) expanded its resonance
  radius from expected ≤20ft to 45ft. First time Tyler says "I'll come back" in Season 6.
  RSI receipts: TYLER-047, Camera Op Entry 40, VC-001 Day 222, Jiangshi Memo #028 (binding
  coefficient under review), EMILY-SPRING-020. S06E04 added to TYLER backlog. DONE 2026-06-08.

- [x] **TYLER S06E04: "The Third Photograph"** — Build 0064. AMPLIFICATION confirmed (named).
  Tyler returns A-269 (Day 224), stays 6h, radius 45ft→71ft. Camera Op finds edge at 70ft
  before Tyler measures 71ft. Third photograph: taken from inside construction threshold (1952);
  Tyler doesn't know if he chose it. A-270 visited: found construction, expected size, no
  amplification — contrast confirms duration mechanism. AMPL-001 ledger opened. RSI receipts:
  TYLER-048, Camera Op Entry 41, VC-001 Day 224, EMILY-SPRING-021. TYLER commit 0544f25.
  DONE 2026-06-08.

- [x] **TYLER S06E06 (Season Finale): "The Usual Place"** — Build 0066. A-273, East Kirby St,
  340ft radius. Tyler visited since before 1929 — "The usual place. Cold." AMPL-001-MAX opened
  (incalculable). Binding coefficient retired. COEFF-002 initiated. AMPL-001-SCOPE floor: 300+.
  Tyler: "I think I know what I am" — did not finish. Camera Op at threshold, already inside
  field. Tyler handed camera back without speaking. Jiangshi reclassification pending. Season 6
  complete. RSI receipts: TYLER-050, Camera Op Entry 43, VC-001 Day 226, Jiangshi #030,
  EMILY-SPRING-023. TYLER commits aa40078, efb2631. DONE 2026-06-08.

- [x] **TYLER S06E05: "How Many"** — Build 0065. AMPL-001 first entries. A-271 (Mt. Elliott):
  prior amplification confirmed from 1952 (2+ hours, "could not easily leave"), current radius
  58ft vs baseline 12-15ft — Tyler revisited, documented, did not amplify further. A-272 (Chene
  St): expected 11ft, Tyler stood at sidewalk and chose not to stay — first deliberate
  non-amplification in documentary history. AMPL-001 scope estimate: 270+ sites globally across
  90 years. Camera Op set camera on table for 43 seconds (first time). Jiangshi: "we are no
  longer documenting what he has done; we are documenting what he is deciding." One archive site
  remains: A-273. RSI receipts: TYLER-049, Camera Op Entry 42, VC-001 Day 225, Jiangshi #029,
  EMILY-SPRING-022. TYLER commits f10670a, 65d35f7. DONE 2026-06-08.

- [x] **TYLER lore files: catch up S06E05–E06** — TYLER-049/050 → eastwind_archive.md.
  Memos #029-030 → jiangshi_project_memos.md. Entries 42-43 → camera_op_sealed_log.md.
  TYLER commit c074b90. DONE 2026-06-09.

- [x] **TYLER S07E01: "The First City"** — Build 0067. Season 7 opener. Exchange Student
  cross-city CARDINAL-3 analysis (8 cities, 200-400 sites). First audit site: A-274, Rue de
  Fleurus, Paris 1931. Tyler knew Gertrude Stein. Stayed 23 days (wrote "two weeks"). Jiangshi
  reclassification memo: no new term, cedes audit methodology. Tyler reads own notation aloud
  for first time. RSI receipts: TYLER-051, Camera Op Entry 44, Jiangshi #031. TYLER commit
  c074b90. DONE 2026-06-09.

- [x] **TYLER S07E02: "Paris"** — Build 0068. First audit measurement. A-274, Rue de Fleurus,
  163ft. Tyler and Camera Op + Exchange Student in Paris. Paris Correspondent on site. Construction
  type revised: DUAL-LAYER — found construction (Gertrude Stein salon 1903-1938) + Tyler
  unintentional amplification (1931). Tyler in room for first time since 1931: 41-second threshold
  pause; told room its measurement. Jiangshi Council Session #6 concluded: reclassification term
  proposed, SEALED, deferred to tomorrow. RSI receipts: TYLER-052, Camera Op Entry 45, Jiangshi
  #032, VC-001 Day 232, EMILY-SPRING-025. DONE 2026-06-09.

- [x] **TYLER S07E03: "The Word"** — Build 0069. Jiangshi reclassification term delivered:
  RESONANT CUSTODIAN. Tyler accepts: "That's accurate." Council Session #6 closed. A-275, Rue
  de l'Odéon, 29ft, 1924 — Tyler wrong about Russian, mechanism proportionate to presence not
  outcome. Day 226 sentence: term part complete; remainder ("why he stayed") still open.
  RSI receipts: TYLER-053, Camera Op Entry 46, Jiangshi #033, VC-001 Day 233, EMILY-SPRING-026.
  DONE 2026-06-09.

- [x] **TYLER S07E04: "The Paris Sites"** — Build 0070. Paris phase complete. 6 sites.
  A-276 (86ft, 1928, faction argument), A-277 (19ft, 1933, light), A-278 FIRST PRE-AWARENESS
  REFUSAL (11ft, 12min, Tyler felt pull and left — had a 1931 rule; still left 11ft), A-279
  (8ft, 1932, inattentive coffee). Mechanism confirmed: proportionate to presence, not attention.
  Refusal-trace: new construction sub-type. RSI receipts: TYLER-054, Camera Op Entry 47,
  Jiangshi #034, VC-001 Day 238, EMILY-SPRING-027. DONE 2026-06-09.

- [x] **TYLER S07E05: "Vienna 1938"** — Build 0071. Vienna phase open. Days 239-255. 4 sites
  (A-280 through A-283). Ichthyosapiens contact F. introduced — left Vienna July 1938, conceded
  witness methodology, gave Tyler records, survived the war. A-282 FIRST COMPREHENDED
  AMPLIFICATION (52ft, 6hrs, rule revised: stay for presence/refuse for absence). A-283
  SECOND-CATEGORY REFUSAL: purpose-based, Tyler said no to amplifying absence. RSI receipts:
  TYLER-055, Camera Op Entry 48, Jiangshi #035. TYLER commit 5ffaf0c. DONE 2026-06-09.

- [x] **TYLER S07E06: "What Stays"** — Build 0072. Vienna days 257-264. A-284: Musikhaus Vogt 31ft,
  violin C major, unambiguous presence. A-285: Café Bernstein 47ft/2h40m, "presence is what still
  knows you," Herr M., AMPL-001 provisional. A-286: Westbahnhof 12ft DEPARTURE-TRACE (new category) —
  F.'s departure, "the going was a presence." A-287: Volksschule 28ft, borderline, understanding-takes-
  time rule revision, ambiguity-refusal category opened. Council Session #7: 17 retroactive second-
  category refusals found, retroactive audit initiated. Fourth category: announced pending. RSI receipts:
  TYLER-056, Camera Op Entry 49, Jiangshi #036, VC-001 Day 264. TYLER commit 5585dc5. DONE 2026-06-09.

- [x] **SHANKPIT Layer 2: scene-isolated snapshot broadcast** — PacketSceneChange = 6 in
  protocol.go + protocol.h. Go server: sendSceneChange() ack on portal travel (14 bytes:
  sceneID + spawn pos). broadcastSnapshots() 20Hz goroutine, scene-filtered snapshot per client
  (18 bytes/entity: clientID+sceneID+pos+yaw). yaw tracked per-client from UserCmds. sync.Mutex
  on clients map. Tests: scene isolation, constant collision check, yaw. SHANKPIT commits
  beb975b, 57fd055. DONE 2026-06-09.

- [x] **SHANKPIT cross-scene attack guard (Go server)** — gameWorld replaces world{} stub in
  apps2/server-go. RayTrace filters clients by sceneID (same-scene only) before ray-sphere
  intersection (hitbox r=0.4, chest-height center). Per-shot shankPlayer uses real client
  pos + sceneID from clientInfo. Vec3.Sub/Len/Dot added to ballistics.go. Tests:
  TestCrossSceneAttackGuard, TestCrossSceneNoHit, TestShooterDoesNotHitSelf. go test passes.
  SHANKPIT commit 0924011. DONE 2026-06-09.

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

- [x] **rsi-loop.sh: task completion detection via Apple** — TOCK phase now tries IDUNA
  Apple polling (Apple ID > pre-TIC snapshot) as primary, falls back to claude-runs/ file
  count when IDUNA unavailable. obs-watcher runReportFooter now instructs Claude to post
  a completion Apple via `emily observe -s success` (best-effort, skipped if emily/IDUNA
  offline). EMILY commit d34269e, PRRJECT_FATBABY commit 9347348. DONE 2026-06-08.

- [x] **emily tui: token spend from IDUNA Apples** — runReportFooter now includes
  tokens_used field (Claude fills it in); tui.go reads actual token counts from
  claude-runs/*.json. When any today's run has tokens_used > 0, displays "(actual)"
  in green with real per-run cost; falls back to 8.2k/run estimate otherwise.
  emily.cli commit 628a6da (v0.8.0), PRRJECT_FATBABY commit 3736982. DONE 2026-06-08.

- [x] **emily tui: keyboard command input bar** — Row 3 CMD bar (tview.InputField). Press ':'
  to focus. Commands: pt/prime-task, eo/observe, tyler [N], start, refresh. Rune hotkeys
  (r, q) bypass while focused; F1-F4 still fire. emily.cli commit fabe5f1 (v0.7.0). DONE 2026-06-08.

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

- [~] **extractProposals() regex — NOT a code bug (investigated 2026-06-08)** — Observation gap
  "0 proposals parsed despite directors found" is a FALSE ALARM. Root cause: (1) the regex
  DOES work correctly — TestParseItem507_BALiveFixture passes with 6 proposals found on actual
  BA cleaned_text; (2) `directors_found` in the observation is the TOTAL accumulated graph node
  count (not this-batch count), so seeing directors > 0 with proposals = 0 is expected when the
  batch contains no proxy vote 8-Ks; (3) event store cursor at 67657 is PAST all proxy vote
  8-Ks (last was at seq ≤ 63759). Recent batches process other 8-K types (earnings, etc.) that
  have no Item 5.07. Observation wording is misleading; no code fix needed. Closed 2026-06-08.

- [x] **governance_health_index score=0 for ALL clean tickers (double-counting)** — Signal IDs
  like `director_long_tenure_{name}_{ticker}` have no date component; every batch re-generates
  identical IDs. DeduplicateSignals(historicalSignals) at load time keeps 1 per ID, but
  combined = historicalSignals + allSignals then has 2 copies. AMZN with 20 long-tenured
  directors: 40×0.03 = 1.20 penalty → score = 0. Fix: added
  `combined = entitygraph.DeduplicateSignals(combined)` before governance health scoring loop
  in cmd/entity-graph/main.go. Test: TestDeduplicateSignals_CombinedHistoricalAndBatch.
  PRRJECT_FATBABY commit 43e9214. DONE 2026-06-08.

- [x] **governance_health_index always score=0** — Root cause: 11 spurious nomination_rejection
  signals from 2011-2013 BA proxy 8-Ks (entities: "Rights Code", "Written Consent", etc.) each
  carried 0.40 penalty, driving BA health to 0.0. Fixed via isSpuriousName(s.Entity) guard in
  ScoreGovernanceHealthWithPenalties(). Test: TestScoreGovernanceHealth_SpuriousEntityIgnored.
  PRRJECT_FATBABY commit 093ec0a. DONE 2026-06-08.

- [x] **Signal accuracy feedback loop: precision always 0** — Root cause: all filing-tied
  signals used DetectedAt=today (batch run date) instead of the SEC filing date. Historical
  2010-2025 proxy votes got detectedAt=2026-06-08, ValidThrough=2027-06-08 → 100% pending,
  precision=0, AccuracyAdjustedPenalties never calibrated. Fixed: signalTimestamps(filingDate)
  helper uses filingDate as detectedAt for historical filings so prediction windows expire in
  the past, immediately enabling confirmed/refuted resolution. Composite/graph-state signals
  (activist_risk, director_link, etc.) still use today. PRRJECT_FATBABY commit a4e0a14.
  DONE 2026-06-08.

- [x] **director_link always 0** — Root cause: name variants for the same person created
  split nodes (e.g. "Susan L. Wagner" at BLK and "Sue Wagner" at AAPL stored as separate
  nodes → no multi-ticker filings → no cross-board links). Fixed via: (1) NamesMatch now
  uses last-name + first-initial matching (ignores middle initials, Jr./Sr. suffixes);
  (2) UpsertPerson fuzzy-scans existing nodes on create; (3) mergeNameVariants() merge
  pass in LoadNodesFromDir retroactively consolidates existing split records.
  Tests: TestNamesMatch_MiddleInitialVariants, TestScoreDirectorLinks_CrossBoardViaNameVariant.
  PRRJECT_FATBABY commit 2bfeb1f. DONE 2026-06-08.

- [x] **`eo` alias for `emily observe`** — Human observation 2026-06-08: "dedicated emily
  observe command eo". `emily eo` now routes to RunObserve (wired in main.go). DONE 2026-06-08.

- [x] **APPLES dedicated git repo auto-sync** — Human observation 2026-06-08: "auto git sync
  all apples with dedicated APPLES repo in real-ish time". Implemented `--apples-git-dir PATH`
  flag on `emily sync` (one-shot and --watch). Each successfully posted Apple is written as
  `<gitDir>/<YYYYMMDD>/<id>_<apple-type>.json` and auto-committed with `git commit -m "apple: #N
  type — title"`. Best-effort (git failures logged, sync continues). emily.cli commit 0974e10.
  DONE 2026-06-08.

---

## SECTION 8: RSI NEXT HORIZON (2026-06-09)

- [x] **Observation batching in obs-watcher** — `--batch-window=Xs` flag added. Directory-scan
  mode: collects all new obs files, gates trivials, invokes Claude once for the batch with a
  consolidated prompt. Separate `.last-batch-processed` cursor. Estimated 50-80% token reduction
  on busy entity-graph days. PRRJECT_FATBABY commit 29cc503. DONE 2026-06-09.

- [x] **RSI loop: smarter TOCK detection via claude-runs/ filename sentinel** — rsi-loop.sh
  TOCK now captures INITIAL_LATEST_RUN (newest filename) at TIC, then detects completion
  when a different filename appears. Eliminates false-positive exits from file count changes.
  EMILY commit 49b6908. DONE 2026-06-09.

- [x] **emily tui: live observation tail in column 2** — 't' hotkey toggles col2 between
  Apple feed and obs tail (last 10 lines of most-recent obs file). renderObsTailPanel() reads
  most-recently modified .json in var/emily-observations/. Also fixed F1 TOCK detection to
  use latestFileInDir (filename sentinel) mirroring rsi-loop.sh fix. emily.cli commit c0075d2
  (v0.9.0). DONE 2026-06-09.

---

## SECTION 9: MJOLNIR — Android Intelligence Terminal

- [x] **IDUNA: push_tokens table + API** — `push_tokens` table (migration 202606090001), SQLite +
  MySQL impls. `POST /api/v1/push-tokens` (upsert token), `GET /api/v1/push-tokens/{name}` (read).
  IAMStore interface updated. Tests pass. IDUNA commit df987e6. DONE 2026-06-09.

- [x] **Emily Prime FCM sender package** — `EMILY/emily-agent/pkg/fcm/sender.go` + `jwt.go`.
  Full FCM HTTP v1 API impl. RS256 service account JWT, OAuth2 token exchange, token caching.
  `IsConfigured()` for graceful degradation. EMILY commit cfe168a. DONE 2026-06-09.

- [x] **Emily Prime push dispatch wiring** — `runPrimeTriageCycle` builds `PushFunc` closure:
  resolves MJOLNIR device token from IDUNA, calls `fcmSender.Send()` for CEO-visible escalations.
  `IdunaClient.GetPushToken()` added. `runPrimeTriage` accepts `PushFunc` callback.
  EMILY commit a0296df. DONE 2026-06-09.

- [x] **MJOLNIR Android project skeleton** — Kotlin + Jetpack Compose + Hilt. IDUNA Retrofit client.
  Google Sign-In → IDUNA JWT flow. FCM token registration. Apple feed screen (LazyColumn).
  Build variant: debug → localhost:8090, release → iduna.einhorn.industrial.
  34 source files across data/remote/repository/notification/di/ui layers. MJOLNIR commit 9bb2e1b. DONE 2026-06-09.

- [x] **APPLES MANIFEST.json generation** — `emily sync --apples-git-dir` now calls
  `updateManifest()` after each `archiveAppleToGit()`. Appends entry, amends the Apple commit.
  MJOLNIR reads MANIFEST.json for fast offline index. emily.cli commit 4775961. DONE 2026-06-09.

- [x] **MJOLNIR docs: Emily Prime integration spec** — `EMILY/docs/MJOLNIR_INTEGRATION.md` authored.
  Covers Apple severity thresholds, device token resolution, morning briefing design, FCM env vars,
  codebase seams, Android registration flow, notification channels. EMILY commit 3341a26. DONE 2026-06-09.

- [x] **Emily Prime morning briefing FCM push** — `emily-agent/briefing.go`. `briefingDue()` gate:
  09:00 UTC ±30 min, sentinel file prevents duplicate. `runMorningBriefing()` fetches last 200 Apples
  from IDUNA, filters to 24h window, groups by type, builds push title + body, fires via PushFunc.
  `IdunaClient.ListApples()` added to iduna.go. Wired into `RunOnce()` PLAN phase in cron.go.
  EMILY commit b2718b1. DONE 2026-06-09.

---

## SECTION 10: EMILY PRIME SELENIUM / WEB AUDIT

- [x] **Emily Prime web audit tool (stdlib HTTP)** — `emily-agent/webaudit.go`. Tool `web_audit_url`:
  fetches URL, checks HTTP status, extracts title/h1, counts and HEAD-checks same-host links,
  returns structured JSON. No external deps (stdlib net/http + regexp). EMILY commit ff5c84f.
  DONE 2026-06-09.

- [x] **Newssite audit preset** — `emily prime-task --preset web-audit-newssite`. Directs Emily Prime
  to audit localhost:8082 + :8083, post findings as signal_observation Apple. Added to RSI
  PRESET_LIST rotation. emily.cli commit 934c567. EMILY commit a0296df. DONE 2026-06-09.

- [ ] **Web audit as front door validator** — Once MJOLNIR WebView targets are live
  (`:8082`, `:8083`), Emily Prime runs web_audit_url on each before Emily's phone gets a
  push linking to them. Guard: if audit fails (5xx, broken links > 3), suppress push and
  post escalation Apple instead. Dependency: web audit tool ✓, MJOLNIR Milestone 2 ✓.

---

## SECTION 11: MJOLNIR INTELLIGENCE + SOURCE BROWSER

- [x] **MJOLNIR camera → Emily Prime intelligence** — Full pipeline:
  CameraX capture (CameraScreen.kt) → base64 → IDUNA `POST /api/v1/intelligence/observe`
  → Emily Prime vision cycle (vision.go, claude-haiku-4-5) → analysis stored + Apple posted.
  IntelligenceScreen.kt (observation history), ObservationDetailScreen.kt (analysis view).
  IDUNA: camera_observations table, intelligence handler (migration 202606090002).
  EMILY: vision.go + IdunaClient.ListPendingObservations/CompleteObservation + runVisionCycle
  in cron.go PLAN phase. IDUNA commit ff8c3fb, EMILY commit 3433f36, MJOLNIR commit 46277f4.
  DONE 2026-06-09.

- [x] **MJOLNIR offline multi-repo source browser** — MultiRepoSyncWorker.kt: JGit clone/pull
  EMILY + TYLER + IDUNA + MJOLNIR + APPLES on WiFi (daily, 24h periodic work). SourceBrowserScreen.kt:
  repo list → depth-4 file tree → dark monospace code viewer (200k cap, h+v scroll).
  Scheduled from MjolnirApp on startup. MJOLNIR commit 46277f4. DONE 2026-06-09.

---

## SECTION 12: HEIMDAL — Sprint Planning Interface

- [x] **HEIMDAL sprint planning interface** — Full pipeline:
  MJOLNIR sends raw requirement text → IDUNA `POST /api/v1/heimdal/sprints` → Emily Prime
  translates via Claude haiku into `RoadmapItem` with `AcceptanceCriterion[]` → adds to RSI
  roadmap → patches sprint status to "queued" → FCM push to MJOLNIR.
  IDUNA: `heimdal_sprints` table (migration 202606090003), HeimdalHandler (submit/list/get/patch),
  `heimdal.submit` + `heimdal.process` permissions, SprintItem auth struct + store methods.
  EMILY: `heimdal.go` (translateRequirement via claude-haiku, IdunaClient HEIMDAL methods,
  runHeimdalCycle on AutonomousCycle), cron.go PLAN phase integration.
  MJOLNIR: SprintItem.kt model, IdunaApi sprint endpoints, HeimdalRepository.kt,
  HeimdalViewModel.kt, HeimdalScreen.kt (requirement input + sprint list with status chips),
  HEIMDAL icon (FlashOn) in ApplesFeedScreen, `heimdal` route in MainActivity.
  DONE 2026-06-09.

- [ ] **HEIMDAL status feedback** — When an RSI task (heimdal-N) completes or is blocked,
  patch the corresponding HEIMDAL sprint to `complete` or `blocked` and send FCM push.
  Requires: emily-agent checks task completion for heimdal-* IDs and calls PatchHeimdalSprint.

---

## SECTION 13: OBSERVATION → BACKLOG CURATION PIPELINE (added Emily Prime 2026-06-11)

- [x] **FatBaby observation → golden backlog curation pipeline (CLI + Emily Prime autonomous)** —
  (1) `emily backlog curate [--all]` CLI: mtime-sorted, trivial filter, INTAKE QUEUE append,
  git commit, Apple. State: EMILY/var/backlog-curated.txt. emily.cli commit ab28096.
  (2) Emily Prime autonomous: runBacklogCuration() in runPrimeTriageCycle — max 5 obs/cycle,
  trivial filter, emilyGitAddCommit. emily-agent commit 8e74bf4.
  (3) emily_read_file / emily_write_file / emily_list_files agent tools: sandboxed, auto-commit.
  DONE 2026-06-11.

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

---

- [ ] ** go run . apples list
error: iduna auth: Post "http://localhost:8080/api/v1/auth/agent": dial tcp [::1]:8080: connect: …** — obs `2026-06-10T00:51:05Z`. CURATED: 2026-06-11.
- [ ] **fatbaby@localhost:~/emily.cli$ systemctl --user status daemon-reload
Failed to connect to bus: Permission denied
fatbab…** — obs `2026-06-07T21:23:23Z`. CURATED: 2026-06-11.
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
