package discourse

import (
	"context"
	"net/url"
)

type SearchResult struct {
	Posts               []SearchPost  `json:"posts"`
	Topics              []SearchTopic `json:"topics"`
	GroupedSearchResult struct {
		Error string `json:"error"`
	} `json:"grouped_search_result"`
}

type SearchPost struct {
	TopicID int    `json:"topic_id"`
	Blurb   string `json:"blurb"`
}

type SearchTopic struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	Title      string `json:"title"`
}

func (c *IdentifiedClient) Search(ctx context.Context, query string) (res *SearchResult, err error) {
	params := url.Values{"q": {query}}
	return res, c.client.Do(ctx, "GET", "/search.json", params, nil, &res)
}
