package handler

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"libs/go/meowbot/feature/db"
	"libs/go/meowbot/feature/state"
	"libs/go/meowbot/util"
	"strconv"
	"strings"
	"time"
)

const leaderboardPageSize = 5

func sendResponseEmbed(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	embed *discordgo.MessageEmbed,
	guildID string,
	commandName string,
) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		util.Cfg.Logger.Error(fmt.Sprintf("❌ Failed to respond to /%s", commandName), "error", err, "guildID", guildID)
	} else {
		util.Cfg.Logger.Info(fmt.Sprintf("💬 Responded to /%s", commandName), "guildID", guildID, "response", embed.Description)
	}
}

// CommandHandler manages slash commands
func CommandHandler(ctx context.Context) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		guildID := i.GuildID
		gs := state.GetOrCreate(ctx, guildID)

		switch i.ApplicationCommandData().Name {
		case "count":
			handleCount(s, i, gs)
		case "highscore":
			handleHighscore(s, i, gs)
		case "stats":
			handleStats(ctx, s, i)
		case "setup":
			handleSetup(ctx, s, i)
		case "leaderboard":
			handleLeaderboard(ctx, s, i)

		default:
			util.Cfg.Logger.Warn("⚠️ Unknown command", "guildID", guildID, "command", i.ApplicationCommandData().Name)
		}
	}
}

func handleCount(s *discordgo.Session, i *discordgo.InteractionCreate, gs *state.GuildState) {
	title := "📈 Meow Count"
	desc := fmt.Sprintf("Current meow count: **%d**", gs.MeowCount)

	embed := formatSimpleEmbed(title, desc)
	sendResponseEmbed(s, i, embed, i.GuildID, "count")
}

func handleHighscore(s *discordgo.Session, i *discordgo.InteractionCreate, gs *state.GuildState) {
	title := "🏆 High Score"
	desc := "😿 No high score yet!"
	if gs.HighScore > 0 {
		desc = fmt.Sprintf("High score: **%d** by <@%s>", gs.HighScore, gs.HighScoreUserID)
	}

	embed := formatSimpleEmbed(title, desc)
	sendResponseEmbed(s, i, embed, i.GuildID, "highscore")
}

func handleStats(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	scope := "guild"
	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "scope":
			scope = opt.StringValue()
		}
	}

	guildID := &i.GuildID
	scopeTitle := "Guild Stats"
	if scope == "global" {
		guildID = nil
		scopeTitle = "Global Stats"
	}

	stats, err := db.GetUserStats(ctx, db.DB, guildID, i.Member.User.ID)
	if err != nil {
		sendErrorEmbed(s, i, "❌ Failed to Fetch Stats", "Couldn't fetch your stats. You might not have any meows yet!", i.GuildID, "stats", err)
		return
	}

	lastMeow := "N/A"
	if stats.LastMeowAt != nil {
		lastMeow = time.Since(*stats.LastMeowAt).Round(time.Second).String() + " ago"
	}

	title := fmt.Sprintf("📊 **Your Meows — %s**", scopeTitle)
	resp := fmt.Sprintf(
		"📈 Total Meows: %d\n"+
			"✅ Successful Meows: %d\n"+
			"❌ Failed Meows: %d\n"+
			"🔁 Highest Streak: %d\n"+
			"🔥 Current Streak: %d\n"+
			"⏱️ Last Meow: %s",
		stats.TotalMeows,
		stats.SuccessfulMeows,
		stats.FailedMeows,
		stats.HighestStreak,
		stats.CurrentStreak,
		lastMeow,
	)

	embed := formatSimpleEmbed(title, resp)
	sendResponseEmbed(s, i, embed, i.GuildID, "stats")
}

