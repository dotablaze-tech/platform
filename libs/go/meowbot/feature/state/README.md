# 🧠 Meowbot State

![Build](https://img.shields.io/github/actions/workflow/status/dotablaze-tech/apps/ci.yml?branch=main)
![Nx](https://img.shields.io/badge/Nx-managed-blue)
![Go Module](https://img.shields.io/badge/Go-Module-brightgreen)

**Meowbot State** is a modular Go library for managing in-memory state tracking within the [🐾 Meow Bot](https://github.com/dotablaze-tech/apps/tree/main/apps/go/meowbot). It maintains per-guild data such as streak counts, user participation, and high score tracking to enable responsive, context-aware bot behavior.

This package is part of the monorepo’s feature set and is managed by [Nx](https://nx.dev).

---

## 📁 Project Structure

```
libs/go/meowbot/feature/state/
├── state.go         # Core state management logic
├── state_test.go    # Unit tests for state behavior
├── go.mod / go.sum  # Go module definition
└── project.json     # Nx project definition
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- [Nx CLI](https://nx.dev)

### Installation

Used internally within Meowbot features. To import:

```go
import "github.com/dotablaze-tech/apps/libs/go/meowbot/feature/state"
```

---

## ✨ Features

- ✅ Tracks per-guild streaks and user activity
- 🚫 Prevents same-user repeat meows
- 🏆 Maintains high score history
- 🔁 Provides reset logic and accessors
- 🔒 Thread-safe for concurrent access

---

## 🧪 Testing

Run unit tests to verify state behavior:

```bash
go test ./libs/go/meowbot/feature/state
```

---

## 🧠 Example Usage

```go
state := NewGuildState()

ok := state.TryMeow("user123")
if ok {
    fmt.Println("Meow accepted! Current streak:", state.MeowCount)
} else {
    fmt.Println("Invalid meow! You can't meow twice in a row.")
}

highScoreUser := state.HighScoreUser
```

---

## 📌 Notes

- The state is ephemeral by default. For persistence, pair it with the [db](../db) feature.
- Designed for single-channel operation per guild.
- This package does not include Discord-specific logic—pure Go state.
