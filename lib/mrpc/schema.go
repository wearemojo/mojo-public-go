package mrpc

import (
	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/slicefn"
	"github.com/xeipuuv/gojsonschema"
)

// CoerceJSONSchemaError converts a JSON Schema validation result into a
// bad_request cher error, with one schema_failure reason per validation
// error.
func CoerceJSONSchemaError(result *gojsonschema.Result) error {
	if result.Valid() {
		return nil
	}

	return cher.New(cher.BadRequest, nil, slicefn.Map(result.Errors(), func(err gojsonschema.ResultError) cher.E {
		return cher.E{
			Code: "schema_failure",
			Meta: cher.M{
				"field":   err.Field(),
				"type":    err.Type(),
				"message": err.Description(),
			},
		}
	})...)
}
