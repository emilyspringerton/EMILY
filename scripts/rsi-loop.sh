#!/bin/bash
# rsi-loop.sh — EMILY OS Recursive Self-Improvement Master Loop
#
# The tic-toc engine. Each cycle:
#   TIC:     Ask Emily Prime (dispatch RSI token-efficiency analysis task via emily CLI)
#   TOCK:    Claude Code acts on the task (observation-watcher → claude)
#   ENTROPY: Tyler runs N build iterations (creative entropy injection)
#   ANALYZE: Emily Prime observes what was built → feeds next TIC
#   REPEAT:  while true
#
# Usage:
#   ./scripts/rsi-loop.sh                    # infinite loop
#   ./scripts/rsi-loop.sh 5                  # 5 iterations then exit
#   TYLER_ITERS=3 ./scripts/rsi-loop.sh      # 3 Tyler builds per cycle
#   SLEEP_BETWEEN=60 ./scripts/rsi-loop.sh   # 60s between cycles
#   DRY_RUN=1 ./scripts/rsi-loop.sh          # preview without firing tasks
#
# Environment:
#   MAX_ITERS       max iterations (0 = infinite, default: arg 1 or 0)
#   TYLER_ITERS     Tyler builds per entropy cycle (default: 2)
#   SLEEP_BETWEEN   seconds between RSI cycles (default: 30)
#   DRY_RUN         set to 1 to preview without writing tasks
#   SKIP_TYLER      set to 1 to skip Tyler entropy (useful for FatBaby-only RSI)
#   SKIP_WAIT       set to 1 to skip waiting for Claude completion (fire-and-forget)
#
# Prerequisites:
#   emily CLI installed (emily --version)
#   ANTHROPIC_API_KEY set in environment
#   observation-watcher running (emily start, or emily start --iduna)

set -euo pipefail

# ── Config ────────────────────────────────────────────────────────────────────
MAX_ITERS="${1:-${MAX_ITERS:-0}}"
TYLER_ITERS="${TYLER_ITERS:-2}"
SLEEP_BETWEEN="${SLEEP_BETWEEN:-30}"
DRY_RUN="${DRY_RUN:-0}"
SKIP_TYLER="${SKIP_TYLER:-0}"
SKIP_WAIT="${SKIP_WAIT:-0}"

EMILY_ROOT="${EMILY_ROOT:-/home/fatbaby/EMILY}"
TYLER_ROOT="${TYLER_ROOT:-/home/fatbaby/TYLER}"
FATBABY_ROOT="${FATBABY_ROOT:-/home/fatbaby/PRRJECT_FATBABY}"
CLAUDE_RUNS_DIR="$FATBABY_ROOT/claude-runs"
RSI_LOG="$EMILY_ROOT/var/logs/rsi-loop.log"
RSI_STATE="$EMILY_ROOT/var/rsi-loop-state.json"

# ── Guards ────────────────────────────────────────────────────────────────────
if ! command -v emily >/dev/null 2>&1; then
  echo "ERROR: emily CLI not found. Install: cd /home/fatbaby/emily.cli && ./scripts/build.sh"
  exit 1
fi

if [ -z "${ANTHROPIC_API_KEY:-}" ]; then
  echo "WARN: ANTHROPIC_API_KEY not set — Claude invocations will fail"
fi

mkdir -p "$EMILY_ROOT/var/logs"

# ── Helpers ───────────────────────────────────────────────────────────────────
log() {
  local ts
  ts=$(date '+%Y-%m-%d %H:%M:%S')
  echo "[$ts] $*" | tee -a "$RSI_LOG"
}

dry_prefix() {
  if [ "$DRY_RUN" = "1" ]; then
    echo "[DRY-RUN] "
  else
    echo ""
  fi
}

run_or_dry() {
  if [ "$DRY_RUN" = "1" ]; then
    echo "  [dry-run] would run: $*"
  else
    "$@"
  fi
}

write_state() {
  local iter="$1" task_id="$2" phase="$3"
  cat >"$RSI_STATE" <<EOF
{
  "iteration": $iter,
  "task_id": "$task_id",
  "phase": "$phase",
  "timestamp": "$(date -u '+%Y-%m-%dT%H:%M:%SZ')",
  "max_iters": $MAX_ITERS,
  "tyler_iters": $TYLER_ITERS
}
EOF
}

