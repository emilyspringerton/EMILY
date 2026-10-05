#!/usr/bin/env bash
# gitops-autodeploy needs write access to the EMILY repo (it commits+pushes tag bumps) and read
# access to Artifact Registry (it lists image tags via Workload Identity, not this secret).
# Same "created out of band, never in git" rule as make-fatbaby-secret.sh.
#
# usage (from a machine with a kubeconfig for the cluster):
#   ./make-gitops-autodeploy-secret.sh <github-PAT-with-repo-scope-on-EMILY> | kubectl apply -f -
set -euo pipefail
TOKEN="${1:?usage: make-gitops-autodeploy-secret.sh <github-PAT>}"
kubectl create secret generic gitops-autodeploy-git \
  --namespace emily \
  --from-literal=token="$TOKEN" \
  --dry-run=client -o yaml
