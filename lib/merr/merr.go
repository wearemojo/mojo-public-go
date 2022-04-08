package merr

import (
	"errors"
	"fmt"
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
//
// Should not be called directly, because direct comparison and unwrapping is left to the caller
func (e E) Is(err error) bool {
	return errors.Is(e.Code, err)
}

func (e E) Unwrap() error {
	return e.Reason
}
