FROM postgres:17.11@sha256:2d2b8998d31037bf721cfdf764d76ba74171b4fab3431b7f72c27c56ddbdf9e3

# POSTGRES_PASSWORD is deliberately not baked in: supply it at run time.
ENV POSTGRES_USER=default_user
ENV POSTGRES_DB=meowbot

COPY --chown=postgres:postgres apps/database/meowbotdb/src/ /docker-entrypoint-initdb.d/

# The official entrypoint supports starting as the postgres user directly.
USER postgres

HEALTHCHECK --interval=10s --timeout=5s --start-period=30s --retries=5 \
  CMD pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" || exit 1
