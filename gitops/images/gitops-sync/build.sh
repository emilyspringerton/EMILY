#!/usr/bin/env bash
# build.sh [TAG] -- build the gitops-sync image: parena-gitops built from PARENA's own committed source
# (same no-box pattern as images/gitops-autodeploy/build.sh; needs libsdl2-dev + libsdl2-ttf-dev for the
# PARENA compiler), relinked statically so it runs on alpine, then Cloud Build (async + status poll, see
# gitops/CI_SETUP.md for why). usage: build.sh [tag]  (default: git short SHA of this EMILY checkout)
#   PARENA_DIR=/path/to/PARENA reuses an existing checkout instead of cloning.
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"
TAG="${1:-$(git -C "$HERE/../../.." rev-parse --short HEAD)}"
PROJECT="${GCP_PROJECT:-project-d24a71e9-2daf-4b2d-917}"
IMAGE="us-central1-docker.pkg.dev/$PROJECT/emily/gitops-sync:$TAG"
CLEANUP=(); cleanup() { for d in "${CLEANUP[@]}"; do rm -rf "$d"; done; }; trap cleanup EXIT
if [ -n "${PARENA_DIR:-}" ]; then PD="$PARENA_DIR"
else PD="$(mktemp -d)"; CLEANUP+=("$PD"); git clone --depth 1 https://github.com/emilyspringerton/PARENA.git "$PD"; fi
make -C "$PD" parena-gitops
# same command as the Makefile's parena-gitops rule, linked -static (the binary only needs libc/libm)
( cd "$PD" && ${CC:-gcc} -static -std=c99 -D_POSIX_C_SOURCE=200809L -I runtime -I tools tools/gitops_host.c \
    runtime/parena_runtime.c -o parena-gitops-static -lm )
CTX="$(mktemp -d)"; CLEANUP+=("$CTX")
cp "$HERE/Dockerfile" "$HERE/entrypoint.sh" "$CTX/"
cp "$PD/parena-gitops-static" "$CTX/parena-gitops"
BUILD_ID=$(gcloud builds submit "$CTX" --tag "$IMAGE" --project "$PROJECT" \
  --gcs-source-staging-dir="gs://${PROJECT}_cloudbuild/source" --async --format="value(id)")
echo "submitted build $BUILD_ID, polling for completion..."
while true; do
  STATUS=$(gcloud builds describe "$BUILD_ID" --project "$PROJECT" --format="value(status)")
  case "$STATUS" in
    SUCCESS) echo "gitops-sync:$TAG"; exit 0 ;;
    FAILURE|INTERNAL_ERROR|TIMEOUT|CANCELLED|EXPIRED) echo "build $BUILD_ID: $STATUS" >&2; exit 1 ;;
    *) sleep 5 ;;
  esac
done