# ── Startup banner ────────────────────────────────────────────────────────────
echo ""
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║  EINHORN INDUSTRIAL ◈ EMILY OS — RSI MASTER LOOP               ║"
echo "╠══════════════════════════════════════════════════════════════════╣"
printf "║  Max iterations: %-10s  Tyler builds/cycle: %-10s    ║\n" "${MAX_ITERS:-∞}" "$TYLER_ITERS"
printf "║  Sleep between:  %-10s  Dry run: %-10s              ║\n" "${SLEEP_BETWEEN}s" "$DRY_RUN"
echo "╠══════════════════════════════════════════════════════════════════╣"
echo "║  TIC  → emily prime-task --preset rsi-token-report             ║"
echo "║  TOCK → observation-watcher dispatches Claude Code             ║"
echo "║  ENTR → TYLER/emily.sh $TYLER_ITERS (entropy injection)             ║"
echo "║  ANLZ → emily observe → Emily Prime reads → next TIC           ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""
log "RSI loop started (max_iters=$MAX_ITERS, tyler_iters=$TYLER_ITERS, sleep=${SLEEP_BETWEEN}s)"

# ── Check if system is running ────────────────────────────────────────────────
OBS_WATCHER_RUNNING=0
if pgrep -f "observation-watcher" >/dev/null 2>&1; then
  OBS_WATCHER_RUNNING=1
  log "observation-watcher: RUNNING ✓"
else
  log "WARN: observation-watcher not detected — tasks will queue but won't auto-dispatch"
  log "      Start with: emily start"
  echo ""
  echo "  TIP: Run 'emily start' first to bring up observation-watcher + emily-agent."
  echo "       The RSI loop will still fire tasks to the queue."
  echo ""
fi

