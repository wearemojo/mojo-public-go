package postmark

import (
	"context"
	"fmt"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/gerrors"
	"github.com/wearemojo/mojo-public-go/lib/httpclient"
	"github.com/wearemojo/mojo-public-go/lib/jsonclient"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/secret"
)

const baseURL = "https://api.postmarkapp.com"

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

func (c *Client) SendWithTemplate(ctx context.Context, req *EmailWithTemplate) (res *Response, err error) {
	err = c.client.Do(ctx, "POST", "/email/withTemplate", nil, req, &res)
	if cerr, ok := gerrors.As[cher.E](err); ok {
		cerr.Code = fmt.Sprintf("postmark_%s", cerr.Code)
		return nil, cerr
	} else if err != nil {
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
