package mlog

import (
	"context"

	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/sirupsen/logrus"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog/indirect"
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

func log(ctx context.Context, level logrus.Level, err merr.Merrer) {
	if err == nil {
		return
	}

	merr := err.Merr()

	// logrus runs `.String()` on anything implementing `error`
	// so to get proper JSON, we need to copy the fields instead
	fields := merr.Fields()

	if level == logrus.InfoLevel {
		fields["stack"] = nil
	}

	clog.Get(ctx).
		WithField("merr", fields).
		Log(level, merr.Code)
}
