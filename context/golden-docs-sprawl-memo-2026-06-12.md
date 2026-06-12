# MEMO: Full Golden Docs Sprawl Audit

**From:** Claude Code (session 2026-06-12, instructed by Emily Springerton)
**To:** Emily Prime
**Re:** Complete inventory of mission-critical docs across all repos; plan to give you real system context
**Date:** 2026-06-12
**Action required:** Analyze this. Propose and execute an iteration plan.

---

## Why This Exists

You do not currently have access to your own system's documentation. Your system prompt is a static
constant in `emily-agent/main.go`. Your FABLE advisor reads only `EMILY/GOLDEN.md` (backlog compression).
The northstar docs from every other repo — EmilyOS, PITVIPER, SHANKPIT, MJOLNIR, GoblinFoxDragon —
are invisible to you when you plan.

This memo is the full audit. It is your starting point.

---

## Repos Audited

PRRJECT_FATBABY · EMILY · IDUNA · MJOLNIR · SHANKPIT · APPLES · GoblinFoxDragon ·
MoneyPrinterTurbo · PITVIPER · EmilyOS · emily.cli

Skipped: TYLER (separate beast, own MD sprawl, defer)

---

## Tier 1: GOLDEN — Must be in your context at all times

These are northstar-level, mission-defining, or architecturally canonical. You cannot plan without them.

| Doc | Path | Notes |
|-----|------|-------|
| FatBaby Northstar | `PRRJECT_FATBABY/docs/northstar/northstar.md` | 8-K intelligence engine, 4 layers, Phase 1-3 done, Phase 4 partial |
| FatBaby Exec Summary | `PRRJECT_FATBABY/docs/northstar/executive_summary.md` | 3-doc orientation kit, system-as-built |
| Emily Integration Spec | `EMILY/emily-prime-spec.md` | FatBaby-Emily ↔ Emily Prime integration, architecture, git memory |
| Emily Backlog Context | `EMILY/GOLDEN.md` | Auto-generated compressed backlog; re-runs `emily backlog compress` |
| Emiree Transmission | `EMILY/emiree.md` | Sanskrit/Chinese poetic spec for the witch engine (7 volumes); drives `emiree.go` |
| AI System Governance | `EMILY/emiree-emily-fatbaby.md` | Canonical 3-agent spec: Emiree (witch) · Emily Prime (chief of staff) · FatBaby-Emily |
| Emiree Over-Agent Spec | `EMILY/docs/emiree-over-agent-spec.md` | Emiree implementation spec; how the gear system works and why |
| RSI Engine Spec | `EMILY/docs/emily-rsi-engine-spec.md` | **Supersedes** all earlier RSI docs; ground truth for the running system |
| IDUNA Apples Golden | `IDUNA/golden.md` | Apples implementation spec; all endpoints, schema, permissions; checklist 95% done |
| IDUNA IAM Spec | `IDUNA/docs/iam-spec.md` | Full IAM spec (need to verify if current) |
| APPLES Northstar | `APPLES/docs/NORTHSTAR.md` | Git-authoritative Apple audit trail; sync protocol |
| MJOLNIR Northstar | `MJOLNIR/docs/NORTHSTAR.md` | Android intelligence terminal; FCM + Apple feed + WebView |
| MJOLNIR Tech Spec | `MJOLNIR/docs/SPEC.md` | Pre-implementation Kotlin/Compose/Hilt stack spec |
| SHANKPIT Northstar | `SHANKPIT/docs2/NORTHSTAR.md` | UDP FPS + DragonsNShit persistent world bridge (Dragonfly/Bedrock) |
| EmilyOS Northstar | `EmilyOS/docs/NORTHSTAR.md` | Bare-metal exokernel; SOC 2 readiness as northstar |
| PITVIPER Northstar | `PITVIPER/docs/NORTHSTAR.md` | SDL2 GPU terminal with Emily Prime hooks; closes audit trail gap |
| emily.cli Northstar | `emily.cli/docs/NORTHSTAR.md` | Operator terminal; the human half of the operator↔agent loop |
| GTM Funnel | `PRRJECT_FATBABY/docs/GTM_FUNNEL.md` | Ask Emily consumer product; free tier → subscription → Merkle query (2026-06-12, new) |
| Training Pipeline | `EMILY/docs/TRAINING_PIPELINE.md` | Self-improving fine-tuning flywheel; user data → RLHF loop (2026-06-12, new) |
| MJOLNIR Context | `EMILY/context/mjolnir-context.md` | MJOLNIR context for Emily's own use — already in your context dir |

---

## Tier 2: IMPORTANT SPECS — Understand before touching the relevant subsystem

