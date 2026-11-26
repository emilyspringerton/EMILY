#!/usr/bin/env python3
"""PROVIDER=vultr post-pass for render.sh: strip what only exists on GKE from clusters/vultr-emily.

GCP was shut off (billing) on 2026-11-26 and the cluster moves to Vultr Kubernetes Engine. The PARENA renderer
already emits Vultr forms for what it owns (PARENA_K8S_PROVIDER=vultr); this handles the hand-written heredocs in
render.sh, which still describe GKE objects:

- BackendConfig / GCPBackendPolicy / HealthCheckPolicy have no Vultr form. Envoy Gateway routes on endpoint
  readiness and takes its WebSocket timeout from the Gateway's BackendTrafficPolicy (91-edge-gateway.yaml).
- loadBalancerIP is dropped: a Vultr LB is given its own address.
- gitops-sync: Role swaps the GKE policy kinds for BackendTrafficPolicy, its ServiceAccount pulls from the Vultr
  Container Registry, and it watches this cluster's directory instead of prrject-fatbaby.

usage: vultr_postprocess.py CLUSTER_DIR
"""
import os
import re
import sys

GKE_ONLY = re.compile(r"^apiVersion: (networking\.gke\.io|cloud\.google\.com)/", re.M)

ROLE_GKE = "  - apiGroups: [networking.gke.io]\n    resources: [healthcheckpolicies, gcpbackendpolicies]\n"
ROLE_EG = "  - apiGroups: [gateway.envoyproxy.io]\n    resources: [backendtrafficpolicies]\n"
SA_END = "  name: gitops-sync\n  namespace: emily\n---\napiVersion: rbac"
SA_PULL = ("  name: gitops-sync\n  namespace: emily\nimagePullSecrets:\n  - name: vultr-cr-credentials\n"
           "---\napiVersion: rbac")
SYNC_C = "        - name: sync\n"
SYNC_ENV = SYNC_C + "          env:\n            - {name: GITOPS_DIR, value: gitops/clusters/vultr-emily}\n"


def only_comments(doc):
    return all(not l.strip() or l.lstrip().startswith("#") for l in doc.splitlines())


def must_replace(text, old, new, what):
    if old not in text:
        sys.exit(f"vultr_postprocess: {what}: expected text not found (render.sh heredoc changed?)")
    return text.replace(old, new, 1)


def main():
    out = sys.argv[1]
    for f in sorted(os.listdir(out)):
        if not f.endswith(".yaml"):
            continue
        path = os.path.join(out, f)
        text = open(path).read()
        docs = [d for d in re.split(r"^---\n", text, flags=re.M)]
        docs = [d for d in docs if not GKE_ONLY.search(d) and not only_comments(d)]
        text = "---\n".join(docs)
        text = re.sub(r"^  loadBalancerIP: .*\n", "", text, flags=re.M)
        if f == "96-gitops-sync.yaml":
            text = must_replace(text, ROLE_GKE, ROLE_EG, "gitops-sync Role")
            text = must_replace(text, SA_END, SA_PULL, "gitops-sync ServiceAccount")
            text = must_replace(text, SYNC_C, SYNC_ENV, "gitops-sync container")
        if text.strip():
            open(path, "w").write(text)
        else:
            os.remove(path)


if __name__ == "__main__":
    main()
