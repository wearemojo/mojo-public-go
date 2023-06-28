package authenforce

import (
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/wearemojo/mojo-public-go/lib/authparsing"
	"github.com/wearemojo/mojo-public-go/lib/bodycontext"
)

func CRPCMiddleware(enforcers Enforcers) crpc.MiddlewareFunc {
	return func(next crpc.HandlerFunc) crpc.HandlerFunc {
		return func(res http.ResponseWriter, req *crpc.Request) error {
			ctx := req.Context()
			authState := authparsing.GetAuthState(ctx)

			body := bodycontext.GetContext(ctx)

			if err := enforcers.Run(ctx, authState, body); err != nil {
				return err
			}

			return next(res, req)
		}
	}
}
