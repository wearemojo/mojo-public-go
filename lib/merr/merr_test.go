package merr

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/wearemojo/mojo-public-go/lib/gerrors"
	"github.com/wearemojo/mojo-public-go/lib/stacktrace"
)

func TestNew(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err := New(ctx, "foo", nil)

	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, nil)
	is.True(strings.HasSuffix(err.Stack[0].File, "/lib/merr/merr_test.go"))
	is.Equal(err.Reason, nil)
}

func TestNewMeta(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err := New(ctx, "foo", M{"a": "b"})

	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, M{"a": "b"})
	is.True(strings.HasSuffix(err.Stack[0].File, "/lib/merr/merr_test.go"))
	is.Equal(err.Reason, nil)
}

func TestWrap(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err1 := errors.New("underlying error") //nolint:goerr113,forbidigo // needed for testing
	err2 := Wrap(ctx, err1, "foo", nil)

	err, ok := gerrors.As[E](err2)

	is.True(ok)
	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, nil)
	is.True(strings.HasSuffix(err.Stack[0].File, "/lib/merr/merr_test.go"))
	is.Equal(err.Reason, errors.New("underlying error")) //nolint:goerr113,forbidigo // needed for testing
}

func TestWrapMeta(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err1 := errors.New("underlying error") //nolint:goerr113,forbidigo // needed for testing
	err2 := Wrap(ctx, err1, "foo", M{"a": "b"})

	err, ok := gerrors.As[E](err2)

	is.True(ok)
	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, M{"a": "b"})
	is.True(strings.HasSuffix(err.Stack[0].File, "/lib/merr/merr_test.go"))
	is.Equal(err.Reason, errors.New("underlying error")) //nolint:goerr113,forbidigo // needed for testing
}

func TestWrapNil(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err := Wrap(ctx, nil, "foo", nil)

	is.Equal(err, nil)
}

func TestEqual(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	//nolint:lll // same line to ensure same stack trace 😅
	err1, err2, err3, err4, err5, err6 := New(ctx, "foo", M{"a": "b"}), New(ctx, "foo", M{"a": "b"}), New(ctx, "foo", M{"a": "c"}), New(ctx, "bar", M{"a": "b"}), New(ctx, "foo", nil), New(ctx, "foo", M{"a": "b"})
	err6.Reason = errors.New("foo") //nolint:goerr113,forbidigo // needed for testing

	is.True(err1.Equal(err2))
	is.True(!err1.Equal(err3))
	is.True(!err1.Equal(err4))
	is.True(!err1.Equal(err5))
	is.True(!err1.Equal(err6))
}

func TestString(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err := New(ctx, "foo", M{"a": "b"})

	err.Stack = []stacktrace.Frame{
		{
			File:     "/lib/foo/foo.go",
			Line:     123,
			Function: "github.com/wearemojo/mojo-public-go/lib/foo.doFoo",
		},
		{
			File:     "/lib/foo/bar.go",
			Line:     456,
			Function: "github.com/wearemojo/mojo-public-go/lib/foo.barThing",
		},
	}

	expected := `foo (map[a:b])

github.com/wearemojo/mojo-public-go/lib/foo.doFoo
	/lib/foo/foo.go:123
github.com/wearemojo/mojo-public-go/lib/foo.barThing
	/lib/foo/bar.go:456
`

	is.Equal(err.String(), expected)
}

func TestEError(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	is.Equal(New(ctx, "foo", nil).Error(), "foo")
	is.Equal(New(ctx, "foo", M{"a": "b"}).Error(), "foo (map[a:b])")
	is.Equal(Wrap(ctx, errors.New("foo"), "bar", nil).Error(), "bar: foo") //nolint:goerr113,forbidigo // needed for testing
	is.Equal(Wrap(ctx, New(ctx, "foo", nil), "bar", nil).Error(), "bar: foo")
	is.Equal(Wrap(ctx, New(ctx, "foo", M{"a": "b"}), "bar", nil).Error(), "bar: foo (map[a:b])")
	is.Equal(Wrap(ctx, New(ctx, "foo", nil), "bar", M{"c": "d"}).Error(), "bar (map[c:d]): foo")
	is.Equal(Wrap(ctx, New(ctx, "foo", M{"a": "b"}), "bar", M{"c": "d"}).Error(), "bar (map[c:d]): foo (map[a:b])")
}

func TestEIsCode(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	errs := []error{
		New(ctx, "foo", nil),
		New(ctx, Code("foo"), nil),
		New(ctx, "foo", M{"a": "b"}),
		New(ctx, Code("foo"), M{"a": "b"}),
		Wrap(ctx, New(ctx, "foo", nil), "bar", nil),
		Wrap(ctx, wrappedError{New(ctx, "foo", nil)}, "bar", nil),
	}

	for _, err := range errs {
		is.True(errors.Is(err, Code("foo")))
	}
}

func TestEIsE(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	errFoo := New(ctx, "foo", nil)

	errs := []error{
		errFoo,
		Wrap(ctx, errFoo, "bar", nil),
		Wrap(ctx, wrappedError{errFoo}, "bar", nil),
	}

	for _, err := range errs {
		is.True(errors.Is(err, errFoo))
	}
}

func TestEUnwrap(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err1 := errors.New("underlying error") //nolint:goerr113,forbidigo // needed for testing
	err2 := Wrap(ctx, err1, "foo", nil)
	err3 := Wrap(ctx, err2, "bar", nil)

	is.Equal(errors.Unwrap(err3), err2)
	is.Equal(errors.Unwrap(err2), err1)
}

func TestEAs(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err := New(ctx, "foo", nil)

	var errFoo E
	is.True(errors.As(err, &errFoo))

	_, ok := gerrors.As[E](err)
	is.True(ok)
}
