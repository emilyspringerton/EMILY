# Emily Agent Communication Protocol
## Canonical wire format for inter-agent messages

**Produced by:** RSI Loop 3 — "Design specialist agent communication protocol"
**Acceptance criteria satisfied:** typed envelope, structured body, routing, audit, versioning

---

## Design Constraints

- Emily (orchestrator) and agents (Bob, DataQuality, Storage, etc.) run as separate processes.
- Messages are HTTP POST to each agent's `/task` endpoint.
- Every message is logged before it's sent and after the result is received. No action without a trace.
- Agents are stateless workers. Emily holds all state. Agents read task, execute, return result.
- The protocol must be extensible without breaking existing agents (versioned envelope).

---

## Message Envelope

All messages use this JSON envelope:

```json
{
  "v":          1,
  "id":         "msg_<unix_nano>",
  "from":       "emily | bob | dataquality | storage | ...",
  "to":         "emily | bob | dataquality | storage | ...",
  "type":       "task | result | escalation | heartbeat",
  "priority":   "critical | high | normal | low",
  "sent_at":    "2026-05-30T14:00:00Z",
  "expires_at": "2026-05-30T14:05:00Z",
  "body":       { ... type-specific payload ... },
  "trace_id":   "trace_<id>",
  "reply_to":   "msg_<id of task this replies to, if type=result>"
}
```

**Field rules:**
- `v`: always 1 for now. Increment only on breaking changes.
- `id`: globally unique. Format: `msg_<unix nanoseconds>`.
- `type`: determines how `body` is parsed.
- `expires_at`: agent must not start task after this time. Emily sets it to `cycle_deadline - 10s`.
- `trace_id`: groups all messages in one Emily cycle. Used for audit log correlation.
- `reply_to`: present on `result` and `escalation` messages only.

---

## Body Schemas

### type: task

```json
{
  "task_id":     "task_<id>",
  "description": "Validate quality of the last Reddit batch",
  "inputs": {
    "data_path": "s3://emily-data/reddit/2026-05-30.jsonl",
    "expected_schema": ["id", "title", "body", "score", "subreddit"]
  },
  "success_criteria": [
    {"name": "schema_valid",  "target": "100% records pass schema"},
    {"name": "no_empty_body", "target": "< 1% records with empty body"},
    {"name": "score_range",   "target": "all scores >= 0"}
  ],
  "timeout_seconds": 120,
  "requires_confirm": false
}
```

**`requires_confirm`**: if `true`, the agent must not execute until it sends back a `result` with `status: awaiting_confirm` and Emily replies with a `task` containing `{"confirm": true}`. Used for destructive operations (DB drops, bulk deletes).

### type: result

```json
{
  "task_id":    "task_<id>",
  "status":     "success | partial | failed | awaiting_confirm",
  "outputs": {
    "records_checked": 10420,
    "schema_valid_pct": 99.8,
    "empty_body_pct": 0.3,
    "score_range_ok": true
  },
  "criteria_results": [
    {"name": "schema_valid",  "passes": true,  "value": "99.8%",  "gap": ""},
    {"name": "no_empty_body", "passes": true,  "value": "0.3%",   "gap": ""},
    {"name": "score_range",   "passes": true,  "value": "all ok", "gap": ""}
  ],
  "all_pass":   true,
  "error":      "",
  "duration_ms": 4230,
  "artifacts":  ["s3://emily-data/reddit/2026-05-30-quality-report.json"]
}
```

### type: escalation

```json
{
  "task_id":     "task_<id>",
  "category":    "legal | budget | strategy | values | novel | emergency | safety",
  "title":       "License ambiguity: cannot confirm commercial use of source X",
  "summary":     "Attempting to add source X to corpus. ToS allows research but not commercial use. Need legal review.",
  "context": {
    "what_happened":  "...",
    "why_escalated":  "Crosses legal boundary; outside agent authority.",
    "what_is_blocked": "Corpus expansion (source X only)"
  },
  "options": [
    {"id": "approve",  "label": "Approve inclusion",   "requires_confirm": true},
    {"id": "exclude",  "label": "Exclude source X",     "requires_confirm": false}
  ],
  "recommended": "exclude",
  "deadline":    null
}
```

