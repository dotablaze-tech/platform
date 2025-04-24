# 🧰 Meowbot Util

![Build](https://img.shields.io/github/actions/workflow/status/dotablaze-tech/platform/ci.yml?branch=main)
![Nx](https://img.shields.io/badge/Nx-managed-blue)
![Go Module](https://img.shields.io/badge/Go-Module-brightgreen)

**Meowbot Util** is a utility library for
the [🐾 Meow Bot](https://github.com/dotablaze-tech/platform/tree/main/apps/go/meowbot), providing reusable helpers,
configuration loading, and emoji constants. It promotes clean separation of concerns and shared logic across Meowbot
features.

This library is managed with [Nx](https://nx.dev) and structured as a standalone Go module for use in feature packages
like `handler`, `state`, and `api`.

---

## 📁 Project Structure

```
libs/go/meowbot/util/
├── config.go           # Bot config loading from env
├── config_test.go      # Tests for config loading
├── emojis.go           # Emoji constants used by Meowbot
├── emojis_test.go      # Emoji tests and verification
├── go.mod / go.sum     # Independent module for utility library
└── project.json        # Nx project definition
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- [Nx CLI](https://nx.dev)

### Installation

This module is used internally in the monorepo. To use it in another Go module:

```go
import "github.com/dotablaze-tech/platform/libs/go/meowbot/util"
```

Make sure your `go.work` and `go.mod` are properly wired if using outside of Nx context.

---

## ✨ Features

- **Emoji Constants**: Standardized emojis for Meowbot interactions (`😺`, `😼`, `😾`, etc.).
- **Config Loader**: Reads bot configuration (e.g., `DISCORD_BOT_TOKEN`) via `config.go`.
- **Test Coverage**: Unit tests included to validate config parsing and emoji availability.

---

## 🧪 Testing

To run tests for the util package:

```bash
go test ./libs/go/meowbot/util
```

---

## 🔧 Usage Example

```go
import (
    "fmt"
    "os"

    "github.com/dotablaze-tech/platform/libs/go/meowbot/util"
)

func main() {
    config, err := util.LoadConfigFromEnv()
    if err != nil {
        panic(err)
    }
    fmt.Println("Bot token:", config.Token)

    fmt.Println("Happy Meow Emoji:", util.EmojiHappy)
}
```

---

## 📌 Notes

- This library is a foundational utility used across all `meowbot` feature packages.
- Consider extending `config.go` if additional environment-based configuration is needed.
