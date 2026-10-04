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
