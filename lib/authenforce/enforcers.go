package authenforce

import (
	"context"
	"errors"
	"sync"

	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
	"github.com/wearemojo/mojo-public-go/lib/slicefn"
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

type enforcerOutcome struct {
	handled bool
	err     error
}

// Run all enforcers in parallel and ensures that none of them error.
func (e Enforcers) Run(ctx context.Context, authState any, req []byte) error {
	var outcomes []enforcerOutcome
	var outcomeMutex sync.Mutex

	var wg sync.WaitGroup

	for _, enforcer := range e {
		wg.Go(func() {
			handled, err := enforcer(ctx, authState, req)

			outcomeMutex.Lock()
			defer outcomeMutex.Unlock()

			outcomes = append(outcomes, enforcerOutcome{
				handled: handled,
				err:     err,
			})
		})
	}

	wg.Wait()

	handled := slicefn.Filter(outcomes, func(outcome enforcerOutcome) bool {
		return outcome.handled
	})

	unhandled := slicefn.Filter(outcomes, func(outcome enforcerOutcome) bool {
		return !outcome.handled && outcome.err != nil
	})

	unhandledErrs := slicefn.Map(unhandled, func(outcome enforcerOutcome) error {
		return outcome.err
	})

	if len(unhandledErrs) != 0 {
		mlog.Warn(ctx, merr.New(ctx, "unhandled_enforcer_errors", nil, unhandledErrs...))
	}

	if len(handled) == 0 {
		return wrapAuthErrType(authState, cher.New("no_suitable_auth_found", nil))
	}

	if len(handled) != 1 {
		return merr.New(ctx, "multiple_enforcers_handled", merr.M{
			"outcomes": outcomes,
		})
	}

	outcome := handled[0]

	if outcome.err == nil {
		return nil
	}

	if cerr, ok := errors.AsType[cher.E](outcome.err); ok {
		return wrapAuthErrType(authState, cerr)
	}

	return merr.New(ctx, "enforcer_failed", nil, outcome.err)
}

func wrapAuthErrType(authState any, err cher.E) error {
	if authState == nil {
		return cher.New(cher.Unauthorized, nil, err)
	}

	return cher.New(cher.AccessDenied, nil, err)
}
