package revenuecat

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/httpclient"
	"github.com/wearemojo/mojo-public-go/lib/jsonclient"
	"github.com/wearemojo/mojo-public-go/lib/secret"
)

const baseURL = "https://api.revenuecat.com/v1"

type Client struct {
	client *jsonclient.Client
}

func NewClient(ctx context.Context, serverTokenSecretID string) (*Client, error) {
	if _, err := secret.Get(ctx, serverTokenSecretID); err != nil {
		return nil, err
	}

	return &Client{
		client: jsonclient.NewClient(
			baseURL,
			httpclient.NewClient(5*time.Second, roundTripper{serverTokenSecretID}),
		),
	}, nil
}

type UserInfoResponse struct {
	Subscriber Subscriber `json:"subscriber"`
}

type Subscriber struct {
	Entitlements     map[string]Entitlement       `json:"entitlements"`
	Subscriptions    map[string]Subscription      `json:"subscriptions"`
	NonSubscriptions map[string][]NonSubscription `json:"non_subscriptions"`
	ManagementURL    *string                      `json:"management_url"`
}

type Entitlement struct {
	ExpiresDate       *time.Time `json:"expires_date"`
	ProductIdentifier string     `json:"product_identifier"`
}

type Subscription struct {
	PeriodType            PeriodType `json:"period_type"`
	UnsubscribeDetectedAt *time.Time `json:"unsubscribe_detected_at"`
	Store                 StoreType  `json:"store"`
}

type NonSubscription struct {
	Store StoreType `json:"store"`
}

type PeriodType string

const (
	PeriodTypeTrial  PeriodType = "trial"
	PeriodTypeIntro  PeriodType = "intro"
	PeriodTypeNormal PeriodType = "normal"
)

// Possible values for store:
// - app_store: The product was purchased through Apple App Store.
// - play_store: The product was purchased through the Google Play Store.
// - stripe: The product was purchased through Stripe.
// - promotional: The product was granted via RevenueCat.
type StoreType string

const (
	StoreTypeAppStore    StoreType = "app_store"
	StoreTypePlayStore   StoreType = "play_store"
	StoreTypeStripe      StoreType = "stripe"
	StoreTypePromotional StoreType = "promotional"
)

// https://www.revenuecat.com/reference/subscribers
// This API either fetches a user or creates one :shrug:
// for now, we only care about the expiry date of the full_access entitlement
// and the period type of the subscription. See useSubscriptionStateIap.ts
func (c *Client) GetOrCreateSubscriberInfo(ctx context.Context, appUserID string) (res *UserInfoResponse, err error) {
	escapedAppUserID := url.PathEscape(appUserID)
	path := fmt.Sprintf("/subscribers/%s", escapedAppUserID)
	return res, c.client.Do(ctx, "GET", path, nil, nil, &res)
}

func (s *Subscriber) ActiveSubscriptionCount() int {
	activeNonSubCount := len(s.NonSubscriptions) // always active
	activeSubCount := 0
	now := time.Now()

	for _, e := range s.Entitlements {
		if e.ExpiresDate != nil && e.ExpiresDate.Before(now) {
			continue
		}

		activeSubCount++
	}

	return activeNonSubCount + activeSubCount
}
