package discourse

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/url"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

type ConnectClient struct {
	client *Client

	connectSecret string
}

func NewConnectClient(client *Client, connectSecret string) *ConnectClient {
	return &ConnectClient{
		client: client,

		connectSecret: connectSecret,
	}
}

func (c *ConnectClient) SyncSSO(ctx context.Context, params url.Values) error {
	// https://meta.discourse.org/t/sync-discourseconnect-user-data-with-the-sync-sso-route/84398

	if params.Get("external_id") == "" {
		return merr.New(ctx, "missing_external_id", nil)
	}

	sso := base64.StdEncoding.EncodeToString([]byte(params.Encode()))

	hmac := hmac.New(sha256.New, []byte(c.connectSecret))
	hmac.Write([]byte(sso))
	sig := hex.EncodeToString(hmac.Sum(nil))

	req := map[string]string{
		"sso": sso,
		"sig": sig,
	}

	return c.client.systemClient().Do(ctx, "POST", "/admin/users/sync_sso", nil, req, nil)
}
