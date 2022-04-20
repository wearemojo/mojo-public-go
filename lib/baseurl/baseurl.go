package baseurl

import (
	"fmt"
	"net/url"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

func ParseBaseURL(in string) (string, error) {
	url, err := url.Parse(in)
	if err != nil {
		return "", err
	}

	if url.Scheme == "" || url.Host == "" {
		return "", merr.New("invalid_url", nil)
	}

	baseURL := fmt.Sprintf("%s://%s", url.Scheme, url.Hostname())

	return baseURL, nil
}

func MustParseBaseURL(in string) string {
	baseURL, err := ParseBaseURL(in)
	if err != nil {
		panic(err)
	}

	return baseURL
}
