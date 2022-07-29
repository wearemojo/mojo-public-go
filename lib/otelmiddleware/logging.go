package otelmiddleware

import (
	"context"
	"fmt"
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
				url := fmt.Sprintf("%s?project=%s&tid=%s", gcpBaseURL, gcpProjectID, spanContext.TraceID())
				res := fmt.Sprintf("projects/%s/traces/%s", gcpProjectID, spanContext.TraceID())

				handleCLogError(ctx, clog.SetFields(ctx, clog.Fields{
					"trace_url": url,
					"trace_id":  spanContext.TraceID().String(),
					"span_id":   spanContext.SpanID().String(),

					"logging.googleapis.com/trace":  res,
					"logging.googleapis.com/spanId": spanContext.SpanID().String(),
				}))
			}

			next.ServeHTTP(res, req)
		})
	}
}

func handleCLogError(ctx context.Context, err error) {
	if err != nil {
		mlog.Warn(ctx, merr.Wrap(ctx, err, "clog_set_fields_failed", nil))
	}
}
