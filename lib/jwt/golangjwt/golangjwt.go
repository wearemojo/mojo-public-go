package golangjwt

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

func HandleVerifyError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return cher.New("token_malformed", nil)
	case errors.Is(err, jwt.ErrTokenUnverifiable):
		return cher.New("token_unverifiable", nil)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return cher.New("token_bad_signature", nil)
	case errors.Is(err, jwt.ErrTokenExpired):
		return cher.New("token_expired", nil)
	case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
		return cher.New("token_used_before_issued", nil)
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return cher.New("token_not_yet_valid", nil)
	default:
		return merr.New(ctx, "unknown_token_validation_error", nil, err)
	}
}
