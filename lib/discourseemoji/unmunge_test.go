package discourseemoji

import (
	"net/url"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func urlMustParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func TestUnmungeCookedHTML(t *testing.T) {
	tests := []struct {
		name         string
		inputBaseURL *url.URL
		inputSource  string
		expected     string
	}{
		{
			"upside_down_face",
			nil,
			`<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">`,
			"🙃",
		},
		{
			"only-emoji end",
			nil,
			`<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="emoji only-emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">`,
			"🙃",
		},
		{
			"only-emoji start",
			nil,
			`<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="only-emoji emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">`,
			"🙃",
		},
		{
			"only-only-emoji",
			nil,
			`<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="only-emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20">`,
			`<img src="https://emoji.discourse-cdn.com/twitter/upside_down_face.png?v=12" title=":upside_down_face:" class="only-emoji" alt=":upside_down_face:" loading="lazy" width="20" height="20"/>`,
		},
		{
			"two paragraphs",
			nil,
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
			nil,
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
		{
			"urls",
			nil,
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
		{
			"preserves links without a base URL",
			nil,
			strings.TrimSpace(`
				<a href="https://example.com">example</a>
				<a href="/relative">relative</a>
				<a href="anchor">anchor</a>
				<a href="mailto:foo@example.com">email</a>
				<a href="?query">query</a>
				<a href="#fragment">fragment</a>
			`),
			strings.TrimSpace(`
				<a href="https://example.com">example</a>
				<a href="/relative">relative</a>
				<a href="anchor">anchor</a>
				<a href="mailto:foo@example.com">email</a>
				<a href="?query">query</a>
				<a href="#fragment">fragment</a>
			`),
		},
		{
			"updates links correctly with a base URL",
			urlMustParse("https://example.invalid/testing/foo"),
			strings.TrimSpace(`
				<a href="https://example.com">example</a>
				<a href="/relative">relative</a>
				<a href="anchor">anchor</a>
				<a href="mailto:foo@example.com">email</a>
				<a href="?query">query</a>
				<a href="#fragment">fragment</a>
			`),
			strings.TrimSpace(`
				<a href="https://example.com">example</a>
				<a href="https://example.invalid/relative">relative</a>
				<a href="https://example.invalid/testing/anchor">anchor</a>
				<a href="mailto:foo@example.com">email</a>
				<a href="https://example.invalid/testing/foo?query">query</a>
				<a href="https://example.invalid/testing/foo#fragment">fragment</a>
			`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)

			res, err := UnmungeCookedHTML(test.inputSource, test.inputBaseURL)

			is.NoErr(err)
			is.Equal(res, test.expected)
		})
	}
}
