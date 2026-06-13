# API Key Unlock Sequence

*One-time ops doc. Run this when ANTHROPIC_API_KEY is first provisioned.*

---

## What's blocked without the key

Emily Prime's 5-minute cron runs in degraded mode without `ANTHROPIC_API_KEY`:

| Component | Without key | With key |
|-----------|-------------|----------|
| emily-agent HTTP server | ✓ starts, all endpoints up | ✓ same |
| Emiree witch engine | ✓ runs (no LLM) | ✓ same |
| GoldenDocCompiler | writes placeholder sections | compresses via haiku |
| FABLE advisor | returns error | sprint recommendations |
| RSI task generation | returns error | haiku task generation |
| HEIMDAL translation | returns error | haiku requirement translation |
| `/chat` endpoint | returns error | full conversation |
| `/api/v1/emily/plan` | returns error | structured sprint planning |
| `emily context build` | writes raw context (large) | compressed bilingual context |
| `emily backlog promote` | heuristic-only mode | full haiku classify + route |

**Cost floor:** All LLM calls use `claude-haiku-4-5-20251001` (MODEL env default).
`SonnetModel` is defined in config but unused — zero sonnet/opus spend at current scale.

---

## Unlock Sequence

### 1. Set the key

Add to the emily-agent's environment. Preferred: systemd drop-in or `.env` file at the agent root.

```bash
# Option A: export in shell before emily start
export ANTHROPIC_API_KEY=sk-ant-...

# Option B: add to EMILY/emily-agent/.env (if your start script sources it)
echo 'ANTHROPIC_API_KEY=sk-ant-...' >> /home/fatbaby/EMILY/emily-agent/.env

# Option C: systemd override (if running as a service)
systemctl --user edit emily-agent.service
# Add: [Service]
#      Environment=ANTHROPIC_API_KEY=sk-ant-...
```

### 2. Restart emily-agent

```bash
emily start   # idempotent — kills existing daemon and starts fresh
# or if running as systemd:
systemctl --user restart emily-agent.service
```

### 3. Compile full system context for the first time

```bash
emily context build
# With key: compresses all 18 Tier 1 golden sources via haiku → full-system-context.md
# Expect: "Sources compiled: 18/18", one haiku call per source (~18 calls total)
# Cached after first run — subsequent builds only recompress changed sources
```

### 4. Run the multilingual compression A/B test (S22-13)

```bash
bash /home/fatbaby/EMILY/scripts/compression-abtest.sh
# Compares bilingual vs English-only GOLDEN.md comprehension via haiku
# If score ≥ 95%: deploy bilingual version (saves ~30-40% tokens per context call)
# If score < 95%: keep English-only, note result
```

### 5. Verify FABLE advice endpoint

```bash
curl -s http://localhost:8086/api/v1/emily/fable/advice | jq '.recommendations[].title'
# Should return 3 sprint recommendations based on full-system-context.md
# If emily-agent isn't running yet: emily start first
```

### 6. Verify the cron cycle fires

Check the log after ~5 minutes:

```bash
tail -f /home/fatbaby/EMILY/var/logs/emily-agent.log | grep -E "goldenbuild|cron|RSI|apple"
# Expect: "goldenbuild: loaded 18 Tier 1 sources from index"
# Expect: "cron: cycle N starting"
# Expect: Apple filed to IDUNA
```

---

## Troubleshooting

**"ANTHROPIC_API_KEY is not set" in logs but key is in env**
→ emily-agent was started before the env var was set. Restart: `emily start`.

**`emily context build` slow on first run**
→ Normal — 18 haiku calls serially (~30s total). Subsequent runs reuse cache for unchanged sources.

**FABLE returns "no API key"**
→ emily-agent is still running the old binary. Check PID: `emily status`. Restart if needed.

**Compression A/B test fails**
→ Check `EMILY/docs/compression-experiment/` for the test harness and logs. The test requires
the key to be set and the bilingual test file to exist at `docs/compression-experiment/GOLDEN_BILINGUAL_TEST.md`.

---

*Written 2026-06-13. Update when key rotation, model upgrades, or new unlock steps are added.*
