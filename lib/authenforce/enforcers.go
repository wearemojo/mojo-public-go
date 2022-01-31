package authenforce

import (
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/authparsing"
)

type Enforcer func(*authparsing.AuthState) error

func UnsafeNoAuthentication(*authparsing.AuthState) error {
	return nil
}

func AllowAny(state *authparsing.AuthState) error {
	if state == nil {
		return cher.New(cher.Unauthorized, nil)
	}

	return nil
}

func RequireS2S(state *authparsing.AuthState) error {
	if state == nil {
		return cher.New(cher.Unauthorized, nil)
	}

	if state.Type != authparsing.S2S {
		return cher.New(cher.AccessDenied, nil)
	}

	return nil
}
