# 🌐 Meowbot API

![Build](https://img.shields.io/github/actions/workflow/status/dotablaze-tech/platform/ci.yml?branch=main)
![Nx](https://img.shields.io/badge/Nx-managed-blue)
![Go Module](https://img.shields.io/badge/Go-Module-brightgreen)

**Meowbot API** is a lightweight HTTP layer for exposing Meowbot statistics via a RESTful interface. Built using the [`chi`](https://github.com/go-chi/chi) router and structured response models, it enables external systems or dashboards to consume streak and score data from Meowbot's backend.

Part of the modular Meowbot architecture, managed via [Nx](https://nx.dev) in the `libs/go` workspace.

---

## 📁 Project Structure

```
libs/go/meowbot/feature/api/
├── api.go              # HTTP router, middleware, and handlers
├── models.go           # Response models for stats and errors
├── go.mod / go.sum     # Go module definition
└── project.json        # Nx project definition
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- [Nx CLI](https://nx.dev)
- PostgreSQL (if using with Meowbot's DB package)

### Installation

This package is intended to be used internally by the Meowbot platform:

```go
import "github.com/dotablaze-tech/platform/libs/go/meowbot/feature/api"
```

---

## 🔌 Usage

Initialize the router with dependency injection:

```go
r := api.NewRouter(api.Options{
    DB: myPostgresConn, // *sql.DB
    Logger: myLogger,   // slog.Logger
})
http.ListenAndServe(":8080", r)
```

---

## ✨ Features

- 🚏 REST-style endpoints for streak and score retrieval
- 📦 Typed JSON responses with helpful error structure
- 🔐 Graceful error handling and logging
- 🔄 Ready for container-based deployment
- 🧩 Easily extendable with new routes or middlewares

---

## 🧪 Testing

*Test coverage coming soon.*

---

## 📌 Notes

- Routes are not versioned yet (`/leaderboard`, `/highscore` etc.).
- Intended to be served behind a proxy or gateway (e.g. Ingress).
- All endpoints return JSON content.

---

## 📘 Example Endpoints (Planned)

- `GET /leaderboard?guild_id=123&sort=streak`
- `GET /highscore?guild_id=123`
- `GET /stats?guild_id=123&user_id=456`
