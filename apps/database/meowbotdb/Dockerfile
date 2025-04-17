FROM postgres:17.4

ENV POSTGRES_USER=default_user
ENV POSTGRES_PASSWORD=default_password
ENV POSTGRES_DB=meowbot

COPY apps/database/meowbotdb/src/ /docker-entrypoint-initdb.d/
