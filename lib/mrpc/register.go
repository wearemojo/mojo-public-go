// Package mrpc is Mojo's RPC protocol: both the client and the server.
//
// The server is an HTTP handler that dispatches POST /<version>/<method>
// requests to explicitly-registered, per-version handlers, with JSON Schema
// validation, authorization enforcement, and cher-shaped error responses.
//
// The client (see Client) is composed into per-service typed clients.
package mrpc

import (
	"context"

	"github.com/wearemojo/mojo-public-go/lib/authenforce"
	"github.com/xeipuuv/gojsonschema"
)

// Register registers a handler taking a request body and returning a response
// body, at a single version.
func Register[Req, Resp any](s *Server, method, version string, schema gojsonschema.JSONLoader, en authenforce.Enforcers, fn func(ctx context.Context, req *Req) (Resp, error)) {
	s.register(method, version, schema, en, adaptReqRes(fn))
}

// RegisterNoRes registers a handler taking a request body and returning
// no response body (204 No Content), at a single version.
func RegisterNoRes[Req any](s *Server, method, version string, schema gojsonschema.JSONLoader, en authenforce.Enforcers, fn func(ctx context.Context, req *Req) error) {
	s.register(method, version, schema, en, adaptReq(fn))
}

// RegisterNoReq registers a handler taking no request body and returning
// a response body, at a single version.
func RegisterNoReq[Resp any](s *Server, method, version string, en authenforce.Enforcers, fn func(ctx context.Context) (Resp, error)) {
	s.register(method, version, nil, en, adaptRes(fn))
}

// RegisterNoReqRes registers a handler taking no request body and
// returning no response body (204 No Content), at a single version.
func RegisterNoReqRes(s *Server, method, version string, en authenforce.Enforcers, fn func(ctx context.Context) error) {
	s.register(method, version, nil, en, adaptBare(fn))
}
