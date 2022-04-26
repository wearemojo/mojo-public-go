package authenforce

import (
	"context"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

type (
	MapRequest map[string]any
	Enforcer   func(context.Context, any, MapRequest) error
	Enforcers  []Enforcer
)

var ErrNotHandled = merr.New("auth_not_handled", nil)

func UnsafeNoAuthentication(_ context.Context, _ any, _ MapRequest) error {
	return nil
}

func AllowAny(_ context.Context, state any, _ MapRequest) error {
	if state == nil {
		return cher.New(cher.Unauthorized, nil)
	}

	return nil
}
