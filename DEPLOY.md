# AeGit Stage 1 — Deployment (finish from Beast)

Image built and pushed: `ghcr.io/agyeman-enterprises/aegit:stage1`

All code is on `main`. Run from Beast via Hetzner web console (amiacoda 5.9.153.215).

## Step 1 — Env file

Get a fresh GitHub token first: run `gh auth token` on aa-bkdc2.

```bash
mkdir -p /opt/aegit-data

cat > /opt/aegit-data/.env << 'ENVEOF'
AEGIT_ADMIN_USER=agyeman
AEGIT_ADMIN_PASSWORD=Gitea2026!xK9
AEGIT_ADMIN_EMAIL=isaalia@gmail.com
GITHUB_TOKEN=PASTE_TOKEN_HERE
SECRET_KEY=
INTERNAL_TOKEN=
ENVEOF
```

## Step 2 — Login to GHCR and deploy

```bash
echo "PASTE_TOKEN_HERE" | podman login ghcr.io -u isaalia --password-stdin

source /opt/aegit-data/.env
podman stop aegit 2>/dev/null || true
podman rm   aegit 2>/dev/null || true

podman run -d \
  --name aegit \
  --restart unless-stopped \
  -p 8001:8001 \
  -p 2222:22 \
  -v /opt/aegit-data:/data \
  -e AEGIT_ADMIN_USER \
  -e AEGIT_ADMIN_PASSWORD \
  -e AEGIT_ADMIN_EMAIL \
  -e GITHUB_TOKEN \
  -e SECRET_KEY \
  -e INTERNAL_TOKEN \
  ghcr.io/agyeman-enterprises/aegit:stage1

sleep 10 && podman logs aegit --tail 50
```

## Step 3 — Cloudflare tunnel

Copy `deploy/cloudflare-tunnel.yaml` to amiacoda and activate:
- `aegit.agyemanenterprises.com` → port 8001
- `ssh.aegit.agyemanenterprises.com` → TCP port 2222

## Step 4 — Smoke tests

Run all 10 tests from `PLAN.md` Step 8. All must pass before OO Stage 1 completion review.

## Notes

- SSH port 22 blocked on Hetzner firewall — use Hetzner web console only
- Admin creds match existing Gitea instance on same server
