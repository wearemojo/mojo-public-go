package cors

import (
	"net/http"
)

func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			header := res.Header()

			header.Add("Vary", "Origin")

			if req.Header.Get("Origin") == "" {
				next.ServeHTTP(res, req)
				return
			}

			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS, POST, PUT, PATCH, DELETE")
			header.Set("Access-Control-Allow-Headers", "authorization, content-type")
			header.Set("Access-Control-Expose-Headers", "trace-id, trace-url, trace-logs-url")
			header.Set("Access-Control-Max-Age", "86400")

			if req.Method == http.MethodOptions {
				res.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
