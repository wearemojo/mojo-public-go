//nolint:gocritic // omitempty required for the Flex API
package flex

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/httpclient"
	"github.com/wearemojo/mojo-public-go/lib/jsonclient"
	"github.com/wearemojo/mojo-public-go/lib/secret"
)

const baseURL = "https://api.withflex.com"

type Client struct {
	client *jsonclient.Client
}

func NewClient(ctx context.Context, apiKeySecretID string) (*Client, error) {
	if _, err := secret.Get(ctx, apiKeySecretID); err != nil {
		return nil, err
	}

	return &Client{
		client: jsonclient.NewClient(
			baseURL,
			httpclient.NewClient(15*time.Second, roundTripper{apiKeySecretID}),
		),
	}, nil
}

type CreateCheckoutSessionRequest struct {
	CheckoutSession CheckoutSessionParams `json:"checkout_session"`
}

type CheckoutSessionParams struct {
	CancelURL         string            `json:"cancel_url"`
	ClientReferenceID string            `json:"client_reference_id,omitempty"`
	Defaults          *Defaults         `json:"defaults,omitempty"`
	LineItems         []LineItem        `json:"line_items"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	Mode              string            `json:"mode"`
	SubscriptionData  *SubscriptionData `json:"subscription_data,omitempty"`
	SuccessURL        string            `json:"success_url"`
}

type SubscriptionData struct {
	TrialEnd          *time.Time `json:"trial_end,omitempty"`
	CancelAt          *time.Time `json:"cancel_at,omitempty"`
	CancelAtPeriodEnd *bool      `json:"cancel_at_period_end,omitempty"`
	TrialPeriodDays   *int       `json:"trial_period_days,omitempty"`
}

type Defaults struct {
	Email string `json:"email,omitempty"`
}

type LineItem struct {
	Price     string     `json:"price,omitempty"`
	PriceData *PriceData `json:"price_data,omitempty"`
	Quantity  int        `json:"quantity"`
}

type PriceData struct {
	Product    string     `json:"product"`
	UnitAmount int64      `json:"unit_amount"`
	Recurring  *Recurring `json:"recurring,omitempty"`
}

type Recurring struct {
	Interval      string `json:"interval"`
	IntervalCount int    `json:"interval_count,omitempty"`
}

type CheckoutSession struct {
	CheckoutSessionID string  `json:"checkout_session_id"`
	Customer          *string `json:"customer"`
	RedirectURL       string  `json:"redirect_url"`
	Status            string  `json:"status"`
	Subscription      *string `json:"subscription"`
}

type checkoutSessionResponse struct {
	CheckoutSession *CheckoutSession `json:"checkout_session"`
}

func (c *Client) CreateCheckoutSession(ctx context.Context, req *CreateCheckoutSessionRequest) (*CheckoutSession, error) {
	var res checkoutSessionResponse
	if err := c.client.Do(ctx, "POST", "/v1/checkout/sessions", nil, req, &res); err != nil {
		return nil, err
	}
	return res.CheckoutSession, nil
}

func (c *Client) GetCheckoutSession(ctx context.Context, sessionID string) (*CheckoutSession, error) {
	var res checkoutSessionResponse
	if err := c.client.Do(ctx, "GET", fmt.Sprintf("/v1/checkout/sessions/%s", sessionID), nil, nil, &res); err != nil {
		return nil, err
	}
	return res.CheckoutSession, nil
}

type Subscription struct {
	SubscriptionID string             `json:"subscription_id"`
	Items          []SubscriptionItem `json:"items"`
}

type SubscriptionItem struct {
	Price *Price `json:"price"`
}

type Price struct {
	ID        string     `json:"id"`
	Recurring *Recurring `json:"recurring"`
}

type subscriptionResponse struct {
	Subscription *Subscription `json:"subscription"`
}

func (c *Client) GetSubscription(ctx context.Context, subscriptionID string) (*Subscription, error) {
	params := url.Values{"expand[]": []string{"items.price"}}
	var res subscriptionResponse
	if err := c.client.Do(ctx, "GET", fmt.Sprintf("/v1/subscriptions/%s", subscriptionID), params, nil, &res); err != nil {
		return nil, err
	}
	return res.Subscription, nil
}
