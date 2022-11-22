package merr

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/wearemojo/mojo-public-go/lib/stacktrace"
	"go.opentelemetry.io/otel/trace"
)

var _ interface {
	error
	Is(error) bool
	Unwrap() error
} = E{}

// Merrer (merr-er) represents a merr-compatible error
//
// It primarily exists to allow `Wrap` to return nil without forcing us to use
// pointers for `E`, but also allows other structs to offer a merr.E option
type Merrer interface {
	error
	Merr() E
}

type E struct {
	Code Code `json:"code"`
	Meta M    `json:"meta"`

	TraceID trace.TraceID `json:"trace_id"`
	SpanID  trace.SpanID  `json:"span_id"`

	Stack []stacktrace.Frame `json:"stack"`

	Reason error `json:"reason"`
}

type M map[string]any

func newE(ctx context.Context, reason error, code Code, meta M) E {
	spanContext := trace.SpanContextFromContext(ctx)

	return E{
		Code: code,
		Meta: meta,

		TraceID: spanContext.TraceID(),
		SpanID:  spanContext.SpanID(),

		Stack: stacktrace.GetCallerFrames(3),

		Reason: reason,
	}
}

func New(ctx context.Context, code Code, meta M) E {
	return newE(ctx, nil, code, meta)
}

// TODO: replace in favor of variadic reasons on `New` in Go 1.20
func Wrap(ctx context.Context, reason error, code Code, meta M) Merrer {
	if reason == nil {
		return nil
	}

	return newE(ctx, reason, code, meta)
}

func (e E) Merr() E {
	return e
}

func (e E) Fields() map[string]any {
	return map[string]any{
		"code":   e.Code,
		"meta":   e.Meta,
		"stack":  e.Stack,
		"reason": e.Reason,
	}
}

func (e E) Equal(e2 E) bool {
	return e.Code == e2.Code &&
		cmp.Equal(e.Meta, e2.Meta) &&
		cmp.Equal(e.Stack, e2.Stack) &&
		cmp.Equal(e.Reason, e2.Reason)
}

func (e E) String() string {
	return e.Error() + "\n\n" + stacktrace.FormatFrames(e.Stack)
}

// Error implements the error interface
//
// Provides a simple string representation of the error, but cannot include the
// complete data contained in the error
//
// No compatibility guarantees are made with its output - it may change at any time
func (e E) Error() string {
	str := string(e.Code)

	if len(e.Meta) > 0 {
		str += fmt.Sprintf(" (%v)", e.Meta)
	}

	if e.Reason != nil {
		str += fmt.Sprintf(": %v", e.Reason)
	}

	return str
}

// Is enables the use of `errors.Is`
func (e E) Is(err error) bool {
	if errors.Is(e.Code, err) {
		return true
	}

	//nolint:errorlint // needed because E is not comparable
	if merr, ok := err.(E); ok {
		return e.Equal(merr)
	}

	return false
}

// Unwrap enables the use of `errors.Unwrap`
func (e E) Unwrap() error {
	return e.Reason
}
