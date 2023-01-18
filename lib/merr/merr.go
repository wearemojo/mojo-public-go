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

func New(ctx context.Context, code Code, meta M, reasons ...error) E {
	spanContext := trace.SpanContextFromContext(ctx)

	for _, reason := range reasons {
		if reason == nil {
			panic("merr.New: nil reason provided for " + code.String())
		}
	}

	// TODO: remove once Go 1.20 is released
	var reason error
	switch len(reasons) {
	case 0:
	case 1:
		reason = reasons[0]
	default:
		panic("merr.New: multiple reasons provided for " + code.String())
	}

	return E{
		Code: code,
		Meta: meta,

		TraceID: spanContext.TraceID(),
		SpanID:  spanContext.SpanID(),

		Stack: stacktrace.GetCallerFrames(2),

		Reason: reason,
	}
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
		str += fmt.Sprintf("\n- %v", e.Reason)
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
