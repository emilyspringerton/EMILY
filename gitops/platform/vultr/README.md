# Vultr Kubernetes Engine — the cluster after GCP

**Why:** GCP was shut down for billing on 2026-11-26. Nothing survived: not the GKE cluster `prrject-fatbaby`, not
Artifact Registry, not the PVC data, not the `gs://` buckets. Founder: *"we need to convert everything to VULTR /
VULTR has kubernetes"*. This is a **fresh start**, not a data migration.

**Status (2026-11-26): repo-side prep is done; nothing is deployed yet.** There is no Vultr API key or cluster yet,
so every manifest here has been rendered and schema-validated, but none has been applied to a real cluster.

## Layout

| Path | What |
|---|---|
| `../../clusters/vultr-emily/` | the cluster, **generated**: `PROVIDER=vultr ../../render.sh`. gitops-sync applies it. |
| `../../clusters/prrject-fatbaby/` | the dead GKE cluster, kept as history. `./render.sh` (PROVIDER=gke) still reproduces it byte-for-byte. |
| `bootstrap.sh` | one-time setup: Envoy Gateway, cert-manager, TLS, registry pull secret, first apply |
| `attach-game-ip.sh` | Reserved IP → game node, plus `game-ip.yaml` to configure it on the node |
| `gatewayclass.yaml`, `cluster-issuer.yaml`, `certificate.yaml`, `game-ip.yaml` | cluster-scoped and platform objects that gitops-sync's namespaced Role can't (and shouldn't) apply |

## What changed from GKE, and why

| GKE | Vultr | Where |
|---|---|---|
| Artifact Registry `us-central1-docker.pkg.dev/…/emily` | Vultr Container Registry, `$VULTR_REGISTRY` (default `ewr.vultrcr.com/emily`), pulled with `vultr-cr-credentials` | `render.sh` swaps the registry. `bootstrap.sh` patches the default ServiceAccount. |
| GKE Gateway `gke-l7-global-external-managed` + Certificate Manager certmap + static global IP | Envoy Gateway (class `eg`) + cert-manager DNS-01 on Cloudflare → Secret `edge-tls`. The Envoy Service gets its own Vultr LB. | PARENA `pod-gateway-vultr-yaml`, `certificate.yaml` |
| GCPBackendPolicy `timeoutSec: 3600` (WebSockets) | one BackendTrafficPolicy on the Gateway, `requestTimeout: 3600s` | same |
| HealthCheckPolicy (collections `/healthz`, deadweight/gfd TCP, signals) | dropped. Envoy routes to ready endpoints and never probes `/`. | `vultr_postprocess.py` |
| 6 UDP + 1 TCP `LoadBalancer` Services sharing `34.63.32.219` | **hostPorts on a `game` node pool carrying one Vultr Reserved IP.** Only `gfd-ssh-tcp` (public :22 → tcp-edge:2222) remains a Vultr LB, because :22 on the node itself is the node's sshd. | PARENA `pod-hostports-container-yaml`, `attach-game-ip.sh` |
| PVCs 1–20Gi on GKE PD | `vultr-block-storage` (VKE default, NVMe). Raised to Vultr's 10 GB minimum. | PARENA host |
| vLLM coder engine on an L4 | **dropped for now** (founder). CODER-AGENT stays parked. | `render.sh` |
| gitops-autodeploy CronJob (Workload Identity → Artifact Registry) | **not rendered.** Its VCR port is still to do. Move tags by hand: `bump-tag.sh`, then `PROVIDER=vultr ./render.sh`. | `render.sh` |
| gitops-sync watching `clusters/prrject-fatbaby` | same image, `GITOPS_DIR=gitops/clusters/vultr-emily`; Role swaps GKE policy kinds for BackendTrafficPolicy | `vultr_postprocess.py` |

### Why game traffic doesn't go through Vultr load balancers

Checked against the Vultr CCM source (`vultr/vultr-cloud-controller-manager`, `vultr/loadbalancers.go`) and Vultr's docs:

- **15 forwarding rules per LB.** `redgarden-rnd` publishes 43 UDP ports and `redgarden-stable` publishes 21, so they would need five LBs, each with its own IP.
- **The default health check is TCP on the first port's NodePort.** A UDP-only Service has no TCP listener there, so every node fails the check.
- **The CCM never reads `externalTrafficPolicy`.** Its health check also can't target the Service's `healthCheckNodePort`.

hostPorts avoid all three problems. They give one IP for every game port, as before, and real client source
addresses (the matchmakers reply per client). They also cost nothing per port. The trade-off is that the IP is tied
to one node, so `attach-game-ip.sh` has to be re-run whenever VKE replaces that node.

## Bring-up, in order

1. **Vultr account:** API key, a VKE cluster, a Container Registry, and a node pool named for games with label
   **`emily/pool=game`** and exactly 1 node. Size the game node for redgarden ×2 + gfd + shankpit + einhorn + tcp-edge.
   The other workloads can share it or use a general pool.
2. **Images.** The Artifact Registry is gone, so all ~25 images must be rebuilt from source and pushed to VCR
   (see *Not done* below). The tags in `tags.env` are what the manifests expect.
3. `PROVIDER=vultr VULTR_REGISTRY=<host>/<name> ./render.sh`. Commit the result if the registry name differs from
   the default.
4. `CLOUDFLARE_API_TOKEN=… ACME_EMAIL=… VCR_HOST=… VCR_USER=… VCR_PASS=… ./bootstrap.sh`
5. `VULTR_API_KEY=… ./attach-game-ip.sh`, then set the DNS A records: game hosts → game IP; web hosts → Envoy LB IP
   (bootstrap.sh prints where to find it).
6. App Secrets, out of band: `../../secrets/README.md`.

## Verify on first deploy (not checkable without a cluster)

- **hostPort works on VKE's CNI.** VKE defaults to Calico, whose standard install chains the `portmap` plugin. To
  test, send `nc -u <game-ip> 6969` traffic and watch SHANKPIT's logs.
- **Whether Vultr's firewall or node security group lets the game ports reach the node.** Open them if not.
- **game-ip DaemonSet:** after attach, `ip addr` on the game node should show the Reserved IP exactly once.
- **Envoy WebSockets** (GFD browser MUD, IDUNA console, DEADWEIGHT bridge) stay up past Envoy's 5-minute default
  stream-idle timeout. If they drop, add a ClientTrafficPolicy `timeout.http.idleTimeout`.
- **`vultr-block-storage` is the cluster's default StorageClass.** The PVCs don't name a class.

## Not done (tracked in EMILY/BACKLOG.md)

- **Image builds → VCR.** Every repo's `scripts/build-image.sh` and `.github/workflows/*` use
  `gcloud builds submit` and Workload Identity Federation, both dead. These need `docker buildx` + VCR credentials as
  GitHub secrets.
- **gitops-autodeploy → VCR** (tag polling through the VCR API instead of Artifact Registry).
- **`emily backup`** (`emily.cli/cmd/backup.go`) still shells out to `gcloud storage cp` to
  `gs://…-backups`. Its replacement is Vultr Object Storage (S3 API).
- **Google *application* APIs are not infrastructure and are untouched:** Google SSO, Drive, Vertex and Gemini in
  IDUNA and emily.cli. Whether their API project survived the billing shutdown is a separate question.
