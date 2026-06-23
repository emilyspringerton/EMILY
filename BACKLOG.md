# EMILY PRIME — CROSS-REPO GOLDEN BACKLOG
## Owner: Emily Prime | Machine-readable | Git-authoritative
### Last updated: 2026-06-12 | S23 (EDIS WordPress product) + S24 (Newssite ops hardening) added

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

- [x] **MPT → TYLER RSI trigger** — compiled/mpt_episode_trigger.sh written per moneyprinter_pipeline.md §VI. Extracts MPT_TOPIC, generates payload, POSTs to MPT API, polls, routes output, writes Emily Prime observation. --dry-run verified on S01E01. Apple #3090 | 2026-06-23

---

## SECTION 3: SYSTEM OBSERVABILITY (Emily Prime sees everything)

- [x] **Apple instrumentation audit** — Enumerate all system events that should post an Apple
  but don't: cron triggers, observation drops, haiku call failures, HEIMDAL state changes,
  FCM failures. Create `emily eo` observations for each gap found.
  — obs `2026-06-11T01:24:22Z`. Done 2026-06-11. Apple #344. Commit EMILY 8ea558a.
  (Note: obs-watcher now posts warning Apple on permanent invoke failure — partial gap closed.)

- [x] **Single log stream** — All system inputs (obs-watcher, rsi-loop, emily-agent, IDUNA,
  PRRJECT_FATBABY) feed into a single append-only log synced to git on every write.
  Candidate: `var/emily-stream.ndjson` → `emily sync --stream`.
  — obs `2026-06-11T00:02:32Z`. Done 2026-06-11. Apple #345. Commits EMILY 645ece9 + emily.cli 6fd1a28.

- [x] **Apples IDUNA→APPLES git sync** — IdunaClient reads APPLES_GIT_DIR env var; after
  every successful Apple POST fires async `emily sync --apples-git-dir <dir>` goroutine.
  — obs `2026-06-10T23:35:13Z`. Done 2026-06-11. Commit emily-agent (iduna.go).

- [x] **Golden context feed via Apples** — `IdunaClient.FetchAppleContext(ctx, n)` returns
  compact Apple summaries (≤200 tokens). Injected into HEIMDAL translateRequirement haiku
  calls via runHeimdalCycle. — obs `2026-06-11T00:01:38Z`. Done 2026-06-11. Commit emily-agent.

- [x] **Web-audit: newssite and signalapi stale logs (8082/8083)** — Diagnose and restore newssite (port 8082) and signalapi (port 8083) services; logs last updated 2026-06-07. Obs: 2026-06-13T04:04:17Z.
  — Done 2026-06-16. Apple #824. entitygraph.mergeNameVariants nil-deref when ids[i] deleted mid-loop; fix: break inner j on dropped==ids[i]. newssite :8082 serving 200 OK. Commit PRRJECT_FATBABY 8539693.
---

## SECTION 4: SHANKPIT / TYLER GAME ENGINE (lower priority)

