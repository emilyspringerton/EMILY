#!/usr/bin/env bash
# build.sh [TAG] — Cloud Build the mixforge-web (nginx) image; site = MIXFORGE/web from the checkout.
set -euo pipefail
D="$(cd "$(dirname "$0")" && pwd)"; PROJECT="${PROJECT:-project-d24a71e9-2daf-4b2d-917}"; TAG="${1:-v1}"
CTX="$(mktemp -d)"; trap 'rm -rf "$CTX"' EXIT
cp "$D/Dockerfile" "$D/default.conf" "$CTX/"; cp -r "${MIXFORGE_DIR:-/home/fatbaby/MIXFORGE}/web" "$CTX/site"; rm -rf "$CTX/site/room-cache"
gcloud builds submit "$CTX" --project "$PROJECT" --tag "us-central1-docker.pkg.dev/$PROJECT/emily/mixforge-web:$TAG"
