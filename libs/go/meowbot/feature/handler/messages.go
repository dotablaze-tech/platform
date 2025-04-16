package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"libs/go/meowbot/feature/state"
	"libs/go/meowbot/util"
	"log/slog"
	"regexp"
	"strings"
)

var meowRegex = regexp.MustCompile(`(?i)^m+e+o+w+$`)

func MessageHandler(logger *slog.Logger, allowedChannel string) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if allowedChannel != "" && m.ChannelID != allowedChannel {
			logger.Debug("Ignored message from unauthorized channel",
				"channelID", m.ChannelID,
				"guildID", m.GuildID,
			)
			return
		}

		guildID := m.GuildID
		content := strings.ToLower(strings.TrimSpace(m.Content))
		gs := state.GetOrCreate(guildID)

		logger.Info("Message received",
			"user", m.Author.Username,
			"userID", m.Author.ID,
			"content", m.Content,
			"channelID", m.ChannelID,
			"guildID", guildID,
		)

		if meowRegex.MatchString(content) {
			if m.Author.ID == gs.LastUserID {
				if _, err := s.ChannelMessageSend(m.ChannelID, "😾 You can't meow twice in a row!"); err != nil {
					logger.Error("Failed to send repeat warning",
						"error", err,
						"guildID", guildID,
						"channelID", m.ChannelID,
					)
				}
				logger.Warn("Repeat meow",
					"user", m.Author.Username,
					"guildID", guildID,
				)
				state.Reset(guildID)
				return
			}

			gs.MeowCount++

			if gs.MeowCount > gs.HighScore {
				gs.HighScore = gs.MeowCount
				gs.HighScoreUser = m.Author.Username

				if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🏆 New high score: %d meows by %s!", gs.HighScore, m.Author.Username)); err != nil {
					logger.Error("Failed to send high score message", "error", err, "guildID", guildID)
				}
				logger.Info("New high score", "user", m.Author.Username, "score", gs.HighScore, "guildID", guildID)
			}

			gs.LastUserID = m.Author.ID
			msg := fmt.Sprintf("%s **meow** x%d!", util.RandomEmoji(), gs.MeowCount)

			if _, err := s.ChannelMessageSend(m.ChannelID, msg); err != nil {
				logger.Error("Failed to send meow message",
					"error", err,
					"guildID", guildID,
					"channelID", m.ChannelID,
				)
			} else {
				logger.Info("Meow counted",
					"user", m.Author.Username,
					"count", gs.MeowCount,
					"guildID", guildID,
				)
			}
		} else {
			if _, err := s.ChannelMessageSend(m.ChannelID, "❌ No meow? Resetting."); err != nil {
				logger.Error("Failed to send reset message",
					"error", err,
					"guildID", guildID,
					"channelID", m.ChannelID,
				)
			}
			state.Reset(guildID)
			logger.Info("Reset triggered",
				"user", m.Author.Username,
				"guildID", guildID,
			)
		}
	}
}
