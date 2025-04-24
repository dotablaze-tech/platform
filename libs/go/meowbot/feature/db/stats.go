package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func UpsertUser(ctx context.Context, db *sql.DB, user User) error {
	query := `
		INSERT INTO users (id, username)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET username = EXCLUDED.username;
	`
	_, err := db.ExecContext(ctx, query, user.ID, user.Username)
	return err
}

func UpsertGuild(ctx context.Context, db *sql.DB, guild Guild) error {
	query := `
		INSERT INTO guilds (id)
		VALUES ($1)
		ON CONFLICT (id) DO NOTHING;
	`

	_, err := db.ExecContext(ctx, query, guild.ID)
	return err
}

func UpsertGuildChannel(ctx context.Context, db *sql.DB, guildID, channelID string) error {
	query := `
		INSERT INTO guild_channels (guild_id, channel_id)
		VALUES ($1, $2)
		ON CONFLICT (guild_id) DO UPDATE SET
			channel_id = EXCLUDED.channel_id;
	`

	_, err := db.ExecContext(ctx, query, guildID, channelID)
	if err != nil {
		return fmt.Errorf("failed to upsert guild channel: %w", err)
	}
	return nil
}

func GetChannelForGuild(ctx context.Context, db *sql.DB, guildID string) (string, error) {
	query := `SELECT channel_id FROM guild_channels WHERE guild_id = $1;`

	var channelID string
	err := db.QueryRowContext(ctx, query, guildID).Scan(&channelID)

	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get guild channel: %w", err)
	}
	return channelID, nil
}

func IncrementMeow(ctx context.Context, db *sql.DB, guildID, userID string, success bool, now time.Time) error {
	successQuery := `
			INSERT INTO user_guild_stats (guild_id, user_id, successful_meows, total_meows, current_streak, highest_streak, last_meow_at)
			VALUES ($1, $2, 1, 1, 1, 1, $3)
			ON CONFLICT (guild_id, user_id)
			DO UPDATE SET
				successful_meows = user_guild_stats.successful_meows + 1,
				total_meows = user_guild_stats.total_meows + 1,
				current_streak = user_guild_stats.current_streak + 1,
				highest_streak = GREATEST(user_guild_stats.highest_streak, user_guild_stats.current_streak + 1),
				last_meow_at = $3;
	`
	failureQuery := `
		INSERT INTO user_guild_stats (guild_id, user_id, failed_meows, total_meows, current_streak, last_failed_meow_at)
		VALUES ($1, $2, 1, 1, 0, $3)
		ON CONFLICT (guild_id, user_id)
		DO UPDATE SET
			failed_meows = user_guild_stats.failed_meows + 1,
			total_meows = user_guild_stats.total_meows + 1,
			current_streak = 0,
			last_failed_meow_at = $3;
		`

	if success {
		_, err := db.ExecContext(ctx, successQuery, guildID, userID, now)
		return err
	}

	_, err := db.ExecContext(ctx, failureQuery, guildID, userID, now)
	return err
}

func GetGuildStreak(ctx context.Context, db *sql.DB, guildID string) (*GuildStreak, error) {
	query := `
		SELECT guild_id, meow_count, last_user_id, high_score, high_score_user_id
		FROM guild_streaks
		WHERE guild_id = $1;
	`

	row := db.QueryRowContext(ctx, query, guildID)

	var gs GuildStreak
	err := row.Scan(&gs.GuildID, &gs.MeowCount, &gs.LastUserID, &gs.HighScore, &gs.HighScoreUserID)
	if err != nil {
		return nil, err
	}
	return &gs, nil
}

func UpsertGuildStreak(ctx context.Context, db *sql.DB, streak GuildStreak) error {
	query := `
		INSERT INTO guild_streaks (guild_id, meow_count, last_user_id, high_score, high_score_user_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (guild_id) DO UPDATE SET
			meow_count = EXCLUDED.meow_count,
			last_user_id = EXCLUDED.last_user_id,
			high_score = EXCLUDED.high_score,
			high_score_user_id = EXCLUDED.high_score_user_id;
	`
	_, err := db.ExecContext(
		ctx,
		query,
		streak.GuildID,
		streak.MeowCount,
		streak.LastUserID,
		streak.HighScore,
		streak.HighScoreUserID,
	)
	return err
}

