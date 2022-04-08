package authparsing

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/wearemojo/mojo-public-go/lib/gerrors"
)

func jsonError(ctx context.Context, res http.ResponseWriter, err error) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	enc := json.NewEncoder(res)
	var encErr error

	if err, ok := gerrors.As[cher.E](err); ok {
		res.WriteHeader(err.StatusCode())
		encErr = enc.Encode(err)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
		encErr = enc.Encode(cher.New(cher.Unknown, cher.M{"error": err}))
	}

	if encErr != nil {
		clog.Get(ctx).WithError(encErr).Error("failed to encode error")
	}
}

func Middleware(parser Parser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			authzHeader := req.Header.Get("Authorization")

			authState, err := parser.Check(ctx, authzHeader)
			if err != nil && !errors.Is(err, ErrNoAuthorization) {
				jsonError(ctx, res, err)

				if cerr, ok := gerrors.As[cher.E](err); ok && cerr.Code == cher.Unauthorized && len(cerr.Reasons) == 1 {
					err = cerr.Reasons[0]
				}
				clog.Get(ctx).WithError(err).Info("auth check failed")

				return
			}

			ctx = SetAuthState(ctx, authState)
			req = req.WithContext(ctx)

			next.ServeHTTP(res, req)
		})
	}
}
