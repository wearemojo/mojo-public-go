package parseip

import (
	"context"

	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

// GetIP returns a v4 or v6 IP address from the remote address in an incoming
// http request.
func GetIP(ctx context.Context) (string, error) {
	rawReq := crpc.GetRequestContext(ctx)
	if rawReq == nil {
		return "", merr.New("ctx_missing_request", nil)
	}

	return StripPort(rawReq.RemoteAddr), nil
}
