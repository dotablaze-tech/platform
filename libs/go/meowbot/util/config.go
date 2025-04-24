package util

import (
	"github.com/jba/slog/handlers/loghandler"
	"log/slog"
	"os"
	"strings"
)

var Cfg = LoadConfig()

type AppConfig struct {
	Mode             string
	Debug            bool
	IsProd           bool
	BotToken         string
	ApiPort          string
	DatabaseURL      string
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	EmojiList        string
	Logger           *slog.Logger
	Whitelist        struct {
		Guilds []string
	}
}

func LoadConfig() AppConfig {
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "production"
	}

	debug := os.Getenv("DEBUG") == "true"

	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	logger := slog.New(loghandler.New(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080"
	}

	guildsCSV := os.Getenv("WHITELISTED_GUILDS")
	var guilds []string
	if guildsCSV != "" {
		guilds = strings.Split(guildsCSV, ",")
	}

	return AppConfig{
		Mode:             mode,
		Debug:            debug,
		IsProd:           mode == "production",
		ApiPort:          apiPort,
		BotToken:         os.Getenv("DISCORD_BOT_TOKEN"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		EmojiList:        os.Getenv("EMOJI_LIST"),
		Logger:           logger,
		Whitelist: struct {
			Guilds []string
		}{Guilds: guilds},
	}
}

func (cfg AppConfig) IsAllowedGuild(guildID string) bool {
	if cfg.IsProd {
		return true
	}
	for _, id := range cfg.Whitelist.Guilds {
		if id == guildID {
			return true
		}
	}
	return false
}
