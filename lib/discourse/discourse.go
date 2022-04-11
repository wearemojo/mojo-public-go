package discourse

import (
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

type Client struct {
	BaseURL string

	apiKey string
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,

		apiKey: apiKey,
	}
}

func (c *Client) client(header http.Header) *jsonclient.Client {
	return jsonclient.NewClient(c.BaseURL, &http.Client{
		Transport: &headerRoundtripper{
			header: header,
		},
	})
}

func (c *Client) usernameClient(username string) *jsonclient.Client {
	return c.client(http.Header{
		"Api-Key":      []string{c.apiKey},
		"Api-Username": []string{username},
	})
}

func (c *Client) systemClient() *jsonclient.Client {
	return c.usernameClient("system")
}
