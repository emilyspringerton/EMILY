# CI -> GCP auth (Workload Identity Federation)

K8S-CD-01: lets GitHub Actions build+push images to Artifact Registry with no stored key --
GCP trusts GitHub's own OIDC token directly. **Live as of 2026-10-05** (created for real, from
Cloud Shell with the founder's own `roles/owner` session, not just documented):

| Resource | Value |
|---|---|
| Workload Identity Pool | `projects/532442865445/locations/global/workloadIdentityPools/github` |
| OIDC Provider | `projects/532442865445/locations/global/workloadIdentityPools/github/providers/github` |
| Trust condition | `assertion.repository_owner == 'emilyspringerton'` -- any repo under the org can mint a token, no per-repo re-setup needed |
| Service account (CI) | `github-ci@project-d24a71e9-2daf-4b2d-917.iam.gserviceaccount.com` -- `roles/cloudbuild.builds.editor` + `roles/artifactregistry.writer`, bound to the pool via `roles/iam.workloadIdentityUser` |
| Service account (in-cluster) | `gitops-autodeploy@project-d24a71e9-2daf-4b2d-917.iam.gserviceaccount.com` -- `roles/artifactregistry.reader`, bound to KSA `emily/gitops-autodeploy` via GKE Workload Identity (unrelated to the GitHub pool above -- this one's for the CronJob reading the registry, not CI pushing to it) |

Cloud Build's own default service account (`532442865445@cloudbuild.gserviceaccount.com`) already
holds `roles/cloudbuild.builds.builder`, which was already enough to push images before any of
this existed (manual builds have worked all along) -- untouched, nothing needed there.

## What every repo's workflow needs (GitHub-side, not done by this pass)

Two repo (or org) **Variables** (Settings -> Secrets and variables -> Actions -> Variables --
these aren't secrets, the values aren't sensitive):

```
GCP_WIF_PROVIDER = projects/532442865445/locations/global/workloadIdentityPools/github/providers/github
GCP_CI_SERVICE_ACCOUNT = github-ci@project-d24a71e9-2daf-4b2d-917.iam.gserviceaccount.com
```

Setting them once at the **organization** level (Settings -> Actions -> Variables, if
`emilyspringerton` is an org account) makes every repo's workflow below work with zero per-repo
config. Per-repo Variables work the same way if it's a personal account instead.

Each workflow then just needs:
```yaml
permissions:
  contents: read
  id-token: write    # required for Workload Identity -- without this the auth step fails
steps:
  - uses: google-github-actions/auth@v2
    with:
      workload_identity_provider: ${{ vars.GCP_WIF_PROVIDER }}
      service_account: ${{ vars.GCP_CI_SERVICE_ACCOUNT }}
  - uses: google-github-actions/setup-gcloud@v2
```

## Recreating this (if the pool/SAs are ever deleted)

```bash
PROJECT=project-d24a71e9-2daf-4b2d-917
gcloud services enable sts.googleapis.com --project "$PROJECT"

gcloud iam workload-identity-pools create github --project "$PROJECT" --location=global \
  --display-name="GitHub Actions"
gcloud iam workload-identity-pools providers create-oidc github --project "$PROJECT" \
  --location=global --workload-identity-pool=github --display-name="GitHub" \
  --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository,attribute.repository_owner=assertion.repository_owner" \
  --attribute-condition="assertion.repository_owner == 'emilyspringerton'" \
  --issuer-uri="https://token.actions.githubusercontent.com"

gcloud iam service-accounts create github-ci --project "$PROJECT"
gcloud projects add-iam-policy-binding "$PROJECT" --condition=None \
  --member="serviceAccount:github-ci@$PROJECT.iam.gserviceaccount.com" --role="roles/cloudbuild.builds.editor"
gcloud projects add-iam-policy-binding "$PROJECT" --condition=None \
  --member="serviceAccount:github-ci@$PROJECT.iam.gserviceaccount.com" --role="roles/artifactregistry.writer"
gcloud iam service-accounts add-iam-policy-binding "github-ci@$PROJECT.iam.gserviceaccount.com" \
  --project "$PROJECT" --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/532442865445/locations/global/workloadIdentityPools/github/attribute.repository_owner/emilyspringerton"
```

(The `gitops-autodeploy` service account + its GKE Workload Identity binding is documented
separately in `gitops/secrets/README.md` -- different purpose, different trust relationship.)

## `gcloud builds submit` needs more than the two base roles (found live, 2026-10-05)

`roles/cloudbuild.builds.editor` + `roles/artifactregistry.writer` are NOT sufficient on their
own for `gcloud builds submit` to work end to end as a non-Viewer/Owner CI SA. Three more real
gaps, found live debugging SHANKPIT's first CI run, now fixed in every `build-image.sh`/
`build.sh` that calls `gcloud builds submit` (SHANKPIT, IDUNA, PRRJECT_FATBABY,
`gitops/images/gitops-autodeploy/build.sh`):

1. **Staging bucket access.** `github-ci` needs write/read access to the Cloud Build default
   staging bucket (`gs://<PROJECT>_cloudbuild`), which has Uniform Bucket-Level Access enabled so
   it must be granted via `gsutil iam ch`, not legacy ACLs:
   ```bash
   PROJECT=project-d24a71e9-2daf-4b2d-917
   SA="github-ci@$PROJECT.iam.gserviceaccount.com"
   for role in storage.objectAdmin storage.legacyBucketReader storage.legacyBucketWriter; do
     gsutil iam ch "serviceAccount:$SA:roles/$role" "gs://${PROJECT}_cloudbuild"
   done
   gcloud projects add-iam-policy-binding "$PROJECT" --condition=None \
     --member="serviceAccount:$SA" --role="roles/serviceusage.serviceUsageConsumer"
   ```

2. **Auto-detecting the staging bucket requires a project-level `storage.buckets.list` call**
   that the bucket-scoped bindings above don't satisfy -- `gcloud builds submit` does
   `GET /storage/v1/b?project=...&prefix=...` to find the bucket, a 403 there is misleadingly
   reported as `"forbidden ... serviceusage.services.use"`. Fix: every build script passes
   `--gcs-source-staging-dir="gs://${PROJECT}_cloudbuild/source"` explicitly, which skips that
   list call entirely. Granting a project-level `roles/storage.admin` (which includes
   `buckets.list`) would also work but is broader than this CI SA should hold.

3. **Cloud Build runs the worker as the Compute Engine default service account**
   (`<PROJECT_NUMBER>-compute@developer.gserviceaccount.com`) unless a custom build SA is
   configured. Submitting a build means "acting as" that SA, so `github-ci` needs:
   ```bash
   gcloud iam service-accounts add-iam-policy-binding \
     "${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" --project "$PROJECT" \
     --member="serviceAccount:$SA" --role="roles/iam.serviceAccountUser"
   ```

4. **Log streaming requires a project Viewer/Owner primitive role**, which `github-ci` is
   deliberately not. `gcloud builds submit` blocks on `--suppress-logs` the same way (still
   polls/waits, still exits 1) -- the real fix is `--async` + manually polling
   `gcloud builds describe <id> --format="value(status)"` for a terminal status, never
   streaming logs. See any of the `build-image.sh` scripts above for the pattern.

## What's still manual

Only the things that need a human with GitHub account access, which no CI identity or GCP
credential can substitute for:
- The two repo/org Variables above.
- A GitHub PAT for `gitops-autodeploy`'s own git-push Secret (`gitops/secrets/README.md`).
- Adding `.github/workflows/deploy.yml` (build+push on push to main) to each additional service
  repo beyond IDUNA/PRRJECT_FATBABY/EMILY as it's ready -- the pattern is the same each time: see
  those repos' own workflow files.
