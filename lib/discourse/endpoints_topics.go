package discourse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/wearemojo/mojo-public-go/lib/slicefn"
)

type GetTopicOptions struct {
	Print bool
}

func (c *IdentifiedClient) GetTopic(ctx context.Context, topicID int, options *GetTopicOptions) (res *PostStreamResult, err error) {
	path := fmt.Sprintf("/t/%d", topicID)
	params := url.Values{}
	if options != nil && options.Print {
		params.Set("print", "true")
	}
	return res, c.client.Do(ctx, "GET", path, params, nil, &res)
}

func (c *IdentifiedClient) ListTopicPostIDs(ctx context.Context, topicID int) (res *PostIDsResult, err error) {
	path := fmt.Sprintf("/t/%d/posts_ids", topicID)
	params := url.Values{"post_number": {"0"}}
	return res, c.client.Do(ctx, "GET", path, params, nil, &res)
}

func (c *IdentifiedClient) ListTopicPostsByIDs(ctx context.Context, topicID int, postIDs []int) (res *PostStreamResult, err error) {
	path := fmt.Sprintf("/t/%d/posts", topicID)
	params := url.Values{"post_ids[]": slicefn.Map(postIDs, strconv.Itoa)}
	return res, c.client.Do(ctx, "GET", path, params, nil, &res)
}
