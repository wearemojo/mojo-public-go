package kmsjwt

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	jwt "github.com/golang-jwt/jwt/v4"
	jwtinterface "github.com/wearemojo/mojo-public-go/lib/jwt"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/ttlcache"
)

var _ jwtinterface.Verifier = (*Verifier)(nil)

type cacheKey struct {
	issuer, keyID string
}

type Verifier struct {
	client    *kms.KeyManagementClient
	projectID string

	publicKeyCache *ttlcache.KeyedCache[cacheKey, *ecdsa.PublicKey]
}

func NewVerifier(client *kms.KeyManagementClient, projectID string) *Verifier {
	return &Verifier{
		client:    client,
		projectID: projectID,

		publicKeyCache: ttlcache.NewKeyed[cacheKey, *ecdsa.PublicKey](time.Minute * 5),
	}
}

func (s *Verifier) getPublicKey(ctx context.Context, issuer, keyID string) (*ecdsa.PublicKey, error) {
	k := cacheKey{issuer, keyID}

	return s.publicKeyCache.GetOrDoE(k, func() (*ecdsa.PublicKey, error) {
		return s.findPublicKey(ctx, issuer, keyID)
	})
}

func (s *Verifier) findPublicKey(ctx context.Context, issuer, keyID string) (*ecdsa.PublicKey, error) {
	env, service, ok := strings.Cut(issuer, ";")
	if !ok {
		return nil, merr.New(ctx, "invalid_issuer", merr.M{"issuer": issuer})
	}

	path := fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%s",
		s.projectID,
		"global",
		"services",
		fmt.Sprintf("%s-%s", env, service),
		keyID,
	)

	res, err := s.client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{
		Name: path,
	})
	if err != nil {
		return nil, err
	}

	if res.Algorithm != kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256 {
		return nil, merr.New(ctx, "unexpected_crypto_key_algorithm", merr.M{"algorithm": res.Algorithm})
	}

	return jwt.ParseECPublicKeyFromPEM([]byte(res.Pem))
}

func (s *Verifier) Verify(ctx context.Context, token string) (claims jwtinterface.Claims, err error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"ES256"}),
		jwt.WithJSONNumber(),
	)
	_, err = parser.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		issuer, _ := claims["iss"].(string)
		keyID, _ := t.Header["kid"].(string)

		if issuer == "" || keyID == "" {
			return nil, merr.New(ctx, "missing_fields", merr.M{"iss": issuer, "kid": keyID})
		}

		return s.getPublicKey(ctx, issuer, keyID)
	})
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return nil, cher.New("token_malformed", nil)
	case errors.Is(err, jwt.ErrTokenUnverifiable):
		return nil, cher.New("token_unverifiable", nil)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, cher.New("token_bad_signature", nil)
	case errors.Is(err, jwt.ErrTokenExpired):
		return nil, cher.New("token_expired", nil)
	case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
		return nil, cher.New("token_used_before_issued", nil)
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, cher.New("token_not_yet_valid", nil)
	case err != nil:
		return nil, merr.New(ctx, "unknown_token_validation_error", nil, err)
	}

	return claims, nil
}
