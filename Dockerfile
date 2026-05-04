# AeGit Stage 1 — sovereign forge image
# Bases on official Gitea nightly (matches our v1.27.0-dev fork).
# Adds: aegit symlink, custom entrypoint, config template, mirror scripts.
# Stage 2 will build from source once AEGIS hooks are in place.

# syntax=docker/dockerfile:1
FROM gitea/gitea:nightly

USER root

# aegit binary alias + runtime tools
RUN apk add --no-cache gettext jq && \
    ln -sf /usr/local/bin/gitea /usr/local/bin/aegit

# Config template — envsubst fills ${SECRET_KEY} and ${INTERNAL_TOKEN} at container start
COPY custom/conf/app.ini.template /app/app.ini.template

# Container runtime scripts (Linux/Alpine — PowerShell exemption granted by OO)
COPY scripts/first-boot.sh    /usr/local/bin/first-boot.sh
COPY scripts/mirror-github.sh /usr/local/bin/mirror-github.sh
COPY entrypoint.sh            /entrypoint.sh

RUN chmod +x /usr/local/bin/first-boot.sh \
             /usr/local/bin/mirror-github.sh \
             /entrypoint.sh

VOLUME ["/data"]
EXPOSE 8001 22

ENTRYPOINT ["/entrypoint.sh"]

# ---- original Gitea build stages follow (unused in Stage 1, kept for source reference) ----
# To build from source in a future stage: docker build --target gitea-build .
FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.26-alpine3.23 AS gitea-build-source
RUN apk --no-cache add build-base git nodejs pnpm
WORKDIR /src
COPY package.json pnpm-lock.yaml .npmrc ./
RUN --mount=type=cache,target=/root/.local/share/pnpm/store pnpm install --frozen-lockfile
COPY --exclude=.git/ . .
RUN make frontend

# Build backend for each target platform
FROM docker.io/library/golang:1.26-alpine3.23 AS build-env

ARG GITEA_VERSION
ARG TAGS="sqlite sqlite_unlock_notify"
ENV TAGS="bindata timetzdata $TAGS"
ARG CGO_EXTRA_CFLAGS

# Build deps
RUN apk --no-cache add \
    build-base \
    git

WORKDIR ${GOPATH}/src/code.gitea.io/gitea
COPY go.mod go.sum ./
RUN go mod download
# Use COPY instead of bind mount as read-only one breaks makefile state tracking and read-write one needs binary to be moved as it's discarded.
# ".git" directory is mounted separately later only for version data extraction.
COPY --exclude=.git/ . .
COPY --from=frontend-build /src/public/assets public/assets

# Build gitea, .git mount is required for version data
RUN --mount=type=cache,target="/root/.cache/go-build" \
    --mount=type=bind,source=".git/",target=".git/" \
    make backend

COPY docker/root /tmp/local

# Set permissions for builds that made under windows which strips the executable bit from file
RUN chmod 755 /tmp/local/usr/bin/entrypoint \
              /tmp/local/usr/local/bin/* \
              /tmp/local/etc/s6/gitea/* \
              /tmp/local/etc/s6/openssh/* \
              /tmp/local/etc/s6/.s6-svscan/* \
              /go/src/code.gitea.io/gitea/gitea

FROM docker.io/library/alpine:3.23 AS gitea

EXPOSE 22 3000

RUN apk --no-cache add \
    bash \
    ca-certificates \
    curl \
    gettext \
    git \
    linux-pam \
    openssh \
    s6 \
    sqlite \
    su-exec \
    gnupg

RUN addgroup \
    -S -g 1000 \
    git && \
  adduser \
    -S -H -D \
    -h /data/git \
    -s /bin/bash \
    -u 1000 \
    -G git \
    git && \
  echo "git:*" | chpasswd -e

COPY --from=build-env /tmp/local /
COPY --from=build-env /go/src/code.gitea.io/gitea/gitea /app/gitea/gitea

ENV USER=git
ENV GITEA_CUSTOM=/data/gitea

VOLUME ["/data"]

# HINT: HEALTH-CHECK-ENDPOINT: don't use HEALTHCHECK, search this hint keyword for more information
ENTRYPOINT ["/usr/bin/entrypoint"]
CMD ["/usr/bin/s6-svscan", "/etc/s6"]
