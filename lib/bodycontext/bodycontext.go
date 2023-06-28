package bodycontext

import (
	"bytes"
	"io"
	"net/http"

	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			mlog.Warn(ctx, merr.New(ctx, "cannot_read_body", nil, err))
			next.ServeHTTP(res, req)
			return
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body))

		ctx = SetContext(ctx, body)
		req = req.WithContext(ctx)

		next.ServeHTTP(res, req)
	})
}
