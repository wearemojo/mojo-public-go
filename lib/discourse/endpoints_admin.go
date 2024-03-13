package discourse

import (
	"context"
	"fmt"
)

func (c *IdentifiedClient) AdminGetUserByID(ctx context.Context, userID int) (res *User, err error) {
	path := fmt.Sprintf("/admin/users/%d", userID)
	return res, c.client.Do(ctx, "GET", path, nil, nil, &res)
}

func (c *IdentifiedClient) AdminAnonymizeUser(ctx context.Context, userID int) error {
	path := fmt.Sprintf("/admin/users/%d/anonymize", userID)
	return c.client.Do(ctx, "PUT", path, nil, nil, nil)
}
