# 🗣️ Meowbot Handler

![Build](https://img.shields.io/github/actions/workflow/status/dotablaze-tech/apps/ci.yml?branch=main)
![Nx](https://img.shields.io/badge/Nx-managed-blue)
![Go Module](https://img.shields.io/badge/Go-Module-brightgreen)

**Meowbot Handler** is a Go package that powers the Discord event and command handling layer of [🐾 Meow Bot](https://github.com/dotablaze-tech/apps/tree/main/apps/go/meowbot). It defines logic for responding to message events, slash commands, and interaction components using the `discordgo` library.

This package is organized under the monorepo’s feature set and managed using [Nx](https://nx.dev).

---

## 📁 Project Structure

```
libs/go/meowbot/feature/handler/
├── commands.go        # Slash command handling logic
├── messages.go        # Regex-based message response logic
├── messages_test.go   # Unit tests for message handling
├── go.mod / go.sum    # Go module definition
└── project.json       # Nx project definition
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- [Nx CLI](https://nx.dev)

### Installation

Used internally by Meowbot’s core app. Import like this:

```go
import "github.com/dotablaze-tech/apps/libs/go/meowbot/feature/handler"
```

---

## ✨ Features

- 🔍 Regex-based detection for “meow” messages (`meooow`, `meeeeow`, etc.)
- 🧩 Handles Discord slash commands such as `/highscore` and `/leaderboard`
- 📬 Responds to Discord message events with embedded state logic
- 📜 Modular message and command routing
- 🧪 Unit tested for behavior correctness

---

## 🧠 Example Usage

```go
handler := handler.NewHandler(deps)

session.AddHandler(handler.HandleMessageCreate)
session.AddHandler(handler.HandleInteractionCreate)
```

---

## 🧪 Testing

Run tests for message event logic:

```bash
go test ./libs/go/meowbot/feature/handler
```

---

## 🔌 Integration Notes

- The handler relies on `state` and `db` libraries for tracking and persistence.
- Slash commands should be registered using `commands.go` definitions during startup.
- Structured logging via `slog` is embedded throughout.

---

## 📌 Notes

- Not a standalone Discord bot. Requires orchestration via the main Meowbot app.
- Commands and interactions are modular—easily extended via additional handlers.
