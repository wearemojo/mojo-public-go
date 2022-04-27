package authenforce

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/wearemojo/mojo-public-go/lib/authparsing"
	"github.com/wearemojo/mojo-public-go/lib/errgroup"
)

func CRPCMiddleware(enforcers []Enforcer) crpc.MiddlewareFunc {
	return func(next crpc.HandlerFunc) crpc.HandlerFunc {
		return func(res http.ResponseWriter, req *crpc.Request) error {
			ctx := req.Context()
			authState := authparsing.GetAuthState(ctx)
			var mapReq map[string]any

			pr, pw := io.Pipe()
			tr := io.TeeReader(req.Body, pw)

			if err := json.NewDecoder(tr).Decode(&mapReq); err != nil {
				return err
			}

			req.Body = io.NopCloser(pr)

			g := errgroup.WithContext(ctx)

			var handleCount uint64

			for _, enforcer := range enforcers {
				enforcer := enforcer

				g.Go(func(ctx context.Context) (err error) {
					handled, err := enforcer(ctx, authState, mapReq)
					if err != nil {
						return
					}

					if handled {
						atomic.AddUint64(&handleCount, 1)
					}

					return
				})
			}

			if err := g.Wait(); err != nil {
				return err
			}

			if handleCount != 1 {
				return cher.New(cher.AccessDenied, nil, cher.New("enforcement_dispute", cher.M{"count": handleCount}))
			}

			return next(res, req)
		}
	}
}
