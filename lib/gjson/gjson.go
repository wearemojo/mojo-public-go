package gjson

import (
	"encoding/json"
)

// Unmarshal parses JSON-encoded data and returns a value of the type T.
//
// Just a generic, type-safe wrapper around json.Unmarshal.
func Unmarshal[T any](data []byte) (res T, err error) {
	return res, json.Unmarshal(data, &res)
}

func MustUnmarshal[T any](data []byte) T {
	res, err := Unmarshal[T](data)
	if err != nil {
		panic(err)
	}
	return res
}
