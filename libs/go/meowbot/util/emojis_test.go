// util/emojis_test.go
package util

import (
	"testing"
)

// helper to save and restore global state
func withEmojiConfig(emojiList string, fn func()) {
	origList := Cfg.EmojiList
	defer func() { Cfg.EmojiList = origList }()
	Cfg.EmojiList = emojiList
	fn()
}

func TestInitEmojis_WithCustomEmojis(t *testing.T) {
	withEmojiConfig("🙂,😀,😁,😂", func() {
		InitEmojis()
		// Now emojis should be exactly the four we set
		want := []string{"🙂", "😀", "😁", "😂"}
		if len(emojis) != len(want) {
			t.Fatalf("expected %d emojis, got %d", len(want), len(emojis))
		}
		for i, e := range want {
			if emojis[i] != e {
				t.Errorf("at index %d, want %q, got %q", i, e, emojis[i])
			}
		}
	})
}

func TestInitEmojis_WithDefaultEmojis(t *testing.T) {
	withEmojiConfig("", func() {
		InitEmojis()
		// Default list is exactly 10 long
		want := []string{"😺", "🐈", "🐾", "😹", "😼", "😻", "😽", "🐅", "🦁", "🐈‍⬛"}
		if len(emojis) != len(want) {
			t.Fatalf("expected %d default emojis, got %d", len(want), len(emojis))
		}
		for i, e := range want {
			if emojis[i] != e {
				t.Errorf("default[%d] = %q; want %q", i, emojis[i], e)
			}
		}
	})
}

func TestRandomEmoji_WithinList(t *testing.T) {
	// Force a small, known emoji list
	emojis = []string{"A", "B", "C", "D"}
	// Run RandomEmoji many times; it should never panic and always return one of A,B,C,D
	for i := 0; i < 100; i++ {
		got := RandomEmoji()
		found := false
		for _, e := range emojis {
			if e == got {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("RandomEmoji returned %q; want one of %v", got, emojis)
		}
	}
}
