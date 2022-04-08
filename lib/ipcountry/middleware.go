package ipcountry

import (
	"net/http"
)

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			ipCountry := req.Header.Get("mojo-ip-country")

			ctx = SetIPCountry(ctx, ipCountry)
			req = req.WithContext(ctx)

			next.ServeHTTP(res, req)
		})
	}
}
