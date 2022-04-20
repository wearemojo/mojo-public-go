package weborigin

import (
	"net/url"
	"testing"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

var tests = []struct {
	Name     string
	URL      string
	Expected string
	Err      merr.Code
}{
	{
		Name:     "with trailing slash",
		URL:      "https://example.com/",
		Expected: "https://example.com",
	},
	{
		Name:     "with path",
		URL:      "https://example.com/foo/bar",
		Expected: "https://example.com",
	},
	{
		Name:     "with path and port",
		URL:      "https://google.com:8443/bar",
		Expected: "https://google.com",
	},
	{
		Name:     "with http and path",
		URL:      "http://google.com/foo",
		Expected: "http://google.com",
	},
	{
		Name:     "it is already perfect",
		URL:      "https://app.mojo.so",
		Expected: "https://app.mojo.so",
	},
	{
		Name:     "gibberish",
		URL:      "alskdjsdkghkjdfg",
		Expected: "",
		Err:      merr.Code("invalid_scheme"),
	},
	{
		Name:     "missing scheme",
		URL:      "google.com",
		Expected: "",
		Err:      merr.Code("invalid_scheme"),
	},
	{
		Name:     "missing hostname",
		URL:      "https://",
		Expected: "",
		Err:      merr.Code("invalid_hostname"),
	},
}

func TestParseWebOrigin(t *testing.T) {
	for _, test := range tests {
		test := test

		t.Run(test.Name, func(t *testing.T) {
			url, err := url.Parse(test.URL)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			actual, err := ParseWebOrigin(url)
			if err != nil && test.Err != "" {
				if merr.IsCode(err, test.Err) {
					return
				}

				t.Errorf("unexpected error: %v", err)
			}
			if actual != test.Expected {
				t.Errorf("%s, expected %s, got %s", test.Name, test.Expected, actual)
			}
		})
	}
}
