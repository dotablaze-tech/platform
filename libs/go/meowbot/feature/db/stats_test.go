package db

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestUpsertGuildChannel(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func(mockDB *sql.DB) {
		err := mockDB.Close()
		if err != nil {

		}
	}(mockDB)

	// Expect an exec with the right SQL and args
	mock.ExpectExec(regexp.QuoteMeta(`
        INSERT INTO guild_channels (guild_id, channel_id)
        VALUES ($1, $2)
        ON CONFLICT (guild_id) DO UPDATE SET
            channel_id = EXCLUDED.channel_id;
    `)).
		WithArgs("guild-123", "chan-456").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = UpsertGuildChannel(context.Background(), mockDB, "guild-123", "chan-456")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetChannelForGuild_NoRow(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		err := mockDB.Close()
		if err != nil {

		}
	}(mockDB)

	// Simulate no row
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT channel_id FROM guild_channels WHERE guild_id = $1;`)).
		WithArgs("guild-foo").
		WillReturnError(sql.ErrNoRows)

	cid, err := GetChannelForGuild(context.Background(), mockDB, "guild-foo")
	require.NoError(t, err)
	require.Equal(t, "", cid)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetChannelForGuild_Found(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		err := mockDB.Close()
		if err != nil {

		}
	}(mockDB)

	rows := sqlmock.NewRows([]string{"channel_id"}).AddRow("chan-789")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT channel_id FROM guild_channels WHERE guild_id = $1;`)).
		WithArgs("guild-foo").
		WillReturnRows(rows)

	cid, err := GetChannelForGuild(context.Background(), mockDB, "guild-foo")
	require.NoError(t, err)
	require.Equal(t, "chan-789", cid)
	require.NoError(t, mock.ExpectationsWereMet())
}
