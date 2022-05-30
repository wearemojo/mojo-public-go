package authparsing

import (
	"context"
)

type contextKey string

const contextKeyAuthState contextKey = "auth_state"

func GetAuthState(ctx context.Context) any {
	return ctx.Value(contextKeyAuthState)
}

func SetAuthState(ctx context.Context, val any) context.Context {
	return context.WithValue(ctx, contextKeyAuthState, val)
}