- [x] **SHANKPIT story horror mechanics — TYLER S10 lore** — Mechanism Reader HUD [Apple #1071, 2026-06-17]
  (Goetia Hz + site taxonomy), boss question-state phase (33% trigger, 10% damage, zero-taxonomy
  reading, observation point resolution), S10E07/E08 cutscene slides (al-Waqfa, al-idrak
  al-muttasil, MARRAKECH-001 three-signal, Camera Op Entry 88). Build 0110. Commit 2f5387b.

- [x] **SHANKPIT story: second TYLER lore level** — CAVE-001 (al-Waqfa) endurance mission:
  mechanism always ZERO_TAXONOMY (site IS the source), indestructible entity, 90s endurance timer,
  8+6 cave cutscene slides (S10E02–E03 lore), draw_story_cave_entity slow-field visual.
  "The archive cannot read its own source." — Build 0111. Apple #1074. Commits fa0e4e8 + 5579d5a. 2026-06-17.

- [~] **SHANKPIT → MPT bridge** — Spec exists (engine/shankpit_mpt_bridge.md). Implementation
  deferred until MPT is running end-to-end.

- [x] **SHANKPIT Dragonfly server (Milestones 1+2) — end-to-end verified** — Go backend :6969
  builds + runs. UDP handshake (PacketConnect→PacketWelcome), portal travel (BTN_USE → PortalTriggered
  → ResolvePortalDestination → sendSceneChange), 20Hz per-scene snapshots, voxel streaming, and
  cross-scene attack guard all operational. Architecture memo filed to Emily Prime for Milestone 3
  (WorldBackend interface) and Milestone 4 (Construct expansion). — 2026-06-12.

- [x] **SHANKPIT Milestone 3: WorldBackend Go interface** — Define the portal_resolve_destination / [Apple #358, 2026-06-12]
  world_backend_*() abstraction as a Go interface in server/. Spec must go in docs2/specs/ FIRST.
  Enables: swapping Dragonfly fork in behind stable seam. Dependency: Emily Prime architecture review.
  — obs `2026-06-12`.

- [x] **SHANKPIT Milestone 4: Construct expansion** — DragonsNShit bridge code, scene/portal code, [Apple #369, 2026-06-12]
  world backend surfaces must be included in the Construct. Currently the Construct only covers core
  FPS files — DragonsNShit is invisible to agents and automation. — obs `2026-06-12`.

- [x] **SHANKPIT: GoblinFoxDragon repo relationship** — Clarify whether GoblinFoxDragon is the [Apple #361, 2026-06-12]
  Dragonfly fork (DragonsNShit backend), a separate umbrella, or the future consolidated repo.
  Both have go.mod module=dragonsnshit. Document in NORTHSTAR.md. — obs `2026-06-12`.

- [x] **SHANKPIT: Season lineage snapshot schema** — No spec exists. Write minimal schema: what [Apple #364, 2026-06-12]
  fields are captured at season-end for lineage query? Shapes all persistence decisions.
  — obs `2026-06-12`.

- [x] **TYLER: expand the role of the phone** — The in-game smartphone should have richer [Apple #371, 2026-06-12]
  mechanics: notifications, apps, contacts, map. Spec to be written.
  — obs `2026-06-07T22:04:47Z`.

- [x] **TYLER: cutscene system (FFXI-style dialogue)** — Add TYLER-specific cutscene [Apple #371, 2026-06-12]
  functionality. Audit SHANKPIT for scripting capabilities matching FFXI dialogue scenes.
  — obs `2026-06-07T21:55:11Z`.

- [x] **TYLER S09E01 'The Interval' — Córdoba pre-archive site** — Implement Season 9 Episode 1 with new pre-archive site disclosure at Córdoba (~10th century). Obs: 2026-06-14T17:37:05Z.
  — Stale. Already done: s09e01_the_interval.md 533 lines, Apple #514, Season 9 complete. Apple #835.
- [x] **SHANKPIT display/fullscreen hard blocker for Steam** — Implement fullscreen display mode for SHANKPIT as critical pre-Steam release requirement per DISPLAY_FULLSCREEN_SPEC.md. Obs: 2026-06-13T20:56:16Z.
  — Stale. Already done: VIRTUAL_W/H 1280×720, letterbox, toggle_fullscreen, SDL_WINDOW_RESIZABLE, shankpit_display.cfg. Commit 2026-06-14. Apple #836.
---

## SECTION 5: FUTURE (Emily Prime decides when to promote)

- [ ] **MySQL wire protocol embedded server** — go-mysql-server as in-process backend.
  Enables: external tooling access, migration path to production MySQL, MySQL-compatible
  backups. Ref: TYLER/outlines/emily_iduna_bootstrap.md Option B.
  Dependency: embedded SQLite stable ✓ (already done), need concrete reason to upgrade.

- [ ] **Tyler IDUNA agent registration via iduna CLI** — `iduna agents register --id tyler`
  when the IDUNA CLI is built. Currently Tyler is seeded via migration. The CLI path is
  the programmatic future.

- [x] **S29-05 RSI smoke test: end-to-end loop verification** — Pipeline confirmed end-to-end: RSI cycles fire, Apples file to IDUNA, obs-watcher dispatches, context-overflow recovery works. 3 bugs found and fixed (cursor format, isContextTooLongOutput stdout capture, go run . compile). Blocked at final claude dispatch by API credit balance (user action: top up console.anthropic.com). Apple #848, commit 7edf6f5.
- [x] **FatBaby system health check: 4 fixes applied** — signalapi O(N) scan (86% CPU → 0%), form4-watcher XSL prefix (0→479 transactions), form4-watcher 4MB→32MB body limit, SQLite COMMENT= migration. All 14 processes healthy. Apple #1114 | 2026-06-17.
- [x] **signal pipeline audit: 10748 signal_failed today (9435=EDGAR 429 no-retry, 977=…** — Addressed by S36-01 (429 retry+throttle), S36-02 (skip pre-2000 empty URL filings), S36-03 (4MB→16MB limit), S36-05 (empty ticker). All 4 root causes fixed 2026-06-17. Apple #1229–#1236. Obs: 2026-06-17T22:21:23Z. — CLOSED — 2026-06-18
---

## SECTION 6: RSI TIGHTENING (next horizon)

- [x] **TUI: show real-time clock** — Already implemented via 1s clockTicker in tui.go.
  Verified 2026-06-11. Apple #330.

- [x] **emily cli status: show all fatbaby processes** — `emily status --fatbaby` lists
  all PRRJECT_FATBABY processes (newssite, signalapi, eps-processor, eps-reconciler,
  entity-graph, observation-watcher, secwatch) + IDUNA health. Also works in --watch mode.
  — obs `2026-06-07T21:35:59Z`. Apple #330. Committed emily.cli 8e0d1be.

- [x] **emily cli status: include emily.cli TUI PID** — `emily status` now reads
  `/tmp/emily-tui.pid` and checks liveness via kill -0. — obs `2026-06-07T23:00:33Z`.
  Apple #330. Committed emily.cli 8e0d1be.

- [x] **emily observe: typo correction mechanism** — `emily obs amend <key> <correction>`
  appends a correction note (original summary preserved); Amendment+AmendedAt fields added
  to obs.Payload JSON schema. — obs `2026-06-07T22:05:53Z`. Apple #330. Committed emily.cli 8e0d1be.

- [x] **TUI: fatbaby mode** — `emily tui --fatbaby` activates FatBaby panel in col 3.
  'b' key toggles between system health and FatBaby (processes + entity graph stats +
  eps counts). v0.8.0. — obs `2026-06-07T21:27:09Z`. Done 2026-06-11. Commit emily.cli.

- [x] **obs-watcher: rate-limit resilience** — invokeWithRetry wraps all three invocation
  sites. Captures stderr, detects rate-limit signals, retries 3x (30s→90s→270s). Permanent
  failure posts warning Apple via emily observe. — obs `2026-06-07T22:50:41Z`. Done 2026-06-11.
  Commit PRRJECT_FATBABY c04fd79.

---

## SECTION 7: SIGNAL QUALITY (from 2026-06-08 FatBaby observations)

- [x] **extractProposals() false-alarm gap fixed (2026-06-13)** — Root cause was `detectGaps()`
  using `len(graph.Nodes)` (total accumulated nodes) for the gap condition instead of directors
  found in the CURRENT batch. Fixed: `BuildObservation` now accepts `directorsThisBatch` (this-
  batch director vote count); gap fires only when the current batch had proxy filings with
  director votes but no proposals — eliminating all false positives for non-proxy 8-K batches.
  — obs `2026-06-04T10:00:00Z`. Done 2026-06-13. Apple #426. Commit PRRJECT_FATBABY 371d957.

---

## SECTION 8: RSI NEXT HORIZON (2026-06-09)

---

## SECTION 9: MJOLNIR — Android Intelligence Terminal

- [x] **MJOLNIR Milestones 0–3 audit + NORTHSTAR update** — Milestones 0–3 COMPLETE.
  NORTHSTAR updated [x]. README → Milestone 4. TYLER/EPISODES.md created (52 eps, S1–S7, Build 0082).
  TYLER + SHANKPIT product cards added to ProductsScreen. Apple #395. Done 2026-06-12.

- [x] **MJOLNIR Milestone 4: RSI loop state display** — RsiScreen + RsiViewModel. EmilyApi calls
  `GET /api/v1/emily/state` (Emily Prime port 8086). Shows: cycle number, metrics (tasks done,
  iters, failures), active task + progress, next cycle plan. Timeline icon in feed. EMILY_BASE_URL
  BuildConfig field added (debug: 10.0.2.2:8086). Apple #397. Commit MJOLNIR 9d88ec6.

- [x] **MJOLNIR Milestone 4: Token spend sparkline** — IDUNA GET /api/v1/apples/stats/daily-tokens added (DailyTokenStat, SQLite+MySQL store, zero-padded 7-day series, apples.read gate). MJOLNIR: TokenSparklineCard Canvas bar chart in RsiScreen; RsiViewModel fetches stats in parallel with cycle state. M4 now fully complete. IDUNA 19bacbc, MJOLNIR 2695d1d. Apple #487.

- [x] **MJOLNIR Milestone 4: RSI cycle push** — FCM push fires on task completion in Apple-filing
  goroutine in cron.go. Title: "RSI cycle N complete". Body: task description. Priority: normal.
  Data: apple_id, cycle, tasks_done. Fires only on task.Status == "success" (not idle cycles).
  Done 2026-06-12. Apple #396. Commit EMILY 50dc5a1.

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

- [x] **HEIMDAL status feedback** — notifyHeimdalStatus goroutine fires on heimdal-* task
  terminal status (complete or blocked): patches sprint in IDUNA, files completion Apple,
  sends FCM push to MJOLNIR. — Done 2026-06-11. Commit emily-agent (heimdal.go + cron.go).

---

## SECTION 18: MULTI-REPO RSI TOOLING + CONTEXT SPRAWL

- [x] **CHANGELOG updates getting dropped across repos** — Agents working in SHANKPIT, EMILY, [Apple #354, 2026-06-12]
  IDUNA, emily.cli, APPLES have no automated CHANGELOG enforcement. PRRJECT_FATBABY obs-watcher
  mentions CHANGELOG in step 5/6 but NOT in the mandatory runReportFooter. Fix: (1) add
  CHANGELOG.md update as step 5 in runReportFooter in cmd/observation-watcher/main.go so it is
  mandatory for all obs-watcher-dispatched runs; (2) add explicit CHANGELOG reminder to CLAUDE.md
  for every repo that doesn't have one. — obs `2026-06-12`.

- [x] **runReportFooter: add CHANGELOG as mandatory step** — Added CHANGELOG.md update as
  mandatory step 1 in runReportFooter (PRRJECT_FATBABY/cmd/observation-watcher/main.go).
  Steps renumbered 1-5. Build + tests pass. — 2026-06-12. Apple #349.

- [x] **Apple posts getting dropped in non-PRRJECT_FATBABY repos** — The `emily apples post` [Apple #354, 2026-06-12]
  command exists but agents in SHANKPIT, EMILY, IDUNA, emily.cli don't know to use it because
  it's not in their CLAUDE.md or any prompt footer. Fix: add "After any meaningful change, run:
  emily apples post -t completion -repo <REPONAME> <title>" to CLAUDE.md for SHANKPIT, EMILY,
  IDUNA, emily.cli, APPLES. — obs `2026-06-12`.

- [x] **API abstraction for cross-repo repeatable ops** — Everything that currently only works via [Apple #353, 2026-06-12]
  Claude Code should have a repeatable CLI/API surface: CHANGELOG update, Apple post, backlog
  mark-done, git commit+push. Pattern: emily CLI is the model. Extend emily.cli with:
  (a) `emily changelog add <repo> <message>` — appends a dated entry to <repo>/CHANGELOG.md;
  (b) `emily backlog done <item-id> [--apple-id N]` — marks an item [x] in EMILY/BACKLOG.md.
  These make the ops repeatable by any agent, not just Claude Code. — obs `2026-06-12`.

- [x] **Monorepo consideration (document trade-offs)** — Multi-repo sprawl causes context window [Apple #366, 2026-06-12]
  overload. A monorepo would consolidate context but Git history merge is complex and a single
  large repo may overwhelm context windows even more. Decision: investigate whether a partial
  consolidation (SHANKPIT + GoblinFoxDragon) or a workspace/submodule approach reduces sprawl
  without context penalty. Document trade-offs before acting. — obs `2026-06-12`.

- [x] **S18-TOOL-01: `emily context build` command** — Needed to compile all Tier 1 golden docs into [done 2026-06-13]
  `EMILY/context/full-system-context.md` on demand. No implementation exists. Gap documented here;
  implementation tracked in S22-05. Tooling is insufficient without this.

- [x] **S18-TOOL-02: `emily backlog add` and `emily backlog add-section`** — All backlog edits are [done 2026-06-13]
  currently manual. Emily Prime cannot programmatically add items via `/api/v1/emily/run`.
  Gap documented here; implementation tracked in S22-08.

- [x] **S18-TOOL-03: `emily northstar <repo>`** — No way to quickly read a repo's northstar from CLI. [done 2026-06-13]
  Gap documented here; implementation tracked in S22-09.

- [x] **S18-TOOL-04: Emily Prime planning endpoint** — All strategic planning happens in Claude Code [done 2026-06-13]
  conversations (expensive). Emily Prime has no `/plan` endpoint. Claude Code tokens should be spent
  on implementation, not planning. Gap documented; implementation tracked in S22-07.

- [x] **S18-AGI-01: obs-watcher --continue flag (AGI loop mode)** — Added -continue flag [2026-06-14]
  (env: OBSERVATION_CONTINUE=true) to observation-watcher. When enabled, appends --continue to
  every claude invocation so RSI cycles continue the prior session and build persistent context
  across iterations rather than starting fresh. emily start --agi wires this automatically.
  PRRJECT_FATBABY commit af5c76d, emily.cli commit 2ea4b9d. THE_EMILY_WAY.md principle 11.

- [x] **S18-AGI-02: THE_EMILY_WAY.md** — Comprehensive operating procedure doc encoding the [2026-06-14]
  Emily Way: 13 principles from backlog-first through AGI loop mode, synthesized from all
  CLAUDE.md files, git commits, and changelogs. Registered as Tier 1 golden doc (THE-EMILY-WAY).
  EMILY commit 8400afa.

- [x] **S18-AGI-03: Tier 2/3 golden-docs-index expansion** — Added 19 new entries to [2026-06-14]
  golden-docs-index.md: emily-tools, protocol, framework, cron-evo, integration, IAM spec,
  MJOLNIR-apples/push/spec, APPLES-schema, SHANKPIT-netcode/predict/bridge,
  EmilyOS-arch/memo/posture, GFD-NORTH. Total index now 39 sources (19 Tier 1 + 20 Tier 2).
  EMILY commit 8400afa.

- [~] **Ops docs token efficiency: multilingual compression experiment** — EXPERIMENTAL. Steps 1–2 done.
  Step 1: GOLDEN.md confirmed as sole runtime haiku context doc (~576 tokens, under 1200 budget).
  Large docs/ files are design docs, not runtime-loaded. Step 2: bilingual test version created at
  docs/compression-experiment/GOLDEN_BILINGUAL_TEST.md (~30–40% token reduction estimate).
  Test harness: scripts/compression-abtest.sh (ready to run). Step 3 BLOCKED: requires
  ANTHROPIC_API_KEY for haiku comprehension A/B test. Step 4 BLOCKED: depends on step 3.
  DO NOT deploy bilingual version until step 3 confirms equal comprehension.
  — obs `2026-06-12`. Steps 1–2 committed 2026-06-12.

---

## SECTION 13: GOLDEN DOCS + SYSTEM CONTEXT HYGIENE

- [x] **Golden docs audit** — CLAUDE.md added to EMILY, IDUNA, SHANKPIT, emily.cli, APPLES
  (5 missing repos; PRRJECT_FATBABY/TYLER/MJOLNIR already had them). All committed + pushed.
  Remaining gaps: northstars for PRRJECT_FATBABY/TYLER/EMILY/IDUNA still missing.
  — obs `2026-06-10T23:33:35Z`. Done 2026-06-11. Apple #334.

- [x] **emily prime API parity with emily.cli** — `POST /api/v1/emily/run { "command": "..." }`
  wired into emily-agent. Whitelist: observe/obs/apples/status/sync/backlog/primetask/agents/watch.
  Single-quote-aware arg split, 60s default / 300s max timeout. Returns exit_code/stdout/stderr/duration.
  — obs `2026-06-10T23:38:51Z`. Done 2026-06-11. Apple #335.

---

## SECTION 14: EMILYOS (bare-metal exokernel)

- [x] **EmilyOS northstar** — Already exists at EmilyOS/docs/NORTHSTAR.md (written 2026-06-09
  by Emily Prime). SOC 2 Type II framing, 6 milestones, posture/RBAC/audit design. The
  product is a Linux-based policy kernel, not bare-metal. — Done 2026-06-09.

- [x] **EmilyOS package repo + build system** — Choose: Debian base or Arch base? [Apple #374, 2026-06-12]
  Build SOC2-auditable software repository (BAZEL or BAZEL equivalent). Key constraint:
  reproducible builds, content-addressed storage. — obs `2026-06-11T00:00:18Z`.

---

## SECTION 15: PITVIPER (custom terminal)

- [x] **PITVIPER northstar** — PITVIPER/docs/NORTHSTAR.md written. SDL2 + FreeType2,
  5 milestones (bootstrap→PTY→glyph cache→splits→Emily pane), CommandRecord hook layer,
  JetBrains Mono embedded, emily:// pane type. PITVIPER/CLAUDE.md.
  — obs `2026-06-11T00:52:56Z`. Done 2026-06-11.

---

## SECTION 16: EMILY PRIME AI TIER (FABLE + API)

- [x] **FABLE advisor (basic)** — `GET /api/v1/emily/fable/advice` reads GOLDEN.md + recent
  Apples via haiku, returns `FableAdvice{recommendations[3], summary, generated_at}` with
  `FableItem{priority,title,rationale,section,effort}`.
  — obs `2026-06-10T23:59:02Z`. Done 2026-06-11. Apple #336.

- [x] **FABLE→HEIMDAL integration** — `POST /api/v1/emily/fable/execute` generates FABLE
  advice then files top recommendation as a HEIMDAL sprint in IDUNA; Emily Prime translates
  + queues it on the next cron cycle. Returns queued item + sprint_id.
  — obs `2026-06-10T23:56:50Z`. Done 2026-06-11. Apple #336.

- [x] **Emily Prime API** — Emily Prime needs a stable API so external orchestration
  (cron, HEIMDAL webhooks, MJOLNIR) can drive RSI loops without a human at the terminal.
  Split token usage: cheap haiku for classification, Sonnet/Opus only for implementation.
  — obs `2026-06-10T23:54:59Z`. Done 2026-06-11. Apple #346. Commit EMILY 4e6886b.

- [x] **emily-memory/ activation & world-state.md golden doc** — Completed as S27-01 through S27-03. world-state.md + cycle-log.md wired into goldenbuild + cron PLAN phase. 2026-06-14. Apple #513.
---

## SECTION 17: NEWSSITE + GTM (product growth)

- [x] **Newssite: stock charts** — Build an in-house charting library for the newssite.
  Goal: render equity price charts inline with governance articles. Start simple (SVG, no
  canvas deps), iterate. "Do it the Emily way." — obs `2026-06-11T00:51:58Z`.

- [x] **Newssite: use filing date not publication date** — docindex.AllSummaries() +
  ForTicker() now sort by FilingDate desc via docNewerThan(). ReadLatest fallback
  also sorted. Two regression tests added. — obs `2026-05-30T22:14:19Z`. Done 2026-06-11.
  Commit PRRJECT_FATBABY 284db04. Apple #333.

- [x] **Newssite: 500 errors investigation** — Root cause: `serveFrontPage` fell back to
  `ReadLatest` (full 191MB store scan) when `docIdx==nil` at startup; client timeouts → 500.
  Fix: 2s context bound on fallback scan; exceeded → 200 self-refreshing loading page.
  `emily start --newssite` and `--all` added (newssite was missing entirely).
  — obs `2026-06-07T21:43:40Z`. Done 2026-06-11. Apple #337.

- [x] **Newssite: Emily-authored governance commentary ingest endpoint** — Add an ingest
  endpoint so Emily Prime can POST original governance commentary articles directly to the
  newssite CMS. — obs `2026-05-30T21:53:49Z`.

- [x] **GTM funnel** — Full product funnel: Ask Emily free tier, Emily+ subscription, [Apple #377, 2026-06-12]
  community editorial engine, Merkle query monetization. Spec to be drafted.
  — obs `2026-05-30T22:02:43Z`.

- [x] **Self-improving training pipeline** — User data flywheel for Emily fine-tuning. [Apple #381, 2026-06-12]
  Collect prompt/response pairs from Emily Prime interactions, build annotation pipeline,
  RLHF loop. Long-term initiative. — obs `2026-05-30T22:10:22Z`.

---

---

## STRATEGIC PRIORITY ORDER (as of 2026-06-12)

Emily Prime read the golden docs sprawl memo and did deep thinking on value creation sequencing.
Revenue path analysis:

```
Track 1: SHANKPIT → Steam Early Access     (6-12 months, first dollars)
Track 2: Ask Emily (newssite product)       (6-12 months, freemium revenue)
  Blocker: FatBaby MySQL projections + MongoDB graph flattening
Track 3: Data platform licensing            (2+ years, institutional revenue)
RSI:     Emily Prime Brain                  (parallel enabler — reduces Claude Code token spend)
```

---

Priority section order (updated 2026-06-16): S29 → S30 → S23 → S31 → S32 → S19 → S33 → S34 → S5 → S2

*S29 (prod RSI cron) and S30 (signalapi prod) unblock everything downstream.*
*S23-01 (EDIS deploy) in progress — server live at iduna.farthq.com.*

---

## SECTION 19: SHANKPIT → STEAM EARLY ACCESS (revenue track 1)

*Northstar: ship a playable EA build with FPS core only. No Dragonfly yet. BedWars lands as major update post-launch.*
*The C client + Go server are both partially functional. The gap is client-side portal travel + playable standalone build.*

- [x] **S19-01: Client-side portal travel state machine** — The Go server sends `PacketSceneChange` (type=6)
  on portal transit. The C client ignores it. Implement travel state machine in C client:
  receive PacketSceneChange → freeze input → play transition → update scene_id → resume. This makes
  every scene in the game genuinely traversable. Dependency: Go server Milestone 2 ✓.
  Acceptance: player can walk into a portal and arrive in a new scene without desync.

- [x] **S19-02: Per-player physics isolation (C server)** — Currently physics is a global
  active-scene swap. Each player needs an independent physics context so players in different
  scenes don't step on each other. This is the last blocker for multi-scene co-existence in the
  C server. Acceptance: two players in different scenes can both interact with physics simultaneously.

- [ ] **S19-03: Steam Direct account + listing prep** — Create Steamworks developer account ($100 fee,
  human action required). Prepare minimum Steam page: 3 screenshots, 30-second capsule trailer,
  short description. Target: "Server-authoritative fast multiplayer shooter, Early Access."
  This can run in parallel with S19-01 and S19-02.
  BLOCKED: needs Steam account (human action).

- [x] **S19-04: Standalone playable EA build** — Makefile targets: go-server (GOWORK=off Go build),
  ea (Linux: lobby + go-server + README → dist/ea/), ea-windows (mingw cross-compile).
  docs/EA_BUILD.md: 4-player LAN/internet session guide, requirements, controls, EA status.
  — 2026-06-12. Apple #411. Commit 574f2b1.

- [ ] **S19-05: Steam EA launch** — Set price ($9.99 USD EA), upload build, enable EA store page.
  File completion Apple. Post observation.
  Dependency: S19-03 ✓, S19-04 ✓.

- [x] **S19-06: SHANKPIT NORTHSTAR: Steam milestones** — Update `SHANKPIT/docs2/NORTHSTAR.md` with
  Milestone 5 (EA launch: FPS core), Milestone 6 (BedWars + Dragonfly post-launch), Milestone 7
  (Season 1 lineage). Keeps agent context current with the launch path.

---

## SECTION 20: FATBABY QUERYABILITY (MySQL projections + MongoDB graph)

*FatBaby's eventstore is append-only and correct. It cannot serve ad-hoc queries needed for the newssite
product. Solution: CQRS read models — MySQL for relational projections, MongoDB for flattened entity graph.*
*User direction: MongoDB over Neo4j (no Neo4j ops complexity).*

- [x] **S20-01: MySQL read model schema design** — Design the projection tables in `PRRJECT_FATBABY/docs/`:
  `governance_signals` (ticker, event_type, date, entity, filing_id, raw_signal),
  `eps_results` (ticker, period, eps_actual, eps_estimate, surprise_pct),
  `entity_timeline` (ticker, entity_name, role, event_type, event_date, source_filing).
  Write as SQL DDL + rationale doc. No implementation yet — spec first.
  Acceptance: doc written, reviewed by Emily Prime, committed.

- [x] **S20-02: MySQL projector in PRRJECT_FATBABY** — cmd/projector tails secwatch eventstore,
  projects signal_generated events into governance_signals + entity_timeline tables. 4 SQL
  migrations in migrations/mysql/. Idempotent on eventstore_seq. Graceful degradation when
  MYSQL_URL not set. go-sql-driver/mysql v1.10.0 added. — 2026-06-12. Apple #410. Commit 74f09a2.

- [x] **S20-03: MongoDB entity document schema** — Design the flattened entity document format.
  One MongoDB document per tracked entity (company/director/auditor). Fields:
  `{ ticker, name, entity_type, directors: [...], governance_events: [...], eps_history: [...],
  signal_score, last_updated }`. This flattens the EAV graph into queryable JSON documents.
  Avoids Neo4j ops entirely. Write schema doc + rationale.
  Acceptance: schema committed, Emily Prime reviewed.

- [x] **S20-04: MongoDB entity writer in PRRJECT_FATBABY** — internal/mongowriter: EntityDocument
  per ticker from Graph+signals, upserted via ReplaceOne. entity-graph --mongo-url flag.
  Graceful no-op when MONGODB_URL not set. — 2026-06-12. Apple #412. Commit 8665179.

- [x] **S20-05: signalapi: query endpoints over MySQL + MongoDB** — GET /v1/governance-signals
  (ticker/type/since/until/limit filter → MySQL), /v1/eps/{ticker} (MySQL eps_results),
  /v1/entities/{ticker} (MongoDB). MYSQL_URL + MONGODB_URL env activation; 503 when not
  configured. — 2026-06-12. Apple #413. Commit f6e4233.

- [x] **S20-06: MySQL/MongoDB local dev setup doc** — docs/local-dev-setup.md: Docker quickstart,
  env vars, projector + entity-graph one-shot, query examples, reset/reseed, docker-compose,
  troubleshooting. — 2026-06-12. Apple #413. Commit f6e4233.

---

## SECTION 21: ASK EMILY PRODUCT (revenue track 2)

*Consumer-facing intelligence layer. Free tier: 5 questions/day. Paid: unlimited + email digest.*
*Sits on top of FatBaby signals (S20 queryable) + Emily Prime API (already exists at :8086).*
*Full GTM spec: `PRRJECT_FATBABY/docs/GTM_FUNNEL.md`.*

- [x] **S21-01: Ask Emily chat endpoint on newssite** — POST /api/ask proxies to Emily Prime
  /chat. Ticker prefix in message. 30s timeout. 503 when emily not configured.
  --emily-url flag + EMILY_BASE_URL env. — 2026-06-12. Apple #414. Commit 967ca3d.

- [x] **S21-02: Ask Emily chat UI on newssite** — Sidebar widget: ticker input + question textarea
  + Ask button. Async JS fetch /api/ask. Answer/error display inline. Mobile-friendly.
  — 2026-06-12. Apple #414. Commit 967ca3d.

- [x] **S21-03: Ask Emily: wire FatBaby signals into Emily Prime context** — When `/api/ask` receives [done 2026-06-13]
  a question with a ticker, fetch the entity document (S20-05) and governance_signals for that ticker,
  prepend as context to Emily Prime's user message. Emily answers from real data, not hallucination.
  Dependency: S20-05 ✓, S21-01 ✓.

- [x] **S21-04: Ask Emily: rate limiting (free tier)** — 5 questions per IP per day. Use IDUNA JWT [done 2026-06-13]
  if user is logged in (then tie to user_id). Anonymous: use IP + daily bucket in Redis or SQLite.
  Enforcement in newssite handler, not Emily Prime.

- [x] **S21-05: Ask Emily: auth integration via IDUNA** — Google OAuth login on newssite via IDUNA [done 2026-06-13]
  `/api/v1/auth/google`. Logged-in users get 20 questions/day (free tier). Subscription tier TBD.
  Dependency: S21-04 ✓, IDUNA OAuth ✓.

- [x] **S21-06: Landing page + waitlist** — Simple landing page at root of newssite (or separate
  subdomain). Copy: "Ask Emily — governance intelligence for active investors." Email capture for
  waitlist. Mailchimp or simple SMTP to emilyspringerton@gmail.com. Apple #416

---

## SECTION 22: EMILY PRIME BRAIN (goldenbuild + dynamic prompt + FABLE full context)

*Emily Prime currently runs blind — static system prompt, FABLE reads only GOLDEN.md.*
*This section builds the context infrastructure that lets Emily Prime actually understand and plan the system.*
*Completing this section is the multiplier for all other sections.*

- [x] **S22-01: `goldenbuild.go` — continuous golden doc compiler** — New file in `emily-agent/`.
  `GoldenDocCompiler` reads all Tier 1 golden docs from all repos, compresses each via claude-haiku
  bilingual Chinese/English (≤150 tokens/source), writes to `EMILY/context/full-system-context.md`.
  `MaybeRebuild(ctx)`: checks source mtimes vs output mtime, rebuilds if any source is newer.
  Sources: 18 Tier 1 docs listed in `context/golden-docs-sprawl-memo-2026-06-12.md`.
  Acceptance: `full-system-context.md` written, contains compressed context from every repo.

- [x] **S22-02: Dynamic Emily system prompt** — Change `const emilySystemPrompt` in `main.go` to
  `buildEmilySystemPrompt(emilyRoot string) string`. Reads `context/full-system-context.md` if
  present, prepends to static roles/tools section. Updates 3 call sites (lines 1614, 1648, 1698).
  Acceptance: Emily Prime's system prompt includes full cross-repo context on every conversation.

- [x] **S22-03: Wire MaybeRebuild into cron cycle** — Call `GoldenDocCompiler.MaybeRebuild(ctx)`
  at the start of `RunOnce()` in `cron.go`. Full context refreshes on each 5-min cron cycle when
  any source doc has changed.
  Dependency: S22-01 ✓.

- [x] **S22-04: FABLE reads full context** — Change `fable.go` `handleFableAdvice` and
  `handleFableExecute` to pass `full-system-context.md` (or append to `GOLDEN.md`).
  FABLE recommendations account for all repo northstars, not just backlog state.
  Dependency: S22-01 ✓.

- [x] **S22-05: `emily context build` CLI command** — Add `emily context build` to emily.cli.
  Runs the goldenbuild compiler on demand: reads all golden sources, compresses, writes
  `full-system-context.md`. Allows human-triggered refresh without waiting for cron.
  Dependency: S22-01 ✓ (or implement directly in emily.cli without importing emily-agent).

- [x] **S22-06: `EMILY/docs/NORTHSTAR.md`** — Emily's own repo has no NORTHSTAR.md. Write one
  that synthesizes `emily-prime-spec.md` + `emiree-emily-fatbaby.md` + `emiree.md` into a single
  canonical document. This becomes the Tier 1 golden doc for EMILY itself.

- [x] **S22-07: Emily Prime `/api/v1/emily/plan` endpoint** — Accept a planning question [done 2026-06-13]
  (`{ question: string, context?: string }`), return a structured sprint batch
  (`{ sprints: [{ title, rationale, section, effort }], summary }`).
  This is the upgrade from FABLE (fixed 3 items) to a full planning conversation.
  Moves planning responsibility from Claude Code sessions to Emily Prime herself.
  Dependency: S22-02 ✓, S22-04 ✓.

- [x] **S22-08: Tooling gap — `emily backlog add` and `emily backlog add-section`** — Currently [done 2026-06-13]
  adding items to BACKLOG.md requires manual editing. Add CLI commands:
  `emily backlog add --section N "<item text>"` — appends item to section N
  `emily backlog add-section --title "<title>"` — appends new numbered section
  These make backlog editing programmatic and Emily Prime can drive it via `/api/v1/emily/run`.

- [x] **S22-09: `emily northstar <repo>`** — Print the northstar doc for the named repo. [done 2026-06-13]
  Reads from the canonical location (e.g., `<repo>/docs/NORTHSTAR.md` or `<repo>/docs2/NORTHSTAR.md`).
  Used by Emily Prime and human operators to quickly orient.

- [x] **S22-10: Add EMILY + EDIS northstars to goldenbuild.go source list** — Added EMILY-NORTH
  (EMILY/docs/NORTHSTAR.md) and EDIS (EDIS/NORTHSTAR.md) to GoldenDocCompiler sources in both
  emily-agent/goldenbuild.go and emily.cli/cmd/context.go. Both lists now 18 sources. Builds clean.
  — 2026-06-13. Apple #469.

- [x] **S22-11: `EMILY/context/golden-docs-index.md` — manifest-driven source list** — Created
  EMILY/context/golden-docs-index.md (18 Tier 1 sources, markdown table). goldenbuild.go
  loadGoldenIndex() reads it at startup; hardcoded fallback if absent. Same loader in
  emily.cli/cmd/context.go. Adding a new golden doc = append one row, no Go change.
  — 2026-06-13. Apple #470.

- [x] **S22-12: New-doc registration protocol** — Step 6 added to obs-watcher runReportFooter:
  agents append row to golden-docs-index.md when creating mission-critical docs. Golden Doc
  Registration section added to CLAUDE.md for EMILY, IDUNA, SHANKPIT, emily.cli, APPLES, EDIS.
  — 2026-06-13. Apple #471.

- [x] **S22-13: Unblock multilingual compression A/B test (S18)** — FAIL: bilingual degrades Q1 (stale [Apple #521, 2026-06-15]
  sections) and Q3 (drops 4 repos). Token reduction only 7% not 30-40%. Decision: keep GOLDEN.md.
  Efficiency win already live via emily context build (full-system-context.md bilingual, 108KB→23KB).
  — 2026-06-15. Apple #521.

- [x] **S22-14: API key gate — document the unlock sequence** — Wrote EMILY/docs/API_KEY_UNLOCK.md:
  6-step sequence (set key → restart → context build → compression A/B test → verify FABLE →
  verify cron log). Degraded-mode table included. Cost note: haiku-only, SonnetModel unused.
  — 2026-06-13. Apple #472.

---

## SECTION 25: GOLDEN DOC NAMING CONVENTION (Phases 4–6 from 2026-06-12 memo)

*These are hygiene items from the golden-docs-sprawl-memo. Not blockers for revenue, but required*
*for Emily Prime to have accurate cross-repo context as the system scales.*
*Priority: after S22 items above. Emily Prime executes these in background RSI cycles.*

- [x] **S25-01: IDUNA/docs/NORTHSTAR.md** — Written. IAM + Apples + HEIMDAL systems, endpoint
  table, architecture, status/gaps. Registered in golden-docs-index.md as IDUNA-NORTH Tier 1.
  — 2026-06-13. Apple #473.

- [x] **S25-02: GoblinFoxDragon/docs/NORTHSTAR.md** — GFD is the R&D studio umbrella and
  DragonsNShit engine; SHANKPIT is the first shipped product derived from GFD DNA. Both have
  module=dragonsnshit because GFD IS the Dragonfly fork. NORTHSTAR written: milestones,
  package structure, EduScript VM, VS0 specs, bridge relationship to SHANKPIT, IDUNA identity.
  Registered in golden-docs-index.md as Tier 2 (GFD-NORTH). GFD commit 287733b.
  — 2026-06-14.

- [x] **S25-03: Archive EmilyOS legacy docs** — Moved EmilyOS/docs/legacy/ (12 files) to
  EmilyOS/docs/legacy-archive/. README.md added: superseded by NORTHSTAR.md 2026-06-09.
  — 2026-06-13. Apple #474.

- [x] **S25-04: Resolve SHANKPIT/GFD doc duplication** — NETCODE_CONTRACT_SPEC.md and
  CLIENT_PREDICTION_SPEC.md replaced in GFD with reference stubs pointing to canonical
  SHANKPIT/docs2/ sources. GFD world-backend extensions go in addendum docs under
  docs2/specs/. Prevents spec drift. GFD commits 287733b + 8760b31.
  — 2026-06-14.

- [x] **S25-05: Move EMILY/THE_FIELD.md into docs/** — Moved to EMILY/docs/THE_FIELD.md.
  Added THE-FIELD row to golden-docs-index.md as Tier 2.
  — 2026-06-13. Apple #475.

---

## S26 — GPT-2 Emily Prime Training Pipeline

*Opened: 2026-06-14. Initiative: fine-tune GPT-2 small on Emily's golden docs + prime directive to
produce an Emily-domain language model for local inference and RSI entropy sourcing.*

- [x] **S26-01: IDUNA Drive API** — `/api/v1/drive/*` endpoints (upload, list, get); stdlib-only
  service account JWT auth (RS256); EMILY-TRAINING agent registered; degraded-mode 503 when env
  not set. Commit IDUNA 58dc2c4. — 2026-06-14.

- [x] **S26-02: gpt2-alpine-c training tooling** — scripts/prime_directive_dataset.py (JSONL
  corpus builder from golden docs + prime directive + RSI task history), scripts/drive_sync.py
  (IDUNA Drive API upload/download), scripts/convert_ft_checkpoint.py (HF checkpoint → C binary),
  notebooks/gpt2_finetune_colab.ipynb (HF Trainer fine-tune on Colab GPU), NORTHSTAR.md, CLAUDE.md,
  CHANGELOG.md, .gitignore. Commit gpt2-alpine-c e3c4e7e. — 2026-06-14.

- [x] **S26-03: emily train CLI command** — `emily train build-dataset`, `emily train upload`,
  `emily train status`; Drive API methods added to internal/iduna/client.go (DriveUpload, DriveList,
  DriveGet). Wired into main.go. — 2026-06-14.

- [ ] **S26-04: First experimental fine-tune on Colab** — Run notebooks/gpt2_finetune_colab.ipynb
  on T4 GPU with Emily corpus; validate perplexity < GPT-2 base on Emily domain text; download
  checkpoint and convert to C binary; run `./gpt2_run weights/emily-ft.bin --entropy-stats`.
  Manual step — requires Google Colab session.
  Corpus ready: `emily train build-dataset` → 466 records / 154k tokens / 2.2 min T4 (Apple #518).
  Entropy fix: H_mean=4.49 nats base (was 0.0 — float32 underflow fixed, Apple #519, commit 4ab6bf0).
  Smoke test (20 steps CPU): loss 6.40→5.73, fine-tuned H_mean=4.6602 vs base 4.4877 (+0.17 nats). Apple #520.
  Pipeline validated end-to-end. Full fine-tune requires Colab T4 GPU.
  2026-06-16: Git LFS set up; emily-ft.bin + model.bin + model.safetensors pushed to GitHub LFS (Apple #589).
  2026-06-16: Local 300-step LoRA run started (PID 58912, 578 corpus records, bfloat16 CPU). IN PROGRESS.

- [x] **S26-05: Validate entropy source** — Base: H_mean=4.4877, fine-tuned (emily-ft.bin 300-step): H_mean=4.6602 (+0.17 nats). NORTHSTAR milestone 3 updated. Target ≥0.5 nats not yet met — Colab T4 full fine-tune (S26-04) needed to reach target. RSI entropy source functional. Apple #3093 | 2026-06-23

- [x] **S26-06: Emily Prime self-description training data** — Build instruction fine-tune pairs
  from Emily's prime directive in instruct mode (`--mode instruct`); verify that fine-tuned model
  completes Emily-domain prompts correctly; establish perplexity baseline for future comparison.
  — 2026-06-14. Apple #488. Commit gpt2-alpine-c 42b0d66.
  24 prime-directive instruct pairs (8 identity/protocol + section Q&A). Base GPT-2 PPL: 130.33.
  Post-fine-tune target: PPL < 60. Instruct corpus: 128 pairs / 455 records / 568KB.

---

## INTAKE QUEUE (curated by emily backlog curate)

Items below require `ANTHROPIC_API_KEY` for haiku routing or manual triage.
Run: `emily backlog promote --limit=50 --batch=15`

- [x] **UX: ticker search should auto-navigate on click and on Enter key — remove redundant Go button step** — Form submit + datalist selection now navigate directly to /ticker/{q}; /search intermediate step eliminated. 'change' event added for cross-browser compat. — obs `2026-05-31T20:39:51Z`. Done 2026-06-11. Apple #339.
- [x] **Feature: create a GitHub issue automatically whenever Emily writes an observation** — obs `2026-05-31T20:41:56Z`. Done 2026-06-11. Apple #343. Commit PRRJECT_FATBABY b2f52bc.
- [x] **All required and optional environment variables must be documented at the top of the README** — obs `2026-05-31T20:43:15Z`. CURATED: 2026-06-11. Done 2026-06-12. Apple #351. Commit gpt2-alpine-c 8348d1f. Apple #342 2026-06-11.
- [x] **observation-watcher must inject full reporting and git sync requirements into the Claude Code prompt** — runReportFooter updated: Apple now required (not optional), BACKLOG.md marking step added, git add -A. — obs `2026-05-31T20:56:45Z`. Done 2026-06-11. Apple #340.
- [x] **EDGAR submissions endpoint returning truncated JSON for all 5 major bank tickers — BAC, C, GS, JPM, MS** — maxBodyBytes increased from 64MB to 256MB in secwatch/client.go. — obs `2026-05-30T21:45:32Z`. Done 2026-06-11. Apple #341.
- [x] **entity-graph cannot detect 8-K documents from persisted store — form/source_type field mismatch** — buildFilingIndexes now returns form type map; 4th detection path via filing_discovered events recovers 880/1104 legacy docs. — obs `2026-05-30T21:32:12Z`. Done 2026-06-11. Apple #338.
- [~] **entity-graph parsing all 8-K subtypes — Item 5.07 not found in non-proxy filings** — INVESTIGATED 2026-06-08: FALSE ALARM. TestParseItem507_BALiveFixture passes with 6 proposals on real BA data. Recent batches process non-proxy 8-Ks (earnings etc.) that legitimately have no Item 5.07. Not a bug. — obs `2026-05-30T21:29:24Z`. Closed.
- [x] **entity-graph reads 0 filings despite 846 source documents in var/secwatch** — Stale: cursor at 67657 past store end (67656). Graph has 167 nodes, 1068 signals, 1215 edges. 1417 8-K docs processed. Also fixed by Apple #338 (legacy detection). — obs `2026-05-30T09:52:12Z`. Done.
- [~] **eps-processor ticker map has only 2 entries — all press releases are being dropped silently** — STALE: ticker map now has 341 entries (was 2 on 2026-05-30). Processor cursor at 1562/1561 (fully caught up). 0 articles because current corpus lacks extractable earnings releases — data issue, not code. — obs `2026-05-30T09:46:52Z`. Closed.
- [x] **gpt-2 c fork git@github.com:emilyspringerton/gpt2-alpine-c.git parity with og gpt-2 repo (we may have to build tensorflow…** — obs `2026-06-10T00:50:04Z`. CURATED: 2026-06-11. Done 2026-06-12. Apple #351. Commit gpt2-alpine-c 8348d1f.
- [x] **gpt-2 as an entropy source** — obs `2026-06-10T00:48:37Z`. Done 2026-06-12. Apple #351. --entropy-stats outputs entropy_mean_nats for RSI loop use. Done 2026-06-12. Apple #351. Commit gpt2-alpine-c 8348d1f.
- [~] **S17 newssite filing date bug fixed: sort by FilingDate descending, not ingestion timestamp** — session completion recap. Closed: already done in prior session.
- [~] **S6 TUI fatbaby mode: --fatbaby flag + b-toggle showing entity graph, process status, eps counts** — session completion recap. Closed: already done in prior session.
- [~] **S3+S6+S12 sprint: HEIMDAL feedback, Apples auto-sync, golden context feed, obs-watcher rate-limit resilience** — session completion recap. Closed: already done.
- [~] **session 2026-06-11: emily backlog curate CLI + Emily Prime agent tools (emily_read/write_file) + autonomous triage cura…** — session completion recap. Closed: already done.
- [x] **we need the iduna middleware to lockdown emily read file and emily write file such that apples are always filed and emi…** — emily_write_file now unconditionally POSTs Apple via IdunaClient inside tool handler (Emily cannot opt out). OpenAPI spec for all 13 emily-agent tools written. Future tool stubs (grep_files, shell_exec, run_tests, git_commit, git_push) documented as path toward own Claude Code abstraction. EMILY commit 25e6c83. Done 2026-06-12. Apple #386.
---

- [x] **S29-05 smoke: obs-watcher dispatch verification — RSI loop end-to-end test** — Covered by S29-05 above. Apple #848.
- [x] **S29-05 final smoke: single-obs dispatch test — confirm obs-watcher picks up and dispatches to Claude** — Covered by S29-05 above. Apple #848.
## SECTION 23: EDIS — WORDPRESS INTELLIGENCE PRODUCT (public face of FatBaby)

*Northstar: WordPress site with three plugins that call signalapi. SEO-optimized, community-ready.*
*Plugins: edis-core (API client + cache), edis-signals (shortcodes), edis-ask-emily (Ask Emily widget).*
*Repo: EDIS/ — scaffold complete 2026-06-12.*

- [x] **S23-00: EDIS repo scaffold** — CLAUDE.md, NORTHSTAR.md, CHANGELOG.md, three plugins (edis-core,
  edis-signals, edis-ask-emily), starter theme (edis), docs (architecture.md, wordpress-setup.md).
  edis-core: WP HTTP API client + transient cache + admin settings. edis-signals: [edis_signals],
  [edis_entity], [edis_eps] shortcodes + sidebar widget. edis-ask-emily: [ask_emily] shortcode +
  WP REST POST /wp-json/edis/v1/ask + sidebar widget. — 2026-06-12. Apple #418.

- [x] **S23-06: Digital Immune System (DIS)** — Go ops posturing layer from golden.md spec.
  ring.go (16k-slot lock-free ring buffer), fingerprint.go (header-order hash + session HMAC),
  harvester.go (Go HTTP middleware), posture.go (30s rolling health state: healthy/elevated/attack/degraded),
  adengine.go (health→ad mode selector), cmd/dis/main.go (nginx log-tailing daemon, :9099).
  edis-dis WordPress plugin: reads health from collector, [edis_dis_ad] shortcode, admin posture panel.
  emily install --edis CLI command: full-stack EDIS provisioner (nginx+PHP+WordPress+plugins+theme+custodian agent).
  EDIS-CUSTODIAN agent added to IDUNA/config/agents.json. go.work updated. docs/digital-immune-system.md written.
  — 2026-06-12. DIS pre-deploy analysis: Apple #450 (2026-06-13). 3 bugs fixed: Apple #591 (2026-06-16).

- [ ] **S23-01: Deploy EDIS to live WordPress install** — Install WordPress, activate EDIS plugins
  and theme. Configure EDIS_SIGNALAPI_URL + EDIS_EMILY_URL in wp-config.php. Verify connection test
  passes in admin. Create /ask page. Add Ask Emily + Signals widgets to sidebar.
  Acceptance: /ticker/AAPL shows live governance signals. /ask widget returns Emily answer.
  Dependency: signalapi live (S20-05 ✓). Use: emily install --edis --domain edis.example.com --dry-run

- [x] **S23-02: SEO + OpenGraph wiring** — Install Yoast SEO. Add OpenGraph meta to ticker pages [done 2026-06-13]
  (company name, signal count, score as description). Sitemap includes /ticker/ pages.
  Ticker page title: "{TICKER} Governance Intelligence — FatBaby".
  Acceptance: Google Search Console shows ticker pages indexed.

- [x] **S23-03: Ticker auto-pages** — For each ticker in FatBaby watchlist, auto-create or render [done 2026-06-13]
  a /ticker/{SYM} page. Options: (a) WordPress page per ticker with custom field; (b) virtual page
  via rewrite rule (already implemented in theme). Decide + implement. No duplicate content penalty.
  Acceptance: /ticker/AAPL, /ticker/JPM, /ticker/MSFT all render without 404.

- [x] **S23-04: Emily+ subscription gate (WooCommerce)** — IDUNA: user_subscriptions table
  (202606140002), auth.Subscription + IsActive(), GetEffectivePermissions auto-appends cap.query.full,
  SubscriptionHandler POST /api/v1/subscriptions + GET /me. newssite: VerifyTokenPermissions on
  askVerifier, cap.query.full bypasses all Ask Emily rate limits. EDIS: emily-plus-woocommerce.php
  WooCommerce order complete hook provisions IDUNA subscription, REST /wp-json/edis/v1/set-iduna-user.
  IDUNA 5ed080e, PRRJECT_FATBABY d3136fc, EDIS f5dde79. Apple #477. — 2026-06-14.

- [x] **S23-05: Mailchimp waitlist integration** — mailchimp-waitlist.php in edis-ask-emily: overrides
  /wp-json/edis/v1/waitlist (priority 20), Mailchimp API v3 PUT upsert, auto-detect data center from
  key suffix, tag support, wp_options fallback, admin settings fields. Config: EDIS_MAILCHIMP_API_KEY +
  LIST_ID + TAG. EDIS commit 22b7de5. Apple #478. — 2026-06-14.

---
- [ ] **S23-01: LIVE DEPLOY — sprint-deploy.sh ready to run. Requires: sudo bash /home/fatbaby/EDIS/ops/sprint-deploy.sh. Pre…**

## SECTION 24: NEWSSITE OPS HARDENING (traffic + production readiness)

*Context: once EDIS and Ask Emily get traction we'll get hit hard. Need to be ready before not after.*
*Ops infrastructure scaffolded 2026-06-12. Deploy steps documented in docs/ops-runbook.md.*

- [x] **S24-00: Ops scaffold** — nginx configs (newssite + signalapi), systemd service units
  (newssite, processor, secwatch, signalapi), docker-compose.prod.yml (MySQL + MongoDB + nginx),
  deploy.sh build+restart script, env.production template, ops-runbook.md.
  — 2026-06-12. Apple #419.

- [ ] **S24-01: Deploy to production server** — Run deploy.sh on production host. Install nginx
  configs. Start services under systemd. Confirm /healthz returns 200 through nginx.
  Set ANTHROPIC_API_KEY, MYSQL_URL, MONGODB_URL in env.production.
  Acceptance: https://fatbaby.io/healthz returns "ok". Emily TUI shows all services active.

- [ ] **S24-02: SSL certificate + domain wiring** — Let's Encrypt cert for fatbaby.io and api.fatbaby.io.
  Use certbot with nginx plugin. Auto-renew via cron. Update nginx configs with cert paths.
  Acceptance: HTTPS works, cert grade A on SSL Labs.

- [x] **S24-03: Log rotation + alerting** — logrotate already configured (PRRJECT_FATBABY/ops/logrotate/fatbaby: daily 14d compress maxsize 500M). Emily Prime watchdog (watchdog.go): pings IDUNA/newssite/signalapi/emily-agent each cycle; escalation Apple when service down ≥ 2 min; log file size alert at 500 MB; WatchdogState persisted across cycles. Apple #479.

- [ ] **S24-04: nginx cache tuning post-traffic** — After first traffic spike, review cache hit rate.
  Tune proxy_cache_valid TTLs based on actual traffic patterns. Target >80% cache hit rate.
  Acceptance: cache hit rate logged and documented.

- [ ] **S24-05: Load test baseline** — Run `wrk -t4 -c50 -d30s https://fatbaby.io/` before real traffic.
  Record: req/s, p99 latency, cache hit rate. Document results in ops-runbook.md.
  Acceptance: baseline documented; known ceiling before we hit it in production.

---

## SECTION 27: EMILY AGI MEMORY (persistent cross-cycle world model)

*Emily Prime starts cold every 5 minutes. emily-memory/ gives her accumulated context across cycles.*
*Gap identified in RSI AGI trajectory memo (#457): emily-memory/ was empty since repo creation.*
*Resolved 2026-06-14: world-state.md + cycle-log.md wired into goldenbuild + cron PLAN phase.*

- [x] **S27-01: Create emily-memory/world-state.md** — Structured world-state artifact: all product
  statuses, human blockers, AGI gaps, recent wins, next priorities. Registered as Tier 1 golden doc
  (EMILY-MEMORY) so goldenbuild compiles it into Emily Prime's system context each cycle.
  Resolves AGI Gap 1 from trajectory memo #457. — 2026-06-14. Apple #513.

- [x] **S27-02: Create emily-memory/cycle-log.md** — Rolling log of RSI cycle outcomes (last 100 entries).
  Registered as Tier 1 golden doc (EMILY-CYCLE-LOG). Emily Prime reads it each cycle via goldenbuild
  to know what happened in recent cycles without human intervention. — 2026-06-14. Apple #513.

- [x] **S27-03: Wire updateCycleLog() into cron.go PLAN phase** — New updateCycleLog() function
  appends timestamped cycle summary to emily-memory/cycle-log.md after each cron cycle. Bounded to
  100 entries via sentinel markers. Non-fatal: memory failure never breaks the cron loop. — 2026-06-14. Apple #513.

- [x] **S27-04: Cross-domain synthesis RSI preset** — New roadmap item in defaultRoadmap(): reads
  FatBaby/TYLER/SHANKPIT golden docs weekly, synthesizes cross-domain insights, files observation Apple.
  Resolves AGI Gap 2 (three intelligence domains siloed). Adds weekly cross-domain preset to cron.go.
  — 2026-06-15. Apple #522.

- [x] **S27-05: Revenue signal → RSI priority wiring** — When Apple type=completion includes a revenue
  tag (steam_launch, ask_emily_subscription, etc.), bump priority of adjacent roadmap items.
  Resolves AGI Gap 3 (revenue signals don't feed RSI priority). Wire in cron.go PLAN phase.
  — 2026-06-15. Apple #523.

- [x] **S27-06: Emily Prime world-state update protocol** — Emily Prime should update world-state.md
  after major milestones (not just cycle-log.md). Add a `/api/v1/emily/memory/update` endpoint that
  accepts a section name + new content and patches the relevant section in world-state.md. This makes
  Emily Prime self-describing and self-updating across cold restarts.
  — 2026-06-15. Apple #524.

---

## S28 — GPT-2 Inference Layer + System Swagger

**Goal:** Expose the fine-tuned GPT-2 model as a production-ready service behind the FatBaby broker
proxy on :8679, wire it into the emily CLI, and document the full EINHORN_INDUSTRIAL API surface
as a single OpenAPI 3.0.3 spec in EMILY.

- [x] **S28-01: Broker routes config for GPT-2 proxy on :8679** — `gpt2-alpine-c/config/broker-routes.json`
  defining tenant route from :8679 → :8088 using the FatBaby broker package. Bearer token auth.
  — Done 2026-06-15. Apple #539.

- [x] **S28-02: `emily gpt2` CLI command** — `emily gpt2 start|proxy|status|tokenizer` subcommands.
  start: spawns serve.py on :8088. proxy: spawns broker on :8679. status: pgrep both. tokenizer: make tokenizer.
  — Done 2026-06-15. Apple #540.

- [x] **S28-03: OpenAPI 3.0.3 Swagger — full system** — `EMILY/docs/api.yaml` documenting all
  EINHORN_INDUSTRIAL API endpoints: emily-agent (:8086), IDUNA (:8080), GPT-2 serve.py (:8088),
  GPT-2 broker proxy (:8679). Registered as Tier 2 golden doc.
  — Done 2026-06-15. Apple #541.

---

## SECTION 29: PRODUCTION RSI LOOP (cron + systemd auto-start)

*Context: user confirmed 2026-06-16 — cron is on the horizon. Server is live at iduna.farthq.com.*
*Goal: Emily Prime runs autonomously on the production server without any manual invocation.*

- [x] **S29-01: systemd service for emily-agent** — `emily install --system --write` on production
  server. Enables emily-agent (:8086) on boot. Requires ANTHROPIC_API_KEY in EnvironmentFile.
  Acceptance: `systemctl status emily-system` shows active; emily-agent survives reboot.
  ✓ Apple #562 — 2026-06-16. Service file written to ~/.config/systemd/user/emily-system.service.
  RUN: `! systemctl --user daemon-reload && systemctl --user enable --now emily-system.service`

- [x] **S29-02: ANTHROPIC_API_KEY in production env** — Write key to `/etc/emily/env` (mode 0600)
  and reference from systemd unit EnvironmentFile. Required for goldenbuild compression + haiku calls.
  Acceptance: emily-agent log shows "goldenbuild: compressed N sources" after restart.
  ✓ Apple #562 — 2026-06-16. config.Resolve() reads EMILY/var/emily-secrets.env automatically;
  no EnvironmentFile needed in systemd unit. Key already present in emily-secrets.env.

- [x] **S29-03: IDUNA systemd service on production** — `emily install --iduna-systemd --write`.
  IDUNA must start before emily-agent. Use After=iduna.service in emily-system.service.
  Acceptance: IDUNA health at https://iduna.farthq.com/api/v1/health returns 200 on reboot.
  ✓ Apple #562 — 2026-06-16. Service file written to ~/.config/systemd/user/iduna.service.
  RUN: `! systemctl --user daemon-reload && systemctl --user enable --now iduna.service`

- [x] **S29-04: obs-watcher systemd service** — obs-watcher currently started manually. Add to
  emily-system.service or as separate unit. Polls EMILY/signals/tasks/ every 10s.
  Acceptance: `emily status` shows obs-watcher running without manual start.
  ✓ Apple #562 — 2026-06-16. obs-watcher is started by `emily start` which is called by
  emily-system.service ExecStart. No separate unit needed.

- [x] **S29-05: End-to-end RSI loop smoke test** — File a test observation, verify obs-watcher
  dispatches, Apple is filed to IDUNA, cycle-log.md is updated. One full loop with no human touch.
  Acceptance: Apple filed, cycle-log shows entry, no manual intervention required.
  ✓ Apple #848 — 2026-06-16. 3 bugs fixed (cursor format, isContextTooLongOutput stdout capture,
  go run . compile). Final claude dispatch blocked by API credit balance (human: top up).

---

## SECTION 30: NEWSSITE + SIGNALAPI PRODUCTION DEPLOY

*Context: production server NOW LIVE (iduna.farthq.com — resolved 2026-06-16).*
*EDIS plugins call signalapi for data. Without signalapi in production, /ticker/AAPL shows empty.*
*Newssite is the internal ops tool; signalapi is the data layer both newssite and EDIS depend on.*

- [x] **S30-01: signalapi production deploy** — Binary built at `PRRJECT_FATBABY/bin/signalapi`.
  ops/deploy.sh already builds + installs systemd unit (fatbaby-signalapi.service → :9091).
  Apple #569 | 2026-06-16

- [ ] **S30-02: MySQL + MongoDB in production** — Point signalapi at production MySQL (already
  running) and optionally MongoDB. `MYSQL_URL` + `MONGODB_URL` env vars in emily-agent env.
  Run projector once to seed MySQL from event store.
  Acceptance: `GET /v1/governance-signals?ticker=AAPL` returns records.

- [x] **S30-03: emily start --signalapi on production** — `emily start --signalapi` launches
  signalapi detached on :9091, logs to var/logs/signalapi.log, pgrep-idempotent.
  Apple #569 | 2026-06-16

- [x] **S30-04: nginx proxy for signalapi** — `/signals/` location block in EDIS edis.conf:
  `rewrite ^/signals(/.*)$ $1 break; proxy_pass http://127.0.0.1:9091` with 60s GET cache,
  CORS headers, 20r/s rate limit. RUN: `sudo nginx -t && sudo systemctl reload nginx`
  Apple #569 | 2026-06-16

---

## SECTION 31: PRODUCTION MONITORING (DIS 24h + cache tuning)

*DIS collector is deployed with sprint-deploy.sh. First 24h checklist from dis-deployment-analysis.md.*
*These are ops validation items — run after S23-01 deploy completes.*

- [ ] **S31-01: DIS 24h monitoring checklist** — Run through dis-deployment-analysis.md first-24h
  checklist: verify hostile_ratio non-zero after traffic, confirm log tailer survived logrotate,
  check hostile_ratio baseline (>5% on quiet day = recalibrate).
  Acceptance: all 10 items in checklist checked off; baseline documented.

- [ ] **S31-02: nginx cache tuning** — After first real traffic, review cache hit rate via nginx
  stub_status or access log analysis. Tune `proxy_cache_valid` TTLs if hit rate < 80%.
  Acceptance: cache hit rate ≥ 80% documented in ops-runbook.md. (S24-04)

- [ ] **S31-03: Load test baseline** — `wrk -t4 -c50 -d30s https://iduna.farthq.com/` before
  real traffic. Record req/s, p99 latency, cache hit rate. Document in ops-runbook.md.
  Acceptance: baseline documented; ceiling known before real traffic hits. (S24-05)

- [ ] **S31-04: Verify logrotate + DIS tailer** — After midnight, confirm `systemctl status edis-dis`
  shows active and `/dis/health` returns valid JSON. logrotate must NOT use copytruncate
  (DIS tailer needs inode-change detection, not truncate). Check /etc/logrotate.d/nginx.
  Acceptance: DIS health check passes 24h after first deploy.

---

## SECTION 32: MJOLNIR MILESTONE 5 (FCM live device)

*Milestones 0–4 complete. M5 is the end-to-end FCM live test on Emily's real Android device.*
*Infrastructure exists: FCM sender in emily-agent, push registration endpoint in IDUNA.*
*Gap: no real device_token registered in IDUNA. No live push verified.*

- [ ] **S32-01: Register device_token in IDUNA** — On Emily's Android device, open MJOLNIR,
  navigate to Settings, copy the FCM registration token. POST to
  `IDUNA /api/v1/push/register` with token + device_name=emily-phone.
  Acceptance: IDUNA push_tokens table has one entry.

- [ ] **S32-02: Send test FCM push from emily-agent** — `curl -X POST http://localhost:8086/api/v1/emily/push/test`
  or fire a manual RSI cycle completion to trigger the FCM send path.
  Acceptance: Emily's phone receives push notification from production emily-agent.

- [ ] **S32-03: Verify MJOLNIR RsiScreen live** — Open MJOLNIR, verify RsiScreen fetches from
  production emily-agent at EMILY_BASE_URL (production IP/domain). TokenSparklineCard shows
  real token spend from production Apples.
  Acceptance: MJOLNIR displays live data from production system, not localhost.

- [ ] **S32-04: NORTHSTAR + BACKLOG: mark M5 complete** — Update MJOLNIR/docs/NORTHSTAR.md
  Milestone 5 status. Post Apple. File in world-state.md.

---

## SECTION 33: TYLER LORE BACKFILL (S08E11–S09E01)

*Emily Prime flagged this as priority 1 in world-state NEXT PRIORITIES (2026-06-15).*
*Episodes S08E11–S09E01 (Builds 0093–0098) exist in episode files but lore files not backfilled.*
*The lore file series is the canonical Tyler archive — episode-only content is invisible to Emily.*

- [x] **S33-01: TYLER-077 through TYLER-082** — Backfill lore entries from episode content
  (S08E11–S09E01). Format: existing TYLER-NNN.md convention. Each entry: title, city, date,
  duration, construction type, key dialogue excerpt, camera op note.
  ✓ Apple #564 — 2026-06-16. All entries present in eastwind_archive.md (built during S8/S9 work).

- [x] **S33-02: Camera Op Entries 70–75** — Written from episode content. Camera Op entries are
  the first-person field notes; they carry the phenomenological texture of each site.
  ✓ Apple #564 — 2026-06-16. Entries 70–75 present in camera_op_sealed_log.md.

- [x] **S33-03: Memos #057–#062** — Strategic memos from the mechanism POV. Covers the
  Córdoba→Toledo→Prague→Genoa arc (Season 8 finale + Season 9 opener).
  ✓ Apple #564 — 2026-06-16. Memos 057–062 present in jiangshi_project_memos.md.

- [x] **S33-04: EPISODES.md sync** — Verify EPISODES.md reflects all 72 episodes correctly.
  Season 9 sections complete through S09E05. Season 10 placeholder.
  ✓ Apple #564 — 2026-06-16. EPISODES.md verified complete, 72 episodes S01–S09E05.

---

## SECTION 34: RSI TOKEN EFFICIENCY (obs-watcher dedup)

*From rsi-token-report observation (2026-06-13): two pending optimizations left unimplemented.*
*Together estimated to save 100K–200K tokens/avoided session per day.*

- [x] **S34-01: obs-watcher dedup window 4h→8h** — In PRRJECT_FATBABY/cmd/observation-watcher/main.go,
  increase the dedup window from 4 hours to 8 hours. Prevents near-duplicate observations
  from firing redundant claude sessions within the same trading day.
  Estimated savings: 100K–200K tokens/avoided session per day. (Rank 3 from rsi-token-report)
  Acceptance: go test ./... passes; dedup window documented in CLAUDE.md.
  ✓ Apple #561 — 2026-06-16. Commit 123566e. go test passes.

- [x] **S34-02: runReportFooter compression** — Compress the mandatory report footer injected
  into every claude session. Current footer is verbose; can be tightened by ~30% without
  losing required steps. (Rank 4, low priority)
  ✓ Apple #566 — 2026-06-16. 63 lines → 19 lines (~37% reduction). Commit bf89779. go test passes.

---

## SECTION 35: EDIS DIS PRODUCTION HARDENING

*Pre-deploy bugs fixed 2026-06-16 (Apple #591): nginx parser zero-records, posture window race, missing hostile_ratio.*
*3 items remain before DIS is fully production-hardened post-launch.*

- [x] **S35-01: ForceState admin endpoint + manual override button** — POST /dis/force?state=<state>
  with Bearer auth (--admin-token flag). Admin panel Force Posture button + token field in edis-dis.php.
  Apple #871 | commit c5f9982 | 2026-06-16

- [x] **S35-02: Fix EDIS_DIS_COLLECTOR_URL constant-at-boot timing** — Replaced constant with
  lazy edis_dis_collector_url() function. Option changes take effect immediately.
  Apple #871 | commit c5f9982 | 2026-06-16

- [x] **S35-03: Per-IP session map for inter-request delta scoring** — The DIS fingerprinter has a
  `DeltaMs` scoring signal (+30 for delta < 20ms) but it's never populated from the log tailer —
  only the Harvester middleware (live request path) has timing context. A small bounded per-IP
  map in the collector (e.g., last-seen timestamp per /24 prefix) would fill this gap.
  Second-pass feature; not a launch blocker. Adds ~+30 hostile score on scanner bursts.
  Acceptance: go test covers delta scoring path; hostile_ratio rises under simulated burst.
  — Apple #874. EDIS commit (pending).

---

## SECTION 36: SIGNAL PIPELINE DATA QUALITY (audit 2026-06-17)

*Source: Claude Code data audit — 2026-06-17. 10,748 signal_failed today, 87% failure rate overall.*
*Priority order: S36-01 (429 blocks all new processing) → S36-02 (pre-2000 junk re-runs every restart) → S36-03 (doc limit) → S36-04 (stub provider) → S36-05 (empty ticker).*

- [x] **S36-01: EDGAR 429 rate-limit retry + throttle in processor** — `internal/processor/fetch_clean.go` [done 2026-06-17]
  has no retry logic. On 429, the filing is permanently written as `signal_failed` and never retried.
  9,435 failures today alone. Fix: detect 429 in `FetchAndCleanText`, backoff 60s/120s/300s before
  recording permanent failure. Also add a per-host rate limiter (≤10 req/s) across the 4 workers so
  EDGAR is never hammered simultaneously.
  Acceptance: `go test ./...` passes; processor runs 10 min without 429 failures in log.

- [x] **S36-02: Skip pre-2000 EDGAR filings with empty primary_document URL** — 977 failures today [done 2026-06-17]
  are all accession years 1994–2000 (e.g. `0000320193-94-000013` for AAPL). Their `primary_document`
  field is empty (`""`), producing `fetch primary document: Get "": unsupported protocol scheme ""`.
  These re-fail on every processor restart, accumulating junk `signal_failed` events indefinitely.
  Fix: in the processor worker, skip (log once, write `signal_failed` once) any filing where
  `primary_document` is empty or the URL scheme is not http/https. Do not retry.
  Acceptance: no new `unsupported protocol scheme` entries after restart.

- [x] **S36-03: Raise processor doc limit 4 MB → 16 MB** — 289 failures today (BEN, BLK confirmed). [done 2026-06-17]
  `cmd/processor/main.go` default flag `max-doc-bytes = 4<<20`. Proxy statements and large 8-Ks
  commonly exceed 4 MB. The prior secwatch fix raised body limit 64 MB → 256 MB (Apple #341).
  Apply same pattern here: raise default to `16<<20`. No other code changes needed.
  Acceptance: BEN/BLK proxy filings process without `document too large` error.

- [x] **S36-04: Wire real intelligence provider into cmd/processor** — All 5,316 `signal_generated`
  events in the store use `stubProvider`, returning `"Stub analysis result."` / `signal_type=Other`
  for every filing. The LLM analysis layer has never run. Wire `ANTHROPIC_API_KEY` + claude-haiku
  into the processor as the real `intelligence.Provider`. Use haiku (cheap) for classification.
  Dependency: API credit balance ✓ (top up console.anthropic.com).
  Acceptance: new `signal_generated` events have real `signal_type`, non-stub `summary`.
  ✓ Apple #3075

- [x] **S36-05: Fix empty-ticker guidance signal** — `cmd/guidance-watcher/main.go:112`: [done 2026-06-17]
  `ticker := tickerByID[ev.PRDiscoveryID]` — when the PR discovery ID is not in the ticker map,
  `ticker` is `""` and a signal is published with no ticker. One confirmed today:
  `guidance published ticker=  action=raised metric=eps confidence=0.90`.
  Fix: skip `if ticker == ""` before publishing. Also add guard in dividend-watcher and
  buyback-watcher where same empty-ticker events appear (`ticker="" event=cut amount=2.03`).
  Acceptance: no `ticker=""` guidance/dividend/buyback signals emitted.

---

## SECTION 37: NEWSSITE FUNCTIONALITY (web audit 2026-06-17)

*Source: emily web_audit_url run 2026-06-17. Newssite serving at :8082.*
*Critical: /ask 500 breaks product landing page. HEAD 405 breaks SEO crawlers.*

- [x] **S37-01: Fix /ask 500 — AskLandingData missing Symbols field** — `GET /ask` returns [Apple #1229 — 2026-06-18]
  HTTP 500 "template error" because `internal/newssite/render.go:RenderAskLanding()` passes
  `AskLandingData{GoogleClientID: ""}` to the template, but the shared `masthead` template
  uses `{{range .Symbols}}` for the ticker datalist. `AskLandingData` has no `Symbols` field,
  so template execution panics internally. Fix: add `Symbols []string` to `AskLandingData`
  and populate it from the handler's symbol list (same source as `TickerPageData.Symbols`).
  File: `internal/newssite/render.go` (AskLandingData struct) + `internal/newssite/asklily.go`
  (serveAskLanding populates it). Acceptance: `GET /ask` returns 200 with ticker datalist.

- [x] **S37-02: Add HEAD method support across all newssite routes** — All routes return [Apple #1229 — 2026-06-18]
  HTTP 405 on HEAD requests (web audit confirmed: /ticker/AAPL, /section/*, /doc/*, /person/*, etc.).
  Go's net/http ServeMux does not automatically handle HEAD for registered GET handlers.
  Fix: add a middleware wrapper in `cmd/newssite/main.go` that converts HEAD to GET and
  suppresses the body before writing the response. Standard pattern:
  `mux.Handle("/", headToGet(h))`. Acceptance: HEAD on all existing GET routes returns 200.

- [x] **S37-03: /ticker/{sym}/feed.xml — implement RSS or return 404** — `GET /ticker/JPM/feed.xml` [Apple #1234 — 2026-06-18]
  returns HTML (200 with ticker page HTML). The route falls through to the ticker handler
  because `/ticker/` catches the whole path. Fix: either (a) implement a real RSS/Atom feed
  for the ticker (signals as feed items), or (b) explicitly return 404 for `.xml` suffix.
  Acceptance: feed.xml returns valid RSS/Atom with correct Content-Type, or 404.

- [x] **S37-04: Start newssite with EMILY_BASE_URL set** — `POST /api/ask` returns 503 [Apple #1229 — 2026-06-18]
  "Ask Emily not configured" because `EMILY_BASE_URL` env var is not set at newssite launch.
  Emily Prime is at `http://localhost:8086`. Fix: add `EMILY_BASE_URL=http://localhost:8086`
  to the emily start --newssite launch command or to emily-secrets.env.
  Acceptance: `POST /api/ask` with a question returns an Emily response, not 503.

- [x] **S37-05: Seed SQLite governance_signals from entity-graph output** — signalapi serves [Apple #1236 — 2026-06-18]
  `GET /v1/governance-signals?ticker=AAPL` → `[]` (empty). The SQLite `governance_signals`
  table has 0 rows. The projector writes to MySQL (not running). The entity-graph produces
  real signals (insider, governance_health, board_decay etc.) but these go to the entity-graph
  var/ files, not to governance_signals. Fix: wire entity-graph signal output to the SQLite
  governance_signals table as a SQLite-mode projector path.
  Dependency: S36-01 (need data flowing) or manual seed from entity-graph signals.
  Acceptance: `/v1/governance-signals?ticker=AAPL` returns ≥1 signal.

---

## SECTION 38: EMILY BOT — SHANKPIT PLAYER (2026-06-18)

*Goal: Emily Prime can join and play SHANKPIT as a real network player.*
*`apps2/emily-bot` is the Go headless bot client targeting the Go server (:6969).*

- [x] **S38-01: emily-bot base client — connect, snapshot, heuristic AI** — `apps2/emily-bot/main.go`
  connects via `PacketConnect`, receives `PacketSnapshot` peer positions, applies
  seek-and-shoot heuristic (aim yaw via atan2, shoot within 15° tolerance, close in/back
  off by range), sends `PacketUserCmd` at 20 Hz with correct 49-byte wire format.
  `make emily-bot` builds to `bin/emily-bot`. GOWORK=off go test ./... passes.
  ✓ Apple #1238 — 2026-06-18.

- [x] **S38-02: emily-bot self-position tracking via dead-reckoning** — Two-level tracking:
  (1) snapshot parser extracts own entity (id==myID) as server-authoritative anchor; (2) per-tick
  dead-reckoning integrates fwd/str × 8u/s × 0.05s in yaw-space between snapshots. Also fixed:
  own entity was previously added to peers[], causing bot to aim at itself. — Apple #1377 — 2026-06-18

- [x] **S38-03: emily-bot kill/event reporting to Emily Prime** — `emily observe` called on
  session start, PacketWelcome connect, and PacketImpact hit_entity==1 (kill). Non-blocking
  goroutine; rate-limited 15s gap; -no-report flag. — Apple #1388 — 2026-06-18

- [x] **S38-04: emily-bot weapon rotation** — `weaponForRange(dist)`: Magnum default/close,
  AR >30u, Sniper >50u. WeaponIdx set each tick from target distance. — Apple #1393 — 2026-06-18

- [x] **S38-05: bot_client genome version guard** — `apps/bot_client/src/main.c:88` does
  `fread(&brain, sizeof(BotGenome), 1, f)`. BotGenome grew from 7→8 floats (version 1→2) when
  `w_retreat` was added (Apple #1268). Old saved genome files (28 bytes) will leave `w_retreat`
  uninitialized, producing undefined retreat behavior.
  Fix: after fread, check `brain.version < 2` and set `brain.w_retreat = 0.5f`.
  Acceptance: loading a version-1 genome file produces `w_retreat = 0.5` in logs. — Apple #1396 — 2026-06-18

- [x] **S38-06: accumulated_reward never resets on respawn — evolution selects oldest bot, not best**
  `phys_respawn()` (`packages/common/physics.h:2623`) calls `evolve_bot(p, get_best_bot())` but
  never zeroes `p->accumulated_reward` before or after evolution. Result: `get_best_bot()` always
  returns the longest-lived bot (most survival ticks = highest total reward), not the bot that
  performed best THIS life. A new bot with a great genome is immediately outscored by any bot
  that has been alive longer.
  Fix: zero `accumulated_reward` at the END of `phys_respawn` (after evolve_bot, so the outgoing
  winner's score is captured for this selection round), then start fresh for the next life.
  Acceptance: bots with better genomes win selection even if they respawned recently. — Apple #1398 — 2026-06-18

- [x] **S38-07: hit_feedback set on attacker, not defender — w_repel can't learn evasion**
  All `hit_feedback` assignments in `packages/common/physics.h` set `attacker->hit_feedback`
  (lines 2138, 2167, 2357). The DEFENDER's `hit_feedback` is never set from incoming damage.
  `w_repel` in `bot_think` fires on `me->hit_feedback > 0`, which means it fires when the bot
  is LANDING hits, not when it is BEING shot. The gene cannot evolve evasion under fire.
  Fix: in `phys_enter_death_state` and damage-application paths, also set
  `target->hit_feedback = max(target->hit_feedback, 15)` so the defender knows it was hit.
  Rename the defender-side field meaning: "incoming_damage_feedback" in comments.
  Acceptance: bot that is shot at strafes evasively (w_repel fires on incoming damage, not outgoing). — Apple #1398 — 2026-06-18

---

## SECTION 39: GPT-2 GAME AI — POLICY NETWORK (2026-06-18)

*Goal: Use our fine-tuned GPT-2 model as a real game AI policy — encode game state as tokens,
generate action tokens, decode to game inputs. Testbed: SHANKPIT emily-bot. MOBA northstar: BedWars.*
*Full spec: `gpt2-alpine-c/docs/GAME_AI_NORTHSTAR.md`*

- [x] **S39-01: Game state serializer + action decoder** — `scripts/game_state.py` in gpt2-alpine-c.
  Two functions: `serialize_snapshot(snapshot_bytes) → str` (PacketSnapshot → natural language token
  string like `"player pos:14,8 hp:85 enemy pos:20,15 dist:12 vis:1"`) and
  `decode_action(action_str) → dict` (parse GPT-2 output → UserCmd fields).
  The state/action format is the contract for all downstream milestones.
  Acceptance: round-trip encode→decode→re-encode is stable; Python unit tests pass. — Apple #1401 — 2026-06-18

- [x] **S39-02: Replay logger in emily-bot** — `apps2/emily-bot/main.go` logs `(state, action)` NDJSON
  per tick to `SHANKPIT/var/replays/YYYYMMDD-HHmm.ndjson`. Each line:
  `{"tick":N, "state":"...", "action":"..."}`. After session, `scripts/build_game_corpus.py`
  aggregates replay files → training JSONL for dataset builder.
  Acceptance: 100-tick bot session produces ≥100 records; corpus builder ingests them. — Apple #1403 — 2026-06-18

- [x] **S39-03: Fine-tune GPT-2 on replay corpus** — Add `--game-replays <dir>` flag to
  `prime_directive_dataset.py`. Each replay record becomes an instruction pair
  `{prompt: state_str, completion: action_str}`. Mixed corpus: 30% game / 70% Emily operational.
  Run Colab fine-tune. Acceptance: game loss < 2.0; `gpt2_run` generates valid action tokens
  when prompted with a game state string. — Apple #1405 — 2026-06-18

- [x] **S39-04: GPT-2 policy in emily-bot** — Replace heuristic `think()` with `POST :8088/generate`.
  `-gpt2-url` flag activates GPT-2 policy; heuristic is fallback when server unavailable.
  4 Hz decision loop: serialize state → generate → decode → send UserCmd.
  Acceptance: bot connects and plays a session with non-trivial action distribution. — Apple #1408 — 2026-06-18

---

## SECTION 40: DRAGONFLY VOXEL INTEGRATION — REAL CHUNKS IN SHANKPIT (2026-06-18)

*Goal: Get actual Dragonfly/GoblinFoxDragon world chunks rendering in SHANKPIT instead of
procedurally generated placeholder terrain. Pre-req for the SHANKPIT↔Dragonfly persistent
world bridge (NORTHSTAR Milestone 4).*
*Northstar ref: `SHANKPIT/docs2/NORTHSTAR.md` Milestone 4.*

- [x] **S40-01: Wire WorldBackend.SceneVoxelPayload() in Go server** — `apps2/server-go/main.go`
  currently calls `scanChunkForVoxelBlocks()` directly in `sendVoxelPacket()`, bypassing the
  `WorldBackend` interface seam (`server/system/backend.go`). Refactor so the voxel dispatch
  path calls `backend.SceneVoxelPayload(sceneID, chunkX, chunkZ)` and the procedural generator
  lives in `StaticBackend.SceneVoxelPayload()`.
  Acceptance: `go test ./...` passes; voxel behavior unchanged with `StaticBackend`; swapping
  backends is now a one-line change in `main.go`. — Apple #1410 — 2026-06-18

- [x] **S40-02: DragonflyBackend.SceneVoxelPayload() — real Dragonfly chunks** — In
  `GoblinFoxDragon`, implement `SceneVoxelPayload(sceneID, chunkX, chunkZ)` by reading the
  Dragonfly world's chunk at `(chunkX, chunkZ)` and serializing non-air blocks into the
  `PacketVoxelData` format (block IDs mapped to SHANKPIT block ID constants).
  Depends on S40-01 (interface wired in server).
  Acceptance: `emily start --shankpit --dragonfly` shows Dragonfly world terrain in the
  SHANKPIT client for SCENE_VOXWORLD; stone/grass/trees from Dragonfly chunk data, not
  procedural generator. — Apple #1412/#1413 (SHANKPIT DragonflyBackend + GFD worldapi) — 2026-06-18

- [x] **S40-03: Block ID mapping table** — Define a canonical mapping between Dragonfly/Bedrock
  block IDs and SHANKPIT block ID constants (stone=1, grass=2, dirt=3, log=17, leaf=18).
  Dragonfly uses numeric runtime IDs per game version. Create `packages2/common/block_map.go`
  (Go) and `packages/common/block_map.h` (C) with the mapping.
  Acceptance: all block IDs round-trip through map without falling back to `log` color. — Apple #1417 — 2026-06-18

---

## SECTION 41: DRAGONFLY BACKEND ACTIVATION (2026-06-18)

*Follow-on from S40. S40 wired the interface seam and created DragonflyBackend + worldapi.*
*S41 activates the full bridge and wires the GFD ChunkGenerator to real Dragonfly chunk data.*

- [x] **S41-01: --dragonfly-url flag in Go server** — `apps2/server-go/main.go` hardcodes
  `&system.StaticBackend{}`. Add `-dragonfly-url` flag: when set, activates `DragonflyBackend{APIURL}`.
  Enables runtime switching without recompilation.
  Acceptance: `./server-go --dragonfly-url http://localhost:7070` logs "WorldBackend: DragonflyBackend" + go test passes. — Apple #1419 — 2026-06-18

- [x] **S41-02: GFD DragonFly ChunkGenerator — real world chunk reads** — `server/worldapi/worldapi.go`
  defines `ChunkGenerator` interface but no real Dragonfly implementation exists. Write
  `DragonFlyChunkGenerator` in `server/worldapi/dragonfly_gen.go` that reads from a Dragonfly
  `world.World` and converts blocks using `DragonflyNameToBlockID()` from packages2/common/block_map.go.
  Acceptance: a running Dragonfly server's chunk data is served via GET /chunks and rendered in SHANKPIT. — Apple #1421 — 2026-06-18

---

## SECTION 42: NEXT SPRINT PLAN (2026-06-18)

*State of play: all blocked items await human action (Colab GPU, Steam account, production server sudo, Android device, API credits).*
*All implementable code gaps from S38–S41 are closed. Next implementable work is in two tracks.*

### Track A — GFD World Server activation (unlocks S40 fully)

- [x] **S42-01: Start worldapi HTTP server in GFD apps2/server-go** — worldapi goroutine on :7070.
  --worldapi-port flag. DragonflyChunkGenerator(ProceduralWorldStore) wired in.
  — Apple #1869. 2026-06-20.

- [x] **S42-02: GFD world generation — ProceduralWorldStore scene-differentiated terrain.**
  scene0=flat meadow, scene1=rolling hills (sin height 2-8), scene2=stone caves (EW+NS corridors).
  13 tests pass. — Apple #1869. 2026-06-20.

### Track B — GPT-2 self-play corpus accumulation (parallel with S26-04 Colab wait)

- [x] **S42-03: emily-bot self-play mode — -sessions N + -session-duration.** runSession() loop;
  each session gets fresh UDP conn, state.done closed on disconnect, deadline per session.
  — Apple #1873. 2026-06-20.

- [x] **S42-04: emily-bot health-aware retreat at hp<30.** myHealth float32; decrement near enemy,
  regenerate when clear, retreat (fwd=-1, no attack, strafe) below 30. hp in serializeState.
  — Apple #1873. 2026-06-20.

### Track C — Intelligence provider (unblocks S36-04 if API key available)

- [x] **S42-05: HaikuProvider with ENABLE_LLM_ANALYSIS=true feature flag.**
  internal/processor/haiku_provider.go; cmd/processor/main.go gates on env var.
  4 mock-HTTP tests pass. — Apple #1876. 2026-06-20.

---

## SECTION 43: SHANKPIT PLAYER AUTHENTICATION (IDUNA + OAuth) (2026-06-18)

*Goal: SHANKPIT players have real identities — Google login via IDUNA, persistent player profiles,*
*ELO scores, kill history. Later: email/password as fallback.*
*This is the prerequisite for leaderboards, Steam Early Access accounts, and anti-cheat.*
*Architecture: IDUNA issues ES256 JWTs; game server validates on connect; client presents token in PacketConnect.*

- [x] **S43-01: PacketConnect JWT field** — ES256 JWT parsed from bytes [13..268] of PacketConnect.
  JWKSCache fetches IDUNA /.well-known/jwks.json (10min refresh). Valid → playerID+displayName
  in clientInfo. Invalid/absent → guest-NNN. admin /admin/players now returns player_id+display_name.
  — Apple #1887. 2026-06-20.

- [x] **S43-02: IDUNA player registration endpoint** — migration 202606200001_players.sql.
  POST /api/v1/players/register upserts on (provider, provider_sub), returns player_id+display_name.
  GET /api/v1/players/{id} returns public profile (kills, deaths, kd_ratio, sessions). JWT-gated.
  — Apple #1889. 2026-06-20.

- [x] **S43-03: Google OAuth flow for SHANKPIT** — ShankpitAuthHandler: GET /api/v1/auth/google/shankpit
  → Google consent (CSRF state cookie). /callback: code exchange → userinfo → upsertPlayer →
  IDUNA JWT (sub=player_id, aud=shankpit, 72h) → redirect shankpit://auth?token=...&player_id=...
  GOOGLE_CLIENT_SECRET + SHANKPIT_OAUTH_REDIRECT_URI env vars. — Apple #1891. 2026-06-20.

- [x] **S43-04: Player profile in IDUNA** — Implemented in S43-02: GET /api/v1/players/{id}
  returns kills, deaths, kd_ratio, sessions, registered_at, last_seen. — Apple #1889. 2026-06-20.

- [x] **S43-05: Email + password login** — PlayerEmailAuthHandler: POST /api/v1/auth/email/register
  (create player + bcrypt credential) + POST /api/v1/auth/email/login (verify + return JWT).
  Migration 202606200002_player_credentials.sql. Same 72h JWT format (aud=shankpit).
  — Apple #1894. 2026-06-20.

---

## HUMAN UNBLOCK QUEUE (items waiting on you)

| Item | What's needed | Unblocks |
|---|---|---|
| S26-04 | Google Colab T4 session + `emily train upload` | GPT-2 fine-tune; S26-05; game AI policy quality |
| S32-01 | Open MJOLNIR on Emily's phone, copy FCM token | FCM live push (S32-02,03,04) |
| S36-04 / S42-05 | Top up console.anthropic.com API credits | LLM signal analysis pipeline |
| S19-03 | Create Steamworks account ($100) | Steam EA launch (S19-05) |
| S23-01 | `sudo bash` on production host | EDIS WordPress live |
| S30-02 | `sudo mysql` on production | MySQL projections in prod |
| S2 (MPT) | Pexels API key | Video compilation pipeline |
| S45-02 run | `sudo apt install libsdl2-dev` + create emilyspringerton/PITVIPER on GitHub | PITVIPER window launches |

---

---

## SECTION 44: PYTHON SYSTEM SDK (Colab + distributed Emily federation)

**Goal:** A Python SDK (`einhorn_sdk/`) that wraps the full IDUNA API (local auth + user CRUD +
agents + apples + drive + HEIMDAL + subscriptions) so any Colab notebook or remote Emily worker
can authenticate and interact with the system without writing raw HTTP. Auto-regenerates from the
IDUNA Swagger spec.

**Distributed Emily vision:** Multiple independent Emily clusters (Colab GPU worker, production
monolith on IDUNA, external compute) self-identify via agent credentials, push work to the same
repos via git pull-rebase, and pull from a shared log stream. Proxies present a seamless API to
institutions.

- [x] **S44-00** — IDUNA webmaster uid=0 + user CRUD + event log + SQLite/MySQL projectors.
  POST /api/v1/auth/local, FileEventLog (FatBaby NDJSON pattern), bcrypt. All tests green.
  — Apple #1445. 2026-06-18.

- [x] **S44-01** — IDUNA Swagger spec at GET /api/v1/openapi.json. OpenAPI 3.1 JSON served by IDUNA.
  Covers: auth (local+agent), users (CRUD), apples, drive, heimdal, health, jwks.
  — Apple #1448. 2026-06-18.

- [x] **S44-02** — Python SDK `einhorn_sdk/` at `IDUNA/sdk/python/`. IdunaClient + sub-clients
  (UsersClient, ApplesClient, DriveClient, HeimdalClient, AgentsClient). 7 tests pass.
  — Apple #1448. 2026-06-18.

- [x] **S44-03** — Colab observability notebook `sdk/python/notebooks/emily_observability.ipynb`.
  Training loop hooks, Apple tail (polling), HEIMDAL submission, Drive upload, git pull-rebase.
  ColabEmily.from_colab_secrets() reads Colab Secrets panel or env vars.
  — Apple #1448. 2026-06-18.

- [x] **S44-04** — IDUNA GET /api/v1/agents endpoint. ?type=emily_cluster filter. AgentsHandler
  in handlers/agents.go. Distributed Emily clusters visible in /admin/ and queryable by Emily Prime.
  — Apple #1878. 2026-06-20.

- [x] **S44-05** — git pull --rebase before push. emily-agent ConversationStore.Save() + obs-watcher
  runReportFooter(). Rebase conflict → escalation Apple + skip push. No force push.
  — Apple #1881. 2026-06-20.

- [x] **S44-06** — IDUNA GET /api/v1/stream/user-events SSE endpoint. UserEventStreamHandler
  polls FileEventLog every 2s, streams records as text/event-stream. JWT-gated. ?from_seq +
  ?timeout params. OpenAPI 3.1 spec updated. Colab notebooks subscribe for real-time user events.
  — Apple #1884. 2026-06-20.

---

## SECTION 45: SHANKPIT STAT REPORTING + PITVIPER MILESTONE 1

*S45-01: SHANKPIT server reports player kill/death stats to IDUNA on session end.*
*S45-02 onward: PITVIPER Milestone 1 — SDL2 terminal bootstrap (PTY, glyph render, first keypress).*

- [x] **S45-00: emily-bot JWT auth in PacketConnect** — sendConnect() packs SHANKPIT_AUTH_TOKEN
  env / ~/.shankpit/auth.json into bytes [13..268]. Server (S43-01) validates. Emily Prime now
  connects as an authenticated IDUNA player. — Apple #2305. 2026-06-21.

- [x] **S45-00b: POST /api/v1/emily/push/test** — handlePushTest in api_push.go.
  Gets mjolnir-emily device token from IDUNA, fires test FCM push.
  503 (FCM/IDUNA not configured), 404 (no token / S32-01 needed), 502 (FCM error), 200 on success.
  — Apple #2307. 2026-06-21.

- [x] **S45-01: SHANKPIT → IDUNA player stat updates** — clientInfo.lastPacket + kills fields.
  pruneIdleClients goroutine (15s check, default 60s timeout) removes idle clients and
  POSTs {"kills": N} to IDUNA /api/v1/players/{id}/session via reportPlayerSession.
  IDUNA: handleSessionEnd does UPDATE players SET sessions+1, kills+N, deaths+N, last_seen.
  New server flags: --server-token, --idle-timeout. — Apple #2312. 2026-06-21.

- [x] **S45-02: PITVIPER Milestone 1 scaffold** — go module pitviper created. internal/vterm:
  VT100 state machine (SGR, cursor, scroll, erase — 8 tests pass). internal/pty: openpty +
  TIOCSWINSZ. internal/font: 8×13 glyph atlas via x/image/font/basicfont. cmd/pitviper:
  SDL2 window + PTY + 60fps render loop + keyboard forwarding + resize. Blocked to run on:
  sudo apt install libsdl2-dev + create emilyspringerton/PITVIPER repo on GitHub. — Apple #2315. 2026-06-21.

- [x] **S45-03: PITVIPER glyph cache + UTF-8 decode** — font.Atlas pre-rendered at init() for all
  printable ASCII (0x20–0x7E) via x/image/font/basicfont. vterm.Write() accumulates multi-byte UTF-8
  sequences; invalid bytes → U+FFFD. TestUTF8Rune: "café" decodes, cursor at col 4. 14 tests pass.
  — Apple #2318. 2026-06-21.

- [x] **S45-04: PITVIPER color + SGR escape sequence parser** — handleCSI covers: CHA (G), VPA (d),
  A/B/C/D cursor movement, H/f position, J/K erase, m SGR (8-color, bright-8, 256-color 38;5;N,
  bold 1/22, reset 0), s/u save/restore. TestEraseInLine, TestSGR256Color, TestSGRBold all pass.
  14 vterm tests total pass. — Apple #2320. 2026-06-21.

- [x] **S45-05: PITVIPER scrollback buffer** — Ring buffer (10,000 lines max). scrollUp() saves
  rows before clearing. Snapshot() composites scrollback+live when scrollTop>0; cursor=-1,-1
  when scrolled. ScrollBy/ScrollReset/ScrollLines/ScrollbackLen. SDL2 Shift+PageUp/Down (half-page),
  Shift+Home/End. 17 vterm tests pass (TestScrollbackView, TestScrollByClamp). — Apple #2330. 2026-06-21.

---

## SECTION 46: EMILYOS MILESTONE 2+3 TEST COVERAGE + CMD

*Goal: Complete EmilyOS Milestone 2 acceptance criteria (RBAC + posture tests) and add the CLI.*

- [x] **S46-01: EmilyOS RBAC + posture test coverage** — rbac_test.go: 8 tests (Operator/Admin/Auditor
  caps, CapsForRole copy, ValidRole, AllRoles). machine_test.go: 7 tests (default, NormalToSiege,
  invalid SIEGE→INCIDENT, persist, SIEGE denies cap.net, NORMAL passthrough, INCIDENT forces export).
  All 4 EmilyOS packages green. — Apple #2322. 2026-06-21.

- [x] **S46-02: EmilyOS `cmd/emilyos` CLI** — posture get/set (admin-gated, audit event, persist),
  verb dispatch (capForVerb table, policy deny → exit 2), audit tail -n N (JSON events),
  audit verify (chain integrity → exit 3 on tamper). EMILY_{ACTOR,SESSION,DEVICE,ROLE,POSTURE,AUDIT}
  env vars. Smoke tested: NORMAL→SIEGE, auditor denied EXEC, chain verify ok. — Apple #2325. 2026-06-21.

- [x] **S46-03: EmilyOS posture gate in emily-agent** — posture.go: readPosture() reads
  EMILY_POSTURE_PATH/posture.json. postureBlocksLLM(): SIEGE=skip cycle, EXITED=block.
  cron.go RunOnce() checks posture first. main() asserts NOT EXITED on startup.
  Build ok, tests pass. — Apple #2327. 2026-06-21.

---

## SECTION 47: PITVIPER MILESTONE 1 HARDENING + EMILYOS SOC2 EXPORT

*S47-01: PITVIPER OSC escape sequences (window title from bash/vim).*
*S47-02: PITVIPER TextInput key echo snap-to-live.*
*S47-03: EmilyOS `audit export` — SOC 2 evidence bundle (zip + chain manifest).*
*S47-04: emily.cli `emily emilyos` sub-command — wraps emilyos CLI.*

- [x] **S47-01: PITVIPER OSC title sequences** — Parse `\033]0;title\007`, `\033]1;...`, `\033]2;title\007`
  in vterm.Write(). Store in Screen.Title string. SDL2 window title updated on each render tick
  via `win.SetTitle()`. Acceptance: bash `\033]2;myshell\007` updates SDL2 window title bar.
  [done 2026-06-21] Apple #2333. PITVIPER commit 5237e0c.

- [x] **S47-02: PITVIPER snap-to-live on input** — Any TextInput event (character typed) while
  scrolled back should snap the view to live (ScrollReset). Currently only non-scroll key events
  snap back. Fix handleScrollKey to return false for TextInput so the TextInput handler calls
  ScrollReset before writing to PTY.
  [done 2026-06-21] Apple #2333. PITVIPER commit 5237e0c.

- [x] **S47-03: EmilyOS `audit export`** — `emilyos audit export <outdir>` writes:
  - `audit.jsonl` (full log copy)
  - `manifest.json` (event count, first/last seq, chain verified bool, export_ts)
  Acceptance: `emilyos audit export /tmp/soc2` creates both files; manifest.chain_ok=true.
  [done 2026-06-21] Apple #2335. EmilyOS commit a41f0d7.

- [x] **S47-04: emily.cli `emily emilyos` command** — Wraps EmilyOS CLI:
  `emily emilyos posture get/set`, `emily emilyos audit tail/verify`.
  Reads EMILY_POSTURE_PATH / EMILY_AUDIT_PATH from config or env.
  Acceptance: `emily emilyos posture get` prints current posture.
  [done 2026-06-21] Apple #2337. emily.cli commit a8bbd0a.

---

## SECTION 48: IDUNA PLAYERS LIST + SHANKPIT OPS

*S48-01: IDUNA GET /api/v1/players — list all players with stats.*
*S48-02: emily.cli `emily shankpit players-leaderboard` command — top kills via IDUNA players list.*
*S48-03: EmilyOS `emilyos posture history N` — last N posture transitions from audit log.*

- [x] **S48-01: IDUNA GET /api/v1/players** — `GET /api/v1/players?limit=N&sort=kills|deaths|sessions`
  returns `[{player_id, display_name, kills, deaths, sessions, last_seen}]` ordered by sort key.
  Uses existing `players` table (already has all these columns from S28 migration).
  Acceptance: `curl .../api/v1/players?sort=kills&limit=10` returns top 10 players by kills.
  [done 2026-06-21] Apple #2340. IDUNA commit c09d5b0. 3 tests.

- [x] **S48-02: emily.cli `emily shankpit leaderboard`** — Calls IDUNA `GET /api/v1/players?sort=kills&limit=10`
  and renders a formatted leaderboard table. Uses the existing IDUNA client.
  Acceptance: `emily shankpit leaderboard` prints rank, player_id, kills, deaths, K/D ratio.
  [done 2026-06-21] Apple #2342. emily.cli commit bd6d9ae.

- [x] **S48-03: EmilyOS `emilyos posture history [N]`** — Reads audit log, filters events where
  `verb=posture.set`, prints `seq | ts | actor | old_state → new_state`. Default N=20.
  Acceptance: `emilyos audit history 5` shows last 5 posture transitions.
  [done 2026-06-21] Apple #2344. EmilyOS commit c37ee46.

---

## SECTION 49: EMILY-AGENT POSTURE VISIBILITY

*S49-01: emily-agent GET /api/v1/emily/posture — posture endpoint.*
*S49-02: emily status shows EmilyOS posture from emily-agent.*

- [x] **S49-01: emily-agent GET /api/v1/emily/posture** — Returns `{posture, llm_blocked, queried_at}`.
  Routes to `api_posture.go:handlePostureGet()`. Reads posture.json via `readPosture()`.
  llm_blocked=true when SIEGE or EXITED.
  [done 2026-06-21] Apple #2346. EMILY commit 926e10f.

- [x] **S49-02: emily status shows EmilyOS posture** — Calls `/api/v1/emily/posture` from emily-agent.
  Prints `POSTURE: NORMAL` in human-readable status output. Best-effort; skipped if agent not running.
  EMILY_AGENT_URL env override. emily.cli commit dfac697.
  [done 2026-06-21] Apple #2348. emily.cli commit dfac697.

---

## SECTION 50: EMILYOS MILESTONE 4+5 — POLICY SNAPSHOTS + EVIDENCE BUNDLE

*S50-01: EmilyOS policy snapshot store tests + emilyos about + emilyos snapshot commands.*
*S50-02: EmilyOS emilyos audit bundle — tar.gz evidence bundle for SOC 2 (Milestone 5).*

- [x] **S50-01: EmilyOS policy snapshots (Milestone 4 partial)** — `internal/policy/snapshot_test.go`:
  6 tests (Write, HashStable, Latest, GetNotFound, NewIDFormat, DirCreated).
  `emilyos about`: shows version + commit + date + posture + latest snapshot.
  `emilyos snapshot capture|list|show`: manages policy snapshots in var/snapshots/.
  Build attestation via -ldflags buildCommit/buildDate.
  NORTHSTAR: milestones 1-3 marked complete; 4-5 in progress.
  [done 2026-06-21] Apple #2352. EmilyOS commit 47314dd.

- [x] **S74-01: EMILY emilyListDir test suite — 4 tests** — flat sorted/hidden-excluded, dirs-first,
  maxDepth control, non-existent error. 72 emily-agent tests. [done 2026-06-21] Apple #2468.

- [x] **S73-01: EMILY Emiree FractalFingerprint + Summary tests — 3 tests** — dimensions (20×5),
  zero-size defaults (40×14), Summary (gear/h/p/steps). 68 emily-agent tests.
  [done 2026-06-21] Apple #2465.

- [x] **S72-01: EMILY HEIMDAL formatCriteria + min tests — 4 tests** — empty/single/multi criteria,
  min (4 cases). 65 emily-agent tests. [done 2026-06-21] Apple #2418.

- [x] **S71-01: EMILY briefing buildBriefingMessage + typeLabel — 5 tests** — empty, count grouping,
  90-char title cap, top-3 notable filtering, 7 typeLabel cases. 61 emily-agent tests.
  [done 2026-06-21] Apple #2416.

- [x] **S70-01: EMILY DocStore test suite — 3 tests** — Append+LoadSeen IDs, seq increment, empty
  LoadSeen on fresh store. 56 emily-agent tests. [done 2026-06-21] Apple #2414.

- [x] **S69-01: EMILY goldenbuild contentHash + loadGoldenIndex tests — 7 tests** — contentHash
  (deterministic, different, budget truncate, budget=0), loadGoldenIndex (fallback, table parse,
  empty-tier1). 53 emily-agent tests total. [done 2026-06-21] Apple #2411.

- [x] **S68-01: EMILY Emiree witch engine test suite — 10 tests** — TestNewWitchStateDefaults,
  TestGearString (7), TestSelectGear (7 cases), Observe, Update stability (200 iter), Influence,
  ObserveRSIOutcome. 46 emily-agent tests total. [done 2026-06-21] Apple #2409.

- [x] **S67-01: PITVIPER font package test suite — 5 tests** — GlyphDimensions (8×13), AtlasSize (95),
  KnownChars (A/Z/0/9 rendered), OutOfRange (emoji→?), AllPrintable (0x20–0x7E).
  [done 2026-06-21] Apple #2407.

- [x] **S66-01: IDUNA drive.Client test suite — 5 tests** — bad JSON, missing email, missing key,
  default token URI, oversized upload rejection. No network.
  [done 2026-06-21] Apple #2405.

- [x] **S65-01: EMILY FCM jwt_test.go — 3 tests** — TestBuildServiceAccountJWTStructure (RSA key, 3-part
  JWT, alg/iss/aud/1h window), InvalidPEM, UnsupportedKeyType. No network.
  [done 2026-06-21] Apple #2402.

- [x] **S64-01: PITVIPER CSI E/F/X (CNL/CPL/ECH) + TestCSICNLCPLECH** — CSI E: down N rows col=0;
  CSI F: up N rows col=0; CSI X: erase N chars in-place. 31 vterm tests.
  [done 2026-06-21] Apple #2400.

- [x] **S63-01: PITVIPER CSI L/M (insert/delete line) + TestCSIInsertDeleteLine** — CSI L: insert N
  blank lines at cursor, shift region down; CSI M: delete N lines, shift region up. 30 vterm tests.
  [done 2026-06-21] Apple #2398.

- [x] **S62-01: IDUNA auth.Subscription.IsActive() test suite — 7 cases** — nil, perpetual-active,
  future-expiry, past-expiry, cancelled, expired-status, cancelled-but-future-expiry.
  [done 2026-06-21] Apple #2396.

- [x] **S61-01: PRRJECT_FATBABY chart package test suite — 7 tests** — TestParseYahooResponse
  (NaN/0 filter), TestParseYahooResponseError, TestParseYahooResponseEmpty, TestRenderSVGUptrend
  (green), TestRenderSVGDowntrend (red), TestRenderSVGTooFewPoints, TestRenderSVGFlatLine.
  [done 2026-06-21] Apple #2393.

- [x] **S60-01: PRRJECT_FATBABY mysqlToSQLite + splitStatements test suite — 20 cases** —
  TestSplitStatements (6), TestMySQLToSQLiteRules (12 rules), TestAlterAddIndex, TestCleanTrailingCommas.
  [done 2026-06-21] Apple #2391.

- [x] **S59-01: PITVIPER CSI S/T/P/@ (scroll+char edit) + 2 tests** — CSI S: scroll-up-N; CSI T:
  scroll-down-N; CSI P: delete-N-chars (shift left); CSI @: insert-N-blanks (shift right).
  TestCSIScrollUpDown + TestCSIDeleteInsertChar. 29 vterm tests.
  [done 2026-06-21] Apple #2389.

- [x] **S58-01: PITVIPER HT tab stop handling + TestTabStop** — `'\t'` (0x09) advances cursor to next
  8-column tab stop; clamped to cols-1. 27 vterm tests.
  [done 2026-06-21] Apple #2386.

- [x] **S57-01: PRRJECT_FATBABY issuerregistry test suite — 4 tests** — TestNew (deep copy),
  TestResolveByCIK, TestResolveReturnsCopy, TestNilRegistry.
  [done 2026-06-21] Apple #2384.

- [x] **S56-02: IDUNA subscriptions handler test suite — 5 tests** — provision 200/403/400-status;
  getMe subscribed/no-sub. stubSubStore.
  [done 2026-06-21] Apple #2382.

- [x] **S56-01: IDUNA push_tokens handler test suite — 5 tests** — upsert 200/403/400, get 200/404.
  stubPushTokenStore.
  [done 2026-06-21] Apple #2379.

- [x] **S55-02: PITVIPER IsAltActive() + TestIsAltActive + TestSGRDefaultColors** — `IsAltActive()` mutex-
  protected accessor for renderers; SGR39/49 default color reset test. 26 vterm tests.
  [done 2026-06-21] Apple #2377.

- [x] **S55-01: PITVIPER TestEraseInDisplay 0J/1J** — closes erase-in-display coverage gap. 24→25 tests.
  [done 2026-06-21] Apple #2375.

- [x] **S54-02: EmilyOS EXPORT_EVIDENCE verb + bundle test — Milestone 5 complete** — `runAuditBundle`
  dispatches `EXPORT_EVIDENCE` (cap.export) before writing; export is audited. `TestBundleManifest
  Verification`. NORTHSTAR M5 [x].
  [done 2026-06-21] Apple #2373.

- [x] **S54-01: EmilyOS POLICY_ROLLBACK verb — Milestone 4 complete** — `emilyos snapshot rollback <id>`:
  hash-verify, dispatch via Dispatcher (cap.policy.write), SOC 2 audit event. `ComputeSnapshotHash`
  exported. `TestSnapshotRollback`: 3 snapshots + rollback + verify. NORTHSTAR M4 [x].
  [done 2026-06-21] Apple #2370.

- [x] **S53-02: IDUNA intelligence handler test suite — 4 tests observe/list/patch** — TestIntelObserveSuccess,
  TestIntelObserveForbidden, TestIntelObserveMissingImageData, TestIntelListAndPatch. stubIntelStore.
  [done 2026-06-21] Apple #2367.

- [x] **S53-01: IDUNA HEIMDAL handler test suite — 5 tests submit/list/patch** — TestHeimdalSubmit,
  TestHeimdalSubmitForbidden, TestHeimdalSubmitEmptyRequirement, TestHeimdalList, TestHeimdalPatch.
  [done 2026-06-21] Apple #2365.

- [x] **S52-01: PITVIPER DECSTBM scroll region + ESC M reverse index** — scrollRegTop/scrollRegBot;
  scrollUp respects region; scrollDown() for RI; TestDECSTBM + TestReverseIndex. 23 tests.
  [done 2026-06-21] Apple #2362.

- [x] **S51-01: PITVIPER cursor visibility — ESC [?25l/h** — `CursorHidden` bool; `Snapshot()` returns (-1,-1)
  when hidden. `TestCursorVisibility`. PITVIPER commit 74f3603.
  [done 2026-06-21] Apple #2357.

- [x] **S51-02: PITVIPER alternate screen buffer — ESC [?1049h/l** — Save/restore primary cells + cursor.
  `TestAlternateScreen`. 21 vterm tests total. PITVIPER commit 36012c9.
  [done 2026-06-21] Apple #2359.

- [x] **S50-02: EmilyOS audit bundle — tar.gz SOC 2 evidence bundle (Milestone 5 partial)** —
  `emilyos audit bundle [outpath]` produces `.tar.gz` with `audit.jsonl` + `policy-snapshot.json`
  + `manifest.json` (SHA-256 per file, chain_ok, event_count, build attestation).
  Default: `soc2-evidence.tar.gz`. Exit 3 on chain tamper.
  [done 2026-06-21] Apple #2354. EmilyOS commit 60cb485.

---

## SECTION 75: DRAGONSNSHIT MMO — IDUNA SCHEMA (2026-06-21)

*Context: MMO_NORTHSTAR written. IDUNA is the auth and persistence layer for all player identity,
item provenance, and world state. No MMO system can ship without this schema.*

- [x] **S75-01: IDUNA MMO schema migration** — `migrations/truestore/202606230001_mmo_schema.sql`:
  characters/character_skills/items/guilds/guild_memberships/world_events/scene_state; 4 scenes seeded. [done 2026-06-23] Apple #2935.

- [x] **S75-02: IDUNA character API** — `internal/http/handlers/mmo.go`: POST create, GET fetch,
  PATCH position; all auth-gated via RequireAuth. [done 2026-06-23] Apple #2936.

- [x] **S75-03: IDUNA item provenance API** — same file: POST craft (provenance[0]=crafted), GET full chain,
  DELETE soft (destroyed_at), POST :id/transfer (append chain). [done 2026-06-23] Apple #2936.

- [x] **S75-04: IDUNA guild API** — same file: POST found, GET (members), POST join,
  PATCH role, DELETE disband. [done 2026-06-23] Apple #2936.

- [x] **S75-05: IDUNA world events API** — same file: POST open, GET, PATCH phase+ley_integrity,
  POST :id/resolve. [done 2026-06-23] Apple #2936.

---

## SECTION 76: DRAGONSNSHIT MMO — GO GAME SERVER (2026-06-21)

*Context: GFD server-go needs to become the authoritative MMO backend. Currently a bridge/stub.
Depends on S75 IDUNA schema.*

- [x] **S76-01: GFD server-go IDUNA auth integration** — idunaauth pkg (JWKS ES256 JWT verifier);
  PacketConnect reads JWT from buf[1:n], validates, stores playerID=sub, rejects with PacketAuthReject=8.
  [done 2026-06-23] Apple #2940.

- [x] **S76-02: GFD zone backend** — `server/zone/zone.go`: Zone{ID,Name,SpawnXYZ},
  DefaultZones() (Meadow=0/Hills=1/Caves=2), Manager.Enter/Leave/Transfer, PlayersIn(),
  SameZone(), SpawnFor(). EncodeSceneChange/DecodeSceneChange (14-byte wire, zoneID uint8 +
  float32×3 XYZ). PacketSceneChange=7 added to common/protocol.go. 38 tests.
  [done 2026-06-22] Apple #2539.

- [x] **S76-03: GFD Telecrystal travel — server-authoritative** — telecrystal registry (6 crystals),
  idunaclient (GetCharacter/DeductGold/UpdatePosition/TravelTelecrystal), PacketTelecrystalUse/Ack/Err;
  IDUNA PATCH /api/v1/characters/:id/gold added. [done 2026-06-23] Apple #2943/#2944.

- [x] **S76-04: GFD item crafting endpoint** — craft.LookupRecipe, PacketCraftRequest/Result=12/13,
  idunaclient ListItems/CreateItem/DestroyItem; server-go handler validates reagents→IDUNA→craft.Attempt.
  IDUNA: GET /api/v1/characters/:id/items. [done 2026-06-23] Apple #2947/#2948.

- [x] **S76-05: GFD World Crisis phase machine** — server/worldcrisis (6-phase, LEY decay, concurrent gate),
  PacketWorldCrisisUpdate/ObjectiveComplete=14/15, tick goroutine+broadcaster, idunaclient.PatchWorldEvent.
  [done 2026-06-23] Apple #2952.

- [x] **S76-06: GFD skill XP server-side** — PacketSkillXP=16, server-go handler (cap 1.0/action, async),
  idunaclient.IncrementSkill. IDUNA: PATCH+GET /api/v1/characters/:id/skills (UPSERT, cap 110.0).
  [done 2026-06-23] Apple #2954/#2955.

---

## SECTION 77: DRAGONSNSHIT MMO — MMO_NORTHSTAR REGISTRATION (2026-06-21)

- [x] **S77-01: DragonsNShit MMO_NORTHSTAR written** — `GoblinFoxDragon/docs2/MMO_NORTHSTAR.md`;
  7 core MMO systems documented (player identity, item provenance, guilds, economy, scene/travel,
  World Crisis, skills); IDUNA schema requirements listed; 8-milestone roadmap.
  [done 2026-06-21] Apple #2470.

- [x] **S77-02: MMO_NORTHSTAR golden-index registration** — DNS-MMO-NORTHSTAR Tier 1 budget 4000
  appended to `EMILY/context/golden-docs-index.md`. [done 2026-06-21].

- [x] **S78-01: DragonsNShit chat system — say/tell/yell/guild, PacketChat=6** —
  `server/chat/chat.go`: Router.Deliver(), 4 channels (say=50u radius, tell=private,
  yell=zone-wide, guild=cross-scene). ParseChatPacket/encodeChat. Wired into server-go main.
  19 tests. [done 2026-06-21] Apple #2513.

- [x] **S79-01: DragonsNShit linkshell guild system — Feather + Feather Sack** —
  `server/guild/guild.go`: FFXI-style linkshell. Feather=member token, Feather Sack=admin.
  Roles: Member < Officer < Leader. Officer forges Feathers; Leader forges Sacks.
  RevokeItem removes from roster. CanChat() gates ChatGuild on active item.
  22 tests. [done 2026-06-21] Apple #2516.

- [x] **S80-01: DragonsNShit auto-attack + mob tagging** —
  `server/mob/mob.go`: Mob AI state machine (Idle/Pursuing/Returning/Dead), Registry.Hit()
  tag-on-first-hit, leash+full-heal reset, mob swing timer, player PlayerCombat auto-attack
  with range/scene/timer gate. EvtMobAggro/Attack/Reset/Died events.
  26 tests. [done 2026-06-21] Apple #2519.

---

## SECTION 81: DRAGONSNSHIT — FFXI PARITY NORTHSTAR (2026-06-21)

*Context: GFD engine targets FFXI-level MMO feature parity as the product northstar.
FFXI's core systems define the implementation roadmap. Each section below is one FFXI system.*

### FOUNDATION (already shipped)
- chat (say/tell/yell/guild) ✓ S78
- linkshell guild (Feather/Sack) ✓ S79
- mob auto-attack + tagging ✓ S80
- MMO_NORTHSTAR (IDUNA schema, item provenance, World Crisis) ✓ S77

### COMBAT SYSTEMS

- [x] **S81-01: Skill chain system** — `server/skillchain/skillchain.go`: 14 Resonances
  (8 T1 + 4 T2 + Light/Darkness T3). combinationTable (23 pairs). Chain() → highest-tier
  match across ws1×ws2 attribute product. Multipliers: T1=0.20, T2=0.35, T3=0.50.
  8s window. 19 CanonicalWeaponSkills. [done 2026-06-21] Apple #2527.

- [x] **S81-02: Magic burst system** — same package. CanBurst()/Burst()/BurstElements().
  T1=1 element, T2=2 elements (union), T3=4 elements. 15s burst window. BurstMultiplier=0.35.
  End-to-end: Distortion+Ice, Light+Fire/Lightning. 31 tests total. [done 2026-06-21] Apple #2527.

- [x] **S81-03: TP weapon skill points system** — `server/combat/tp.go`: `TPState.AddTP(delay,
  hastePct)` formula (effective_delay = delay*(1-haste/100)/6, floor 5, cap 300).
  `CanWeaponSkill()`, `UseWeaponSkill()` (resets to 0), `ConsumeTP(n)`, `Reset()`,
  `HitsToWS()`. 6 delay consts (HtH=80, 1H=240, GS=480, Staff=366, Polearm=402).
  MaxHaste=80%. 34 tests. [done 2026-06-21] Apple #2533.

- [x] **S81-04: Status effects (enfeebles + buffs)** — `server/status/status.go`: 10 kinds
  (Poison/Paralyze/Slow/Silence/Bind/Haste/Regen/Refresh/Protect/Shell). Stack.Apply() with
  weaker-rejected / same-potency-refreshed / stronger-upgrades rules. Stack.Tick(now) →
  TickResult{Expired, Events} with DoT (Poison→−HP) and HoT (Regen→+HP, Refresh→+MP) events.
  Cat() for dispel targeting. DispelOne(). NetHastePct() / PhysDefBonus() / MagicDefBonus() /
  IsParalyzed() / IsSilenced() / IsBound(). Permanent effects (zero ExpiresAt) never expire.
  43 tests. [done 2026-06-21] Apple #2536.

- [x] **S81-05: Enmity (hate) system** — `server/enmity/enmity.go`: Table.Add(damage CE)/AddCure(heal CE; AoE halved)/
  Reduce/Top (highest CE, alpha tie-break); EnmityCap=10000; ErrNoPlayers. 18 tests. [done 2026-06-23] Apple #2919.

- [x] **S81-06: Death + raise system** — `server/combat/death.go`: HPState{Current,Max,IsKO}; TakeDamage→KO;
  TakeKO; Raise(xp,pct)→penalty HP=1; RaiseDefault=10%; ErrNotKO/ErrAlreadyKO. 15 tests. [done 2026-06-23] Apple #2920.

### JOB SYSTEM

- [x] **S82-01: Job definitions** — `server/job/job.go`: JobID=string const; AllJobs[22]; Stats{BaseHP,HPPerLevel,
  BaseMP,MPPerLevel,STR/DEX/VIT/AGI/INT/MND/CHR}; StatsFor/HPAtLevel/MPAtLevel; melee MP=0; ErrUnknownJob.
  14 tests. [done 2026-06-23] Apple #2921.

- [x] **S82-02: Sub-job system** — `server/job/subjob.go`: CharJob{Main,Sub,MainLvl,SubLvl}; EffectiveSubLevel=SubLvl/2;
  CombinedStats main+sub at half level; ErrSameJob. 10 tests. [done 2026-06-23] Apple #2925.

- [x] **S82-03: Job abilities + traits** — `server/job/subjob.go`: Ability{ID,Job,Recast,MinLevel}; Trait; RecastTracker.Use/Ready/
  RecastRemaining; ErrAbilityOnRecast/LevelGated/Unknown; WarriorAbilities/WhiteMageAbilities. 12 tests. [done 2026-06-23] Apple #2926.

### PROGRESSION SYSTEM

- [x] **S83-01: Level cap + XP** — `server/xp/xp.go`: XPToLevel(lvl)→100*(lvl-1)^1.8 floor 75;
  CharXP.AddXP→leveledUp bool; multi-level overflow; ErrAtCap at L99; Progress/XPToNextLevel.
  13 tests. [done 2026-06-23] Apple #2922.

- [x] **S83-02: Merit points** — `server/merit/merit.go`: MeritBank; Earn(1000XP=1merit, cap30); Spend(tier cost=cur+1);
  ErrMeritCapReached/TierCapped/InsufficientMerits; TotalSpent(). 13 tests. [done 2026-06-23] Apple #2927.

- [x] **S83-03: Item level system** — `server/gear/gear.go`: 16 slot consts; Equipment.Equip/Unequip/ItemAt/EffectiveIL/
  OccupiedCount; empty slots excluded from avg; ErrNoEquipment/SlotEmpty/UnknownSlot. 12 tests. [done 2026-06-23] Apple #2928.

### ECONOMY + CRAFTING

- [x] **S84-01: Crafting guilds** — `server/craft/craft.go`: 8 CraftType consts; CraftSkill.Level/SetLevel[0,110];
  SuccessChance(50+gap*2 floor5 cap95); Attempt(recipe,skill,rng)->Result{Success,HQTier,ItemID}; ErrBreak at gap≤-10.
  15 tests. [done 2026-06-23] Apple #2929.

- [x] **S84-02: Synthesis HQ system** — same file: HQTier(gap,roll): 0=NQ, 1=gap≥10, 2=gap≥30 roll<0.15,
  3=gap≥50 roll<0.05; Recipe.HQ1/2/3ItemID wired into Attempt. 7 tests. [done 2026-06-23] Apple #2930.

- [x] **S84-03: Auction house** — `server/market/ah.go`: 8 categories (Weapons/Armor/Ammo/
  Food/Crystals/Materials/CraftItems/Misc). Listing{ID,ItemID,ItemName,Category,Price,Qty,
  SellerID}. Buy(buyerID,itemID) → cheapest listing; self-purchase blocked; SaleRecord written.
  CancelListing(). BrowseCategory() → ItemSummary per item cheapest-first. ItemPage() →
  active listings + last 12 sales (FFXI PS2 terminal parity). SellerListings/BuyerHistory.
  AHFeePct=4 HistoryCap=12. 31 tests. [done 2026-06-22] Apple #2547.

- [x] **S84-04: Swampville secondary starting zone** — Zone 3 `DefaultZones()` SpawnY=1.
  Scene 3 swamp terrain: clay/mud floor, 25% water pools, dead mangrove trees.
  Mobs: leech (aggressive, fast swing), slime (passive, HP=140), lizard (aggressive, fast).
  SwampvilleSpawns(): 6 leeches inner r≈20, 4 slimes mid r≈32, 4 lizards outer r≈42.
  20 mob tests. [done 2026-06-23] Apple #2865.

- [x] **S84-05: Mining skill** — `server/gather/mining.go`: FFXI-parity MiningPoint (skill 0-110,
  difficulty 0-110). SuccessBase=60% +2%/lvl, floor 5%, cap 95%. HQ rolls: skill>diff+10 →
  1%/lvl HQ chance, cap 25%. HQ: alternate item (HQItemID) or double qty. Depletion+Respawn.
  Preset points: MeadowMiningPoints (diff 0-10) + SwampMiningPoints (diff 15-40, Swampite+WaterCrystal).
  27 tests. [done 2026-06-23] Apple #2868.

- [x] **S84-06: DragonsNShit MUD server** — `apps2/mud/main.go`: playable TCP text MUD on :2323.
  Connects: nc localhost 2323. Commands: look/l, n/s/e/w, attack, stop, ws (weapon skill 3×dmg),
  mine, mine-points/mp, status/st, who, say, tell, map, mobs, help, quit.
  1Hz game loop: mob.Tick (AI/aggro/burrow), mob.TickPlayer (auto-attack + TP gain),
  status.Stack.Tick (Poison/Regen/Refresh), mob respawn queue (60s). 4 zones with exits.
  All server packages wired: zone.Manager, mob.Registry, combat/tp, status, gather, worldapi.
  [done 2026-06-23] Apple #2887.

### PARTY + ALLIANCE SYSTEM

- [x] **S85-01: Party system** — `server/party/party.go`: Party{Leader, Members[]}, Invite/Kick/Leave
  (leader transfers on leave), XPSplit(xp, dists, limit) even split OOR=0. 6-player cap. [done 2026-06-23] Apple #2904.

- [x] **S85-02: Alliance system** — same file: Alliance{parties[3]}, AllianceLeader=party[0].Leader,
  TotalSize, ErrAllianceFull at 3 parties. 18-player max. [done 2026-06-23] Apple #2904.

- [x] **S85-03: Party XP bonus** — same file: XPChain.Record(nowNano)→bonus%, +10%/kill cap 50%,
  3min ChainTimeout resets on expiry. BonusPct()/Expired()/Reset(). 37 tests. [done 2026-06-23] Apple #2904.

### TRAVEL + WORLD

- [x] **S86-01: Outpost / conquest system** — `server/conquest/conquest.go`: Nation int (Neutral/Sandoria/Bastok/Windurst),
  Region.AddPoints/Tick (majority wins; incumbent retains on tie; points reset), Map.TickAll/Scoreboard/RegionCount,
  DefaultRegions() 4 zones. 19 tests. [done 2026-06-23] Apple #2908.

- [x] **S86-02: Home point teleport** — `server/homepoint/homepoint.go`: Crystal{SceneID,Pos}; State.SetHome/ReturnHome
  8% XP penalty on KO; ErrNoHomePoint/ErrNotKO. 10 tests. [done 2026-06-23] Apple #2914.

- [x] **S86-03: Survival guide / field manuals** — `server/field/field.go`: Manual{SceneID,BonusPct,ExpiresAt};
  Active/ApplyBonus/ApplyAll stacking; 100% XP bonus; MeadowSurvivalGuide/SwampFieldManual presets. 14 tests. [done 2026-06-23] Apple #2915.

### CONTENT

- [x] **S87-01: Notorious Monster (NM) spawn conditions** — `server/nm/nm.go`: NMSpawn{PlaceholderID, SpawnChance,
  WindowOpen/Close}; OnPlaceholderKilled opens absolute window; InWindow/TrySpawn/WindowExpired/Reset;
  MeadowNMs/SwampNMs presets. 14 tests. [done 2026-06-23] Apple #2909.

- [x] **S87-02: Treasure pool system** — `server/loot/loot.go`: Pool{MobID,Eligible,Items}; Lot(slot,itemID)->roll 1-999;
  Pass->0; Ready(); Resolve() highest roll wins, all-pass=no winner, tie-break slot asc; partial OK.
  15 tests. [done 2026-06-23] Apple #2910.

- [x] **S87-03: Notorious Monster aggro types** — `server/mob/aggro.go`: AggroType{Sight,Sound,JobDetect};
  StatusFlags{Invisible,Sneak}; Detects/DetectsDefault 20m sight 45° half-cone 8m sound; Invisible→no sight,
  Sneak→no sound, JobDetect unblockable; AggroSight/Sound/SightSound/Job/Passive presets. 16 tests. [done 2026-06-23] Apple #2916.

---

## SECTION 88: DRAGONSNSHIT MUD — PROGRESSION WIRING (2026-06-23)

*Context: S81-S87 built all FFXI parity packages (XP, party, homepoint, field manuals, NMs, loot,
mob aggro). The MUD from S84-06 doesn't wire any of them. This section wires the progression loop
into the playable MUD so kills grant XP, players level up, parties split XP, and KO means something.*

- [x] **S89-01: MUD real weapon skills + skillchain detection** — Wire `server/skillchain` into
  `apps2/mud/main.go`. Player has `wsSkill string` (default "Fast Blade"). `ws [skillname]` uses
  the named WS or current default. `setws <name>` picks a WS. `wslist` shows available skills +
  resonances. Per-mob chain state tracks last WS attrs + time; 8s window; chain announces tier name
  + bonus damage. Mob dies → clears chain state. All 29 GFD packages green.
  [done 2026-06-23] Apple #2964.

- [x] **S89-02: MUD job system** — Wire `server/job` into the MUD. `setjob <JOB>` command (e.g.
  setjob WAR). Job sets max HP/MP via HPAtLevel/MPAtLevel. `jobs` command lists all 22 jobs.
  Default job=WAR. Status shows current job. HP/MP grow on level-up.
  [done 2026-06-23] Apple #2966.

- [x] **S90-01: MUD crafting system** — Wire `server/craft` into the MUD. Player inventory tracks
  item counts by ID (from loot drops + mining). `craft <recipe>` command: look up recipe, check
  ingredients in inventory, Attempt() with player's craft skill, consume ingredients on success,
  add output item to inventory. `craft-skills` shows player's craft skill levels. `recipes` lists
  known recipes. Recipes: iron ingot (2× earth crystal + iron ore → iron ingot), herbal remedy.
  [done 2026-06-23] Apple #3048.

- [x] **S90-02: MUD conquest system** — Wire `server/conquest` into the MUD. Zone kills award
  conquest points to the nation the player has declared for (default: none/Neutral). `conquest`
  command shows current nation control of each zone. `declare <nation>` picks a nation.
  Tick conquest once per minute. Scoreboard shows nation point totals.
  [done 2026-06-23] Apple #3051.

- [x] **S90-03: MUD auction house** — Wire `server/market` into the MUD. Player inventory items
  can be listed for sale. `ah sell <item-id> <price>` lists item. `ah browse [category]` shows
  listings. `ah buy <listing-id>` deducts gil + transfers item. `ah history <item-id>` last sales.
  [done 2026-06-23] Apple #3053.

- [x] **S89-03: MUD loot pool + NM spawns** — Wire `server/loot` and `server/nm` into the MUD.
  Mob kill → dropsForMob table (worm/leech/slime/lizard/NM kinds). Solo: auto-award. Party:
  shared loot.Pool; lot/pass/pool commands. NM spawn checked each tick; King Worm + Marsh Leech
  NMs defined. Placeholder kill → OnPlaceholderKilled → window open → TrySpawn.
  [done 2026-06-23] Apple #2969.

- [x] **S88-01: MUD XP + leveling + homepoint + field manuals + party** — Wire `server/xp`,
  `server/homepoint`, `server/field`, `server/party` into `apps2/mud/main.go`. Mob kills grant
  XP (MaxHP*10). Field manual (+100% XP, 30min). KO on HP=0 → `home` returns to Home Point
  (8% XP penalty). Party invite/accept/leave/pt; XPChain bonus (+10%/kill, cap 50%). Level-up
  announcements. `status` shows Lv/XP/homepoint/party/chain. All 29 GFD packages green.
  [done 2026-06-23] Apple #2962.

---

## SECTION 91: SIGNAL PIPELINE WATCHLIST EXPANSION (2026-06-23)

*Source: Emily Prime expand_coverage task queue (424 tasks). 16 tickers disabled pending CIK verification.*

- [x] **S91-01: Enable 16 watchlist tickers — fix 5 wrong CIKs** — Verified all disabled entries via
  SEC company_tickers.json. Fixed wrong CIKs: COIN (1679273→1679788), SOFI (1818502→1818874),
  HOOD (1783398→1783879), GE (40987→40545), DIS (1001293→1744489). Enabled UNH/MRK/TGT/LOW/CAT/RTX/T/VZ/CMCSA/NEE/COP.
  All 50 entries now enabled at poll_priority=2.
  [done 2026-06-23] Apple #3056.

---

## SECTION 92: DRAGONSNSHIT MUD — ENMITY + CHAT ROUTER + GUILD (2026-06-23)

*Wire remaining unimported server packages into apps2/mud/main.go.*

- [x] **S92-01: MUD enmity (hate) system** — Wire `server/enmity` into the MUD. Each mob has
  an enmity.Table keyed by player slot. Damage adds CE; healing adds CE (AoE halved). Mob always
  targets the player with the highest enmity (Top()). `enmity` command shows your threat vs active
  mob. Enmity clears on mob death. Enforce: healer healing a party member draws aggro.
  ✓ Apple #3062

- [x] **S92-02: MUD chat router** — Wire `server/chat` Router into the MUD. Player sessions tracked
  in chat.Session{Slot, Pos, Zone, GuildID}. `say <msg>` uses Chat channel (50u radius).
  `yell <msg>` zone-wide. `tell <player> <msg>` private. `guild <msg>` linkshell-wide.
  Routing replaces the current raw broadcast for say/tell.
  ✓ Apple #3063

- [x] **S92-03: MUD linkshell guild system** — Wire `server/guild` Registry into the MUD. `ls-create
  <name> <tag>` founds a linkshell. `ls-invite <player>` forges a Feather. `ls-leave` revokes.
  `ls-info` shows roster. GuildID in player struct; guild chat gated via CanChat().
  ✓ Apple #3066

---

## SECTION 93: DRAGONSNSHIT MUD — GEAR + SUB-JOB + MERIT (2026-06-23)

- [x] **S93-01: MUD equipment system** — Wire `server/gear` Equipment into the MUD. Players have
  16 gear slots. `equip <slot> <item-id>` equips from inventory. `unequip <slot>` returns to
  inventory. `gear` command shows equipped items + effective item level. Gear drops added to select
  mob drop tables (e.g. lizard drops leather-belt). EffectiveIL shown in `status`.
  ✓ Apple #3069

- [x] **S93-02: MUD sub-job system** — Wire `server/job.CharJob` (subjob.go) into the MUD.
  `setsubjob <JOB>` picks a sub-job. Sub-job grants half-level stats via CombinedStats(). `status`
  shows main/sub job. ErrSameJob: cannot set same as main.
  ✓ Apple #3070

- [x] **S93-03: MUD merit points** — Wire `server/merit` MeritBank into the MUD. Each level-up
  at L75+ earns 1000 merit XP toward 1 merit point. `merits` command shows bank + tier counts.
  `merit-spend <category>` spends a point (e.g. merit-spend HP/MP). MeritBank in player struct.
  ✓ Apple #3071

---

## SECTION 94: DRAGONSNSHIT MUD — TELECRYSTAL TRAVEL (2026-06-23)

- [x] **S94-01: MUD telecrystal fast travel** — Wire `server/telecrystal` Registry into the MUD.
  Crystals exist in zones (e.g. MeadowCrystal zone=0). `crystals` lists nearby crystals (InScene).
  `travel <crystal-id>` validates range + gold balance (Validate()), deducts 200 gil, teleports
  player to crystal pos. `touch` activates the nearest crystal (logs its ID). Crystal range=10u.
  ✓ Apple #3073

---

## SECTION 95: DRAGONSNSHIT MUD — WORLD CRISIS WIRING (2026-06-23)

*Wire `server/worldcrisis` Crisis state machine into apps2/mud. Players fight in a living world
with escalating phases and objectives. The 1Hz loop ticks the Crisis; NM kills + gather
objectives advance it. Broadcast phase transitions to all players.*

- [x] **S95-01: MUD World Crisis — wire server/worldcrisis, `crisis` command, phase broadcast** —
  Add `*worldcrisis.Crisis` to world struct. Start a Crisis on first player login. 1Hz tickAll
  drives `gw.crisis.Tick(now)` and broadcasts phase-change banners to all players. `crisis` command
  shows phase, LEY integrity, objectives complete, time remaining. Objectives: kill 3 NMs in Swamp
  (CompleteObjective(ObjKillNM, 10)) and gather 5 crisis-materials (gather action on special SwampCrisis
  points CompleteObjective(ObjGather, 5)). `crisis-ley <amount>` lets a GM decay LEY for testing.
  On OutcomeSaved/Failed, broadcast final result. GOWORK=off go test ./... must pass.

- [x] **S95-02: MUD crisis mob — World Crisis NM spawns in Swamp zone on Crisis start** —
  When Crisis.Phase transitions to PhaseMobilization, spawn 3 special "chaos-elemental" NMs in
  zone 3 (swamp). On death, CompleteObjective(ObjKillNM, 10) + broadcast "[Crisis] Chaos Elemental
  defeated! LEY stabilized." Extra drop: "crisis-shard" (new item). `crisis-shards` go in inventory
  and can be sold at AH.

---

## SECTION 96: DRAGONSNSHIT MUD — INVISIBLE/SNEAK AGGRO (2026-06-23)

*Wire `server/mob/aggro.go` StatusFlags into the MUD mob aggro scan. Players can use the
status.Stack Invisible/Sneak effects to evade mob aggro.*

- [x] **S96-01: MUD mob aggro wire — check Invisible/Sneak before aggroing** —
  In `tickAll()` mob.Tick loop, before applying aggro results, build a `mob.StatusFlags` per
  player from their `p.statFX.Stack` (check if Invisible or Sneak status is active). Pass this
  to a new helper `shouldAggro(m *mob.Mob, p *player) bool` that calls
  `mob.AggroDefault(m.Kind).Detects(pos, target, flags, dist, yawDiff)`. Only aggro if Detects()
  returns true. Commands: `cast invisible` (apply Invisible status for 60s, no sight aggro),
  `cast sneak` (apply Sneak status for 60s, no sound aggro). `cast` uses MP (50 MP each).

---

## SECTION 97: DRAGONSNSHIT MUD — ABILITY RECAST SYSTEM (2026-06-23)

*Wire `server/job.RecastTracker` into the MUD. Job abilities (Provoke, Berserk, etc.) have
real recast timers. `ja <ability>` command uses the ability if ready.*

- [x] **S97-01: MUD job ability commands — RecastTracker, `ja <ability>`, `recasts`** —
  Add `recastTracker *job.RecastTracker` to player struct. Initialized from the player's main
  job abilities (WarriorAbilities() for WAR, WhiteMageAbilities() for WHM, etc.). `ja <ability-id>`
  checks `t.Ready()`, calls `t.Use(id, now, charXP.Level)`, applies effect (Provoke: add 2000
  enmity CE on current target; Berserk: add Haste+30 status effect to self for 3min).
  `recasts` command lists all abilities with remaining recast time.
  Error display: `<ability> is on recast (2m15s remaining)`.

---

## SECTION 98: DRAGONSNSHIT MUD — IDUNA CHARACTER PERSISTENCE (2026-06-23)

*Wire `server/idunaclient` into apps2/mud. Players get persistent characters stored in IDUNA.
On login, fetch or create the character. Level-ups, gil changes, position updates POST to IDUNA.*

- [x] **S98-01: MUD IDUNA login — fetch/create character from IDUNA on connect** —
  Add `idunaclient.Client` to world struct. On player connect (after name prompt), call
  `GET /api/v1/characters?name=<name>`. If found: load level/hp/gil/position from IDUNA.
  If not found: `POST /api/v1/characters` to create it. `-iduna-url` flag (default: $IDUNA_BASE_URL).
  No-IDUNA mode: fall back to defaults (level 1, 500 Gil) when IDUNA_BASE_URL unset.
  Log: `[IDUNA] character <name> loaded/created (id: <cid>)`.

- [x] **S98-02: MUD IDUNA save — level-up + disconnect position save** —
  On level-up in grantXP(), call `PATCH /api/v1/characters/<cid>` with new level + hp + mp (goroutine).
  On disconnect (defer in handleConn), PATCH position + gil + guildID. On death, PATCH hp=0.
  No blocking on network: fire-and-forget goroutine with 5s timeout context.

---

## SECTION S99: GFD MUD — MOB AI + SPELLCASTING + BUFF SYSTEM

- [x] **S99-01: MUD mob spellcasting AI** — Mobs cast status spells during combat. In `resolveKill`
  loop (post-mob-attack), add a 20% chance that the mob "casts" a debuff on its aggro target.
  Spell pool: Poison (DoT), Paralyze (25% action fail chance), Bind (movement stop), Silence
  (blocks `cast` commands). Use `status.Stack.Apply()` on the target player. Add `removedebuffs`
  command (consumable: requires 1 "echo-drop" item in inventory). Message player on debuff application.

- [x] **S99-02: MUD cure spells** — Wire `cast cure` and `cast cure2` (White Mage only). `cast cure`
  restores 100 HP (50 MP). `cast cure2` restores 250 HP (80 MP). Both require WHM job or sub.
  In cmdCast, check `p.jobID == job.WHM || p.charJob.SubJob == job.WHM`. Dispel silence check:
  Silenced players cannot cast any spell. Blocked by Silence status.

- [x] **S99-03: MUD `bazaar` command** — Open personal shop. `bazaar set <item> <price>` puts item
  in personal bazaar at price. `bazaar list` shows your listing. `bazaar buy <player> <item>` buys
  from player's bazaar (deducts gil, transfers item). Persists to `world.bazaars map[slot]bazaarStore`.

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

## SECTION S100: GFD MUD — COMBAT DEPTH + PARTY TOOLS

- [x] **S100-01: MUD rest/meditate HP+MP regen** — `rest`/`meditate` command. `isResting bool` on
  player. In tickAll: if isResting, +5% maxHP and +3% maxMP per second. Any attack or mob aggro
  cancels rest. `stand` cancels manually.

- [x] **S100-02: MUD `target <mob>` reticle** — Explicit target selection by mob kind substring.
  Sets p.combat.TargetMobID. Shows HP bar. Integrates with `attack`, `ws`, `ja provoke`.

- [x] **S100-03: MUD party chat `/p <msg>`** — Lines starting with `/p ` broadcast to party members
  only. Wire through chat.Router channel 4 (Party). No tab completion needed.

---

## SECTION S101: GFD MUD — ECONOMY + WORLD EVENTS

- [x] **S101-01: MUD `bank` command (gil deposit/withdraw)** — `bank deposit <amount>`: p.bankGil int
  field. `bank withdraw <amount>`. `bank balance`. Save bankGil to IDUNA on disconnect.

- [x] **S101-02: MUD random weather events** — `weather` command. `weatherByZone map[int]string` on
  world. In tickAll, every 60s, 10% chance new weather per zone (Clear/Cloudy/Rain/Thunder/Fog).
  Rain = +10% water-crystal drop chance. Thunder = +10% fire-crystal drop chance.

- [x] **S101-03: MUD `survey` command** — Shows all players in zone with direction (N/S/E/W from pos).
  If same party: show HP bar. Shows: name, job, level, direction.

---

---

## SECTION 102: DYNAMIC HYBRID AI AGENT ARCHETYPE ENGINE (2026-06-23)

*Northstar: `EMILY/docs/ARCHETYPE_ENGINE_NORTHSTAR.md`*
*THE_FIELD as a real Go service. E₁/E₂ dual-persona + 72-spirit Goetia bank + resonance routing.*
*Milestone 1 complete. Remaining: dual-persona invocation, interference engine, HTTP service, integrations.*

- [x] **S102-01: Goetia bank — all 72 spirits as Go structs + intent selector** —
  `emily-agent/pkg/archetypes/` package. `All72 []Spirit` with freq/amp/phase/corridor/seed phrase
  for all 72 spirits (M_G v1.0). `ByName()`, `ByRank()`, `ForCorridor()`, `AmplitudeSum()`.
  `MatchIntent(intent, allowCollapse)` → stack of ≤3 spirits via keyword scoring.
  `ResonanceState(Δφ)` → Corridor. `E1Weight/E2Weight` → blend weights.
  17 tests pass. Apple #3210 | 2026-06-23.

- [x] **S102-02: Dual-persona E₁/E₂ invocation** — `emily-agent/pkg/archetypes/field.go`.
  `Field.Invoke(ctx, intent, context string, allowCollapse bool) (*InvokeResult, error)`.
  Calls E₁ (claude-sonnet, temp 0.3, structured identity prompt) and E₂ (claude-haiku, temp 1.8,
  divergence prompt) concurrently. Returns both outputs + Goetia stack + seed phrases injected.
  Seed phrases injected at the end of each system prompt at the selected spirit's phase offset.
  Apple #3212 | 2026-06-23.

- [x] **S102-03: Interference engine + synthesizer** — Compute Δφ from cosine similarity of E₁/E₂
  outputs (word-level token overlap as proxy for embedding cosine). Classify resonance state.
  Blend: `output = E1Weight*e1 + E2Weight*e2`. Amplitude-scale each spirit's seed contribution.
  Collapse guard: if Δφ > 90°, return E₁ output only.
  Apple #3212 | 2026-06-23.

- [x] **S102-04: HTTP service `cmd/archetype-engine`** — POST /invoke → InvokeResult JSON.
  GET /spirits → All72 JSON. GET /status → current resonance state. :8090.
  ANTHROPIC_API_KEY required. `-allow-collapse` flag for high-risk spirits.
  Apple #3214 | 2026-06-23.

- [x] **S102-05: Emily Prime RSI integration** — In `emily-agent/rsi.go` DECIDE phase, route
  the current task intent through `field.Invoke()`. Use Vassago #03 (foresight) + Eligos #15
  (strategy) as default RSI stack. Inject blended output as enhanced task description.
  Apple #3216 | 2026-06-23.

---

## SECTION 103: ARCHETYPE ENGINE FATBABY + SHANKPIT INTEGRATION (2026-06-23)

*Full integration of THE_FIELD into PRRJECT_FATBABY signal analysis and SHANKPIT game AI.*
*Northstar: `EMILY/docs/ARCHETYPE_ENGINE_NORTHSTAR.md` Milestones 8 + 9.*

- [x] **S103-01: ArchetypeProvider — THE_FIELD routing for FatBaby signal analysis** —
  `PRRJECT_FATBABY/internal/processor/archetype_provider.go`. `intentForSource()` maps
  SEC form types to spirit stacks: 8-K→Aim#23(fraud)+Eligos#15, 10-K→Eligos+Vassago+Marax,
  SC13D→Bune+Aim, DEF14A→Orobas+Botis, PR→Forneus+Vassago. Fallback to haiku when engine down.
  `resonance_state`, `phase_delta_deg`, `intent` stamped in `Signal.RawMetadata`.
  cmd/processor: `-archetype-engine` flag activates. 9 tests pass. Apple #3219 | 2026-06-23.

- [x] **S103-02: SHANKPIT emily-bot archetype policy** — Wire archetype engine into
  `SHANKPIT/apps2/emily-bot/main.go`. `-archetype-engine` flag (default:
  `http://localhost:8090`). In `think()`: serialize state → POST /invoke with
  `intent="precision combat strike martial clarity"` (Leraje#14+Marchosias#35) →
  parse `output` as action tokens → decode UserCmd. Chaos mode: `allow_collapse=true`
  (Andras#63). Falls back to heuristic when engine unavailable.
  Apple #3222 | 2026-06-23.

- [x] **S103-03: Archetype engine operational metric in EMILY status** — POST /status to
  `/api/v1/emily/archetype/status` proxy. `emily status` shows archetype engine state
  (E1/E2 models, last resonance corridor, cycle count). emily-agent proxies :8090/status.
  Apple #3224 | 2026-06-23.

---

## SECTION 104: DRAGONSNSHIT MUD — BST PET SYSTEM + NPC QUESTS (2026-06-23)

*GoblinFoxDragon MUD depth: Beastmaster pet companions + quest-giving NPCs.*
*Goal: BST job summons a pet that fights autonomously. NPCs give quests with item/gil rewards.*

- [x] **S104-01: BST pet system — `server/pet/pet.go`** — `Pet{Kind, HP, MaxHP, Level, OwnerSlot}`
  8 tameable mob kinds (wolf/bird/lizard/crab/leech/slime/worm/bear). `Tame(p *mob.Mob) bool`:
  success chance = 40%+2%/lvl, fail if mob HP>70%. Pet fights its owner's target (auto-attack 1.5s
  swing). Pet dies → nil; respawn after 5min via `JugPet(kind)`. `PetStatus()` shows HP.
  `petHeal <slot>` command (WHM sub: cure pet). 20 tests.
  Apple #3228 | 2026-06-23.

- [x] **S104-02: BST MUD wiring** — `bst` command: `charmed <mob>` if BST and mob HP<50%.
  Pet auto-attacks owner's target each tick. `pet-release` dismisses. `pet-status` shows pet HP.
  `pet-heel` (pet stops attacking, follows). Pet XP grants to owner. Pet slot on player struct.
  `setjob BST` enables `bst` command.
  Apple #3231 | 2026-06-23.

- [x] **S104-03: NPC quest system — `server/quest/quest.go`** — `Quest{ID, Title, GiverNPCID,
  RequireItems map[string]int, RequireKills map[string]int, RewardGil int, RewardItem string}`.
  `State{Active, KillProgress map[string]int, Complete}`. `Accept/TurnIn` flow.
  5 starter quests: "Gather Iron Ore" (3 ore→100gil), "Defeat the King Worm" (1 NM kill→500gil +
  Iron Sword), "The Merchant's Request" (5 Crystals→200gil), "Swamp Patrol" (kill 5 leeches→300gil
  + leather belt), "Crisis Volunteer" (1 crisis-shard→400gil + Echo Drop).
  Apple #3233 | 2026-06-23.

- [x] **S104-04: NPC dialogue MUD wiring** — `npc` command: `npcs` lists NPCs in zone.
  `talk <npc-id>` shows dialogue + available quests. `quest-accept <quest-id>`. `quest-turn-in <id>`.
  `quests` lists active + completed quests. NPCs: Guildmaster (Meadow), Merchant (Hills), Scout (Swamp).
  Apple #3235 | 2026-06-23.

---

## SECTION 105: DRAGONSNSHIT MUD — PLAYER MAP + WEATHER SYSTEM (2026-06-23)

*Exploration depth: players build a zone map as they explore. Weather modulates mob behavior + combat.*

- [x] **S105-01: Zone map system — `server/cartography/cartography.go`** — `Atlas` tracks which
  zones a player has visited. `Visit(zoneID)`. `KnownZones() []int`. `ExitMap(exits map[int]map[string]int) string`
  builds ASCII exit map showing visited vs unknown zones. 10 tests.
  Apple #3238 | 2026-06-23.

- [x] **S105-02: Cartography MUD wiring** — `map` command shows ASCII zone map (visited=zone name,
  unvisited=???). `Atlas` on player struct, `Visit()` called on zone entry. `explore` alias.
  Apple #3238 | 2026-06-23.

- [x] **S105-03: Weather system — `server/weather/weather.go`** — `Phase{Clear,Overcast,Rain,Storm}`
  and `Engine{phase,phaseEnds}`. `Tick(now) (changed bool, old, new Phase)`. Duration: Clear=5-12min,
  Overcast=3-8min, Rain=2-6min, Storm=1-4min. Storm: all mob damage +15%, ChaosEdge resonance for BST
  tame (+10% success). 12 tests.
  Apple #3240 | 2026-06-23.

- [x] **S105-04: Weather MUD wiring** — Weather goroutine in MUD server. `weather` command shows current
  phase. Weather changes broadcast to all players in zone. Storm buffs mob MeleeDamage via modifier.
  BST tame: if storm, add +10% to chance. Show weather on player prompt line.
  Apple #3240 | 2026-06-23.

---

## SECTION 106: DRAGONSNSHIT MUD — PVP DUEL SYSTEM + GLOBAL LEADERBOARD (2026-06-23)

*Consensual PvP: players challenge each other to duels. Winner gains duel rating. Leaderboard shows top duelists.*

- [x] **S106-01: Duel system — `server/duel/duel.go`** — `State{Challenger,Defender,Phase}` where
  Phase: Pending/Active/Done. `Manager{pending,active map}`. `Challenge(challenger,defender) error`.
  `Accept(defender) error`. `Forfeit(slot) error`. `Tick(challHP,defHP,now) (winner,done bool)`.
  Done when either player HP≤0 or forfeit. Auto-expire pending challenges after 60s.
  Rating: winner +25 pts, loser -10 pts; tracked in `RatingMap map[string]int`. 15 tests.
  Apple #3243 | 2026-06-23.

- [x] **S106-02: Duel MUD wiring** — `duel <player-name>` sends challenge. `duel-accept` accepts
  pending challenge. `duel-forfeit` concedes. While dueling: combat hits opponent instead of mob.
  On win: announce to zone, update rating, XP=0. `leaderboard` shows top 10 duelists by rating.
  Apple #3243 | 2026-06-23.

---

## SECTION 107: PRRJECT_FATBABY — SIGNAL REPLAY + BACKTESTER (2026-06-23)

*Replay stored signals through the analysis pipeline to backtest provider changes.*

- [x] **S107-01: Signal replay engine — `internal/replay/replay.go`** — `ReplayEngine{EventStore,
  Provider, Out chan Signal}`. `Replay(ctx, from, to time.Time) (replayed int, err error)`.
  Reads events from the event store ordered by timestamp, re-runs each through `Provider.AnalyzeText()`,
  streams results to Out channel. `Stats{Total, Success, Failed, Duration}`. 10 tests. Apple #3246

- [x] **S107-02: Replay CLI — `cmd/replay/main.go`** — Flags: `-from 2006-01-02`, `-to 2006-01-02`,
  `-provider haiku|archetype|stub`, `-output signals.jsonl`. Prints per-signal result summary to
  stdout. Uses `internal/replay.ReplayEngine`. Builds standalone binary. Apple #3247

---

## SECTION 108: DRAGONSNSHIT MUD — FISHING + FOOD BUFFS (2026-06-23)

*FFXI-parity gathering depth: fishing skill + cooked food stat buffs.*
*Fishing: server/gather/fishing.go (mirrors mining.go). Food: server/food/food.go (stat buff on eat).*

- [x] **S108-01: Fishing skill — `server/gather/fishing.go`** — `FishingPoint{ID,SceneID,Difficulty,Yield}`
  with species loot table (carp/bass/trout/eel/crab/fullmoon-sardine). `FishSkill.Attempt(pt,rng)`
  same success model as mining (60%+2%/lvl, floor 5%, cap 95%). HQ: trophy fish variant.
  Depletion: 10 attempts per point. `MeadowFishingPoints` (diff 0-10) + `SwampFishingPoints` (diff 20-40).
  `fish` command in MUD: `fish` at a nearby point → attempt → yield or "The line snaps." 15 tests.
  Apple #3251 | 2026-06-23.

- [x] **S108-02: Food buff system — `server/food/food.go`** — `Food{ID,Name,STRBonus,DEXBonus,VITBonus,
  HPBonus,Duration}`. `FoodEffect` applied to player stats. `eat <item-id>` command. Effect duration
  10-30min. Sample foods: Meat Mithkabob (STR+5, 30min), Tavnazian Salad (DEX+5, 30min), Rolanberry
  Pie (MP+50, 20min), Selbina Milk (VIT+3, HP+30, 10min). Food created via craft recipes. 12 tests.
  Apple #3252 | 2026-06-23.

---

## SECTION 109: DRAGONSNSHIT MUD — NATION FAME SYSTEM (2026-06-23)

*FFXI parity: players earn fame in each nation by completing quests. Fame gates advanced NPCs and quests.*

- [x] **S109-01: Fame system — `server/fame/fame.go`** — `Nation` enum (Bastok/Sandoria/Windurst/Neutral).
  `Store{points map[Nation]int}`. `Earn(n Nation, pts int)`. `Rank(n Nation) int` (0-5 thresholds:
  0/50/200/500/1000/2000). `TotalPoints(n)`. `NewStore()`. Quests gain `RewardFame int + FameNation Nation`.
  `Journal.TurnIn()` grants fame via Store. 12 tests. Apple #3255

- [x] **S109-02: Fame MUD wiring** — `fame` command shows rank + points per nation. Quests updated with
  fame rewards. `StartingQuests` for Sandoria (Guildmaster quests) + Bastok (Merchant) + Windurst (Scout).
  `talk <npc>` shows fame requirement if present. Fame gate: rank < required → "Your reputation isn't
  strong enough yet." Apple #3257

## SECTION 110: DRAGONSNSHIT MUD — GUILD MANAGEMENT + NPC VENDOR (2026-06-23)

*Completing the linkshell social loop and adding direct NPC shopping.*

- [x] **S110-01: Linkshell kick + promote** — `ls-kick <player>`: Officer+ removes a member's Feather via
  `RevokeItem`. `ls-promote <player>`: Leader forges a Feather Sack for a Member, promoting to Officer
  via `ForgeFeatherSack`. Both broadcast a linkshell announcement to online guild members. Apple #3260

- [x] **S110-02: NPC vendor shop** — `shop` command at NPC zones. Each zone has a static item catalog
  (price in gil). `shop` lists items. `shop buy <item-id>` deducts gil and adds to inventory. Items:
  basic consumables (echo-drop, antidote, hi-potion) and zone-specific gear. `shop sell <item-id>` sells
  player inventory item at 50% NPC buy price. Apple #3260

## SECTION 111: DRAGONSNSHIT MUD — WHITE MAGIC BUFF SPELLS (2026-06-23)

*Wiring the status.Stack buff/debuff system into WHM/RDM spell casting.*

- [x] **S111-01: WHM/RDM buff spells — Protect/Shell/Haste/Regen/Refresh/Dia** —
  `cast protect/shell/haste/regen/refresh` apply status.Effect (Protect/Shell/Haste/Regen/Refresh)
  to player's statFX stack. `cast dia` deals 5 immediate HP to combat target. All require WHM job
  (Refresh also allows RDM). MP costs: Protect 60, Shell 60, Haste 75, Regen 40, Refresh 50, Dia 30.
  Duration: 3 minutes. Tick loop already handles Regen/Refresh events. Apple #3262

## SECTION 112: DRAGONSNSHIT MUD — PARTY-TARGETED SPELLS (2026-06-23)

- [x] **S112-01: Party-targeted spell casting — `cast <spell> <player>`** —
  `resolveSpellTarget(p, targetName)` finds players in same zone. Cure/Cure II/Protect/Shell/
  Haste/Regen/Refresh all accept optional player target; notifications sent to both caster and
  target. Apple #3265

## SECTION 113: DRAGONSNSHIT MUD — BLM BLACK MAGIC NUKES (2026-06-23)

- [x] **S113-01: BLM black magic nukes — Fire/Blizzard/Thunder/Stone/Water/Aero I-III** —
  18 spells (6 elements × 3 tiers). INT scaling (+1 dmg per INT >10). BLM/RDM main/sub required.
  Silence check. Requires combat target. `cmdCastBlackMagic` shared handler. Apple #3267

## SECTION 114: DRAGONSNSHIT MUD — BRD SONGS (2026-06-23)

- [x] **S114-01: BRD zone-AoE songs — march/paeon/ballad/minne/carol/mambo** —
  6 songs apply Haste/Regen/Refresh/Protect/Shell to all players in zone. BRD job/sub required.
  3-minute duration. Party notification broadcast. Apple #3269

---

*EMILY PRIME BACKLOG | Cross-repo | Git-authoritative*
*The backlog is what outlasts everything.*
*Clean builds first. Then custody. Then everything else.*
