package authparsing

import (
	"context"
)

type contextKey string

const authStateKey contextKey = "auth_state"

func GetAuthState(ctx context.Context) (val *AuthState) {
	val, _ = ctx.Value(authStateKey).(*AuthState)
	return
}

func SetAuthState(ctx context.Context, val *AuthState) context.Context {
	return context.WithValue(ctx, authStateKey, val)
}
