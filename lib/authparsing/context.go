package authparsing

import (
	"context"
)

type contextKey string

const authStateKey contextKey = "auth_state"

func GetAuthState(ctx context.Context) *AuthState {
	if val, ok := ctx.Value(authStateKey).(*AuthState); ok {
		return val
	}

	return nil
}

func SetAuthState(ctx context.Context, val *AuthState) context.Context {
	return context.WithValue(ctx, authStateKey, val)
}
