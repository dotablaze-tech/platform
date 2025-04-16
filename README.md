# Dotablaze Platform Monorepo

[![CI](https://github.com/dotablaze-tech/platform/actions/workflows/ci.yml/badge.svg)](https://github.com/dotablaze-tech/platform/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/Go-1.24-blue)
![Nx](https://img.shields.io/badge/Nx-monorepo-blue)
![MIT License](https://img.shields.io/badge/License-MIT-yellow.svg)

The Dotablaze Tech Platform Monorepo is a multi-language, multi-project repository that houses all application code,
reusable libraries, configuration, and tooling used across the Dotablaze ecosystem. This repository is organized to
support scalable development with a clean separation of concerns and structured sharing of features, utilities, and
infrastructure.

## 📁 Structure Overview

- **apps/** – Full-stack applications grouped by language or runtime (e.g., Go bots and services).
- **libs/** – Reusable code libraries organized by app, language, or shared domain:
    - `feature/` – Domain-specific handlers, components, or state.
    - `data-access/` – Integration with APIs or external services.
    - `util/` – Common utility functions, models, and helpers.
    - `ui/` – Visual or bot-rendered components (when applicable).
- **tools/** – Dev tooling for formatting, local dev, automation, and updates.

## 🗂️ Directory Structure

```
.
├── apps/                    # Complete applications
│   └── go/                  # Grouped by framework / language
│       ├── meowbot/         # Specific application
│       └── barkbot/
│
├── libs/                    # Reusable libraries
│   ├── go/                  # Grouped by framework
│   │   ├── meowbot/         # App-specific libraries
│   │   │   ├── feature/
│   │   │   └── util/
│   │   ├── shared/          # Framework-wide shared
│   │   │   ├── ui/
│   │   │   ├── util/
│   │   │   └── data-access/
│   └── shared/              # Cross-framework shared
│       └── utils/
│
└── tools/                   # Monorepo tooling
```

### 🔑 Key Architectural Principles

- **Language-Scoped**: Projects and libraries are grouped by runtime (e.g., `go/`, `angular/`, etc.).
- **App Isolation**: Features unique to a specific app are namespaced under that app.
- **Shared Logic**:
    - App-specific: Scoped to one app only (`libs/go/meowbot/*`).
    - Framework-wide: Usable across projects using the same runtime.
    - Cross-platform: Reusable by any project (e.g., `libs/shared`).
- **Library Types**:
    - `feature/` – Domain-level logic or bot commands.
    - `data-access/` – Service interfaces and clients.
    - `util/` – Reusable helpers, emoji maps, etc.
    - `ui/` – Discord message formatting or future rich UIs.

## 🚀 Nx Task Execution

Nx provides efficient task orchestration for builds, tests, and automation.

```bash
npx nx <target> <project> [options]
```

**Examples:**

```bash
npx nx build go-meowbot
npx nx test go-meowbot-feature-handler
npx nx run-many -t build -p meowbot
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

## 📚 Library Strategy

This monorepo encourages reuse and clean separation of concerns by organizing libraries in layers:

- `feature/` – Commands, event handlers, and scoped state.
- `data-access/` – Service and bot integrations (coming soon).
- `util/` – Emoji helpers, string formatters, etc.
- `ui/` – For any visual rendering components.

## 🧪 CI/CD & Automation

CI pipelines live in `.github/workflows/`, and use Nx to build affected projects only, improving feedback time.

## 📌 About

This repo is the central platform for building, deploying, and scaling Dotablaze Tech applications — including bots,
services, and eventually frontend dashboards. Managed with Nx and Go Workspaces for flexibility and performance.

### 👤 Maintainer

- **Jake Willmsen**