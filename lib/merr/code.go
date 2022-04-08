package merr

import (
	"errors"
)

var _ error = Code("")

type Code string

func (c Code) Error() string {
	return string(c)
}

// IsCode is a convenience wrapper for `errors.Is` - as it works when called with strings directly
func IsCode(err error, code Code) bool {
	return errors.Is(err, code)
}
