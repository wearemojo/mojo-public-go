package discourse

import (
	"context"
	"fmt"
	"time"
)

type PostNew struct {
	TopicID int `json:"topic_id"`

	Raw string `json:"raw"`
}

type Post struct {
	PostNew

	ID        int       `json:"id"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	UserID    int       `json:"user_id"`
}

func (c *Client) CreatePost(ctx context.Context, username string, req *PostNew) (res *Post, err error) {
	return res, c.usernameClient(username).Do(ctx, "POST", "/posts", nil, req, &res)
}

func (c *Client) GetPostByID(ctx context.Context, postID int) (res *Post, err error) {
	path := fmt.Sprintf("/posts/%d.json", postID)
	return res, c.systemClient().Do(ctx, "GET", path, nil, nil, &res)
}
