#!/bin/sh
set -eu

# Generate secrets if not pre-set (use /dev/urandom — openssl absent in Alpine)
SECRET_KEY="${SECRET_KEY:-$(head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n')}"
INTERNAL_TOKEN="${INTERNAL_TOKEN:-$(head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n')}"
export SECRET_KEY INTERNAL_TOKEN

# Write final config from template
mkdir -p /data/gitea/conf
envsubst < /app/app.ini.template > /data/gitea/conf/app.ini

# Create admin user (CLI, pre-server � writes direct to SQLite)
aegit migrate --config /data/gitea/conf/app.ini
/usr/local/bin/first-boot.sh

# Start AeGit in background
aegit web --config /data/gitea/conf/app.ini &
GITEA_PID=$!
trap 'kill $GITEA_PID 2>/dev/null' INT TERM

# Wait for server ready (max 60s)
i=0
while [ $i -lt 60 ]; do
  curl -sf "http://localhost:8085/" > /dev/null 2>&1 && echo "[entrypoint] Ready after ${i}s." && break || true
  sleep 1
  i=$((i+1))
done

# Mirror all GitHub repos now that API is live
/usr/local/bin/mirror-github.sh

# Foreground
wait $GITEA_PID
