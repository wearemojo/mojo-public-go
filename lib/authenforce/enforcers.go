package authenforce

import (
	"context"
	"sync/atomic"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/errgroup"
	"github.com/wearemojo/mojo-public-go/lib/gerrors"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

type (
	// Enforcers is a set of Enforcer functions that are run together.
	//
	// Enforcers are given auth state (any) and must first make a determination
	// to handle that auth type. If an enforcer is not applicable, it should
	// return false with no error. An enforcer should return false with an error
	// if it is unable to run.
	//
	// From this stage, an enforcer should return handled as true.
	// Once checks have completed, return no error to indicate the request
	// should proceed, or an error to deny access.
	//
	// It is important for an enforcer to be tightly scoped when it is
	// determining if it is applicable, as if multiple enforcers return true,
	// the request will be rejected.
	Enforcers []Enforcer

	Enforcer func(context.Context, any, []byte) (handled bool, err error)
)

func UnsafeNoAuthentication(_ context.Context, _ any, _ []byte) (bool, error) {
	return true, nil
}

func UnsafeAllowAny(_ context.Context, state any, _ []byte) (bool, error) {
	if state == nil {
		return true, cher.New("auth_not_provided", nil)
	}

	return true, nil
}

// Run all enforcers in parallel and ensures that none of them error.
func (e Enforcers) Run(ctx context.Context, authState any, req []byte) error {
	g := errgroup.WithContext(ctx)

	var handleCount uint64

	for _, enforcer := range e {
		enforcer := enforcer

		g.Go(func(ctx context.Context) error {
			handled, err := enforcer(ctx, authState, req)
			if err != nil {
				return err
			}

			if handled {
				atomic.AddUint64(&handleCount, 1)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if cerr, ok := gerrors.As[cher.E](err); ok {
			return wrapAuthErrType(authState, cerr)
		}

		return merr.Wrap(ctx, err, "enforcer_failed", nil)
	}

	switch handleCount {
	case 0:
		return wrapAuthErrType(authState, cher.New("no_suitable_auth_found", nil))
	case 1: // all is safe, exactly one enforcer ran
		return nil
	default:
		return merr.New(ctx, "multiple_enforcers_handled", merr.M{"count": handleCount})
	}
}

func wrapAuthErrType(authState any, err cher.E) error {
	if authState == nil {
		return cher.New(cher.Unauthorized, nil, err)
	}

	return cher.New(cher.AccessDenied, nil, err)
}
