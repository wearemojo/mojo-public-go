package discourse

import (
	"maps"
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

	maps.Copy(req.Header, r.header)

	return http.DefaultTransport.RoundTrip(req)
}
