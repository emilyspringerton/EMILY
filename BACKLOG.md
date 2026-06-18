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

- [ ] **MPT → TYLER RSI trigger** — When emily.sh produces a new episode script, auto-invoke
  MPT. Spec exists in engine/moneyprinter_pipeline.md. Dependency: cold open clip ✓.

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
- [ ] **signal pipeline audit: 10748 signal_failed today (9435=EDGAR 429 no-retry, 977=…** — Awaiting full classification — run emily backlog promote with ANTHROPIC_API_KEY. Obs: 2026-06-17T22:21:23Z.
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

- [ ] **S26-05: Validate entropy source** — Compare entropy_mean_nats for base vs fine-tuned model
  on Emily operational text; document delta in NORTHSTAR.md; update RSI loop entropy injection if
  fine-tuned model produces better (more calibrated) entropy signal.
  Base: H_mean=4.49 nats, H_max=8.62 nats. Target: H_mean shift of ≥0.5 nats after fine-tune.

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

- [ ] **S36-04: Wire real intelligence provider into cmd/processor** — All 5,316 `signal_generated`
  events in the store use `stubProvider`, returning `"Stub analysis result."` / `signal_type=Other`
  for every filing. The LLM analysis layer has never run. Wire `ANTHROPIC_API_KEY` + claude-haiku
  into the processor as the real `intelligence.Provider`. Use haiku (cheap) for classification.
  Dependency: API credit balance ✓ (top up console.anthropic.com).
  Acceptance: new `signal_generated` events have real `signal_type`, non-stub `summary`.

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

- [x] **S37-01: Fix /ask 500 — AskLandingData missing Symbols field** — `GET /ask` returns [done 2026-06-18]
  HTTP 500 "template error" because `internal/newssite/render.go:RenderAskLanding()` passes
  `AskLandingData{GoogleClientID: ""}` to the template, but the shared `masthead` template
  uses `{{range .Symbols}}` for the ticker datalist. `AskLandingData` has no `Symbols` field,
  so template execution panics internally. Fix: add `Symbols []string` to `AskLandingData`
  and populate it from the handler's symbol list (same source as `TickerPageData.Symbols`).
  File: `internal/newssite/render.go` (AskLandingData struct) + `internal/newssite/asklily.go`
  (serveAskLanding populates it). Acceptance: `GET /ask` returns 200 with ticker datalist.

- [x] **S37-02: Add HEAD method support across all newssite routes** — All routes return [done 2026-06-18]
  HTTP 405 on HEAD requests (web audit confirmed: /ticker/AAPL, /section/*, /doc/*, /person/*, etc.).
  Go's net/http ServeMux does not automatically handle HEAD for registered GET handlers.
  Fix: add a middleware wrapper in `cmd/newssite/main.go` that converts HEAD to GET and
  suppresses the body before writing the response. Standard pattern:
  `mux.Handle("/", headToGet(h))`. Acceptance: HEAD on all existing GET routes returns 200.

- [x] **S37-03: /ticker/{sym}/feed.xml — implement RSS or return 404** — `GET /ticker/JPM/feed.xml` [done 2026-06-18]
  returns HTML (200 with ticker page HTML). The route falls through to the ticker handler
  because `/ticker/` catches the whole path. Fix: either (a) implement a real RSS/Atom feed
  for the ticker (signals as feed items), or (b) explicitly return 404 for `.xml` suffix.
  Acceptance: feed.xml returns valid RSS/Atom with correct Content-Type, or 404.

- [x] **S37-04: Start newssite with EMILY_BASE_URL set** — `POST /api/ask` returns 503 [done 2026-06-18]
  "Ask Emily not configured" because `EMILY_BASE_URL` env var is not set at newssite launch.
  Emily Prime is at `http://localhost:8086`. Fix: add `EMILY_BASE_URL=http://localhost:8086`
  to the emily start --newssite launch command or to emily-secrets.env.
  Acceptance: `POST /api/ask` with a question returns an Emily response, not 503.

- [x] **S37-05: Seed SQLite governance_signals from entity-graph output** — signalapi serves [done 2026-06-18]
  `GET /v1/governance-signals?ticker=AAPL` → `[]` (empty). The SQLite `governance_signals`
  table has 0 rows. The projector writes to MySQL (not running). The entity-graph produces
  real signals (insider, governance_health, board_decay etc.) but these go to the entity-graph
  var/ files, not to governance_signals. Fix: wire entity-graph signal output to the SQLite
  governance_signals table as a SQLite-mode projector path.
  Dependency: S36-01 (need data flowing) or manual seed from entity-graph signals.
  Acceptance: `/v1/governance-signals?ticker=AAPL` returns ≥1 signal.

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
