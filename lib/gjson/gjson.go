package gjson

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// Unmarshal parses JSON-encoded data and returns a value of the type T.
//
// Just a generic, type-safe wrapper around json.Unmarshal.
func Unmarshal[T any](data []byte) (res T, err error) {
	return res, json.Unmarshal(data, &res)
}

// UnmarshalDecode reads the next JSON value from dec and returns a value of
// the type T.
//
// Just a generic, type-safe wrapper around json.UnmarshalDecode. Prefer this
// over Unmarshal inside an UnmarshalJSONFrom method, so the caller's options
// reach the decoded value.
func UnmarshalDecode[T any](dec *jsontext.Decoder) (res T, err error) {
	return res, json.UnmarshalDecode(dec, &res)
}

func MustUnmarshal[T any](data []byte) T {
	res, err := Unmarshal[T](data)
	if err != nil {
		panic(err)
	}
	return res
}

func Remarshal[T any](in T) (res T, err error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return res, err
	}
	return Unmarshal[T](bytes)
}
