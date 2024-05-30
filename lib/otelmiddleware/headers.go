package otelmiddleware

import (
	"net/http"
	"time"
)

const (
	headerKeyID      = "Trace-Id"
	headerKeyURL     = "Trace-Url"
	headerKeyLogsURL = "Trace-Logs-Url"
)

func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		traceID := getTraceID(req)

		if traceID != "" {
			res.Header().Set(headerKeyID, traceID)
		}

		next.ServeHTTP(res, req)
	})
}

func TraceIDWithGCPURLs(gcpProjectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			traceID := getTraceID(req)

			if traceID != "" {
				res.Header().Set(headerKeyID, traceID)
				res.Header().Set(headerKeyURL, getGCPTraceURL(gcpProjectID, traceID))
				res.Header().Set(headerKeyLogsURL, GetGCPTraceLogsURL(gcpProjectID, traceID, time.Now()))
			}

			next.ServeHTTP(res, req)
		})
	}
}
