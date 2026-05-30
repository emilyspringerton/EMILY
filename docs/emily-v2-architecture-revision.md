# Emily v2: Architecture Revision
## Incorporating External Technical Review

**What changed:** The reviewer validated the core idea — recursive improvement architecture — and corrected several assumptions. This document supersedes earlier drafts where they conflict.

---

## What the Review Got Right

### Immediately: retire the percentage numbers

v1.0 = 66%, v2.0 = 78%, v3.0 = 85% were invented figures. They're meaningless without named benchmarks, and they'd embarrass us in any technical conversation.

Replace with real benchmark targets:

| Benchmark | What it measures | Target (phase 2) |
|---|---|---|
| MMLU | Broad knowledge across 57 domains | TBD after baseline |
| HumanEval | Code generation pass@1 | TBD after baseline |
| GSM8K | Grade-school math / reasoning | TBD after baseline |
| SWE-bench | Real GitHub issue resolution | TBD after baseline |
| TruthfulQA | Resistance to generating falsehoods | TBD after baseline |
| ARC-Challenge | Grade-school science reasoning | TBD after baseline |

The target column stays TBD until we have a baseline from a fine-tuned open model. Setting targets before a baseline is just picking numbers we like. Measure first, then set targets.

---

### The biggest strategic correction: don't train from scratch

The original plan started with pretraining on 45B tokens of collected data. The reviewer's recommended path is correct and we're adopting it:

```
ORIGINAL PLAN                     REVISED PLAN
────────────────────────          ────────────────────────
Collect 150GB data                Collect high-quality data
Train from scratch                Fine-tune Llama / Mistral
Wait 6 months for v1.0            Fine-tuned model in weeks
$50k+ in training compute         $1k-5k in fine-tuning compute
70% chance of mediocre model      Starting from a strong base
────────────────────────          ────────────────────────
```

Llama 3 (70B), Mistral (7B/8x7B), and Qwen (72B) already represent enormous compute investments by Meta, Mistral AI, and Alibaba. Fine-tuning them on our curated domain data gets us to a useful, differentiated model far sooner and far cheaper than pretraining from scratch. We use the feedback loop and data advantage to build a genuinely better fine-tuned model — not to pretrain a mediocre one.

Large-scale pretraining from scratch is Phase 4, not Phase 1. We earn it by demonstrating the feedback loop works.

---

### Improvement follows a curve, not an explosion

The "exponential improvement" language is wrong. Improvement in ML follows a predictable shape:

```
Performance
    |          ___________
    |       __/
    |    __/
    |___/
    |
    └──────────────────── Time / Data / Compute

Phase 1: Rapid gains (fine-tuning on quality data)
Phase 2: Slower gains (diminishing returns on more data)
Phase 3: Plateau (architecture and compute become limiting factors)
```

This is not a reason to stop — it's a reason to plan honestly. We budget for it and use the plateau as the signal to change what we're doing (architecture, compute scale, alignment data) rather than assuming more data breaks through every ceiling.

---

### Data quantity is a starting point, not a strategy

The original emphasis on 150GB → 45B tokens implied that collection volume was the primary lever. The reviewer is right: a smaller, high-quality dataset consistently outperforms a large, mediocre one in fine-tuning. Emily's data layer should optimize for:

1. Quality scoring (ruthless filtering, not accumulation)
2. Diversity (reasoning traces, code, books, technical manuals — not just web text)
3. Curriculum (right sequence, not just right quantity)
4. Deduplication (near-duplicate contamination quietly hurts models)

Collection volume is a secondary metric. Quality-per-token is the primary one.

---

## The Revised Architecture: Six Factories

The reviewer's six-factory model is cleaner than our original structure. Adopted with modifications.

