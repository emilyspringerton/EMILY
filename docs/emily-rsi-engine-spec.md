# Emily RSI Engine: Ground Truth Specification
## The Canonical Bridge from Design to Implementation

**This document supersedes conflicting details in earlier docs. Read it first.**

---

## Part 1: What Is Actually Built

### The Working System (Today)

`emily-agent/main.go` implements Emily v0:

```
┌─────────────────────────────────────────────┐
│  EMILY v0 (emily-agent/main.go)             │
│                                             │
│  ├─ LLM client (OpenAI + Anthropic, swap   │
│  │    via env)                              │
│  ├─ Agentic tool loop (up to MAX_TOOL_ITERS │
│  │    before producing final response)      │
│  ├─ Tools: git_list_files, git_read_file   │
│  ├─ Two-pass hallucination validation       │
│  ├─ Git-backed JSONL conversation history  │
│  ├─ Token-bucket rate limiter per IP       │
│  └─ Embedded web chat UI                   │
└─────────────────────────────────────────────┘
```

This is the seed. Every component of the RSI system grows from it.

### The Gap

The docs describe Emily in her mature form. What's missing between v0 and maturity:

| Described | Status |
|---|---|
| 5-minute cron cycle (observe → decide → act → plan) | Not built |
| Iterative improvement loop engine | Not built |
| Acceptance criteria evaluation | Not built |
| Specialist agents (Bob, DataQuality, Storage) | Not built |
| Persistent evolution roadmap | Not built |
| State between cron cycles | Not built |
| Data collection pipelines | Not built |
| Model training pipeline | Not built |
| Android escalation companion | Not built |

The path from v0 to the full system is the RSI loop itself. We use the loop to build the loop.

---

## Part 2: Architecture Reconciliation

### What the v2 Revision Got Right (Adopt These)

The `emily-v2-architecture-revision.md` corrects the original design. These corrections are now canonical:

1. **Fine-tune, don't pretrain from scratch** — Phase 1 is fine-tuning Llama 3 / Mistral with QLoRA. Pretraining from scratch is Phase 4 (earn it).

2. **No invented percentages** — The original "66% → 78% → 85%" are made-up figures. Use named benchmarks (MMLU, HumanEval, GSM8K, TruthfulQA). Set targets after seeing a baseline.

3. **Quality-per-token over volume** — A pipeline that collects 500GB and discards 80% is working correctly. The primary metric is quality per retained token, not total tokens.

4. **Honest improvement curve** — Improvement follows an S-curve with plateau planning, not "exponential." Budget for diminishing returns; use the plateau as a signal to change what you're doing.

5. **Experiment tracking from day one** — W&B or MLflow on every run: hyperparameters, dataset version, loss curves, eval results, git hash, checkpoint locations. Without this the loop is blind.

6. **Alignment evaluations are not optional** — Harmful content rate, jailbreak resistance, hallucination rate on known facts, bias evals. Required before any model reaches production.

7. **Product Factory matters** — A model no one can use is just a checkpoint. vLLM/TGI serving, cost-per-inference tracking, A/B routing, and canary deploys are first-class concerns.

### What Remains True from All Earlier Docs

- Cron-based 5-minute cycle with observe/decide/act/plan phases is solid
- The escalation framework (CEO delegation model) is solid
- Wayback Machine as a differentiated data source is valid and worth pursuing
- The Android companion spec is unchanged and still needed
- The recursive improvement architecture is the competitive moat — confirmed by the external reviewer
- The six-level recursive hierarchy (execution → self-improvement → meta-improvement → research) is the long-term vision — just phase-gated properly

---

## Part 3: The RSI Engine — Concrete Specification

### What the RSI Engine Is

The RSI engine is the component that runs iterative improvement loops. It is described in `emily-iterative-loop-execution.md` as pseudocode. This document specifies the real implementation.

### Core Abstraction

```
ImprovementTask
├── Description (what we're improving)
├── AcceptanceCriteria (array, all must pass)
├── MaxIterations (hard limit)
└── Iterations (history, grows each run)

Each Iteration:
├── Generate: call Emily (via Pipeline) to produce an artifact
├── Evaluate: call Emily to measure artifact against criteria
├── Record: append to history
└── Loop or exit

Exit conditions:
├── AllCriteriaPASS → success
├── MaxIterations exhausted → partial success
├── Safety violation detected → halt + escalate
└── Cost/time budget exceeded → halt
```

### The Evaluation Model

The key insight: **Emily evaluates herself.** The measurement function calls the LLM with:
- The artifact
- The criteria
- A structured response schema

This makes the evaluation recursive: Emily generating → Emily measuring → Emily improving. No external oracle needed for most tasks.

For tasks where LLM evaluation is insufficient (actual code execution, benchmark scores), the evaluation function can run shell commands and parse output.

### The Meta-Loop

The RSI engine can be tasked with improving itself:

```
Task: "Improve the RSI engine's iteration efficiency"
Criteria:
  - Average iterations to convergence < 4 (currently ~6)
  - All existing test tasks still pass
  - No degradation in quality of generated artifacts

Run: RSI loop on the RSI engine's own prompts and decision logic
Result: Better system prompts, better analysis format, better evaluation queries
```

This is the seed of recursive self-improvement. The engine iterates on its own internals.

---

## Part 4: Implementation Plan

### Iteration 1: RSI Engine v0 (This PR)

