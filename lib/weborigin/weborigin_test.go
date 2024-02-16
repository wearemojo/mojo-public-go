package weborigin

import (
	"context"
	"net/url"
	"testing"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

var baseTests = []struct {
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

func TestGetOrigin(t *testing.T) {
	tests := append([]struct {
		Name     string
		URL      string
		Expected string
		Err      merr.Code
	}{
		{
			Name:     "no delimiters",
			URL:      "blah-foo",
			Expected: "blah-foo://",
			Err:      merr.Code("invalid_scheme"),
		},
		{
			Name:     "colon only",
			URL:      "blah-foo:",
			Expected: "blah-foo://",
		},
		{
			Name:     "one slash",
			URL:      "blah-foo:/",
			Expected: "blah-foo://",
		},
		{
			Name:     "two slashes",
			URL:      "blah-foo://",
			Expected: "blah-foo://",
		},
		{
			Name:     "two slashes and more",
			URL:      "blah-foo://foo/bar/baz",
			Expected: "blah-foo://foo",
		},
		{
			Name:     "two slashes and more with port",
			URL:      "blah-foo://foo:8080/bar/baz",
			Expected: "blah-foo://foo",
		},
		{
			Name:     "three slashes",
			URL:      "blah-foo:///",
			Expected: "blah-foo://",
		},
		{
			Name:     "three slashes and more",
			URL:      "blah-foo:///foo/bar/baz",
			Expected: "blah-foo://",
		},
		{
			Name:     "three slashes and more with port",
			URL:      "blah-foo://:8080/foo/bar/baz",
			Expected: "blah-foo://",
		},
		{
			Name:     "four slashes",
			URL:      "blah-foo:////",
			Expected: "blah-foo://",
		},
	}, baseTests...)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			url, err := url.Parse(test.URL)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			actual, err := GetOrigin(context.Background(), url)
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

func TestGetWebOrigin(t *testing.T) {
	for _, test := range baseTests {
		t.Run(test.Name, func(t *testing.T) {
			url, err := url.Parse(test.URL)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			actual, err := GetWebOrigin(context.Background(), url)
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
