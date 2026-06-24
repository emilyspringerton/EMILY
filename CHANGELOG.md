## 2026-06-24

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

