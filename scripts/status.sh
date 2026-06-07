#!/bin/bash
# status.sh — Emily Prime cross-repo status: git + backlog + last apple per agent
#
# Usage:
#   ./scripts/status.sh
#
# Requires:
#   IDUNA_BASE_URL  (default: http://localhost:8080)
#   IDUNA_AGENT_NAME + IDUNA_AGENT_SECRET (auto-sourced from IDUNA var/ if present)

set -euo pipefail

BASE_URL="${IDUNA_BASE_URL:-http://localhost:8080}"
AGENT_NAME="${IDUNA_AGENT_NAME:-EMILY-PRIME}"
AGENT_SECRET="${IDUNA_AGENT_SECRET:-}"

REPOS=(
  "/home/fatbaby/TYLER"
  "/home/fatbaby/EMILY"
  "/home/fatbaby/IDUNA"
  "/home/fatbaby/PRRJECT_FATBABY"
  "/home/fatbaby/SHANKPIT"
)

if [ -z "$AGENT_SECRET" ]; then
  SECRETS="${IDUNA_ROOT:-/home/fatbaby/IDUNA}/var/agent-secrets.env"
  if [ -f "$SECRETS" ]; then
    # shellcheck disable=SC1090
    source "$SECRETS"
    AGENT_SECRET="${IDUNA_SECRET_EMILY_PRIME:-}"
  fi
fi

echo ""
echo "◈ EMILY OS — SYSTEM STATUS | $(date '+%Y-%m-%d %H:%M')"
echo "══════════════════════════════════════════════════════"

# ── Git Status ────────────────────────────────────────────────────────────────

echo ""
echo "  GIT REPOS"
echo "  ─────────"

for repo in "${REPOS[@]}"; do
  name=$(basename "$repo")
  if [ ! -d "$repo/.git" ]; then
    printf "  %-20s  (no git)\n" "$name"
    continue
  fi

  cd "$repo"
  branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "?")
  last_commit=$(git log -1 --pretty="%h %s" 2>/dev/null | cut -c1-60 || echo "(no commits)")
  dirty=$(git status --porcelain 2>/dev/null | wc -l | tr -d ' ')
  pending_items=$(grep -c '^\- \[ \]' BACKLOG.md 2>/dev/null || echo 0)
  done_items=$(grep -c '^\- \[x\]' BACKLOG.md 2>/dev/null || echo 0)

  dirty_tag=""
  if [ "$dirty" -gt 0 ]; then
    dirty_tag=" [+${dirty} dirty]"
  fi

  printf "  %-20s  %s%s\n" "$name" "$branch" "$dirty_tag"
  printf "  %-20s  %s\n" "" "$last_commit"
  if [ -f "BACKLOG.md" ]; then
    printf "  %-20s  backlog: %s done / %s pending\n" "" "$done_items" "$pending_items"
  fi
  echo ""
done

# ── IDUNA ─────────────────────────────────────────────────────────────────────

IDUNA_OK=0
if [ -z "$AGENT_SECRET" ]; then
  echo "  IDUNA: (no credentials — skip)"
else
  TOKEN=$(curl -sf -X POST "${BASE_URL}/api/v1/auth/agent" \
    -H "Content-Type: application/json" \
    -d "{\"agent_name\":\"${AGENT_NAME}\",\"agent_secret\":\"${AGENT_SECRET}\"}" 2>/dev/null \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])" 2>/dev/null) || true

  if [ -n "$TOKEN" ]; then
    IDUNA_OK=1
    APPLES_RAW=$(curl -sf "${BASE_URL}/api/v1/apples?limit=50" \
      -H "Authorization: Bearer $TOKEN" 2>/dev/null || echo '{}')

    echo "  IDUNA: ${BASE_URL} ✓"
    echo "  ──────"
    echo "$APPLES_RAW" | python3 -c "
import sys, json
d = json.load(sys.stdin)
apples = d.get('apples', d) if isinstance(d, dict) else d
if not isinstance(apples, list):
    print('  (no apples)')
    sys.exit(0)

# Group by source_repo, find latest per repo
latest = {}
for a in apples:
    repo = a.get('source_repo', 'UNKNOWN')
    if repo not in latest:
        latest[repo] = a

print('  Last Apple per repo:')
for repo, a in sorted(latest.items()):
    ts = a.get('recorded_at', '?')[:16].replace('T', ' ')
    atype = a.get('apple_type', '?')
    title = a.get('title', '?')[:45]
    print(f'    [{repo:<12}]  {ts}  #{a[\"id\"]}  {title}')
print(f'  Total apples: {len(apples)}')
"
  else
    echo "  IDUNA: ${BASE_URL} (offline or auth failed)"
  fi
fi

echo ""
echo "──────────────────────────────────────────────────────"
echo ""