func handleSetup(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID

	if i.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		embed := formatSimpleEmbed("🚫 Permission Denied", "You need to be a server admin to use this command.")
		sendResponseEmbed(s, i, embed, guildID, "setup")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) != 1 || options[0].Name != "channel" {
		embed := formatSimpleEmbed("⚠️ Invalid Usage", "You must provide a channel using `/setup channel:#channel-name`.", 0xffff00)
		sendResponseEmbed(s, i, embed, guildID, "setup")
		return
	}

	channelOpt := options[0]
	channelID := channelOpt.ChannelValue(s).ID

	err := db.UpsertGuildChannel(ctx, db.DB, guildID, channelID)
	if err != nil {
		sendErrorEmbed(s, i, "❌ Failed to Set Channel", "Failed to set meow channel. Try again later.", guildID, "setup", err)
		return
	}

	title := "⚙ Setup Complete"
	resp := fmt.Sprintf("✅ Meow channel has been set to <#%s>", channelID)
	sendSuccessEmbed(s, i, title, resp, guildID, "setup")
}

func handleLeaderboard(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Default options
	scope := "guild"
	metric := "total"
	page := 1

	// Parse options
	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "scope":
			scope = opt.StringValue()
		case "metric":
			metric = opt.StringValue()
		case "page":
			page = int(opt.IntValue())
		}
	}

	if page < 1 {
		page = 1
	}

	// Determine metric column
	var column string
	switch metric {
	case "success":
		column = "successful_meows"
	case "fail":
		column = "failed_meows"
	default:
		column = "total_meows"
	}

	// Set scope
	var guildID *string
	if scope == "guild" {
		id := i.GuildID
		guildID = &id
	}

	// Fetch leaderboard data from DB
	_, totalCount, err := db.GetLeaderboard(ctx, db.DB, guildID, column, 0, 0)
	if err != nil {
		sendErrorEmbed(s, i, "❌ Failed to Fetch Leaderboard", "Something went wrong while retrieving leaderboard data.", i.GuildID, "leaderboard", err)
		return
	}

	maxPages := (totalCount + leaderboardPageSize - 1) / leaderboardPageSize
	if maxPages < 1 {
		maxPages = 1
	}

	if page > maxPages {
		page = maxPages
	}

	offset := (page - 1) * leaderboardPageSize
	entries, _, err := db.GetLeaderboard(ctx, db.DB, guildID, column, leaderboardPageSize, offset)
	if err != nil {
		sendErrorEmbed(s, i, "❌ Failed to Fetch Leaderboard", "Something went wrong while retrieving leaderboard data.", i.GuildID, "leaderboard", err)
		return
	}

	if len(entries) == 0 {
		embed := formatSimpleEmbed("📉 Empty Leaderboard", "No one has meowed yet! Be the first.", 0xFEE75C)
		sendResponseEmbed(s, i, embed, i.GuildID, "leaderboard")
		return
	}

	// Fetch user's rank if interaction is from a user
	userRank, rankErr := db.GetUserRank(ctx, db.DB, i.Member.User.ID, guildID, column)

	// Format embed and buttons
	embed := formatLeaderboardEmbed(entries, scope, metric, page, totalCount, userRank, rankErr, i.Member.User.ID)
	components := renderLeaderboardButtons(scope, metric, page, totalCount)

	// Respond with leaderboard
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		util.Cfg.Logger.Error("❌ Failed to send leaderboard response:", "error", err)
	}
}

// UnregisterCommands unregisters slash commands with Discord
func UnregisterCommands(sess *discordgo.Session) error {
	cmds, _ := sess.ApplicationCommands(sess.State.User.ID, "")
	for _, cmd := range cmds {
		err := sess.ApplicationCommandDelete(sess.State.User.ID, "", cmd.ID)
		if err != nil {
			return err
		}
		util.Cfg.Logger.Info("✅ Unregistered command", "command", cmd.Name)
	}
	return nil
}

