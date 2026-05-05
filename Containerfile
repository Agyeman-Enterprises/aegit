# AeGit — Sovereign Forge — multi-stage Podman build from source
# Stage 1: Node 22 frontend (pnpm + vite)
# Stage 2: Go 1.26 binary (CGO/SQLite)
# Stage 3: Alpine runtime (non-root git user)

# ── Stage 1: Frontend ───────────────────────────────────────────────────────
FROM node:22-alpine AS frontend

WORKDIR /source
COPY package.json pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile

COPY . .
RUN pnpm exec vite build

# ── Stage 2: Go binary ──────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache build-base git

WORKDIR /source

# Download deps in a separate layer for cache reuse
COPY go.mod go.sum ./
RUN go mod download

# Copy Go source + built frontend assets
COPY main.go Makefile ./
COPY build      ./build
COPY cmd        ./cmd
COPY models     ./models
COPY modules    ./modules
COPY routers    ./routers
COPY services   ./services
COPY templates  ./templates
COPY options    ./options
COPY --from=frontend /source/public ./public

RUN GOEXPERIMENT=jsonv2 \
    CGO_ENABLED=1 \
    CGO_CFLAGS="-DSQLITE_MAX_VARIABLE_NUMBER=32766" \
    go build \
      -tags 'sqlite sqlite_unlock_notify' \
      -ldflags '-s -w -X "main.Version=stage2" -X "main.Tags=sqlite,sqlite_unlock_notify"' \
      -o aegit

# ── Stage 3: Runtime ────────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache \
    bash git curl \
    openssh-client openssh-server \
    gettext jq tzdata ca-certificates && \
    addgroup -S -g 1000 git && \
    adduser  -S -H -D -u 1000 -G git -s /bin/bash git && \
    mkdir -p /app/aegit /data /tmp/aegit && \
    chown -R git:git /app/aegit /data /tmp/aegit

# Binary
COPY --chown=git:git --from=builder  /source/aegit                    /app/aegit/aegit
# Runtime assets (templates, public, options, migration schemas)
COPY --chown=git:git --from=frontend /source/public                   /app/aegit/public
COPY --chown=git:git --from=builder  /source/templates                /app/aegit/templates
COPY --chown=git:git --from=builder  /source/options                  /app/aegit/options
COPY --chown=git:git --from=builder  /source/modules/migration/schemas /app/aegit/modules/migration/schemas

# Symlinks: aegit is primary; gitea stays for compatibility
RUN ln -sf /app/aegit/aegit /usr/local/bin/aegit && \
    ln -sf /app/aegit/aegit /usr/local/bin/gitea && \
    chmod +x /app/aegit/aegit

# Runtime scripts and config template
COPY custom/conf/app.ini.template /app/app.ini.template
COPY scripts/first-boot.sh        /usr/local/bin/first-boot.sh
COPY scripts/mirror-github.sh     /usr/local/bin/mirror-github.sh
COPY entrypoint.sh                /entrypoint.sh
RUN chmod +x /usr/local/bin/first-boot.sh \
             /usr/local/bin/mirror-github.sh \
             /entrypoint.sh

ENV GITEA_WORK_DIR=/app/aegit

USER git
WORKDIR /app/aegit

VOLUME ["/data"]
EXPOSE 8085 2222

ENTRYPOINT ["/entrypoint.sh"]
