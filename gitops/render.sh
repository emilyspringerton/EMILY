#!/usr/bin/env bash
# Regenerate gitops/clusters/prrject-fatbaby/*.yaml with PARENA's renderer (make -C PARENA parena-k8s-render).
# Manifests are generated, never hand-edited: change the flags here, rerun, commit.
set -euo pipefail
R="${PARENA_K8S_RENDER:-/home/fatbaby/PARENA/parena-k8s-render}"
OUT="$(cd "$(dirname "$0")" && pwd)/clusters/prrject-fatbaby"
REG="us-central1-docker.pkg.dev/project-d24a71e9-2daf-4b2d-917/emily"
TAG="${COLLECTIONS_TAG:-0.1.0}"
printf 'apiVersion: v1\nkind: Namespace\nmetadata:\n  name: emily\n' > "$OUT/00-namespace.yaml"
# The ONE app with --host owns the cluster's single Ingress (one external LB; see networking cost rules).
"$R" --name collections-server --namespace emily --image "$REG/collections-server:$TAG" \
     --port 8087 --cpu-milli 100 --memory-mi 128 --host golden.okemily.com > "$OUT/10-collections-server.yaml"
