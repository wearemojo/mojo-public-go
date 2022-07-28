package mlog

import (
	"context"
	"fmt"

	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/sirupsen/logrus"
	"github.com/wearemojo/mojo-public-go/lib/gcp"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog/indirect"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	indirect.Debug = Debug
	indirect.Info = Info
	indirect.Warn = Warn
	indirect.Error = Error
}

// Debug is to help with tracing the system's behavior
//
// e.g. logic has been evaluated to determine some action
func Debug(ctx context.Context, err merr.Merrer) {
	log(ctx, logrus.DebugLevel, err)
}

// Info is for informational messages which are not errors
//
// e.g. a system is starting up
func Info(ctx context.Context, err merr.Merrer) {
	log(ctx, logrus.InfoLevel, err)
}

// Warn covers issues which were handled and do not require specific action
// from an engineer, but which should be fixed at some point
//
// Only for issues which could occur indefinitely without serious consequences
//
// e.g. a system was unavailable, but failed gracefully
func Warn(ctx context.Context, err merr.Merrer) {
	log(ctx, logrus.WarnLevel, err)
}

// Error represents an issue which requires prompt and individual action from
// an engineer to resolve
//
// e.g. a data integrity issue has been identified which needs to be fixed
func Error(ctx context.Context, err merr.Merrer) {
	log(ctx, logrus.ErrorLevel, err)
}

const gcpBaseURL = "https://console.cloud.google.com/traces/list"

func log(ctx context.Context, level logrus.Level, err merr.Merrer) {
	if err == nil {
		return
	}

	merr := err.Merr()

	// logrus runs `.String()` on anything implementing `error`
	// so to get proper JSON, we need to copy the merrFields instead
	merrFields := merr.Fields()

	if level == logrus.InfoLevel {
		merrFields["stack"] = nil
	}

	fields := logrus.Fields{
		"merr": merrFields,
	}

	sc := trace.SpanContextFromContext(ctx)

	if sc.IsValid() {
		fields["trace_id"] = sc.TraceID().String()
		fields["span_id"] = sc.SpanID().String()
		fields["logging.googleapis.com/spanId"] = sc.SpanID().String()

		gcpProjectID, err := gcp.GetProjectID(ctx)

		if err == nil {
			url := fmt.Sprintf("%s?project=%s&tid=%s", gcpBaseURL, gcpProjectID, sc.TraceID())
			res := fmt.Sprintf("projects/%s/traces/%s", gcpProjectID, sc.TraceID())

			fields["trace_url"] = url
			fields["logging.googleapis.com/trace"] = res
		}
	}

	clog.Get(ctx).
		WithFields(fields).
		Log(level, merr.Code)
}
