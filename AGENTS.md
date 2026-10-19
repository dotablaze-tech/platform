# AGENTS.md

Canonical instructions for AI agents in this repo. `CLAUDE.md` only imports
this file; edit this one. Structure and Nx usage: [README.md](README.md);
scripts: [tools/README.md](tools/README.md).

## Repo facts

- Nx 20 + Go workspace (`go.work`). Projects: `meowbot` (Discord bot, the only
  deployed app), `meowbotdb` (SQL), and `go-meowbot-*` libs. `npx nx show projects`
  lists them.
- npm only (`package-lock.json`); install with `npm ci`.

## Commands

```bash
npm ci
npx nx format:check                       # CI gate; fix with npx nx format:write
npx nx affected -t lint test build        # what CI runs per PR
npx nx run-many -t lint test build -p meowbot
npx nx reset                              # after editing project.json targets
```

Run `npx nx format:write` before committing (CI checks Markdown too).

## Release flow

Push to `main` runs `.github/workflows/ci.yml`: `check` (format, lint, test,
build) then, on `main` only, `release`:

1. `@jscutlery/semver` versions affected projects. Only `fix` and `feat`
   commits bump; `chore`, `docs`, `ci` etc. release nothing.
2. It pushes the `meowbot-X.Y.Z` tag, then runs `build-image` (pushes
   `dotablaze/meowbot:X.Y.Z` and `latest`).
3. `update-app` runs `tools/update-app.sh`, which pushes the new `appVersion`
   to `charts/meowbot/Chart.yaml` on `dotablaze-tech/deployments` `main`.
4. ArgoCD rolls the new version to non-prod automatically. Prod is pinned in
   `deployments` `charts/meowbot/values-prd.yaml`; promote by PR.

Never run `nx affected -t version` locally: tags and pushes are irreversible.
Renovate bumps for `gomod` and `dockerfile` use `fix(deps)` so they release.

## Branches, commits, PRs

- Never work on or push to `main`; use a worktree. Branch
  `type/DOTA-N-kebab-case`, e.g. `docs/DOTA-10-add-agents-md`.
- Conventional Commits, enforced by commitlint in the `commit-msg` hook. No
  `--no-verify`. Atomic commits; rebase before pushing, never merge commits.
- Merges are rebase-merge only. Renovate auto-merges minor/patch updates.
- PR title under 70 chars; body about 150 words: Why / Needs attention /
  Risk / Verified. Wait for green CI.
- AI attribution goes only in commit trailers (`Co-Authored-By`,
  `Assisted-by`), never in files or PR text. Org rules:
  [dotablaze-tech/.github CONTRIBUTING.md](https://github.com/dotablaze-tech/.github/blob/main/CONTRIBUTING.md).

## Secrets

`DOCKERHUB_*` and deploy-key secrets exist only for the `release` job. Never
hardcode or print a credential; the bot reads `DISCORD_BOT_TOKEN` from the
environment.
