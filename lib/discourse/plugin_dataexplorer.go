package discourse

import (
	"context"
	"fmt"
)

// https://github.com/discourse/discourse-data-explorer

type PluginDataExplorerRunQueryOptions struct {
	// https://meta.discourse.org/t/120063
	// https://github.com/discourse/discourse-data-explorer/blob/2f1044820c479424d29d94df389360b1d9dee871/app/controllers/discourse_data_explorer/query_controller.rb#L137

	Params   map[string]any
	Explain  bool
	Download bool
	LimitAll bool
}

type PluginDataExplorerQueryResult struct {
	// Errors indicates what went wrong with the query.
	//
	// At the time of writing, it is always empty when the request is successful.
	//
	// https://github.com/discourse/discourse-data-explorer/blob/2f1044820c479424d29d94df389360b1d9dee871/app/controllers/discourse_data_explorer/query_controller.rb#L193
	Errors []string `json:"errors"`

	Duration     float64                     `json:"duration"`
	ResultCount  int                         `json:"result_count"`
	Params       map[string]any              `json:"params"`
	Columns      []string                    `json:"columns"`
	DefaultLimit int                         `json:"default_limit"`
	Explain      *string                     `json:"explain"`
	Relations    map[string][]map[string]any `json:"relations"`
	ColRender    map[string]string           `json:"col_render"`
	Rows         [][]any                     `json:"rows"`
}

func (c *IdentifiedClient) PluginDataExplorerRunQuery(ctx context.Context, queryID int, options *PluginDataExplorerRunQueryOptions) (res *PluginDataExplorerQueryResult, err error) {
	path := fmt.Sprintf("/admin/plugins/explorer/queries/%d/run", queryID)
	body := map[string]any{}
	if options != nil {
		body["params"] = options.Params
		body["explain"] = options.Explain
		body["download"] = options.Download
		if options.LimitAll {
			body["limit"] = "ALL"
		}
	}
	return res, c.client.Do(ctx, "POST", path, nil, body, &res)
}
