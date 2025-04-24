# 🐘 Meowbot DB

![Build](https://img.shields.io/github/actions/workflow/status/dotablaze-tech/platform/ci.yml?branch=main)
![Nx](https://img.shields.io/badge/Nx-managed-blue)
![Go Module](https://img.shields.io/badge/Go-Module-brightgreen)

**Meowbot DB** is the database persistence layer for [🐾 Meow Bot](https://github.com/dotablaze-tech/platform/tree/main/apps/go/meowbot). It provides typed access to PostgreSQL-backed guild and user statistics, including streaks, high scores, and usage tracking. This package wraps raw SQL interactions with clean Go functions and models.

Part of the `meowbot` suite, it is managed via [Nx](https://nx.dev) in the `libs/go` workspace.

---

## 📁 Project Structure

```
libs/go/meowbot/feature/db/
├── connection.go      # Establishes DB connection with pooling and logging
├── models.go          # Structs for DB rows and query results
├── stats.go           # Core DB access functions for stats read/write
├── stats_test.go      # Unit tests for DB logic using mock/stub data
├── go.mod / go.sum    # Go module files
└── project.json       # Nx project definition
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- A running PostgreSQL instance
- [Nx CLI](https://nx.dev)

### Installation

This package is used internally by Meowbot. Import like so:

```go
import "github.com/dotablaze-tech/platform/libs/go/meowbot/feature/db"
```

---

## 🔌 Usage

Initialize with your own `*sql.DB` connection:

```go
conn := db.NewConnection("postgres://user:pass@host:5432/dbname")
err := conn.Ping()
if err != nil {
    log.Fatalf("DB unreachable: %v", err)
}

stats, err := db.GetUserGuildStats(conn, guildID, userID)
```

---

## ✨ Features

- ⚙️ Connection management with custom config
- 📊 Fetch and update guild/user streak statistics
- 🧾 Models for `users`, `guilds`, and `user_guild_stats` tables
- 🔐 Explicit, type-safe SQL operations
- 🧪 Tests for core stat logic and edge cases

---

## 🧪 Testing

```bash
go test ./libs/go/meowbot/feature/db
```

---

## 🧠 Schema Overview

This package assumes the following schema (defined in `apps/database/meowbotdb/src/00_schema.sql`):

- `guilds (id TEXT PRIMARY KEY)`
- `users (id TEXT PRIMARY KEY)`
- `user_guild_stats (guild_id TEXT, user_id TEXT, meow_count INT, high_score INT, PRIMARY KEY (guild_id, user_id))`

---

## 📌 Notes

- This package avoids global state; all DB functions are explicitly passed a `*sql.DB` instance.
- Intended for use in stateless containers with Postgres connectivity.
- Error handling is left to the caller—check all return values!
