package merr

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

var _ interface {
	error
	Is(error) bool
	Unwrap() error
} = E{}

type E struct {
	Code Code `json:"code"`
	Meta M    `json:"meta"`

	Reason error `json:"reason"`
}

type M map[string]any

func New(code Code, meta M) E {
	return E{
		Code: code,
		Meta: meta,
	}
}

func Wrap(reason error, code Code, meta M) error {
	if reason == nil {
		return nil
	}

	return E{
		Code: code,
		Meta: meta,

		Reason: reason,
	}
}

func (e E) Equal(e2 E) bool {
	return e.Code == e2.Code &&
		e.Reason == e2.Reason &&
		cmp.Equal(e.Meta, e2.Meta)
}

// Error implements the error interface
//
// Provides a simple string representation of the error, but lacks some detail
//
// No compatibility guarantees are made with its output - it may change at any time
func (e E) Error() string {
	s := e.Code.Error()

	if len(e.Meta) > 0 {
		s += fmt.Sprintf(" (%v)", e.Meta)
	}

	if e.Reason != nil {
		s += fmt.Sprintf(": %v", e.Reason)
	}

	return s
}

// Is enables the use of `errors.Is`
func (e E) Is(err error) bool {
	if errors.Is(e.Code, err) {
		return true
	}

	// needed because E is not comparable
	if merr, ok := err.(E); ok {
		return e.Equal(merr)
	}

	return false
}

// Unwrap enables the use of `errors.Unwrap`
func (e E) Unwrap() error {
	return e.Reason
}
