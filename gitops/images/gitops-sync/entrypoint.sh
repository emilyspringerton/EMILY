#!/bin/sh
# gitops-sync: long-running `parena-gitops` (replaces the old sync.sh CronJob). Clones the public EMILY
# repo once, then every GITOPS_INTERVAL seconds fast-forwards it and kubectl-applies the cluster dir
# when HEAD moved. Failed revisions back off and HOLD after 6 attempts (a new commit starts clean).
# State lives in an emptyDir, so a restart re-applies once (idempotent, no prune).
set -eu
REPO="${GITOPS_REPO:-https://github.com/emilyspringerton/EMILY.git}"
DIR="${GITOPS_DIR:-gitops/clusters/prrject-fatbaby}"
W=/tmp/gitops
mkdir -p "$W"
[ -d "$W/repo/.git" ] || git clone --depth 1 --quiet "$REPO" "$W/repo"
# -x 00-: the Namespace is cluster-scoped; the namespaced Role cannot apply it.
exec parena-gitops -r "$W/repo" -p "$DIR" -s "$W/state" -i "${GITOPS_INTERVAL:-180}" -x 00-
