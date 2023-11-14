package datahappy

import (
	"time"

	"github.com/wearemojo/mojo-public-go/lib/httpclient"
	"github.com/wearemojo/mojo-public-go/lib/jsonclient"
	"github.com/wearemojo/mojo-public-go/lib/version"
)

var library = &Library{
	Name:    "github.com/wearemojo/mojo-public-go/lib/datahappy",
	Version: version.Revision,
}

type Client struct {
	client *jsonclient.Client

	AuthToken string
}

func NewClient(baseURL, authToken string) *Client {
	return &Client{
		client: jsonclient.NewClient(
			baseURL,
			httpclient.NewClient(5*time.Second, nil),
		),

		AuthToken: authToken,
	}
}
