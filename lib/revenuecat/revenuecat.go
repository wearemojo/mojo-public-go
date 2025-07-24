package revenuecat

import (
	"context"
	"fmt"
	"maps"
	"net/url"
	"slices"
	"time"

	"github.com/igrmk/decimal"
	"github.com/samber/lo"
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
	ExpiresDate            *time.Time `json:"expires_date"`
	GracePeriodExpiresDate *time.Time `json:"grace_period_expires_date"`
	ProductIdentifier      string     `json:"product_identifier"`
	PurchaseDate           time.Time  `json:"purchase_date"`
}

type Subscription struct {
	AutoResumeDate          *time.Time    `json:"auto_resume_date"`
	BillingIssuesDetectedAt *time.Time    `json:"billing_issues_detected_at"`
	DisplayName             string        `json:"display_name"`
	ExpiresDate             time.Time     `json:"expires_date"`
	GracePeriodExpiresDate  *time.Time    `json:"grace_period_expires_date"`
	IsSandbox               bool          `json:"is_sandbox"`
	OriginalPurchaseDate    time.Time     `json:"original_purchase_date"`
	OwnershipType           OwnershipType `json:"ownership_type"`
	PeriodType              PeriodType    `json:"period_type"`
	Price                   Price         `json:"price"`
	PurchaseDate            time.Time     `json:"purchase_date"`
	RefundedAt              *time.Time    `json:"refunded_at"`
	Store                   StoreType     `json:"store"`
	StoreTransactionID      string        `json:"store_transaction_id"`
	UnsubscribeDetectedAt   *time.Time    `json:"unsubscribe_detected_at"`
}

type NonSubscription struct {
	DisplayName  string    `json:"display_name"`
	ID           string    `json:"id"`
	IsSandbox    bool      `json:"is_sandbox"`
	Price        Price     `json:"price"`
	PurchaseDate time.Time `json:"purchase_date"`
	Store        StoreType `json:"store"`
}

type Price struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

type OwnershipType string

const (
	OwnershipTypePurchased    OwnershipType = "PURCHASED"
	OwnershipTypeFamilyShared OwnershipType = "FAMILY_SHARED"
)

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

func (s *Subscriber) ActiveSubscriptions() map[string]Subscription {
	now := time.Now()
	return lo.PickBy(s.Subscriptions, func(key string, value Subscription) bool {
		return value.ExpiresDate.After(now)
	})
}

func (s *Subscriber) ActiveCount() int {
	nonSub := lo.Flatten(slices.Collect(maps.Values(s.NonSubscriptions)))
	return len(nonSub) + len(s.ActiveSubscriptions())
}
