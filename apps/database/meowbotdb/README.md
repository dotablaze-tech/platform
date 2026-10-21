# MeowbotDB

![Build](https://img.shields.io/github/actions/workflow/status/dotablaze-tech/apps/ci.yml?branch=main)
![Docker Image Version](https://img.shields.io/docker/v/dotablaze/meowbotdb)
![Docker Image Size](https://img.shields.io/docker/image-size/dotablaze/meowbotdb)
![Docker Downloads](https://img.shields.io/docker/pulls/dotablaze/meowbotdb?label=downloads)
![Nx](https://img.shields.io/badge/Nx-managed-blue)

**MeowbotDB** is the PostgreSQL database service used by Meowbot — a multi-server Discord bot for community interaction, reactions, and playful engagement. It provides the schema and seed data required for bot state, logs, preferences, and other persistent features.

This service is containerized and managed via [Nx](https://nx.dev), supporting multi-arch Docker builds and local development with optional volume caching.

---

## 📁 Project Structure

```text
database/meowbotdb/
├── Dockerfile           # Defines the Postgres image and entrypoint
├── project.json         # Nx configuration and targets
└── src/
    ├── 00_schema.sql    # SQL schema definitions for bot state and logs
    └── 01_data.sql      # Initial data (e.g., test configs, emoji presets)
```

---

## 🚀 Getting Started

### Prerequisites

- Docker
- [Nx CLI](https://nx.dev)

### Run Locally (Ephemeral)

```bash
nx run meowbotdb:serve
```

Runs PostgreSQL in a disposable Docker container on port `5432`.

### Run with Persistent Volume

```bash
nx run meowbotdb:serve-cache
```

Creates and mounts a Docker volume (`meowbotdb-data`) for persistent local storage.

### Clear Persistent Volume

```bash
nx run meowbotdb:clear-cache
```

Removes the `meowbotdb-data` volume to reset local database state.

---

## 🔨 Build

Build Docker images with semantic versioning and multi-platform support.

### CI/CD Multi-Arch Build

```bash
nx run meowbotdb:build-image
```

Builds and pushes images for `linux/amd64` and `linux/arm64`, tagged as `latest` and with semantic version.

### Local Build Only

```bash
nx run meowbotdb:local-build-image
```

Builds a local image using host architecture for development/testing.

---

## 🗃️ Database Credentials

Default credentials after container start:

```text
Host:     localhost
Port:     5432
Database: meowbot
User:     default_user
Password: default_password
```

---

## 📄 SQL Files

- `00_schema.sql` – Table definitions, indexes, relationships
- `01_data.sql` – Seed data for dev/testing (emojis, preferences, etc.)

These scripts are executed automatically at container start.

---

## 🧪 Versioning

Handled via [Conventional Commits](https://www.conventionalcommits.org/) + [`@jscutlery/semver`](https://github.com/jscutlery/semver):

```bash
nx run meowbotdb:version
```

Automates version bumping, changelog, tagging, and build/push.

---

## 📦 Deployment

Used internally by Meowbot services across environments, including CI pipelines and local dev.

---

## 📌 Notes

- This is a purpose-built, stateful service for Discord bot functionality.
- Migration tooling is not currently integrated — for production schema evolution, consider [Flyway](https://flywaydb.org/) or [Sqitch](https://sqitch.org/).
