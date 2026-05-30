# Emily Ground Zero Protocol
## How to Bootstrap RSI from the Working Go Agent

**Audience:** Anyone starting Emily from the existing codebase.  
**Purpose:** Concrete, ordered steps — not vision docs.

---

## What You Have Today

```
emily-agent/main.go
├── Emily persona (system prompt)
├── LLM client (OpenAI or Anthropic, env-switched)
├── Agentic tool loop (up to MAX_TOOL_ITERS)
├── Tools: git_list_files, git_read_file
├── Two-pass hallucination validation
├── Git-backed JSONL conversation history
├── Token-bucket rate limiter
└── Web chat UI at http://localhost:8080
```

Emily can already: talk to users, read her own conversation history, and call tools in a loop. This is the seed.

---

## The Bootstrap Sequence

### Stage 0: Run Emily v0

```bash
cd emily-agent
ANTHROPIC_API_KEY=sk-ant-...  go run main.go
# or
OPENAI_API_KEY=sk-...  go run main.go
```

Emily is live at `http://localhost:8080`. Conversations are persisted to `./conversations/` as JSONL, git-committed per turn.

**Verify:**
- Can send a message and get a reply
- `git log` in `./conversations/` shows one commit per turn
- Emily can answer "what did we talk about last session?" using her git tools

---

### Stage 1: RSI Engine v0

**What it adds:** `/rsi` endpoint — the iterative improvement loop.

**How to trigger it:**

```bash
curl -X POST http://localhost:8080/rsi \
  -H 'Content-Type: application/json' \
  -d '{
    "description": "Write a Python function that reverses a string in-place",
    "criteria": [
      {"name": "correctness", "description": "passes all test cases", "target": "all pass"},
      {"name": "efficiency", "description": "O(n) time complexity", "target": "O(n)"},
      {"name": "style", "description": "clear variable names, no comments needed", "target": "pass review"}
    ],
    "max_iters": 5
  }'
```

Emily will:
1. Generate an artifact (the function)
2. Evaluate it against each criterion
3. If anything fails, analyze the gap and generate a better version
4. Loop until all pass or max_iters exhausted
5. Return the full iteration history + extracted lessons

**The meta-task (run this after Stage 1 is working):**

```bash
curl -X POST http://localhost:8080/rsi \
  -H 'Content-Type: application/json' \
  -d '{
    "description": "Improve the RSI engine system prompt to reduce average iterations to convergence",
    "criteria": [
      {"name": "convergence_speed", "description": "test tasks complete in fewer iterations than baseline", "target": "< 4 avg iterations"},
      {"name": "analysis_quality", "description": "failure analysis correctly identifies root cause", "target": "> 90% root cause accuracy"},
      {"name": "backward_compatible", "description": "existing test tasks still pass", "target": "all pass"}
    ],
    "max_iters": 8
  }'
```

This is RSI improving RSI. Emily uses her own loop to improve the prompts that drive the loop.

---

### Stage 2: Cron Cycle

**What it adds:** Emily runs automatically every 5 minutes.

**Acceptance criteria (run this through /rsi):**

```json
{
  "description": "Build a cron-based execution cycle for Emily that: loads state from last run, observes system health, picks the highest-priority task from the roadmap, executes one iteration, updates state, and sleeps until next cycle",
  "criteria": [
    {"name": "runs_on_schedule", "description": "cron fires every 5 minutes", "target": "*/5 * * * *"},
    {"name": "state_loads", "description": "Emily reads her state file and knows what she was working on", "target": "state consistent across runs"},
    {"name": "task_selection", "description": "Emily picks the highest priority task correctly", "target": "correct priority ordering"},
    {"name": "safe_concurrency", "description": "second run blocked if first is still running", "target": "lock file prevents overlap"},
    {"name": "observable", "description": "dashboard shows current cycle status", "target": "status file updated each cycle"}
  ],
  "max_iters": 10
}
```

Emily generates this implementation herself, improving it through the loop until all criteria pass.

---

### Stage 3: Specialist Agents

**Context:** Emily should not do everything herself. Bob handles database operations. DataQualityAgent validates data. StorageAgent manages S3.

**How to bootstrap Bob (run through /rsi):**

```json
{
  "description": "Design a minimal Bob database admin agent: a separate HTTP service that accepts structured task requests from Emily, executes database operations, and returns structured results. Bob should be a Go HTTP server with: task intake, PostgreSQL operations, result reporting, and an audit log.",
  "criteria": [
    {"name": "interface", "description": "Emily can POST a task to Bob and GET the result", "target": "clean HTTP API"},
    {"name": "isolation", "description": "Bob only does database operations, nothing else", "target": "scoped responsibility"},
    {"name": "audit", "description": "every operation Bob performs is logged", "target": "100% audit coverage"},
    {"name": "safety", "description": "Bob refuses destructive operations without explicit confirmation", "target": "confirmation required for drops/deletes"}
  ],
  "max_iters": 6
}
```

