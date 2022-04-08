package merr

import (
	"errors"
)

var _ error = Code("")

type Code string

func (c Code) Error() string {
	return string(c)
}

// IsCode is a convenience wrapper for `errors.Is`
//
// `errors.Is(err, "foo")` does not work
//
// `IsCode(err, "foo")` does - as the string is automatically converted
func IsCode(err error, code Code) bool {
	return errors.Is(err, code)
}
