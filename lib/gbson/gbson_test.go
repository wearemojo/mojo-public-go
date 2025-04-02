package gbson

import (
	"testing"

	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson"
)

type testInterface interface {
	code() string
}

type testStruct struct {
	Foo string `bson:"foo"`
}

func (testStruct) code() string { return "bar" }

func TestUnmarshalToConcrete(t *testing.T) {
	is := is.New(t)

	bsonBytes, err := bson.Marshal(testStruct{"1"})
	is.NoErr(err)
	res, err := Unmarshal[testStruct](bsonBytes)
	is.NoErr(err)
	is.Equal(res, testStruct{"1"})
}

func TestUnmarshalToInterface(t *testing.T) {
	is := is.New(t)

	var res any
	bsonBytes, err := bson.Marshal(testStruct{"2"})
	is.NoErr(err)
	res, err = Unmarshal[testStruct](bsonBytes)
	is.NoErr(err)
	is.Equal(res, testStruct{"2"})
}

func TestUnmarshalToConstrainedInterface(t *testing.T) {
	is := is.New(t)

	var res testInterface
	bsonBytes, err := bson.Marshal(testStruct{"3"})
	is.NoErr(err)
	res, err = Unmarshal[testStruct](bsonBytes)
	is.NoErr(err)
	is.Equal(res, testStruct{"3"})
}