func GetUserStats(
	ctx context.Context,
	db *sql.DB,
	guildID *string, // nil means global
	userID string,
) (UserGuildStats, error) {
	var (
		query string
		args  []any
	)

	// If a specific guild, just pull that row.
	if guildID != nil {
		query = `
            SELECT
                guild_id,
                user_id,
                successful_meows,
                failed_meows,
                total_meows,
                current_streak,
                highest_streak,
                last_meow_at,
                last_failed_meow_at
            FROM user_guild_stats
            WHERE guild_id = $1 AND user_id = $2
        `
		args = []any{*guildID, userID}

		// Otherwise aggregate globally:
	} else {
		query = `
            SELECT
                '' AS guild_id,
                user_id,
                COALESCE(SUM(successful_meows),0)    AS successful_meows,
                COALESCE(SUM(failed_meows),0)        AS failed_meows,
                COALESCE(SUM(total_meows),0)         AS total_meows,
                COALESCE(MAX(current_streak),0)      AS current_streak,
                COALESCE(MAX(highest_streak),0)      AS highest_streak,
                MAX(last_meow_at)                    AS last_meow_at,
                MAX(last_failed_meow_at)             AS last_failed_meow_at
            FROM user_guild_stats
            WHERE user_id = $1
            GROUP BY user_id
        `
		args = []any{userID}
	}

	var stats UserGuildStats
	err := db.QueryRowContext(ctx, query, args...).Scan(
		&stats.GuildID,
		&stats.UserID,
		&stats.SuccessfulMeows,
		&stats.FailedMeows,
		&stats.TotalMeows,
		&stats.CurrentStreak,
		&stats.HighestStreak,
		&stats.LastMeowAt,
		&stats.LastFailedMeowAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) && guildID != nil {
			// no entry for this guild/user yet
			return UserGuildStats{
				GuildID:          *guildID,
				UserID:           userID,
				SuccessfulMeows:  0,
				FailedMeows:      0,
				TotalMeows:       0,
				CurrentStreak:    0,
				HighestStreak:    0,
				LastMeowAt:       nil,
				LastFailedMeowAt: nil,
			}, nil
		}
		return UserGuildStats{}, fmt.Errorf("GetUserStats: %w", err)
	}
	return stats, nil
}

func GetGlobalStats(ctx context.Context, db *sql.DB) (*GlobalStats, error) {
	guildsQuery := `SELECT COUNT(*) FROM guilds`
	usersQuery := `SELECT COUNT(*) FROM users`
	usersGuildStatsQuery := `SELECT COALESCE(SUM(total_meows), 0) FROM user_guild_stats`
	var stats GlobalStats

	err := db.QueryRowContext(ctx, guildsQuery).Scan(&stats.TotalGuilds)
	if err != nil {
		return nil, err
	}

	err = db.QueryRowContext(ctx, usersQuery).Scan(&stats.TotalUsers)
	if err != nil {
		return nil, err
	}

	err = db.QueryRowContext(ctx, usersGuildStatsQuery).Scan(&stats.TotalMeows)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func GetUserGlobalStats(ctx context.Context, db *sql.DB, userID string) (UserGlobalStats, error) {
	query := `
		SELECT 
			u.id,
			u.username,
			u.created_at,
			COALESCE(SUM(ugs.successful_meows), 0),
			COALESCE(SUM(ugs.failed_meows), 0),
			COALESCE(SUM(ugs.total_meows), 0),
			COALESCE(MAX(ugs.highest_streak), 0)
		FROM users u
		LEFT JOIN user_guild_stats ugs ON u.id = ugs.user_id
		WHERE u.id = $1
		GROUP BY u.id, u.username, u.created_at
	`

	var res UserGlobalStats
	err := db.QueryRowContext(ctx, query, userID).Scan(
		&res.UserID,
		&res.Username,
		&res.CreatedAt,
		&res.SuccessfulMeows,
		&res.FailedMeows,
		&res.TotalMeows,
		&res.HighestStreak,
	)
	if err != nil {
		return UserGlobalStats{}, err
	}

	return res, nil
}

func GetLeaderboard3(ctx context.Context, db *sql.DB, limit int) (entries []LeaderboardEntry, err error) {
	query := `
		SELECT u.id, u.username, u.created_at, SUM(ugs.total_meows) as total
		FROM user_guild_stats ugs
		JOIN users u ON ugs.user_id = u.id
		GROUP BY u.id, u.username, u.created_at
		ORDER BY total DESC
		LIMIT $1;
	`

	rows, err := db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("query leaderboard: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close rows: %w", closeErr)
		}
	}()

	for rows.Next() {
		var user User
		var entry LeaderboardEntry

		err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt, &entry.TotalMeows)
		if err != nil {
			return nil, fmt.Errorf("scan leaderboard row: %w", err)
		}

		entry.User = &user
		entries = append(entries, entry)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("iterate leaderboard rows: %w", rowsErr)
	}

	return entries, nil
}

