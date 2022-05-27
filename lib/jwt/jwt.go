package jwt

import (
	"context"
	"time"
)

type Signer interface {
	Sign(ctx context.Context, expiresAt *time.Time, customClaims Claims) (token string, err error)
}

type Verifier interface {
	Verify(ctx context.Context, token string) (claims Claims, err error)
}
