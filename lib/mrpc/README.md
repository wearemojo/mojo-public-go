# mrpc

mrpc is Mojo's RPC protocol — one package holding both the client and the
server, plus shared helpers.

It is heavily influenced by [net/rpc](https://golang.org/pkg/net/rpc/) and
[Monzo Typhon](https://godoc.org/github.com/monzo/typhon).

## Server

`NewServer` returns an `http.Handler` that dispatches `POST /<version>/<method>`
to handlers registered explicitly per version via `Register` / `RegisterNoReq`
/ `RegisterNoRes` / `RegisterNoReqRes`. It applies JSON Schema validation,
authorization enforcement, and cher-shaped error responses. There is no implied
inheritance between versions — each version lists exactly the methods it serves.

## Client

`Client` is not used directly; it is composed into a per-service typed client
(see the generated `svc/*/client.go`).

## CoerceJSONSchemaError

A helper that converts a JSON Schema validation result into a cher
`bad_request` error with one `schema_failure` reason per validation error.
