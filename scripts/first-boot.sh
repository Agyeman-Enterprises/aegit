#!/bin/sh
set -eu
FLAG="/data/.first-boot-done"
if [ -f "$FLAG" ]; then echo "[first-boot] Already done."; exit 0; fi
echo "[first-boot] Creating admin: ${AEGIT_ADMIN_USER}"
aegit admin user create \
  --username "$AEGIT_ADMIN_USER" \
  --password "$AEGIT_ADMIN_PASSWORD" \
  --email "$AEGIT_ADMIN_EMAIL" \
  --admin \
  --must-change-password=false \
  --config /data/gitea/conf/app.ini
touch "$FLAG"
echo "[first-boot] Done."
