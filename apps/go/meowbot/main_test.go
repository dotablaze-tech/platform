package main

import (
	"context"
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
	t.Skip("Skipping integration test: no local DB running")
}
