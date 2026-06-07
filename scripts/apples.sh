#!/bin/bash
# apples.sh — Emily Prime's system view: query IDUNA Apples log
#
# Usage:
#   ./scripts/apples.sh              # last 20 apples, all agents
#   ./scripts/apples.sh tyler        # filter by agent name (partial match on source_repo)
#   ./scripts/apples.sh --full       # include body text
#
# Requires:
#   IDUNA_BASE_URL  (default: http://localhost:8080)
#   IDUNA_AGENT_NAME + IDUNA_AGENT_SECRET (for auth)
#   OR: source var/agent-secrets.env first
#
# If IDUNA is not reachable, exits 1 with a clear message.

set -euo pipefail

BASE_URL="${IDUNA_BASE_URL:-http://localhost:8080}"
AGENT_NAME="${IDUNA_AGENT_NAME:-EMILY-PRIME}"
AGENT_SECRET="${IDUNA_AGENT_SECRET:-}"
LIMIT="${APPLES_LIMIT:-20}"
FULL_BODY=0
FILTER=""

for arg in "$@"; do
  case "$arg" in
    --full) FULL_BODY=1 ;;
    --limit=*) LIMIT="${arg#--limit=}" ;;
    --*) echo "unknown flag: $arg" >&2; exit 1 ;;
    *) FILTER="$arg" ;;
  esac
done

# ── Auth ──────────────────────────────────────────────────────────────────────

if [ -z "$AGENT_SECRET" ]; then
  # Try sourcing the secrets file from IDUNA repo root
  SECRETS="${IDUNA_ROOT:-/home/fatbaby/IDUNA}/var/agent-secrets.env"
  if [ -f "$SECRETS" ]; then
    # shellcheck disable=SC1090
    source "$SECRETS"
    AGENT_SECRET="${IDUNA_SECRET_EMILY_PRIME:-}"
  fi
fi

if [ -z "$AGENT_SECRET" ]; then
  echo "ERROR: IDUNA_AGENT_SECRET not set and var/agent-secrets.env not found." >&2
  echo "  Run: source /home/fatbaby/IDUNA/var/agent-secrets.env" >&2
  exit 1
fi

AUTH_RESPONSE=$(curl -sf -X POST "${BASE_URL}/api/v1/auth/agent" \
  -H "Content-Type: application/json" \
  -d "{\"agent_name\":\"${AGENT_NAME}\",\"agent_secret\":\"${AGENT_SECRET}\"}" 2>/dev/null) || {
  echo "ERROR: IDUNA not reachable at ${BASE_URL}" >&2
  echo "  Start IDUNA: cd /home/fatbaby/IDUNA && go run ." >&2
  exit 1
}

TOKEN=$(echo "$AUTH_RESPONSE" | python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])")

# ── Query ─────────────────────────────────────────────────────────────────────

QUERY="${BASE_URL}/api/v1/apples?limit=${LIMIT}"
if [ -n "$FILTER" ]; then
  QUERY="${QUERY}&source_repo=${FILTER}"
fi

APPLES=$(curl -sf "${QUERY}" \
  -H "Authorization: Bearer ${TOKEN}" 2>/dev/null) || {
  echo "ERROR: failed to query apples" >&2
  exit 1
}

COUNT=$(echo "$APPLES" | python3 -c "
import sys, json
d = json.load(sys.stdin)
lst = d.get('apples', d) if isinstance(d, dict) else d
print(len(lst) if isinstance(lst, list) else 0)
" 2>/dev/null || echo "0")

# ── Display ───────────────────────────────────────────────────────────────────

echo ""
echo "◈ EMILY OS — APPLES LOG | ${BASE_URL} | $(date '+%Y-%m-%d %H:%M')"
echo "  Showing ${COUNT} apple(s)$([ -n "$FILTER" ] && echo " (filter: ${FILTER})" || true)"
echo ""

if [ "$COUNT" -eq 0 ]; then
  echo "  (no apples)"
  echo ""
  exit 0
fi

echo "$APPLES" | python3 -c "
import sys, json

d = json.load(sys.stdin)
apples = d.get('apples', d) if isinstance(d, dict) else d
if not isinstance(apples, list):
    print('  (unexpected response)')
    sys.exit(0)

full_body = $FULL_BODY

for a in apples:
    apple_id = a.get('id', '?')
    agent    = a.get('agent_id', '?')[:8]
    repo     = a.get('source_repo', '?')
    run_id   = a.get('run_id', '?')
    atype    = a.get('apple_type', '?')
    title    = a.get('title', '?')
    ts       = a.get('recorded_at', '?')[:16].replace('T', ' ')

    print(f'  #{apple_id:>4}  {ts}  [{repo:<12}]  {atype:<25}  {title[:55]}')
    if full_body:
        body = a.get('body', '')
        if body:
            for line in body.split(chr(10))[:5]:
                print(f'         {line}')
        print()

print()
"
