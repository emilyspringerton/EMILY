#!/usr/bin/env bash
# gitops-autodeploy needs write access to the EMILY repo (it commits+pushes tag bumps) and read
# access to Artifact Registry (it lists image tags via Workload Identity, not this secret).
# Same "created out of band, never in git" rule as make-fatbaby-secret.sh.
#
# Switched from a GitHub PAT to an SSH deploy key (founder, 2026-10-05: "use the box ssh key why
# are you using a PAT?") -- this is the box's own account-level key, the same one every `git
# push` in this repo's own operator/CI flow already uses, so no new GitHub-side credential to
# create or scope.
#
# usage (from a machine with a kubeconfig for the cluster):
#   ./make-gitops-autodeploy-secret.sh ~/.ssh/id_ed25519 | kubectl apply -f -
set -euo pipefail
KEY_FILE="${1:?usage: make-gitops-autodeploy-secret.sh <path-to-ssh-private-key>}"
kubectl create secret generic gitops-autodeploy-git \
  --namespace emily \
  --from-file=id_ed25519="$KEY_FILE" \
  --dry-run=client -o yaml
