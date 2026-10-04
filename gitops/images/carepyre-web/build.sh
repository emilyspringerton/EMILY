#!/usr/bin/env bash
# build.sh [TAG] — Cloud Build the carepyre-web image. Site source: SITE_DIR (default: the live box copy /var/www/carepyre,
# which is what currently serves carepyre.org; it differs from the CarePyre repo and includes provisioning/).
set -euo pipefail
D="$(cd "$(dirname "$0")" && pwd)"; PROJECT="${PROJECT:-project-d24a71e9-2daf-4b2d-917}"; TAG="${1:-v1}"
CTX="$(mktemp -d)"; trap 'rm -rf "$CTX"' EXIT
cp "$D/Dockerfile" "$D/default.conf" "$CTX/"; cp -r "${SITE_DIR:-/var/www/carepyre}" "$CTX/site"
gcloud builds submit "$CTX" --project "$PROJECT" --tag "us-central1-docker.pkg.dev/$PROJECT/emily/carepyre-web:$TAG"
