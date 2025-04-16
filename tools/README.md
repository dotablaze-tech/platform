# 🛠️ Dotablaze Tech Monorepo Tools

This directory contains helper scripts and tooling used across the Dotablaze Tech platform monorepo. These tools assist
with version management, refactoring, and automation of common developer workflows in a multi-language, multi-project
environment.

## 📂 Directory Contents

```
tools/
├── README.md           # This file
└── update-app.sh       # Git-based app version bump script for deployments
```

## 📜 Scripts

### `update-app.sh`

A Git automation script to bump the `appVersion` field in Kubernetes manifests or Helm chart files inside
the [dotablaze-tech/deployments](https://github.com/dotablaze-tech/deployments) repository.

This tool clones the remote deployments repo, updates the version for a specific app, commits the change with a standard
message, and pushes it to the remote.

> ✅ **Use this script to consistently version your apps for CI/CD deployment.**

#### 📦 What it does:

- Clones the deployments Git repository to a temporary folder
- Validates the target file is tracked and clean
- Updates the `appVersion:` field with the new version
- Commits and pushes the change with a conventional commit message

#### ✅ Usage

```bash
./tools/update-app.sh <file_path> <version_number> <project_name>
```

**Example:**

```bash
./tools/update-app.sh charts/meowbot/Chart.yaml 1.3.2 meowbot
```

This command updates the version in `Chart.yaml` to `1.3.2` for the `meowbot` project and pushes the change to GitHub.

> ⚠️ **Warning:** This script will attempt to push changes. Make sure your Git credentials and SSH keys are set up
> properly.

---

## 📌 Future Tools (Ideas)

This directory is designed to grow with the platform. Planned or proposed tooling includes:

- 🔄 Changelog automation and conventional commits enforcement
- 🧪 Test coverage reporters
- 🐳 Docker helpers for local multi-app development
- 🧹 Nx/Go workspace cleanup and refactor tools
- 📊 Dependency analysis and graph auditing

## 👀 See Also

- Platform Monorepo: [`../README.md`](../README.md)
- Nx CLI Reference: [https://nx.dev/cli](https://nx.dev/cli)
- Deployment Repository: [dotablaze-tech/deployments](https://github.com/dotablaze-tech/deployments)

---

Maintained by **@jdwillmsen** – Built to support all internal Dotablaze platform applications.
