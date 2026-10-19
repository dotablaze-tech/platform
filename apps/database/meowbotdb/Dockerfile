FROM postgres:17.11@sha256:2d2b8998d31037bf721cfdf764d76ba74171b4fab3431b7f72c27c56ddbdf9e3

ENV POSTGRES_USER=default_user
ENV POSTGRES_PASSWORD=default_password
ENV POSTGRES_DB=meowbot

COPY apps/database/meowbotdb/src/ /docker-entrypoint-initdb.d/
