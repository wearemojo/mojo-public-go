package otelmiddleware

import (
	"context"
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
	"go.opentelemetry.io/otel/trace"
)

func SetCLogFields(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		spanContext := trace.SpanContextFromContext(ctx)

		if spanContext.IsValid() {
			handleCLogError(ctx, clog.SetFields(ctx, clog.Fields{
				"trace_id": spanContext.TraceID().String(),
				"span_id":  spanContext.SpanID().String(),
			}))
		}

		next.ServeHTTP(res, req)
	})
}

func SetCLogFieldsForGCP(gcpProjectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			spanContext := trace.SpanContextFromContext(ctx)

			if spanContext.IsValid() {
				traceID := spanContext.TraceID().String()
				spanID := spanContext.SpanID().String()

				handleCLogError(ctx, clog.SetFields(ctx, clog.Fields{
					"trace_url": getGCPTraceURL(gcpProjectID, traceID),
					"trace_id":  traceID,
					"span_id":   spanID,

					"logging.googleapis.com/trace":  getGCPTracePath(gcpProjectID, traceID),
					"logging.googleapis.com/spanId": spanID,
				}))
			}

			next.ServeHTTP(res, req)
		})
	}
}

func handleCLogError(ctx context.Context, err error) {
	if err != nil {
		mlog.Warn(ctx, merr.New(ctx, "clog_set_fields_failed", nil, err))
	}
}
