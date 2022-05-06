package gcppubsub

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/wearemojo/mojo-public-go/lib/authenforce"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed receive_event.json
var rootSchemaBytes []byte

var (
	errorType   = reflect.TypeOf((*error)(nil)).Elem()
	contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
	rootSchema  = gojsonschema.NewBytesLoader(rootSchemaBytes)
)

type IncomingMessage struct {
	Subscription string  `json:"subscription"`
	Message      Message `json:"message"`
}

type Message struct {
	Attributes  map[string]any `json:"attributes"`
	Data        string         `json:"data"`
	MessageID   string         `json:"message_id"`
	PublishTime time.Time      `json:"publish_time"`
}

type Server struct {
	*crpc.Server
}

func New() *Server {
	server := crpc.NewServer(authenforce.CRPCMiddleware(authenforce.Enforcers{
		authenforce.UnsafeNoAuthentication,
	}))

	return &Server{
		server,
	}
}

func (s *Server) RegisterEventHandler(method, version string, schema gojsonschema.JSONLoader, fn any) {
	wrapped := MustWrap(schema, fn)

	s.Register(method, version, rootSchema, wrapped)
}

func MustWrap(schema gojsonschema.JSONLoader, fn any) any {
	v := reflect.ValueOf(fn)
	t := v.Type()

	switch {
	case t.Kind() != reflect.Func:
		panic(merr.New("fn_not_function", merr.M{"type": t.Kind()}))
	case t.NumIn() < 1 || t.NumIn() > 2:
		panic(merr.New("fn_bad_arg_count", merr.M{"args": t.NumIn()}))
	case t.NumOut() != 1:
		panic(merr.New("fn_outputs_invalid", merr.M{"args": t.NumOut()}))
	case !t.In(0).Implements(contextType):
		panic(merr.New("fn_inputs_invalid", merr.M{"first_arg": t.In(0)}))
	case !t.Out(0).Implements(errorType):
		panic(merr.New("fn_outputs_invalid", merr.M{"out_type": t.Out(0)}))
	}

	var reqT reflect.Type

	if t.NumIn() == 2 {
		reqT = t.In(1).Elem()

		switch {
		case t.In(1).Kind() != reflect.Ptr:
			panic(merr.New("fn_inputs_invalid", merr.M{"second_arg": t.In(1)}))
		case reqT.Kind() != reflect.Struct:
			panic(merr.New("fn_inputs_invalid", merr.M{"second_arg": reqT.Kind()}))
		}
	}

	var err error
	var compiledSchema *gojsonschema.Schema

	if reqT != nil {
		compiledSchema, err = gojsonschema.NewSchemaLoader().Compile(schema)
		if err != nil {
			panic(fmt.Sprintf("json schema error: %s", err))
		}
	}

	handler := func(ctx context.Context, im *IncomingMessage) (err error) {
		inputs := []reflect.Value{reflect.ValueOf(ctx)}

		if reqT != nil && compiledSchema != nil {
			bytes, err := base64.StdEncoding.DecodeString(im.Message.Data)
			if err != nil {
				return err
			}

			ld := gojsonschema.NewBytesLoader(bytes)

			result, err := compiledSchema.Validate(ld)
			if err != nil {
				return merr.Wrap(err, "cannot_validate_message", nil)
			}

			if err = crpc.CoerceJSONSchemaError(result); err != nil {
				return err
			}

			req := reflect.New(reqT)
			if err := json.Unmarshal(bytes, req.Interface()); err != nil {
				return err
			}

			inputs = append(inputs, req)
		}

		res := v.Call(inputs)
		if err := res[0]; !err.IsNil() {
			e, _ := err.Interface().(error)
			return e
		}

		return nil
	}

	return handler
}
