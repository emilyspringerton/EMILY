#!/usr/bin/env bash
# Regenerate gitops/clusters/prrject-fatbaby/*.yaml with PARENA's renderers (make -C PARENA parena-k8s-render parena-pod-render).
# Manifests are generated, never hand-edited: change the flags here, rerun, commit.
set -euo pipefail
R="${PARENA_K8S_RENDER:-/home/fatbaby/PARENA/parena-k8s-render}"
OUT="$(cd "$(dirname "$0")" && pwd)/clusters/prrject-fatbaby"
REG="us-central1-docker.pkg.dev/project-d24a71e9-2daf-4b2d-917/emily"
TAG="${COLLECTIONS_TAG:-0.1.0}"
printf 'apiVersion: v1\nkind: Namespace\nmetadata:\n  name: emily\n' > "$OUT/00-namespace.yaml"
# No app renders its own Ingress: every Ingress object is another billed external LB, so all hosts are
# rules in the ONE edge Ingress (specs/edge.ingress -> 90-edge-ingress.yaml). Backends share a namespace.
"$R" --name collections-server --namespace emily --image "$REG/collections-server:$TAG" \
     --port 8087 --cpu-milli 100 --memory-mi 128 > "$OUT/10-collections-server.yaml"

# FatBaby: the whole pipeline as one multi-container pod (unix sockets in a shared emptyDir).
P="${PARENA_POD_RENDER:-/home/fatbaby/PARENA/parena-pod-render}"
FB_TAG="${FATBABY_TAG:-0.1.0}"
"$P" "$(dirname "$0")/specs/fatbaby-core.pod" | sed "s/:IMAGE_TAG\$/:$FB_TAG/" > "$OUT/20-fatbaby-core.yaml"
"$P" "$(dirname "$0")/specs/edge.ingress" > "$OUT/90-edge-ingress.yaml"
