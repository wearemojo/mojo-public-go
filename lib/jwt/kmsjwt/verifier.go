package kmsjwt

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	kmsapi "cloud.google.com/go/kms/apiv1"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/wearemojo/mojo-public-go/lib/gerrors"
	jwtinterface "github.com/wearemojo/mojo-public-go/lib/jwt"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/ttlcache"
	kms "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

var _ jwtinterface.Verifier = (*Verifier)(nil)

type cacheKey struct {
	issuer, keyID string
}

type Verifier struct {
	client    *kmsapi.KeyManagementClient
	projectID string

	publicKeyCache *ttlcache.KeyedCache[cacheKey, *ecdsa.PublicKey]
}

func NewVerifier(client *kmsapi.KeyManagementClient, projectID string) *Verifier {
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
	path := fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%s",
		s.projectID,
		"global",
		"services",
		issuer,
		keyID,
	)

	res, err := s.client.GetPublicKey(ctx, &kms.GetPublicKeyRequest{
		Name: path,
	})
	if err != nil {
		return nil, err
	}

	if res.Algorithm != kms.CryptoKeyVersion_EC_SIGN_P256_SHA256 {
		return nil, merr.New("unexpected_crypto_key_algorithm", merr.M{"algorithm": res.Algorithm})
	}

	return jwt.ParseECPublicKeyFromPEM([]byte(res.Pem))
}

func (s *Verifier) Verify(ctx context.Context, token, allowedTypeVersion string) (jwtinterface.Claims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"ES256"}),
		jwt.WithJSONNumber(),
	)

	claims := jwtinterface.Claims{}

	_, err := parser.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		issuer, _ := claims["iss"].(string)
		tokenType, _ := claims["t"].(string)
		version, _ := claims["v"].(string)
		keyID, _ := t.Header["kid"].(string)

		if issuer == "" || keyID == "" || tokenType == "" || version == "" {
			return nil, merr.New("missing_fields", merr.M{"iss": issuer, "kid": keyID, "type": tokenType, "version": version})
		}

		return s.getPublicKey(ctx, issuer, keyID)
	})
	if vErr, ok := gerrors.As[*jwt.ValidationError](err); ok {
		switch {
		case vErr.Errors&jwt.ValidationErrorIssuedAt != 0:
			return nil, cher.New("token_used_before_issued", nil)
		case vErr.Errors&jwt.ValidationErrorNotValidYet != 0:
			return nil, cher.New("token_not_yet_valid", nil)
		case vErr.Errors&jwt.ValidationErrorExpired != 0:
			return nil, cher.New("token_expired", nil)
		}
	}

	tokenType, _ := claims["t"].(string)
	version, _ := claims["v"].(string)

	tokenTypeVersion := fmt.Sprintf("%s_%s", tokenType, version)

	if tokenTypeVersion != allowedTypeVersion {
		return nil, cher.New("token_type_version_mismatch", cher.M{
			"token_type_version":   tokenTypeVersion,
			"allowed_type_version": allowedTypeVersion,
		})
	}

	return claims, nil
}
