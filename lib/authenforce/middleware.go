package authenforce

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/wearemojo/mojo-public-go/lib/authparsing"
	"github.com/wearemojo/mojo-public-go/lib/errgroup"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

func (e Enforcers) CRPCMiddleware() crpc.MiddlewareFunc {
	return func(next crpc.HandlerFunc) crpc.HandlerFunc {
		return func(res http.ResponseWriter, req *crpc.Request) error {
			ctx := req.Context()
			authState := authparsing.GetAuthState(ctx)
			mapReq := map[string]any{}

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return merr.Wrap(err, "cannot_read_body", nil)
			}

			if len(body) > 0 {
				if err := json.Unmarshal(body, &mapReq); err != nil {
					return err
				}
			}

			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			g := errgroup.WithContext(ctx) //nolint:varnamelen

			var handleCount uint64

			for _, enforcer := range e {
				enforcer := enforcer

				g.Go(func(ctx context.Context) (err error) {
					err = enforcer(ctx, authState, mapReq)
					if err == nil {
						atomic.AddUint64(&handleCount, 1)
						return
					}

					if errors.Is(err, ErrNotHandled) {
						err = nil
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
