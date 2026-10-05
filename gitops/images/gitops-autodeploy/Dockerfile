# gitops-autodeploy (K8S-CD-01 Phase 2) -- runs as an in-cluster CronJob, same spirit as
# images/gitops-sync but pointed the other direction (detect a new pushed tag -> regenerate +
# commit, instead of pull a commit -> kubectl apply).
FROM google/cloud-sdk:alpine
RUN apk add --no-cache git bash openssh-client
COPY parena-k8s-render parena-pod-render /usr/local/bin/
COPY entrypoint.sh /app/entrypoint.sh
# GitHub's own published SSH host keys (docs.github.com/.../githubs-ssh-key-fingerprints), baked
# in at build time rather than trusted-on-first-connect at runtime -- so a clone over SSH (see
# entrypoint.sh) can pin UserKnownHostsFile without relying on whatever the Secret happens to
# carry or a live ssh-keyscan each run.
RUN mkdir -p /etc/gitops-autodeploy/ssh \
    && printf '%s\n' \
       'github.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCj7ndNxQowgcQnjshcLrqPEiiphnt+VTTvDP6mHBL9j1aNUkY4Ue1gvwnGLVlOhGeYrnZaMgRK6+PKCUXaDbC7qtbW8gIkhL7aGCsOr/C56SJMy/BCZfxd1nWzAOxSDPgVsmerOBYfNqltV9/hWCqBywINIR+5dIg6JTJ72pcEpEjcYgXkE2YEFXV1JHnsKgbLWNlhScqb2UmyRkQyytRLtL+38TGxkxCflmO+5Z8CSSNY7GidjMIZ7Q4zMjA2n1nGrlTDkzwDCsw+wqFPGQA179cnfGWOWRVruj16z6XyvxvjJwbz0wQZ75XK5tKSb7FNyeIEs4TT4jk+S4dhPeAUC5y+bDYirYgM4GC7uEnztnZyaVWQ7B381AK4Qdrwt51ZqExKbQpTUNn+EjqoTwvqNj4kqx5QUCI0ThS/YkOxJCXmPUWZbhjpCg56i+2aB6CmK2JGhn57K5mj0MNdBXA4/WnwH6XoPWJzK5Nyu2zB3nAZp+S5hpQs+p1vN1/wsjk=' \
       'github.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=' \
       'github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl' \
       > /etc/gitops-autodeploy/ssh/known_hosts
RUN chmod +x /usr/local/bin/parena-k8s-render /usr/local/bin/parena-pod-render /app/entrypoint.sh \
    && addgroup -g 10001 gitops && adduser -D -u 10001 -G gitops gitops \
    && chown -R 10001:10001 /etc/gitops-autodeploy
USER 10001
ENTRYPOINT ["/app/entrypoint.sh"]
