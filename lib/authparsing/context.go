package authparsing

import (
	"context"
)

type contextKey string

const contextKeyAuthState contextKey = "auth_state"

func GetAuthState(ctx context.Context) (val any) {
	val = ctx.Value(contextKeyAuthState)
	return
}

func SetAuthState(ctx context.Context, val any) context.Context {
	return context.WithValue(ctx, contextKeyAuthState, val)
}
