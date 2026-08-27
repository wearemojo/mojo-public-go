package gjson

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"time"
)

// Duration is a time.Duration that carries through JSON as an integer count of
// nanoseconds.
//
// encoding/json/v2 refuses a bare time.Duration outright - it has no default
// representation, and no `format` tag option covers it - so a named type with
// its own methods is the only way to carry one through JSON. The failure is at
// the level of the whole enclosing value, so a single unusable field breaks
// every other field alongside it.
//
// Nanoseconds are what encoding/json v1 produced for a bare time.Duration,
// since the underlying type is int64, so the wire format is unchanged.
type Duration time.Duration

// Duration returns the value as a time.Duration.
//
// A defined type does not inherit the methods of its underlying type, so this
// is needed wherever the standard library expects a time.Duration.
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d Duration) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, int64(d))
}

func (d *Duration) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var nanos int64
	if err := json.UnmarshalDecode(dec, &nanos); err != nil {
		return err
	}

	*d = Duration(nanos)

	return nil
}
