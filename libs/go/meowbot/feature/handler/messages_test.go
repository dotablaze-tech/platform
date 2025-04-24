package handler

import (
	"strings"
	"testing"
)

func TestMeowRegex(t *testing.T) {
	cases := []struct {
		input  string
		should bool
	}{
		{"meow", true},
		{"MEOW", true},
		{"MEEEOOOWWWW", true},
		{"mEoW", true},
		{" meow  ", true}, // whitespace is trimmed in your handler
		{"meow1", false},
		{"m e o w", false},
		{"moo", false},
		{"woof", false},
		{"", false},
	}

	for _, c := range cases {
		// simulate your handler's normalization
		norm := strings.ToLower(strings.TrimSpace(c.input))
		matched := meowRegex.MatchString(norm)
		if matched != c.should {
			t.Errorf("meowRegex.MatchString(%q) = %v, want %v", c.input, matched, c.should)
		}
	}
}
