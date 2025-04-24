package util

import (
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	t.Setenv("MODE", "")
	t.Setenv("DEBUG", "")
	t.Setenv("API_PORT", "")
	t.Setenv("WHITELISTED_GUILDS", "")
	t.Setenv("DISCORD_BOT_TOKEN", "test_token")
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("EMOJI_LIST", "😺,😸")

	cfg := LoadConfig()

	if cfg.Mode != "production" {
		t.Errorf("Expected mode=production, got %s", cfg.Mode)
	}
	if cfg.ApiPort != "8080" {
		t.Errorf("Expected default API_PORT=8080, got %s", cfg.ApiPort)
	}
	if !cfg.IsProd {
		t.Error("Expected IsProd to be true in production mode")
	}
	if cfg.EmojiList != "😺,😸" {
		t.Error("Emoji list not loaded correctly")
	}
}

func TestLoadConfig_WhitelistParsing(t *testing.T) {
	t.Setenv("WHITELISTED_GUILDS", "123,456,789")
	cfg := LoadConfig()

	expected := []string{"123", "456", "789"}
	if len(cfg.Whitelist.Guilds) != 3 {
		t.Errorf("Expected 3 whitelisted guilds, got %d", len(cfg.Whitelist.Guilds))
	}
	for i, id := range expected {
		if cfg.Whitelist.Guilds[i] != id {
			t.Errorf("Expected guild id %s, got %s", id, cfg.Whitelist.Guilds[i])
		}
	}
}

func TestIsAllowedGuild(t *testing.T) {
	cfg := AppConfig{
		IsProd: false,
		Whitelist: struct {
			Guilds []string
		}{
			Guilds: []string{"111", "222"},
		},
	}

	tests := []struct {
		guildID string
		allowed bool
	}{
		{"111", true},
		{"333", false},
	}

	for _, test := range tests {
		if cfg.IsAllowedGuild(test.guildID) != test.allowed {
			t.Errorf("Expected IsAllowedGuild(%s) = %v", test.guildID, test.allowed)
		}
	}

	// Production override
	cfg.IsProd = true
	if !cfg.IsAllowedGuild("anything") {
		t.Error("Expected IsAllowedGuild to return true in prod mode")
	}
}
