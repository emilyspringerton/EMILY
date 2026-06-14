## 2026-06-14
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

