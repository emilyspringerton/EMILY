#!/usr/bin/env bash
# build.sh [TAG] — Cloud Build the einhorn-survival image (JRE + /run-paper.sh). World/jars/plugins are on PVC einhorn-var, not baked.
set -euo pipefail
D="$(cd "$(dirname "$0")" && pwd)"; PROJECT="${PROJECT:-project-d24a71e9-2daf-4b2d-917}"; TAG="${1:-v1}"
gcloud builds submit "$D" --project "$PROJECT" --tag "us-central1-docker.pkg.dev/$PROJECT/emily/einhorn-survival:$TAG"
