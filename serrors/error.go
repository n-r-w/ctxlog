// Package serrors provides error wrapping with additional context about the function call.
package serrors

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

const (
	callerFuncLvl = 3
	maxFrames     = 4
)

// Errorf creates an error in the format "{pkg.[type.]callerFuncName} " + fmt.Errorf().String().
func Errorf(format string, a ...any) error {
	return errorf(callerFuncLvl, maxFrames, format, a...)
}

// Error - analog of serrors.Errorf("%w", err).
func Error(err error) error {
	return errorf(callerFuncLvl, maxFrames, "%w", err)
}

// Join - analog of errors.Join(err1, err2) with added information about the function call in err2.
func Join(err1, err2 error) error {
	if err2 == nil {
		return err1
	}

	return errors.Join(err1, errorf(callerFuncLvl, maxFrames, "%w", err2))
}

// Joinf - join with formatting.
func Joinf(err error, format string, a ...any) error {
	return errors.Join(err, errorf(callerFuncLvl, maxFrames, format, a...))
}

// New - analog of errors.New().
func New(msg string) error {
	return errorf(callerFuncLvl, maxFrames, "%s", msg)
}

//nolint:unparam // maxFrames is kept for possible future use.
func errorf(callerFuncLvl, maxFrames int, format string, a ...any) error {
	pc := make([]uintptr, maxFrames)
	n := runtime.Callers(callerFuncLvl, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	parts := strings.Split(frame.Function, "/")
	if len(parts) == 0 {
		return fmt.Errorf(format, a...)
	}
	return fmt.Errorf("{"+parts[len(parts)-1]+"."+strconv.Itoa(frame.Line)+"} "+format, a...)
}
