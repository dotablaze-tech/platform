package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"
	"libs/go/meowbot/util"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"libs/go/meowbot/feature/db"
)

func New(db *sql.DB, sess *discordgo.Session) *Server {
	return &Server{
		Logger:  util.Cfg.Logger,
		DB:      db,
		Session: sess,
	}
}

func (s *Server) Start(ctx context.Context) {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/liveness", s.livenessHandler)
	mux.HandleFunc("/readiness", s.readinessHandler)

	// Stats
	mux.HandleFunc("/stats", s.statsHandler)
	mux.HandleFunc("/guilds/", s.guildStatsRouter)
	mux.HandleFunc("/users/", s.userStatsRouter)
	mux.HandleFunc("/leaderboard", s.leaderboardHandler)
	mux.HandleFunc("/users", s.usersHandler)
	mux.HandleFunc("/guilds", s.guildsHandler)

	addr := ":" + util.Cfg.ApiPort

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.Logger.Info("🌐 Starting API server", "addr", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Error("❌ API server failed", "error", err)
		}
	}()

	<-ctx.Done()
	s.Logger.Info("🛑 Shutting down API server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		s.Logger.Error("❌ Failed to shutdown API server", "error", err)
	}
}

func (s *Server) livenessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dbHealthy := s.DB.PingContext(ctx) == nil
	discordHealthy := s.Session != nil &&
		s.Session.State != nil &&
		s.Session.State.User != nil &&
		s.Session.State.User.ID != ""

	status := "healthy"
	if !dbHealthy || !discordHealthy {
		status = "unhealthy"
	}

	dbStatus := "unhealthy"
	if dbHealthy {
		dbStatus = "healthy"
	}

	discordStatus := "unhealthy"
	if discordHealthy {
		discordStatus = "healthy"
	}

	s.writeJSON(w, HealthResponse{
		Status:   status,
		Database: dbStatus,
		Discord:  discordStatus,
	})
}

func (s *Server) readinessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dbHealthy := s.DB.PingContext(ctx) == nil
	discordHealthy := s.Session != nil &&
		s.Session.State != nil &&
		s.Session.State.User != nil &&
		s.Session.State.User.ID != ""

	status := "ready"
	if !dbHealthy || !discordHealthy {
		status = "not ready"
	}

	dbStatus := "not ready"
	if dbHealthy {
		dbStatus = "ready"
	}

	discordStatus := "not ready"
	if discordHealthy {
		discordStatus = "ready"
	}

	s.writeJSON(w, HealthResponse{
		Status:   status,
		Database: dbStatus,
		Discord:  discordStatus,
	})
}

func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := db.GetGlobalStats(ctx, s.DB)
	if err != nil {
		s.Logger.Error("❌ Failed to fetch stats", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}

func (s *Server) guildStatsRouter(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "guilds" || parts[2] != "stats" {
		http.NotFound(w, r)
		return
	}
	s.guildStatsHandler(w, r, parts[1])
}

func (s *Server) guildStatsHandler(w http.ResponseWriter, r *http.Request, guildID string) {
	ctx := r.Context()
	streak, err := db.GetGuildStats(ctx, s.DB, guildID)
	if err != nil {
		s.writeError(w, http.StatusNotFound, "guild not found", err)
		return
	}

	s.writeJSON(w, streak)
}

func (s *Server) userStatsRouter(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "users" || parts[2] != "stats" {
		http.NotFound(w, r)
		return
	}
	s.userStatsHandler(w, r, parts[1])
}

func (s *Server) userStatsHandler(w http.ResponseWriter, r *http.Request, userID string) {
	stats, err := db.GetUserGlobalStats(r.Context(), s.DB, userID)
	if err != nil {
		s.writeError(w, http.StatusNotFound, "user not found", err)
		return
	}

	perGuildStats, err := db.GetUserPerGuildStats(r.Context(), s.DB, userID)
	if err != nil {
		s.Logger.Warn("⚠ Failed to fetch per-guild stats", slog.String("user_id", userID), slog.Any("error", err))
	}
	stats.GuildStats = perGuildStats

	s.writeJSON(w, stats)
}

func (s *Server) leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := db.GetLeaderboard3(r.Context(), s.DB, 10)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "internal error", err)
		return
	}

	s.writeJSON(w, entries)
}

func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetAllUsers(r.Context(), s.DB)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "failed to fetch users", err)
		return
	}
	s.writeJSON(w, users)
}

func (s *Server) guildsHandler(w http.ResponseWriter, r *http.Request) {
	guilds, err := db.GetAllGuilds(r.Context(), s.DB)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "failed to fetch guilds", err)
		return
	}
	s.writeJSON(w, guilds)
}

func (s *Server) writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.Logger.Error("❌ Failed to encode response", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (s *Server) writeError(w http.ResponseWriter, status int, msg string, err error) {
	s.Logger.Error("❌ "+msg, "error", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": msg}); encodeErr != nil {
		s.Logger.Error("❌ Failed to write error JSON", "error", encodeErr)
	}
}
