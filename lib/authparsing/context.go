package authparsing

import (
	"context"
)

type contextKey string

const contextKeyAuthState contextKey = "auth_state"

func GetAuthState(ctx context.Context) (val *AuthState) {
	val, _ = ctx.Value(contextKeyAuthState).(*AuthState)
	return
}

func SetAuthState(ctx context.Context, val *AuthState) context.Context {
	return context.WithValue(ctx, contextKeyAuthState, val)
}
