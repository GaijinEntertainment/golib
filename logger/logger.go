package logger

import (
	"math"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
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

	adapter Adapter
}

// New creates new [Logger]. Consult [Logger] docs for more info about
// parameters.
func New(adapter Adapter, maxLevel int) Logger {
	return Logger{
		maxLevel: maxLevel,
		adapter:  adapter,
	}
}

// NewNop creates a new logger that does nothing.
func NewNop() Logger {
	return Logger{
		adapter:  NopAdapter{},
		maxLevel: DefaultLogLevel,
	}
}

// IsNop returns true if the logger's adapter is a [NopAdapter].
func (l Logger) IsNop() bool {
	_, ok := l.adapter.(NopAdapter)

	return ok
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
func (l Logger) Log(level int, msg string, err error, fs ...fields.Field) {
	if level > l.maxLevel {
		return
	}

	l.adapter.Log(level, msg, err, fs...)
}

// WithFields returns a new child-logger with the given fields attached to it.
func (l Logger) WithFields(fs ...fields.Field) Logger {
	//revive:disable-next-line:modifies-value-receiver
	l.adapter = l.adapter.WithFields(fs...)

	return l
}

// WithStackTrace returns a new child-logger with the stack trace attached to it.
func (l Logger) WithStackTrace(_ uint) Logger {
	//revive:disable-next-line:modifies-value-receiver
	l.adapter = l.adapter.WithStackTrace("")

	return l
}

// WithName returns a new child-logger with the given name assigned to it.
func (l Logger) WithName(name string) Logger {
	//revive:disable-next-line:modifies-value-receiver
	l.adapter = l.adapter.WithName(name)

	return l
}

// Flush flushes the underlying logger adapter, allowing buffered adapters to
// write logs to the output.
//
// It is the application's responsibility to call [Logger.Flush] before exiting.
func (l Logger) Flush() error {
	return l.adapter.Flush() //nolint:wrapcheck
}

// IsZero returns true if the logger is a zero-value structure.
func (l Logger) IsZero() bool {
	return l.adapter == nil && l.maxLevel == 0
}

// NewErrorLogger creates new [e.ErrorLogger] that allows to log errors with
// given log-level.
//
// The function is most usesful for scenarios where loglevel is configured after
// logger creation, such as in a middlewares.
func NewErrorLogger(log Logger, level int) e.ErrorLogger {
	return func(msg string, err error, fs ...fields.Field) {
		log.Log(level, msg, err, fs...)
	}
}
