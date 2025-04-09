package gbson

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Unmarshal parses BSON-encoded data and returns a value of the type T.
//
// Just a generic, type-safe wrapper around bson.Unmarshal.
func Unmarshal[T any](data []byte) (res T, err error) {
	return res, bson.Unmarshal(data, &res)
}

func MustUnmarshal[T any](data []byte) T {
	res, err := Unmarshal[T](data)
	if err != nil {
		panic(err)
	}
	return res
}
