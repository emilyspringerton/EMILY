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
"$R" --name collections-server --namespace emily --image "$REG/collections-server:$TAG" \
     --port 8087 --cpu-milli 100 --memory-mi 128 > "$OUT/10-collections-server.yaml"

# FatBaby: the whole pipeline as one multi-container pod (unix sockets in a shared emptyDir).
P="${PARENA_POD_RENDER:-/home/fatbaby/PARENA/parena-pod-render}"
FB_TAG="${FATBABY_TAG:-0.1.0}"
"$P" "$(dirname "$0")/specs/fatbaby-core.pod" | sed "s/:IMAGE_TAG\$/:$FB_TAG/" > "$OUT/20-fatbaby-core.yaml"

# REDGARDEN game servers (stateless pods, UDP LoadBalancer each; images from REDGARDEN/scripts/build-image.sh).
RG_TAG="${RG_STABLE_TAG:-3e55798}"; RGR_TAG="${RG_RND_TAG:-3e55798}"
"$P" "$(dirname "$0")/specs/redgarden-stable.pod" | sed "s/:IMAGE_TAG\$/:$RG_TAG/" > "$OUT/30-redgarden-stable.yaml"
"$P" "$(dirname "$0")/specs/redgarden-rnd.pod" | sed "s/:IMAGE_TAG\$/:$RGR_TAG/" > "$OUT/31-redgarden-rnd.yaml"
"$P" "$(dirname "$0")/specs/gfd-wsrelay.pod" > "$OUT/32-gfd-wsrelay.yaml"
"$P" "$(dirname "$0")/specs/edge.gateway" > "$OUT/91-edge-gateway.yaml"

# GFD (DragonsNShit MUD + server-go): one pod on a PVC, TCP+UDP LoadBalancers, IDUNA via the PARENA secure-channel sidecar.
GFD_TAG="${GFD_TAG:-v2}"
"$P" "$(dirname "$0")/specs/gfd-core.pod" | sed "s/:IMAGE_TAG\$/:$GFD_TAG/" > "$OUT/33-gfd-core.yaml"

# SHANKPIT (UDP FPS server + zombie sandbox + queue bot pool): one stateless pod, UDP LoadBalancer (image: SHANKPIT/scripts/build-image.sh).
SP_TAG="${SHANKPIT_TAG:-v1}"
"$P" "$(dirname "$0")/specs/shankpit.pod" | sed "s/:IMAGE_TAG\$/:$SP_TAG/" > "$OUT/34-shankpit.yaml"

# IDUNA (K8S-MV-01): the IAM hub on a PVC; ClusterIP only (front-door pod fronts it). Image: IDUNA/scripts/build-image.sh.
ID_TAG="${IDUNA_TAG:-v2}"
"$P" "$(dirname "$0")/specs/iduna.pod" | sed "s/:IMAGE_TAG\$/:$ID_TAG/" > "$OUT/40-iduna.yaml"

# IDUNA public front door (K8S-MV-01): iam./console.okemily.com on the edge Gateway (cert: okemily-com-wild in certmap
# edge-certs). The renderer emits one HTTPRoute per service, so this one (two hosts + iam's "/" -> SSO login rewrite,
# as a 302 to the SSO login instead of nginx's internal rewrite - GKE Gateway only allows ReplacePrefixMatch rewrites) is written here. 1h backend timeout for console websockets.
cat > "$OUT/92-iduna-routes.yaml" <<'YAML'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: iduna
  namespace: emily
spec:
  parentRefs:
    - name: edge-gw
  hostnames:
    - iam.okemily.com
    - console.okemily.com
  rules:
    - backendRefs:
        - name: iduna
          port: 8080
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: iduna-iam-root
  namespace: emily
spec:
  parentRefs:
    - name: edge-gw
  hostnames:
    - iam.okemily.com
  rules:
    - matches:
        - path: {type: Exact, value: /}
      filters:
        - type: RequestRedirect
          requestRedirect:
            path: {type: ReplaceFullPath, replaceFullPath: /api/v1/auth/sso/login}
            statusCode: 302
