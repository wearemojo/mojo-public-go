package googleiid

import (
	"context"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/errgroup"
	"github.com/wearemojo/mojo-public-go/lib/httpclient"
	"github.com/wearemojo/mojo-public-go/lib/jsonclient"
	"github.com/wearemojo/mojo-public-go/lib/secret"
)

const baseURL = "https://iid.googleapis.com"

type Client struct {
	client *jsonclient.Client
}

func NewClient(ctx context.Context, serverKeySecretID, vapidPublicKeySecretID string) (*Client, error) {
	g := errgroup.WithContext(ctx)

	g.Go(func(ctx context.Context) (err error) {
		_, err = secret.Get(ctx, serverKeySecretID)
		return err
	})

	g.Go(func(ctx context.Context) (err error) {
		_, err = secret.Get(ctx, vapidPublicKeySecretID)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &Client{
		client: jsonclient.NewClient(
			baseURL,
			httpclient.NewClient(5*time.Second, roundTripper{
				ServerKeySecretID:      serverKeySecretID,
				VAPIDPublicKeySecretID: vapidPublicKeySecretID,
			}),
		),
	}, nil
}

type APNSRequest struct {
	Application string   `json:"application"`
	Sandbox     bool     `json:"sandbox"`
	APNSTokens  []string `json:"apns_tokens"`
}

type APNSResponse struct {
	// order does not relate to the request ordering - match by token
	Results []APNSResponseResult `json:"results"`
}

type APNSResponseResult struct {
	RegistrationToken string `json:"registration_token"`
	APNSToken         string `json:"apns_token"`
	Status            string `json:"status"`
}

func (r APNSResponseResult) Valid() bool {
	// we've also noticed `INVALID_ARGUMENT` and `INTERNAL` on the status field
	return r.Status == "OK" && r.APNSToken != "" && r.RegistrationToken != ""
}

type WebPushRequest struct {
	Endpoint string `json:"endpoint"`

	Keys WebPushRequestKeys `json:"keys"`
}

type WebPushRequestKeys struct {
	Auth   string `json:"auth"`
	P256DH string `json:"p256dh"`
}

type WebPushResponse struct {
	Token string `json:"token"`
}

func (c *Client) ImportAPNSTokens(ctx context.Context, req *APNSRequest) (res *APNSResponse, err error) {
	// https://web.archive.org/web/20220407013020/https://developers.google.com/instance-id/reference/server#create_registration_tokens_for_apns_tokens
	return res, c.client.Do(ctx, "POST", "iid/v1:batchImport", nil, req, &res)
}

// Deprecated: https://firebase.google.com/support/faq#fcm-depr-features
//
// supported equivalent: https://firebase.google.com/docs/cloud-messaging/js/client#access_the_registration_token
// or rip from: https://github.com/firebase/firebase-js-sdk/blob/master/packages/messaging/src/internals/requests.ts#L39
func (c *Client) ImportWebPushSubscription(ctx context.Context, req *WebPushRequest) (res *WebPushResponse, err error) {
	// https://web.archive.org/web/20220407013020/https://developers.google.com/instance-id/reference/server#import_push_subscriptions
	return res, c.client.Do(ctx, "POST", "v1/web/iid", nil, req, &res)
}
