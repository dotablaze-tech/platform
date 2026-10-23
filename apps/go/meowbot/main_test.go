package main

import (
	"context"
	"libs/go/meowbot/util"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestRunHealthcheck(t *testing.T) {
	status := http.StatusOK
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/liveness" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(status)
	}))
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	port := u.Port()

	if err := runHealthcheck(context.Background(), port); err != nil {
		t.Fatalf("expected healthy, got %v", err)
	}
	status = http.StatusServiceUnavailable
	if err := runHealthcheck(context.Background(), port); err == nil {
		t.Fatal("expected error for non-200")
	}
}

func TestRunHealthcheckConnectionRefused(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	u, _ := url.Parse(srv.URL)
	srv.Close()
	if err := runHealthcheck(context.Background(), u.Port()); err == nil {
		t.Fatal("expected error when nothing is listening")
	}
}
