package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"libs/go/meowbot/feature/api"
	"libs/go/meowbot/feature/db"
	"libs/go/meowbot/feature/handler"
	"libs/go/meowbot/util"
	"os"
	"os/signal"
	"syscall"
)

func Run(ctx context.Context, cfg util.AppConfig) error {
	util.Cfg.Logger.Info("🚀 Booting up Meow bot...",
		"mode", util.Cfg.Mode,
		"debug", util.Cfg.Debug,
	)

	// Initialize emojis and DB connection
	util.InitEmojis()
	if err := db.InitDB(ctx); err != nil {
		return err
	}
	util.Cfg.Logger.Info("✅ Connected to meowbot PostgreSQL!")

	// Create Discord session
	sess, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return err
	}

	// Set up intents and handlers
	sess.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
	apiServer := api.New(db.DB, sess)
	apiCtx, apiCancel := context.WithCancel(ctx)
	defer apiCancel()
	go apiServer.Start(apiCtx)

	// Add message handlers
	sess.AddHandler(handler.MessageHandler(ctx))
	sess.AddHandler(handler.CommandHandler(ctx))
	sess.AddHandler(handler.ComponentHandler(ctx))

	// Open Discord session
	if err := sess.Open(); err != nil {
		return err
	}
	defer func() {
		if err := sess.Close(); err != nil {
			util.Cfg.Logger.Error("❌ Failed to close Discord session", "error", err)
		} else {
			util.Cfg.Logger.Info("✅ Successfully closed Discord session.")
		}
	}()

	// Register commands
	if err := handler.RegisterCommands(sess); err != nil {
		return err
	}

	util.Cfg.Logger.Info("🐱 Meow bot is online!")

	// Wait for termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	// Graceful shutdown
	if err := db.CloseDB(); err != nil {
		return err
	}

	util.Cfg.Logger.Info("👋 Meow bot has shut down gracefully.")
	return nil
}

func main() {
	// Ensure bot token is available before proceeding
	if util.Cfg.BotToken == "" {
		util.Cfg.Logger.Error("❌ DISCORD_BOT_TOKEN not set. Exiting.")
		os.Exit(1)
	}

	// Run the bot and handle errors
	if err := Run(context.Background(), util.Cfg); err != nil {
		util.Cfg.Logger.Error("❌ Error while running bot", "error", err)
		os.Exit(1)
	}
}
