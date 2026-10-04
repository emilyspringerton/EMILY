FROM alpine:3.20
RUN apk add --no-cache git kubectl ca-certificates
COPY sync.sh /sync.sh
RUN chmod +x /sync.sh && adduser -D -u 10001 gitops
USER 10001
ENTRYPOINT ["/sync.sh"]
