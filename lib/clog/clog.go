package clog

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/servicecontext"
	"github.com/wearemojo/mojo-public-go/lib/version"
)

type contextKey string

type Fields map[string]any

// LoggerKey is the key used for request-scoped loggers in a requests context map
const loggerKey contextKey = "clog"

const (
	// ServiceKey is the log entry key for the name of the crpc service
	ServiceKey = "_service"

	// VersionKey is the log entry key for the current version of the codebase
	VersionKey = "_commit_hash"

	// LevelKey is the log entry key for the log level
	LevelKey = "severity"

	// TimeKey is the log entry key for the timestamp
	TimeKey = "time"
)

// Config allows services to configure the logging format, level and storage options
// for Logrus logging.
type Config struct {
	// Format configures the output format. Possible options:
	//   - text - logrus default text output, good for local development
	//   - json - fields and message encoded as json, good for storage in e.g. cloudwatch
	Format string `json:"format"`

	// Debug enables debug level logging, otherwise INFO level
	Debug bool `json:"debug"`
}

// Configure applies standard Logging structure options to a logrus Entry.
func (c Config) Configure(ctx context.Context) *logrus.Entry {
	var serviceName string
	if svc := servicecontext.GetContext(ctx); svc != nil {
		serviceName = svc.Service
	}

	log := logrus.WithFields(logrus.Fields{
		ServiceKey: serviceName,
		VersionKey: version.Revision,
	})

	switch c.Format {
	case "json", "logstash":
		log.Logger.Formatter = &fallbackFormatter{
			baseFormatter: &logrus.JSONFormatter{
				DisableHTMLEscape: true,
				TimestampFormat:   time.RFC3339Nano,
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime:  TimeKey,
					logrus.FieldKeyLevel: LevelKey,
					// msg isn't overridden because it's not usually useful
					// and makes other data fields harder to find
					// logrus.FieldKeyMsg:   "message",
				},
			},
		}

	default:
		log.Logger.Formatter = &logrus.TextFormatter{}
	}

	if c.Debug {
		log.Logger.Level = logrus.DebugLevel
		log.Debug("debug logging enabled")
	} else {
		log.Logger.Level = logrus.InfoLevel
	}

	return log
}

type fallbackFormatter struct {
	baseFormatter logrus.Formatter
}

func (f *fallbackFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// try to use the normal one
	output, err := f.baseFormatter.Format(entry)
	if err == nil {
		return output, nil
	}

	// fall back to pretty printing the entry data
	entry = entry.Dup()
	entry.Data = logrus.Fields{"_fallback_data": pretty.Sprint(entry.Data), "_format_error": pretty.Sprint(err)}
	output, err = f.baseFormatter.Format(entry)
	if err == nil {
		return output, nil
	}

	// then fall back to pretty printing the entire entry
	return []byte(pretty.Sprint(entry)), nil
}

// ContextLogger wraps logrus Entry to allow field mutation, which means the
// context itself can store a pointer to a ContextLogger, so it doesn't need
// replacing each time new fields are added to the logger
type ContextLogger struct {
	entry            *logrus.Entry
	timeoutsAsErrors bool
}

// NewContextLogger creates a new (mutable) ContextLogger instance from an (immutable) logrus Entry
func NewContextLogger(log *logrus.Entry) *ContextLogger {
	return &ContextLogger{entry: log}
}

// GetLogger returns (an immutable) logrus entry from a (mutable) ContextLogger instance
func (l *ContextLogger) GetLogger() *logrus.Entry {
	return l.entry
}

// SetField updates the internal field map
func (l *ContextLogger) SetField(field string, value any) *ContextLogger {
	l.entry = l.entry.WithField(field, value)
	return l
}

// SetFields updates the internal field map with multiple fields at a time
func (l *ContextLogger) SetFields(fields logrus.Fields) *ContextLogger {
	l.entry = l.entry.WithFields(fields)
	return l
}

// SetError updates the internal error
func (l *ContextLogger) SetError(err error) *ContextLogger {
	l.entry = l.entry.WithError(err)
	return l
}

// getContextLogger retrieves the ContextLogger instance from the context
func getContextLogger(ctx context.Context) *ContextLogger {
	ctxLogger, _ := ctx.Value(loggerKey).(*ContextLogger)
	return ctxLogger
}

func mustGetContextLogger(ctx context.Context) *ContextLogger {
	ctxLogger := getContextLogger(ctx)
	if ctxLogger != nil {
		return ctxLogger
	}

	panic("no clog exists in the context")
}

// Set sets a persistent, mutable logger for the request context.
func Set(ctx context.Context, log *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, NewContextLogger(log))
}

// Get retrieves the logrus Entry from the ContextLogger in a context
// and returns a new logrus Entry if none is found
func Get(ctx context.Context) *logrus.Entry {
	ctxLogger := getContextLogger(ctx)
	if ctxLogger != nil {
		return ctxLogger.GetLogger()
	}

	logger := logrus.NewEntry(logrus.New())

	logger.Warn("no clog exists in the context")

	return logger
}

// SetField adds or updates a field to the ContextLogger in a context
func SetField(ctx context.Context, field string, value any) {
	mustGetContextLogger(ctx).SetField(field, value)
}

// SetFields adds or updates fields to the ContextLogger in a context
func SetFields(ctx context.Context, fields Fields) {
	mustGetContextLogger(ctx).SetFields(logrus.Fields(fields))
}

// SetError adds or updates an error to the ContextLogger in a context
func SetError(ctx context.Context, err error) {
	ctxLogger := mustGetContextLogger(ctx)

	ctxLogger.SetError(err)
}

// ConfigureTimeoutsAsErrors changes to default behavior of logging timeouts as info, to log them as errors
func ConfigureTimeoutsAsErrors(ctx context.Context) {
	mustGetContextLogger(ctx).timeoutsAsErrors = true
}

// TimeoutsAsErrors determines whether ConfigureTimeoutsAsErrors was called on the context
func TimeoutsAsErrors(ctx context.Context) bool {
	return mustGetContextLogger(ctx).timeoutsAsErrors
}

// DetermineLevel returns a suggested logrus Level type for a given error
func DetermineLevel(err error, timeoutsAsErrors bool) logrus.Level {
	if cherError, ok := errors.AsType[cher.E](err); ok {
		if cherError.StatusCode() >= 500 {
			return logrus.ErrorLevel
		}

		switch cherError.Code {
		case cher.ContextCanceled:
			return levelForContextCancelation(timeoutsAsErrors)

		default:
			return logrus.WarnLevel
		}
	}

	if strings.Contains(err.Error(), "canceling statement due to user request") {
		return levelForContextCancelation(timeoutsAsErrors)
	}

	// non-cher errors are "unhandled" so warrant an error
	return logrus.ErrorLevel
}

func levelForContextCancelation(timeoutsAsErrors bool) logrus.Level {
	if timeoutsAsErrors {
		return logrus.ErrorLevel
	}

	return logrus.InfoLevel
}
