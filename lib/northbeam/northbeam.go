package northbeam

import (
	"time"

	"github.com/wearemojo/mojo-public-go/lib/httpclient"
	"github.com/wearemojo/mojo-public-go/lib/jsonclient"
)

type Client struct {
	client *jsonclient.Client
}

func NewClient(baseURL, dataClientID, apiKey string) *Client {
	return &Client{
		client: jsonclient.NewClient(
			baseURL,
			httpclient.NewClient(5*time.Second, roundTripper{
				DataClientID: dataClientID,
				APIKey:       apiKey,
			}),
		),
	}
}
