# AeGit Stage 1 — sovereign forge image
# Based on official Gitea nightly (matches v1.27.0-dev fork).
# Stage 2 will build from source once AEGIS hooks are in place.

FROM gitea/gitea:latest

USER root

# aegit binary alias + runtime tools
RUN apk add --no-cache gettext jq && \
    ln -sf /usr/local/bin/gitea /usr/local/bin/aegit

# Config template — secrets filled at container start via envsubst
COPY custom/conf/app.ini.template /app/app.ini.template

# Container runtime scripts
COPY scripts/first-boot.sh    /usr/local/bin/first-boot.sh
COPY scripts/mirror-github.sh /usr/local/bin/mirror-github.sh
COPY entrypoint.sh            /entrypoint.sh

RUN chmod +x /usr/local/bin/first-boot.sh \
             /usr/local/bin/mirror-github.sh \
             /entrypoint.sh

VOLUME ["/data"]
EXPOSE 8001 22

ENTRYPOINT ["/entrypoint.sh"]