### type: heartbeat

```json
{
  "status":        "healthy | degraded | down",
  "uptime_seconds": 3600,
  "tasks_pending":  2,
  "tasks_complete": 47,
  "last_error":     ""
}
```

---

## Agent HTTP Interface

Every agent exposes these endpoints:

```
POST /task          Accept a task message (envelope + body type=task)
                    Returns 202 Accepted immediately
                    Result delivered async via Emily's POST /result endpoint

GET  /task/{id}     Check status of a running task
                    Returns result envelope if complete, or {"status":"running"}

GET  /health        Returns heartbeat body (no envelope)

POST /confirm/{id}  Deliver confirmation for awaiting_confirm tasks
```

Emily's server exposes:

```
POST /result        Agents POST their completed results here
POST /escalation    Agents POST escalations here (triggers Android push)
```

---

## Delivery Guarantees

**At-least-once:** Emily retries failed deliveries up to 3 times with exponential backoff (1s, 4s, 16s). After 3 failures, the task is marked failed and the roadmap item is re-queued.

**Idempotency:** Tasks carry a `task_id`. If an agent receives the same `task_id` twice, it returns the cached result without re-executing. Agents store the last 100 task results in memory.

**Ordering:** Emily sends at most one task per agent per cycle. Agents process tasks serially. No parallelism within an agent by default.

---

## Audit Log Format

Emily writes one line per message event:

```jsonl
{"event":"sent","trace_id":"trace_42","msg_id":"msg_1717000000000","from":"emily","to":"bob","type":"task","task_id":"task_001","at":"2026-05-30T14:00:00Z"}
{"event":"received","trace_id":"trace_42","msg_id":"msg_1717000000001","from":"bob","to":"emily","type":"result","task_id":"task_001","status":"success","duration_ms":4230,"at":"2026-05-30T14:00:04Z"}
```

Every send and receive is logged before/after the HTTP call. The log is append-only and never deleted.

---

## Implementation Notes

**Bob (database admin agent):**
- Written in Go, separate binary
- Inputs: SQL query or named operation (backup, health_check, migration)
- Outputs: rows (SELECT), affected count (UPDATE/DELETE), or error
- `requires_confirm: true` on all DDL and DELETE operations

**DataQualityAgent:**
- Written in Python (natural fit for data processing)
- Inputs: data path (S3 or local), expected schema, quality thresholds
- Outputs: per-record quality scores, summary statistics, pass/fail per criterion

**StorageAgent:**
- Written in Python
- Inputs: operation (upload, download, move, delete), paths, metadata
- Outputs: s3 URIs, sizes, checksums
- `requires_confirm: true` on delete

**CommunicationAgent (future):**
- Sends escalations to the Android app
- Inputs: escalation body (category, title, summary, options)
- Outputs: acknowledgment, human decision when it arrives

---

## Error Handling

**Task timeout:** If agent does not respond within `timeout_seconds`, Emily marks task failed, increments `consecutive_failures`, and re-queues on the roadmap after a cooldown.

**Agent down:** If `/health` returns non-200 or times out, Emily stops sending tasks to that agent and creates an escalation of category `emergency`.

**Protocol version mismatch:** If `v != 1`, Emily logs a warning and attempts to process anyway. If the body cannot be parsed, it is logged and discarded with error.

---

## Adding a New Agent

1. Implement `POST /task`, `GET /task/{id}`, `GET /health`, `POST /confirm/{id}`
2. Register in Emily's agent registry (JSON config, not hardcoded)
3. Emily begins health-checking on the next cron cycle
4. Add roadmap items with `agent: <agent_name>` in inputs to route tasks

No code changes to Emily needed for a new agent — only config.
