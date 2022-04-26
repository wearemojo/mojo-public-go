package authparsing

import (
	"context"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

const ErrNoAuthorization = merr.Code("no_authorization")

type Handler func(context.Context, string) (context.Context, error)

type Parser struct {
	// keys must be lowercase
	Handlers map[string]Handler
}

func (a Parser) Check(ctx context.Context, authorizationHeader string) (context.Context, error) {
	if authorizationHeader == "" {
		return nil, ErrNoAuthorization
	}

	parts := strings.SplitN(authorizationHeader, " ", 2)
	if len(parts) != 2 {
		return nil, cher.New(cher.Unauthorized, nil, cher.New("invalid_authorization", nil))
	}
	kind, token := parts[0], parts[1]

	h, ok := a.Handlers[strings.ToLower(kind)]
	if !ok {
		return nil, cher.New(cher.Unauthorized, nil, cher.New("unknown_authorization_type", nil))
	}

	return h(ctx, token)
}
