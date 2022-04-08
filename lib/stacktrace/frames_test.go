package stacktrace

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

// xxxxx - do not move the following code

var thisLineNumber = 12

func func1() []Frame {
	return GetCallerFrames(1)
}

func func2() []Frame {
	return func1()
}

func TestGetCallerFrames(t *testing.T) {
	is := is.New(t)

	frames := func2()

	// xxxxx - okay, you can change things now

	is.True(len(frames) > 3)

	is.True(strings.HasSuffix(frames[0].File, "/lib/stacktrace/frames_test.go"))
	is.Equal(frames[0].Line, thisLineNumber+3)
	is.Equal(frames[0].Function, "github.com/wearemojo/mojo-public-go/lib/stacktrace.func1")

	is.True(strings.HasSuffix(frames[1].File, "/lib/stacktrace/frames_test.go"))
	is.Equal(frames[1].Line, thisLineNumber+7)
	is.Equal(frames[1].Function, "github.com/wearemojo/mojo-public-go/lib/stacktrace.func2")

	is.True(strings.HasSuffix(frames[2].File, "/lib/stacktrace/frames_test.go"))
	is.Equal(frames[2].Line, thisLineNumber+13)
	is.Equal(frames[2].Function, "github.com/wearemojo/mojo-public-go/lib/stacktrace.TestGetCallerFrames")
}

var expectedFormat = `github.com/wearemojo/mojo-public-go/lib/foo.doFoo
	/lib/foo/foo.go:123
github.com/wearemojo/mojo-public-go/lib/foo.barThing
	/lib/foo/bar.go:456
`

func TestFormatFrames(t *testing.T) {
	is := is.New(t)

	frames := []Frame{
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

	is.Equal(FormatFrames(frames), expectedFormat)
}
