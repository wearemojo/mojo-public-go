package discourse

import (
	"context"
)

type PostNew struct {
	TopicID int `json:"topic_id"`

	Raw string `json:"raw"`
}

type Post struct {
	PostNew

	ID int `json:"id"`
}

func (c *Client) CreatePost(ctx context.Context, username string, req *PostNew) (res *Post, err error) {
	return res, c.usernameClient(username).Do(ctx, "POST", "/posts", nil, req, &res)
}
