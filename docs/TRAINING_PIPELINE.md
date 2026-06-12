# Self-Improving Training Pipeline
## Emily Prime Fine-Tuning Flywheel | Einhorn Industrial

*Written: 2026-06-12 | Status: spec — long-term initiative*

---

## Why This Matters

Emily Prime currently runs on Anthropic's Claude models (haiku for classification, Sonnet for
implementation). Fine-tuning on EINHORN_INDUSTRIAL operational data would produce a model that:
- Understands the domain vocabulary (FATBABY signals, Apple filing, RSI loop, sector terminology)
- Produces lower-latency responses on routine classification tasks
- Embeds the specific conventions (backlog format, Apple body structure, CHANGELOG entries) without
  needing them in the system prompt on every call

This is the long-term play: Emily Prime as a fine-tuned specialist model, not just a general model
with a long system prompt.

---

## Data Collection Layer

### What to collect

1. **RSI task interactions** — Emily Prime generates an RSI task, Claude Code executes it,
   produces output. The (task_spec → implementation_output) pair is the training signal.
   Source: EMILY/signals/tasks/ + EMILY/var/rsi-loop-state.json

2. **Apple classifications** — Every Apple filed has a `apple_type` field. The (event_description
   → apple_type) classification is a straightforward supervised task.
   Source: IDUNA Apple log (POST /api/v1/apples payloads)

3. **HEIMDAL translations** — MJOLNIR sends a product requirement; Emily Prime translates it to
   an RSI roadmap item. The (requirement → roadmap_item) pair is a structured generation task.
   Source: EMILY/var/emily-state/ HEIMDAL cycle logs

4. **Triage decisions** — FatBaby observations arrive; Emily Prime classifies them as
   `backlog` (new work), `noise` (dismiss), or `observation` (Apple, no action).
   Source: emily-agent cron.go triage cycle logs

5. **Briefing generations** — Morning briefing prompt → output. Short-form structured
   generation from multi-source context.
   Source: emily-agent briefing.go outputs

### Data format

Training pairs are stored as JSONL in `EMILY/var/training-data/`:

```json
{"id":"uuid","source":"rsi_task","created_at":"2026-06-12T00:00:00Z","prompt":"...","completion":"...","metadata":{"apple_id":123,"task_id":"...","model":"claude-haiku-4-5"}}
```

Fields: `id`, `source` (which collection), `created_at`, `prompt`, `completion`, `metadata`.
The `metadata` links back to the original Apple/task/event so examples can be filtered by outcome.

---

## Annotation Pipeline

Raw data is not training data. Annotation determines which examples to keep:

1. **Outcome filter**: RSI tasks where the downstream build passed tests → positive examples.
   Tasks that required a retry or produced a broken build → negative examples (exclude or label
   for DPO negative).

2. **Quality gate**: Any example where the completion was manually corrected by Claude Code
   (edit vs. accept unchanged) gets flagged for human review before inclusion.

3. **Deduplication**: Exact-match dedup on prompt; near-dedup via MinHash on completion.

4. **Annotation UI (future)**: A lightweight web interface (at `:8087`) where Emily / human
   reviewers label examples as `keep` / `discard` / `needs_edit`. Writes to IDUNA as
   `training_annotation` Apple type.

---

## Fine-Tuning Approach

**Phase 1 (immediate):** Collect data passively. No training yet. Build the pipeline.

**Phase 2 (6-12 months):** Supervised fine-tuning on Anthropic's fine-tuning API
(currently in beta) using the collected prompt/completion pairs. Start with Apple
classification (simplest, lowest stakes) and HEIMDAL translation.

**Phase 3 (12+ months):** RLHF loop. Use outcome signals (did the RSI task produce a
passing build? did the triage decision prove correct?) as reward signals for preference
tuning (DPO or RLHF depending on available tooling).

**Model strategy:** Fine-tune haiku-tier first (cheapest, most frequent calls). Sonnet-tier
fine-tuning deferred until haiku results validate the approach.

---

## Privacy and Safety

- Training data never includes: user passwords, API keys, private keys, session tokens.
- The system prompt is stripped from collected pairs — only the task-specific portion is kept.
- Any collected example can be deleted via `emily primetask delete-training --id <uuid>`.
- Training data collection is opt-in by default; controlled by `EMILY_COLLECT_TRAINING=1` env var.

---

## Implementation Plan

| Phase | What ships | Timeline |
|-------|-----------|----------|
| 1a | EMILY_COLLECT_TRAINING flag + JSONL writer in emily-agent | Sprint 1 |
| 1b | Apple classification data capture (triage cycle) | Sprint 1 |
| 1c | RSI task outcome linking (pass/fail from test results) | Sprint 2 |
| 2a | Annotation UI stub (`:8087`) | Sprint 3 |
| 2b | Dedup + quality gate pipeline | Sprint 4 |
| 3a | Anthropic fine-tuning API integration | When API available |
| 3b | DPO pipeline from outcome signals | Long-term |

---

## Relationship to RSI Loop

The training pipeline IS an RSI loop at a longer time scale:
- Immediate loop (5 min): Emily observes → Claude Code acts → Apple filed
- Training loop (months): Emily observes patterns in outcomes → fine-tuned model acts more
  precisely → better outcomes → better training data → better model

The backlog item for each RSI iteration is the forward pass. The training pipeline is the
weight update.
