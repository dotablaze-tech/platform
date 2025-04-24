package db

import (
	"time"
)

type Guild struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type UserGuildStats struct {
	GuildID          string     `json:"guild_id"`
	UserID           string     `json:"user_id"`
	SuccessfulMeows  int        `json:"successful_meows"`
	FailedMeows      int        `json:"failed_meows"`
	TotalMeows       int        `json:"total_meows"`
	CurrentStreak    int        `json:"current_streak"`
	HighestStreak    int        `json:"highest_streak"`
	LastMeowAt       *time.Time `json:"last_meow_at,omitempty"`
	LastFailedMeowAt *time.Time `json:"last_failed_meow_at,omitempty"`
}

type GuildStreak struct {
	GuildID         string  `json:"guild_id"`
	MeowCount       int     `json:"meow_count"`
	LastUserID      *string `json:"last_user_id,omitempty"`
	HighScore       int     `json:"high_score"`
	HighScoreUserID *string `json:"high_score_user_id,omitempty"`
}

type GlobalStats struct {
	TotalGuilds int `json:"total_guilds"`
	TotalUsers  int `json:"total_users"`
	TotalMeows  int `json:"total_meows"`
}

type LeaderboardEntry struct {
	User            *User `json:"user"`
	TotalMeows      int   `json:"total_meows"`
	SuccessfulMeows int   `json:"successful_meows"`
	FailedMeows     int   `json:"failed_meows"`
}

type GuildStats struct {
	Guild           *Guild `json:"guild"`
	CurrentStreak   int    `json:"current_streak"`
	LastUser        *User  `json:"last_user,omitempty"`
	HighScore       int    `json:"high_score"`
	HighScoreUser   *User  `json:"high_score_user,omitempty"`
	TotalMeows      int    `json:"total_meows"`
	SuccessfulMeows int    `json:"successful_meows"`
	FailedMeows     int    `json:"failed_meows"`
}

type UserGlobalStats struct {
	UserID          string           `json:"user_id"`
	Username        string           `json:"username"`
	CreatedAt       time.Time        `json:"created_at"`
	SuccessfulMeows int              `json:"successful_meows"`
	FailedMeows     int              `json:"failed_meows"`
	TotalMeows      int              `json:"total_meows"`
	HighestStreak   int              `json:"highest_streak"`
	GuildStats      []UserGuildStats `json:"guild_stats,omitempty"`
}
