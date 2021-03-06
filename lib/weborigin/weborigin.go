package weborigin

import (
	"context"
	"fmt"
	"net/url"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

// ParseWebOrigin takes a URL and returns a string to use to match web origins.
// The definition of the web origin is quite complex so you should refer to the
// test cases of this lib to see how it works.
func GetWebOrigin(ctx context.Context, url *url.URL) (string, error) {
	if url.Scheme != "http" && url.Scheme != "https" {
		return "", merr.New(ctx, "invalid_scheme", nil)
	}

	if url.Hostname() == "" {
		return "", merr.New(ctx, "invalid_hostname", nil)
	}

	webOrigin := fmt.Sprintf("%s://%s", url.Scheme, url.Hostname())

	return webOrigin, nil
}

func MustGetWebOrigin(ctx context.Context, in *url.URL) string {
	webOrigin, err := GetWebOrigin(ctx, in)
	if err != nil {
		panic(err)
	}

	return webOrigin
}
