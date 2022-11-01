package selforigin

import "context"

type contextKey string

const contextKeySelfOrigin contextKey = "self_origin"

func ContextWithSelfOrigin(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, contextKeySelfOrigin, val)
}

func ContextSelfOrigin(ctx context.Context) (val string) {
	val, _ = ctx.Value(contextKeySelfOrigin).(string)
	return
}
