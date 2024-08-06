package crpc

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/wearemojo/mojo-public-go/lib/gerrors"
	"github.com/wearemojo/mojo-public-go/lib/jsonclient"
	"github.com/wearemojo/mojo-public-go/lib/servicecontext"
	"github.com/wearemojo/mojo-public-go/lib/version"
)

const (
	userAgentTemplate            = "crpc/%s (+https://github.com/wearemojo/mojo-public-go/tree/main/lib/crpc)"
	userAgentTemplateWithService = "crpc/%s (+https://github.com/wearemojo/mojo-public-go/tree/main/lib/crpc) [%s/%s]"
)

// Client represents a crpc client. It builds on top of jsonclient, so error
// variables/structs and the authenticated round tripper live there.
type Client struct {
	client *jsonclient.Client
}

// NewClient returns a client configured with a transport scheme, remote host
// and URL prefix supplied as a URL <scheme>://<host></prefix>
func NewClient(ctx context.Context, baseURL string, c *http.Client) *Client {
	jcc := jsonclient.NewClient(baseURL, c)

	svc := servicecontext.GetContext(ctx)
	if svc != nil {
		jcc.UserAgent = fmt.Sprintf(userAgentTemplateWithService, version.Truncated, svc.Service, svc.Env)
	} else {
		jcc.UserAgent = fmt.Sprintf(userAgentTemplate, version.Truncated)
	}

	return &Client{jcc}
}

// Do executes an RPC request against the configured server.
func (c *Client) Do(ctx context.Context, method, version string, src, dst any) error {
	err := c.client.Do(ctx, "POST", path.Join(version, method), nil, src, dst)

	if err == nil {
		return nil
	}

	if err, ok := gerrors.As[jsonclient.ClientTransportError](err); ok {
		return ClientTransportError{method, version, err.ErrorString, err.Cause()}
	}

	return err
}

// ClientTransportError is returned when an error related to
// executing a client request occurs.
type ClientTransportError struct {
	Method, Version, ErrorString string

	cause error
}

// Cause returns the causal error (if wrapped) or nil
func (cte ClientTransportError) Cause() error {
	return cte.cause
}

func (cte ClientTransportError) Error() string {
	if cte.cause != nil {
		return fmt.Sprintf("%s/%s %s: %s", cte.Version, cte.Method, cte.ErrorString, cte.cause.Error())
	}

	return fmt.Sprintf("%s/%s %s", cte.Version, cte.Method, cte.ErrorString)
}
