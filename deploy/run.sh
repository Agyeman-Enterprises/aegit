#!/bin/sh
set -eu
# AeGit Stage 1 deploy script for amiacoda
# Secrets sourced from /opt/aegit-data/.env before calling this script.
IMAGE="ghcr.io/agyeman-enterprises/aegit:stage1"
CONTAINER="aegit"
DATA_DIR="/opt/aegit-data"

docker pull "$IMAGE"
docker stop "$CONTAINER" 2>/dev/null || true
docker rm   "$CONTAINER" 2>/dev/null || true

docker run -d \
  --name "$CONTAINER" \
  --restart unless-stopped \
  -p 8001:8001 \
  -p 2222:22 \
  -v "$DATA_DIR:/data" \
  -e AEGIT_ADMIN_USER \
  -e AEGIT_ADMIN_PASSWORD \
  -e AEGIT_ADMIN_EMAIL \
  -e GITHUB_TOKEN \
  -e SECRET_KEY \
  -e INTERNAL_TOKEN \
  "$IMAGE"

echo "AeGit running. Web: http://localhost:8001 | SSH: localhost:2222"
