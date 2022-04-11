package discourse

import (
	"net/http"
)

type headerRoundtripper struct {
	header http.Header
}

func (h *headerRoundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header == nil {
		req.Header = http.Header{}
	}

	for k, v := range h.header {
		req.Header[k] = v
	}

	return http.DefaultTransport.RoundTrip(req)
}