func GetGuildStats(ctx context.Context, db *sql.DB, guildID string) (*GuildStats, error) {
	query := `
	SELECT
		g.id, g.created_at,
		gs.meow_count,
		lu.id, lu.username, lu.created_at,
		gs.high_score,
		hu.id, hu.username, hu.created_at,
		COALESCE(SUM(ugs.total_meows), 0) as total_meows,
		COALESCE(SUM(ugs.successful_meows), 0) as successful_meows,
		COALESCE(SUM(ugs.failed_meows), 0) as failed_meows
	FROM guild_streaks gs
	JOIN guilds g ON g.id = gs.guild_id
	LEFT JOIN users lu ON lu.id = gs.last_user_id
	LEFT JOIN users hu ON hu.id = gs.high_score_user_id
	LEFT JOIN user_guild_stats ugs ON gs.guild_id = ugs.guild_id
	WHERE gs.guild_id = $1
	GROUP BY g.id, g.created_at, gs.meow_count,
	         lu.id, lu.username, lu.created_at,
	         gs.high_score,
	         hu.id, hu.username, hu.created_at;
	`

	row := db.QueryRowContext(ctx, query, guildID)

	var stats GuildStats
	var guild Guild
	var lastUser, highScoreUser User
	var lastUserID, highScoreUserID sql.NullString
	var lastUsername, highScoreUsername sql.NullString
	var lastCreatedAt, highScoreCreatedAt sql.NullTime

	err := row.Scan(
		&guild.ID, &guild.CreatedAt,
		&stats.CurrentStreak,
		&lastUserID, &lastUsername, &lastCreatedAt,
		&stats.HighScore,
		&highScoreUserID, &highScoreUsername, &highScoreCreatedAt,
		&stats.TotalMeows,
		&stats.SuccessfulMeows,
		&stats.FailedMeows,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan guild stats: %w", err)
	}

	stats.Guild = &guild

	if lastUserID.Valid && lastUsername.Valid && lastCreatedAt.Valid {
		lastUser = User{
			ID:        lastUserID.String,
			Username:  lastUsername.String,
			CreatedAt: lastCreatedAt.Time,
		}
		stats.LastUser = &lastUser
	}

	if highScoreUserID.Valid && highScoreUsername.Valid && highScoreCreatedAt.Valid {
		highScoreUser = User{
			ID:        highScoreUserID.String,
			Username:  highScoreUsername.String,
			CreatedAt: highScoreCreatedAt.Time,
		}
		stats.HighScoreUser = &highScoreUser
	}

	return &stats, nil
}

func GetAllUsers(ctx context.Context, db *sql.DB) ([]*User, error) {
	query := `SELECT id, username, created_at FROM users`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close rows: %w", closeErr)
		}
	}()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &u)
	}
	return users, err
}

func GetAllGuilds(ctx context.Context, db *sql.DB) ([]*Guild, error) {
	query := `SELECT id, created_at FROM guilds`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query guilds: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close rows: %w", closeErr)
		}
	}()

	var guilds []*Guild
	for rows.Next() {
		var g Guild
		if err := rows.Scan(&g.ID, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan guild: %w", err)
		}
		guilds = append(guilds, &g)
	}
	return guilds, rows.Err()
}

