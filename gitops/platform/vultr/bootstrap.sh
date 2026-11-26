#!/usr/bin/env bash
# bootstrap.sh -- one-time platform setup for a fresh Vultr Kubernetes Engine cluster, then the first apply of
# clusters/vultr-emily. After this, gitops-sync (96-gitops-sync.yaml) keeps the cluster at EMILY HEAD.
# Idempotent: every step is an install-or-upgrade / apply, so re-running on a half-done cluster finishes it.
#
# Needs: kubectl pointed at the VKE cluster (KUBECONFIG from the Vultr console), helm 3, curl, jq.
# env (required): CLOUDFLARE_API_TOKEN  Zone:DNS:Edit + Zone:Zone:Read on okemily.com and carepyre.org
#                 ACME_EMAIL            Let's Encrypt account contact
#                 VCR_HOST VCR_USER VCR_PASS   Vultr Container Registry pull credentials (console -> registry ->
#                                       Docker credentials); VCR_HOST e.g. ewr.vultrcr.com
# Versions pinned 2026-11-26 to releases older than two weeks.
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"
CLUSTER="$HERE/../../clusters/vultr-emily"
EG_VERSION=v1.9.3          # Envoy Gateway, released 2026-10-25
CM_VERSION=v1.21.2         # cert-manager, released 2026-09-10
: "${CLOUDFLARE_API_TOKEN:?}" "${ACME_EMAIL:?}" "${VCR_HOST:?}" "${VCR_USER:?}" "${VCR_PASS:?}"
for t in kubectl helm jq; do command -v "$t" >/dev/null || { echo "bootstrap: $t not on PATH" >&2; exit 1; }; done
[ -d "$CLUSTER" ] || { echo "bootstrap: $CLUSTER missing -- run PROVIDER=vultr ../../render.sh first" >&2; exit 1; }

echo "== namespace"
kubectl apply -f "$CLUSTER/00-namespace.yaml"

echo "== game node pool"
if [ "$(kubectl get nodes -l emily/pool=game -o name | wc -l)" -lt 1 ]; then
  echo "bootstrap: no node labelled emily/pool=game. Add a VKE node pool with label emily/pool=game" >&2
  echo "  (Vultr console -> Kubernetes -> Node Pools, or the API's node-pools 'labels' field), then re-run." >&2
  exit 1
fi

echo "== Envoy Gateway $EG_VERSION (also installs the Gateway API CRDs)"
helm upgrade --install eg oci://docker.io/envoyproxy/gateway-helm --version "$EG_VERSION" \
  -n envoy-gateway-system --create-namespace --wait
kubectl apply -f "$HERE/gatewayclass.yaml"

echo "== cert-manager $CM_VERSION"
helm repo add jetstack https://charts.jetstack.io >/dev/null 2>&1 || true
helm repo update jetstack >/dev/null
helm upgrade --install cert-manager jetstack/cert-manager --version "$CM_VERSION" \
  -n cert-manager --create-namespace --set crds.enabled=true --wait
kubectl -n cert-manager create secret generic cloudflare-api-token \
  --from-literal=api-token="$CLOUDFLARE_API_TOKEN" --dry-run=client -o yaml | kubectl apply -f -
sed "s/\${ACME_EMAIL}/$ACME_EMAIL/" "$HERE/cluster-issuer.yaml" | kubectl apply -f -
kubectl apply -f "$HERE/certificate.yaml"

echo "== registry pull secret"
kubectl -n emily create secret docker-registry vultr-cr-credentials \
  --docker-server="$VCR_HOST" --docker-username="$VCR_USER" --docker-password="$VCR_PASS" \
  --dry-run=client -o yaml | kubectl apply -f -
# Every rendered workload runs as the namespace's default ServiceAccount, so one patch covers them all.
kubectl -n emily patch serviceaccount default -p '{"imagePullSecrets":[{"name":"vultr-cr-credentials"}]}'

echo "== first apply (gitops-sync takes over from here)"
for f in "$CLUSTER"/*.yaml; do
  [ "$(basename "$f")" = 00-namespace.yaml ] && continue
  kubectl apply -f "$f"
done

cat <<'NEXT'
== done. Still yours, out of band (never in git):
  - app Secrets: ../../secrets/README.md (fatbaby-env, edge-relay-env, iduna, ...); every pod mounts its Secret
    as optional, so pods start without them
  - game IP: VULTR_API_KEY=... ./attach-game-ip.sh, then point the game hosts' A records at it
  - web IP: kubectl -n envoy-gateway-system get svc  (EXTERNAL-IP of the envoy-emily-edge-gw-* Service);
    point okemily.com / *.okemily.com / carepyre.org hosts at it
  - kubectl -n emily get certificate edge-tls   (READY=True once DNS-01 completes, usually ~2 min)
NEXT
