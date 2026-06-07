#!/bin/bash
# sync-fatbaby.sh — Post FatBaby-Emily observations to IDUNA as Apples
#
# Usage:
#   ./scripts/sync-fatbaby.sh              # post new observations only (up to --limit)
#   ./scripts/sync-fatbaby.sh --all        # backfill all (skips already-posted via state file)
#   ./scripts/sync-fatbaby.sh --dry-run    # show what would be posted, don't POST
#   ./scripts/sync-fatbaby.sh --limit=N    # max N files to process (default 10)
#
# State: EMILY/var/fatbaby-synced.txt — one filename per line, already posted.
#
# Requires:
#   IDUNA_BASE_URL  (default: http://localhost:8080)
#   IDUNA_AGENT_NAME + IDUNA_AGENT_SECRET (auto-sourced from IDUNA var/ if present)

set -euo pipefail

OBS_DIR="${FATBABY_OBS_DIR:-/home/fatbaby/PRRJECT_FATBABY/var/emily-observations}"
STATE_FILE="${EMILY_STATE:-/home/fatbaby/EMILY/var/fatbaby-synced.txt}"
BASE_URL="${IDUNA_BASE_URL:-http://localhost:8080}"
AGENT_NAME="${IDUNA_AGENT_NAME:-EMILY-PRIME}"
AGENT_SECRET="${IDUNA_AGENT_SECRET:-}"
ALL=0
DRY_RUN=0
LIMIT=10

for arg in "$@"; do
  case "$arg" in
    --all) ALL=1 ;;
    --dry-run) DRY_RUN=1 ;;
    --limit=*) LIMIT="${arg#--limit=}" ;;
  esac
done

# ── Auth ──────────────────────────────────────────────────────────────────────

if [ -z "$AGENT_SECRET" ]; then
  SECRETS="${IDUNA_ROOT:-/home/fatbaby/IDUNA}/var/agent-secrets.env"
  if [ -f "$SECRETS" ]; then
    # shellcheck disable=SC1090
    source "$SECRETS"
    AGENT_SECRET="${IDUNA_SECRET_EMILY_PRIME:-}"
  fi
fi

if [ -z "$AGENT_SECRET" ]; then
  echo "ERROR: IDUNA_AGENT_SECRET not set." >&2; exit 1
fi

TOKEN=""
if [ "$DRY_RUN" -eq 0 ]; then
  TOKEN=$(curl -sf -X POST "${BASE_URL}/api/v1/auth/agent" \
    -H "Content-Type: application/json" \
    -d "{\"agent_name\":\"${AGENT_NAME}\",\"agent_secret\":\"${AGENT_SECRET}\"}" 2>/dev/null \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])") || {
    echo "ERROR: IDUNA not reachable at ${BASE_URL}" >&2; exit 1
  }
fi

# ── State ─────────────────────────────────────────────────────────────────────

mkdir -p "$(dirname "$STATE_FILE")"
touch "$STATE_FILE"

already_posted() {
  grep -qxF "$1" "$STATE_FILE" 2>/dev/null
}

mark_posted() {
  echo "$1" >> "$STATE_FILE"
}

# ── Post ──────────────────────────────────────────────────────────────────────

POSTED=0
SKIPPED=0
TMPFILE=$(mktemp /tmp/fatbaby-apple-XXXXXX.json)
trap "rm -f $TMPFILE" EXIT

echo ""
echo "◈ EMILY OS — FATBABY OBSERVATION SYNC | $(date '+%Y-%m-%d %H:%M')"

for obs_file in $(ls -1t "$OBS_DIR"/*.json 2>/dev/null | head -"$LIMIT"); do
  fname=$(basename "$obs_file")
  [ "$fname" = "latest.json" ] && continue

  if already_posted "$fname"; then
    SKIPPED=$((SKIPPED + 1))
    continue
  fi

  python3 - "$obs_file" > "$TMPFILE" 2>/dev/null <<'PYEOF'
import json, sys

obs_path = sys.argv[1]
fname = obs_path.rsplit('/', 1)[-1].replace('.json', '')

with open(obs_path) as f:
    d = json.load(f)

summary  = d.get('summary', '(no summary)')
severity = d.get('severity', 'info')
ts       = d.get('timestamp', '?')
findings = d.get('findings', '')
fix      = d.get('suggested_fix', '')

title = f"FatBaby obs [{severity}]: {summary}"[:100]
body_parts = [f"severity: {severity}", f"timestamp: {ts}"]
if findings:
    body_parts.append(f"\nfindings:\n{findings[:800]}")
if fix:
    body_parts.append(f"\nsuggested_fix:\n{fix[:400]}")
body = "\n".join(body_parts)

payload = {
    "apple_type": "signal_observation",
    "title": title,
    "body": body,
    "source_repo": "PRRJECT_FATBABY",
    "run_id": fname,
}
print(json.dumps(payload))
PYEOF

  if [ ! -s "$TMPFILE" ]; then
    echo "  SKIP: $fname (parse error)"; continue
  fi

  title=$(python3 -c "import sys,json; print(json.load(sys.stdin)['title'])" < "$TMPFILE" 2>/dev/null || echo "$fname")

  if [ "$DRY_RUN" -eq 1 ]; then
    echo "  [DRY-RUN] $fname → $title"
    POSTED=$((POSTED + 1))
    continue
  fi

  RESULT=$(curl -sf -X POST "${BASE_URL}/api/v1/apples" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    --data-binary "@$TMPFILE" 2>/dev/null) || { echo "  FAIL: $fname"; continue; }

  apple_id=$(echo "$RESULT" | python3 -c "import sys,json; print(json.load(sys.stdin).get('id','?'))" 2>/dev/null || echo "?")
  echo "  Apple #${apple_id} ← $fname"
  mark_posted "$fname"
  POSTED=$((POSTED + 1))
done

echo ""
echo "  Done. Posted: ${POSTED} | Skipped (already synced): ${SKIPPED}"
echo ""
