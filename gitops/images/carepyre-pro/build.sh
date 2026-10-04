#!/usr/bin/env bash
# build.sh [TAG] — Cloud Build the idunapro image from IDUNA_PRO's git-tracked files (no var/, no keys).
set -euo pipefail
D="$(cd "$(dirname "$0")" && pwd)"; SRC="${IDUNA_PRO_DIR:-$HOME/IDUNA_PRO}"; PROJECT="${PROJECT:-project-d24a71e9-2daf-4b2d-917}"
TAG="${1:-$(git -C "$SRC" rev-parse --short HEAD)}"
CTX="$(mktemp -d)"; trap 'rm -rf "$CTX"' EXIT
git -C "$SRC" ls-files -z -- . ':(exclude)var' ':(exclude)docs' | (cd "$SRC" && xargs -0 -r tar cf - 2>/dev/null) | tar xf - -C "$CTX"
cp "$D/Dockerfile" "$CTX/Dockerfile"
gcloud builds submit "$CTX" --project "$PROJECT" --tag "us-central1-docker.pkg.dev/$PROJECT/emily/carepyre-pro:$TAG"
echo "carepyre-pro:$TAG"
