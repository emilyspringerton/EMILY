## 2026-09-03

- docs/THE_EMILY_WAY.md: new Principle 19 ("A Big, Unscoped Ask Gets Scoped, Not Swallowed Whole"), kanban priority-queue card GOLDENOPS-001 — the standing default for a big, unscoped ask on the priority/cruise queue: investigate, write a real phased NORTHSTAR, register it as a golden doc, then return the real, resulting sub-tasks to BACKLOG.md as their own section/sub-items rather than marking the whole kanban card fully done the moment a plan exists. Explicitly does not change how a genuinely "plan/scope this" card closes (that's a real, complete deliverable on its own, per Principle 3). Apple #17524. (sess-20260902-2008-ed50169e)

## 2026-08-28

- BACKLOG.md: logged S202-54 (PARENA construct-split Spotlight mod), Apple #16540 (sess-20260825-1938-f6bd411e)

## 2026-08-25
- S202-19: construct CI step added to all 9 gap repos (EDIS, MJOLNIR, gpt2-alpine-c, GOLDENBAND, SKULDMARK, GTA7, PARENA, REDGARDEN, WEAKNIGHT_BEDROCK_RACERS), all live-verified green in real GitHub Actions CI. Apple #16033. (sess-20260825-1938-f6bd411e)
- GoldenDocCompiler rewritten pure-CLI (no LLM, no ANTHROPIC_API_KEY): deterministic header+lead-line extractive summary replaces claude-haiku compression, closing the HITL-11 stale-cache side-finding structurally (no more real-vs-degraded distinction to lose track of). emily context build's one-shot path live-verified: 47/47 sources compile with the key unset. Daemon-side restart queued (sudo-queue/28, needs fatbaby -- treeiii can't signal fatbaby's PID). Commits EMILY 1822e22, emily.cli 79bc797+b60ca6e. (sess-20260825-1938-f6bd411e)
- reboot-recovery session: confirmed S176-34's Google OAuth client-ID blocker is a genuine gcloud-CLI limitation (IAP oauth-brands API deprecated/wrong-tool, never enabled on einhorn-mjolnir); clientg_id.tct's client ID traced to an unrelated GCP project (278374120873, not einhorn-mjolnir) and not used (sess-20260825-1938-f6bd411e)
- SECTION 191: security audit + IDUNA admin session logout fix — read-only sweep found no signs of intrusion (permission-limited on auth.log/fail2ban, consistent with existing no-sudo constraint), root-caused and fixed the admin session's silent 8h hard-cutoff logout (S191-02, IDUNA Apple #15775), PARENA intrusion-detection tooling deferred to S190-01's already-tracked gap (S191-03). (sess-20260825-0828-cc32a704)
- S139-02: toothbrush manufacturing vendor research (STINKIES COMMISSAIRE + ULTRA) — Tier 1 well-served by existing OEM/ODM toothbrush makers, Tier 2 (brass ULTRA) found to be a real 3-vendor/3-step supply chain (handle CNC/PVD shop + head-mold OEM + assembly), not a single vendor pick. Apple #15767. (sess-20260825-0828-cc32a704)

- Wrote up a research note on an external AI's hallucinated 'EmilyOS'/'Smilium' architecture claims, tracing what's real (public repo names, a real Erlang idea) vs fabricated. docs/EXTERNAL_AI_HALLUCINATION_CASE_STUDY_EMILYOS.md, commit 166f8a5. (sess-20260824-2252-ce890e4f)

## 2026-08-20
- golden-docs-index 新增 LADYBUG-NORTH(ladybug/README.md)——PARENA 的官方 BDD 測試框架獨立發布為新 repo (sess-20260820-0649-a3f19d93)
- 發布反思文章「Nine Hours In, The Emily Way」——PARENA 從 VS0 domain 1 落地(11:13 UTC)到現在約 9 小時的真實建置歷程,透過 The Emily Way 操作哲學(backlog-first、Apple-before-done、誠實缺口清單)檢視,提及 public domain(Unlicense)授權 (sess-20260820-0649-a3f19d93)
- 發布 Tyler 更正文章「Correction: It's SAND」——PARENA 原生 IDE 最終命名確定為 SAND(S-expression And Not Dumbshit),取代發布不到一小時前的 JEWEL;誠實記錄命名演變過程(IRONCLAD→GPT-2 嘗試→JEWEL→SAND) (sess-20260820-0649-a3f19d93)
- 發布 Tyler 命名儀式部落格文章「Naming JEWEL」——PARENA 原生 IDE 正式命名為 JEWEL(jewel beetle,呼應 PITVIPER 既有的 F11 shiny font 玩笑),誠實記錄 IRONCLAD 被否決 + GPT-2 推論真實結果(4 次生成都是不相關的技術雜訊,含一次把 PITVIPER 誤植成 POTVIPER-INCLUDE) (sess-20260820-0649-a3f19d93)
- 發布 Tyler 視角部落格文章「What a Legendary Product Designer Sees in PARENA」,同系列 keynote-followup 格式,聚焦 PARENA 對整個生態系的產品價值(mod-surface API 哲學、今天的真實技術進展、誠實缺口清單) (sess-20260820-0649-a3f19d93)
- 發布 FATBABY_NEWSWIRE 新聞稿「PARENA Eats the Codebase From the Outside In」到 OKEMILY blog,涵蓋 PARENA compiler/stdlib/哲學/Bazel/後續規劃 (sess-20260820-0649-a3f19d93)
- Added `scripts/gmail_imap_fetch.py` — on-demand IMAP read tool using the SAME
  `GMAIL_SMTP_ADDRESS`/`GMAIL_SMTP_PASSWORD` app-password credential as the SMTP send path (app
  passwords aren't scoped per-protocol). Founder, real-time: confirmed they want to use their app
  password to get email read access specifically to retrieve the language-spec email they sent to
  emilyspringerton@gmail.com. Deliberately a plain script, not a service — matches the founder's
  own stated scope: "for now it needs to be treated as a queryable source not something you are
  totally chunking through yet." Honors the explicit no-auto-logging-email-content boundary from
  the same conversation: prints to stdout for the operator, writes nothing to Apples/BACKLOG/
  CHANGELOG. Stdlib only (`imaplib`/`email`), no dependencies. Syntax-verified
  (`python3 -m py_compile`) but NOT yet run against a real mailbox — the founder hadn't set the
  credential yet as of this commit (`emily key set GMAIL_SMTP_ADDRESS`/`GMAIL_SMTP_PASSWORD`).
  (sess-20260813-2154-dda37e8b)
- feat(emily-agent): added an SMTP App Password fallback path to the Gmail integration
  (`emily-agent/gmail.go`). Founder: "we need to get the email integration fixed" — the existing
  OAuth2 refresh-token flow (Path 1) needs a one-time interactive browser consent flow only the
  founder can do; the founder offered a Google Account App Password instead. New send-only Path 2
  (`GMAIL_SMTP_ADDRESS`/`GMAIL_SMTP_PASSWORD`, real `net/smtp` + STARTTLS to `smtp.gmail.com:587`,
  no third-party dependency) — `buildGmailClient()` tries Path 1 first, falls back to Path 2.
  Credentials go through the *existing* generic `emily key set <NAME> <VALUE>` command (no new
  CLI subcommand needed — the founder's own ask for a fast-tracked key command was already
  satisfied by prior infrastructure). `ReadInbox` still requires Path 1 (real Gmail API) and now
  fails with a clear error instead of a confusing one if only Path 2 is configured. Explicit
  founder privacy boundary honored: no code path here logs email content into Apples/CHANGELOG/
  BACKLOG. `go test ./...` (emily-agent module) passes clean. (sess-20260813-2154-dda37e8b)

## 2026-08-18
- Published blog post 'Tuxedo Duck Is Not Duck Tuxedo' on the semantics/ontology thread, live at okemily.com/blog/tuxedo-duck-is-not-duck-tuxedo/ (sess-20260813-2154-dda37e8b)

- NORTHSTAR_PROMPT_O_VERSE.md §9: mashup/hybrid discovery requirements, deferred pending semantic-judgment mechanism (sess-20260813-2154-dda37e8b)

## 2026-08-17

- docs(blog): 'The Backoff That Remembers It Was A Different Process' -- write-up of the promptoverse adaptive cross-invocation backoff pattern (persisted failure state consulted preemptively before a stateless CLI's first request of a new run, not just mid-run retries) (sess-20260813-2154-dda37e8b)

## 2026-08-15

- emily-agent/integration.go: runBacklogCuration截斷summary改用rune-safe([]rune),修復真正的無效UTF-8根因(不是emily.cli,是RSI cron自己獨立實作的一份重複邏輯)。新增regression test。commit 1a00e2e。 (sess-20260813-2154-dda37e8b)

## 2026-08-14
- goldenbuild: fallback entries(500字元截斷版)從沒被自動重試的bug修復,71%的golden context來源受影響。IsFallback欄位+Build/MaybeRebuild兩處判斷修正。commit 70b1857。 (sess-20260813-2154-dda37e8b)
- S175-04 完成:開放資料釋出 teaser 部落格文章已發布(繁中/梵文/英譯三語),https://okemily.com/blog/the-thread-not-yet-cut/ (sess-20260813-2154-dda37e8b)
- docs/NORTHSTAR.md: added Claire section (Tier 2 golden doc, excluded from GoldenDocCompiler by tier filter, real emily claire CLI + claire-log.md audit trail); Key Files table updated. (sess-20260813-2154-dda37e8b)

- HITL-11 re-confirmed dead 2026-08-14; corrected 19 misleading SECTION 5 backlog entries pointing at a missing-env-var explanation that was false (key is set, just out of credit). (sess-20260813-2154-dda37e8b)

## 2026-08-13

- REBOOT_RUNBOOK.md updated: gpt2-serve + fatbaby-broker now real systemd units, not manual-restart; documented broker's hardcoded 30s timeout vs serve.py's real multi-minute cold-inference latency (sess-20260813-2154-dda37e8b)

## 2026-08-11

- Pizza 部落格文章發布：'Nobody Checked'，涵蓋本次 session 的 REDGARDEN/GTA7 更新 (sess-20260810-0505-a53abca2)

## 2026-08-10

- CLAUDE.md 規範宣導：session tag commit 規則推廣到全 monorepo 21 個 sub-repo + APPLES + root (sess-20260810-0505-a53abca2)

## 2026-08-02
- fix(emily-agent): RSI loop was re-issuing the same ~9 directed tasks into `signals/tasks/`
  forever. Founder: "the rsi loop keeps putting the same stale 7ish tasks into the backlog."
  Root cause: `runPrimeTriage` re-scores the same `ReadObservations(10)` window every cycle with
  no cursor tracking which observations had already been triaged for tasking -- `WriteTask`'s
  own `recentDuplicateExists` (a rolling 4h content-match window) only ever rate-limited this to
  once per ~4h12m (window expires -> next 15-min poll finds no recent match -> writes a fresh
  duplicate -> clock resets), not stopped it. 260 duplicate copies of the same 9 tasks had
  accumulated in `signals/tasks/` over roughly 45 days. Fixed with a `.task-cursor` file mirroring
  the existing (correctly-working) `.escalation-cursor` pattern -- once an observation's been
  triaged for tasking, it's never re-decided. 2 new regression tests, one of which specifically
  backdates a task file past the 4h window to prove the fix (not the pre-existing dedup layer) is
  what prevents the duplicate; both confirmed to fail without the fix, pass with it.
  `emily-agent/integration.go`.

## 2026-07-25
- OKEMILY: rebranded redgarden.html + redgarden-wishlist.html to Knights of the Void per real-time founder direction ('update the redgarden landing page to be knights of the void', 'current status download from artifacts on github instructions mailing list for knights of the void wishlist on steam'). New download section with real GitHub Actions artifact instructions checked against ci.yml. Deployed live. OKEMILY c9922d6, Apple #10737. (sess-20260723-2347-df115bd5)
- backlog: ✓ S170-96 -- REDGARDEN hero name labels above floating health bars. arena_ai_bridge.c wasn't linked into the arena client build at all -- fixed. REDGARDEN e53ee5f, Apple #10714.
- backlog: ✓ S170-90 -- REDGARDEN bots bunching up, fixed. Move-target logic sent the nearest enemy's exact coords instead of a spread approach angle. Verified live with a real 20/20 match. REDGARDEN b22ee89, Apple #10712.
- backlog: ✓ S170-76 -- REDGARDEN 'Knights of the Void' naming, found already complete (blog post + window title rename both done, item was just left open). Closed with evidence.

- backlog: ✓ S170-87 -- REDGARDEN capture nodes rendering compressed onto one point, fixed. ArenaSnapshotMsg never included node data; added ArenaNodeSnapshot to the wire format + owner-based node coloring. Wire-format change, restarted all three live systemd units on the new build, verified clean. REDGARDEN 36f868e, Apple #10708.

## 2026-07-24
- backlog: ✓ S170-50/51 -- REDGARDEN Arathi Basin channel capture (redesigned from S170-46) + territorial jungle creeps + bot names. Real WoW Arathi Basin rules built exactly as specified: stealth capture, stealth-breaks-on-channel-start, damage-interrupts-capture. Jungle creeps tie territory control to combat rewards, countering turtling comps. Activated real WOTAN bot stats tracking (was silently disabled all session) + memorable bot names. 28 new tests (251 total), verified live. REDGARDEN 2cf6cdd, Apple #10654.
- Verified TYLER's IDUNA agent registration was already complete (stale backlog premise) -- closed SECTION 5 item and HITL-09 with DB evidence rather than building a redundant CLI. IDUNA Apple #10648.

- backlog: ✓ S170-49 — REDGARDEN The Donkey, Paper Glide ability (docs only). Founder direction: launch airborne, fold into a paper airplane, glide clear of danger, flying over terrain and immune to ground-based CC while airborne. Added as Q in `docs/HEROES_VS0.md`, consistent with the existing Indirect-Control identity (auto-triggered, not player-cast). Docs only -- The Donkey stays blocked on a non-piloted-unit system that doesn't exist in the sim yet, flagged explicitly. REDGARDEN `5c2a54c`, Apple #10646.
- backlog: ✓ S170-48 — REDGARDEN The Courier (Ratatoskr), eleventh hero, roster 10 → 11. TYLER #32's messenger-between-two-fixed-points framing (eagle at Yggdrasil's crown, Nidhogg at its root) maps directly onto the arena's two existing `ArenaNode` positions -- W is a pure fixed-geography teleport, distinct from every other hero's ally/foe-relative teleport. Q is a dash-strike whose landed cast also cleanses The Courier's own debuffs; R is a flat life-drain execute. 7 new tests (223 total), verified live across a real 22-bot pool after cleaning up a stray port conflict. REDGARDEN `d01eef8`, Apple #10645.
- backlog: ✓ S170-46/S170-47 — REDGARDEN territory/node system + five new heroes (Tree, Pizza, Flamel absorbing Druid, Morrigan, Dagda). Founder-picked territory/resource economy unblocks the most heroes at once; `arena_tick_nodes` turns the two decorative `ArenaNode` markers into a real capture contest. Mid-build redirect merged Druid into Flamel (TYLER lore confirmed Druid had no named-character backing); two more founder additions (Morrigan "meta jungler," Dagda "two-natured hammer") both real TYLER entries (#68/#69). Roster 5 → 10, 62 new tests (216 total), verified live with the persistent bot pool relaunched on the freshest build and left running. REDGARDEN `acdd8ce`, Apple #10644.
- backlog: ✓ S170-14 (3/3, now fully closed) — ranked matchmaking design pass. Plain ELO over Glicko/TrueSkill (this pool is 1v1-only and symmetric, doesn't need their uncertainty modeling yet), a new `redgarden_ranked_stats` table kept separate from casual stats, and a widening-search-window queue design scoped as its own future build pass. `docs/RANKED_MATCHMAKING.md`, golden-indexed as REDGARDEN-RANKED. All three matchmaking pools the founder asked for are now either built+verified or properly designed. REDGARDEN `8b52e54`, Apple #10640.
- backlog: ✓ S148-00 — wired `/chat` to the Apple ledger. Per-turn, new `conversation` Apple type (added to Back Office filter dropdown too, IDUNA `d2b9b6e`). `buildChatTurnApple` (pure payload builder) + `postChatTurnApple` (non-fatal PostApple wrapper), same pattern as `buildCycleApple`. 5 new tests. EMILY `b76f254`, Apple #10638.
- Published 'RED GARDEN: Allies, and the Wall Worth Stopping For' — the roster-audit wall (every remaining hero blocked on a missing system), the founder's build-allies decision, `arena_nearest_ally` unblocking Ghost/Frog/Doc Wheel, and the honest bug-in-my-own-test story. https://okemily.com/blog/red-garden-allies-and-doc-wheel/ — Apple #10637.
- backlog: ✓ S170-45 — REDGARDEN arena allies + Doc Wheel (fifth hero, first ally-only kit). Founder decision: build allies over territory or declaring the roster complete. Added `arena_nearest_ally`, unblocked Ghost's ally-heal + Frog's Borrowed Time, shipped Doc Wheel's full kit (heal/cleanse, teleport, teamwide heal — shield simplified to heal, flagged not faked). Found and fixed a real test bug (an assertion passing by coincidence, not correctness). 16 new tests, verified live across two real bot matches. REDGARDEN `df58b13`/`c6f9f8c`, Apple #10635.
- backlog: S170-14 (2/3) — REDGARDEN player-only matchmaking pool verified live. `scripts/launch_arena_pools.sh` runs a second, separate matchmaker instance (port 7779, 1v1, zero bots ever pointed at it); two real human `--queue` clients matched into a genuine 1v1, cross-checked clean isolation both directions. Bot games (1/3) was already S170-43. Ranked (3/3) stays explicitly undesigned. REDGARDEN `c31e90e`, Apple #10628.
- fix(naming): resolved the FABLE naming collision flagged by the 2026-07-18 SAGA audit before starting SECTION 145 work — `emily-agent/fable.go` (Emily Prime's claude-haiku backlog advisor) renamed to `mimir.go`/MIMIR (Norse: the wise counsel-giver, matching NORN/SAGA/FATES), routes `/api/v1/emily/fable/*` → `/api/v1/emily/mimir/*`. HQ-SPEC-AI-103's sovereign model-line keeps the FABLE name. `docs/NORTHSTAR.md`/`docs/API_KEY_UNLOCK.md` updated to match. Found a second, bigger collision while doing this — this environment's real Claude model lineup includes an actual model called Fable (`claude-fable-5`) — flagged in BACKLOG.md SECTION 145 for the founder, not resolved unilaterally (renaming an entire spec's branding is a product decision, not a mechanical fix).
- backlog: ✓ S144-01 — GOLDEN BAND v0 (HQ-SPEC-SIM-100 §8 build step 1). New standalone repo `GOLDENBAND`: `.gband` binary format + JSON manifest, ~90-line C sampler, self-contained sha256, Go `gbtool` (BVH import/hash/validate). 13 tests, full pipeline smoke-tested. glTF import + SHANKPIT integration explicitly deferred. GOLDENBAND `a9de6b1`, golden-indexed as GBAND-FORMAT, Apple #10612. One real gap: no GitHub remote yet (same as PITVIPER's S127-04), committed locally only.
- Published 'TYLER, Series X: The Long Quiet / Recruitment Drive' — teaser covering the two newest Series X interludes: an intentionally plotless infrastructure-drama episode turning the real emily-agent 2-day OOM incident into an in-universe Emily OS silence, and a RED GARDEN in-game cutscene dramatizing the roster's nine-hero gang-formation backstory. https://okemily.com/blog/tyler-series-x-the-long-quiet-and-recruitment-drive/ — Apple #10610.
- backlog: ✓ S143-01 — SAGA frontmatter schema (`docs/hq-specs/SAGA_SCHEMA.md`) + `emily saga lint` built (HQ-SPEC-DOC-102 build-sequence step 1). Retrofitted real frontmatter + inline claim citations onto all 7 live HQ-SPEC docs, DOC-102 first per its own instruction. Honest mixed reality-binding: PRIME-101's `CheckLineage`/`pkg/apples` and AI-103's entropy source marked `verified` (real running code), INFRA-105's MJOLNIR BuildConfig gap marked `diverged` (a real found gap). `emily saga lint` reports ALL CLEAN. emily.cli `92c0780`, Apple #10607.
- Published 'The Morning Report, Early' — founder asked for today's scheduled morning briefing done ahead of time as a blog post rather than waiting on the automated 09:00 UTC FCM push; same last-24h Apple window `briefing.go` would use, written longhand instead of compressed into a push body. https://okemily.com/blog/morning-report-early-2026-07-24/ — Apple #10601.

## 2026-07-23

- golden-index: GFD-HERO-FRAMEWORK added — `GoblinFoxDragon/docs2/HERO_CONTENT_FRAMEWORK.md`, the story-first process for turning a TYLER-HEROES entry into a dungeon/NM-boss/loot-drop, grounded in GFD's real mob/NM/itemdef/loot systems. Five worked examples, zero numbers yet. Also: live-played the freshly-deployed GFD MUD tonight and found + fixed two real new-player-blocking bugs (combat unreachable from spawn; a lethal worm Poison proc) — Apples #10514/#10516 (GFD), #10515 (IDUNA status-page addition), #10517 (blog post 'The First Ten Minutes'). Apple #10518 (framework doc).
- golden-index: TYLER-HEROES added — `TYLER/multiverse_heroes.md`, a 110-entry hero lore compendium for a League of Legends/Diablo II-style multiverse game (11 factions: TYLER canon + reframed real mythology + real historical figures + folklore/scripture), lore/archetype only per docs-before-software discipline. RED GARDEN cross-referenced as the closest existing game home. Apple #10511 (TYLER), #10512 (IDUNA blog post 'Ten Heroes Worth a Closer Look').
- New Principle 16 (`THE_EMILY_WAY.md`): "The Inbox — Capture Now, Triage Later." `EMILY/inbox/INBOX.md`, a zero-bar-to-entry append-only log for anything shared with no task attached — distinct from `BACKLOG.md`'s `INTAKE QUEUE` (curated, en route to becoming backlog items) since the inbox is upstream of any triage at all. Built after the founder shared a link with no clear intent and asked for exactly this. Apple #10510.
- backlog: ✓ S141-05 — append-only amendment notes added to `HQ-SPEC-SIM-100` §6 and `HQ-SPEC-FIN-099` §6, mapping each spec's own bespoke gate/approval rules onto `pkg/norn`'s Artifact/Oracle/CheckLineage per PRIME-101 §6's instantiation table. Apple #10509.
- backlog: SECTION 23/S23-01b — redispatched fable-next-backlog entry #1 (IDUNA front door funnel) alone per the 2026-07-19 postmortem's own instruction; landed same day as `IDUNA/docs/kikoryu/FRONT_DOOR_FUNNEL.md` (agents get a real lifecycle, PROPOSED→CUSTODIED→SCOPED→LIVE, but no ceremony — the owner is the accountable party) and resolves the nginx root-path collision blocking the okemily.com⇄EDIS domain merge. Also closed two stale checkboxes found while scanning (S23-01, S24-01/02 both referenced `fatbaby.io`, which never resolved — superseded by `news.okemily.com`, already live). Apple #10508.
- backlog: ✓ S158-03 — investigated (not blind-reverted) uncommitted drift on IDUNA's already-applied `202606180001_local_users.sql` migration. Every write path to `local_users.updated_at` sets the value explicitly from Go at whole-second precision, so the found `TIMESTAMP` → `TIMESTAMP(6)` edit had no functional benefit anywhere — reverted rather than completed as a new migration. IDUNA `43a930f`, Apple #10507.
- backlog: OpenClaw integration work audited on founder request ("saw some openclaw stuff mid-restart, make sure it's persisted") — confirmed nothing was lost: northstar (`docs/NORTHSTAR_OPENCLAW_INTEGRATION.md`), Apple #10285, EMILY commit `4ea6ab7`, and golden-index entry OPENCLAW-NORTH were all already fully committed and pushed on 2026-07-19. The one real loose end was S170-03, a vague bundled "both northstars' VS0 builds" placeholder — split into S170-03a (OpenClaw VS0, now a real trackable item, blocked on a founder deployment-isolation call per the northstar's own §4) and S170-03b (IDUNA Vault VS0).
- feat(watchdog): CheckCheckpointHealth — freshness alerting for FatBaby's four SQLite index checkpoints (signalapi-index.db, newssite-index.db, entity-graph's filings-index.db + accuracy-index.db). Reads meta.snapshot_at read-only, fires an escalation Apple if it hasn't advanced within 5 minutes, using the same debounce pattern as CheckServiceHealth/CheckPollerHealth. Closes SECTION 1 Phase 3 (docs/northstar/replay-fragility.md in PRRJECT_FATBABY) alongside PRRJECT_FATBABY's companion heartbeat-writing change (Apple #10504). Added modernc.org/sqlite v1.52.0 dependency (pure Go, version-pinned to match PRRJECT_FATBABY/IDUNA's existing usage). go test ./... green, 4 new tests. Apple #10505.
- backlog: closed out SECTION 1 Phase 2 / 2c (entity-graph graph-lifetime hoist) retroactively — code had shipped in PRRJECT_FATBABY d450635 (2026-07-19) without the required Apple/CHANGELOG/checkbox close-out. Apple #10503.

## 2026-07-21

- Activated the AGI continuous loop (emily start --agi) — observation-watcher now passes --continue to every claude invocation, RSI cycles accumulate context instead of starting fresh. Persisted into ~/.config/systemd/user/emily-system.service's ExecStart for future restarts. Found a gap: that unit file has no source-of-truth copy in any repo (unlike the FatBaby watcher units under ops/systemd/) — worth adding one

## 2026-07-19
- RED GARDEN cloned into the ecosystem (github.com/emilyspringerton/REDGARDEN + wiki) — elevated to core product status, northstar written capturing real-time scoping (Clash Royale card-RTS model, TrapX prototype link, multiplayer + Android/iOS/Desktop target, honor-code-compliant engagement design)
- OpenClaw integration northstar written — MIT-licensed 20+-channel self-hosted chat gateway (WhatsApp/Telegram/Slack/Discord/Signal/iMessage/etc), researched via WebSearch not guessed; scoped as a channel front-end calling Emily Prime's existing API, not a replacement for her agent loop
- FCM push to MJOLNIR on new mailing-list signups (observeMailingListSignups, OBSERVE phase) — inert until Firebase project + FCM env vars exist
- Rewrote README.md as a real technical overview (three agents, RSI cron cycle, The Emily Way, repo layout, how to run, related repos) — old README was a stale Feb 2026 brand-pivot investor memo, archived at docs/archive/brand-pivot-memo-2026-02-04.md
- Added SECTION 163 (STINKIES COMMISSAIRE VS0 hoodie funnel page) — okemily.com stinkies.html shipped (Apple #10177, OKEMILY 01c1832), live deploy + vendor selection + WooCommerce listing queued as follow-ups
- backlog: correct stale S23-01 (EDIS is actually live, HTTPS-only gap), add S23-01b domain-merge decision item, flag HITL-03 GFD theme/plugin gap and HITL-04 deploy.sh systemd-unit collision risk found during 2026-07-19 reboot session

- docs(runbook): reflect S152-03 systemd units (secwatch/prwatch/prwatch-body/processor/eps-reconciler/newssite/signalapi) + shankpit460-server.service; document eps-reconciler double-start risk found during 2026-07-19 reboot recovery

## 2026-07-18
- backlog: SECTION 1 top-priority done — replay-fragility northstar landed via Fable dispatch #6 (`PRRJECT_FATBABY/docs/northstar/replay-fragility.md`, PRRJECT_FATBABY `a628947`, golden-indexed as FATBABY-REPLAY). Decision: streaming eventstore `Scan` API + per-process SQLite snapshot-plus-tail checkpoints; the observed signalapi stall was O(n²) full-file re-reads in `FileStore.ReadFrom` against a 354MB journal, not a mystery. fable-next-backlog entry #6 moved to Done. Apple #9986.
- backlog: S153-12 done — newssite systemd-hardened + linked from okemily.com (`/news/`), fixed a broken stale drafted unit that would have failed to start; "Enter the void" CTA to farthq.com added. Apple #9961.
- backlog: S153-10 done — real ecosystem status page live at `okemily.com/status.html`, real uptime history from day one, fixed a live self-check race bug. Apple #9960.
- backlog: S153-03 done — Mailchimp live-verified end-to-end (real subscribe → encrypted store → `mailchimp_synced=1`). SECTION 153 now fully closed except S153-06 (parked) and S153-09 (spec reconciliation). Apple #9958.
- backlog: S153-08 done — API playground (Swagger UI) live at `okemily.com/api-playground.html`, fixed a real bug (spec's `servers[]` was localhost-only). S153-09 queued (spec is known-stale, missing blog/mailing-list endpoints, plus a second unreconciled `openapi.yaml`). Apple #9955.
- backlog: S153-07 done — static-HTML blog for okemily.com (no PHP/MySQL, memory-constrained box), `IDUNA/internal/blog`, commit `5bcd730`. 4 real posts published (Emily Way reflection, IAM/governance-as-moat, small-teams-security framing, Stillness-outline-sourced piece), linked from the footer. Apple #9954.
- backlog: S153-05 done — `emily key` generalized to write named secrets to any target env file (`--target emily|iduna`), emily.cli commit `6d53eba`. Apple #9952.
- backlog: SECTION 153 — okemily.com launched (landing page + never-at-rest-unencrypted mailing-list signup via `IDUNA/internal/mailinglist`). New repo `OKEMILY`. S153-01/02 done, Apple #9950; S153-03 (Mailchimp account), S153-04 (HTTPS), S153-05 (`emily key` CLI command), S153-06 (board/foundation structure, parked) still open.

## 2026-07-17
- backlog: SECTION 152 closed out — S152-01 (poller watchdog), S152-02 (prime-directive amendment), and S152-03 (systemd supervision for secwatch/processor/prwatch/prwatch-body/eps-reconciler, PRRJECT_FATBABY commits `266dad7`/`5fce9ef`) all done. Apples #9943/#9945/#9946.
- docs(process): `THE_EMILY_WAY.md` Principle 15 (Operational Health Is Not Optional) + a new Section 0 in `docs/emily-prime-directive-data-collection.md` — permanent operating principle that infra health (services + headless pollers) must be actively, automatically verified every cycle, never discovered by luck. Direct response to S152-01's incident (secwatch/eps-reconciler silently down for hours) and the founder's own framing: FatBaby's SEC/PR pipeline is a multi-year data asset ("2 years of good data"), and a silent ingestion gap cannot be recovered once it has passed. Apple #9945.
- fix(watchdog): `CheckPollerHealth` + `PollerConfig` in `emily-agent/watchdog.go` — monitors headless FatBaby pollers (secwatch, processor, prwatch, prwatch-body, eps-reconciler) via log-file freshness, since none expose an HTTP health endpoint for the existing `CheckServiceHealth` to ping. Root cause: secwatch was OOM-killed and stayed down ~10h, eps-reconciler ~22h, both with zero automated detection — only found via a manual "check on all the fatbaby data" audit prompted by an OOM-kill message seen in tmux on a phone. Reuses `WatchdogState`'s debounce/escalation pattern, keyed `poller:<name>`. Wired into the cron cycle. 4 new tests, full suite green.
- feat(S147-02): `emily-agent/enrichworker.go` — `ApplEnrichWorker`, the GPT-2 tower Apple enrichment worker, same shape as `CheckinAlertWorker`. Polls IDUNA for Apples missing `gpt2_fingerprint` (via the new `has_gpt2_fingerprint` list field), generates via `serve.py` (`:8088`), computes the fingerprint in-process via `gpt2-alpine-c/pkg/towerprint.Compute()`, PATCHes it back — async, caller-side, per `TOWERPRINT.md` §5's decision. Wired into `cron.go` alongside the alert worker. Rebuilt+restarted the live IDUNA service to pick up the list-endpoint change, started `serve.py` (base model, ~400MB RSS — light enough for this memory-constrained box, unlike training), and manually walked the full pipeline against a real Apple (#9932) end-to-end: real generation, real fingerprint, real PATCH, confirmed `has_gpt2_fingerprint` flips `false`→`true` on the list endpoint afterward. No unit tests added for the HTTP-calling logic itself — matches this codebase's existing precedent for `IdunaClient`/`CheckinAlertWorker` (neither has any), and live verification already caught what a mock-based test wouldn't have (this session's NORN `run_id` bug came from exactly this class of live check, not a unit test). `serve.py` is running but not yet systemd-supervised — flagged as a follow-up, not blocking.
- docs: `HQ-SPEC-INFRA-105-fates-dns.md` (FATES, the name layer) — DNS northstar for `farthq.com`. Positions taken: Cloudflare stays authoritative (no self-hosted NS at single-host scale; sovereignty = zone-as-code in `IDUNA/ops/dns/` as source of truth, Cloudflare as a replaceable projection), one-subdomain-per-product-surface with paths within a surface (IDUNA's drafted path split untouched; `gate.farthq.com` reserved but owned by the front-door funnel prompt), health-check-gated `dns-apply` (no `--force`), three-rung DNS-blindness self-diagnosis ladder (loopback-only inner loop → pinned-resolver probe → MJOLNIR client-side dead-man), sketched `biometric`-tier NORN row for zone diffs. Real divergence found while writing it: MJOLNIR prod BuildConfig points `EMILY_BASE_URL` at `https://iduna.farthq.com` with no nginx route to `:8086` and no cert. Golden-indexed (HQ-FATES, tier 2); build sequence at BACKLOG SECTION 151 with the Cloudflare-token/registrar-custody human unblock row.

## 2026-07-16
- docs(process): `THE_EMILY_WAY.md` gains Principle 14 (Continuity Report) — full git sync sweep + changelog digest + in-flight/blocked/queued state + next steps, stored at `EMILY/continuity/YYYY-MM-DD-HHMMSS.md`. First report caught two unpushed commits (the SAGA VS-spec audit had landed locally in IDUNA and EMILY but never reached origin).
- backlog: Act II arc — SECTIONS 141-146 (NORN loop kernel, KAREN ledger, SAGA doc curation, GOLDEN BAND animation, FABLE model line, gpt2-alpine-c corpus evolution bridge), sequenced NORN-first per the six founder-uploaded `HQ-SPEC-*.md` drafts. All six registered in `golden-docs-index.md` (`docs/hq-specs/`).
- docs: intook `iduna_roadmap.md` → `IDUNA/docs/NORTHSTAR_KIKORYU.md` (KIKORYU product roadmap, tier 1) and the fourteen founder-written 2020 `vs0.md`-`vs13.md` originals → a full SAGA-style reality audit (`IDUNA/docs/VS_REALITY_AUDIT.md` + `docs/kikoryu/*.md`, 16 golden-index rows) — found VS0's identity gate live-but-undocumented, VS2 (tournaments) unbuilt with every prerequisite in place, several VS's reincarnated elsewhere under different names.
- backlog: SECTION 147 (Apple enrichment — GPT-2 fingerprint via a recovered 2020 personal repo's squish/tower/gematria pipeline, ported as `gpt2-alpine-c/pkg/towerprint`; model-fingerprint provenance; astrology field left open) — S147-01/03/05 done (Apples #9915, #9910), S147-02 (the emily-agent enrichment worker) and S147-04 (astrology data source) still open. Also surfaced and fixed a real IDUNA bug while verifying this live: 9226 of 9908 Apples had never synced to the APPLES git mirror (backfilled, `APPLES` commit `699bdd5`; race condition fixed, IDUNA commit `c9217df`).
- backlog: SECTION 148 (chat-to-ledger + personal predictive corpus + Slack firehose — gated on an undispatched governance design pass, explicitly not started blind given real third-party-data questions), SECTION 149 (email operational fabric — AM/PM digests, directive intake, Q&A by email; found `emily-agent/gmail.go`'s `GmailClient` already fully built but dormant, no credentials configured), SECTION 150 (towerprint training-data integration + a real vector-cache component — corrected an initial "nothing exists" finding once the founder supplied a working FAISS/embedding reference implementation, `gpt2-alpine-c/docs/reference/vector_cache.md`).
- backlog: SECTION 146 EPS snapshot pipeline (S146-01..04) implemented directly — `eps_headlines_to_records()`, immutable content-addressed snapshots, eval tombstoning, new `--fable-eps` preset in `prime_directive_dataset.py`. Apple #9913, gpt2-alpine-c commit `458756c`. S146-06 corpus rebuild refreshed the perplexity baseline (116.76→166.56, expected from corpus growth, not a regression) — Apple #9888.
- ops: PRRJECT_FATBABY's core ingestion (`secwatch`/`processor`/`prwatch`/`prwatch-body`) had been down 29 days, found during a data-resale-readiness audit; restarted with founder confirmation — Apple #9912.
- process: full Apple/CHANGELOG audit pass found real gaps — S147-01 had landed and been marked `[x]` with no Apple filed at all (a Fable dispatch was explicitly told not to file one, which was wrong; filed retroactively as Apple #9915), S146-01/02/03 and S147-03 were marked `[x]` citing no Apple ID despite one existing, and IDUNA's own CHANGELOG.md was missing an entry for the S147-02/03/05 commit entirely. All corrected same-day.
- docs: Act III synthesis dispatch failed twice on transient `529 Overloaded` (server-side, not session-limit). Per founder instruction, not auto-retried — intent documented as a standalone dispatch-ready prompt at `docs/fable-prompts/act3-synthesis.md`, tracked "ready to dispatch, paused" in `fable-next-backlog.md`.
- feat(S10-01): front-door audit gate before the morning briefing push — `auditFrontDoor()` in `emily-agent/webaudit.go` runs the existing `web_audit_url` tool against newssite (`:8082`) and signalapi (`:9091` — corrected from `:8083`, a pre-existing port discrepancy between MJOLNIR's displayed URL and where the service actually listens, found while implementing this). Fails on 5xx or >3 broken links; wired into `runMorningBriefing`, suppressing the push and filing an escalation Apple instead on failure. Apple #9917.
- research(S135-02): VS0 sticker vendor comparison — Sticker Mule, StickerApp, StickerGiant, StickerYou compared against the brief's specs; Sticker Mule and StickerApp had extractable per-unit pricing (~$0.23-0.30/unit at 500-1000, well under the $4 target), StickerGiant/StickerYou pricing is behind JS calculators. No vendor auto-selected — the brief's own schedule names this a human decision; added to HUMAN UNBLOCK QUEUE. Apple #9918.
- docs: two TYLER HQ-SPEC docs landed via GitHub web upload (`HQ-SPEC-LORE-104-QUEEN-SALLY-DOCTRINE.md`, `HQ-CANON-TYLER-105-EPOCH-EXTINCTION.md`) — golden-indexed as LORE-104/CANON-TYLER-105 (tier 2), TYLER CHANGELOG updated.
- feat(S141-01): `pkg/norn` kernel library — new top-level repo `NORN/` (added to `go.work`), five interfaces (Artifact, Proposer, Oracle, Gate, Registry) per HQ-SPEC-PRIME-101 §3, `NDJSONRegistry` (mutex-guarded append-only NDJSON + in-memory index), `CheckLineage` (mechanical Löbian hazard check), `DefaultGate` reference implementation. 10 tests, all §8.1 properties covered, green under `-race`. `HQ-SPEC-PRIME-097-fixed-points.md` (NORN's cited mathematical foundation, previously absent from the repo — a known blocking gap) landed the same day; reconciled point-by-point against what was built, zero code changes needed. PRIME-097 golden-indexed (tier 1); SECTION 141 blocking note resolved. Apple #9920.
- ops: `emilyspringerton/NORN` GitHub remote created (founder) and wired up (`git@github.com:emilyspringerton/NORN.git`, SSH — matching every other repo's convention); pushed. HUMAN UNBLOCK QUEUE row resolved. NORN naming confirmation row remains open.
- feat(S141-04, partial): Apples hookup done and verified live — `NORN/pkg/apples`, `NORN` registered as a real IDUNA agent, two real Apples filed (#9926, #9927) from production data. CLI/daemon/Back Office deliberately deferred, not attempted shallow. Apple #9929.
- **incident (fully recovered)**: testing S141-04 live surfaced that `IDUNA/cmd/bootstrap`'s `writeSecretsEnv` had been silently destroying previously-provisioned agents' plaintext secrets on every run that provisioned ≥1 new agent — EMILY-PRIME/FATBABY-EMILY/EMIREE/JON/BOB/TYLER all lost their plaintext (DB hash intact, nothing was actually broken for already-running processes) the moment NORN/EDIS-WOOCOMMERCE/EMILY-TRAINING/EDIS-CUSTODIAN needed provisioning for the first time. EMILY-PRIME recovered from a live process's environment; the other five had no recoverable copy anywhere (verified via grep + full `/proc` scan) and were deliberately rotated. All ten agents verified live post-recovery. Root cause fixed at the source (merge instead of overwrite, 6 regression tests); a related `.gitignore` bug that could have hidden the fix's own test file was also found and fixed. Apple #9930 (escalation, full writeup).
- feat(S141-03): recon matcher migrated to NORN (`PRRJECT_FATBABY/internal/entitygraph/norngate`) — the kernel's second real instantiation and first with genuine production history (~32.7K resolved accuracy records vs. S141-02's fixture-only eval). Real finding surfaced: entity-graph's overall signal precision is 11.9%, with three signal types at exactly 0% — flagged as a follow-up, not fixed here. GameEvolutionEngine leg investigated and found unbuildable as a migration target (no such system exists in the codebase); disclosed as a gap. `norn.GradeAndPromote` extracted into the kernel after the bootstrap-or-gate branch was about to be written a second time by hand. Apple #9924.
- feat(S141-02): first real NORN instantiation — `PRRJECT_FATBABY/internal/eps/norngate` wraps the EPS headline extractor as a `norn.Oracle`, bootstrap-promoted via `cmd/norn-eps-migrate`. Investigated the live oracle store before building the "replay historical decisions" acceptance proof PRIME-101 §8.2 asks for: zero real historical grading decisions exist system-wide (reconciler has never matched a pending case against a filed 8-K, ever — verified, not assumed), so the determinism proof is built against the 4 ground-truth fixtures instead, with the gap disclosed in BACKLOG.md rather than silently closed. Found and fixed a real `pkg/norn` bug live: `Registry.Record` wasn't stamping `Timestamp`. Apple #9922.

## 2026-07-15
- docs(ops): REBOOT_RUNBOOK.md — self-contained VM reboot/hydration runbook for a fresh Claude Code session; linger + iduna.service/emily-system.service enabled for boot-time auto-start; deep audit (Fable) surfaced and fixed a `--all`-starts-SHANKPIT bug, silent exit-0 on child start failure, an over-broad pgrep pattern, a PATH-blind `emily` exec, world-readable JWT secret, and missing `emily start` coverage for entity-graph/eps-reconciler/eps-processor (Apples #9823-9825)

## 2026-06-28
- S148-02: CARTEL opening night — Thursday, Detroit techno, UR lineage, $10/$5 High Life, tortilla truck 10PM-2AM

- S156-01: breakfast taco truck morning mode $4/3for$10; S157-01: punch card buy-10-get-1, no app, paper card

## 2026-06-27
- S155-02: truck-to-truck Tier 0 — peer B2B, spot+weekly rates, 10 trucks with our number in 90 days
- S155-01: tortilla networking strategy — restaurant B2B wedge, event catering, community anchors, media
- S154-01: NORTHSTAR_STINKIES — full bootstrap roadmap VS0 hoodie → Store 0 → CARTEL → Store 1, zero outside capital
- S146-02: VS0 redefined as STINKIES hoodie; stickers moved to VS1; version sequence locked
- S153-01: lighters — STINKIES disposable $3 + EINHORN R&D refillable $18 (antler mark, laser-etched)
- S152-02: strict TP rule locked (Scott only, empty shelf before substitution); house-brand TP roadmap
- S152-01: Scott TP $1.50/roll; store services doc updated with TP + jalapeño line in product matrix
- S150-01: jalapeño line — stuffed hot dogs, pickles, chips, JALAPEÑO KIT $26; S151-01: ROCKET energy drink (lion's mane 500mg, exclusive STINKIES, $4/can)
- S149-01: THE FEED — free meal, no questions, Store 0 Pontiac, day one operational
- S148-01: CARTEL nightclub brief — Pontiac MI venue, full bar, 250-400 cap, Thu-Sat, $1.5M-$3M gross target
- S147-01: tortilla truck brief — fresh corn masa, Store 0 exterior, $0.50/tortilla, 85%+ margin
- S146-01: apparel brief (hoodie + trucker hat STORE 0 · PONTIAC MI) + Store 0 services (public restroom, Michigan Lottery)
- S145-01: free soap brief — lye bars, pallet quantity ≤$0.12/unit, in every bag/on every counter
- S144-01: STINKIES Store 0 — Pontiac MI; beer brief (Miller High Life 6/12/30-pack); Store 0 build-out scope, layout, Michigan licenses, Emily Prime inventory/POS role
- S142-01: STINKIES hot dogs brief (all-beef natural casing 8-pack $12, nacho kit bundle); S143-01: cigarettes brief (matte black, retail-only distribution, regulatory gates)
- S141-01: STINKIES cheese sauce brief — Nacho + Jalapeño ($8/jar, 12oz glass, real cheddar, shelf-stable, $13/mo sub)
- S140-01 update: Powderhorn instant espresso jar ($22) + STINKIES singles ($1, loss leader CAC not margin)
- S140-01 update: Powderhorn is Colombian single-origin (Huila/Nariño), price 17 vs 14 Mountain Man, sub 29/mo
- S140-01: STINKIES COMMISSAIRE coffee brief — Mountain Man (dark/Brazil+Sumatra) + Powderhorn (medium-light/Colombia+Ethiopia); $14/bag, $24/mo subscription, ≥67% margin
- S139-01: VS1 toothbrush brief — STINKIES COMMISSAIRE ($9 entry) + ULTRA ($68 brass handle); both on 90-day head subscription; Supply Chain AGI reorder wired
- S138-01/05: EINHORN INDEX NORTHSTAR + KnowledgeQuery tool (kgraph.go); golden-docs registered
- S137-01/03/04: Research tool (research.go): SHA256 cache, multi-domain fetch, HTML extraction, research_log Apple
- S136-04/05: supply_chain_research + supply_chain_draft_po tools in emilytools.go; registered in main.go dispatcher
- S136-01: Supply Chain AGI NORTHSTAR at docs/NORTHSTAR_SUPPLY_CHAIN.md; golden-docs-index registered
- S135-01: VS0 sticker design brief at docs/merch/stickers_vs0_brief.md — Emily Prime Mark, Wordmark, Logotype; specs, schedule, success metrics

- S131-04/05/06: Slack on HEIMDAL sprint complete/blocked; Slack on escalation apples; emily-prime-cron IDUNA monitor auto-created + check-in at RunOnce end

## 2026-06-25
- Full Slack integration: SlackNotifier (webhook), alerting.go (CheckinAlertWorker polls IDUNA overdue monitors), watchdog alerts → Slack, cron cadence slowed to 15m
- Slow Emily Prime cron from 5m to 15m cadence; comment updates in cron.go, integration.go
- feat: S128-05 federated task router — MaybeRouteTask, ListActiveClusters, routeTask + 5 tests (Apple #3867)
- feat: S128-03 cluster heartbeat loop — SendHeartbeat, detectCapabilities, startHeartbeatLoop in NewAutonomousCycle (Apple #3864)

- feat: S128-01 earnings section in morning briefing — reads earnings-calendar/dates.ndjson, appends EARNINGS THIS WEEK to FCM push body (Apple #3860)

## 2026-06-24
- feat: S125-08 weekly memory consolidation — consolidateMemoryFragments in PLAN phase, 7-day sentinel (Apple #3667)
- feat: S125-07 GET /health endpoint — {ok, service, gear} for gfdapi + PitViper polling (Apple #3664)
- feat: S125-06 GPT-2 inference fallback — gpt2_available in CycleMetrics, planWithFallback + gpt2Fallback via $GPT2_INFER_BIN (Apple #3662)
- feat: S126-12 velocity alerts — observeVelocityAlerts in OBSERVE phase, Apple escalation + FCM push (Apple #3562)
- feat: S126-07 haiku→sonnet escalation — StuckCycles/EscalatedAt on ImprovementTask, escalateTaskWithSonnet in DECIDE phase (Apple #3547)
- feat: S126-06 EMILY↔TRAPX event bridge — DragonBroadcastWorldEvent + TrapXAppleWatcher in dragon.go, wired to main.go startup (Apple #3545)
- feat: S126-05 observation digest — obsdigest.go, cron PLAN phase wire, emily-memory/observations-digest.json (Apple #3541)
- feat: S125-10 FatBaby signal → morning brief — fetchTopSignals from signalapi /v1/data-quality, confidence>=0.75, top-5 appended to push body (Apple #3520)
- S121-05: Archetype Engine Dragon routing — DragonArchetypeAugment() FIELD invocation per Dragon city decision; archetype corridor+spirit stack in Dragon Apple bodies
- S121-03+04: Dragon ACT phase — dragonDecide() escalation rules + dragonACT() fires city events + Dragon Apples per event in RSI cycle
- S121-02: Dragon observer — DragonObserve() reads TRAPX city state in RSI OBSERVE; dragon_observe stream entry; 14 tests
- S123 stubs: TYLER×TRAPX district scenes 200-207, receipt bridge, multi-timeline branch system, flip phone (5-tab diegetic device), VS0 Detroit 2-scene loop — Apple #3359

- TRAPX S121-S122 stubs: Dragon GM spike (Emily Prime as city intelligence, Archetype Engine routing) + GFD urban fantasy crossover (scene cluster 200-299, Watcher/Enforcement packages, 8 quest-gated class chains, TRAPX faction reputation) — Apple #3337

## 2026-06-23
- S103-03: api_archetype.go — archetype engine status+spirits proxy endpoints; /api/v1/emily/archetype/status+spirits
- S102-05: field_bridge.go — THE_FIELD AugmentTaskWithField wired into RSI DECIDE phase; archetypesBridge adapts AnthropicClient; 30s timeout; resonance header prepended to new task descriptions
- S102-04: cmd/archetype-engine HTTP service :8090 — POST /invoke (E1/E2 dual-persona), GET /spirits (all 72), GET /status; builds clean
- S102-01: pkg/archetypes — All72 Spirit structs (M_G v1.0), MatchIntent keyword selector, ResonanceState/E1-E2 weights, amplitude normalization; 17 tests; BACKLOG S102 opened
- ARCHETYPE_ENGINE_NORTHSTAR.md: Dynamic Hybrid AI Agent Archetype Engine northstar — THE_FIELD implementation, dual-persona E1/E2, full 72-spirit Goetia v1.0, resonance corridor router; golden-docs-index registered
- backlog: S76-01/03/04/05/06 all done — IDUNA MMO API + GFD server-go auth/telecrystal/crafting/worldcrisis/skill-xp
- backlog: S75-01/02/03/04/05 IDUNA MMO schema + API marked done (Apples #2935-2936)
- backlog: S82-02/03/S83-02/03/S84-01/02 sub-job, recast, merit, gear IL, crafting marked done (Apples #2925-2930)
- backlog: S81-05/06/S82-01/S83-01 enmity+death+jobs+XP marked done (Apples #2919-2922)
- backlog: S86-01 Conquest, S87-01 NM spawns, S87-02 Treasure pool marked done (Apples #2908/#2909/#2910)
- backlog: S85-01/02/03 Party+Alliance+XPChain marked done (Apple #2904)
- backlog: S84-06 DragonsNShit MUD (Apple #2887)

- backlog: S84-04 Swampville + S84-05 mining skill marked done (Apple #2865, #2868)

## 2026-06-21
- test: S74-01 emilyListDir 4 tests; 72 total (Apple #2467)
- test: S73-01 Emiree FractalFingerprint+Summary; 68 total (Apple #2420)
- test: S72-01 HEIMDAL formatCriteria+min tests; 65 total (Apple #2418)
- test: S71-01 briefing buildBriefingMessage+typeLabel; 61 total (Apple #2416)
- test: S70-01 DocStore append+LoadSeen tests; 56 total (Apple #2413)
- test: S69-01 goldenbuild contentHash+loadGoldenIndex tests; 53 total (Apple #2411)
- test: S68-01 Emiree WitchState test suite; 10 tests, 46 total (Apple #2409)
- test: S65-01 FCM jwt_test.go 3 tests (Apple #2402)
- feat: S49-01 emily-agent GET /api/v1/emily/posture endpoint (Apple #2345)
- feat: S46-03 emily-agent posture gate — SIEGE skips LLM, EXITED blocks startup (Apple #2327)

- feat: S32-02 POST /api/v1/emily/push/test FCM test endpoint (Apple #2306)

## 2026-06-20

- feat: S44-05 emily-agent ConversationStore git pull --rebase before push (Apple #1881)

## 2026-06-18

- feat(watchdog): S42-A SHANKPIT added to Emily Prime service watchdog — pings :6970/healthz (Apple #1436)

## 2026-06-17
- web audit via emily web_audit_url: found /ask 500 (Symbols field), HEAD 405 sitewide, /api/ask 503 (EMILY_BASE_URL), governance-signals empty; signal pipeline: 87% failure rate, all signals stubs, 429/empty-URL/4MB issues; 10 items added to S36+S37

- gpt2 skill registered: gpt2_generate, gpt2_health, gpt2_start tools in emily-agent ToolDispatcher; emily gpt2 generate + emily gpt2 health subcommands in emily.cli. Apple #945.

## 2026-06-16
- FatBaby Emily sub-agent: Emily Prime chats to FatBaby Emily via chat_to_fatbaby_emily; exchange visible in web chat UI as blue sub-chat block (Apple #587)

- feat(ops): S29 production systemd services — emily-system.service + iduna.service written; S29-01/02/03/04 complete. Apples #561 #562.

## 2026-06-15
- feat(api_gpt2): POST /api/v1/gpt2/generate + GET /api/v1/gpt2/health proxy to :8088

- S27-04/05/06: cross-domain synthesis preset + revenue priority wiring + memory/update endpoint

## 2026-06-14
- add TYLER-BIBLE and TYLER-ENGINE to golden-docs-index.md (Tier 2); registered for training corpus pickup
- feat(emily-memory): world-state.md + cycle-log.md created; emily-memory/ activated as persistent world model across cold-start cycles; registered as Tier 1 golden docs (EMILY-MEMORY + EMILY-CYCLE-LOG) in golden-docs-index.md; updateCycleLog() wired into cron.go PLAN phase — each cycle appends outcome to cycle-log.md (bounded to 100 entries via sentinel markers); BACKLOG.md S27 section added (S27-01..06, S27-01..03 complete). Resolves AGI Gap 1 from RSI trajectory memo #457.
- feat(emily-agent): watchdog.go — CheckServiceHealth pings IDUNA :8080, newssite :8082, signalapi :9091, emily-agent :8086 every cycle; escalation Apple when service down ≥ 2 min; checkLogSizes alerts when any *.log in var/logs/ exceeds 500 MB; WatchdogState persisted as watchdog-state.json across cycles (S24-03)
- feat(emily-agent): cron.go — wire CheckServiceHealth into RunOnce immediately after goldenbuild; watchdog alerts post as escalation Apples if IDUNA is reachable
- feat(IDUNA): user_subscriptions table + Subscription type + GetEffectivePermissions injects cap.query.full for active subscribers; POST/GET /api/v1/subscriptions; EDIS-WOOCOMMERCE agent (S23-04 IDUNA 5ed080e)
- feat(PRRJECT_FATBABY): VerifyTokenPermissions on askVerifier interface; iamguard.Guard impl; newssite Emily+ subscriber bypass for cap.query.full permission (S23-04 PRRJECT_FATBABY d3136fc)
- feat(EDIS): emily-plus-woocommerce.php — WooCommerce order completed hook provisions cap.query.full in IDUNA; set-iduna-user REST endpoint; EDIS-WOOCOMMERCE JWT transient (S23-04 EDIS f5dde79)
- feat(EDIS): mailchimp-waitlist.php — replaces wp_options waitlist with Mailchimp API v3 PUT upsert; auto-detects data center; tag support; wp_options fallback (S23-05 EDIS 22b7de5)
- feat(gpt2-alpine-c): NORTHSTAR.md, CLAUDE.md, CHANGELOG.md, .gitignore, scripts/prime_directive_dataset.py, scripts/drive_sync.py, scripts/convert_ft_checkpoint.py, notebooks/gpt2_finetune_colab.ipynb (S26-01..03)
- feat(emily.cli): internal/iduna/client.go DriveUpload/DriveList/DriveGet; cmd/train.go build-dataset/upload/status subcommands (S26-03)
- docs(emily): THE_EMILY_WAY.md — comprehensive operating procedure doc encoding how work is done across all repos; RSI AGI loop discipline, commit/Apple/CHANGELOG protocol, document hierarchy
- golden-index: +19 Tier 2 entries (emily-tools, protocol, framework, cron-evo, integration, IAM, MJOLNIR-apples/push/spec, APPLES-schema, SHANKPIT-netcode/predict/bridge, EmilyOS-arch/memo/posture, GFD-NORTH placeholder); THE-EMILY-WAY as Tier 1
- backlog: ✓ S25-02 GFD NORTHSTAR.md (GFD commit 287733b); ✓ S25-04 SHANKPIT/GFD deduplication (GFD commit 8760b31); S25 section fully complete

## 2026-06-13
- feat(plan): POST /api/v1/emily/plan — haiku sprint-batch planner; accept {question, context?}, return {sprints, summary} (S22-07)
- perf(goldenbuild): per-source SHA-256 content cache (context/golden-cache.json) — only changed sources trigger haiku recompression; saves ~14 haiku calls/rebuild (~63K haiku tokens/day at 5 rebuilds/day)
- perf(rsi): cap iteration history in buildGenerationPrompt to last 3 iterations — prevents prompt bloat on long tasks (~150 tokens saved per omitted iteration record)

## 2026-06-12
- feat(goldenbuild): GoldenDocCompiler — compresses Tier 1 golden docs from all repos into context/full-system-context.md via haiku bilingual compression; MaybeRebuild wired into RunOnce; dynamic buildEmilySystemPrompt() replaces static const; FABLE reads full-system-context.md with GOLDEN.md fallback; EMILY/docs/NORTHSTAR.md written

- emily_write_file: unconditional Apple filing enforced at platform layer via IdunaClient; buildCycleApple now includes tokens_used in metadata; EMILY_PRIME_TOOLS_OPENAPI.yaml spec added

