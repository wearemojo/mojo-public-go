package kmsjwt

import (
	"context"
	"crypto/sha256"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/chebyrash/promise"
	jwt "github.com/golang-jwt/jwt/v5"
	jwtinterface "github.com/wearemojo/mojo-public-go/lib/jwt"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/ttlcache"
	"google.golang.org/api/iterator"
)

var _ jwtinterface.Signer = (*Signer)(nil)

type Signer struct {
	client      *kms.KeyManagementClient
	projectID   string
	env         string
	serviceName string

	keyVersionCache *ttlcache.SingularCache[string]
	promise         *promise.Promise[string]
	mutex           sync.Mutex
}

func NewSigner(client *kms.KeyManagementClient, projectID, env, serviceName string) *Signer {
	return &Signer{
		client:      client,
		projectID:   projectID,
		env:         env,
		serviceName: serviceName,

		keyVersionCache: ttlcache.NewSingular[string](time.Minute * 5),
	}
}

func (s *Signer) getKeyVersion(ctx context.Context) (string, error) {
	return s.keyVersionCache.GetOrDoE(func() (string, error) {
		s.setPromise(ctx)
		defer func() { s.promise = nil }()

		keyVersion, err := s.promise.Await(ctx)
		if err != nil {
			return "", err
		}

		return *keyVersion, nil
	})
}

func (s *Signer) setPromise(ctx context.Context) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.promise == nil {
		s.promise = promise.New(func(resolve func(string), reject func(error)) {
			keyVersion, err := s.findKeyVersion(ctx)
			if err != nil {
				reject(err)
			}

			resolve(keyVersion)
		})
	}
}

func (s *Signer) findKeyVersion(ctx context.Context) (string, error) {
	path := fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		s.projectID,
		"global",
		"services",
		fmt.Sprintf("%s-%s", s.env, s.serviceName),
	)

	res, err := s.client.ListCryptoKeyVersions(ctx, &kmspb.ListCryptoKeyVersionsRequest{
		Parent:   path,
		PageSize: 1,
		Filter:   "state=ENABLED",
		OrderBy:  "name desc",
	}).Next()
	if errors.Is(err, iterator.Done) {
		return "", merr.New(ctx, "missing_crypto_key_version", merr.M{"path": path})
	} else if err != nil {
		return "", err
	}

	if res.Algorithm != kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256 {
		return "", merr.New(ctx, "unexpected_crypto_key_algorithm", merr.M{"algorithm": res.Algorithm})
	}

	i := strings.LastIndex(res.Name, "/")
	displayName := res.Name[i+1:]

	return displayName, nil
}

func (s *Signer) Sign(ctx context.Context, expiresAt time.Time, customClaims jwtinterface.Claims) (string, error) {
	if _, ok := customClaims["v"].(string); !ok {
		return "", merr.New(ctx, "required_claim_missing", merr.M{"claim": "v"})
	}

	if _, ok := customClaims["t"].(string); !ok {
		return "", merr.New(ctx, "required_claim_missing", merr.M{"claim": "t"})
	}

	claims := jwtinterface.Claims{
		"iss": fmt.Sprintf("%s;%s", s.env, s.serviceName),
		"iat": time.Now().Unix(),
		"exp": expiresAt.Unix(),
	}

	for k, v := range customClaims {
		if _, ok := claims[k]; ok {
			return "", merr.New(ctx, "claim_unoverridable", merr.M{"claim": k})
		}

		claims[k] = v
	}

	keyVersion, err := s.getKeyVersion(ctx)
	if err != nil {
		return "", err
	}

	signingMethod := jwtSigningMethodSign{
		ctx:        ctx,
		signer:     s,
		keyVersion: keyVersion,
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	token.Header["kid"] = keyVersion

	return token.SignedString(ctx)
}

var _ jwt.SigningMethod = (*jwtSigningMethodSign)(nil)

type jwtSigningMethodSign struct {
	//nolint:containedctx // we have to conform to jwt.SigningMethod
	ctx        context.Context
	signer     *Signer
	keyVersion string
}

func (s jwtSigningMethodSign) Alg() string {
	return "ES256"
}

func (s jwtSigningMethodSign) Verify(signingString string, signature []byte, key any) error {
	return merr.New(s.ctx, "not_implemented", nil)
}

func (s jwtSigningMethodSign) Sign(signingString string, key any) ([]byte, error) {
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

	res, err := s.signer.client.AsymmetricSign(s.ctx, &kmspb.AsymmetricSignRequest{
		Name: path,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return reencodeSignature(res.Signature, jwt.SigningMethodES256)
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

	//nolint:makezero // this is correct for this use case
	return append(rBytesPadded, sBytesPadded...), nil
}