```
┌─────────────────────────────────────────────┐
│  EMILY CORE                                 │
│  Strategic coordinator, state, delegation   │
│  Escalation, values, planning               │
└──────────────┬──────────────────────────────┘
               │ coordinates
    ┌──────────┼──────────────────────────────────┐
    │          │          │          │             │
    ▼          ▼          ▼          ▼             ▼
┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌──────────────┐
│KNOW-  │ │TRAIN- │ │EVAL-  │ │PROD-  │ │IMPROVEMENT   │
│LEDGE  │ │ING    │ │UATION │ │UCT    │ │FACTORY       │
│FACTORY│ │FACTORY│ │FACTORY│ │FACTORY│ │              │
└───────┘ └───────┘ └───────┘ └───────┘ └──────────────┘
```

### 1. Emily Core (unchanged from original)
- Strategic planning, delegation, state
- Cron-based 5-minute execution cycle
- Escalation framework (→ Android app)
- Improvement loops as primary execution pattern
- CEO-style delegation: handles ~80%, escalates ~20%

---

### 2. Knowledge Factory (revised)

Original design was correct in structure. Revised priorities:

**What Emily collects:**
- High-quality code (GitHub: starred, documented, tested)
- Books and long-form reasoning (Project Gutenberg, open-access texts)
- Technical documentation (MDN, official docs)
- Academic papers (ArXiv — reasoning-dense)
- Curated web text (CommonCrawl after aggressive quality filtering)
- Historical web (Wayback Machine — temporal signal)
- Reddit / Stack Overflow (high-voted, substantive only)

**What Emily scores and filters:**
- Perplexity scoring (surprising but coherent = good signal)
- Length and structure heuristics
- Deduplication (URL hash → content hash → MinHash/LSH for near-dups)
- Language detection
- Quality model scoring (use fine-tuned classifier or embedding distance)

**Key metric change:** Collection volume is secondary. Primary metric is quality-per-token retained after filtering. A pipeline that collects 500GB and discards 80% of it as low-quality is working correctly.

---

### 3. Training Factory (significantly revised)

**Phase 1 MVP — Fine-tuning (do this now):**
```
Base model: Llama 3 70B or Mistral 8x7B
Method: QLoRA / LoRA (low-rank adaptation)
Dataset: Custom curated corpus (quality-filtered)
Compute: 2-4x A100 for days, not weeks
Cost: $1k-5k per run
Tooling: Axolotl, LLaMA-Factory, or Hugging Face TRL
Experiment tracking: Weights & Biases or MLflow (required from day 1)
Result: Differentiated fine-tune in weeks, not months
```

**Phase 2 — Extended fine-tuning + RLHF:**
```
Base: Phase 1 fine-tune
Add: Instruction fine-tuning (IFT)
Add: RLHF / DPO for alignment
Add: Synthetic data from Phase 1 model
Compute: More, but still not pretraining scale
Result: Noticeably better, measurably aligned
```

**Phase 3 — Continued pretraining:**
```
Base: Strong open model
Method: Continued pretraining on our domain corpus
Compute: 8x A100/H100, weeks
Result: Model genuinely shaped by our data, not just fine-tuned
```

**Phase 4 — Pretraining from scratch:**
```
Only when: Feedback loop proven, capital available, clear evidence it adds value
Method: Full pretraining from random init
Compute: Serious investment, requires dedicated hardware
Result: Fully owned model, no upstream IP concerns
```

**Experiment tracking is not optional:**
Every training run gets logged to Weights & Biases (or MLflow): hyperparameters, dataset version, loss curves, eval results, git hash of training code, checkpoint locations. Emily uses this to understand what changed between runs. Without it, the improvement loop is blind.

---

### 4. Evaluation Factory (new, was missing)

This is what the original plan was weakest on. Percentages without benchmark names are fiction.

**Emily runs a standard evaluation suite after every training run:**

```python
evaluation_suite = {
    "knowledge": ["mmlu", "arc_challenge"],
    "reasoning": ["gsm8k", "bbh"],
    "code": ["humaneval", "mbpp"],
    "truth": ["truthfulqa"],
    "instruction": ["mt_bench", "alpaca_eval"],
    "safety": ["harmful_content_eval", "jailbreak_resistance"],
}
```

