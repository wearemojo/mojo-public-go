package discourseemoji

import (
	"testing"

	"github.com/matryer/is"
)

func TestShortcodeToEmoji(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"smile", "😄"},
		{"smiley", "😃"},
		{"grinning", "😀"},
		{"blush", "😊"},
		{"woman_pilot", "👩‍✈️"},
		{"woman_pilot:t2", "👩🏻‍✈️"},
		{"woman_pilot:t3", "👩🏼‍✈️"},
		{"woman_pilot:t4", "👩🏽‍✈️"},
		{"woman_pilot:t5", "👩🏾‍✈️"},
		{"woman_pilot:t6", "👩🏿‍✈️"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)

			is.Equal(ShortcodeToEmoji(test.name), test.expected)
		})
	}
}

func TestExpandShortcodes(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"I am a :smile: human", "I am a 😄 human"},
		{":smile: :smiley: :grinning: :blush:", "😄 😃 😀 😊"},
		{"Fitzpatrick scale - :woman_pilot:t2::woman_pilot:t3: gap :woman_pilot:t4::woman_pilot:t5::woman_pilot:t6:", "Fitzpatrick scale - 👩🏻‍✈️👩🏼‍✈️ gap 👩🏽‍✈️👩🏾‍✈️👩🏿‍✈️"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)

			is.Equal(ExpandShortcodes(test.name), test.expected)
		})
	}
}
