package request

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/clog"
	"github.com/wearemojo/mojo-public-go/lib/gerrors"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
	"github.com/wearemojo/mojo-public-go/lib/slicefn"
)

type responseWriter struct {
	http.ResponseWriter

	Status int
	Bytes  int64
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.Status == 0 {
		rw.Status = code
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(bytes []byte) (int, error) {
	if rw.Status == 0 {
		rw.Status = http.StatusOK
		rw.WriteHeader(http.StatusOK)
	}

	rw.Bytes += int64(len(bytes))

	return rw.ResponseWriter.Write(bytes)
}

// Logger returns a middleware handler that wraps subsequent middleware/handlers and logs
// request information AFTER the request has completed. It also injects a request-scoped
// logger on the context which can be set, read and updated using clog lib
//
// Included fields:
//   - Request ID                (request_id)
//   - HTTP Method               (http_method)
//   - HTTP Path                 (http_path)
//   - HTTP Protocol Version     (http_proto)
//   - Remote Address            (http_remote_addr)
//   - User Agent Header         (http_user_agent)
//   - Referer Header            (http_referer)
//   - Duration with unit        (http_duration)
//   - Duration in microseconds  (http_duration_us)
//   - HTTP Status Code          (http_status)
//   - Response in bytes         (http_response_bytes)
//   - Client Version header     (http_client_version)
//   - User Agent header         (http_user_agent)
func Logger(log *logrus.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// create a mutable logger instance which will persist for the request
			// inject pointer to the logger into the request context
			ctx = clog.Set(ctx, log)
			r = r.WithContext(ctx)

			// panics inside handlers will be logged to standard before propagation
			defer clog.HandlePanic(ctx, true)

			clog.SetFields(ctx, clog.Fields{
				"http_origin":            r.Header.Get("Origin"),
				"http_referer":           r.Header.Get("Referer"),
				"mojo_rn_client_version": r.Header.Get("Mojo-Rn-Client-Version"),
			})

			// wrap given response writer with one that tracks status code/bytes written
			res := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(res, r)

			err := getError(clog.Get(ctx))
			if err == nil {
				mlog.Info(ctx, merr.New(ctx, "request_completed", nil))
				return
			}
			var fn func(context.Context, merr.Merrer)
			switch clog.DetermineLevel(err, clog.TimeoutsAsErrors(ctx)) {
			case
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel:
				fn = mlog.Error
			case logrus.WarnLevel:
				fn = mlog.Warn
			case
				logrus.InfoLevel,
				logrus.DebugLevel,
				logrus.TraceLevel:
				fn = mlog.Info
			}

			if mErr, ok := gerrors.As[merr.E](err); ok {
				fn(ctx, mErr)
			} else if cErr, ok := gerrors.As[cher.E](err); ok {
				reasons := slicefn.Map(cErr.Reasons, func(r cher.E) error { return r })
				// If the cher error has no reasons, add the cher error itself
				if len(reasons) == 0 {
					reasons = append(reasons, cErr)
				}
				fn(ctx, merr.New(ctx, merr.Code(cErr.Code), merr.M(cErr.Meta), reasons...))
			} else {
				fn(ctx, merr.New(ctx, "unexpected_request_failure", nil, err))
			}
		})
	}
}

// getError returns the error if one is set on the log entry
func getError(l *logrus.Entry) error {
	if erri, ok := l.Data[logrus.ErrorKey]; ok {
		if err, ok := erri.(error); ok {
			return err
		}
	}

	return nil
}
