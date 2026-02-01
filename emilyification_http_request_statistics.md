# 🧠 EMILYIFICATION OF HTTP REQUEST STATISTICS
## A Digital Immune System for Sovereign Infrastructure
**Authoritative Mega Meta Specification (v0.9 – Survival Grade)**

## CORE AXIOMS (NON-NEGOTIABLE)
- **The application must never know this exists.**
  - Observability is parasitic, not symbiotic. The host must live even if the parasite dies.
- **Every byte collected must justify its existence in court.**
  - Legal, technical, moral. If it doesn’t prevent downtime or fund survival, it’s gone.
- **All intelligence is local-first.**
  - No third-party beacons. No SaaS umbilical cords. No “just one script tag.”
- **Fail open. Always.**
  - Logging failure must degrade intelligence, never availability.

---

# PART 1 — THE UNIVERSAL HARVESTER
## Backend Telemetry Without Performance Sin

### 1.1 Architectural Overview
**Mental Model:**
- A write-only nervous system that listens, compresses, and forgets fast.

**Flow (All Environments):**

```
[ Request Enters ]
        |
        v
[ Ultra-thin Interceptor ]
        |
        |-- capture metadata (NO BODY)
        |
        |-- enqueue -> lock-free ring buffer
        |
        v
[ Application Logic Continues ]
        |
        v
[ Response Leaves ]
(Parallel, async)
[ Ring Buffer ] -> [ Aggregator Worker ] -> [ Batched Emit (UDP/Unix Socket) ]
```

**Rules**
- No synchronous I/O.
- No disk writes on the hot path.
- No allocations after warm-up.

### 1.2 Environment-Specific Strategies
#### 🐍 Python (WSGI / ASGI)
**Mechanism**
- Single middleware layer
- Capture at `environ` / `scope` boundary
- Use `time.monotonic_ns()` only
- Push struct into `collections.deque(maxlen=N)`

**Rules**
- No JSON encoding in-request
- No logging module usage
- One function call, one struct write

#### 🐹 Golang (`net/http`)
**Mechanism**
- Thin `http.Handler` wrapper
- Record timestamps + pointer-stable strings
- Send struct into buffered channel
- Dedicated goroutine batches + emits

**Rules**
- Channel must never block (drop if full)
- No mutexes on request path
- p99 added latency ≤ 2ms

#### 🪵 Apache
**Mechanism A — Custom LogFormat**
- Log to named pipe
- External daemon parses + aggregates
- `buffered=true`, `flush=off`

**Mechanism B — `mod_lua`**
- Capture request metadata only
- Push to local UDP socket

**Rules**
- No synchronous disk I/O
- Apache must not notice failure

#### 🌊 Nginx
**Mechanism**
- `lua-nginx-module`
- `log_by_lua_block`
- Emit binary-packed record over Unix socket

**Rules**
- Lua code ≤ 50 LOC
- No string concatenation
- No table growth per request

### 1.3 “Soul of the Request” (Captured Fields)
```json
{
  "ts_ns": 0,
  "method": "GET",
  "path_hash": "u64",
  "status": 200,
  "req_bytes": 512,
  "resp_bytes": 2048,
  "latency_us": 1530,
  "ip_prefix": "203.0.113.0/24",
  "tls": true,
  "http_version": "2"
}
```

**Explicitly Excluded**
- Query params
- Headers (raw)
- Cookies
- Request body
- User IDs

### 1.4 Acceptance Criteria (Emily Standard)
- ✅ Added latency < 2ms p99
- ✅ Zero allocations after warm-up
- ✅ Drops telemetry under pressure, never blocks
- ✅ Can be disabled at runtime without restart

### 1.5 Failure State
If telemetry fails:
- Drop records
- Increment internal counter
- Never log the failure
- Never panic
- Never retry synchronously

---

# PART 2 — THE FINGERPRINT
## Privacy-Preserving Threat Discrimination

### 2.1 The Threat Reality
- IP addresses are weather, not identity.
- Bots rotate. Humans hesitate.

