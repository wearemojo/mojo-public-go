package authparsing

import (
	"encoding/json"
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/clog"
)

func jsonError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	e := json.NewEncoder(w)

	if err, ok := err.(cher.E); ok {
		w.WriteHeader(err.StatusCode())
		_ = e.Encode(err)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_ = e.Encode(cher.New(cher.Unknown, cher.M{"error": err}))
	}
}

func Middleware(parser Parser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authzHeader := r.Header.Get("Authorization")

			authState, err := parser.Check(ctx, authzHeader)
			if err != nil {
				jsonError(w, err)

				if cerr, ok := err.(cher.E); ok && cerr.Code == cher.Unauthorized && len(cerr.Reasons) == 1 {
					err = cerr.Reasons[0]
				}
				clog.Get(ctx).WithError(err).Info("auth check failed")

				return
			}

			ctx = SetAuthState(ctx, authState)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
