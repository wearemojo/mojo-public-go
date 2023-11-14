package request

import (
	"net"
	"net/http"
	"strings"
)

// ClientIPHeader is the customer header, set by our network edge,
// that is expected to be the clients real IP address when X-Forwarded-For
// cannot be trusted.
const ClientIPHeader = `Infra-Client-Ip`

// ClientIP respects the infra-client-ip header containing the clients
// real IP address.
func ClientIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cci := strings.TrimSpace(r.Header.Get(ClientIPHeader))
		if cci != "" {
			// assert valid IP address
			ip := net.ParseIP(cci)
			if ip != nil {
				// marshal from parse string to restrict to supported encoding
				r.RemoteAddr = ip.String()
			}
		}

		next.ServeHTTP(w, r)
	})
}
