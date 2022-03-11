package gcpsecretprovider

import (
	"context"
	"fmt"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/cuvva/cuvva-public-go/lib/servicecontext"
	"github.com/wearemojo/mojo-public-go/lib/gcp"
	"github.com/wearemojo/mojo-public-go/lib/secret"
	"github.com/wearemojo/mojo-public-go/lib/ttlcache"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var _ secret.Provider = (*GCPSecretProvider)(nil)

type GCPSecretProvider struct {
	projectID string

	cache *ttlcache.KeyedCache[string, string]
}

func New(ctx context.Context) (*GCPSecretProvider, error) {
	projectID, err := gcp.GetProjectID(ctx)
	if err != nil {
		return nil, err
	}

	return &GCPSecretProvider{
		projectID: projectID,

		cache: ttlcache.NewKeyed[string, string](time.Minute),
	}, nil
}

func (p *GCPSecretProvider) Get(ctx context.Context, secretID string) (string, error) {
	return p.cache.GetOrDoE(secretID, func() (string, error) {
		return p.load(ctx, secretID)
	})
}

func (p *GCPSecretProvider) load(ctx context.Context, secretID string) (secret string, err error) {
	sm, err := secretmanager.NewClient(ctx)
	if err != nil {
		return
	}

	env := servicecontext.GetContext(ctx).Environment
	res, err := sm.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s-%s/versions/latest", p.projectID, env, secretID),
	})
	if err != nil {
		return
	}

	return string(res.Payload.Data), nil
}
