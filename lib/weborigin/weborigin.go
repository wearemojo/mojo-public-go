package weborigin

import (
	"fmt"
	"net/url"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

func ParseWebOrigin(url *url.URL) (string, error) {
	if url.Scheme != "http" && url.Scheme != "https" {
		return "", merr.New("invalid_scheme", nil)
	}

	if url.Hostname() == "" {
		return "", merr.New("invalid_hostname", nil)
	}

	baseURL := fmt.Sprintf("%s://%s", url.Scheme, url.Hostname())

	return baseURL, nil
}

func MustParseWebOrigin(in *url.URL) string {
	webOrigin, err := ParseWebOrigin(in)
	if err != nil {
		panic(err)
	}

	return webOrigin
}
