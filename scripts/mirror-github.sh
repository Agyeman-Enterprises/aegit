#!/bin/sh
set -eu
FLAG="/data/.mirrored"
if [ -f "$FLAG" ]; then echo "[mirror] Already done."; exit 0; fi

AUTH="$(printf '%s:%s' "$AEGIT_ADMIN_USER" "$AEGIT_ADMIN_PASSWORD" | base64 | tr -d '\n')"
API="http://localhost:8085/api/v1"

for ORG in Agyeman-Enterprises isaalia imho-media; do
  echo "[mirror] Ensuring org: $ORG"
  curl -sf -X POST "$API/orgs" \
    -H "Authorization: Basic $AUTH" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$ORG\",\"visibility\":\"private\"}" > /dev/null 2>&1 || true

  PAGE=1
  while true; do
    REPOS=$(curl -sf \
      -H "Authorization: Bearer $GITHUB_TOKEN" \
      "https://api.github.com/orgs/$ORG/repos?per_page=100&type=all&page=$PAGE" \
      | jq -r '.[].name' 2>/dev/null)
    [ -z "$REPOS" ] && break
    for REPO in $REPOS; do
      echo "[mirror] $ORG/$REPO"
      curl -sf -X POST "$API/repos/migrate" \
        -H "Authorization: Basic $AUTH" \
        -H "Content-Type: application/json" \
        -d "{\"clone_addr\":\"https://github.com/$ORG/$REPO\",\"auth_token\":\"$GITHUB_TOKEN\",\"mirror\":true,\"private\":true,\"repo_name\":\"$REPO\",\"repo_owner\":\"$ORG\",\"service\":\"github\"}" \
        > /dev/null 2>&1 \
        && echo "[mirror] OK: $ORG/$REPO" \
        || echo "[mirror] SKIP (exists?): $ORG/$REPO"
    done
    PAGE=$((PAGE+1))
  done
done

touch "$FLAG"
echo "[mirror] Complete."
