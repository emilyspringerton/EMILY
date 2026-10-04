#!/bin/sh
# gitops-sync: pull the public EMILY repo and kubectl-apply gitops/clusters/<cluster>/ into this namespace.
# Idempotent (apply of an unchanged manifest is a no-op), no prune (a deleted file never deletes a live object),
# 00-namespace is skipped (cluster-scoped; the Role is namespaced). Exits 1 if any manifest failed so the Job history shows it.
set -u
REPO="${GITOPS_REPO:-https://github.com/emilyspringerton/EMILY.git}"
DIR="${GITOPS_DIR:-gitops/clusters/prrject-fatbaby}"
W=$(mktemp -d); trap 'rm -rf "$W"' EXIT
git clone --depth 1 --quiet "$REPO" "$W/r" || { echo "gitops-sync: clone failed"; exit 1; }
cd "$W/r/$DIR" || exit 1
echo "gitops-sync: $(git -C "$W/r" rev-parse --short HEAD) $(git -C "$W/r" log -1 --format=%s)"
rc=0
for f in *.yaml; do
  case "$f" in 00-*) continue;; esac
  out=$(kubectl apply -f "$f" 2>&1) || { echo "FAIL $f: $out"; rc=1; continue; }
  echo "$out" | grep -v ' unchanged$' || true
done
echo "gitops-sync: done rc=$rc"
exit $rc
