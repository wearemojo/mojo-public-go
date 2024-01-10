package revenuecat

import (
	"fmt"
	"net/http"

	"github.com/wearemojo/mojo-public-go/lib/secret"
)

type roundTripper struct {
	SecretID string
}

func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	req = req.Clone(ctx)

	apiKey, err := secret.Get(ctx, r.SecretID)
	if err != nil {
		return nil, err
	}

	if req.Header == nil {
		req.Header = http.Header{}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	return http.DefaultTransport.RoundTrip(req)
}
