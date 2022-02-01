package kmsjwt

import (
	"context"
	"crypto/sha256"
	"encoding/asn1"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	kmsapi "cloud.google.com/go/kms/apiv1"
	jwt "github.com/golang-jwt/jwt/v4"
	jwtinterface "github.com/wearemojo/mojo-public-go/lib/jwt"
	"google.golang.org/api/iterator"
	kms "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

var _ jwtinterface.Signer = (*Signer)(nil)

type Signer struct {
	client      *kmsapi.KeyManagementClient
	projectID   string
	env         string
	serviceName string

	keyVersionCache      string
	keyVersionCacheSetAt time.Time
	keyVersionCacheLock  sync.RWMutex
}

func NewSigner(client *kmsapi.KeyManagementClient, projectID, env, serviceName string) *Signer {
	return &Signer{
		client:      client,
		projectID:   projectID,
		env:         env,
		serviceName: serviceName,
	}
}

// TODO: Go 1.18 generics
func (s *Signer) keyVersionCacheGet() (string, time.Time) {
	s.keyVersionCacheLock.RLock()
	defer s.keyVersionCacheLock.RUnlock()

	return s.keyVersionCache, s.keyVersionCacheSetAt
}

func (s *Signer) keyVersionCacheSet(keyVersion string) {
	s.keyVersionCacheLock.Lock()
	defer s.keyVersionCacheLock.Unlock()

	s.keyVersionCache = keyVersion
	s.keyVersionCacheSetAt = time.Now()
}

func (s *Signer) getKeyVersion(ctx context.Context) (string, error) {
	if cache, setAt := s.keyVersionCacheGet(); cache != "" {
		if time.Since(setAt) < time.Minute*5 {
			return cache, nil
		}
	}

	res, err := s.findKeyVersion(ctx)
	if err != nil {
		return "", err
	}

	s.keyVersionCacheSet(res)

	return res, nil
}

func (s *Signer) findKeyVersion(ctx context.Context) (string, error) {
	path := fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		s.projectID,
		"global",
		"services",
		fmt.Sprintf("%s-%s", s.env, s.serviceName),
	)

	res, err := s.client.ListCryptoKeyVersions(ctx, &kms.ListCryptoKeyVersionsRequest{
		Parent:   path,
		PageSize: 1,
		Filter:   "state=ENABLED",
		OrderBy:  "name desc",
	}).Next()
	if err == iterator.Done {
		return "", fmt.Errorf("no crypto key versions found")
	} else if err != nil {
		return "", err
	}

	if res.Algorithm != kms.CryptoKeyVersion_EC_SIGN_P256_SHA256 {
		return "", fmt.Errorf("unexpected algorithm: %s", res.Algorithm)
	}

	i := strings.LastIndex(res.Name, "/")
	displayName := res.Name[i+1:]

	return displayName, nil
}

func (s *Signer) Sign(ctx context.Context, customClaims jwtinterface.Claims) (string, error) {
	if _, ok := customClaims["v"].(string); !ok {
		return "", fmt.Errorf("version claim is required")
	}

	if _, ok := customClaims["t"].(string); !ok {
		return "", fmt.Errorf("type claim is required")
	}

	claims := jwtinterface.Claims{
		"iss": fmt.Sprintf("%s-%s", s.env, s.serviceName),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}

	for k, v := range customClaims {
		if _, ok := claims[k]; ok {
			return "", fmt.Errorf("claim %s cannot be overridden", k)
		}

		claims[k] = v
	}

	keyVersion, err := s.getKeyVersion(ctx)
	if err != nil {
		return "", err
	}

	sm := jwtSigningMethodSign{
		ctx:        ctx,
		signer:     s,
		keyVersion: keyVersion,
	}

	token := jwt.NewWithClaims(sm, claims)
	token.Header["kid"] = keyVersion

	return token.SignedString(ctx)
}

var _ jwt.SigningMethod = (*jwtSigningMethodSign)(nil)

type jwtSigningMethodSign struct {
	ctx        context.Context
	signer     *Signer
	keyVersion string
}

func (s jwtSigningMethodSign) Alg() string {
	return "ES256"
}

func (s jwtSigningMethodSign) Verify(signingString, signature string, key any) error {
	return fmt.Errorf("not implemented")
}

func (s jwtSigningMethodSign) Sign(signingString string, key any) (string, error) {
	path := fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%s",
		s.signer.projectID,
		"global",
		"services",
		fmt.Sprintf("%s-%s", s.signer.env, s.signer.serviceName),
		s.keyVersion,
	)

	h := sha256.New()
	h.Write([]byte(signingString))
	digest := h.Sum(nil)

	res, err := s.signer.client.AsymmetricSign(s.ctx, &kms.AsymmetricSignRequest{
		Name: path,
		Digest: &kms.Digest{
			Digest: &kms.Digest_Sha256{
				Sha256: digest,
			},
		},
	})
	if err != nil {
		return "", err
	}

	sig, err := reencodeSignature(res.Signature, jwt.SigningMethodES256)
	if err != nil {
		return "", err
	}

	return jwt.EncodeSegment(sig), nil
}

func reencodeSignature(sig []byte, method *jwt.SigningMethodECDSA) ([]byte, error) {
	var parsed struct{ R, S *big.Int }
	_, err := asn1.Unmarshal(sig, &parsed)
	if err != nil {
		return nil, err
	}

	keyBytes := method.CurveBits / 8
	if method.CurveBits%8 > 0 {
		keyBytes++
	}

	rBytes := parsed.R.Bytes()
	rBytesPadded := make([]byte, keyBytes)
	copy(rBytesPadded[keyBytes-len(rBytes):], rBytes)

	sBytes := parsed.S.Bytes()
	sBytesPadded := make([]byte, keyBytes)
	copy(sBytesPadded[keyBytes-len(sBytes):], sBytes)

	return append(rBytesPadded, sBytesPadded...), nil
}