// RegisterCommands registers slash commands with Discord
func RegisterCommands(sess *discordgo.Session) error {
	commands := []discordgo.ApplicationCommand{
		{
			Name:        "count",
			Description: "Check the current meow count for this server",
		},
		{
			Name:        "highscore",
			Description: "Check the highest meow streak for this server",
		},
		{
			Name:        "stats",
			Description: "Check your personal meow stats",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "scope",
					Description: "Whether to show the guild or global stats",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Guild", Value: "guild"},
						{Name: "Global", Value: "global"},
					},
				},
			},
		},
		{
			Name:        "setup",
			Description: "Configure Meow Bot for this server",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Channel where Meow Bot should listen for meows",
					Required:    true,
				},
			},
		},
		{
			Name:        "leaderboard",
			Description: "Show the top meowers",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "scope",
					Description: "Whether to show the guild or global leaderboard",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Guild", Value: "guild"},
						{Name: "Global", Value: "global"},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "metric",
					Description: "Leaderboard metric",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Total Meows", Value: "total"},
						{Name: "Successful Meows", Value: "success"},
						{Name: "Failed Meows", Value: "fail"},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "page",
					Description: "Page number of the leaderboard",
					Required:    false,
				},
			},
		},
	}

	for _, cmd := range commands {
		if util.Cfg.IsProd {
			_, err := sess.ApplicationCommandCreate(sess.State.User.ID, "", &cmd)
			if err != nil {
				util.Cfg.Logger.Error("❌ Failed to register command", "command", cmd.Name, "error", err)
				return err
			}
			util.Cfg.Logger.Info("✅ Registered command", "command", cmd.Name)
		} else {
			for _, guildId := range util.Cfg.Whitelist.Guilds {
				_, err := sess.ApplicationCommandCreate(sess.State.User.ID, guildId, &cmd)
				if err != nil {
					util.Cfg.Logger.Error("❌ Failed to register command", "command", cmd.Name, "guildID", guildId, "error", err)
					return err
				}
				util.Cfg.Logger.Info("✅ Registered command", "command", cmd.Name, "guildID", guildId)
			}
		}
	}

	return nil
}

func renderLeaderboardButtons(scope string, metric string, page, total int) []discordgo.MessageComponent {
	totalPages := (total + leaderboardPageSize - 1) / leaderboardPageSize

	if totalPages <= 1 {
		return nil
	}

	firstDisabled := page == 1
	prevDisabled := page <= 1
	nextDisabled := page >= totalPages
	lastDisabled := page == totalPages

	firstStyle := discordgo.PrimaryButton
	prevStyle := discordgo.PrimaryButton
	nextStyle := discordgo.PrimaryButton
	lastStyle := discordgo.PrimaryButton

	if firstDisabled {
		firstStyle = discordgo.SecondaryButton
	}
	if prevDisabled {
		prevStyle = discordgo.SecondaryButton
	}
	if nextDisabled {
		nextStyle = discordgo.SecondaryButton
	}
	if lastDisabled {
		lastStyle = discordgo.SecondaryButton
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "⏮️ First",
					Style:    firstStyle,
					CustomID: fmt.Sprintf("lb_goto:1:%s:%s", scope, metric),
					Disabled: firstDisabled,
				},
				discordgo.Button{
					Label:    "◀️ Prev",
					Style:    prevStyle,
					CustomID: fmt.Sprintf("lb_prev:%d:%s:%s", page, scope, metric),
					Disabled: prevDisabled,
				},
				discordgo.Button{
					Label:    "Next ▶️",
					Style:    nextStyle,
					CustomID: fmt.Sprintf("lb_next:%d:%s:%s", page, scope, metric),
					Disabled: nextDisabled,
				},
				discordgo.Button{
					Label:    "Last ⏭️",
					Style:    lastStyle,
					CustomID: fmt.Sprintf("lb_goto:%d:%s:%s", totalPages, scope, metric),
					Disabled: lastDisabled,
				},
			},
		},
	}
}

