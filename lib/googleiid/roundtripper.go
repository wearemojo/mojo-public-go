package googleiid

import (
	"context"
	"net/http"

	"github.com/wearemojo/mojo-public-go/lib/errgroup"
	"github.com/wearemojo/mojo-public-go/lib/secret"
)

type roundTripper struct {
	ServerKeySecretID      string
	VAPIDPublicKeySecretID string
}

func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	req = req.Clone(ctx)

	var serverKey string
	var vapidPublicKey string

	g := errgroup.WithContext(ctx)

	g.Go(func(ctx context.Context) (err error) {
		serverKey, err = secret.Get(ctx, r.ServerKeySecretID)
		return
	})

	g.Go(func(ctx context.Context) (err error) {
		vapidPublicKey, err = secret.Get(ctx, r.VAPIDPublicKeySecretID)
		return
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if req.Header == nil {
		req.Header = http.Header{}
	}

	// https://web.archive.org/web/20221206045856/https://firebase.google.com/docs/cloud-messaging/auth-server#authorize-http-requests
	req.Header.Set("Authorization", "key="+serverKey)

	// https://developers.google.com/instance-id/reference/server#parameters_5
	req.Header.Set("Crypto-Key", "p256ecdsa="+vapidPublicKey)

	return http.DefaultTransport.RoundTrip(req)
}
