package jwt

import (
	"context"
)

type Signer interface {
	Sign(context.Context, Claims) (string, error)
}

type Verifier interface {
	Verify(context.Context, string) (Claims, error)
}
