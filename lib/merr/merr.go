package merr

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/wearemojo/mojo-public-go/lib/stacktrace"
)

var _ interface {
	error
	Is(error) bool
	Unwrap() error
} = E{}

// EInterface exists to allow `Wrap` to return nil
// without forcing us to use pointers for `E`
type EInterface interface {
	error
	GetConcrete() E
}

type E struct {
	Code Code `json:"code"`
	Meta M    `json:"meta"`

	Stack []stacktrace.Frame `json:"stack"`

	Reason error `json:"reason"`
}

type M map[string]any

func New(code Code, meta M) E {
	return E{
		Code: code,
		Meta: meta,

		Stack: stacktrace.GetCallerFrames(2),
	}
}

func Wrap(reason error, code Code, meta M) EInterface {
	if reason == nil {
		return nil
	}

	return E{
		Code: code,
		Meta: meta,

		Stack: stacktrace.GetCallerFrames(2),

		Reason: reason,
	}
}

func (e E) GetConcrete() E {
	return e
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
// Provides a simple string representation of the error, but lacks some detail
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

	// needed because E is not comparable
	if merr, ok := err.(E); ok { // nolint:errorlint
		return e.Equal(merr)
	}

	return false
}

// Unwrap enables the use of `errors.Unwrap`
func (e E) Unwrap() error {
	return e.Reason
}
