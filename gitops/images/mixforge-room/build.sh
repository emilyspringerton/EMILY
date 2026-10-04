#!/usr/bin/env bash
# build.sh [TAG] — Cloud Build the mixforge-room image from the MIXFORGE checkout (server/ + web/, no node_modules, no cache).
set -euo pipefail
D="$(cd "$(dirname "$0")" && pwd)"; PROJECT="${PROJECT:-project-d24a71e9-2daf-4b2d-917}"; TAG="${1:-v1}"
SRC="${MIXFORGE_DIR:-/home/fatbaby/MIXFORGE}"
CTX="$(mktemp -d)"; trap 'rm -rf "$CTX"' EXIT
cp "$D/Dockerfile" "$CTX/"; mkdir "$CTX/server" "$CTX/web"
cp "$SRC/server/"{package.json,package-lock.json,room_server.mjs} "$CTX/server/"; cp -r "$SRC/web/." "$CTX/web/"; rm -rf "$CTX/web/room-cache"
gcloud builds submit "$CTX" --project "$PROJECT" --tag "us-central1-docker.pkg.dev/$PROJECT/emily/mixforge-room:$TAG"
