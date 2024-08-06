package servicecontext

import (
	"context"
)

// Info type holds useful info about the currently-running service
type Info struct {
	System  string
	Env     string
	Service string
}

type contextKey string

var infoContextKey = contextKey("info")

// SetContext wraps the context with the service info
func SetContext(ctx context.Context, info Info) context.Context {
	return context.WithValue(ctx, infoContextKey, info)
}

// GetContext retrieves the service info from the context
func GetContext(ctx context.Context) *Info {
	if val, ok := ctx.Value(infoContextKey).(Info); ok {
		return &val
	}

	return nil
}
