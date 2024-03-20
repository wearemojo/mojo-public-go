package discourse

import (
	"net/http"
	"net/url"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/httpclient"
	"github.com/wearemojo/mojo-public-go/lib/jsonclient"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

const ErrEmptyParam = merr.Code("empty_param")

type Client struct {
	BaseURL *url.URL

	apiKey string
}

func NewClient(baseURL *url.URL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,

		apiKey: apiKey,
	}
}

type IdentifiedClient struct {
	client *jsonclient.Client
}

func (c *Client) identifiedClient(header http.Header) *IdentifiedClient {
	return &IdentifiedClient{
		client: jsonclient.NewClient(
			c.BaseURL.String(),
			httpclient.NewClient(10*time.Second, roundTripper{header}),
		),
	}
}

func (c *Client) AsUsername(username string) *IdentifiedClient {
	return c.identifiedClient(http.Header{
		"Api-Key":      []string{c.apiKey},
		"Api-Username": []string{username},
	})
}

func (c *Client) AsSystem() *IdentifiedClient {
	return c.AsUsername("system")
}
