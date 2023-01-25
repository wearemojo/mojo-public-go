package datahappy

import (
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
	"github.com/cuvva/cuvva-public-go/lib/version"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
		client: jsonclient.NewClient(baseURL, &http.Client{
			Timeout: 5 * time.Second,

			Transport: otelhttp.NewTransport(nil),
		}),

		AuthToken: authToken,
	}
}