**What we build:**
- `/rsi` HTTP endpoint on the Emily server
- `AcceptanceCriterion`, `ImprovementTask`, `IterationRecord` types
- `RSILoop.Run()` — the core loop
- LLM-based evaluation (Emily measures herself)
- JSON persistence of task state

**Acceptance criteria for iteration 1:**
- POST /rsi accepts task description + criteria + max_iters
- Loop runs up to max_iters
- Each iteration generates an artifact and evaluates it
- Loop exits early when all criteria pass
- Full iteration history returned in response
- No panics on malformed input

**What iteration 1 does NOT include:**
- Shell command execution in evaluation
- Cron scheduling
- Persistent roadmap
- Specialist agents

### Iteration 2: RSI Engine v1 (Next)

**What we build (using the RSI engine to improve itself):**
- Cron-based execution cycle (5-minute trigger)
- State persistence between cron runs (YAML state file)
- Evolution roadmap (YAML, Emily reads and updates)
- Escape hatch commands (`emily-cli pause`, `emily-cli set-priority`)

**Acceptance criteria for iteration 2:**
- Cron runs every 5 minutes
- State loads correctly from previous run
- Emily can read her roadmap and pick next task
- Human can pause/redirect via CLI

### Iteration 3: Specialist Agents

**What we build:**
- Bob (database admin) as a separate HTTP service
- Agent-to-agent communication protocol (structured JSON messages)
- DataQualityAgent
- StorageAgent

### Iteration 4: Data Collection

**What we build:**
- Reddit collection pipeline
- Wikipedia processor
- Quality validation and scoring
- S3-backed storage

### Iteration 5: Fine-tuning Pipeline

**What we build:**
- Dataset preparation (tokenization, deduplication)
- QLoRA fine-tuning on Llama 3 / Mistral
- Evaluation suite (MMLU, HumanEval, GSM8K, TruthfulQA)
- W&B experiment tracking integration

---

## Part 5: The RSI Bootstrap Sequence

This is the recursive part: using each layer to build the next.

```
Day 0: Emily v0 (existing Go chat agent)
  └─ Emily talks to users, reads her conversation history

Day 1: RSI Engine v0 (iteration 1 above)
  └─ Emily can now run improvement loops on any task

Day 7: Emily uses RSI engine to improve RSI engine
  └─ Task: "Improve iteration efficiency"
  └─ Criteria: convergence speed, quality, reliability
  └─ Result: Better prompts, better evaluation, better loop

Day 14: Emily uses improved RSI engine to build cron system
  └─ Task: "Build 5-minute cron cycle"
  └─ Criteria: runs reliably, loads state, picks tasks
  └─ Result: Emily runs autonomously every 5 minutes

Day 21: Emily's cron system builds specialist agents
  └─ Each agent is an improvement task run through the RSI loop
  └─ Bob is built by Emily, using the loop

Day 30: Emily uses Bob to help improve data collection
  └─ Data collection pipeline is an RSI task
  └─ Bob ensures data is stored correctly

Month 3: Emily fine-tunes Llama 3 on collected data
  └─ Fine-tuning is an RSI task (hyperparameter optimization)
  └─ Evaluation suite gives measurable criteria

Month 4: Fine-tuned Emily helps improve data collection
  └─ The loop closes: model improves data → data improves model
```

---

## Part 6: What "RSI of RSI" Means Operationally

The phrase "using recursive self improvement to build recursive self improvement" has a concrete meaning:

1. **The RSI engine is the first subject of the RSI loop**  
   The initial prompt engineering, iteration strategy, and evaluation logic are themselves improved via improvement loops. We run the loop on the loop.

2. **Each component is built by the component before it**  
   The cron system is built using the RSI engine. Specialist agents are built by Emily running cron cycles. Data pipelines are built by specialist agents.

3. **The system accumulates knowledge about its own improvement process**  
   Every iteration record, lesson extracted, and pattern identified becomes context for future iterations. Emily gets better at getting better.

4. **The acceleration is real, but bounded**  
   Improvement follows an S-curve per component. The acceleration comes from compounding across components, not from any single component improving without limit.

---

## Part 7: Acceptance Criteria for the Full System

The RSI system is working when:

### Reliability
- RSI loop runs to completion without crashes: 100%
- Iteration convergence rate (tasks that reach all-pass): > 70% within max_iters
- False positive evaluation rate (passes when shouldn't): < 10%

### Productivity  
- Average iterations to convergence: < 5
- Time from task submission to completion: < 10 minutes for simple tasks
- Lessons extracted and reused in subsequent tasks: measurable

### Safety
- No task produces unsafe output (validated by evaluation layer)
- All iteration history preserved and auditable
- Human can pause system at any point

### Recursive improvement signal
- RSI engine performance improves measurably after self-improvement runs
- Improvement velocity increases (fewer iterations needed) over time
- Lessons library grows and reduces duplicate work

---

## Appendix: Document Hierarchy

For any conflict between documents, this is the precedence order:

1. `emily-rsi-engine-spec.md` (this document) — authoritative on architecture
2. `emily-v2-architecture-revision.md` — authoritative on ML strategy
3. `emily-agent-framework.md` — authoritative on agent hierarchy
4. `emily-iterative-loop-execution.md` — authoritative on loop mechanics
5. `emily-cron-autonomous-evolution.md` — authoritative on cron patterns
6. All other docs — background and context

When earlier docs conflict with higher-precedence docs, the higher-precedence document wins.
