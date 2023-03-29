package jwt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/slicefn"
)

type TypeVersion struct {
	Type    string
	Version string
}

func (p TypeVersion) String() string {
	return fmt.Sprintf("%s_%s", p.Type, p.Version)
}

func TypeVersionFromString(in string) (tv TypeVersion, ok bool) {
	tv.Type, tv.Version, ok = strings.Cut(in, "_")
	return
}

func Sign(ctx context.Context, expiresAt time.Time, customClaims Claims) (token string, err error) {
	signer := ContextSigner(ctx)
	if signer == nil {
		return "", merr.New(ctx, "ctx_missing_signer", nil)
	}

	return signer.Sign(ctx, expiresAt, customClaims)
}

func SignWithPrefix(ctx context.Context, expiresAt time.Time, customClaims Claims, typeVersion TypeVersion) (token string, err error) {
	claims := Claims{
		"t": typeVersion.Type,
		"v": typeVersion.Version,
	}

	for k, v := range customClaims {
		if _, ok := claims[k]; ok {
			return "", merr.New(ctx, "claim_unoverridable", merr.M{"claim": k})
		}

		claims[k] = v
	}

	token, err = Sign(ctx, expiresAt, claims)
	if err != nil {
		return
	}

	token = fmt.Sprintf("%s.%s", typeVersion, token)
	return
}

func Verify(ctx context.Context, token string) (claims Claims, err error) {
	verifier := ContextVerifier(ctx)
	if verifier == nil {
		err = merr.New(ctx, "ctx_missing_verifier", nil)
		return
	}

	return verifier.Verify(ctx, token)
}

func VerifyWithPrefix(ctx context.Context, token string, allowed []TypeVersion) (typeVersion TypeVersion, claims Claims, err error) {
	typeVersionStr, token, ok := strings.Cut(token, ".")
	if !ok {
		err = cher.New("missing_token_type_version", nil)
		return
	}

	if typeVersion, ok = TypeVersionFromString(typeVersionStr); !ok {
		err = cher.New("invalid_token_type_version", cher.M{"token_type_version": typeVersionStr})
		return
	}

	if _, ok = slicefn.Find(allowed, func(t TypeVersion) bool { return typeVersion == t }); !ok {
		err = cher.New("token_type_version_not_allowed", cher.M{"token_type_version": typeVersion})
		return
	}

	if claims, err = Verify(ctx, token); err != nil {
		return
	}

	if claims["t"] != typeVersion.Type || claims["v"] != typeVersion.Version {
		err = cher.New("token_type_version_mismatch", cher.M{"token_type_version": typeVersion})
		return
	}

	return typeVersion, claims, nil
}
