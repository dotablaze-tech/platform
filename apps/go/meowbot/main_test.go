package main

import (
	"context"
	"errors"
	"libs/go/meowbot/util"
	"testing"
)

func TestRunWithInvalidBotToken(t *testing.T) {
	cfg := util.AppConfig{
		BotToken: "",
		ApiPort:  "0",
	}

	err := Run(context.Background(), cfg)
	if err == nil {
		t.Error("Expected error due to empty bot token, got nil")
	}
}

// Example for successful run test (would need deeper mocks/stubs)
func TestRunBootsUpBot(t *testing.T) {
	cfg := util.AppConfig{
		BotToken: "test_token",
		ApiPort:  "0",
		// Set up logger, etc.
	}

	// NOTE: You’d ideally use dependency injection to stub discordgo.New and db.InitDB
	// For now, this can be smoke-tested manually unless you abstract those away
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// simulate shutdown
		cancel()
	}()

	err := Run(ctx, cfg)
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Errorf("Unexpected error running bot: %v", err)
	}
}
