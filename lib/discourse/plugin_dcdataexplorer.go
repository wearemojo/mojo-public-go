package discourse

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/slicefn"
)

// https://github.com/discourse/discourse-data-explorer

type PluginDCDataExplorerRunQueryOptions struct {
	// https://meta.discourse.org/t/120063
	// https://github.com/discourse/discourse-data-explorer/blob/2f1044820c479424d29d94df389360b1d9dee871/app/controllers/discourse_data_explorer/query_controller.rb#L137

	Params   map[string]any
	Explain  bool
	Download bool
	LimitAll bool
}

type PluginDCDataExplorerQueryResult struct {
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
	ColRender    map[string]string           `json:"colrender"`
	Rows         [][]any                     `json:"rows"`
}

func PluginDCDataExplorerQueryResultUnmarshal[T any](ctx context.Context, res *PluginDCDataExplorerQueryResult) ([]T, error) {
	colLen := len(res.Columns)

	mapped, err := slicefn.MapE(res.Rows, func(row []any) (map[string]any, error) {
		if len(row) != colLen {
			return nil, merr.New(ctx, "column_count_mismatch", merr.M{"expected": colLen, "actual": len(row)})
		}

		out := map[string]any{}
		for i, col := range res.Columns {
			out[col] = row[i]
		}
		return out, nil
	})
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(mapped)
	if err != nil {
		return nil, merr.New(ctx, "json_marshal_failed", nil, err)
	}

	var out []T
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, merr.New(ctx, "json_unmarshal_failed", nil, err)
	}

	return out, nil
}

func (c *IdentifiedClient) PluginDCDataExplorerRunQuery(ctx context.Context, queryID int, options *PluginDCDataExplorerRunQueryOptions) (res *PluginDCDataExplorerQueryResult, err error) {
	path := fmt.Sprintf("/admin/plugins/discourse-user-data-explorer/queries/%d/run", queryID)
	body := map[string]any{}
	if options != nil {
		body["explain"] = strconv.FormatBool(options.Explain)
		body["download"] = options.Download
		if options.Params != nil {
			data, err := json.Marshal(options.Params)
			if err != nil {
				return nil, merr.New(ctx, "json_marshal_failed", nil, err)
			}
			body["params"] = string(data)
		}
		if options.LimitAll {
			body["limit"] = "ALL"
		}
	}
	return res, c.client.Do(ctx, "POST", path, nil, body, &res)
}