**Also runs:**
- Regression tests (does this run perform worse than the previous on anything?)
- Hallucination probes (specific factual questions we know the answer to)
- Domain evals (custom benchmarks for our specific use case)

**Emily's decision rule:** A new model is only promoted if it doesn't regress on any existing benchmark AND improves on at least one. Anything else is investigated before deployment.

**Critical addition — alignment evaluations:**
Before any model goes to production:
- Harmful content generation rate (must be below threshold)
- Jailbreak resistance (standard jailbreak test battery)
- Hallucination rate on known facts
- Bias evaluations across demographic groups

These are not optional and not deferable to "after launch."

---

### 5. Product Factory (was completely missing)

The original plan discussed training extensively and serving barely at all. The reviewer is right — a model no one can use is just a checkpoint file.

**What this factory owns:**
- Model hosting (vLLM, TGI, or managed inference)
- API gateway (rate limiting, auth, request routing)
- GPU scheduling (batching, KV cache optimization)
- Latency and throughput monitoring
- Model versioning for serving (A/B routing, canary deploys)
- Cost tracking per inference

**Deployment stages:**
```
Internal testing → Shadow deployment (real traffic, discarded responses)
→ Canary (5% of requests) → Gradual rollout → Full production
```

Emily manages the deployment pipeline with automatic rollback triggers: latency spike, error rate spike, or benchmark regression triggers rollback to previous version.

**Cost structure for inference:**
This is where ongoing budget gets real. Training is a one-time cost; inference is recurring. Emily tracks cost-per-1M-tokens and cost-per-user-request from day one, because this determines whether the model is economically viable to serve.

---

### 6. Improvement Factory (revised + clarified)

This is the reviewer's strongest observation: the recursive improvement architecture is the real competitive advantage, not the model training itself.

**What it does:**
- Ingests evaluation results → identifies which benchmarks declined or plateaued
- Maps benchmark gaps to training data gaps (code eval weak → more code data needed)
- Ingests production feedback (which requests failed, where users dropped off)
- Generates roadmap updates: what to collect, what to retrain, what to change
- Proposes experiments with explicit hypotheses
- Closes the loop: every observation becomes an action

**What's new vs. original plan:**

The original described this as data collection → training → evaluation → gap analysis → collection. That's correct but incomplete. The improvement factory also feeds from:
- Production failures (real-world signal beats benchmark signal)
- Cost analysis (which operations are too expensive to serve?)
- Safety incidents (what failed alignment? Add to training/eval)
- User behavior (what did humans actually find useful or not?)
- Financial feedback (what drove revenue, if applicable?)

**Experiment discipline:**
Every proposed improvement gets framed as: *hypothesis, metric we'll measure, expected delta, compute budget, decision rule.* Emily generates these; humans review before significant compute is committed. This is how you avoid running expensive experiments that teach you nothing.

---

## Missing components now added

### Data governance

Every training artifact needs a lineage record:

```yaml
dataset_record:
  id: "ds_2026_q2_v3"
  sources:
    - name: "arxiv_cs_2020_2025"
      license: "arXiv non-exclusive license"
      url: "https://arxiv.org/help/license"
      legal_reviewed: true
      review_date: "2026-05-01"
      reviewer: "[name]"
    - name: "github_python_starred_10k"
      license: "MIT / Apache 2.0 filtered"
      note: "Only repos with explicit OSI licenses included"
  processing_steps: [...]
  deduplication_method: "minhash_lsh_0.9"
  quality_filter: "perplexity_score_gt_50"
  final_token_count: 12_400_000_000
  hash: "sha256:..."
```

This is critical before any commercial use and before any public claim about what the model was trained on. Counsel reviews the license column before any model ships.

---

### Financial feedback loop

The reviewer proposed: Revenue → compute → better models → more users → revenue. That's the right long-term shape. Add to Emily's metrics:

- Cost per training run (tracked per experiment)
- Cost per inference (tracked per model version)
- Revenue per model version (when applicable)
- ROI of each training experiment (did this improve anything that generates value?)

