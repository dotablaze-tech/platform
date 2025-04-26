package handler

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"libs/go/meowbot/feature/db"
	"libs/go/meowbot/feature/state"
	"libs/go/meowbot/util"
	"regexp"
	"strings"
	"time"
)

var meowRegex = regexp.MustCompile(`(?i)^m+e+o+w+$`)

func sendMessage(s *discordgo.Session, channelID, message, guildID string) error {
	if !util.Cfg.IsAllowedGuild(guildID) {
		util.Cfg.Logger.Debug("✉️ [DEV] Skipped sending message", "guildID", guildID, "channelID", channelID, "message", message)
		return nil
	}
	_, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		util.Cfg.Logger.Error("❌ Failed to send message", "guildID", guildID, "channelID", channelID, "error", err)
	}
	return err
}

func safeReact(s *discordgo.Session, channelID, messageID, emoji, guildID string) {
	if !util.Cfg.IsAllowedGuild(guildID) {
		util.Cfg.Logger.Debug("🔕 [DEV] Skipped reaction", "guildID", guildID, "emoji", emoji)
		return
	}
	if err := s.MessageReactionAdd(channelID, messageID, emoji); err != nil {
		util.Cfg.Logger.Warn("⚠️ Failed to react", "emoji", emoji, "channelID", channelID, "messageID", messageID, "error", err)
	}
}

func upsertEntities(ctx context.Context, user *discordgo.User, guildID string) {
	if err := db.UpsertGuild(ctx, db.DB, db.Guild{ID: guildID}); err != nil {
		util.Cfg.Logger.Error("❌ Failed to upsert guild", "guildID", guildID, "error", err)
	}
	if err := db.UpsertUser(ctx, db.DB, db.User{
		ID:       user.ID,
		Username: user.Username,
	}); err != nil {
		util.Cfg.Logger.Error("❌ Failed to upsert user", "userID", user.ID, "username", user.Username, "error", err)
	}
}

func incrementMeow(ctx context.Context, guildID string, userID string, isMeow bool, timestamp time.Time) {
	if err := db.IncrementMeow(ctx, db.DB, guildID, userID, isMeow, timestamp); err != nil {
		util.Cfg.Logger.Error("❌ Failed to increment meow", "guildID", guildID, "userID", userID, "error", err)
	}
}

func isInAllowedChannel(ctx context.Context, m *discordgo.MessageCreate) bool {
	allowedChannelID, err := db.GetChannelForGuild(ctx, db.DB, m.GuildID)
	if err != nil {
		util.Cfg.Logger.Error("❌ Could not fetch allowed channel", "guildID", m.GuildID, "channelID", m.ChannelID, "error", err)
		return false
	}
	if allowedChannelID == "" || m.ChannelID != allowedChannelID {
		util.Cfg.Logger.Debug("🚫 Message in unauthorized channel", "guildID", m.GuildID, "channelID", m.ChannelID, "allowedChannelID", allowedChannelID)
		return false
	}
	return true
}

func processMeowMessage(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.ToLower(strings.TrimSpace(m.Content))
	guildID := m.GuildID
	user := m.Author
	gs := state.GetOrCreate(ctx, guildID)

	util.Cfg.Logger.Info("📬 Message received", "guildID", guildID, "channelID", m.ChannelID, "userID", user.ID, "username", user.Username, "content", m.Content)

	if meowRegex.MatchString(content) {
		handleMeow(ctx, s, m, gs)
	} else {
		handleNonMeow(ctx, s, m)
	}
}

func handleMeow(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, gs *state.GuildState) {
	user := m.Author
	guildID := m.GuildID

	if user.ID == gs.LastUserID {
		incrementMeow(ctx, guildID, user.ID, false, m.Timestamp)
		safeReact(s, m.ChannelID, m.ID, "❌", guildID)
		err := sendMessage(s, m.ChannelID, "😾 You can't meow twice in a row!", guildID)
		if err != nil {
			return
		}
		util.Cfg.Logger.Warn("🔂 Repeat meow", "guildID", guildID, "userID", user.ID)
		state.Reset(guildID)
		return
	}

	gs.MeowCount++
	if gs.MeowCount > gs.HighScore {
		gs.HighScore = gs.MeowCount
		gs.HighScoreUserID = user.ID
		err := sendMessage(s, m.ChannelID, fmt.Sprintf("🏆 New high score: %d meows by %s!", gs.HighScore, user.Username), guildID)
		if err != nil {
			return
		}
		util.Cfg.Logger.Info("🏆 New high score", "guildID", guildID, "userID", user.ID, "score", gs.HighScore)
	}

	gs.LastUserID = user.ID
	incrementMeow(ctx, guildID, user.ID, true, m.Timestamp)
	err := sendMessage(s, m.ChannelID, fmt.Sprintf("%s **meow** x%d!", util.RandomEmoji(), gs.MeowCount), guildID)
	if err != nil {
		return
	}
	safeReact(s, m.ChannelID, m.ID, "🐱", guildID)

	err = db.UpsertGuildStreak(ctx, db.DB, db.GuildStreak{
		GuildID:         guildID,
		MeowCount:       gs.MeowCount,
		LastUserID:      &gs.LastUserID,
		HighScore:       gs.HighScore,
		HighScoreUserID: &gs.HighScoreUserID,
	})
	if err != nil {
		util.Cfg.Logger.Error("❌ Failed to upsert guild streak", "guildID", guildID, "error", err)
		return
	}
}

func handleNonMeow(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author
	guildID := m.GuildID

	safeReact(s, m.ChannelID, m.ID, "❌", guildID)
	err := sendMessage(s, m.ChannelID, "❌ No meow? Resetting.", guildID)
	if err != nil {
		return
	}
	incrementMeow(ctx, guildID, user.ID, false, m.Timestamp)
	state.Reset(guildID)

	util.Cfg.Logger.Info("🔄 Reset triggered", "guildID", guildID, "userID", user.ID)
}

func logIgnoreBotMessage(m *discordgo.MessageCreate) {
	util.Cfg.Logger.Debug("🤖 Ignored bot message", "guildID", m.GuildID, "channelID", m.ChannelID)
}

func MessageHandler(ctx context.Context) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// skip bot messages
		if m.Author.Bot {
			logIgnoreBotMessage(m)
			return
		}

		if !isInAllowedChannel(ctx, m) {
			return
		}

		// Upsert user + guild
		upsertEntities(ctx, m.Author, m.GuildID)

		processMeowMessage(ctx, s, m)
	}
}
