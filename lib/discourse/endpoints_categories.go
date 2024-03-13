package discourse

import (
	"context"
	"net/url"
)

type ListCategoriesOptions struct {
	IncludeSubcategories bool
}

func (c *IdentifiedClient) ListCategories(ctx context.Context, options *ListCategoriesOptions) (res *CategoryListResult, err error) {
	params := url.Values{}
	if options != nil && options.IncludeSubcategories {
		params.Set("include_subcategories", "true")
	}
	return res, c.client.Do(ctx, "GET", "/categories", params, nil, &res)
}
