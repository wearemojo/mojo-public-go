package authenforce

import (
	"bytes"
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

			body, err := io.ReadAll(req.Body)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				if err := json.Unmarshal(body, &mapReq); err != nil {
					return err
				}
			}

			req.Body = io.NopCloser(bytes.NewBuffer(body))

			if err := enforcers.Run(ctx, authState, mapReq); err != nil {
				return err
			}

			return next(res, req)
		}
	}
}
