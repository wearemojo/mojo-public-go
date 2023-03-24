//nolint:tagliatelle // datahappy uses camel case
package datahappy

import (
	"context"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/ptr"
)

type RequestType string

const (
	RequestTypeTrack  RequestType = "track"
	RequestTypeUpdate RequestType = "update"
	RequestTypeDelete RequestType = "delete"
)

type TrackRequest struct {
	Type         RequestType    `json:"type"`
	AnonymousID  string         `json:"anonymousId,omitempty"`
	UserID       string         `json:"userId,omitempty"`
	Event        string         `json:"event"`
	MessageID    string         `json:"messageId,omitempty"`
	Properties   map[string]any `json:"properties,omitempty"`
	Context      *Context       `json:"context,omitempty"`
	Integrations *Integrations  `json:"integrations,omitempty"`
	Timestamp    *time.Time     `json:"timestamp,omitempty"`
	Channel      string         `json:"channel,omitempty"`
	AuthToken    string         `json:"authToken"`
}

func (c *Client) Track(ctx context.Context, req *TrackRequest) error {
	// ensure no mutation of the original request
	req = ptr.ShallowCopy(req)

	if req.Type == "" {
		req.Type = RequestTypeTrack
	}

	if req.Context == nil {
		req.Context = &Context{}
	}

	if req.Context.Library == nil {
		req.Context = ptr.ShallowCopy(req.Context)
		req.Context.Library = library
	}

	if req.AuthToken == "" {
		req.AuthToken = c.AuthToken
	}

	return c.client.Do(ctx, "POST", "/v1/t", nil, req, nil)
}
