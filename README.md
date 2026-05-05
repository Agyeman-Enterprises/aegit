# AeGit — The Sovereign Forge

[![](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT "License: MIT")

AeGit is a self-hosted Git service built by Agyeman Enterprises. It is a sovereign fork of Gitea, extended with AE-specific branding and operational improvements, and deployed on the AE infrastructure stack (amiacoda / Srvrsup / Podman).

## What it is

- Full-featured Git forge: repositories, issues, pull requests, milestones, organizations
- SQLite-backed for low-maintenance deployment on sovereign hardware
- SSH + HTTP access on non-privileged ports (2222 / 8085)
- GitHub mirror sync via API on startup
- Packaged as a multi-stage Podman OCI image — no Docker, no Coolify

## What is CODEFLAG'd (phase3-remove)

The following subsystems are present in the source but marked for future removal:

- **Actions** — CI/CD runner (services/actions/) — AE uses external CI
- **Packages** — Package registry (services/packages/) — AE uses AEnIO/MinIO
- **Projects** — Kanban boards (models/project/) — AE uses dedicated planning tools

These subsystems are disabled in app.ini (`ENABLED = false`) but not yet removed from source to avoid touching too many surfaces in one sprint.

## Deployment

AeGit runs on amiacoda via Srvrsup. Docker is banned. Podman only.

```
podman build -f Containerfile -t ghcr.io/agyeman-enterprises/aegit:stage2 .
```

Live at: https://aegit.agyemanenterprises.com

## Upstream credit

AeGit is a fork of [Gitea](https://gitea.com) (MIT License). The Gitea project is maintained by the Gitea organization. This fork is not affiliated with or endorsed by the Gitea project.