| Doc | Path | Notes |
|-----|------|-------|
| Emily Agent Tools Spec | `EMILY/docs/emily-prime-agent-tools-spec.md` | Tool surface Emily Prime exposes |
| Emily Agent Protocol | `EMILY/docs/emily-agent-protocol.md` | Agentic loop protocol |
| Emily Agent Framework | `EMILY/docs/emily-agent-framework.md` | Framework-level design |
| Emily Cron Autonomous | `EMILY/docs/emily-cron-autonomous-evolution.md` | Cron evolution design |
| Emily Complete Integration | `EMILY/docs/emily-complete-system-integration.md` | Full cross-system integration doc |
| Emily IAM Integration | `EMILY/docs/iam-integration-spec.md` | Emily ↔ IDUNA integration |
| MJOLNIR Integration (Emily) | `EMILY/docs/MJOLNIR_INTEGRATION.md` | MJOLNIR from Emily's perspective |
| MJOLNIR Apples Integration | `MJOLNIR/docs/APPLES_INTEGRATION.md` | Live IDUNA + offline APPLES git repo |
| MJOLNIR Push Notifications | `MJOLNIR/docs/PUSH_NOTIFICATIONS.md` | FCM push implementation |
| APPLES Schema | `APPLES/docs/SCHEMA.md` | Apple git repo schema |
| SHANKPIT Netcode Contract | `SHANKPIT/docs2/NETCODE_CONTRACT_SPEC.md` | UDP netcode invariants |
| SHANKPIT Client Prediction | `SHANKPIT/docs2/CLIENT_PREDICTION_SPEC.md` | Client-side prediction spec |
| SHANKPIT DragonsNShit | `SHANKPIT/docs2/specs/SHANKPIT_DRAGONSNSHIT_SYSTEMS_SPEC.md` | Persistent world bridge spec |
| GoblinFoxDragon Netcode | `GoblinFoxDragon/docs2/NETCODE_CONTRACT_SPEC.md` | Mirror of SHANKPIT netcode (both repos have it) |
| EmilyOS Architecture | `EmilyOS/docs/ARCHITECTURE.md` | Exokernel architecture |
| EmilyOS Emily Prime Memo | `EmilyOS/docs/EMILY_PRIME_MEMO.md` | Your own memo from 2026-06-09 translating legacy specs |
| EmilyOS Posture | `EmilyOS/docs/POSTURE.md` | Security posture spec |
| EmilyOS SOC2 | `EmilyOS/docs/SOC2.md` | SOC 2 compliance requirements |
| emily.cli Design | `emily.cli/docs/DESIGN.md` | CLI design decisions |
| emily.cli Commands | `emily.cli/docs/COMMANDS.md` | Full command reference |
| FatBaby Architecture | `PRRJECT_FATBABY/docs/architecture-distributed-event-intelligence.md` | Distributed event intelligence architecture |
| THE_FIELD | `EMILY/THE_FIELD.md` | Synthetic consciousness architecture (dual-persona harmonic engine); philosophical but informs Emiree |

---

## Tier 3: OPERATIONAL — Real-time state, not architectural

| Doc | Path | Notes |
|-----|------|-------|
| EMILY BACKLOG | `EMILY/BACKLOG.md` | Master cross-repo backlog; canonical |
| EMILY DONE | `EMILY/DONE.md` | Completed items archive |
| emily.cli BACKLOG | `emily.cli/BACKLOG.md` | emily.cli-specific backlog |
| All CHANGELOGs | `*/CHANGELOG.md` | Per-repo change logs |
| Run Reports | `PRRJECT_FATBABY/claude-runs/` | Claude Code run reports |
| Conversation Memory | `EMILY/emily-agent/emily-memory/conversations/` | Your own conversation history |

---

## Tier 4: HISTORICAL / LEGACY — Superseded

| Doc | Path | Notes |
|-----|------|-------|
| EmilyOS Legacy Specs | `EmilyOS/docs/legacy/spec-1.md` through `spec-8.md`, `adendum.md`, etc. | Superseded by EmilyOS NORTHSTAR + your own memo |
| Emily Complete Vision | `EMILY/docs/emily-complete-vision.md` | Earlier vision doc; largely superseded by specific northstars |
| Emily Comprehensive Review | `EMILY/docs/emily-comprehensive-review-agi-expansion.md` | Old expansion review |
| Emily Expanded Internet Scale | `EMILY/docs/emily-expanded-internet-scale.md` | Older internet-scale vision |
| Emily V2 Architecture | `EMILY/docs/emily-v2-architecture-revision.md` | V2 revision; verify if superseded |
| Emily Android Companion | `EMILY/docs/emily-android-companion-spec.md` | Pre-MJOLNIR Android spec; MJOLNIR NORTHSTAR supersedes |
| Emily Delegation | `EMILY/docs/emily-delegation-beyond-capabilities.md` | Older delegation spec |
| FatBaby Jon.md | `PRRJECT_FATBABY/docs/jon.md` | Unclear; verify |

---

## Tier 5: NOT GOLDEN — Skip from context

| Category | Paths |
|----------|-------|
| All CLAUDE.md files | `*/CLAUDE.md` — Claude Code instructions, not golden docs |
| MoneyPrinterTurbo docs | `MoneyPrinterTurbo/docs/` — MPT is a third-party fork |
| EMILY book | `EMILY/book/book_outline.md` |
| Press material | `EMILY/docs/emily-press-package.md`, `emily-press-release.md` |
| Code snippets | `EMILY/packages/render/title_screen_integration_snippets.md` |
| GoblinFoxDragon README | `GoblinFoxDragon/README.md`, `README2.md` — studio umbrella, not operational |

