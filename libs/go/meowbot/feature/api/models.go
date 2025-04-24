package api

import (
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

type Server struct {
	Logger  *slog.Logger
	DB      *sql.DB
	Session *discordgo.Session
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Discord  string `json:"discord"`
}
