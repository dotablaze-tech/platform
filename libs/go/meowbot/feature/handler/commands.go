package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"libs/go/meowbot/feature/state"
	"log/slog"
)

func CommandHandler(logger *slog.Logger) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		guildID := i.GuildID
		gs := state.GetOrCreate(guildID)

		switch i.ApplicationCommandData().Name {
		case "meowcount":
			resp := fmt.Sprintf("Current meow count: %d", gs.MeowCount)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: resp,
				},
			})
			if err != nil {
				logger.Error("Failed to respond to /meowcount", "error", err, "guildID", guildID)
			} else {
				logger.Info("Responded to /meowcount", "guildID", guildID, "count", gs.MeowCount)
			}

		case "highscore":
			resp := "No high score yet!"
			if gs.HighScore > 0 {
				resp = fmt.Sprintf("🏆 High score: %d by %s", gs.HighScore, gs.HighScoreUser)
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: resp,
				},
			})
			if err != nil {
				logger.Error("Failed to respond to /highscore", "error", err, "guildID", guildID)
			} else {
				logger.Info("Responded to /highscore", "guildID", guildID, "score", gs.HighScore)
			}
		}
	}
}
