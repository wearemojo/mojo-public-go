package authparsing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

func TestJSONError(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		Name   string
		Err    error
		Status int
		Body   string
	}{
		{
			"cher",
			cher.New(cher.Unauthorized, nil, cher.New("invalid_authorization", nil)),
			http.StatusUnauthorized,
			`{"code":"unauthorized","reasons":[{"code":"invalid_authorization"}]}`,
		},
		{
			// the error is recorded by the caller, not echoed to the client
			"not cher",
			merr.New(ctx, "some_internal_failure", nil),
			http.StatusInternalServerError,
			`{"code":"unknown"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			is := is.New(t)

			rec := httptest.NewRecorder()
			jsonError(ctx, rec, test.Err)

			is.Equal(test.Status, rec.Code)
			is.Equal(test.Body+"\n", rec.Body.String())
			is.Equal("application/json; charset=utf-8", rec.Header().Get("Content-Type"))
		})
	}
}
