#!/bin/bash
# compression-abtest.sh — A/B test haiku comprehension: GOLDEN.md vs GOLDEN_BILINGUAL_TEST.md
#
# Requires: ANTHROPIC_API_KEY
# Usage: ./scripts/compression-abtest.sh
#
# Tests three comprehension questions against both versions.
# Prints: token counts, scores, pass/fail recommendation.

set -euo pipefail

API_KEY="${ANTHROPIC_API_KEY:-}"
if [ -z "$API_KEY" ]; then
  echo "ERROR: ANTHROPIC_API_KEY not set" >&2
  exit 1
fi

EMILY_ROOT="$(dirname "$0")/.."
GOLDEN="${EMILY_ROOT}/GOLDEN.md"
BILINGUAL="${EMILY_ROOT}/docs/compression-experiment/GOLDEN_BILINGUAL_TEST.md"
MODEL="claude-haiku-4-5-20251001"

Q1="Which sections have open items? List section numbers and counts."
Q2="What is the RSI prime directive pipeline in order? List each step."
Q3="Name all the repos (code repositories) mentioned in this system."

# Expected answers for scoring
EXPECTED_Q1="S2:3 S3:2 S4:3 S5:2 S10:1 S13:2 S14:2 S15:1 S16:3 S17:5"
EXPECTED_Q2="raw thought → emily eo → INTAKE QUEUE → emily backlog promote → numbered section → HEIMDAL sprint → implementation"
EXPECTED_Q3="EMILY PRRJECT_FATBABY IDUNA TYLER SHANKPIT MJOLNIR APPLES emily.cli"

query_haiku() {
  local doc_path="$1"
  local question="$2"
  local doc_content
  doc_content=$(cat "$doc_path")

  curl -sf https://api.anthropic.com/v1/messages \
    -H "x-api-key: ${API_KEY}" \
    -H "anthropic-version: 2023-06-01" \
    -H "content-type: application/json" \
    -d "$(python3 -c "
import json, sys
doc = open('${doc_path}').read()
q = '''${question}'''
payload = {
  'model': '${MODEL}',
  'max_tokens': 200,
  'messages': [{
    'role': 'user',
    'content': 'Context:\n\n' + doc + '\n\nQuestion: ' + q + '\n\nAnswer concisely.'
  }]
}
print(json.dumps(payload))
")" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['content'][0]['text'])"
}

token_count() {
  # Rough estimate: 1 token ≈ 4 chars for English, ~1.5 chars per token for Chinese
  local f="$1"
  local chars
  chars=$(wc -c < "$f")
  echo $(( chars / 4 ))
}

echo "=== GOLDEN.md vs GOLDEN_BILINGUAL_TEST.md Compression A/B Test ==="
echo ""
echo "GOLDEN.md tokens (estimate): $(token_count "$GOLDEN")"
echo "BILINGUAL tokens (estimate): $(token_count "$BILINGUAL")"
echo "Token reduction: ~$(( ( $(token_count "$GOLDEN") - $(token_count "$BILINGUAL") ) * 100 / $(token_count "$GOLDEN") ))%"
echo ""

PASS=0
TOTAL=3

for q_num in 1 2 3; do
  eval "Q=\$Q${q_num}"
  eval "EXPECTED=\$EXPECTED_Q${q_num}"

  echo "--- Q${q_num}: $Q ---"

  echo "GOLDEN.md response:"
  GOLDEN_RESP=$(query_haiku "$GOLDEN" "$Q")
  echo "$GOLDEN_RESP"

  echo ""
  echo "BILINGUAL response:"
  BILINGUAL_RESP=$(query_haiku "$BILINGUAL" "$Q")
  echo "$BILINGUAL_RESP"

  echo ""
  echo "(Manual review needed — expected key terms: $EXPECTED)"
  echo ""
done

echo "=== RESULT ==="
echo "Review responses above. If all 3 questions answered equivalently:"
echo "  → PASS: replace GOLDEN.md with GOLDEN_BILINGUAL_TEST.md"
echo "  → Update emily backlog compress to output bilingual format"
echo "If any question degrades:"
echo "  → FAIL: keep GOLDEN.md, document failure mode in analysis.md"
