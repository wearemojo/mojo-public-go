package bodycontext

import (
	"context"
)

type contextKey string

var bodyContextKey = contextKey("body")

// SetContext wraps the context with the body bytes
func SetContext(ctx context.Context, val []byte) context.Context {
	return context.WithValue(ctx, bodyContextKey, val)
}

// GetContext retrieves the body  bytes from the context
func GetContext(ctx context.Context) []byte {
	if val, ok := ctx.Value(bodyContextKey).([]byte); ok {
		return val
	}

	return nil
}
