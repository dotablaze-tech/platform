FROM postgres:17.11@sha256:b8faeff7712663c2773a190181e1246e37e364b887611c04658c28b44e2408fe

# POSTGRES_PASSWORD is deliberately not baked in: supply it at run time.
ENV POSTGRES_USER=default_user
ENV POSTGRES_DB=meowbot

COPY --chown=postgres:postgres apps/database/meowbotdb/src/ /docker-entrypoint-initdb.d/

# The official entrypoint supports starting as the postgres user directly.
USER postgres

HEALTHCHECK --interval=10s --timeout=5s --start-period=30s --retries=5 \
  CMD pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" || exit 1
