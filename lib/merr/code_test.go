package merr

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

type wrappedError struct {
	reason error
}

func (e wrappedError) Error() string {
	return "wrappedError"
}

func (e wrappedError) Unwrap() error {
	return e.reason
}

func TestCodeError(t *testing.T) {
	is := is.New(t)

	err := Code("foo")

	is.Equal(err.Error(), "foo")
}

func TestCodeComparable(t *testing.T) {
	is := is.New(t)

	err1 := Code("foo")
	err2 := Code("foo")

	is.True(err1 == err2)
}

func TestIsCode(t *testing.T) {
	is := is.New(t)

	errs := []error{
		Code("foo"),
		New("foo", nil),
		New(Code("foo"), nil),
		New("foo", M{"a": "b"}),
		New(Code("foo"), M{"a": "b"}),
		Wrap(New("foo", nil), "bar", nil),
		Wrap(wrappedError{New("foo", nil)}, "bar", nil),
	}

	for _, err := range errs {
		is.True(IsCode(err, "foo"))
		is.True(IsCode(err, Code("foo")))
		is.True(errors.Is(err, Code("foo")))
	}
}