---
apiVersion: networking.gke.io/v1
kind: GCPBackendPolicy
metadata:
  name: iduna
  namespace: emily
spec:
  default:
    timeoutSec: 3600
  targetRef:
    group: ""
    kind: Service
    name: iduna
YAML

# DEADWEIGHT + WOTAN (K8S-MV-04 / K8S-MV-03). Images: DEADWEIGHT/scripts/build-image.sh, WOTAN/scripts/build-image.sh.
DW_TAG="${DEADWEIGHT_TAG:-v1}"; WO_TAG="${WOTAN_TAG:-v4}"
"$P" "$(dirname "$0")/specs/deadweight.pod" | sed "s/:IMAGE_TAG\$/:$DW_TAG/" > "$OUT/41-deadweight.yaml"
"$P" "$(dirname "$0")/specs/wotan.pod" | sed "s/:IMAGE_TAG\$/:$WO_TAG/" > "$OUT/42-wotan.yaml"
# wotan.okemily.com on the edge Gateway: path routes replace nginx-wotan.conf's proxy locations (/api -> IDUNA, /DEADWEIGHT/ws -> bridge).
cat > "$OUT/93-wotan-routes.yaml" <<'YAML'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: wotan
  namespace: emily
spec:
  parentRefs:
    - name: edge-gw
  hostnames:
    - wotan.okemily.com
  rules:
    - matches:
        - path: {type: PathPrefix, value: /api/}
      backendRefs:
        - name: iduna
          port: 8080
    - matches:
        - path: {type: PathPrefix, value: /DEADWEIGHT/ws}
      backendRefs:
        - name: deadweight
          port: 8765
    - backendRefs:
        - name: wotan
          port: 80
---
apiVersion: networking.gke.io/v1
kind: GCPBackendPolicy
metadata:
  name: deadweight
  namespace: emily
spec:
  default:
    timeoutSec: 3600
  targetRef:
    group: ""
    kind: Service
    name: deadweight
---
apiVersion: networking.gke.io/v1
kind: HealthCheckPolicy
metadata:
  name: deadweight
  namespace: emily
spec:
  default:
    config:
      type: TCP
      tcpHealthCheck:
        port: 8765
  targetRef:
    group: ""
    kind: Service
    name: deadweight
YAML

# tcp-edge: the single TCP LB rule for every raw-TCP service (gfd 2323/2222/7171/7070, iduna tunnel 8443, deadweight 7180).
TE_TAG="${TCP_EDGE_TAG:-v1}"
"$P" "$(dirname "$0")/specs/tcp-edge.pod" | sed "s/:IMAGE_TAG\$/:$TE_TAG/" > "$OUT/43-tcp-edge.yaml"

# collections-server answers 404 on / (the Gateway's default probe) -> probe /healthz so the backend is healthy.
cat > "$OUT/94-collections-health.yaml" <<'YAML'
apiVersion: networking.gke.io/v1
kind: HealthCheckPolicy
metadata:
  name: collections-server
  namespace: emily
spec:
  default:
    config:
      type: HTTP
      httpHealthCheck:
        port: 8087
        requestPath: /healthz
  targetRef:
    group: ""
    kind: Service
    name: collections-server
YAML

# okemily.com apex + www (K8S-MV-02): IDUNA API paths -> iduna:8080, /gfd-ws/<port> -> gfd-wsrelay (path mode), /news -> news.okemily.com,
# everything else -> the nginx sidecar (iduna:80) that serves the static site off IDUNA's PVC.
cat > "$OUT/95-okemily-apex.yaml" <<'YAML'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: okemily-apex
  namespace: emily
