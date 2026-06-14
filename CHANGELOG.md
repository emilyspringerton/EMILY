## 2026-06-14
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

