package merr

import (
	"errors"
	"testing"

	"github.com/matryer/is"
	"github.com/wearemojo/mojo-public-go/lib/gerrors"
)

func TestNew(t *testing.T) {
	is := is.New(t)

	err := New("foo", nil)

	is.Equal(err, E{
		Code: Code("foo"),
	})
}

func TestNewMeta(t *testing.T) {
	is := is.New(t)

	err := New("foo", M{"a": "b"})

	is.Equal(err, E{
		Code: Code("foo"),
		Meta: M{"a": "b"},
	})
}

func TestWrap(t *testing.T) {
	is := is.New(t)

	err1 := errors.New("underlying error")
	err2 := Wrap(err1, "foo", nil)

	is.Equal(err2, E{
		Code:   Code("foo"),
		Reason: errors.New("underlying error"),
	})
}

func TestWrapMeta(t *testing.T) {
	is := is.New(t)

	err1 := errors.New("underlying error")
	err2 := Wrap(err1, "foo", M{"a": "b"})

	is.Equal(err2, E{
		Code: Code("foo"),
		Meta: M{"a": "b"},

		Reason: errors.New("underlying error"),
	})
}

func TestWrapNil(t *testing.T) {
	is := is.New(t)

	err := Wrap(nil, "foo", nil)

	is.Equal(err, nil)
}

func TestEqual(t *testing.T) {
	is := is.New(t)

	err1 := New("foo", M{"a": "b"})
	err2 := New("foo", M{"a": "b"})
	err3 := New("foo", M{"a": "c"})
	err4 := New("bar", M{"a": "b"})
	err5 := New("foo", nil)
	err6 := New("foo", M{"a": "b"})
	err6.Reason = errors.New("foo")

	is.True(err1.Equal(err2))
	is.True(!err1.Equal(err3))
	is.True(!err1.Equal(err4))
	is.True(!err1.Equal(err5))
	is.True(!err1.Equal(err6))
}

func TestEError(t *testing.T) {
	is := is.New(t)

	is.Equal(New("foo", nil).Error(), "foo")
	is.Equal(New("foo", M{"a": "b"}).Error(), "foo (map[a:b])")
	is.Equal(Wrap(errors.New("foo"), "bar", nil).Error(), "bar: foo")
	is.Equal(Wrap(New("foo", nil), "bar", nil).Error(), "bar: foo")
	is.Equal(Wrap(New("foo", M{"a": "b"}), "bar", nil).Error(), "bar: foo (map[a:b])")
	is.Equal(Wrap(New("foo", nil), "bar", M{"c": "d"}).Error(), "bar (map[c:d]): foo")
	is.Equal(Wrap(New("foo", M{"a": "b"}), "bar", M{"c": "d"}).Error(), "bar (map[c:d]): foo (map[a:b])")
}

func TestEIsCode(t *testing.T) {
	is := is.New(t)

	errs := []error{
		New("foo", nil),
		New(Code("foo"), nil),
		New("foo", M{"a": "b"}),
		New(Code("foo"), M{"a": "b"}),
		Wrap(New("foo", nil), "bar", nil),
		Wrap(wrappedError{New("foo", nil)}, "bar", nil),
	}

	for _, err := range errs {
		is.True(errors.Is(err, Code("foo")))
	}
}

func TestEIsE(t *testing.T) {
	is := is.New(t)

	errFoo := New("foo", nil)

	errs := []error{
		errFoo,
		Wrap(errFoo, "bar", nil),
		Wrap(wrappedError{errFoo}, "bar", nil),
	}

	for _, err := range errs {
		is.True(errors.Is(err, errFoo))
	}
}

func TestEUnwrap(t *testing.T) {
	is := is.New(t)

	err1 := errors.New("underlying error")
	err2 := Wrap(err1, "foo", nil)
	err3 := Wrap(err2, "bar", nil)

	is.Equal(errors.Unwrap(err3), err2)
	is.Equal(errors.Unwrap(err2), err1)
}

func TestEAs(t *testing.T) {
	is := is.New(t)

	err := New("foo", nil)

	var errFoo E
	is.True(errors.As(err, &errFoo))

	_, ok := gerrors.As[E](err)
	is.True(ok)
}
