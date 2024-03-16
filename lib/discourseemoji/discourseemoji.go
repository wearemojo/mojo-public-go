package discourseemoji

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/wearemojo/mojo-public-go/lib/slicefn"
)

// https://github.com/discourse/discourse/blob/d0d659e7330e1221f57cf6202a8de5c1604556f6/app/models/emoji.rb

var shortcodeRegex = regexp.MustCompile(`:([\w\-+]*(?::t\d)?):`)

var fitzpatrickScale = []rune{
	'\U0001F3FB', // type 1-2
	'\U0001F3FC', // type 3
	'\U0001F3FD', // type 4
	'\U0001F3FE', // type 5
	'\U0001F3FF', // type 6
}

//go:embed db.json
var dbJSON []byte
var db processedDB = processDB()

type processedDB struct {
	Emojis map[string]string
}

func processDB() processedDB {
	//nolint:tagliatelle // Discourse uses camel case
	var raw struct {
		Emojis []struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"emojis"`
		TonableEmojis []string `json:"tonableEmojis"`
	}
	if err := json.Unmarshal(dbJSON, &raw); err != nil {
		panic(err)
	}

	emojis := make(map[string]string, len(raw.Emojis))
	tonable := mapset.NewThreadUnsafeSet(raw.TonableEmojis...)

	for _, emoji := range raw.Emojis {
		base := slicefn.Map(strings.Split(emoji.Code, "-"), hexToRune)
		emojis[emoji.Name] = string(base)

		if tonable.Contains(emoji.Name) {
			for i, tone := range fitzpatrickScale {
				name := fmt.Sprintf("%s:t%d", emoji.Name, i+2)
				emojis[name] = string(append([]rune{base[0], tone}, base[1:]...))
			}
		}
	}

	return processedDB{
		Emojis: emojis,
	}
}

func hexToRune(hexStr string) rune {
	padded := strings.Repeat("0", 8-len(hexStr)) + hexStr
	decoded, err := hex.DecodeString(padded)
	if err != nil {
		panic(err)
	}

	return rune(decoded[0])<<24 |
		rune(decoded[1])<<16 |
		rune(decoded[2])<<8 |
		rune(decoded[3])
}

func ShortcodeToEmoji(shortcode string) string {
	return db.Emojis[shortcode]
}

func ExpandShortcodes(str string) string {
	return shortcodeRegex.ReplaceAllStringFunc(str, func(shortcode string) string {
		return ShortcodeToEmoji(shortcode[1 : len(shortcode)-1])
	})
}
