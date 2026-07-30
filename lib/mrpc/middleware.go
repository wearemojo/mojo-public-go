package mrpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/authenforce"
	"github.com/wearemojo/mojo-public-go/lib/authparsing"
	"github.com/wearemojo/mojo-public-go/lib/bodycontext"
	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/clog"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
	"github.com/xeipuuv/gojsonschema"
)

// Logger inherits the context logger and reports RPC request success/failure.
func Logger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) error {
			ctx := req.Context()

			if info := getInfo(ctx); info != nil {
				clog.SetFields(ctx, clog.Fields{
					"rpc_version": info.Version,
					"rpc_method":  info.Method,
				})
			}

			tStart := time.Now()
			err := next(res, req)
			tEnd := time.Now()

			clog.SetFields(ctx, clog.Fields{
				"rpc_duration":    tEnd.Sub(tStart).String(),
				"rpc_duration_us": int64(tEnd.Sub(tStart) / time.Microsecond),
			})

			if err == nil {
				return nil
			}

			// rewrite common errors to internal error standard
			switch {
			case errors.Is(err, io.EOF):
				err = cher.New(cher.EOF, nil)
			case errors.Is(err, io.ErrUnexpectedEOF):
				err = cher.New(cher.UnexpectedEOF, nil)
			case errors.Is(err, context.Canceled):
				err = cher.New(cher.ContextCanceled, nil)
			}

			clog.SetError(ctx, err)

			return err
		}
	}
}

// validate buffers the JSON body and applies a JSON Schema validation.
func validate(schema *gojsonschema.Schema) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) error {
			ctx := req.Context()

			body, err := io.ReadAll(req.Body)
			if err != nil {
				if netErr, ok := errors.AsType[net.Error](err); ok {
					mlog.Warn(ctx, merr.New(ctx, "request_body_read_network_error", nil, netErr))
					return io.ErrUnexpectedEOF
				}

				return merr.New(ctx, "request_body_read_failed", nil, err)
			}

			result, err := schema.Validate(gojsonschema.NewBytesLoader(body))
			if err != nil {
				if errors.Is(err, io.EOF) {
					return cher.New("invalid_json", nil)
				}

				return merr.New(ctx, "request_body_validation_failed", nil, err)
			}

			if err := CoerceJSONSchemaError(result); err != nil {
				return err
			}

			req.Body = io.NopCloser(bytes.NewReader(body))

			return next(res, req)
		}
	}
}

// enforcerMiddleware runs a set of authorization enforcers against the parsed
// auth state and request body from the context.
func enforcerMiddleware(enforcers authenforce.Enforcers) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) error {
			ctx := req.Context()

			authState := authparsing.GetAuthState(ctx)
			body := bodycontext.GetContext(ctx)

			if err := enforcers.Run(ctx, authState, body); err != nil {
				return err
			}

			return next(res, req)
		}
	}
}
