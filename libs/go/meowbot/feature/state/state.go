package state

import (
	"context"
	"libs/go/meowbot/feature/db"
	"sync"
)

type GuildState struct {
	MeowCount       int
	LastUserID      string
	HighScore       int
	HighScoreUserID string
}

var (
	getGuildStreak = func(ctx context.Context, guildID string) (*db.GuildStreak, error) {
		return db.GetGuildStreak(ctx, db.DB, guildID)
	}

	mu    sync.Mutex
	store = make(map[string]*GuildState)
)

func GetOrCreate(ctx context.Context, guildID string) *GuildState {
	mu.Lock()
	defer mu.Unlock()

	if gs, ok := store[guildID]; ok {
		return gs
	}

	dbStreak, err := getGuildStreak(ctx, guildID)
	if err != nil {
		dbStreak = &db.GuildStreak{} // fallback
	}

	gs := &GuildState{
		MeowCount:       dbStreak.MeowCount,
		LastUserID:      deref(dbStreak.LastUserID),
		HighScore:       dbStreak.HighScore,
		HighScoreUserID: deref(dbStreak.HighScoreUserID),
	}
	store[guildID] = gs
	return gs
}

func Reset(guildID string) {
	mu.Lock()
	defer mu.Unlock()
	if gs, ok := store[guildID]; ok {
		gs.MeowCount = 0
		gs.LastUserID = ""
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