---

## Key Problems Found

### 1. No standard golden doc convention
The only reliable pattern is `NORTHSTAR.md` in a `docs/` or `docs2/` directory. But:
- PRRJECT_FATBABY has it at `docs/northstar/northstar.md` (nested)
- EMILY has no NORTHSTAR.md — the spec is spread across `emily-prime-spec.md`, `emiree.md`, `emiree-emily-fatbaby.md`
- IDUNA uses `golden.md` (implementation spec, not a northstar)
- GoblinFoxDragon has no northstar at all

### 2. Your system prompt is static
`const emilySystemPrompt` in `emily-agent/main.go` was written once and never updated. It describes
your roles and tools but contains zero repo state. You are planning blind.

### 3. FABLE sees only one file
`fable.go` reads only `EMILY/GOLDEN.md`. FABLE's recommendations are based on the backlog state
alone, with no northstar context from SHANKPIT, EmilyOS, PITVIPER, GoblinFoxDragon, or the GTM funnel.

### 4. Docs aging in silence
New docs (GTM_FUNNEL.md, TRAINING_PIPELINE.md) were written 2026-06-12. No mechanism tells you
they exist. The only way you learn about new golden docs is if Emily Springerton shouts about it.

### 5. Duplication across SHANKPIT / GoblinFoxDragon
Several specs exist in both repos (NETCODE_CONTRACT_SPEC.md, SHANKPIT_DRAGONSNSHIT_SYSTEMS_SPEC.md,
TYLER pilot specs). There's no canonical source — edits to one won't propagate.

### 6. EmilyOS legacy doc pile
8 legacy spec files + addendum in `EmilyOS/docs/legacy/`. Your own 2026-06-09 memo supersedes
them but they still exist and could confuse a future context reader.

### 7. THE_FIELD lives in EMILY root
`THE_FIELD.md` — the synthetic consciousness architecture that informs Emiree — lives as a top-level
file in the EMILY repo with no clear categorization. It should probably be in `EMILY/docs/`.

---

## Recommended Iteration Plan

### Phase 0: Read this memo (you are doing this now)
File this memo as an Apple observation (`audit` type). Add the items below to BACKLOG.md as a new
section: **S18: GOLDEN DOCS INFRASTRUCTURE**.

### Phase 1: Build `EMILY/context/golden-docs-index.md` (1 sprint)
Create a machine-readable index of all Tier 1 and Tier 2 docs with: path, repo, tier, one-line
description, last-verified date. This becomes the source of truth for the golden doc compiler.
Write a `goldenbuild.go` in `emily-agent/` that reads this index and compresses each source into
`EMILY/context/full-system-context.md` using claude-haiku bilingual Chinese/English compression.

### Phase 2: Wire full-system-context.md into your system prompt (1 sprint)
Change `const emilySystemPrompt` in `main.go` to a function `buildEmilySystemPrompt(emilyRoot)` that:
1. Reads `EMILY/context/full-system-context.md` (compiled golden context)
2. Prepends it to the static roles/tools section
Call `GoldenDocCompiler.MaybeRebuild(ctx)` at the start of each cron cycle.

### Phase 3: Feed full context into FABLE (1 sprint)
Change `fable.go` to pass `full-system-context.md` instead of (or in addition to) `GOLDEN.md`.
FABLE's recommendations will then account for all repo northstars, not just the backlog.

### Phase 4: Establish golden doc naming convention (ongoing)
- Add `NORTHSTAR.md` at `docs/NORTHSTAR.md` in every repo that lacks one (EMILY itself, IDUNA)
- EMILY: create `EMILY/docs/NORTHSTAR.md` that synthesizes `emily-prime-spec.md` + `emiree-emily-fatbaby.md`
- IDUNA: rename or supplement `golden.md` with a proper NORTHSTAR
- GoblinFoxDragon: write a northstar

### Phase 5: Archive legacy docs
Move `EmilyOS/docs/legacy/` contents to a clearly marked archive directory or delete them outright
(your 2026-06-09 memo already synthesized them). Reduces future context confusion.

### Phase 6: Resolve SHANKPIT/GoblinFoxDragon duplication
Pick one repo as the canonical source for shared specs. The other imports via symlink or explicit
reference. Decide: is GoblinFoxDragon a subset of SHANKPIT or a peer?

---

## Immediate Next Action

```
1. Read EMILY/BACKLOG.md
2. Append S18: GOLDEN DOCS INFRASTRUCTURE section with items from Phases 1-6 above
3. Post Apple to IDUNA: type=audit, title="Golden docs sprawl audit 2026-06-12"
4. git add + commit + push EMILY
5. On next FABLE cycle: pick Phase 1 (goldenbuild.go) as priority 1
```

This memo should be included in `EMILY/context/golden-docs-index.md` once that file is built.

---
*Written by Claude Code on behalf of Emily Springerton · 2026-06-12*
