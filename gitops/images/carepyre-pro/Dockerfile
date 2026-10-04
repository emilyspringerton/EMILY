# idunapro (CarePyre console backend): pure-Go (modernc sqlite). State on the PVC at /app/var; migrations ship in the image.
FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOWORK=off go build -o /out/idunapro .
FROM debian:12-slim
RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /out/idunapro /app/idunapro
COPY migrations /app/migrations
ENV IDUNA_PRO_ROOT=/app
CMD ["/app/idunapro"]
