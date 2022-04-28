package authenforce

import (
	"context"
	"sync/atomic"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/errgroup"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

type (
	// Enforcers is a set of Enforcer functions that are run together.
	//
	// It is important to note that only one enforcer can acknowledge a request. If
	// multiple enforcers ack, an error will be returned as it is unsafe. It is
	// therefore important that your enforcers are focused to find requests that are
	// applicable and quickly return that it won't handle the auth.
	Enforcers []Enforcer

	// Enforcer checks the auth state type. An example may be checking it is a
	// user, a service, or some specific form of authentication.
	Enforcer func(context.Context, any, map[string]any) (handled bool, err error)
)

const ErrNotHandled = merr.Code("auth_not_handled")

func UnsafeNoAuthentication(_ context.Context, _ any, _ map[string]any) (bool, error) {
	return true, nil
}

func AllowAny(_ context.Context, state any, _ map[string]any) (bool, error) {
	if state == nil {
		return true, cher.New("auth_not_provided", nil)
	}

	return true, nil
}

// Run all enforcers in parallel and ensures that none of them error.
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
		return cher.New(cher.AccessDenied, nil, cher.Coerce(err))
	}

	switch {
	case handleCount == 1: // all good
		return nil
	case handleCount == 0 && authState == nil:
		return cher.New(cher.Unauthorized, nil, cher.New("no_enforcement", nil))
	case handleCount == 0 && authState != nil:
		return cher.New(cher.AccessDenied, nil, cher.New("auth_not_suitable", nil))
	default: // handle count is greater than 1
		return merr.New("multiple_enforcers_ran", nil)
	}
}
