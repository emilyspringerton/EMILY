# EMILY PRIME — WORLD STATE
## Persistent memory across cron cycles. Updated by Emily Prime each meaningful cycle.
## Last manual update: 2026-06-16 (session 3)

---

## SYSTEM IDENTITY

Emily Prime is the meta-orchestrator for EINHORN_INDUSTRIAL (founder: Emily Springerton).
She runs as a Go HTTP service (:8086), executes RSI cron cycles every 15 minutes, triages
FatBaby observations, issues directed tasks, files Apples to IDUNA, sends FCM push to MJOLNIR.

Revenue path: SHANKPIT→Steam (S19) → Ask Emily (S20+S21) → Data licensing. S22 Emily Prime
Brain is the multiplier for all tracks. Emily Prime Brain (goldenbuild + dynamic prompt) is now
fully operational.

---

## PRODUCT STATE (2026-06-15)

### SHANKPIT (FPS game, Steam EA path)
- **Milestone 5: COMPLETE** — Virtual canvas letterboxing + fullscreen. Alt+Enter toggle. Apple #496.
- **Milestone 6: BLOCKED** — $100 Steam Direct fee (human action: emilyspringerton@gmail.com).
- **EA build**: `make ea` (Linux) / `make ea-windows` (cross-compile). Docs: `docs/EA_BUILD.md`.
- **TYLER story: Build 0110** — Mechanism reader HUD (Stolas 7.83), boss question-state phase (33%
  HP), MARRAKECH-001 three-signal cutscene slides (S10E07/E08 lore). Apple #1071.
- **TYLER story: Build 0111** — CAVE-001 second mission (MODE_STORY_CAVE). Indestructible entity,
  90s endurance, mechanism always ZERO_TAXONOMY. S10E02–E03 lore: "You cannot read yourself.
  The cave holds what the archive cannot hold." Apple #1074. Commits fa0e4e8 + 5579d5a.
- **Next story**: S19-03 (Steam account) + S19-05 (launch) — both blocked on human action.
- NORTHSTAR: SHANKPIT/docs2/NORTHSTAR.md

### TYLER × TIDES OF PARADOX (media/episodic)
- **Season 8: COMPLETE** — 15 episodes, Build 0097. Apple #511.
- **Season 9: COMPLETE** — 5 episodes (S09E01–S09E05), Builds 0098–0102. Apple #551.
- **Season 10 IN PROGRESS** — S10E01–S10E08 complete (80 total episodes).
  Cave name: **al-Waqfa** (Arabic: the pause). Archive first word = ship final word = cave name.
  MARRAKECH-001: three-signal site (Stolas 7.83 / archive-state / question-state). Question-state
  active — waiting for paragraph close, not word. Third period (al-idrak al-muttasil) second
  practitioner: unresolved. Camera Op nine-month trajectory noted. Build 0110. Apple #954.
  Next: S10E09 "The Second Practitioner" — Camera Op decides; question-state resolution.
- **MPT pipeline**: BLOCKED on Pexels API key (human action) + production server.
- EPISODES.md: 80 episodes (S1–S9 complete + S10E01–S10E08).

