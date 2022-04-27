package stacktrace

import (
	"fmt"
	"runtime"
)

func GetCallerCodePath(fallback string) string {
	return GetCallerCodePathWithSkip(1, fallback)
}

func GetCallerCodePathWithSkip(skip int, fallback string) string {
	_, file, line, ok := runtime.Caller(1 + skip)
	if !ok || file == "" {
		return fallback
	}
	return fmt.Sprintf("%s:%d", file, line)
}
