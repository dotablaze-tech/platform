# 🐾 Meow Bot

![Build](https://img.shields.io/github/actions/workflow/status/dotablaze-tech/platform/ci.yml?branch=main)
![Docker Image Version](https://img.shields.io/docker/v/dotablaze/meowbot)
![Docker Image Size](https://img.shields.io/docker/image-size/dotablaze/meowbot)
![Docker Downloads](https://img.shields.io/docker/pulls/dotablaze/meowbot?label=downloads)
![Nx](https://img.shields.io/badge/Nx-managed-blue)

**Meow Bot** is a lightweight and fun Discord bot built with Go and powered by `discordgo`. It tracks consecutive “meow”
messages in a single channel, maintaining streaks, preventing duplicate users, and celebrating high scores. Ideal for
community engagement, cat lovers, and general chaos.

This project is managed with [Nx](https://nx.dev) and includes Docker support for local and CI/CD workflows.

---

## 📁 Project Structure

```
apps/go/meow-bot/
├── Dockerfile              # Production image
├── Dockerfile.local        # Local development Dockerfile
├── README.md               # This file
├── go.mod / go.sum         # Module definitions
├── main.go                 # Main bot entrypoint
├── main_test.go            # Unit tests
└── assets/                 # Bot avatars/icons
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- [Nx CLI](https://nx.dev)
- Docker (optional for containerized runs)
- Discord bot token

### Local Run (Token via ENV)

```bash
DISCORD_TOKEN=your-token-here nx run meow-bot:serve
```

### Local Dev with Volume (Docker)

```bash
nx run meow-bot:serve-cache
```

Creates a volume (`meow-bot-data`) for persistent local development (e.g., future DB use).

---

## ✨ Features

- Regex-based “meow” detection (e.g., `meooow`, `MEEEEOW`)
- Streak counter per guild
- Prevents same user from meowing twice in a row
- Tracks and announces high scores
- Slash command `/highscore` to show current record
- Guild-specific in-memory state tracking
- `slog`-based structured logging

---

## 🐳 Docker

### Build Production Image

```bash
npx nx run meow-bot:build-image
```

### Build Local Dev Image

```bash
npx nx run meow-bot:local-build-image
```

---

## 📦 Deployment

This bot is designed to run as a stateless container with ephemeral memory or backed by a persistent volume. Future
plans may include Redis or embedded DB support.

---

## 🔧 Commands

- `/highscore` – Shows the current top meow streak and who set it.

---

## 🐱 Example Interaction

```text
User1: meow
Bot: 😺 meow x1!
User2: meeeow
Bot: 😼 meow x2!
User1: meow
Bot: 😻 meow x3!
User1: meow
Bot: 😾 You can't meow twice in a row!
```

---

## 🧪 Testing

```bash
go test ./apps/go/meow-bot
```

---

## 📌 Notes

- You must set the bot token via `DISCORD_TOKEN` env var or update `main.go` to read from config.
- The bot is intended for one channel per guild.
- You can customize behavior (e.g., emojis, reset behavior) in the handler and state packages.

---

## 📷 Assets

A few example icons are included in `apps/go/meow-bot/assets/` for use when registering your bot on Discord.

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