ITER=0
# ── Main loop ─────────────────────────────────────────────────────────────────
while true; do
  ITER=$((ITER + 1))
  echo ""
  log "═══ ITERATION $ITER $([ "$MAX_ITERS" -gt 0 ] && echo "/ $MAX_ITERS" || echo "/ ∞") ═══"
  write_state "$ITER" "pending" "tic"

  # ── TIC: Ask Emily Prime ─────────────────────────────────────────────────
  log "[TIC] Dispatching RSI token-efficiency task to Emily Prime..."
  TASK_ID="dry-run-$ITER"
  if [ "$DRY_RUN" != "1" ]; then
    TASK_RESULT=$(emily prime-task --preset rsi-token-report --json 2>/dev/null || echo "{}")
    TASK_ID=$(echo "$TASK_RESULT" | grep -o '"task_id":"[^"]*"' | cut -d'"' -f4 || echo "task-$ITER")
    APPLE_ID=$(echo "$TASK_RESULT" | grep -o '"apple_id":[0-9]*' | cut -d: -f2 || echo "")
    log "  task written: $TASK_ID${APPLE_ID:+ (Apple #$APPLE_ID)}"
    log "  obs-watcher picks up within 10s → claude --dangerously-skip-permissions"
  else
    log "  [dry-run] emily prime-task --preset rsi-token-report --json"
    TASK_ID="dry-run-$ITER"
  fi
  write_state "$ITER" "$TASK_ID" "tock"

  # ── TOCK: Wait for Claude to complete ────────────────────────────────────
  if [ "$SKIP_WAIT" = "1" ]; then
    log "[TOCK] Skipping wait (SKIP_WAIT=1)"
  elif [ "$DRY_RUN" = "1" ]; then
    log "[TOCK] [dry-run] would wait for claude-runs/ new file (max 3 min)"
  else
    log "[TOCK] Waiting for Claude Code to complete task (max 3 min)..."
    WAIT_START=$(date +%s)
    MAX_WAIT=180
    INITIAL_RUNS=$(ls "$CLAUDE_RUNS_DIR" 2>/dev/null | wc -l || echo 0)
    PROGRESS_DOTS=0
    while true; do
      CURRENT_RUNS=$(ls "$CLAUDE_RUNS_DIR" 2>/dev/null | wc -l || echo 0)
      if [ "$CURRENT_RUNS" -gt "$INITIAL_RUNS" ]; then
        LATEST=$(ls -t "$CLAUDE_RUNS_DIR" 2>/dev/null | head -1 || echo "unknown")
        log "  ✓ Claude run detected → $LATEST (total: $CURRENT_RUNS)"
        break
      fi
      NOW=$(date +%s)
      ELAPSED=$((NOW - WAIT_START))
      if [ "$ELAPSED" -ge "$MAX_WAIT" ]; then
        log "  ⚠ Timeout after ${MAX_WAIT}s (obs-watcher running: $OBS_WATCHER_RUNNING) — proceeding"
        break
      fi
      if [ "$OBS_WATCHER_RUNNING" = "0" ]; then
        log "  ⚠ obs-watcher not running — task queued but not dispatched. Skipping wait."
        break
      fi
      PROGRESS_DOTS=$((PROGRESS_DOTS + 1))
      printf "  ... ${ELAPSED}s elapsed\r"
      sleep 10
    done
    echo ""
  fi
  write_state "$ITER" "$TASK_ID" "entropy"

  # ── ENTROPY: Tyler builds ─────────────────────────────────────────────────
  if [ "$SKIP_TYLER" = "1" ]; then
    log "[ENTROPY] Skipping Tyler (SKIP_TYLER=1)"
  elif [ ! -f "$TYLER_ROOT/emily.sh" ]; then
    log "[ENTROPY] WARN: $TYLER_ROOT/emily.sh not found — skipping entropy injection"
  elif [ "$DRY_RUN" = "1" ]; then
    log "[ENTROPY] [dry-run] would run: cd $TYLER_ROOT && ./emily.sh $TYLER_ITERS"
  else
    log "[ENTROPY] Running Tyler ($TYLER_ITERS iterations) for creative entropy..."
    TYLER_PENDING=$(grep -c '^\- \[ \]' "$TYLER_ROOT/BACKLOG.md" 2>/dev/null || echo 0)
    if [ "$TYLER_PENDING" -eq 0 ]; then
      log "  Tyler backlog empty — skipping (no tasks to build)"
    else
      log "  Tyler backlog: $TYLER_PENDING tasks pending"
      cd "$TYLER_ROOT"
      ./emily.sh "$TYLER_ITERS" 2>&1 | while IFS= read -r line; do
        log "  [TYLER] $line"
      done || log "  [TYLER] returned non-zero (continuing RSI loop)"
      cd "$EMILY_ROOT"
    fi
  fi
  write_state "$ITER" "$TASK_ID" "analyze"

  # ── ANALYZE: Emily Prime observes the completed cycle ─────────────────────
  log "[ANALYZE] Posting RSI iteration observation to Emily OS..."
  ITER_MSG="RSI loop iteration $ITER complete — FatBaby+Tyler tic-toc cycle done"
  ITER_FINDINGS="TIC: prime-task $TASK_ID dispatched. TOCK: Claude Code processed. ENTROPY: Tyler ran $TYLER_ITERS build(s). System advancing via recursive self-improvement."
  if [ "$DRY_RUN" = "1" ]; then
    log "  [dry-run] emily observe -s info '$ITER_MSG'"
  else
    emily observe -s info "$ITER_MSG" \
      --findings "$ITER_FINDINGS" \
      2>/dev/null && log "  ✓ Observation posted → obs-watcher will pick up" \
      || log "  WARN: observe failed (IDUNA offline?) — continuing"
  fi
  write_state "$ITER" "$TASK_ID" "complete"
  log "✓ Iteration $ITER complete"

  # ── Stop condition ────────────────────────────────────────────────────────
  if [ "$MAX_ITERS" -gt 0 ] && [ "$ITER" -ge "$MAX_ITERS" ]; then
    echo ""
    log "═══ RSI loop finished after $ITER iteration(s) ═══"
    echo ""
    echo "  emily status    — check repo + IDUNA state"
    echo "  emily apples list -n 10 — see last 10 Apples"
    echo ""
    break
  fi

  # ── Sleep ─────────────────────────────────────────────────────────────────
  log "Sleeping ${SLEEP_BETWEEN}s before iteration $((ITER + 1))..."
  sleep "$SLEEP_BETWEEN"
done
