package jwt

import (
	"context"
)

type contextKey string

const (
	contextKeySigner   contextKey = "signer"
	contextKeyVerifier contextKey = "verifier"
)

func ContextWithSigner(ctx context.Context, val Signer) context.Context {
	return context.WithValue(ctx, contextKeySigner, val)
}

func ContextWithVerifier(ctx context.Context, val Verifier) context.Context {
	return context.WithValue(ctx, contextKeyVerifier, val)
}

func ContextSigner(ctx context.Context) (val Signer) {
	val, _ = ctx.Value(contextKeySigner).(Signer)
	return
}

func ContextVerifier(ctx context.Context) (val Verifier) {
	val, _ = ctx.Value(contextKeyVerifier).(Verifier)
	return
}
