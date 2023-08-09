package httpclient

import (
	"errors"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewClient(timeout time.Duration, transport http.RoundTripper) *http.Client {
	return &http.Client{
		Timeout:       timeout,
		CheckRedirect: CheckRedirect,
		Transport:     otelhttp.NewTransport(transport),
	}
}

// CheckRedirect provides a custom redirect policy that is more conservative/
// strict than the default.
//
// Normally Go follows redirects similarly to how a browser would, but that
// means sometimes changing the method (e.g. POST -> GET), dropping the body, or
// removing auth-related headers.
//
// The goal of this function is to ensure we only execute redirects that would
// be materially equivalent to the original request.
//
// Go's internal behavior is already set up to protect auth data, including
// cookies, by ensuring they can only be sent to the original host (or
// subdomains). If Go has seen fit to remove the auth data, we will consider the
// new request to be materially different, and stop following redirects.
func CheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		//nolint:forbidigo,goerr113 // ripped directly from stdlib
		return errors.New("stopped after 10 redirects")
	}

	// don't follow redirects if we can't send an equivalent request
	ireq := via[0]
	if req.Method != ireq.Method ||
		req.ContentLength != ireq.ContentLength ||
		// we currently only consider the Authorization header, as we don't use
		// cookies or any other headers that would be protected by Go's internal
		// behavior
		req.Header.Get("Authorization") != ireq.Header.Get("Authorization") {
		return http.ErrUseLastResponse
	}

	return nil
}
