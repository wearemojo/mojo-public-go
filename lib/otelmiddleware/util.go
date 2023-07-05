package otelmiddleware

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/otel/trace"
)

const gcpBaseURL = "https://console.cloud.google.com"

func getTraceID(req *http.Request) string {
	ctx := req.Context()
	id := trace.SpanContextFromContext(ctx).TraceID()

	if !id.IsValid() {
		return ""
	}

	return id.String()
}

func getGCPTracePath(gcpProjectID, traceID string) string {
	return fmt.Sprintf("projects/%s/traces/%s", gcpProjectID, traceID)
}

func getGCPTraceURL(gcpProjectID, traceID string) string {
	params := url.Values{
		"project": []string{gcpProjectID},
		"tid":     []string{traceID},
	}

	return fmt.Sprintf("%s/traces/list?%s", gcpBaseURL, params.Encode())
}

func getGCPTraceLogsURL(gcpProjectID, traceID string, refTime time.Time) string {
	tracePath := getGCPTracePath(gcpProjectID, traceID)
	query := fmt.Sprintf("(trace=\"%s\" OR labels.\"appengine.googleapis.com/trace_id\"=\"%s\")", tracePath, traceID)
	timeRange := fmt.Sprintf("%s/%s--PT5M", refTime.Format(time.RFC3339), refTime.Format(time.RFC3339))

	params := url.Values{
		"project":   []string{gcpProjectID},
		"query":     []string{query},
		"timeRange": []string{timeRange},
	}

	return fmt.Sprintf("%s/logs/query?%s", gcpBaseURL, params.Encode())
}
