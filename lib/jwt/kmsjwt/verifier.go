package kmsjwt

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"
	"time"

	kmsapi "cloud.google.com/go/kms/apiv1"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	jwt "github.com/golang-jwt/jwt/v4"
	jwtinterface "github.com/wearemojo/mojo-public-go/lib/jwt"
	kms "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type cachedPublicKey struct {
	retrievedAt time.Time

	data *ecdsa.PublicKey
}

var _ jwtinterface.Verifier = (*Verifier)(nil)

type Verifier struct {
	client    *kmsapi.KeyManagementClient
	projectID string

	publicKeyCache     map[string]cachedPublicKey
	publicKeyCacheLock sync.RWMutex
}

func NewVerifier(client *kmsapi.KeyManagementClient, projectID string) *Verifier {
	return &Verifier{
		client:    client,
		projectID: projectID,

		publicKeyCache: map[string]cachedPublicKey{},
	}
}

// TODO: Go 1.18 generics
func (v *Verifier) publicKeyCacheGet(key string) *cachedPublicKey {
	v.publicKeyCacheLock.RLock()
	defer v.publicKeyCacheLock.RUnlock()

	if result, ok := v.publicKeyCache[key]; ok {
		return &result
	}

	return nil
}

func (v *Verifier) publicKeyCacheSet(key string, data *ecdsa.PublicKey) {
	v.publicKeyCacheLock.Lock()
	defer v.publicKeyCacheLock.Unlock()

	v.publicKeyCache[key] = cachedPublicKey{
		retrievedAt: time.Now(),
		data:        data,
	}
}

func (s *Verifier) getPublicKey(ctx context.Context, issuer, keyID string) (*ecdsa.PublicKey, error) {
	cacheKey := issuer + "/" + keyID

	if cache := s.publicKeyCacheGet(cacheKey); cache != nil {
		if time.Since(cache.retrievedAt) < time.Minute*5 {
			return cache.data, nil
		}
	}

	res, err := s.findPublicKey(ctx, issuer, keyID)
	if err != nil {
		return nil, err
	}

	s.publicKeyCacheSet(cacheKey, res)

	return res, nil
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
		return nil, fmt.Errorf("unexpected algorithm: %s", res.Algorithm)
	}

	return jwt.ParseECPublicKeyFromPEM([]byte(res.Pem))
}

func (s *Verifier) Verify(ctx context.Context, token string) (claims jwtinterface.Claims, err error) {
	parser := jwt.Parser{
		ValidMethods:  []string{"ES256"},
		UseJSONNumber: true,
	}
	_, err = parser.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		issuer, _ := claims["iss"].(string)
		keyID, _ := t.Header["kid"].(string)

		if issuer == "" || keyID == "" {
			return nil, fmt.Errorf("missing issuer or key ID")
		}

		return s.getPublicKey(ctx, issuer, keyID)
	})
	if vErr, ok := err.(*jwt.ValidationError); ok {
		switch {
		case vErr.Errors&jwt.ValidationErrorIssuedAt != 0:
			err = cher.New("token_used_before_issued", nil)
		case vErr.Errors&jwt.ValidationErrorNotValidYet != 0:
			err = cher.New("token_not_yet_valid", nil)
		case vErr.Errors&jwt.ValidationErrorExpired != 0:
			err = cher.New("token_expired", nil)
		}
	}
	return
}
