package hmac

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/wearemojo/mojo-public-go/lib/secret"
	"github.com/wearemojo/mojo-public-go/lib/secret/mocksecretprovider"
)

var ctx = secret.ContextWithProvider(context.Background(), mocksecretprovider.New(map[string]string{
	"valid":          "d182058da62f51e05a725774030d6e182f9cd1da05f0d67b00a806e5ae40102d",
	"invalid_length": "d182058da62f51e05a725774030d6e182f9cd1da05f0d67b00a806e5ae40102d00",
	"invalid_format": "zxcvbn",
}))

func TestNew(t *testing.T) {
	is := is.New(t)

	hmac, err := New(ctx, "valid")

	is.NoErr(err)
	is.True(hmac != nil)
}

func TestNewInvalidLength(t *testing.T) {
	is := is.New(t)

	hmac, err := New(ctx, "invalid_length")

	is.True(err != nil)
	is.Equal(hmac, nil)
}

func TestNewInvalidFormat(t *testing.T) {
	is := is.New(t)

	hmac, err := New(ctx, "invalid_format")

	is.True(err != nil)
	is.Equal(hmac, nil)
}

func TestGenerate(t *testing.T) {
	is := is.New(t)

	hmac, _ := New(ctx, "valid")
	hash, err := hmac.Generate(ctx, "test")

	is.NoErr(err)
	is.Equal(hash, "685b08ca6aa65f0e96627692da5230e5f48d508c49182c0bdba0a7e8ab866caa")
}

func TestCheck(t *testing.T) {
	is := is.New(t)

	hmac, _ := New(ctx, "valid")
	res, err := hmac.Check(ctx, "test", "685b08ca6aa65f0e96627692da5230e5f48d508c49182c0bdba0a7e8ab866caa")

	is.NoErr(err)
	is.Equal(res, true)
}

func TestCheckInvalid(t *testing.T) {
	is := is.New(t)

	hmac, _ := New(ctx, "valid")
	res, err := hmac.Check(ctx, "test", "685b08ca6aa65f0e96627692da5230e5f48d508c49182c0bdba0a7e8ab866cab")

	is.NoErr(err)
	is.Equal(res, false)
}
