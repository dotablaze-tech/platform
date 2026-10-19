# Dotablaze Platform Monorepo

[![CI](https://github.com/dotablaze-tech/platform/actions/workflows/ci.yml/badge.svg)](https://github.com/dotablaze-tech/platform/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/Go-1.26-blue)
![Nx](https://img.shields.io/badge/Nx-monorepo-blue)
![MIT License](https://img.shields.io/badge/License-MIT-yellow.svg)

The Dotablaze Tech Platform Monorepo is a multi-language, multi-project repository that houses all application code,
reusable libraries, configuration, and tooling used across the Dotablaze ecosystem. This repository is organized to
support scalable development with a clean separation of concerns and structured sharing of features, utilities, and
infrastructure.

## 📁 Structure Overview

- **apps/** – Deployable applications, grouped by language or runtime.
- **libs/** – Go libraries used by `meowbot`, under `libs/go/meowbot/` (`feature/api`, `feature/db`,
  `feature/handler`, `feature/state`, `util`).
- **tools/** – Automation scripts, such as `update-app.sh` (see [tools/README.md](tools/README.md)).
- **packages/** – npm workspace packages (currently empty).

```
.
├── apps/
│   ├── database/
│   │   └── meowbotdb/       # Postgres image with the meowbot schema
│   └── go/
│       └── meowbot/         # Discord bot
├── libs/
│   └── go/
│       └── meowbot/
│           ├── feature/     # api, db, handler, state
│           └── util/
└── tools/                   # update-app.sh
```

Go modules are tied together by `go.work`; each project has its own `go.mod` and Nx `project.json`.

## 🚀 Nx Task Execution

Nx provides efficient task orchestration for builds, tests, and automation.

```bash
npx nx <target> <project> [options]
```

**Examples:**

```bash
npx nx build meowbot
npx nx test go-meowbot-feature-handler
npx nx run-many -t build -p meowbot meowbotdb
```

See the full guide: [https://nx.dev/features/run-tasks](https://nx.dev/features/run-tasks)

## 🌐 Visualize Dependencies

Generate the project graph with:

```bash
npx nx graph
```

This helps you visualize app/library dependencies and identify opportunities for optimization or sharing.

## 📦 Deployment & Infrastructure

Deployment and Kubernetes infrastructure for Dotablaze services are managed in separate repositories:

- **Deployment Configs**: [https://github.com/dotablaze-tech/deployments](https://github.com/dotablaze-tech/deployments)
- **Infra/Kubernetes (WIP)**: [https://github.com/jdwillmsen/jdw-kube](https://github.com/jdwillmsen/jdw-kube)

Container images are hosted at:

- **Docker Hub**: [https://hub.docker.com/u/dotablaze](https://hub.docker.com/u/dotablaze)

## 🧪 CI/CD & Automation

CI pipelines live in `.github/workflows/`, and use Nx to build affected projects only, improving feedback time.

`ci.yml` has two jobs. `check` runs format check, lint, test, component-test, e2e and build (affected projects only) on
every PR and push, with read-only permissions. `release` runs only on pushes to `main` after `check` passes and is the
only job that receives the Docker Hub and deploy-key secrets.

### Release Flow

1. A push to `main` runs the CI `release` job, which builds affected projects and runs `nx affected -t version`.
2. `@jscutlery/semver` (conventional commits) bumps only for `fix` and `feat`; `docs`, `chore`, `ci` and similar
   commits do not release (`skipCommitTypes` in `nx.json`). The `go-meowbot-*` libs are versioned too, so a lib fix
   tags the lib and patches `meowbot`.
3. A bump creates a release commit and a tag such as `meowbot-X.Y.Z`.
4. The `postTargets` run: `build-image` pushes `dotablaze/meowbot:X.Y.Z` and `:latest` to Docker Hub for linux/amd64 and
   linux/arm64.
5. `update-app` runs [`tools/update-app.sh`](tools/README.md), which pushes the new `appVersion` to
   `charts/meowbot/Chart.yaml` on `main` of the deployments repo. ArgoCD auto-rolls `non`.
6. `prd` pins `image.tag` in the deployments repo's `charts/meowbot/values-prd.yaml` and is promoted by a PR there.

`meowbotdb` versions the same way (image only, no `update-app`).

Renovate auto-merges minor and patch updates once CI passes. Go module and Dockerfile updates are committed as
`fix(deps)` so they release; a merge therefore deploys a new image to `non`, never to `prd`.

## 📌 About

This repo is the central platform for building, deploying, and scaling Dotablaze Tech applications — including bots,
services, and eventually frontend dashboards. Managed with Nx and Go Workspaces for flexibility and performance.

### 👤 Maintainer

- **Jake Willmsen**
