package authenforce

import (
	"context"
	"sync/atomic"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/errgroup"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

type (
	Enforcers   []Enforcer
	Enforcer    func(context.Context, any, map[string]any) (bool, error)
	SubEnforcer func(context.Context, any, map[string]any) error
)

const ErrNotHandled = merr.Code("auth_not_handled")

func UnsafeNoAuthentication(_ context.Context, _ any, _ map[string]any) (bool, error) {
	return true, nil
}

func AllowAny(_ context.Context, state any, _ map[string]any) (bool, error) {
	if state == nil {
		return true, cher.New(cher.Unauthorized, nil)
	}

	return true, nil
}

// Run runs all enforcers in parallel and ensures that none of them error.
//
// It is important to note that only one enforcer can ackwoledge a request. If
// multiple enforcers ack, an error will be returned as it is unsafe. It is
// therefore important that your enforcers are focused to find requests that are
// applicable and quickly return that it won't handle the auth.
func (e Enforcers) Run(ctx context.Context, authState any, mapReq map[string]any) error {
	g := errgroup.WithContext(ctx)

	var handleCount uint64

	for _, enforcer := range e {
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

	return nil
}
