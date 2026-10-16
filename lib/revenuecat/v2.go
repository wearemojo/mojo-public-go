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

const baseURLV2 = "https://api.revenuecat.com/v2"

// ClientV2 needs a v2 secret API key; the legacy v1 key isn't accepted.
type ClientV2 struct {
	client *jsonclient.Client
}

func NewClientV2(ctx context.Context, secretAPIKeySecretID string) (*ClientV2, error) {
	if _, err := secret.Get(ctx, secretAPIKeySecretID); err != nil {
		return nil, err
	}

	return &ClientV2{
		client: jsonclient.NewClient(
			baseURLV2,
			httpclient.NewClient(5*time.Second, roundTripper{secretAPIKeySecretID}),
		),
	}, nil
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type listProjectsResponse struct {
	Items []Project `json:"items"`
}

// https://www.revenuecat.com/docs/api-v2/project#tag/Project/operation/list-projects
// A v2 secret API key belongs to one project, so this returns just that project.
func (c *ClientV2) ListProjects(ctx context.Context) ([]Project, error) {
	var res listProjectsResponse
	if err := c.client.Do(ctx, "GET", "/projects", nil, nil, &res); err != nil {
		return nil, err
	}

	return res.Items, nil
}

type Product struct {
	ID              string           `json:"id"`
	StoreIdentifier string           `json:"store_identifier"`
	IndicativePrice *IndicativePrice `json:"indicative_price"`
}

type IndicativePrice struct {
	Currency     string `json:"currency"`
	Country      string `json:"country"`
	AmountMicros int64  `json:"amount_micros"`
}

type listProductsResponse struct {
	Items    []Product `json:"items"`
	NextPage *string   `json:"next_page"`
}

// https://www.revenuecat.com/docs/api-v2/product#tag/Product/operation/list-products
// Needs the `project_configuration:products:read` permission.
func (c *ClientV2) ListProducts(ctx context.Context, projectID string) ([]Product, error) {
	path := fmt.Sprintf("/projects/%s/products", url.PathEscape(projectID))
	params := url.Values{
		"limit":  {"100"},
		"expand": {"items.indicative_price"},
	}

	var products []Product
	for {
		var res listProductsResponse
		if err := c.client.Do(ctx, "GET", path, params, nil, &res); err != nil {
			return nil, err
		}

		products = append(products, res.Items...)
		if res.NextPage == nil || len(res.Items) == 0 {
			return products, nil
		}

		params.Set("starting_after", res.Items[len(res.Items)-1].ID)
	}
}
