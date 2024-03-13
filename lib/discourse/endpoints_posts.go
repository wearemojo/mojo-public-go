package discourse

import (
	"context"
	"fmt"
)

func (c *IdentifiedClient) CreatePost(ctx context.Context, req *PostNew) (res *Post, err error) {
	return res, c.client.Do(ctx, "POST", "/posts", nil, req, &res)
}

func (c *IdentifiedClient) GetPostByID(ctx context.Context, postID int) (res *Post, err error) {
	path := fmt.Sprintf("/posts/%d.json", postID)
	return res, c.client.Do(ctx, "GET", path, nil, nil, &res)
}
