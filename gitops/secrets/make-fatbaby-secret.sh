#!/usr/bin/env bash
# Print a Secret manifest (namespace emily, name fatbaby-env) built from one or more KEY=VALUE env files.
# Never commit the output. Pipe it to `kubectl apply -f -`. Missing files are skipped (the pod's
# secretRef is optional). Values are base64'd without being echoed to the terminal.
set -euo pipefail
echo "apiVersion: v1
kind: Secret
metadata:
  name: fatbaby-env
  namespace: emily
type: Opaque
data:"
for f in "$@"; do
  [ -f "$f" ] || { echo "skip $f (not found)" >&2; continue; }
  while IFS='=' read -r k v; do
    [[ -z "$k" || "$k" == \#* ]] && continue
    printf '  %s: %s\n' "$k" "$(printf %s "$v" | base64 -w0)"
  done < "$f"
done
