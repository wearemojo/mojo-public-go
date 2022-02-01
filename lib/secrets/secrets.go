package secrets

import (
	"context"
	"fmt"
)

type Provider interface {
	Get(ctx context.Context, secretID string) (string, error)
}

type contextKey string

const providerKey contextKey = "provider"

func getProvider(ctx context.Context) (val Provider) {
	val, _ = ctx.Value(providerKey).(Provider)
	return
}

func ContextWithProvider(ctx context.Context, val Provider) context.Context {
	return context.WithValue(ctx, providerKey, val)
}

func Get(ctx context.Context, secretID string) (string, error) {
	p := getProvider(ctx)
	if p == nil {
		return "", fmt.Errorf("secrets provider not set")
	}

	return p.Get(ctx, secretID)
}
