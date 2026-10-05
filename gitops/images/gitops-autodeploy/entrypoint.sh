#!/usr/bin/env bash
# gitops-autodeploy: runs IN the cluster (K8S-CD-01 Phase 2), not on the box -- the box is
# ephemeral, so nothing about this deploy loop's ability to run can depend on it staying up.
# Clones a fresh, writable checkout of EMILY, then hands off to gitops/watch-tags.sh, which
# polls Artifact Registry for each service's newest pushed tag and, if anything moved,
# bumps gitops/tags.env + re-renders + commits + pushes. gitops-sync (the OTHER CronJob already
# live in this namespace) picks the push up and applies it within 3 min -- same pull-based,
# no-inbound-creds design as that job, just pointed the other direction (detect + regenerate
# instead of apply).
set -euo pipefail
REPO="${GITOPS_REPO:-https://github.com/emilyspringerton/EMILY.git}"
: "${GITHUB_TOKEN:?GITHUB_TOKEN env var required (from the gitops-autodeploy-git Secret)}"
W="$(mktemp -d)"
trap 'rm -rf "$W"' EXIT

AUTH_REPO="$(echo "$REPO" | sed "s#https://#https://x-access-token:${GITHUB_TOKEN}@#")"
git clone --depth 1 --quiet "$AUTH_REPO" "$W/EMILY"
cd "$W/EMILY"
git config user.name "gitops-autodeploy"
git config user.email "gitops-autodeploy@users.noreply.github.com"
# Keep the push remote authenticated but never print it (it would leak the token into job logs).
git remote set-url origin "$AUTH_REPO" >/dev/null

# The render binaries are baked into this image (see build.sh -- built from PARENA's own
# committed source at image-build time, no box involved); the checkout itself never carries them.
export PARENA_K8S_RENDER=/usr/local/bin/parena-k8s-render
export PARENA_POD_RENDER=/usr/local/bin/parena-pod-render

exec gitops/watch-tags.sh
