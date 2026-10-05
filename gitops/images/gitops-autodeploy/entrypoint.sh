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
REPO="${GITOPS_REPO:-git@github.com:emilyspringerton/EMILY.git}"
SSH_KEY="${GITOPS_SSH_KEY:-/etc/gitops-autodeploy/ssh/id_ed25519}"
[ -f "$SSH_KEY" ] || { echo "missing SSH key at $SSH_KEY (from the gitops-autodeploy-git Secret)" >&2; exit 1; }
W="$(mktemp -d)"
trap 'rm -rf "$W"' EXIT

# Switched from a GitHub PAT to this box's own account-level SSH deploy key (founder, 2026-10-05:
# "use the box ssh key why are you using a PAT?") -- same key every `git push` in this repo's own
# CI/operator flow already uses, so no new GitHub-side credential to create or scope. Pin the host
# key rather than StrictHostKeyChecking=no: avoids trusting whatever github.com happens to answer
# with on a given run.
export GIT_SSH_COMMAND="ssh -i $SSH_KEY -o IdentitiesOnly=yes -o UserKnownHostsFile=/etc/gitops-autodeploy/ssh/known_hosts -o StrictHostKeyChecking=yes"
git clone --depth 1 --quiet "$REPO" "$W/EMILY"
cd "$W/EMILY"
git config user.name "gitops-autodeploy"
git config user.email "gitops-autodeploy@users.noreply.github.com"

# The render binaries are baked into this image (see build.sh -- built from PARENA's own
# committed source at image-build time, no box involved); the checkout itself never carries them.
export PARENA_K8S_RENDER=/usr/local/bin/parena-k8s-render
export PARENA_POD_RENDER=/usr/local/bin/parena-pod-render

exec gitops/watch-tags.sh
