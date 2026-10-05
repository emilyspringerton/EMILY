#!/usr/bin/env bash
# Build+push the gitops-autodeploy image. parena-k8s-render/parena-pod-render build from
# COMMITTED source (PARENA's own tools/k8s_render_host.c + tools/pod_render_host.c,
# stdlib/k8s/{k8s,pod,gitops}.prn) -- no box, no pre-built binary required anywhere. Verified
# 2026-10-05: a from-source build's parena-k8s-render output is byte-identical to the committed
# manifests. Needs libsdl2-dev + libsdl2-ttf-dev (the base `parena` compiler binary links SDL2
# unconditionally for its own editor features, unrelated to this tool -- real, pre-existing PARENA
# build dependency, not introduced here).
# usage: images/gitops-autodeploy/build.sh [tag]   (default: git short SHA of this EMILY checkout)
#   PARENA_DIR=/path/to/existing/PARENA  -- reuse an already-cloned checkout instead of cloning fresh
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"
TAG="${1:-$(git -C "$HERE/../../.." rev-parse --short HEAD)}"
PROJECT="${GCP_PROJECT:-project-d24a71e9-2daf-4b2d-917}"
IMAGE="us-central1-docker.pkg.dev/$PROJECT/emily/gitops-autodeploy:$TAG"

CLEANUP=()
cleanup() { for d in "${CLEANUP[@]}"; do rm -rf "$d"; done; }
trap cleanup EXIT

if [ -n "${PARENA_DIR:-}" ]; then
  PD="$PARENA_DIR"
else
  PD="$(mktemp -d)"; CLEANUP+=("$PD")
  git clone --depth 1 https://github.com/emilyspringerton/PARENA.git "$PD"
fi
make -C "$PD" parena-k8s-render parena-pod-render

CTX="$(mktemp -d)"; CLEANUP+=("$CTX")
cp "$HERE/Dockerfile" "$HERE/entrypoint.sh" "$CTX/"
cp "$PD/parena-k8s-render" "$PD/parena-pod-render" "$CTX/"
# watch-tags.sh/bump-tag.sh/render.sh are deliberately NOT baked in here -- entrypoint.sh clones
# EMILY fresh on every run and execs them from that checkout, so the CronJob always runs
# whatever's currently committed, never a stale copy frozen at image-build time.

# Found live (2026-10-05, SHANKPIT): the CI SA (github-ci, see EMILY/gitops/CI_SETUP.md) is
# deliberately scoped to cloudbuild.builds.editor/artifactregistry.writer, not a project
# Viewer/Owner. Two consequences: (1) unpinned staging dir triggers a project-scoped
# storage.buckets.list call the CI SA's bucket-scoped bindings don't satisfy (403, misleadingly
# reported as a serviceusage error) -- --gcs-source-staging-dir skips that list; (2) `builds
# submit` can't stream the default (outside-the-project) logs bucket -- --async skips the wait,
# we poll `builds describe` (status only, never logs) ourselves.
BUILD_ID=$(gcloud builds submit "$CTX" --tag "$IMAGE" --project "$PROJECT" \
  --gcs-source-staging-dir="gs://${PROJECT}_cloudbuild/source" \
  --async --format="value(id)")

echo "submitted build $BUILD_ID, polling for completion..."
while true; do
  STATUS=$(gcloud builds describe "$BUILD_ID" --project "$PROJECT" --format="value(status)")
  case "$STATUS" in
    SUCCESS) echo "gitops-autodeploy:$TAG"; exit 0 ;;
    FAILURE|INTERNAL_ERROR|TIMEOUT|CANCELLED|EXPIRED) echo "build $BUILD_ID: $STATUS" >&2; exit 1 ;;
    *) sleep 5 ;;
  esac
done
