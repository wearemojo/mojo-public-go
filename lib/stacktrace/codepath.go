package stacktrace

import (
	"fmt"
	"runtime"
)

func GetCallerCodePath(fallback string) string {
	_, file, line, ok := runtime.Caller(1)
	if !ok || file == "" {
		return fallback
	}
	return fmt.Sprintf("%s:%d", file, line)
}
