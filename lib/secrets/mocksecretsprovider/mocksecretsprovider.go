package mocksecretsprovider

import (
	"context"
	"fmt"

	"github.com/wearemojo/mojo-public-go/lib/secrets"
)

var _ secrets.Provider = (*MockSecretsProvider)(nil)

type MockSecretsProvider struct {
	secrets map[string]string
}

func New(secrets map[string]string) *MockSecretsProvider {
	return &MockSecretsProvider{
		secrets: secrets,
	}
}

func (p MockSecretsProvider) Get(ctx context.Context, secretID string) (string, error) {
	if value, ok := p.secrets[secretID]; ok {
		return value, nil
	}

	return "", fmt.Errorf("secret_not_found")
}
