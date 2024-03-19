package discourse

import (
	"context"
	"fmt"
	"net/url"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

// https://github.com/discourse/discourse-reactions

type PluginDCReactionsReactionType string

const (
	PluginDCReactionsReactionTypeEmoji PluginDCReactionsReactionType = "emoji"
)

type PluginDCReactionsPostReaction struct {
	ID    string                        `json:"id"`
	Type  PluginDCReactionsReactionType `json:"type"`
	Count int                           `json:"count"`
}

type PluginDCReactionsPostReactionCurrentUser struct {
	ID   string                        `json:"id"`
	Type PluginDCReactionsReactionType `json:"type"`
}

func (c *IdentifiedClient) PluginDCReactionsToggleReaction(ctx context.Context, postID int, reactionID string) (res *Post, err error) {
	if reactionID == "" {
		return nil, merr.New(ctx, ErrEmptyParam, merr.M{"param": "reactionID"})
	}

	path := fmt.Sprintf("/discourse-reactions/posts/%d/custom-reactions/%s/toggle.json", postID, url.PathEscape(reactionID))
	return res, c.client.Do(ctx, "PUT", path, nil, nil, &res)
}
