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

Two one-time setup steps, neither can be done from a sandbox without GCP/GitHub account access --
the CronJob (`gitops/render.sh`'s own `96b-gitops-autodeploy.yaml` output) will `CrashLoopBackOff`
without both:

1. **Git push access** (a GitHub PAT with `repo` scope on `emilyspringerton/EMILY`):
   ```bash
   ./make-gitops-autodeploy-secret.sh <the-PAT> | kubectl apply -f -
   ```
2. **Artifact Registry read access** (Workload Identity -- the `96b-gitops-autodeploy.yaml`
   ServiceAccount is pre-annotated with a `CHANGEME` GCP service account email; create that GSA,
   grant it `roles/artifactregistry.reader` on the project, bind it to the KSA):
   ```bash
   gcloud iam service-accounts create gitops-autodeploy
   gcloud projects add-iam-policy-binding project-d24a71e9-2daf-4b2d-917 \
     --member="serviceAccount:gitops-autodeploy@project-d24a71e9-2daf-4b2d-917.iam.gserviceaccount.com" \
     --role="roles/artifactregistry.reader"
   gcloud iam service-accounts add-iam-policy-binding \
     gitops-autodeploy@project-d24a71e9-2daf-4b2d-917.iam.gserviceaccount.com \
     --role="roles/iam.workloadIdentityUser" \
     --member="serviceAccount:project-d24a71e9-2daf-4b2d-917.svc.id.goog[emily/gitops-autodeploy]"
   ```
   then replace the `CHANGEME@...` placeholder in the rendered manifest with the real GSA email
   (or template it into `render.sh` once it's created, so it survives the next `render.sh` run).
