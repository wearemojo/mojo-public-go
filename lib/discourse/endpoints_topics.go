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

func (c *IdentifiedClient) GetTopic(ctx context.Context, topicID int, options *GetTopicOptions) (res *TopicResult, err error) {
	path := fmt.Sprintf("/t/%d", topicID)
	params := url.Values{}
	if options != nil && options.Print {
		params.Set("print", "true")
	}
	return res, c.client.Do(ctx, "GET", path, params, nil, &res)
}

type ListTopicPostIDsOptions struct {
	PostNumber int
}

func (c *IdentifiedClient) ListTopicPostIDs(ctx context.Context, topicID int, options *ListTopicPostIDsOptions) (res *PostIDsResult, err error) {
	path := fmt.Sprintf("/t/%d/post_ids", topicID)
	params := url.Values{}
	if options != nil {
		params.Set("post_number", strconv.Itoa(options.PostNumber))
	}
	return res, c.client.Do(ctx, "GET", path, params, nil, &res)
}

func (c *IdentifiedClient) ListTopicPostsByIDs(ctx context.Context, topicID int, postIDs []int) (res *PostStreamResult, err error) {
	path := fmt.Sprintf("/t/%d/posts", topicID)
	params := url.Values{"post_ids[]": slicefn.Map(postIDs, strconv.Itoa)}
	return res, c.client.Do(ctx, "GET", path, params, nil, &res)
}
