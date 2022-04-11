package stacktrace

import (
	"fmt"
	"runtime"
	"strings"
)

type Frame struct {
	File string `json:"file"`
	Line int    `json:"line"`

	Function string `json:"function"`
}

// GetCallerFrames returns a slice of Frame objects representing the stack
//
// The skip argument is the number of frames to skip before starting to collect
// frames, with 0 identifying the frame for GetCallerFrames itself, and 1
// identifying the caller of GetCallerFrames
//
// No more than 100 frames will be collected
func GetCallerFrames(skip int) []Frame {
	stackBuf := make([]uintptr, 100)
	length := runtime.Callers(skip+1, stackBuf[:]) //nolint:gocritic // required to use as a buffer
	stack := stackBuf[:length]

	framesIter := runtime.CallersFrames(stack)
	frames := make([]Frame, 0, length)

	for {
		frame, more := framesIter.Next()

		frames = append(frames, Frame{
			File: frame.File,
			Line: frame.Line,

			Function: frame.Function,
		})

		if !more {
			break
		}
	}

	return frames
}

func (f Frame) String() string {
	return fmt.Sprintf("%s\n\t%s:%d", f.Function, f.File, f.Line)
}

func FormatFrames(frames []Frame) string {
	var sb strings.Builder

	for _, frame := range frames {
		sb.WriteString(frame.String())
		sb.WriteString("\n")
	}

	return sb.String()
}