func GetUserPerGuildStats(ctx context.Context, db *sql.DB, userID string) ([]UserGuildStats, error) {
	query := `
		SELECT 
			guild_id,
			successful_meows,
			failed_meows,
			total_meows,
			current_streak,
			highest_streak,
			last_meow_at,
			last_failed_meow_at
		FROM user_guild_stats
		WHERE user_id = $1
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close rows: %w", closeErr)
		}
	}()

	var stats []UserGuildStats
	for rows.Next() {
		var s UserGuildStats
		err := rows.Scan(
			&s.GuildID,
			&s.SuccessfulMeows,
			&s.FailedMeows,
			&s.TotalMeows,
			&s.CurrentStreak,
			&s.HighestStreak,
			&s.LastMeowAt,
			&s.LastFailedMeowAt,
		)
		if err != nil {
			return nil, err
		}
		s.UserID = userID
		stats = append(stats, s)
	}

	return stats, nil
}

func GetLeaderboard(
	ctx context.Context,
	db *sql.DB,
	guildID *string, // nil means global
	metric string, // "total_meows", "successful_meows", or "failed_meows"
	limit, offset int,
) ([]LeaderboardEntry, int, error) {
	var (
		query      string
		countQuery string
		queryArgs  []any
		countArgs  []any
	)

	// Use aggregation
	baseSelect := fmt.Sprintf(`
		SELECT u.id, u.username, u.created_at, SUM(ugs.%s) AS value
		FROM user_guild_stats ugs
		JOIN users u ON u.id = ugs.user_id
	`, metric)

	baseGroupOrder := ` GROUP BY u.id ORDER BY value DESC LIMIT %d OFFSET %d`
	baseCount := `SELECT COUNT(*) FROM (
		SELECT ugs.user_id
		FROM user_guild_stats ugs
		%s
		GROUP BY ugs.user_id
	) AS subquery`

	if guildID != nil {
		whereClause := `WHERE ugs.guild_id = $1`
		query = baseSelect + " " + whereClause + fmt.Sprintf(baseGroupOrder, limit, offset)
		countQuery = fmt.Sprintf(baseCount, whereClause)
		queryArgs = []any{*guildID}
		countArgs = []any{*guildID}
	} else {
		query = baseSelect + fmt.Sprintf(baseGroupOrder, limit, offset)
		countQuery = fmt.Sprintf(baseCount, "")
		queryArgs = []any{}
		countArgs = []any{}
	}

	rows, err := db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query leaderboard: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close rows: %w", closeErr)
		}
	}()

	var entries []LeaderboardEntry
	for rows.Next() {
		var user User
		var entry LeaderboardEntry
		var value int

		if err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt, &value); err != nil {
			return nil, 0, fmt.Errorf("scan leaderboard row: %w", err)
		}

		entry.User = &user
		switch metric {
		case "total_meows":
			entry.TotalMeows = value
		case "successful_meows":
			entry.SuccessfulMeows = value
		case "failed_meows":
			entry.FailedMeows = value
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	var total int
	err = db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count leaderboard: %w", err)
	}

	return entries, total, nil
}

func GetUserRank(ctx context.Context, db *sql.DB, userID string, guildID *string, column string) (int, error) {
	// Validate column
	validColumns := map[string]bool{
		"total_meows":      true,
		"successful_meows": true,
		"failed_meows":     true,
		"highest_streak":   true,
		"current_streak":   true,
	}
	if !validColumns[column] {
		return 0, fmt.Errorf("invalid column name: %s", column)
	}

	var (
		query string
		args  []any
	)

	if guildID != nil {
		query = fmt.Sprintf(`
			SELECT rank FROM (
				SELECT user_id, RANK() OVER (ORDER BY %s DESC) AS rank
				FROM user_guild_stats
				WHERE guild_id = $1
			) ranked WHERE user_id = $2
		`, column)
		args = []any{*guildID, userID}
	} else {
		query = fmt.Sprintf(`
			SELECT rank FROM (
				SELECT user_id, RANK() OVER (ORDER BY SUM(%s) DESC) AS rank
				FROM user_guild_stats
				GROUP BY user_id
			) ranked WHERE user_id = $1
		`, column)
		args = []any{userID}
	}

	var rank int
	err := db.QueryRowContext(ctx, query, args...).Scan(&rank)
	return rank, err
}
