package cher

import (
	"errors"
)

type timeoutError interface {
	Error() string
	Timeout() bool
}

func CoerceThirdPartyTimeout(err error) error {
	if netErr, ok := errors.AsType[timeoutError](err); ok && netErr.Timeout() {
		return New(ThirdPartyTimeout, M{
			"error": netErr.Error(),
		})
	}

	return err
}
