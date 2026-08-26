package authparsing

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"net/http"

	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/clog"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
)

// json/v2 marshals maps in a non-deterministic order by default, whereas v1
// sorted the keys - error metadata is a map, and a stable wire format keeps
// responses diffable and cacheable
var deterministic = json.Deterministic(true)

func jsonError(ctx context.Context, res http.ResponseWriter, err error) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	enc := jsontext.NewEncoder(res)
	var encErr error

	if err, ok := errors.AsType[cher.E](err); ok {
		res.WriteHeader(err.StatusCode())
		encErr = json.MarshalEncode(enc, err, deterministic)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
		encErr = json.MarshalEncode(enc, cher.New(cher.Unknown, cher.M{"error": err}), deterministic)
	}

	if encErr != nil {
		mlog.Error(ctx, merr.New(ctx, "error_encode_failed", nil, encErr))
	}
}

func Middleware(parser Parser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			authzHeader := req.Header.Get("Authorization")

			authState, err := parser.Check(ctx, authzHeader)
			if err != nil && !errors.Is(err, ErrNoAuthorization) {
				clog.SetError(ctx, err)
				jsonError(ctx, res, err)

				if cerr, ok := errors.AsType[cher.E](err); ok && cerr.Code == cher.Unauthorized && len(cerr.Reasons) == 1 {
					err = cerr.Reasons[0]
				}
				mlog.Info(ctx, merr.New(ctx, "auth_check_failed", nil, err))

				return
			}

			ctx = SetAuthState(ctx, authState)
			req = req.WithContext(ctx)

			next.ServeHTTP(res, req)
		})
	}
}
