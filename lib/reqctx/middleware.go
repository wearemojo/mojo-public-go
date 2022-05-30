package reqctx

import (
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = SetRequest(ctx, req)
		req = req.WithContext(ctx)

		next.ServeHTTP(res, req)
	})
}
