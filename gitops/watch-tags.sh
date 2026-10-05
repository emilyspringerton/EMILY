#!/usr/bin/env bash
# watch-tags.sh -- K8S-CD-01 Phase 2 ("manifests regenerate off the new tags, automated").
#
# Polls Artifact Registry for each service's newest pushed image tag; whenever one has moved
# since the last check, calls bump-tag.sh (updates tags.env, re-renders every manifest) and
# commits+pushes the result. gitops-sync (in-cluster CronJob, every 3 min, K8S-CD-01 Phase 1)
# then applies it -- the two halves together are the full auto-deploy loop.
#
# Deliberately does NOT build or push images itself. Each service's own build/push process is
# still whatever it already is today (manual for most; build-image.sh for iduna/fatbaby, see
# their own repos). This script only closes the "a new tag exists -> deploy it" half; wiring
# per-service CI to build+push automatically on commit is real, separate, not done here (see
# BACKLOG K8S-CD-01's own Phase 2 note and this script's own header above).
#
# Runs IN the cluster, as the gitops-autodeploy CronJob (images/gitops-autodeploy/), NOT on the
# box: the box is ephemeral, so nothing about this loop's ability to run can depend on it staying
# up. The image bakes in the PARENA render binaries (copied from a box that has them built, at
# image-build time only); gcloud auth comes from Workload Identity and git push creds from a
# Secret -- see gitops/secrets/README.md's own "gitops-autodeploy" section for both one-time
# setup steps.
#
# Caveat: the exact `gcloud artifacts docker tags list` invocation below is written from gcloud's
# documented flags, not live-tested (no gcloud in the authoring sandbox) -- verify it against the
# real registry before trusting the timer unattended; `bash gitops/watch-tags.sh` run by hand is
# the way to check.
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"
REG="us-central1-docker.pkg.dev/project-d24a71e9-2daf-4b2d-917/emily"

# VAR image [image2 ...] -- one VAR may drive more than one image; they're always released
# together today (true for every row below -- verified against gitops/specs/*.pod).
#
# NOT listed here because they're still hardcoded directly in gitops/specs/*.pod, not driven by
# any tags.env var yet (bump these by hand in the spec file until they get their own var):
# dwlogs, sc-tunnel, okemily-web (all in iduna.pod / deadweight.pod), gfd-wsudprelay
# (gfd-wsrelay.pod), carepyre-web (carepyre.pod).
MAP="
COLLECTIONS_TAG collections-server
FATBABY_TAG fatbaby
RG_STABLE_TAG redgarden-stable
RG_RND_TAG redgarden-rnd
GFD_TAG gfd
SHANKPIT_TAG shankpit
IDUNA_TAG iduna
DEADWEIGHT_TAG deadweight
WOTAN_TAG wotan
TCP_EDGE_TAG tcp-edge
GITOPS_SYNC_TAG gitops-sync
CAREPYRE_TAG carepyre-pro
EINHORN_TAG einhorn-survival
MIXFORGE_TAG mixforge-room mixforge-web
"

newest_tag() {
  gcloud artifacts docker tags list "$REG/$1" \
    --sort-by=~CREATE_TIME --limit=1 --format='value(TAG)' 2>/dev/null
}

changed=0
while read -r var images; do
  [ -z "$var" ] && continue
  current="$(grep "^${var}=" "$HERE/tags.env" | head -1 | cut -d= -f2-)"
  first_image="$(awk '{print $1}' <<< "$images")"
  latest="$(newest_tag "$first_image")"
  if [ -n "$latest" ] && [ "$latest" != "$current" ]; then
    echo "watch-tags.sh: $var $current -> $latest ($images)"
    "$HERE/bump-tag.sh" "$var" "$latest"
    changed=1
  fi
done <<< "$MAP"

if [ "$changed" = 1 ]; then
  cd "$HERE/.."
  git add gitops/tags.env gitops/clusters
  git commit -m "$(cat <<'MSG'
feat(gitops): auto-bump tag(s) from newest Artifact Registry push (K8S-CD-01)

automated-by: gitops/watch-tags.sh
MSG
)"
  git push
else
  echo "watch-tags.sh: no new tags"
fi
