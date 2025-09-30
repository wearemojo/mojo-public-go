package merr

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/wearemojo/mojo-public-go/lib/gerrors"
	"github.com/wearemojo/mojo-public-go/lib/gjson"
	"github.com/wearemojo/mojo-public-go/lib/stacktrace"
)

func TestNew(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

	err := New(ctx, "foo", nil)

	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, nil)
	is.True(strings.HasSuffix(err.Stack[0].File, "/lib/merr/merr_test.go"))
	is.Equal(err.Reasons, nil)
}

func TestNewMeta(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

	err := New(ctx, "foo", M{"a": "b"})

	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, M{"a": "b"})
	is.True(strings.HasSuffix(err.Stack[0].File, "/lib/merr/merr_test.go"))
	is.Equal(err.Reasons, nil)
}

func TestWrap(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

	err1 := errors.New("underlying error") //nolint:err113,forbidigo // needed for testing
	err2 := New(ctx, "foo", nil, err1)

	err, ok := gerrors.As[E](err2)

	is.True(ok)
	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, nil)
	is.True(strings.HasSuffix(err.Stack[0].File, "/lib/merr/merr_test.go"))
	is.Equal(len(err.Reasons), 1)
	is.Equal(err.Reasons[0], errors.New("underlying error")) //nolint:err113,forbidigo // needed for testing
}

func TestWrapMeta(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

	err1 := errors.New("underlying error") //nolint:err113,forbidigo // needed for testing
	err2 := New(ctx, "foo", M{"a": "b"}, err1)

	err, ok := gerrors.As[E](err2)

	is.True(ok)
	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, M{"a": "b"})
	is.True(strings.HasSuffix(err.Stack[0].File, "/lib/merr/merr_test.go"))
	is.Equal(len(err.Reasons), 1)
	is.Equal(err.Reasons[0], errors.New("underlying error")) //nolint:err113,forbidigo // needed for testing
}

func TestEqual(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

	//nolint:lll // same line to ensure same stack trace 😅
	err1, err2, err3, err4, err5, err6 := New(ctx, "foo", M{"a": "b"}), New(ctx, "foo", M{"a": "b"}), New(ctx, "foo", M{"a": "c"}), New(ctx, "bar", M{"a": "b"}), New(ctx, "foo", nil), New(ctx, "foo", M{"a": "b"})
	err6.Reasons = []error{errors.New("foo")} //nolint:err113,forbidigo // needed for testing

	is.True(err1.Equal(err2))
	is.True(!err1.Equal(err3))
	is.True(!err1.Equal(err4))
	is.True(!err1.Equal(err5))
	is.True(!err1.Equal(err6))
}

func TestString(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

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
	ctx := t.Context()

	is.Equal(New(ctx, "foo", nil).Error(), "foo")
	is.Equal(New(ctx, "foo", M{"a": "b"}).Error(), "foo (map[a:b])")
	is.Equal(New(ctx, "bar", nil, errors.New("foo")).Error(), "bar\n- foo") //nolint:err113,forbidigo // needed for testing
	is.Equal(New(ctx, "bar", nil, New(ctx, "foo", nil)).Error(), "bar\n- foo")
	is.Equal(New(ctx, "bar", nil, New(ctx, "foo", M{"a": "b"})).Error(), "bar\n- foo (map[a:b])")
	is.Equal(New(ctx, "bar", M{"c": "d"}, New(ctx, "foo", nil)).Error(), "bar (map[c:d])\n- foo")
	is.Equal(New(ctx, "bar", M{"c": "d"}, New(ctx, "foo", M{"a": "b"})).Error(), "bar (map[c:d])\n- foo (map[a:b])")
}

func TestEIsCode(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

	errs := []error{
		New(ctx, "foo", nil),
		New(ctx, Code("foo"), nil),
		New(ctx, "foo", M{"a": "b"}),
		New(ctx, Code("foo"), M{"a": "b"}),
		New(ctx, "bar", nil, New(ctx, "foo", nil)),
		New(ctx, "bar", nil, wrappedError{New(ctx, "foo", nil)}),
	}

	for _, err := range errs {
		is.True(errors.Is(err, Code("foo")))
	}
}

func TestEIsE(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

	errFoo := New(ctx, "foo", nil)

	errs := []error{
		errFoo,
		New(ctx, "bar", nil, errFoo),
		New(ctx, "bar", nil, wrappedError{errFoo}),
	}

	for _, err := range errs {
		is.True(errors.Is(err, errFoo))
	}
}

func TestEUnwrap(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

	err1 := errors.New("underlying error") //nolint:err113,forbidigo // needed for testing
	err2 := New(ctx, "foo", nil, err1)
	err3 := New(ctx, "bar", nil, err2)

	// Unwrap returns nil for multi-error wrapping
	is.Equal(errors.Unwrap(err3), nil)
	is.Equal(errors.Unwrap(err2), nil)

	is.Equal(err3.Unwrap(), []error{err2})
	is.Equal(err2.Unwrap(), []error{err1})
}

func TestEAs(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

	err := New(ctx, "foo", nil)

	var errFoo E
	is.True(errors.As(err, &errFoo))

	_, ok := gerrors.As[E](err)
	is.True(ok)
}

type unmarshallableError struct {
	GoUnmarshalYourself func(ctx context.Context) error
}

func (ue unmarshallableError) Error() string {
	return "literally unmarshallable"
}

func TestReasonsMarshallingEnforcing(t *testing.T) {
	is := is.New(t)

	ctx := t.Context()

	reason := unmarshallableError{
		GoUnmarshalYourself: func(context.Context) error {
			//nolint:err113,forbidigo // testing be crazy
			return errors.New("good luck seeing me in the logs")
		},
	}

	errFoo := New(ctx, "foo", nil, reason)

	output, err := json.Marshal(errFoo)
	is.NoErr(err)

	mapped, err := gjson.Unmarshal[map[string]any](output)
	is.NoErr(err)

	reasons, ok := mapped["reasons"].([]any)
	is.True(ok)
	is.Equal(len(reasons), 1)
	is.Equal(reasons[0], "merr.unmarshallableError{GoUnmarshalYourself:func(context.Context) error {...}}")
}
