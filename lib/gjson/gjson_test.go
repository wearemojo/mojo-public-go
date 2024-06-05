package gjson

import (
	"testing"

	"github.com/matryer/is"
)

type testInterface interface {
	code() string
}

type testStruct struct {
	Foo string `json:"foo"`
}

func (testStruct) code() string { return "bar" }

func TestUnmarshalToConcrete(t *testing.T) {
	is := is.New(t)

	res, err := Unmarshal[testStruct]([]byte(`{"foo":"1"}`))
	is.NoErr(err)
	is.Equal(res, testStruct{"1"})
}

func TestUnmarshalToInterface(t *testing.T) {
	is := is.New(t)

	var res any
	res, err := Unmarshal[testStruct]([]byte(`{"foo":"2"}`))
	is.NoErr(err)
	is.Equal(res, testStruct{"2"})
}

func TestUnmarshalToConstrainedInterface(t *testing.T) {
	is := is.New(t)

	var res testInterface
	res, err := Unmarshal[testStruct]([]byte(`{"foo":"3"}`))
	is.NoErr(err)
	is.Equal(res, testStruct{"3"})
}
