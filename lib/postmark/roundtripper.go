package postmark

import (
	"net/http"
)

type roundTripper struct {
	serverToken string
}

func (h roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header == nil {
		req.Header = http.Header{}
	}

	req.Header.Set("X-Postmark-Server-Token", h.serverToken)

	return http.DefaultTransport.RoundTrip(req)
}
