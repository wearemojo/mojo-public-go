package postmark

import (
	"context"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
	"github.com/wearemojo/mojo-public-go/lib/httpclient"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/secret"
)

//nolint:tagliatelle // postmark uses title case
type Response struct {
	To          string    `json:"To"`
	SubmittedAt time.Time `json:"SubmittedAt"`
	MessageID   string    `json:"MessageID"`
	ErrorCode   int       `json:"ErrorCode"`
	Message     string    `json:"Message"`
}

//nolint:tagliatelle // postmark uses title case
type EmailWithTemplate struct {
	MessageStream string            `json:"MessageStream"`
	From          string            `json:"From"`
	To            string            `json:"To"`
	TemplateAlias string            `json:"TemplateAlias"`
	TemplateModel map[string]string `json:"TemplateModel"`
}

type Client struct {
	BaseURL string

	secretID string
}

func NewClient(ctx context.Context, baseURL, secretID string) (*Client, error) {
	if _, err := secret.Get(ctx, secretID); err != nil {
		return nil, err
	}

	return &Client{
		BaseURL: baseURL,

		secretID: secretID,
	}, nil
}

func (c *Client) client(ctx context.Context) (*jsonclient.Client, error) {
	apiKey, err := secret.Get(ctx, c.secretID)
	if err != nil {
		return nil, err
	}

	return jsonclient.NewClient(
		c.BaseURL,
		httpclient.NewClient(5*time.Second, roundTripper{apiKey}),
	), nil
}

func (c *Client) SendWithTemplate(ctx context.Context, req *EmailWithTemplate) (res *Response, err error) {
	jsonClient, err := c.client(ctx)
	if err != nil {
		return nil, err
	}

	if err = jsonClient.Do(ctx, "POST", "email/withTemplate", nil, req, &res); err != nil {
		return
	}

	if res.ErrorCode != 0 {
		err = merr.New(ctx, "postmark_error", merr.M{
			"error_code": res.ErrorCode,
			"message":    res.Message,
			"message_id": res.MessageID,
		})
	}

	return
}
