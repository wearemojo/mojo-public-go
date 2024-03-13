package discourse

import (
	"context"
	"fmt"
	"net/url"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

func (c *IdentifiedClient) GetUserByUsername(ctx context.Context, username string) (res *UserResult, err error) {
	if username == "" {
		return nil, merr.New(ctx, ErrEmptyParam, merr.M{"param": "username"})
	}

	path := fmt.Sprintf("/users/%s", url.PathEscape(username))
	return res, c.client.Do(ctx, "GET", path, nil, nil, &res)
}

func (c *IdentifiedClient) GetUserByExternalID(ctx context.Context, externalID string) (res *UserResult, err error) {
	if externalID == "" {
		return nil, merr.New(ctx, ErrEmptyParam, merr.M{"param": "externalID"})
	}

	path := fmt.Sprintf("/users/by-external/%s", url.PathEscape(externalID))
	return res, c.client.Do(ctx, "GET", path, nil, nil, &res)
}