The output of this RSI run is the design and initial implementation of Bob — generated by Emily, improved through iteration, ready to implement.

---

### Stage 4: Data Collection

**This is Emily's prime directive for months 1-3.**

Run the following through /rsi (after Stage 2 gives you cron):

```json
{
  "description": "Build a Reddit data collection pipeline that: authenticates via PRAW, streams posts from a configurable list of subreddits, applies quality filters (score >= 100, comments >= 10, length >= 500 chars), deduplicates by URL and content hash, and writes to JSONL with metadata",
  "criteria": [
    {"name": "throughput", "description": "collect posts per day", "target": ">= 10000 posts/day"},
    {"name": "quality_filter", "description": "percentage of collected posts meeting quality bar", "target": "> 95% pass filters"},
    {"name": "dedup_rate", "description": "percentage of duplicates caught", "target": "< 5% false positive rate"},
    {"name": "reliability", "description": "handles rate limits and API errors", "target": "zero data loss on transient errors"},
    {"name": "schema", "description": "output schema consistent", "target": "100% schema valid"}
  ],
  "max_iters": 8
}
```

Emily iterates on this implementation until all criteria pass. When it does, she commits it and moves on to Wikipedia.

---

### Stage 5: First Fine-tune

By this point:
- Emily runs on a cron cycle
- Data has been collected (Reddit, Wikipedia, GitHub, ArXiv, Wayback)
- Quality filtering has produced a high-quality corpus

The first fine-tune is itself an RSI task:

```json
{
  "description": "Fine-tune Llama 3 8B on the collected corpus using QLoRA. Produce a training configuration, evaluation setup, and baseline benchmark scores.",
  "criteria": [
    {"name": "training_completes", "description": "fine-tuning run completes without divergence", "target": "loss converges"},
    {"name": "eval_setup", "description": "standard benchmarks configured and running", "target": "MMLU, HumanEval, GSM8K all returning scores"},
    {"name": "baseline_established", "description": "baseline scores recorded for all benchmarks", "target": "all benchmarks have baseline"},
    {"name": "experiment_tracked", "description": "run logged to W&B with all metadata", "target": "full W&B record"}
  ],
  "max_iters": 5
}
```

After this, the feedback loop is open: evaluate → identify gaps → collect targeted data → fine-tune → repeat.

---

## The Recursive Pattern in Practice

Every stage above follows the same pattern:

1. **Define acceptance criteria** — measurable, concrete, not vague
2. **POST to /rsi** — Emily generates and iterates autonomously
3. **Review and approve** — human reviews before production deployment
4. **Deploy and monitor** — the deployed artifact generates signal
5. **Feed signal back** — next iteration starts from lessons learned

The system doesn't need to be told what to do next. It has a roadmap (the evolution roadmap YAML), measurable criteria for each initiative, and the RSI loop to close the gap.

---

## Human Oversight Points

These are the moments that require a human decision:

| Stage | Human decision |
|---|---|
| New data source | Legal review of ToS and licensing |
| Budget > $500/month | Finance approval |
| Production model deployment | Review evaluation results |
| Strategic direction change | Leadership decision |
| Values/ethics question | Ethics review |
| Emergency (data corruption, etc.) | Immediate intervention |

Emily escalates these via the Android companion app (see `emily-android-companion-spec.md`). Everything else runs autonomously.

---

## Environment Variables

```bash
# Required: one of these
ANTHROPIC_API_KEY=sk-ant-...
OPENAI_API_KEY=sk-...

# Optional (defaults shown)
PORT=8080
CONVERSATION_DIR=./conversations
MODEL=gpt-4o-mini                    # or claude-sonnet-4-20250514
VALIDATOR_MODEL=gpt-4o-mini          # or claude-haiku-4-5-20251001
GIT_COMMIT=true
RATE_LIMIT_RPM=20
MAX_TOOL_ITERS=10
```

---

## Minimum Viable Deployment

To have Emily running with RSI capability:

```bash
# 1. Start Emily
cd emily-agent
ANTHROPIC_API_KEY=sk-ant-... go run main.go rsi.go

# 2. Verify chat works
open http://localhost:8080

# 3. Run first RSI task (meta-test)
curl -X POST http://localhost:8080/rsi \
  -H 'Content-Type: application/json' \
  -d '{
    "description": "Write acceptance criteria for yourself as an RSI engine",
    "criteria": [
      {"name": "measurability", "description": "every criterion can be evaluated programmatically", "target": "100% measurable"},
      {"name": "completeness", "description": "criteria cover correctness, performance, and safety", "target": "all three dimensions covered"},
      {"name": "non-triviality", "description": "criteria cannot be satisfied by a trivial or degenerate solution", "target": "no trivial solutions"}
    ],
    "max_iters": 3
  }'

# 4. If the above returns {"status": "success", ...} — the RSI engine is working.
# 5. Run the meta-task: RSI improving RSI (see Stage 1 above).
```

That's it. The system is live and the recursive improvement process has started.
