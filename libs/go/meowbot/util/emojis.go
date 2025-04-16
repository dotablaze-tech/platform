package util

import (
	"log/slog"
	"math/rand"
	"os"
	"strings"
)

var emojis []string

func InitEmojis(logger *slog.Logger) {
	if emojiEnv := os.Getenv("EMOJI_LIST"); emojiEnv != "" {
		emojis = strings.Split(emojiEnv, ",")
		logger.Info("Using custom emojis", "emojis", emojis)
	} else {
		emojis = []string{"😺", "🐈", "🐾", "😹", "😼", "😻"}
		logger.Info("Using default emojis", "emojis", emojis)
	}
}

func RandomEmoji() string {
	return emojis[rand.Intn(len(emojis))]
}
