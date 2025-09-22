package reqctx

import (
	"context"
	"net/http"
)

type contextKey string

const (
	contextKeyRequest contextKey = "request"
)

func GetRequest(ctx context.Context) (req *http.Request) {
	req, _ = ctx.Value(contextKeyRequest).(*http.Request)
	return req
}

func SetRequest(ctx context.Context, val *http.Request) context.Context {
	return context.WithValue(ctx, contextKeyRequest, val)
}
