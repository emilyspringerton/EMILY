# Secrets are NOT in this repo

GitOps here is pull-based: anything committed is applied, so a Secret value in git is a leaked secret.
The `fatbaby-env` Secret that `fatbaby-core` mounts (`envFrom`, **optional**) is created out of band:

```bash
# from the box/operator machine with a kubeconfig for the cluster
./make-fatbaby-secret.sh ~/.config/fatbaby/env ~/.config/fatbaby-newssite/env | kubectl apply -f -
```

State of play (2026-10-04): neither env file exists on the current box — the pipeline runs with no
secrets file today — so the Secret is optional and the pod starts without it. `keys.txt` lists the
variables the code reads, so a value can be added without reading Go. The real fix (GCP Secret Manager +
CSI driver, per-container least privilege instead of one shared Secret) is K8S `S207-06`, not done.

## gitops-autodeploy (K8S-CD-01 Phase 2)

The CronJob (`gitops/render.sh`'s own `96b-gitops-autodeploy.yaml` output) needs two credentials.

1. **Artifact Registry read access** -- DONE 2026-10-05 (live, real GCP resources, not just
   documented here): service account `gitops-autodeploy@project-d24a71e9-2daf-4b2d-917
   .iam.gserviceaccount.com` exists, holds `roles/artifactregistry.reader` on the project, and is
   bound via Workload Identity to the in-cluster KSA (`project-d24a71e9-2daf-4b2d-917.svc.id.goog
   [emily/gitops-autodeploy]`). The ServiceAccount in the rendered manifest carries the real
   annotation, not a placeholder.
2. **Git push access** (a GitHub PAT with `repo` scope on `emilyspringerton/EMILY`) -- **not yet
   done**, needs a real GitHub account action no sandbox/CI identity can perform. Create the PAT,
   then:
   ```bash
   ./make-gitops-autodeploy-secret.sh <the-PAT> | kubectl apply -f -
   ```
   The pod `CrashLoopBackOff`s every 5 min until this exists -- expected, bounded
   (`failedJobsHistoryLimit: 3`), and resolves itself the next scheduled run once the Secret lands.

Image build is automatic (`.github/workflows/build-gitops-autodeploy.yml`, CI auth via the
`github` Workload Identity pool -- see `gitops/CI_SETUP.md`), and was also verified by a real,
live `gcloud builds submit` + push during this work (`gitops-autodeploy:v1` is live in Artifact
Registry right now).
