#!/usr/bin/env bash
# bump-tag.sh VAR_NAME NEW_TAG -- updates one service's tag in tags.env (the source of truth,
# see its own header) and regenerates every manifest from it. Commit + push the result and
# gitops-sync (in-cluster CronJob, every 3 min) deploys it -- no cluster credentials needed here.
#
# Called by hand for a manual release, or by .github/workflows/bump-deploy.yml on a
# repository_dispatch from a service repo's own CI once it has pushed a new image (K8S-CD-01).
#
# usage: gitops/bump-tag.sh SHANKPIT_TAG v5
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"
VAR="${1:?usage: bump-tag.sh VAR_NAME NEW_TAG}"
NEW="${2:?usage: bump-tag.sh VAR_NAME NEW_TAG}"

grep -q "^${VAR}=" "$HERE/tags.env" || {
  echo "bump-tag.sh: no '$VAR' in tags.env -- known vars:" >&2
  grep -oE '^[A-Z_]+=' "$HERE/tags.env" | tr -d '=' >&2
  exit 1
}

OLD="$(grep "^${VAR}=" "$HERE/tags.env" | head -1 | cut -d= -f2-)"
if [ "$OLD" = "$NEW" ]; then
  echo "bump-tag.sh: $VAR already $NEW, nothing to do"
  exit 0
fi

sed -i "s/^${VAR}=.*/${VAR}=${NEW}/" "$HERE/tags.env"
"$HERE/render.sh"
echo "bump-tag.sh: $VAR $OLD -> $NEW"
