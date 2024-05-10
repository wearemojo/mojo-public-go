package stripeclient

import (
	"context"
	"time"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
	"github.com/wearemojo/mojo-public-go/lib/httpclient"
	"github.com/wearemojo/mojo-public-go/lib/secret"
)

// Stripe's default is 80s, but 15s is the highest in our system so far
const timeout = 15 * time.Second

func New(ctx context.Context, keySecretID string) (*client.API, error) {
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

	return client.New(key, backends), nil
}
