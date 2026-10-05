#!/usr/bin/env bash
# Build+push the gitops-autodeploy image. Run from a box that has the PARENA render binaries
# built (parena-k8s-render/parena-pod-render) -- they're copied INTO this image's build context
# so the image itself is self-contained afterward (the box's own lifecycle doesn't matter once
# this image lives in Artifact Registry and the CronJob is pulling it in-cluster).
# usage: images/gitops-autodeploy/build.sh [tag]   (default: git short SHA of this EMILY checkout)
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"
TAG="${1:-$(git -C "$HERE/../../.." rev-parse --short HEAD)}"
PROJECT="${GCP_PROJECT:-project-d24a71e9-2daf-4b2d-917}"
IMAGE="us-central1-docker.pkg.dev/$PROJECT/emily/gitops-autodeploy:$TAG"

PARENA_K8S_RENDER="${PARENA_K8S_RENDER:-$HOME/PARENA/parena-k8s-render}"
PARENA_POD_RENDER="${PARENA_POD_RENDER:-$HOME/PARENA/parena-pod-render}"
[ -x "$PARENA_K8S_RENDER" ] || { echo "build.sh: $PARENA_K8S_RENDER not found/executable" >&2; exit 1; }
[ -x "$PARENA_POD_RENDER" ] || { echo "build.sh: $PARENA_POD_RENDER not found/executable" >&2; exit 1; }

CTX="$(mktemp -d)"; trap 'rm -rf "$CTX"' EXIT
cp "$HERE/Dockerfile" "$HERE/entrypoint.sh" "$CTX/"
cp "$PARENA_K8S_RENDER" "$PARENA_POD_RENDER" "$CTX/"
# watch-tags.sh/bump-tag.sh/render.sh are deliberately NOT baked in here -- entrypoint.sh clones
# EMILY fresh on every run and execs them from that checkout, so the CronJob always runs
# whatever's currently committed, never a stale copy frozen at image-build time.
gcloud builds submit "$CTX" --tag "$IMAGE" --project "$PROJECT"
echo "gitops-autodeploy:$TAG"
