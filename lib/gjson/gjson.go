package gjson

import (
	"bytes"
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

// DecoderFor returns a decoder over an already-read JSON value, carrying the
// given options.
//
// Polymorphic unmarshalling has to read the discriminator before the concrete
// type is known, which consumes the value from the parent decoder. Wrapping the
// bytes back up in a decoder lets the concrete decode receive the parent's
// options, which it would not if handed the bytes directly.
//
// Pass dec.Options() as the first option, then any override the call site needs.
// Later options win.
func DecoderFor(data jsontext.Value, opts ...json.Options) *jsontext.Decoder {
	return jsontext.NewDecoder(bytes.NewReader(data), opts...)
}

func MustUnmarshal[T any](data []byte) T {
	res, err := Unmarshal[T](data)
	if err != nil {
		panic(err)
	}
	return res
}

func Remarshal[T any](in T) (res T, err error) {
	data, err := json.Marshal(in)
	if err != nil {
		return res, err
	}
	return Unmarshal[T](data)
}