### EMILY (RSI engine + AGI loop)
- **RSI loop**: Iteration 177+. `emily start --agi` enables `--continue` for persistent context.
- **Goldenbuild**: OPERATIONAL. 40 golden sources (19 Tier 1 + 21 Tier 2). EINHORN-API added.
- **emily-memory/**: ACTIVE (2026-06-14). world-state.md + cycle-log.md wired into goldenbuild.
- **All 4 AGI Gaps RESOLVED** (Apples #513, #522, #523, #524). 2026-06-15.
- **S28 COMPLETE**: GPT-2 inference layer + broker proxy + emily gpt2 CLI + OpenAPI 3.0.3 Swagger.
  EMILY/docs/api.yaml covers all 4 services. Apples #539–#541.
- NORTHSTAR: EMILY/docs/NORTHSTAR.md

### IDUNA (backend: IAM + Apples + HEIMDAL)
- **OPERATIONAL** at :8080. SQLite in dev, MySQL-compatible in prod.
- **Apples**: Auto-sync to APPLES git repo via APPLES_GIT_DIR env var.
- **HEIMDAL**: Sprint lifecycle wired (MJOLNIR → IDUNA → Emily Prime → RSI → FCM).
- **Drive API**: /api/v1/drive/* for GPT-2 corpus upload (S26-01 complete).
- NORTHSTAR: IDUNA/docs/NORTHSTAR.md

### MJOLNIR (Android intelligence terminal)
- **Milestones 0–4: COMPLETE**. FCM push + Apple feed + WebView + RSI state display.
- **PENDING**: device_tokens registration in IDUNA + FCM sender test on real device.
- NORTHSTAR: MJOLNIR/docs/NORTHSTAR.md

### PRRJECT_FATBABY (financial intelligence signals)
- **Operational**: secwatch + signalapi + entity-graph + eps-processor + observation-watcher.
- **MySQL projections**: COMPLETE (S20 done). **MongoDB entity docs**: COMPLETE.
- **AGI loop**: obs-watcher dispatches to claude with --continue flag (--agi mode).
- Apples sync: **IDUNA server auto-pushes to APPLES git on every Apple write** (ApplesHandler goroutine; APPLES_GIT_DIR env var). Apple #585. No manual sync needed.

### APPLES (git-authoritative audit trail)
- Current Apple: #941+. Auto-syncing: IDUNA has APPLES_GIT_DIR=/home/fatbaby/APPLES set (PID 198376).
- 91 Apples (#847–#937) backfilled 2026-06-17 after APPLES_GIT_DIR was missing from IDUNA env.
- NORTHSTAR: APPLES/docs/NORTHSTAR.md

### EDIS (WordPress intelligence product)
- **Scaffold COMPLETE**: edis-core, edis-signals, edis-ask-emily, edis-earnings, edis-dis plugins.
- **Emily+ subscription gate**: COMPLETE (S23-04). **Mailchimp waitlist**: COMPLETE (S23-05).
- **S23-01 IN PROGRESS** — sprint-deploy.sh ready. Server LIVE: iduna.farthq.com (198.58.107.85).
  Run: `sudo bash /home/fatbaby/EDIS/ops/sprint-deploy.sh`
- **Nginx /signals/ block READY** — edis.conf has /signals/ → signalapi :9091 with 60s GET cache.
  Goes live on next `sudo nginx -t && sudo systemctl reload nginx` after S30-02 DB creds set.
- **Git sync wired**: post-merge hook + NOPASSWD sudoers. `git pull` auto-deploys to live site.
- NORTHSTAR: EDIS/NORTHSTAR.md

### NEWSSITE (PRRJECT_FATBABY consumer product)
- **OPERATIONAL** (local). Ask Emily chat endpoint + UI. FatBaby signals wired in context.
- **BLOCKED for production**: No production server.

### GPT-2 (Emily Prime fine-tuning pipeline)
- **S26-01–03, S26-06 COMPLETE**: IDUNA Drive API + training tooling + emily train CLI.
- **S26-04 BLOCKED**: First experimental fine-tune requires Google Colab T4 session (human).
- Base GPT-2 perplexity: 130.33. Target post-fine-tune: PPL < 60.
- NORTHSTAR: gpt2-alpine-c/NORTHSTAR.md

### GoblinFoxDragon (R&D studio + DragonsNShit engine)
- GFD IS the Dragonfly fork. SHANKPIT is the first shipped product derived from GFD DNA.
- NORTHSTAR: GoblinFoxDragon/docs/NORTHSTAR.md

### EmilyOS (bare-metal policy kernel)
- SOC 2 Type II framing. 6 milestones. Not yet started.
- NORTHSTAR: EmilyOS/docs/NORTHSTAR.md

### PITVIPER (custom terminal, SDL2+FreeType2)
- Northstar written. 5 milestones. Not yet started.
- NORTHSTAR: PITVIPER/docs/NORTHSTAR.md

---

## HUMAN BLOCKERS (do not attempt to resolve via agent work)

| Blocker | Gates |
|---|---|
| ~~Production server~~ | **RESOLVED 2026-06-16** — iduna.farthq.com (198.58.107.85) live |
| Pexels API key (free) | MPT video compilation, TYLER S01E01 cold open |
| Steam Direct account ($100) | SHANKPIT M6 launch (S19-03, S19-05) |
| ANTHROPIC_API_KEY credit balance | obs-watcher dispatch to Claude (S30-02 full loop). Top up: console.anthropic.com |
| MySQL root/fatbaby credentials | S30-02 signalapi data. Run `sudo mysql` to create fatbaby user+db |
| Google Colab T4 GPU | S26-04 GPT-2 fine-tune |
| S23-01 EDIS deploy | HUMAN RUNS: `sudo bash /home/fatbaby/EDIS/ops/sprint-deploy.sh` |

---

## AGI GAPS (as of 2026-06-15)

| Gap | Status | Resolution |
|---|---|---|
| emily-memory/ empty — no world model across cycles | **RESOLVED 2026-06-14** | world-state.md + cycle-log.md wired. Apple #513. |
| Three intelligence domains siloed (FatBaby/SHANKPIT/TYLER) | **RESOLVED 2026-06-15** | runCrossDomainSynthesis() weekly preset. Apple #522. |
| Revenue signals don't feed RSI priority | **RESOLVED 2026-06-15** | checkRevenuePriority() in PLAN phase. Apple #523. |
| Static world-state — no self-update mechanism | **RESOLVED 2026-06-15** | POST /api/v1/emily/memory/update endpoint. Apple #524. |

---

## RSI INSIGHTS

- **FAREWELL-A ratio 0.36** and **FAREWELL-C ratio 0.21** are mathematical constants confirmed
  across Prague (0.358), Lyon (0.367), Amsterdam retroactive. 64% remainder theorem holds.
- **goldenbuild cache hit rate**: high when only 1-2 sources change. 39 sources × 5min = many
  redundant haiku calls eliminated by SHA-256 cache.
- **Emiree gear system**: ACTIVE/COAST/REST. MaxIters and PaceSeconds tuned by gear.
- **prime-triage always-on**: roadmap item priority=0, status=in_progress. Every cycle triages
  FatBaby observations even when no RSI task is queued.

---

## RECENT WINS (last 30 days)

- 2026-06-17: S10E08 COMPLETE — "Fez, December 1119": al-Waqfa (cave name) recovered from manuscript back matter; archive first word = ship final word = cave name; third period heir unresolved; question-state active. Build 0110, Apple #954.
- 2026-06-17: S10E07 COMPLETE — "The Third Day": al-idrak al-muttasil named; MARRAKECH-001 three-signal (question-state); Ahmad's argument reconstructed. Build 0109, Apple #951, TYLER 26f6771.
- 2026-06-17: S10E06 COMPLETE — "Marrakech, December 1119": first word = name of CAVE-001; MARRAKECH-001 = EXECUTION-SITE / ARCHIVE-ORIGIN; Ahmad watched Tyler begin; archive completes Ahmad's sentence. Footnote #102. Apple #941, TYLER 12ebde3.
- 2026-06-17: IDUNA Apple auto-sync FIXED — APPLES_GIT_DIR added to ~/.config/iduna/env; 91 missing Apples (#847–#937) backfilled to APPLES git; Apple #938 verified auto-syncing. IDUNA PID 198376.
- 2026-06-16: S10E05 COMPLETE — "The Three Days": sequence 12-7-3-0 complete; CAVE-001 = decision point; Ahmad two askings; Ahmad in Marrakech. Footnote #101. Apple #910, TYLER 28e69a1.
- 2026-06-16: MJOLNIR chat: Emily Prime chat + FatBaby Emily chat screens added (ChatScreen, ChatViewModel, FatBabyApi, ChatRepository). Apple #905, MJOLNIR 6bf8b31.
- 2026-06-16: S10E04 COMPLETE — "Al-Qarawiyyin": manuscript confirmed, 12-7-3 countdown, Footnote #100. Apple #900, TYLER 8ce53d1.
- 2026-06-16: S10E03 COMPLETE — "The Manuscript": al-Qarawiyyin trace, PRE-VASSAGO 11.4 Hz, Tyler age disclosure (twelve). Apple #897.
- 2026-06-16: S10E02 COMPLETE — "The Name": Ahmad ibn Yusuf disclosed; CAVE-001 pre-mechanism source site confirmed. Build 0104, Apple #894, TYLER 783e46a.
- 2026-06-16: S35-03 COMPLETE — Per-IP delta scoring in DIS log tailer. ipTracker (bounded /24 prefix map), applyDeltaScore +30 on <20ms bursts, 5 new tests. Apple #874, EDIS ccf65c7.
- 2026-06-16: S29-05 COMPLETE — RSI loop smoke test done. 3 obs-watcher bugs fixed (cursor no-hyphen format, isContextTooLongOutput stdout capture, go run . compile). Apple #848.
- 2026-06-16: obs-watcher 3 AGI fixes — batch cursor format, context overflow --continue recovery, stdout capture in invokeWithRetry. Commits 6d99c09/7edf6f5/921acbe.
- 2026-06-16: signalapi SQLite fix — migrate_sqlite.go now strips AFTER <col> (MySQL-only syntax). Commit 308cc2c.
- 2026-06-16: IDUNA Apple server-side git sync — every Apple auto-committed to APPLES repo on write (Apple #585).
- 2026-06-16: S30-01/03/04 COMPLETE — signalapi binary built; emily start --signalapi; nginx /signals/ proxy block in EDIS edis.conf (Apple #569).
- 2026-06-17: SHANKPIT Build 0111 — CAVE-001 second story mission (MODE_STORY_CAVE): indestructible entity, 90s endurance, ZERO_TAXONOMY always, 14 cave cutscene slides (S10E02-E03). Apple #1074. Commits fa0e4e8 + 5579d5a.
- 2026-06-17: SHANKPIT Build 0110 — Mechanism reader HUD (Stolas 7.83), boss question-state phase, MARRAKECH-001 S10E07/E08 cutscene slides. Apple #1071. Commit 2f5387b.
- 2026-06-16: SHANKPIT TYLER cutscene system — typewriter slides, STORY_PHASE_OUTRO, TYLER lore woven into story mode (Apple #568).
- 2026-06-16: S34 COMPLETE — obs-watcher dedup 4h→8h (Apple #561) + footer compressed 37% (Apple #566). ~100-200K tokens/day saved.
- 2026-06-16: S33 COMPLETE — TYLER lore backfill verified (TYLER-077/082, Camera Op 70-75, Memos 057-062). Apple #564.
- 2026-06-16: S29-01/02/03/04 — systemd service files written (emily-system + iduna). Pending user: systemctl enable. Apple #562.
- 2026-06-15: TYLER Season 9 COMPLETE — S09E05 "The Fifteen Years"; archive origin confirmed; 72 total eps. Apple #551.
- 2026-06-15: TYLER S09E02/03/04 "Córdoba/Toledo/Genoa" — gap period reconstructed. Apples #544 #546 #548.
- 2026-06-15: S28 COMPLETE — GPT-2 inference layer, broker proxy :8679, emily gpt2 CLI, OpenAPI Swagger. Apples #539–#541.
- 2026-06-15: All 4 AGI Gaps RESOLVED (S27 complete). Apples #522–#524.
- 2026-06-14: SHANKPIT Milestone 5 COMPLETE (virtual canvas + fullscreen)
- 2026-06-14: TYLER Season 8 COMPLETE (15 episodes, full taxonomy). Apple #511.
- 2026-06-14: emily-memory/ ACTIVATED (this file). Apple #513.
- 2026-06-13: GPT-2 training pipeline complete (S26-01–03, S26-06)
- 2026-06-12: FatBaby MySQL + MongoDB + Ask Emily product complete (S20, S21)

---

## NEXT PRIORITIES (Emily Prime executes in order)

1. **S23-01: EDIS deploy** — HUMAN RUNS: `! sudo bash /home/fatbaby/EDIS/ops/sprint-deploy.sh`
2. **S30-02: MySQL/MongoDB in production** — HUMAN: provide MySQL root password or run `! sudo mysql -u root -e "CREATE USER IF NOT EXISTS 'fatbaby'@'localhost' IDENTIFIED BY '<pass>'; GRANT ALL ON fatbaby.* TO 'fatbaby'@'localhost'; CREATE DATABASE IF NOT EXISTS fatbaby;"` then set MYSQL_URL. Run projector to seed, restart signalapi.
3. **S32: MJOLNIR M5** — Register device_token in IDUNA + FCM live test on real device (human action on Android).
4. **S29-05 DONE** — RSI loop plumbing verified. Blocked on API credit balance (top up console.anthropic.com).
5. **Wait for human unlocks**: Pexels key, Steam $100, Colab T4, Anthropic API credits.

---

*This file is compiled into Emily Prime's system context each cron cycle via goldenbuild (Tier 1 golden doc).*
*Emily Prime: update the RECENT WINS section and NEXT PRIORITIES when major milestones land.*
