package gcppubsub

import (
	"context"
	"embed"
	"encoding/base64"
	"encoding/json"

	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/cuvva/cuvva-public-go/lib/jsonschema"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/xeipuuv/gojsonschema"
	"google.golang.org/api/pubsub/v1"
)

//go:embed *.json
var schemaFS embed.FS

var (
	schema             = jsonschema.NewFS(schemaFS).LoadJSONExt
	ReceiveEventSchema = schema("receive_event")
)

type IncomingMessage struct {
	Subscription string                `json:"subscription"`
	Message      *pubsub.PubsubMessage `json:"message"`
}

func (im IncomingMessage) GetPayload(ctx context.Context, schema *gojsonschema.Schema, req any) error {
	bytes, err := base64.StdEncoding.DecodeString(im.Message.Data)
	if err != nil {
		return err
	}

	ld := gojsonschema.NewBytesLoader(bytes)

	result, err := schema.Validate(ld)
	if err != nil {
		return merr.Wrap(err, "cannot_validate_message", nil)
	}

	if err = crpc.CoerceJSONSchemaError(result); err != nil {
		return err
	}

	return json.Unmarshal(bytes, req)
}

func MakeSchema(schema gojsonschema.JSONLoader, pointer *gojsonschema.Schema) (*gojsonschema.Schema, error) {
	if pointer != nil {
		return pointer, nil
	}

	return gojsonschema.NewSchemaLoader().Compile(schema)
}
