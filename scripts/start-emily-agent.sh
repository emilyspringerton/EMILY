#!/bin/bash
# start-emily-agent.sh — Launch the EMILY RSI engine (emily-agent/cron.go) with IDUNA wired.
#
# This is the correct way to start the emily-agent so that each RSI iteration
# posts an Apple to IDUNA. Without these env vars, the IdunaClient is nil and
# Apple submission is silently skipped.
#
# Usage:
#   ./scripts/start-emily-agent.sh           # daemon mode (runs every 5 minutes)
#   ./scripts/start-emily-agent.sh --once    # run exactly one cycle
#
# Prerequisites:
#   ANTHROPIC_API_KEY must be set in the environment.
#   IDUNA must be running (go run . in IDUNA repo).

set -euo pipefail

EMILY_DIR="/home/fatbaby/EMILY"
IDUNA_SECRETS="${IDUNA_ROOT:-/home/fatbaby/IDUNA}/var/agent-secrets.env"

# ── Load IDUNA credentials ────────────────────────────────────────────────────

if [ -f "$IDUNA_SECRETS" ]; then
  # shellcheck disable=SC1090
  source "$IDUNA_SECRETS"
else
  echo "WARNING: IDUNA secrets not found at $IDUNA_SECRETS"
  echo "  IDUNA Apple submission will be disabled for this run."
fi

export IDUNA_BASE_URL="${IDUNA_BASE_URL:-http://localhost:8080}"
export IDUNA_AGENT_NAME="${IDUNA_AGENT_NAME:-EMILY-PRIME}"
export IDUNA_AGENT_SECRET="${IDUNA_AGENT_SECRET:-${IDUNA_SECRET_EMILY_PRIME:-}}"

# ── Verify IDUNA is reachable ─────────────────────────────────────────────────

if [ -n "$IDUNA_AGENT_SECRET" ]; then
  if curl -sf "${IDUNA_BASE_URL}/api/v1/agents" \
       -H "Content-Type: application/json" \
       --max-time 3 >/dev/null 2>&1; then
    echo "IDUNA: reachable at ${IDUNA_BASE_URL} (Apple submission enabled)"
  else
    echo "IDUNA: not reachable at ${IDUNA_BASE_URL} (Apple submission will fail — non-fatal)"
  fi
else
  echo "IDUNA: no agent secret — Apple submission disabled"
fi

# ── Launch emily-agent ────────────────────────────────────────────────────────

echo ""
echo "Starting EMILY RSI engine..."
echo "  IDUNA_BASE_URL:   ${IDUNA_BASE_URL}"
echo "  IDUNA_AGENT_NAME: ${IDUNA_AGENT_NAME}"
echo "  Mode: ${1:---daemon}"
echo ""

cd "$EMILY_DIR/emily-agent"

if [ "${1:-}" = "--once" ]; then
  exec go run main.go rsi.go cron.go iduna.go -- --cron
else
  exec go run main.go rsi.go cron.go iduna.go -- --daemon
fi
