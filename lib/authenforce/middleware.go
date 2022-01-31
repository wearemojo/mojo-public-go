package authenforce

import (
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/wearemojo/mojo-public-go/lib/authparsing"
)

func (e Enforcer) CRPCMiddleware() crpc.MiddlewareFunc {
	return func(next crpc.HandlerFunc) crpc.HandlerFunc {
		return func(res http.ResponseWriter, req *crpc.Request) error {
			state := authparsing.GetAuthState(req.Context())
			err := e(state)
			if err != nil {
				return err
			}

			return next(res, req)
		}
	}
}
