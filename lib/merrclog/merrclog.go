package merrclog

import (
	"context"

	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/sirupsen/logrus"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

func Info(ctx context.Context, err merr.EInterface) {
	log(ctx, logrus.InfoLevel, err)
}

func Warn(ctx context.Context, err merr.EInterface) {
	log(ctx, logrus.WarnLevel, err)
}

func Error(ctx context.Context, err merr.EInterface) {
	log(ctx, logrus.ErrorLevel, err)
}

func log(ctx context.Context, level logrus.Level, err merr.EInterface) {
	if err == nil {
		return
	}

	merr := err.GetConcrete()

	// logrus runs `.String()` on anything implementing `error`
	// so to get proper JSON, we need to copy the fields instead
	fields := merr.Fields()

	clog.Get(ctx).
		WithField(logrus.ErrorKey, fields).
		Log(level, merr.Code)
}
