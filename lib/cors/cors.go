package cors

import (
	"net/http"
)

func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := w.Header()

			header.Add("Vary", "Origin")

			if r.Header.Get("Origin") == "" {
				next.ServeHTTP(w, r)
				return
			}

			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS, POST, PUT, PATCH, DELETE")
			header.Set("Access-Control-Allow-Headers", "authorization, content-type")
			header.Set("Access-Control-Expose-Headers", "request-id")
			header.Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
