package util

import (
	"math/rand"
	"strings"
)

var emojis []string

func InitEmojis() {
	if emojiEnv := Cfg.EmojiList; emojiEnv != "" {
		emojis = strings.Split(emojiEnv, ",")
		Cfg.Logger.Info("✨ Using custom emojis", "emojis", emojis)
	} else {
		emojis = []string{"😺", "🐈", "🐾", "😹", "😼", "😻", "😽", "🐅", "🦁", "🐈‍⬛"}
		Cfg.Logger.Info("🐾 Using default emojis", "emojis", emojis)
	}
}

func RandomEmoji() string {
	return emojis[rand.Intn(len(emojis))]
}