func formatLeaderboardEmbed(entries []db.LeaderboardEntry, scope, metric string, page, total int, userRank int, rankErr error, currentUserID string) *discordgo.MessageEmbed {
	var sb strings.Builder
	startRank := (page-1)*leaderboardPageSize + 1

	for i, entry := range entries {
		rank := startRank + i
		count := getCountByMetric(entry, metric)

		// Highlight current user
		line := fmt.Sprintf("**%2d.** <@%s> — %d\n", rank, entry.User.ID, count)
		if entry.User.ID == currentUserID {
			line = fmt.Sprintf("**%2d.** 👑 <@%s> — %d\n", rank, entry.User.ID, count)
		}

		sb.WriteString(line)
	}

	title := buildTitle(scope, metric)
	start := (page-1)*leaderboardPageSize + 1
	end := start + len(entries) - 1

	footerText := fmt.Sprintf("📄 Page %d — Showing ranks %d–%d of %d", page, start, end, total)
	if rankErr == nil {
		footerText += fmt.Sprintf(" | Your Rank: #%d", userRank)
	}

	// Set color based on metric
	var color int
	switch metric {
	case "success":
		color = 0x00cc66 // green
	case "fail":
		color = 0xcc3300 // red
	case "total":
		color = 0x3399ff // blue
	default:
		color = 0xaaaaaa // gray fallback
	}

	return &discordgo.MessageEmbed{
		Title:       title,
		Description: sb.String(),
		Color:       color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: footerText,
		},
	}
}

func formatSimpleEmbed(title, description string, color ...int) *discordgo.MessageEmbed {
	embedColor := 0x5865F2 // Default Discord blurple
	if len(color) > 0 {
		embedColor = color[0]
	}
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       embedColor,
	}
}

func sendErrorEmbed(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	title, message string,
	guildID string,
	context string,
	err error,
) {
	util.Cfg.Logger.Error(title, "guildID", guildID, "error", err)
	embed := formatSimpleEmbed(title, fmt.Sprintf("%s\n\n```%s```", message, err.Error()), 0xED4245) // red
	sendResponseEmbed(s, i, embed, guildID, context)
}

func sendSuccessEmbed(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	title, message string,
	guildID string,
	context string,
) {
	embed := formatSimpleEmbed(title, message, 0x57F287) // green
	sendResponseEmbed(s, i, embed, guildID, context)
}

func ComponentHandler(ctx context.Context) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			handleLeaderboardPagination(ctx, s, i)
			return
		}
	}
}

func getCountByMetric(e db.LeaderboardEntry, metric string) int {
	switch metric {
	case "success":
		return e.SuccessfulMeows
	case "fail":
		return e.FailedMeows
	default:
		return e.TotalMeows
	}
}

func buildTitle(scope, metric string) string {
	var scopeLabel, metricLabel string

	// Scope context
	if scope == "global" {
		scopeLabel = "Global Leaderboard 🌐"
	} else {
		scopeLabel = "Guild Leaderboard 🏠"
	}

	// Metric context
	switch metric {
	case "success":
		metricLabel = "Most Successful Meows"
	case "fail":
		metricLabel = "Most Failed Meows"
	default:
		metricLabel = "Most Total Meows"
	}

	return fmt.Sprintf("🏆 %s — %s", metricLabel, scopeLabel)

}

func handleLeaderboardPagination(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := strings.Split(i.MessageComponentData().CustomID, ":")
	if len(data) != 4 {
		// Invalid format
		return
	}

	action := data[0]
	pageStr := data[1]
	scope := data[2]
	metric := data[3]

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Update page number
	switch action {
	case "lb_prev":
		page--
	case "lb_next":
		page++
	case "lb_goto":
	default:
		return
	}
	if page < 1 {
		page = 1
	}

	// Determine metric column
	var column string
	switch metric {
	case "success":
		column = "successful_meows"
	case "fail":
		column = "failed_meows"
	default:
		column = "total_meows"
	}

	// Set scope
	var guildID *string
	if scope == "guild" {
		id := i.GuildID
		guildID = &id
	}

	offset := (page - 1) * leaderboardPageSize

	entries, totalCount, err := db.GetLeaderboard(ctx, db.DB, guildID, column, leaderboardPageSize, offset)

	userRank, rankErr := db.GetUserRank(ctx, db.DB, i.Member.User.ID, guildID, column)

	if err != nil || len(entries) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ Couldn't load that page of the leaderboard.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	embed := formatLeaderboardEmbed(entries, scope, metric, page, totalCount, userRank, rankErr, i.Member.User.ID)
	components := renderLeaderboardButtons(scope, metric, page, totalCount)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}
