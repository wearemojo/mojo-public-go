package authenforce

import (
	"bytes"
	"encoding/json"
	"errors"
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

			var buf bytes.Buffer
			tee := io.TeeReader(req.Body, &buf)

			if err := json.NewDecoder(tee).Decode(&mapReq); err != nil && !errors.Is(err, io.EOF) {
				return err
			}

			req.Body = io.NopCloser(&buf)

			if err := enforcers.Run(ctx, authState, mapReq); err != nil {
				return err
			}

			return next(res, req)
		}
	}
}
