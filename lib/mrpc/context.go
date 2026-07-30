package mrpc

import (
	"context"
	"net/http"
)

// HandlerFunc is like http.HandlerFunc, but also allows error returns.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Middleware allows HandlerFunc wrapping in the usual Go style.
type Middleware func(next HandlerFunc) HandlerFunc

// requestInfo describes the RPC request currently being served. It is attached to
// the request context and readable via getInfo.
type requestInfo struct {
	// Version is the version segment as requested: a date, "latest", or
	// "preview".
	Version string

	// Method is the requested method name.
	Method string
}

type contextKey string

const infoKey contextKey = "mrpc_info"

// getInfo returns the RPC requestInfo from the context, or nil outside an RPC
// request.
func getInfo(ctx context.Context) *requestInfo {
	if val, ok := ctx.Value(infoKey).(*requestInfo); ok {
		return val
	}

	return nil
}

func setInfo(ctx context.Context, info *requestInfo) context.Context {
	return context.WithValue(ctx, infoKey, info)
}