We fingerprint behavior, not people.

### 2.2 Minimum Viable Fingerprint (MVF)
**Collected Signals (Ephemeral)**

| Layer | Signal | Reason |
| --- | --- | --- |
| TLS | JA3 hash | Client stack fingerprint |
| HTTP | Header order hash | Bots lie badly |
| TCP | Window size | Automation artifact |
| Timing | Inter-request deltas | Humans breathe |
| Method mix | GET/POST ratio | Crawlers skew |

### 2.3 Session Hash Construction
```
session_seed = HMAC(
    rotating_daily_key,
    JA3 ||
    header_order_hash ||
    tcp_window ||
    timing_bucket
)
```
- Rotates every 24h
- Cannot be correlated cross-day
- Useless outside this system

### 2.4 Heuristic Classification
```
if ja3 in known_bad: +50
if header_order == "alphabetical": +20
if delta < 20ms repeatedly: +30
if method skew > 95% GET: +10
if score >= 60 => hostile
```

- No ML.
- Deterministic. Explainable. Auditable.

### 2.5 Acceptance Criteria
- ✅ Selenium bot flagged ≤ 3 requests
- ✅ False positive rate < 0.1%
- ✅ Session hashes unrecoverable after 24h
- ✅ No cross-site correlation possible

### 2.6 Failure State
If fingerprinting fails:
- Treat as unknown
- Apply conservative rate limits
- Never block legitimate traffic preemptively

---

# PART 3 — THE EMILY AD-INFRASTRUCTURE
## Monetization as a Survival Reflex

### 3.1 Core Principle
Ads are pressure valves, not trackers.

The server’s health decides what is shown.

### 3.2 Architecture Flow
```
[ Request ]
|
v
[ Health Evaluator ]
|
+-- under load? --> [ DEFENSE RESPONSE ]
|
+-- healthy? -----> [ CONTEXTUAL AD ]
```

### 3.3 Server Health Inputs
- Load average
- Request error rate
- Telemetry drop rate
- Active hostile sessions

### 3.4 Ad Modes
**🟥 Defensive Mode**
- CAPTCHA
- Static SVG
- Plaintext message
- Zero JS
- Zero images

**🟩 Healthy Mode**
- Vector asset (SVG/WebGPU)
- Contextual copy (page category)
- No cookies
- No user storage

### 3.5 Context Model (Not User Model)
```json
{
  "path_class": "docs|commerce|media",
  "server_health": "green|yellow|red",
  "request_intent": "read|write|browse"
}
```
Ads match intent, not identity.

### 3.6 Acceptance Criteria
- ✅ Ad rendering adds < 1ms
- ✅ Defensive ads reduce hostile RPS ≥ 40%
- ✅ No persistent identifiers stored
- ✅ Revenue covers infra cost at scale

### 3.7 Failure State
If ad system fails:
- Serve nothing
- Never block content
- Never cascade into request failure

---

# PART 4 — THE GOLDEN RECORD
## Unified, Minimal, Defensible

### 4.1 Canonical Record Schema
```json
{
  "ts": 0,
  "path_hash": "u64",
  "latency_us": 0,
  "status": 200,
  "session_class": "human|bot|unknown",
  "health_state": "green|yellow|red",
  "action_taken": "allow|throttle|challenge"
}
```

**Retention:**
- Raw: 24h
- Aggregates: 30d
- Nothing permanent

### 4.2 System-Wide Acceptance Criteria
- ✅ Zero external dependencies
- ✅ Operates during partial outages
- ✅ Can be audited line-by-line
- ✅ Cannot be repurposed for surveillance

### 4.3 Absolute Failure Doctrine
The application must survive even if this entire system is on fire.

- No hard dependencies
- No startup coupling
- No shutdown hooks
- No retries on request path

---

## FINAL NOTE (Architect to Architect)
This is not analytics.

This is infrastructure immunity.

It observes only to defend.

It monetizes only to survive.

It forgets by design.

If someone asks “where’s the dashboard?”

They’ve already missed the point.

**Proceed to implementation.**
