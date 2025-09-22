package ipcountry

import (
	"context"
)

type contextKey string

const contextKeyIPCountry contextKey = "ip_country"

func GetIPCountry(ctx context.Context) (val string) {
	val, _ = ctx.Value(contextKeyIPCountry).(string)
	return val
}

func SetIPCountry(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, contextKeyIPCountry, val)
}
