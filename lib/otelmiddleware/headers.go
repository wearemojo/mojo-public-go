package otelmiddleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

const (
	headerKeyID  = "Trace-Id"
	headerKeyURL = "Trace-Url"
	gcpBaseURL   = "https://console.cloud.google.com/traces/list"
)

func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		id := getTraceID(req)

		if id != "" {
			res.Header().Set(headerKeyID, id)
		}

		next.ServeHTTP(res, req)
	})
}

func TraceIDWithGCPURL(gcpProjectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			id := getTraceID(req)

			if id != "" {
				url := fmt.Sprintf("%s?project=%s&tid=%s", gcpBaseURL, gcpProjectID, id)

				res.Header().Set(headerKeyID, id)
				res.Header().Set(headerKeyURL, url)
			}

			next.ServeHTTP(res, req)
		})
	}
}

func getTraceID(req *http.Request) string {
	ctx := req.Context()
	id := trace.SpanContextFromContext(ctx).TraceID()

	if !id.IsValid() {
		return ""
	}

	return id.String()
}
