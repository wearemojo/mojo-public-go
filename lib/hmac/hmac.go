package hmac

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/wearemojo/mojo-public-go/lib/secrets"
)

const keyLength = 32

type HMAC struct {
	secretID string
}

func New(ctx context.Context, secretID string) (*HMAC, error) {
	h := HMAC{secretID: secretID}

	if _, err := h.getSecret(ctx); err != nil {
		return nil, err
	}

	return &h, nil
}

func (h HMAC) getSecret(ctx context.Context) ([]byte, error) {
	keyHex, err := secrets.Get(ctx, h.secretID)
	if err != nil {
		return nil, err
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, err
	}

	if len(key) != keyLength {
		return nil, fmt.Errorf("invalid hmac key")
	}

	return key, nil
}

func (h HMAC) Generate(ctx context.Context, message string) (string, error) {
	key, err := h.getSecret(ctx)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))

	return hex.EncodeToString(mac.Sum(nil)), nil
}

func (h HMAC) Check(ctx context.Context, message, signature string) (bool, error) {
	correct, err := h.Generate(ctx, message)
	if err != nil {
		return false, err
	}

	return hmac.Equal([]byte(signature), []byte(correct)), nil
}
