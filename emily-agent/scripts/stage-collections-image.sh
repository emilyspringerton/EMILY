#!/usr/bin/env bash
# Assemble the Docker build context for Dockerfile.collections.
# usage: scripts/stage-collections-image.sh <out-dir> [golden-docs-checkout]
set -euo pipefail
OUT="${1:?out dir}"; GOLD="${2:-/home/fatbaby/GOLDEN_DOCS}"
HERE="$(cd "$(dirname "$0")/.." && pwd)"
rm -rf "$OUT"; mkdir -p "$OUT/emily-agent" "$OUT/golden"
git -C "$HERE" ls-files -z | grep -z -v '^emily-memory/\|^emily-state/\|^fartco-memory/\|^conversations/' | (cd "$HERE" && xargs -0 -I{} cp --parents {} "$OUT/emily-agent/")
cp "$HERE/Dockerfile.collections" "$OUT/Dockerfile"
cp -r "$GOLD/context" "$GOLD/docs" "$OUT/golden/"
echo "staged $OUT ($(du -sh "$OUT" | cut -f1)); build: docker build -t collections-server $OUT"
