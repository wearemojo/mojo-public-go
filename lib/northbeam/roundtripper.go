package northbeam

import (
	"net/http"
)

type roundTripper struct {
	DataClientID string
	APIKey       string
}

func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	req = req.Clone(ctx)

	if req.Header == nil {
		req.Header = http.Header{}
	}

	req.Header.Set("Data-Client-Id", r.DataClientID)
	req.Header.Set("Authorization", r.APIKey)

	return http.DefaultTransport.RoundTrip(req)
}
