# gitops-autodeploy (K8S-CD-01 Phase 2) -- runs as an in-cluster CronJob, same spirit as
# images/gitops-sync but pointed the other direction (detect a new pushed tag -> regenerate +
# commit, instead of pull a commit -> kubectl apply).
FROM google/cloud-sdk:alpine
RUN apk add --no-cache git bash
COPY parena-k8s-render parena-pod-render /usr/local/bin/
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /usr/local/bin/parena-k8s-render /usr/local/bin/parena-pod-render /app/entrypoint.sh \
    && adduser -D -u 10001 gitops
USER 10001
ENTRYPOINT ["/app/entrypoint.sh"]
