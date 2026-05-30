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

---

# PART 5 — REFERENCE IMPLEMENTATION (MINIMAL, DROP-SAFE)
## The Parasite That Never Bites the Host

### 5.1 Common Wire Format (Binary, Fixed-Size)
**Goal:** zero allocations, no JSON, stable ABI across languages.

**Record (little-endian, 64-bit aligned):**
```
struct TelemetryRecord {
  u64 ts_ns;
  u64 path_hash;
  u32 method;        // enum
  u16 status;
  u16 http_version;  // enum
  u32 req_bytes;
  u32 resp_bytes;
  u32 latency_us;
  u32 ip_prefix;     // IPv4 /24 packed; IPv6 -> 0
  u8  tls;           // 0/1
  u8  reserved[7];   // pad to 64 bytes
}
```

**Rules**
- Fixed 64 bytes for SIMD-friendly batching.
- Path hashing uses a stable, non-cryptographic 64-bit hash (e.g., xxHash64).
- Method + HTTP version are enums, not strings.
- IPv6 is supported via separate optional record type (not on hot path).

### 5.2 Runtime Controls (Fail-Open)
**Shared knobs**
- `EMILY_TELEMETRY_ENABLED=1|0` (checked once per second by worker)
- `EMILY_TELEMETRY_RATE=0-100` (sample percent, drop others)
- `EMILY_RING_SIZE=N` (power of two)

**Guarantees**
- Disable at runtime without restart.
- Drop on overflow.
- Never log on failure.

### 5.3 Python (WSGI / ASGI) — Minimal Middleware
```python
# emily_telemetry.py
import time
import collections
import threading
import os
import socket
import struct
import xxhash

RING = collections.deque(maxlen=int(os.getenv("EMILY_RING_SIZE", "65536")))
ENABLED = True
LOCK = threading.Lock()  # used only by worker, never on request path

def _hash_path(path: str) -> int:
    return xxhash.xxh64(path, seed=0).intdigest()

def _record(method, path, status, req_bytes, resp_bytes, latency_us, ip_prefix, tls, http_version):
    ts_ns = time.monotonic_ns()
    return struct.pack(
        "<QQIHHIIIIB7x",
        ts_ns,
        _hash_path(path),
        method,
        status,
        http_version,
        req_bytes,
        resp_bytes,
        latency_us,
        ip_prefix,
        1 if tls else 0,
    )

def middleware(app):
    def _inner(environ, start_response):
        start = time.monotonic_ns()
        status_holder = {}

        def _sr(status, headers, exc_info=None):
            status_holder["status"] = int(status.split()[0])
            return start_response(status, headers, exc_info)

        result = app(environ, _sr)
        end = time.monotonic_ns()
        try:
            path = environ.get("PATH_INFO", "/")
            method = _method_enum(environ.get("REQUEST_METHOD", "GET"))
            req_bytes = int(environ.get("CONTENT_LENGTH") or 0)
            resp_bytes = int(environ.get("CONTENT_LENGTH") or 0)
            latency_us = int((end - start) / 1000)
            ip_prefix = _ip_prefix_24(environ.get("REMOTE_ADDR", "0.0.0.0"))
            tls = environ.get("wsgi.url_scheme") == "https"
            http_version = _http_enum(environ.get("SERVER_PROTOCOL", "HTTP/1.1"))
            RING.append(_record(method, path, status_holder.get("status", 0), req_bytes,
                                resp_bytes, latency_us, ip_prefix, tls, http_version))
        except Exception:
            pass  # fail open
        return result
    return _inner

def _worker():
    sock = socket.socket(socket.AF_UNIX, socket.SOCK_DGRAM)
    sock.connect("/var/run/emily.sock")
    while True:
        if not RING:
            time.sleep(0.001)
            continue
        batch = []
        while RING and len(batch) < 256:
            batch.append(RING.popleft())
        try:
            sock.send(b"".join(batch))
        except Exception:
            pass

threading.Thread(target=_worker, daemon=True).start()
```

**Notes**
- The request path does a single `deque.append` and no logging.
- Any exception is swallowed on the hot path.

### 5.4 Golang (net/http) — Zero-Block Handler
```go
// emily_telemetry.go
type TelemetryRecord struct {
	TsNs       uint64
	PathHash   uint64
	Method     uint32
	Status     uint16
	HttpVer    uint16
	ReqBytes   uint32
	RespBytes  uint32
	LatencyUs  uint32
	IpPrefix   uint32
	Tls        uint8
	_          [7]byte
}

var ring = make(chan TelemetryRecord, 65536)

func Telemetry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rw, r)
		rec := TelemetryRecord{
			TsNs:      uint64(time.Now().UnixNano()),
			PathHash:  xxhash.Sum64String(r.URL.Path),
			Method:    methodEnum(r.Method),
			Status:    uint16(rw.status),
			HttpVer:   httpEnum(r.Proto),
			ReqBytes:  uint32(r.ContentLength),
			RespBytes: uint32(rw.bytes),
			LatencyUs: uint32(time.Since(start).Microseconds()),
			IpPrefix:  ipPrefix24(r.RemoteAddr),
			Tls:       boolToByte(r.TLS != nil),
		}
		select {
		case ring <- rec:
		default:
			// drop on pressure
		}
	})
}

func telemetryWorker(sockPath string) {
	conn, _ := net.Dial("unixgram", sockPath)
	buf := make([]byte, 0, 256*64)
	for {
		buf = buf[:0]
		for i := 0; i < 256; i++ {
			rec := <-ring
			buf = append(buf, (*(*[64]byte)(unsafe.Pointer(&rec)))[:]...)
		}
		conn.Write(buf)
	}
}
```

### 5.5 Nginx (lua-nginx-module)
```nginx
log_by_lua_block {
  local ts = ngx.now() * 1000000000
  local path = ngx.var.uri
  local status = ngx.status
  local req = tonumber(ngx.var.request_length) or 0
  local resp = tonumber(ngx.var.bytes_sent) or 0
  local latency = tonumber(ngx.var.request_time) * 1000000
  local tls = ngx.var.https == "on" and 1 or 0
  local record = emily.pack(ts, path, status, req, resp, latency, tls, ngx.req.http_version())
  emily.send(record)
}
```

**Rules**
- `emily.pack` must use a preallocated buffer.
- `emily.send` must be non-blocking UDP or unixgram.

### 5.6 Apache (Named Pipe)
**LogFormat**
```
LogFormat "%t %m %U %>s %{req_bytes}n %{resp_bytes}n %D" emily
CustomLog "|/usr/local/bin/emily_pipe" emily
```

**External daemon**
- Parse tokens into fixed-size record.
- Emit over UDP/unix socket.
- Drop on backpressure.

### 5.7 Aggregator Worker (All Envs)
**Responsibilities**
- Batch 256 records at a time.
- Optional: aggregate counts into 1-second windows.
- Emit over UDP or unixgram.
- Never write to disk in hot path.

### 5.8 On-Host Collector (Local-First)
**Inputs**
- UDP/unix socket from workers.

**Outputs**
- In-memory aggregates (24h rolling)
- Optional: rotate daily snapshots to local disk *outside* request path

**Never**
- Call external services
- Expose raw IPs
- Persist beyond policy
