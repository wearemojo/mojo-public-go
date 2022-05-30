package authparsing

import (
	"context"
	"net/http"
)

type contextKey string

const (
	contextKeyAuthState contextKey = "auth_state"
	contextKeyRequest   contextKey = "auth_request"
)

func GetAuthState(ctx context.Context) any {
	return ctx.Value(contextKeyAuthState)
}

func SetAuthState(ctx context.Context, val any) context.Context {
	return context.WithValue(ctx, contextKeyAuthState, val)
}

func GetRequest(ctx context.Context) (req *http.Request) {
	req, _ = ctx.Value(contextKeyRequest).(*http.Request)
	return
}

func SetRequest(ctx context.Context, val *http.Request) context.Context {
	return context.WithValue(ctx, contextKeyRequest, val)
}
