# Meowbot

![Build](https://img.shields.io/github/actions/workflow/status/dotablaze-tech/platform/ci.yml?branch=main)
![Docker Image Version](https://img.shields.io/docker/v/dotablaze/meowbot)
![Docker Image Size](https://img.shields.io/docker/image-size/dotablaze/meowbot)
![Docker Downloads](https://img.shields.io/docker/pulls/dotablaze/meowbot?label=downloads)
![Nx](https://img.shields.io/badge/Nx-managed-blue)

Meowbot is a Discord bot built with Go and `discordgo`. It tracks consecutive "meow" messages in a channel per guild,
keeps streaks and high scores in PostgreSQL, and exposes slash commands and a small HTTP stats API.

## Project Structure

```
apps/go/meowbot/
├── Dockerfile          # Production image (copies a prebuilt binary from dist/)
├── Dockerfile.local    # Local image that builds the binary inside Docker
├── go.mod / go.sum     # Module definition
├── main.go             # Entrypoint
├── main_test.go        # Unit tests
├── project.json        # Nx targets
└── assets/             # Bot icons
```

Handlers, database access, the HTTP API and config live in `libs/go/meowbot/` (`feature/handler`, `feature/db`,
`feature/api`, `feature/state`, `util`). The schema is in [`apps/database/meowbotdb`](../../database/meowbotdb).

## Getting Started

Prerequisites: Go (version in `go.mod`), Node.js with `npm ci` run at the repo root, a Discord bot token, a reachable
PostgreSQL database (Docker optional).

```bash
# Start a local database on :5432 (see apps/database/meowbotdb)
npx nx run meowbotdb:serve

# Run the bot
DISCORD_BOT_TOKEN=your-token npx nx run meowbot:serve
```

With `DATABASE_URL` unset, the bot connects to
`postgres://default_user:default_password@127.0.0.1:5432/meowbot?sslmode=disable`.

## Configuration

| Variable             | Default      | Purpose                                                                        |
| -------------------- | ------------ | ------------------------------------------------------------------------------ |
| `DISCORD_BOT_TOKEN`  | none         | Required. The bot exits at startup without it.                                 |
| `DATABASE_URL`       | unset        | See the note below.                                                            |
| `DATABASE_USER`      | none         | Postgres user, used when `DATABASE_URL` is set.                                |
| `DATABASE_PASSWORD`  | none         | Postgres password, used when `DATABASE_URL` is set.                            |
| `DATABASE_HOST`      | none         | Postgres host, used when `DATABASE_URL` is set.                                |
| `DATABASE_PORT`      | none         | Postgres port, used when `DATABASE_URL` is set.                                |
| `DATABASE_NAME`      | none         | Postgres database name, used when `DATABASE_URL` is set.                       |
| `MODE`               | `production` | `production` allows every guild; any other value enforces the guild whitelist. |
| `WHITELISTED_GUILDS` | empty        | Comma-separated guild IDs allowed when `MODE` is not `production`.             |
| `DEBUG`              | `false`      | `true` enables debug logging.                                                  |
| `API_PORT`           | `8080`       | Port for the HTTP API.                                                         |
| `EMOJI_LIST`         | none         | Emoji list read by the emoji helper.                                           |

`DATABASE_URL` is only tested for emptiness: when non-empty, the connection string is built as
`postgres://USER:PASSWORD@HOST:PORT/NAME` from the five `DATABASE_*` variables and the URL's own value is ignored
(`sslmode` is not set). The Helm chart sets it to a placeholder JDBC URL for this reason.

## Commands

| Command        | Description                                                               |
| -------------- | ------------------------------------------------------------------------- |
| `/count`       | Current meow count for this server                                        |
| `/highscore`   | Highest meow streak for this server                                       |
| `/stats`       | Your personal meow stats; optional `scope` (`guild` or `global`)          |
| `/setup`       | Choose the channel Meowbot listens in (`channel` option)                  |
| `/leaderboard` | Top meowers; optional `scope`, `metric` (`total` or `success`) and `page` |

A message matching `(?i)^m+e+o+w+$` (for example `meooow`) counts as a meow. The same user cannot meow twice in a row.

## HTTP API

Served on `API_PORT`: `/liveness`, `/readiness`, `/stats`, `/leaderboard`, `/users`, `/users/...`, `/guilds` and
`/guilds/...`. The chart's probes use `/liveness` and `/readiness`.

## Nx Targets

```bash
npx nx build meowbot
npx nx test meowbot
npx nx lint meowbot
npx nx run meowbot:local-build-image   # builds dotablaze/meowbot:latest locally
npx nx run meowbot:serve-container     # runs the image; forwards DISCORD_BOT_TOKEN from the environment
```

## Release and Deploy

Releases are automated; do not run `version` or `build-image` by hand (`build-image` pushes to Docker Hub).

1. A `fix:` or `feat:` commit merged to `main` makes the CI `release` job run `nx affected -t version`.
2. `@jscutlery/semver` bumps the version, tags `meowbot-X.Y.Z` and pushes `dotablaze/meowbot:X.Y.Z` and `:latest`
   (linux/amd64 and linux/arm64) to Docker Hub.
3. The `update-app` target sets `appVersion` in `charts/meowbot/Chart.yaml` in
   [dotablaze-tech/deployments](https://github.com/dotablaze-tech/deployments), and ArgoCD rolls the new version to
   `non`.
4. `prd` pins `image.tag` in `charts/meowbot/values-prd.yaml` and is promoted by a PR in the deployments repo.

Other commit types (`docs`, `chore`, `ci`, ...) do not release. See the [root README](../../../README.md#release-flow).

## License

MIT, as declared in the root `package.json`. The repository has no `LICENSE` file.
