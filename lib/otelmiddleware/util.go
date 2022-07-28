package otelmiddleware

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

const gcpBaseURL = "https://console.cloud.google.com/traces/list"

func getTraceID(req *http.Request) string {
	ctx := req.Context()
	id := trace.SpanContextFromContext(ctx).TraceID()

	if !id.IsValid() {
		return ""
	}

	return id.String()
}
