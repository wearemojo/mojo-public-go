package secret

import (
	"context"
	"fmt"
)

type Provider interface {
	Get(ctx context.Context, secretID string) (string, error)
}

type contextKey string

const contextKeyProvider contextKey = "provider"

func getProvider(ctx context.Context) (val Provider) {
	val, _ = ctx.Value(contextKeyProvider).(Provider)
	return
}

func ContextWithProvider(ctx context.Context, val Provider) context.Context {
	return context.WithValue(ctx, contextKeyProvider, val)
}

func Get(ctx context.Context, secretID string) (string, error) {
	p := getProvider(ctx)
	if p == nil {
		return "", fmt.Errorf("secret provider not set")
	}

	return p.Get(ctx, secretID)
}