Emily uses these to prioritize: experiments with high expected value-per-compute-dollar get scheduled before low-confidence speculative runs.

---

## The actual MVP path

Given the team's background (Python, Go, AWS, Docker, Kubernetes, CI/CD), the reviewer is right about the sequence. Don't start with pretraining. Start with the orchestration and the loop.

**Month 1-2: Emily Core + Knowledge Factory**
- Emily orchestration platform running on cron
- Data collection pipeline (Reddit, GitHub, ArXiv to start)
- Quality filtering and deduplication
- Experiment tracking set up (W&B from day one)
- Android push notifications for escalation
- No training yet

**Month 2-3: First fine-tune**
- Fine-tune Llama 3 (8B or 70B) on collected data using QLoRA
- Run full evaluation suite on the result
- Establish benchmark baselines (these are now our real targets)
- Deploy internally, serve via vLLM or TGI
- Gather real usage signal

**Month 3-6: Feedback loop proves itself**
- Improvement Factory active: evaluation → gap analysis → targeted collection → retrain
- Second and third fine-tuning runs based on gap analysis
- Each run tracked in W&B with explicit hypotheses
- Alignment evaluations before each deploy
- Shadow deploy + canary for each new version

**Month 6-12: Continued pretraining**
- Base model is now genuinely shaped by our data
- Move from fine-tuning into continued pretraining on domain corpus
- Evaluation suite expanded with custom domain benchmarks
- Product Factory mature: serving, cost tracking, A/B routing

**Month 12+: Evaluate pretraining from scratch**
- Decision gate: has the feedback loop produced measurable compound improvement?
- Is fine-tuning + continued pretraining hitting a ceiling our data could break?
- Does the capital make sense given inference economics?
- Only proceed if the answer to all three is yes

---

## What remains true from the original plan

- The cron-based 5-minute cycle and improvement loop pattern are solid
- The escalation framework and CEO delegation model are solid
- Wayback Machine as a data source is genuinely differentiated and worth pursuing
- The Android companion spec is unchanged
- The recursive improvement architecture is the real moat — the reviewer confirmed this
- The mystery-register press materials are appropriate for current stage
- The six-level recursive hierarchy (execution → self-improvement → meta-improvement → research) remains the long-term vision — it's just preceded by more careful phase-gating

---

## Summary of changes

| Original plan | Revised plan |
|---|---|
| Train from scratch (Phase 1) | Fine-tune open model (Phase 1) |
| Invented percentages | Named benchmarks, baselines first |
| "Exponential improvement" | Honest S-curve with plateau planning |
| 45B tokens = strong model | Quality-per-token over volume |
| Training-focused | Training + Serving equally important |
| No alignment/safety layer | Alignment evaluations required before deploy |
| No experiment tracking | W&B/MLflow from day one |
| No data governance | Full lineage records, legal review |
| No financial loop | Cost and ROI tracked per experiment |
| Five-layer agent hierarchy | Six factories + Emily Core |
| Single improvement loop | Factory-level specialization |
| Pretraining is the goal | Pretraining is a Phase 4 option |

---

## What the reviewer got slightly wrong (or we'd push back on)

**"Fully autonomous AI company is not a solved problem"** — True, and we agree. But the design includes human oversight: bounded cycles, full audit logging, explicit escalation for legal/budget/strategy, and the Android app for real-time human decisions. This is "AI that dramatically reduces workload" not "AI that runs without humans." The original documents had this right; the summary language about "no manual effort" was imprecise.

**The six factories are implicit, not absent** — Knowledge Factory, Training Factory, and parts of Improvement Factory were described in detail. What was genuinely missing: Product Factory (inference/serving), formal Evaluation Factory with named benchmarks, alignment layer, data governance, and experiment tracking. Those gaps are now filled.

**The core insight is validated** — "The strongest idea in your document is not actually the model training. It's the recursive improvement architecture." That's what Emily is. The training capability is downstream of the orchestration, observation, and loop-closing that make it genuinely self-improving.
