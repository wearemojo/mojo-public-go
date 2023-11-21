package otelmiddleware

import (
	"net/http"

	"github.com/wearemojo/mojo-public-go/lib/clog"
	"go.opentelemetry.io/otel/trace"
)

func SetCLogFields(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		spanContext := trace.SpanContextFromContext(ctx)

		if spanContext.IsValid() {
			clog.SetFields(ctx, clog.Fields{
				"trace_id": spanContext.TraceID().String(),
				"span_id":  spanContext.SpanID().String(),
			})
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
				clog.SetFields(ctx, clog.Fields{
					"logging.googleapis.com/trace":  getGCPTracePath(gcpProjectID, spanContext.TraceID().String()),
					"logging.googleapis.com/spanId": spanContext.SpanID().String(),
				})
			}

			next.ServeHTTP(res, req)
		})
	}
}
