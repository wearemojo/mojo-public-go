package secret

import (
	"context"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

const ErrSecretNotFound = merr.Code("secret_not_found")

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
		return "", merr.New(ctx, "unset_secret_provider", nil)
	}

	return p.Get(ctx, secretID)
}
