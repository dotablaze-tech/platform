package main

import (
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	logger := slog.Default()
	token := os.Getenv("DISCORD_BOT_TOKEN")

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error("Failed to create session", "error", err)
		os.Exit(1)
	}

	err = sess.Open()
	if err != nil {
		logger.Error("Failed to open session", "error", err)
		os.Exit(1)
	}
	defer sess.Close()

	// Grab the app ID directly from the session state — no env var needed
	appID := sess.State.User.ID

	cmds, err := sess.ApplicationCommands(appID, "")
	if err != nil {
		logger.Error("Failed to fetch commands", "error", err)
		os.Exit(1)
	}

	for _, cmd := range cmds {
		err := sess.ApplicationCommandDelete(appID, "", cmd.ID)
		if err != nil {
			logger.Error("Failed to delete", "name", cmd.Name, "error", err)
		} else {
			logger.Info("Deleted", "name", cmd.Name)
		}
	}
}
