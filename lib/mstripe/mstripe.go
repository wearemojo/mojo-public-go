package mstripe

import (
	"context"
	"time"

	stripeV78 "github.com/stripe/stripe-go/v78"
	clientV78 "github.com/stripe/stripe-go/v78/client"
	"github.com/stripe/stripe-go/v85"
	"github.com/wearemojo/mojo-public-go/lib/httpclient"
	"github.com/wearemojo/mojo-public-go/lib/secret"
)

// Stripe's default is 80s, but 15s is the highest in our system so far
const timeout = 15 * time.Second

func NewClientV78(ctx context.Context, keySecretID string) (*clientV78.API, error) {
	key, err := secret.Get(ctx, keySecretID)
	if err != nil {
		return nil, err
	}

	backends := stripeV78.NewBackendsWithConfig(&stripeV78.BackendConfig{
		HTTPClient: httpclient.NewClient(timeout, nil),

		LeveledLogger: &stripeV78.LeveledLogger{
			// Stripe regularly logs e.g. expected 404s as the highest error level
			// all actual errors are handled normally through err returns, so there's
			// currently no value in its logging
			Level: stripeV78.LevelNull,
		},
	})

	return clientV78.New(key, backends), nil
}

func NewClient(ctx context.Context, keySecretID string) (*stripe.Client, error) {
	key, err := secret.Get(ctx, keySecretID)
	if err != nil {
		return nil, err
	}

	backends := stripe.NewBackendsWithConfig(&stripe.BackendConfig{
		HTTPClient: httpclient.NewClient(timeout, nil),

		LeveledLogger: &stripe.LeveledLogger{
			// Stripe regularly logs e.g. expected 404s as the highest error level
			// all actual errors are handled normally through err returns, so there's
			// currently no value in its logging
			Level: stripe.LevelNull,
		},
	})

	return stripe.NewClient(key, stripe.WithBackends(backends)), nil
}

func Collect[T any](seq stripe.Seq2[*T, error]) ([]T, error) {
	res := []T{}
	for item, err := range seq {
		if err != nil {
			return nil, err
		}

		res = append(res, *item)
	}

	return res, nil
}
