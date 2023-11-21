package jsonschema

import (
	"github.com/wearemojo/mojo-public-go/lib/ksuid"
	"github.com/xeipuuv/gojsonschema"
)

func RegisterKSUIDFormat() {
	gojsonschema.FormatCheckers.Add("ksuid", ksuidFormatChecker{})
}

type ksuidFormatChecker struct{}

func (f ksuidFormatChecker) IsFormat(input any) bool {
	str, ok := input.(string)
	if !ok {
		return false
	}

	_, err := ksuid.Parse(str)
	return err == nil
}
