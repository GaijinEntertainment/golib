package logger

import (
	"math"

	"github.com/GaijinEntertainment/golib/fields"
)

const (
	LevelError = iota*10 + 10
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace

	DefaultLogLevel = math.MaxInt
)

type Logger struct {
	// maxLevel is a maximum log-level of logger, assuming that log-levels are
	// ordered from the most important to the least important, meaning the higher
	// log-level - the less important a log message is.
	//
	// In case log-level is higher than defined maximum, it won't be passed to
	// adapter.
	maxLevel int
}

// New creates new [Logger]. Consult [Logger] docs for more info about
// parameters.
func New(maxLevel int) Logger {
	return Logger{
		maxLevel: maxLevel,
	}
}

// Error logs a message with the [LevelError] log-level.
//
// Use Error to log any unrecoverable error, such as a database query failure
// where the application cannot continue. It's OK to pass nil as the error.
// To attach a stack trace, use [Logger.WithStackTrace].
func (l Logger) Error(msg string, err error, fs ...fields.Field) {
	l.Log(LevelError, msg, err, fs...)
}

// Warning logs a message with the [LevelWarning] log-level.
//
// Use WarningE to log any recoverable error, such as an error during a remote
// API call where the service did not respond and the application will retry.
func (l Logger) Warning(msg string, fs ...fields.Field) {
	l.Log(LevelWarning, msg, nil, fs...)
}

// WarningE logs a message with the [LevelWarning] log-level and the provided
// error.
//
// Use WarningE to log any recoverable error, such as an error during a remote
// API call where the service did not respond and the application will retry.
func (l Logger) WarningE(msg string, err error, fs ...fields.Field) {
	l.Log(LevelWarning, msg, err, fs...)
}

// Info logs a message with the [LevelInfo] log-level.
//
// Use Info to log informational messages that highlight the progress of the
// application.
func (l Logger) Info(msg string, fs ...fields.Field) {
	l.Log(LevelInfo, msg, nil, fs...)
}

// InfoE logs a message with the [LevelInfo] log-level and the provided error.
//
// Use InfoE to log informational messages that highlight the progress of the
// application along with an error.
func (l Logger) InfoE(msg string, err error, fs ...fields.Field) {
	l.Log(LevelInfo, msg, err, fs...)
}

// Debug logs a message with the [LevelDebug] log-level.
//
// Use Debug to log detailed information that is useful during development and
// debugging.
func (l Logger) Debug(msg string, fs ...fields.Field) {
	l.Log(LevelDebug, msg, nil, fs...)
}

// DebugE logs a message with the [LevelDebug] log-level and the provided error.
//
// Use DebugE to log detailed information that is useful during development and
// debugging along with an error.
func (l Logger) DebugE(msg string, err error, fs ...fields.Field) {
	l.Log(LevelDebug, msg, err, fs...)
}

// Trace logs a message with the [LevelTrace] log-level.
//
// Use Trace to log very detailed information, typically of interest only when
// diagnosing problems.
func (l Logger) Trace(msg string, fs ...fields.Field) {
	l.Log(LevelTrace, msg, nil, fs...)
}

// TraceE logs a message with the [LevelTrace] log-level and the provided error.
//
// Use TraceE to log very detailed information, typically of interest only when
// diagnosing problems along with an error.
func (l Logger) TraceE(msg string, err error, fs ...fields.Field) {
	l.Log(LevelTrace, msg, err, fs...)
}

// Log logs a message with given log-level, optional error and fields.
func (l Logger) Log(level int, _ string, _ error, _ ...fields.Field) {
	if level > l.maxLevel {
		return
	}
}

// WithFields returns a new child-logger with the given fields attached to it.
func (l Logger) WithFields(_ ...fields.Field) Logger {
	// ToDo: implement attaching fields to the logger
	return l
}

// WithStackTrace returns a new child-logger with the stack trace attached to it.
func (l Logger) WithStackTrace(_ uint) Logger {
	// ToDo: implement attaching actual stack trace to the logger
	return l
}

// WithName returns a new child-logger with the given name assigned to it.
func (l Logger) WithName(_ string) Logger {
	// ToDo: implement assigning a name to the logger
	return l
}

// Flush flushes the underlying logger adapter, allowing buffered adapters to
// write logs to the output.
//
// It is the application's responsibility to call [Logger.Flush] before exiting.
func (Logger) Flush() error {
	return nil
}
