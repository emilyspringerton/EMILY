# gitops-sync: PARENA's own pull-based reconciler (parena-gitops, stdlib/k8s/gitops.prn decides,
# tools/gitops_host.c plumbs) + git + kubectl. parena-gitops is statically linked (see build.sh) so it runs on alpine.
FROM alpine:3.20
RUN apk add --no-cache git kubectl ca-certificates
COPY parena-gitops /usr/local/bin/parena-gitops
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /usr/local/bin/parena-gitops /app/entrypoint.sh && adduser -D -u 10001 gitops
USER 10001
ENTRYPOINT ["/app/entrypoint.sh"]
