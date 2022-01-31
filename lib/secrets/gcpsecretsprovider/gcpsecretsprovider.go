package gcpsecretsprovider

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/servicecontext"
	"github.com/wearemojo/mojo-public-go/lib/gcp"
	"github.com/wearemojo/mojo-public-go/lib/secrets"
	"google.golang.org/api/secretmanager/v1"
)

type cachedResult struct {
	retrievedAt time.Time

	data string
}

var _ secrets.Provider = (*GCPSecretsProvider)(nil)

type GCPSecretsProvider struct {
	projectID string

	cache      map[string]cachedResult
	cacheMutex sync.RWMutex
}

func New(ctx context.Context) (*GCPSecretsProvider, error) {
	projectID, err := gcp.GetProjectID(ctx)
	if err != nil {
		return nil, err
	}

	return &GCPSecretsProvider{
		projectID: projectID,
		cache:     map[string]cachedResult{},
	}, nil
}

// TODO: Go 1.18 generics
func (p *GCPSecretsProvider) cacheGet(key string) *cachedResult {
	p.cacheMutex.RLock()
	defer p.cacheMutex.RUnlock()

	if result, ok := p.cache[key]; ok {
		return &result
	}

	return nil
}

func (p *GCPSecretsProvider) cacheSet(key string, data string) {
	p.cacheMutex.Lock()
	defer p.cacheMutex.Unlock()

	p.cache[key] = cachedResult{
		retrievedAt: time.Now(),
		data:        data,
	}
}

func (p *GCPSecretsProvider) Get(ctx context.Context, secretID string) (string, error) {
	if result := p.cacheGet(secretID); result != nil {
		if time.Since(result.retrievedAt) < time.Minute {
			return result.data, nil
		}
	}

	result, err := p.load(ctx, secretID)
	if err != nil {
		return "", err
	}

	p.cacheSet(secretID, result)

	return result, nil
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
