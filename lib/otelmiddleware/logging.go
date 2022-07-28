package otelmiddleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
)

func SetCLogFields(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id := getTraceID(req)

		if id != "" {
			handleCLogError(ctx, clog.SetFields(ctx, clog.Fields{
				"trace_id": id,
			}))
		}

		next.ServeHTTP(res, req)
	})
}

func SetCLogFieldsForGCP(gcpProjectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			id := getTraceID(req)

			if id != "" {
				url := fmt.Sprintf("%s?project=%s&tid=%s", gcpBaseURL, gcpProjectID, id)
				res := fmt.Sprintf("projects/%s/traces/%s", gcpProjectID, id)

				handleCLogError(ctx, clog.SetFields(ctx, clog.Fields{
					"trace_id":  id,
					"trace_url": url,

					"logging.googleapis.com/trace": res,
				}))
			}

			next.ServeHTTP(res, req)
		})
	}
}

func handleCLogError(ctx context.Context, err error) {
	if err != nil {
		mlog.Warn(ctx, merr.Wrap(err, "clog_set_fields_failed", nil))
	}
}
