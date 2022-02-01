package gcpsecretsprovider

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/servicecontext"
	"github.com/wearemojo/mojo-public-go/lib/gcp"
	"github.com/wearemojo/mojo-public-go/lib/secrets"
	"github.com/wearemojo/mojo-public-go/lib/ttlcache"
	"google.golang.org/api/secretmanager/v1"
)

var _ secrets.Provider = (*GCPSecretsProvider)(nil)

type GCPSecretsProvider struct {
	projectID string

	cache *ttlcache.KeyedCache[string]
}

func New(ctx context.Context) (*GCPSecretsProvider, error) {
	projectID, err := gcp.GetProjectID(ctx)
	if err != nil {
		return nil, err
	}

	return &GCPSecretsProvider{
		projectID: projectID,

		cache: ttlcache.NewKeyed[string](time.Minute),
	}, nil
}

func (p *GCPSecretsProvider) Get(ctx context.Context, secretID string) (string, error) {
	return p.cache.GetOrDoE(secretID, func() (string, error) {
		return p.load(ctx, secretID)
	})
}

func (p *GCPSecretsProvider) load(ctx context.Context, secretID string) (secret string, err error) {
	sm, err := secretmanager.NewService(ctx)
	if err != nil {
		return
	}

	env := servicecontext.Get().Environment
	path := fmt.Sprintf("projects/%s/secrets/%s-%s/versions/latest", p.projectID, env, secretID)
	s, err := sm.Projects.Secrets.Versions.Access(path).Context(ctx).Do()
	if err != nil {
		return
	}

	data, err := base64.StdEncoding.DecodeString(s.Payload.Data)
	if err != nil {
		return
	}

	return string(data), nil
}
