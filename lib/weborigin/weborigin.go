package weborigin

import (
	"context"
	"fmt"
	"net/url"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

// GetOrigin is like GetWebOrigin but it also accepts non-web URLs, like custom
// URL schemes for apps.
func GetOrigin(ctx context.Context, url *url.URL) (string, error) {
	if url.Scheme == "" {
		return "", merr.New(ctx, "invalid_scheme", nil)
	}

	if url.Scheme == "http" || url.Scheme == "https" {
		if url.Hostname() == "" {
			return "", merr.New(ctx, "invalid_hostname", nil)
		}
	}

	return fmt.Sprintf("%s://%s", url.Scheme, url.Hostname()), nil
}

// GetWebOrigin takes a URL and returns a string to use to match web origins.
// The definition of the web origin is quite complex so you should refer to the
// test cases of this lib to see how it works.
func GetWebOrigin(ctx context.Context, url *url.URL) (string, error) {
	if url.Scheme != "http" && url.Scheme != "https" {
		return "", merr.New(ctx, "invalid_scheme", nil)
	}

	return GetOrigin(ctx, url)
}

func MustGetWebOrigin(ctx context.Context, in *url.URL) string {
	webOrigin, err := GetWebOrigin(ctx, in)
	if err != nil {
		panic(err)
	}
	return webOrigin
}
