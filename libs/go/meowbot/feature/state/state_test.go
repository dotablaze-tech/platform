package state

import (
	"context"
	"errors"
	"libs/go/meowbot/feature/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to get a *string
func strPtr(s string) *string { return &s }

func TestDeref(t *testing.T) {
	assert.Equal(t, "", deref(nil))
	assert.Equal(t, "foo", deref(strPtr("foo")))
}

func TestGetOrCreate_CachesState(t *testing.T) {
	// ensure fresh store
	store = make(map[string]*GuildState)

	// stub out DB call that should never be called
	getGuildStreak = func(ctx context.Context, guildID string) (*db.GuildStreak, error) {
		t.Fatal("getGuildStreak should not be called when state already exists")
		return nil, nil
	}

	// pre-populate
	expected := &GuildState{MeowCount: 42, LastUserID: "u", HighScore: 7, HighScoreUserID: "u2"}
	store["g1"] = expected

	gs := GetOrCreate(context.Background(), "g1")
	assert.Same(t, expected, gs)
}

func TestGetOrCreate_LoadsFromDB(t *testing.T) {
	// clear store
	store = make(map[string]*GuildState)

	// stub DB call
	getGuildStreak = func(_ context.Context, guildID string) (*db.GuildStreak, error) {
		assert.Equal(t, "g2", guildID)
		return &db.GuildStreak{
			GuildID:         "g2",
			MeowCount:       5,
			LastUserID:      strPtr("u3"),
			HighScore:       10,
			HighScoreUserID: strPtr("u4"),
		}, nil
	}

	gs := GetOrCreate(context.Background(), "g2")
	assert.NotNil(t, gs)
	assert.Equal(t, 5, gs.MeowCount)
	assert.Equal(t, "u3", gs.LastUserID)
	assert.Equal(t, 10, gs.HighScore)
	assert.Equal(t, "u4", gs.HighScoreUserID)
}

func TestGetOrCreate_DBErrorFallsBack(t *testing.T) {
	store = make(map[string]*GuildState)

	getGuildStreak = func(_ context.Context, guildID string) (*db.GuildStreak, error) {
		return nil, errors.New("boom")
	}

	gs := GetOrCreate(context.Background(), "g3")
	assert.NotNil(t, gs)
	// fallback -> zero values
	assert.Equal(t, 0, gs.MeowCount)
	assert.Equal(t, "", gs.LastUserID)
}

func TestReset(t *testing.T) {
	store = make(map[string]*GuildState)
	store["g4"] = &GuildState{MeowCount: 9, LastUserID: "u5"}
	Reset("g4")
	gs := store["g4"]
	assert.Equal(t, 0, gs.MeowCount)
	assert.Equal(t, "", gs.LastUserID)
}
