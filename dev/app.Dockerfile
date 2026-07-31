# Dev image for local docker compose (make dev-docker).
# Includes Go 1.26+ and Node 20+ (required for maily-builder / Vite 6).
FROM golang:1.26.1-bookworm

ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV PATH="${GOPATH}/bin:/usr/local/go/bin:${PATH}"

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl \
  && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
  && apt-get install -y --no-install-recommends nodejs \
  && corepack enable \
  && corepack prepare yarn@1.22.22 --activate \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
CMD ["sleep", "infinity"]