spec:
  parentRefs:
    - name: edge-gw
  hostnames:
    - okemily.com
    - www.okemily.com
  rules:
    - matches:
        - path: {type: Exact, value: /.well-known/jwks.json}
        - path: {type: PathPrefix, value: /api}
        - path: {type: PathPrefix, value: /services}
        - path: {type: PathPrefix, value: /admin}
        - path: {type: PathPrefix, value: /portal}
        - path: {type: PathPrefix, value: /q}
        - path: {type: PathPrefix, value: /nock}
      backendRefs:
        - name: iduna
          port: 8080
    - matches:
        - path: {type: PathPrefix, value: /gfd-ws}
      backendRefs:
        - name: gfd-wsrelay
          port: 8081
    - matches:
        - path: {type: PathPrefix, value: /news}
      filters:
        - type: RequestRedirect
          requestRedirect:
            hostname: news.okemily.com
            path: {type: ReplacePrefixMatch, replacePrefixMatch: /}
            statusCode: 301
    - backendRefs:
        - name: iduna
          port: 80
YAML

# gitops-sync (K8S-CD-01): auto-deploy. A CronJob pulls the public EMILY repo and kubectl-applies this very directory
# (minus 00-namespace) every 3 min. Namespaced Role, named kinds only (no secrets), no prune. Image: images/gitops-sync/build.sh.
SYNC_TAG="${GITOPS_SYNC_TAG:-v1}"
cat > "$OUT/96-gitops-sync.yaml" <<YAML
apiVersion: v1
kind: ServiceAccount
metadata:
  name: gitops-sync
  namespace: emily
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: gitops-sync
  namespace: emily
rules:
  - apiGroups: [""]
    resources: [services, persistentvolumeclaims]
    verbs: [get, list, create, update, patch]
  - apiGroups: [apps]
    resources: [deployments]
    verbs: [get, list, create, update, patch]
  - apiGroups: [batch]
    resources: [cronjobs]
    verbs: [get, list, create, update, patch]
  - apiGroups: [gateway.networking.k8s.io]
    resources: [gateways, httproutes]
    verbs: [get, list, create, update, patch]
  - apiGroups: [networking.gke.io]
    resources: [healthcheckpolicies, gcpbackendpolicies]
    verbs: [get, list, create, update, patch]
  - apiGroups: [rbac.authorization.k8s.io]
    resources: [roles, rolebindings]
    verbs: [get, list, create, update, patch]
  - apiGroups: [""]
    resources: [serviceaccounts]
    verbs: [get, list, create, update, patch]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: gitops-sync
  namespace: emily
subjects:
  - kind: ServiceAccount
    name: gitops-sync
    namespace: emily
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: gitops-sync
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: gitops-sync
  namespace: emily
spec:
  schedule: "*/3 * * * *"
  concurrencyPolicy: Forbid
  successfulJobsHistoryLimit: 2
  failedJobsHistoryLimit: 3
  jobTemplate:
    spec:
      activeDeadlineSeconds: 240
      backoffLimit: 0
      template:
        spec:
          serviceAccountName: gitops-sync
          restartPolicy: Never
          containers:
            - name: sync
              image: $REG/gitops-sync:$SYNC_TAG
              resources:
                requests: {cpu: 20m, memory: 64Mi, ephemeral-storage: 256Mi}
                limits: {memory: 128Mi, ephemeral-storage: 512Mi}
YAML

# CarePyre (K8S-MV-05): carepyre.org + www on the edge Gateway (cert carepyre-org-wild in certmap edge-certs). Site + idunapro in one pod.
CP_TAG="${CAREPYRE_TAG:-v1}"
"$P" "$(dirname "$0")/specs/carepyre.pod" | sed "s/:IMAGE_TAG\$/:$CP_TAG/" > "$OUT/44-carepyre.yaml"
cat > "$OUT/97-carepyre-routes.yaml" <<'YAML'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: carepyre
  namespace: emily
spec:
  parentRefs:
    - name: edge-gw
  hostnames:
    - carepyre.org
    - www.carepyre.org
  rules:
    - backendRefs:
        - name: carepyre
          port: 80
YAML
