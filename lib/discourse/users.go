package discourse

import (
	"context"
	"fmt"
	"regexp"

	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/ksuid"
)

var usernameRegex = regexp.MustCompile(`^[\w.\-]+$`)

type UserResult struct {
	User *User `json:"user"`
}

type User struct {
	ID                 int                 `json:"id"`
	Username           string              `json:"username"`
	AvatarTemplate     *string             `json:"avatar_template"`
	SingleSignOnRecord *SingleSignOnRecord `json:"single_sign_on_record"`
}

type SingleSignOnRecord struct {
	ExternalID ksuid.ID `json:"external_id"`
}

func (c *Client) GetUserByID(ctx context.Context, userID int) (res *User, err error) {
	path := fmt.Sprintf("/admin/users/%d.json", userID)

	err = c.systemClient().Do(ctx, "GET", path, nil, nil, &res)
	return
}

func (c *Client) GetUserByUsername(ctx context.Context, username string) (res *User, err error) {
	if !usernameRegex.MatchString(username) {
		return nil, cher.New("invalid_username", cher.M{"username": username})
	}

	path := fmt.Sprintf("/users/%s", username)

	return c.getUserByPath(ctx, path)
}

func (c *Client) GetUserByExternalID(ctx context.Context, externalID string) (res *User, err error) {
	path := fmt.Sprintf("/users/by-external/%s", externalID)

	return c.getUserByPath(ctx, path)
}

func (c *Client) getUserByPath(ctx context.Context, path string) (res *User, err error) {
	var res2 UserResult
	err = c.systemClient().Do(ctx, "GET", path, nil, nil, &res2)
	res = res2.User
	return
}

func (c *Client) AnonymizeUser(ctx context.Context, userID int) error {
	path := fmt.Sprintf("/admin/users/%d/anonymize", userID)

	return c.systemClient().Do(ctx, "PUT", path, nil, nil, nil)
}
