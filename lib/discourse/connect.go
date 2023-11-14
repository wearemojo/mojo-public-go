package discourse

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/url"

	"github.com/wearemojo/mojo-public-go/lib/cher"
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

	sso, sig := c.StringifyAndSign(params)

	req := map[string]string{
		"sso": sso,
		"sig": sig,
	}

	return c.client.systemClient().Do(ctx, "POST", "/admin/users/sync_sso", nil, req, nil)
}

func (c *ConnectClient) sign(sso string) string {
	h := hmac.New(sha256.New, []byte(c.connectSecret))
	h.Write([]byte(sso))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *ConnectClient) StringifyAndSign(params url.Values) (sso, sig string) {
	sso = base64.StdEncoding.EncodeToString([]byte(params.Encode()))
	sig = c.sign(sso)

	return
}

func (c *ConnectClient) ParseAndVerify(ctx context.Context, sso, sig string) (url.Values, error) {
	expectedSig := c.sign(sso)

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return nil, cher.New("invalid_signature", nil)
	}

	bytes, err := base64.StdEncoding.DecodeString(sso)
	if err != nil {
		return nil, merr.New(ctx, "cannot_decode", nil, err)
	}

	params, err := url.ParseQuery(string(bytes))
	if err != nil {
		return nil, merr.New(ctx, "cannot_parse_query", nil, err)
	}

	return params, nil
}
