package mrpc

import (
	"encoding/json/v2"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

// errorResponseBody decides what every mrpc handler error becomes on the wire,
// and had no coverage before. These cases pin the mapping down, including that
// a cher.E is found through arbitrary wrapping rather than only at the top.
func TestErrorResponseBody(t *testing.T) {
	var target struct {
		A int `json:"a"`
	}

	semanticErr := json.Unmarshal([]byte(`{"a":"nope"}`), &target)
	require.Error(t, semanticErr)

	syntaxErr := json.Unmarshal([]byte(`{`), &target)
	require.Error(t, syntaxErr)

	t.Run("a cher.E passes straight through", func(t *testing.T) {
		body := errorResponseBody(cher.New("some_code", cher.M{"k": "v"}))

		require.Equal(t, "some_code", body.Code)
		require.Equal(t, cher.M{"k": "v"}, body.Meta)
	})

	t.Run("a wrapped cher.E is still found", func(t *testing.T) {
		ctx := t.Context()
		inner := merr.New(ctx, "inner", nil, cher.New("deep", cher.M{"n": 1}))
		body := errorResponseBody(merr.New(ctx, "outer", nil, inner))

		require.Equal(t, "deep", body.Code)
		require.Equal(t, cher.M{"n": 1}, body.Meta)
	})

	t.Run("a cher.E keeps its status code", func(t *testing.T) {
		body := errorResponseBody(cher.New(cher.BadRequest, nil))

		require.Equal(t, http.StatusBadRequest, body.StatusCode())
	})

	t.Run("a semantic json error becomes invalid_json with the pointer", func(t *testing.T) {
		body := errorResponseBody(semanticErr)

		require.Equal(t, "invalid_json", body.Code)
		require.Equal(t, "/a", body.Meta["name"])
	})

	t.Run("a syntactic json error becomes invalid_json with an offset", func(t *testing.T) {
		body := errorResponseBody(syntaxErr)

		require.Equal(t, "invalid_json", body.Code)
		require.Contains(t, body.Meta, "offset")
	})

	t.Run("anything else becomes unknown", func(t *testing.T) {
		body := errorResponseBody(merr.New(t.Context(), "not_a_cher_error", nil))

		require.Equal(t, cher.Unknown, body.Code)
		require.Nil(t, body.Meta)
	})
}
