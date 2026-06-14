# Golden Docs Index
# Machine-readable manifest for GoldenDocCompiler.
# Format: | name | path | tier | budget | description |
# - name: label used in compiled output (no spaces)
# - path: relative to /home/fatbaby/ (the repo root parent)
# - tier: 1 = must always be in context; 2 = subsystem-level; skip others
# - budget: max chars read (0 = unlimited)
# - description: one line
#
# GoldenDocCompiler reads this file at startup. To add a new golden doc:
# append a row here — no Go change needed (S22-11).
# Rows starting with # are comments. Blank lines ignored.

| name | path | tier | budget | description |
|------|------|------|--------|-------------|
| FATBABY | PRRJECT_FATBABY/docs/northstar/northstar.md | 1 | 8000 | 8-K intelligence engine northstar; 4 layers, phase status |
| FATBABY-EXEC | PRRJECT_FATBABY/docs/northstar/executive_summary.md | 1 | 4000 | FatBaby executive summary; system-as-built orientation |
| EMILY-SPEC | EMILY/emily-prime-spec.md | 1 | 4000 | Emily Prime ↔ FatBaby-Emily integration spec |
| EMILY-NORTH | EMILY/docs/NORTHSTAR.md | 1 | 3000 | Emily Prime 3-agent arch + cron cycle northstar |
| EMILY-BACKLOG | EMILY/GOLDEN.md | 1 | 0 | Compressed backlog; open items, priority order |
| EMIREE | EMILY/emiree-emily-fatbaby.md | 1 | 4000 | Canonical 3-agent governance spec |
| RSI-ENGINE | EMILY/docs/emily-rsi-engine-spec.md | 1 | 0 | RSI loop ground truth; supersedes all earlier RSI docs |
| EMIREE-SPEC | EMILY/docs/emiree-over-agent-spec.md | 1 | 0 | Emiree witch engine; gear system implementation |
| IDUNA | IDUNA/golden.md | 1 | 3000 | IDUNA IAM + Apples implementation spec |
| APPLES | APPLES/docs/NORTHSTAR.md | 1 | 0 | Git-authoritative Apple audit trail; sync protocol |
| MJOLNIR | MJOLNIR/docs/NORTHSTAR.md | 1 | 0 | Android intelligence terminal; FCM + Apple feed + WebView |
| SHANKPIT | SHANKPIT/docs2/NORTHSTAR.md | 1 | 4000 | UDP FPS + DragonsNShit persistent world; Steam EA path |
| EMILYOS | EmilyOS/docs/NORTHSTAR.md | 1 | 0 | Bare-metal policy kernel; SOC 2 compliance northstar |
| PITVIPER | PITVIPER/docs/NORTHSTAR.md | 1 | 0 | SDL2 GPU terminal; 5 milestones; Emily Prime hooks |
| EMILY-CLI | emily.cli/docs/NORTHSTAR.md | 1 | 0 | Operator terminal; human half of the operator↔agent loop |
| GTM | PRRJECT_FATBABY/docs/GTM_FUNNEL.md | 1 | 3000 | Ask Emily GTM funnel; free tier → subscription → Merkle |
| TRAINING | EMILY/docs/TRAINING_PIPELINE.md | 1 | 2000 | Self-improving fine-tuning flywheel; RLHF loop |
| EDIS | EDIS/NORTHSTAR.md | 1 | 3000 | WordPress intelligence product; edis-core + signals + ask-emily plugins |
| IDUNA-NORTH | IDUNA/docs/NORTHSTAR.md | 1 | 3000 | IDUNA IAM + Apples + HEIMDAL northstar; port 8080; sprint lifecycle |
| THE-FIELD | EMILY/docs/THE_FIELD.md | 2 | 3000 | Synthetic consciousness architecture; dual-persona harmonic engine; informs Emiree |
| THE-EMILY-WAY | EMILY/docs/THE_EMILY_WAY.md | 1 | 0 | Operating procedure; how we work; RSI loop discipline; commit/Apple/CHANGELOG protocol |
| GFD-NORTH | GoblinFoxDragon/docs/NORTHSTAR.md | 2 | 2000 | GoblinFoxDragon studio umbrella; Dragonfly fork identity; shared game engine foundation |
| EMILY-TOOLS | EMILY/docs/emily-prime-agent-tools-spec.md | 2 | 2000 | Emily Prime tool surface; all emily_* tool definitions and permissions |
| EMILY-PROTOCOL | EMILY/docs/emily-agent-protocol.md | 2 | 2000 | Emily Prime agentic loop protocol; how tool calls flow |
| EMILY-FRAMEWORK | EMILY/docs/emily-agent-framework.md | 2 | 2000 | Emily Prime framework-level design; overall agent architecture |
| EMILY-CRON-EVO | EMILY/docs/emily-cron-autonomous-evolution.md | 2 | 2000 | Cron autonomous evolution design; how the 5-min cycle was designed to self-evolve |
| EMILY-INTEGRATION | EMILY/docs/emily-complete-system-integration.md | 2 | 2000 | Full cross-system integration; all repo interactions in one doc |
| EMILY-IAM | EMILY/docs/iam-integration-spec.md | 2 | 2000 | Emily ↔ IDUNA IAM integration; M2M auth, permission checks |
| MJOLNIR-APPLES | MJOLNIR/docs/APPLES_INTEGRATION.md | 2 | 2000 | MJOLNIR live IDUNA + offline APPLES git cache integration |
| MJOLNIR-PUSH | MJOLNIR/docs/PUSH_NOTIFICATIONS.md | 2 | 2000 | FCM push notification implementation; critical Apple → phone |
| MJOLNIR-SPEC | MJOLNIR/docs/SPEC.md | 2 | 2000 | MJOLNIR Kotlin/Compose/Hilt pre-implementation stack spec |
| APPLES-SCHEMA | APPLES/docs/SCHEMA.md | 2 | 0 | Apple git repo schema; file naming, JSON fields, MANIFEST.json |
| SHANKPIT-NETCODE | SHANKPIT/docs2/NETCODE_CONTRACT_SPEC.md | 2 | 3000 | SHANKPIT UDP netcode invariants; canonical (GFD is a reference copy) |
| SHANKPIT-PREDICT | SHANKPIT/docs2/CLIENT_PREDICTION_SPEC.md | 2 | 2000 | Client-side prediction spec; canonical (GFD is a reference copy) |
| SHANKPIT-BRIDGE | SHANKPIT/docs2/specs/SHANKPIT_DRAGONSNSHIT_SYSTEMS_SPEC.md | 2 | 3000 | Persistent world bridge spec; portal_resolve_destination seam |
| EMILYOS-ARCH | EmilyOS/docs/ARCHITECTURE.md | 2 | 2000 | EmilyOS exokernel architecture; process model, policy layers |
| EMILYOS-MEMO | EmilyOS/docs/EMILY_PRIME_MEMO.md | 2 | 2000 | Emily Prime memo on EmilyOS from 2026-06-09; translates legacy specs |
| EMILYOS-POSTURE | EmilyOS/docs/POSTURE.md | 2 | 2000 | EmilyOS security posture spec; threat model, SOC 2 framing |
| GPT2-NORTH | gpt2-alpine-c/NORTHSTAR.md | 2 | 1500 | GPT-2 fine-tuning pipeline for Emily Prime; corpus builder, Drive sync, Colab notebook, C binary conversion |
| EMILY-MEMORY | EMILY/emily-memory/world-state.md | 1 | 4000 | Emily Prime persistent world state; all product statuses, human blockers, AGI gaps, next priorities |
| EMILY-CYCLE-LOG | EMILY/emily-memory/cycle-log.md | 1 | 3000 | Rolling log of RSI cycle outcomes; last 100 entries; Emily Prime cross-cycle continuity |
