package ipcountry

import (
	"net/http"
)

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ipCountry := r.Header.Get("mojo-ip-country")

			ctx = SetIPCountry(ctx, ipCountry)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
