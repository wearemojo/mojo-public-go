package mrpc

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

// adaptReqRes adapts a func(ctx, *Req) (Resp, error) to a HandlerFunc: it
// decodes the JSON request body into *Req, invokes fn, and JSON-encodes the
// response.
func adaptReqRes[Req, Resp any](fn func(ctx context.Context, req *Req) (Resp, error)) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		req, err := decodeRequest[Req](r)
		if err != nil {
			return err
		}

		res, err := fn(r.Context(), req)
		if err != nil {
			return err
		}

		return encodeResponse(w, res)
	}
}

// adaptReq adapts a func(ctx, *Req) error to a HandlerFunc. The handler
// produces no response body (204 No Content).
func adaptReq[Req any](fn func(ctx context.Context, req *Req) error) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		req, err := decodeRequest[Req](r)
		if err != nil {
			return err
		}

		if err := fn(r.Context(), req); err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}

// adaptRes adapts a func(ctx) (Resp, error) to a HandlerFunc. The handler
// takes no request body.
func adaptRes[Resp any](fn func(ctx context.Context) (Resp, error)) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := rejectRequestBody(r); err != nil {
			return err
		}

		res, err := fn(r.Context())
		if err != nil {
			return err
		}

		return encodeResponse(w, res)
	}
}

// adaptBare adapts a func(ctx) error to a HandlerFunc. The handler takes no
// request body and produces no response body (204 No Content).
func adaptBare(fn func(ctx context.Context) error) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := rejectRequestBody(r); err != nil {
			return err
		}

		if err := fn(r.Context()); err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}

// json/v2 marshals maps in a non-deterministic order by default, whereas v1
// sorted the keys - responses carry maps (error metadata, free-form fields),
// and a stable wire format keeps them diffable and cacheable
var deterministic = json.Deterministic(true)

func decodeRequest[Req any](r *http.Request) (*Req, error) {
	if r.Body == nil {
		return nil, cher.New(cher.BadRequest, nil, cher.New("missing_request_body", nil))
	}

	req := new(Req)

	err := json.UnmarshalDecode(jsontext.NewDecoder(r.Body), req)
	if errors.Is(err, io.EOF) {
		return nil, cher.New(cher.BadRequest, nil, cher.New("missing_request_body", nil))
	} else if err != nil {
		// This happens after we've run the JSON schema validation, so we already
		// know it's valid JSON in a structure we're expecting. If this does occur,
		// it means the Go struct doesn't match the JSON schema, so is an internal
		// issue, hence the use of merr.New() instead of cher.New().
		return nil, merr.New(r.Context(), "request_body_decode_failed", nil, err)
	}

	return req, nil
}

// rejectRequestBody enforces that a no-input handler received no request body.
func rejectRequestBody(r *http.Request) error {
	if r.Body == nil {
		return nil
	}

	buf := make([]byte, 1)
	i, err := r.Body.Read(buf)

	if i != 0 || !errors.Is(err, io.EOF) {
		return cher.New(cher.BadRequest, nil, cher.New("unexpected_request_body", nil))
	}

	return nil
}

func encodeResponse[Resp any](w http.ResponseWriter, res Resp) error {
	// jsontext does not escape HTML by default
	enc := jsontext.NewEncoder(w)

	if err := json.MarshalEncode(enc, res, deterministic); err != nil {
		if strings.Contains(err.Error(), "broken pipe") {
			return nil
		}

		return err
	}

	return nil
}
