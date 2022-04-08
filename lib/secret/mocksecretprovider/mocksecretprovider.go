package mocksecretprovider

import (
	"context"

	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/secret"
)

var _ secret.Provider = (*MockSecretProvider)(nil)

type MockSecretProvider struct {
	secrets map[string]string
}

func New(secrets map[string]string) *MockSecretProvider {
	return &MockSecretProvider{
		secrets: secrets,
	}
}

func (p MockSecretProvider) Get(ctx context.Context, secretID string) (string, error) {
	if value, ok := p.secrets[secretID]; ok {
		return value, nil
	}

	return "", merr.New(secret.ErrSecretNotFound, merr.M{"secret_id": secretID})
}
