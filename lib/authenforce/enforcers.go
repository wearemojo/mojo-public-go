package authenforce

import (
	"context"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

type (
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
