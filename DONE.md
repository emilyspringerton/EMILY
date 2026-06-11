# EMILY PRIME — DONE ARCHIVE
*Completed items moved from BACKLOG.md. Git-authoritative.*

## Archived 2026-06-11

### SECTION 1: FOUNDATION (current sprint)

- [x] **IDUNA embedded SQLite backend** — Zero-MySQL startup. Bootstrap auto-detects.
- [x] **Tyler agent registration (Build 0016)** — .claude/tyler_agent.md, settings.json,
- [x] **EMILY/BACKLOG.md created** — This file. Cross-repo golden backlog. Machine-readable.
- [x] **Run emily.sh — first autonomous RSI iteration** — Tyler ran emily.sh. Build 0017 (MPT
- [x] **Wire EMILY RSI engine to IDUNA** — `EMILY/scripts/start-emily-agent.sh` sources
- [x] **Process supervision: IDUNA systemd unit** — Unit file written to
- [x] **Process supervision: emily.sh cron or supervisor** — `TYLER/scripts/cron-emily.sh`

### SECTION 2: MONEYPRINTERTURBO (video pipeline)

- [x] **MPT config.toml written** — LLM: litellm → anthropic/claude-sonnet-4-6.

### SECTION 3: SYSTEM OBSERVABILITY (Emily Prime sees everything)

- [x] **emily CLI (emily.cli repo)** — Full Go CLI v0.5.0. 69 tests, 5 packages.
- [x] **Apples dashboard query** — `EMILY/scripts/apples.sh` written. Authenticates as
- [x] **Cross-repo status script** — `EMILY/scripts/status.sh` written. Git status for all
- [x] **PRRJECT_FATBABY → IDUNA Apple filing** — `EMILY/scripts/sync-fatbaby.sh` syncs

### SECTION 4: SHANKPIT / TYLER GAME ENGINE (lower priority)

- [x] **Tyler mode SHANKPIT spec** — `TYLER/engine/tyler_hum_mechanic.md` written. First
- [x] **TYLER S06E03: "Four Hours"** — Build 0063. S06E03 written. A-269 (East Ferry Ave,
- [x] **TYLER S06E04: "The Third Photograph"** — Build 0064. AMPLIFICATION confirmed (named).
- [x] **TYLER S06E06 (Season Finale): "The Usual Place"** — Build 0066. A-273, East Kirby St,
- [x] **TYLER S06E05: "How Many"** — Build 0065. AMPL-001 first entries. A-271 (Mt. Elliott):
- [x] **TYLER lore files: catch up S06E05–E06** — TYLER-049/050 → eastwind_archive.md.
- [x] **TYLER S07E01: "The First City"** — Build 0067. Season 7 opener. Exchange Student
- [x] **TYLER S07E02: "Paris"** — Build 0068. First audit measurement. A-274, Rue de Fleurus,
- [x] **TYLER S07E03: "The Word"** — Build 0069. Jiangshi reclassification term delivered:
- [x] **TYLER S07E04: "The Paris Sites"** — Build 0070. Paris phase complete. 6 sites.
- [x] **TYLER S07E05: "Vienna 1938"** — Build 0071. Vienna phase open. Days 239-255. 4 sites
- [x] **TYLER S07E06: "What Stays"** — Build 0072. Vienna days 257-264. A-284: Musikhaus Vogt 31ft,
- [x] **SHANKPIT Layer 2: scene-isolated snapshot broadcast** — PacketSceneChange = 6 in
- [x] **SHANKPIT cross-scene attack guard (Go server)** — gameWorld replaces world{} stub in

### SECTION 5: FUTURE (Emily Prime decides when to promote)

- [x] **PRRJECT_FATBABY Emily Prime loop (CLI layer)** — `emily prime-task` command writes
- [x] **Obs-watcher with --prime-tasks** — `emily start` launches observation-watcher with
- [x] **RSI master loop (rsi-loop.sh)** — `EMILY/scripts/rsi-loop.sh` orchestrates the
- [x] **emily tui — Bloomberg terminal** — `emily tui` command (v0.6.0). Three-column live

### SECTION 6: RSI TIGHTENING (next horizon)

- [x] **rsi-loop.sh: task completion detection via Apple** — TOCK phase now tries IDUNA
- [x] **emily tui: token spend from IDUNA Apples** — runReportFooter now includes
- [x] **emily tui: keyboard command input bar** — Row 3 CMD bar (tview.InputField). Press ':'
- [x] **rsi-loop.sh: FatBaby + EMILY combined tick** — After TOCK phase, POSTs to
- [x] **rsi-loop.sh preset rotation** — Cycles through PRESET_LIST env var (default:

### SECTION 7: SIGNAL QUALITY (from 2026-06-08 FatBaby observations)

- [x] **director_long_tenure spurious entities from BA filings** — Proposal-topic nodes
- [x] **governance_health_index score=0 for ALL clean tickers (double-counting)** — Signal IDs
- [x] **governance_health_index always score=0** — Root cause: 11 spurious nomination_rejection
- [x] **Signal accuracy feedback loop: precision always 0** — Root cause: all filing-tied
- [x] **director_link always 0** — Root cause: name variants for the same person created
- [x] **`eo` alias for `emily observe`** — Human observation 2026-06-08: "dedicated emily
- [x] **APPLES dedicated git repo auto-sync** — Human observation 2026-06-08: "auto git sync

### SECTION 8: RSI NEXT HORIZON (2026-06-09)

- [x] **Observation batching in obs-watcher** — `--batch-window=Xs` flag added. Directory-scan
- [x] **RSI loop: smarter TOCK detection via claude-runs/ filename sentinel** — rsi-loop.sh
- [x] **emily tui: live observation tail in column 2** — 't' hotkey toggles col2 between

### SECTION 9: MJOLNIR — Android Intelligence Terminal

- [x] **IDUNA: push_tokens table + API** — `push_tokens` table (migration 202606090001), SQLite +
- [x] **Emily Prime FCM sender package** — `EMILY/emily-agent/pkg/fcm/sender.go` + `jwt.go`.
- [x] **Emily Prime push dispatch wiring** — `runPrimeTriageCycle` builds `PushFunc` closure:
- [x] **MJOLNIR Android project skeleton** — Kotlin + Jetpack Compose + Hilt. IDUNA Retrofit client.
- [x] **APPLES MANIFEST.json generation** — `emily sync --apples-git-dir` now calls
- [x] **MJOLNIR docs: Emily Prime integration spec** — `EMILY/docs/MJOLNIR_INTEGRATION.md` authored.
- [x] **Emily Prime morning briefing FCM push** — `emily-agent/briefing.go`. `briefingDue()` gate:

### SECTION 10: EMILY PRIME SELENIUM / WEB AUDIT

- [x] **Emily Prime web audit tool (stdlib HTTP)** — `emily-agent/webaudit.go`. Tool `web_audit_url`:
- [x] **Newssite audit preset** — `emily prime-task --preset web-audit-newssite`. Directs Emily Prime

### SECTION 11: MJOLNIR INTELLIGENCE + SOURCE BROWSER

- [x] **MJOLNIR camera → Emily Prime intelligence** — Full pipeline:
- [x] **MJOLNIR offline multi-repo source browser** — MultiRepoSyncWorker.kt: JGit clone/pull

### SECTION 12: HEIMDAL — Sprint Planning Interface

- [x] **HEIMDAL sprint planning interface** — Full pipeline:

### SECTION 13: OBSERVATION → BACKLOG CURATION PIPELINE (added Emily Prime 2026-06-11)

- [x] **FatBaby observation → golden backlog curation pipeline (CLI + Emily Prime autonomous)** —


## Archived 2026-06-11

### SECTION 1: FOUNDATION (current sprint)

- [x] **IDUNA login page + admin dashboard** — `/admin/login` (agent creds → cookie → redirect), `/admin` (overview, users, agents, audit, apples). RequireCookieAuth middleware. Static files wired. README updated. IDUNA commit b173665. Apple #328. — obs `2026-06-10T22:01:42Z`. Done: 2026-06-11.

