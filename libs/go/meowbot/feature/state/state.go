package state

import (
	"sync"
)

type GuildState struct {
	MeowCount     int
	LastUserID    string
	HighScore     int
	HighScoreUser string
}

var (
	mu    sync.Mutex
	store = make(map[string]*GuildState)
)

// GetOrCreate returns the guild state, creating it if needed.
func GetOrCreate(guildID string) *GuildState {
	mu.Lock()
	defer mu.Unlock()

	gs, ok := store[guildID]
	if !ok {
		gs = &GuildState{}
		store[guildID] = gs
	}
	return gs
}

// Reset clears state for a guild.
func Reset(guildID string) {
	mu.Lock()
	defer mu.Unlock()

	if gs, ok := store[guildID]; ok {
		gs.MeowCount = 0
		gs.LastUserID = ""
	}
}

func ResetAll(guildID string) {
	mu.Lock()
	defer mu.Unlock()
	store[guildID] = &GuildState{}
}
