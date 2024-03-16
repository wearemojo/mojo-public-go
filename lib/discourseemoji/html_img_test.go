package discourseemoji

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestReplaceHTMLImagesWithEmojis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"upside_down_face",
			`<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">`,
			"🙃",
		},
		{
			"only-emoji end",
			`<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="emoji only-emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">`,
			"🙃",
		},
		{
			"only-emoji start",
			`<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="only-emoji emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">`,
			"🙃",
		},
		{
			"only-only-emoji",
			`<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="only-emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">`,
			`<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="only-emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20"/>`,
		},
		{
			"two paragraphs",
			strings.TrimSpace(`
				<p>
					<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/twitter/cold_face.png?v=12" title=":cold_face:" class="emoji" alt=":cold_face:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/twitter/raised_back_of_hand/2.png?v=12" title=":raised_back_of_hand:t2:" class="emoji" alt=":raised_back_of_hand:t2:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/twitter/vulcan_salute/3.png?v=12" title=":vulcan_salute:t3:" class="emoji" alt=":vulcan_salute:t3:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/twitter/vulcan_salute.png?v=12" title=":vulcan_salute:" class="emoji" alt=":vulcan_salute:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/twitter/vulcan_salute/6.png?v=12" title=":vulcan_salute:t6:" class="emoji" alt=":vulcan_salute:t6:" loading="lazy" width="20" height="20">
				</p>
				<p>
					<img src="https://emoji.discourse-cdn.com/apple/upside_down_face.png?v=12" title=":upside_down_face:" class="emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/apple/cold_face.png?v=12" title=":cold_face:" class="emoji" alt=":cold_face:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/apple/raised_back_of_hand/2.png?v=12" title=":raised_back_of_hand:t2:" class="emoji" alt=":raised_back_of_hand:t2:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/apple/vulcan_salute/3.png?v=12" title=":vulcan_salute:t3:" class="emoji" alt=":vulcan_salute:t3:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/apple/vulcan_salute.png?v=12" title=":vulcan_salute:" class="emoji" alt=":vulcan_salute:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/apple/vulcan_salute/6.png?v=12" title=":vulcan_salute:t6:" class="emoji" alt=":vulcan_salute:t6:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/apple/vulcan_salute.png?v=12" title=":vulcan_salute:" class="emoji" alt=":vulcan_salute:" loading="lazy" width="20" height="20">
				</p>
			`),
			strings.TrimSpace(`
				<p>
					🙃
					🥶
					🤚🏻
					🖖🏼
					🖖
					🖖🏿
				</p>
				<p>
					🙃
					🥶
					🤚🏻
					🖖🏼
					🖖
					🖖🏿
					🖖
				</p>
			`),
		},
		{
			"unrecognized",
			strings.TrimSpace(`
				<p>
					<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/blah/raised_back_of_hand/2.png?v=12" title=":blah:t2:" class="emoji" alt=":blah:t2:" loading="lazy" width="20" height="20">
					<img src="https://emoji.discourse-cdn.com/apple/cold_face.png?v=12" title=":cold_face:" class="emoji" alt=":cold_face:" loading="lazy" width="20" height="20">
				</p>
			`),
			strings.TrimSpace(`
				<p>
					🙃
					<img src="https://emoji.discourse-cdn.com/blah/raised_back_of_hand/2.png?v=12" title=":blah:t2:" class="emoji" alt=":blah:t2:" loading="lazy" width="20" height="20"/>
					🥶
				</p>
			`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)

			res, err := ReplaceHTMLImagesWithEmojis(test.input)

			is.NoErr(err)
			is.Equal(res, test.expected)
		})
	}
}
