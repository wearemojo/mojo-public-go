package discourse

import (
	"net/http"
)

type roundTripper struct {
	header http.Header
}

func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())

	if req.Header == nil {
		req.Header = http.Header{}
	}

	for k, v := range r.header {
		req.Header[k] = v
	}

	return http.DefaultTransport.RoundTrip(req)
}
