package authenforce

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/wearemojo/mojo-public-go/lib/authparsing"
)

func CRPCMiddleware(enforcers Enforcers) crpc.MiddlewareFunc {
	return func(next crpc.HandlerFunc) crpc.HandlerFunc {
		return func(res http.ResponseWriter, req *crpc.Request) error {
			ctx := req.Context()
			authState := authparsing.GetAuthState(ctx)
			var mapReq map[string]any

			pr, pw := io.Pipe() // TODO: test thoroughly
			tr := io.TeeReader(req.Body, pw)

			if err := json.NewDecoder(tr).Decode(&mapReq); err != nil {
				return err
			}

			req.Body = io.NopCloser(pr)

			if err := enforcers.Run(ctx, authState, mapReq); err != nil {
				return err
			}

			return next(res, req)
		}
	}
}
