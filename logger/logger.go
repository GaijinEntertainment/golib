package logger

import (
	"fmt"
	"math"
	"reflect"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/stacktrace"
)

const (
	LevelError = iota*10 + 10
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
)

// Option is a functional option for configuring [Logger] behavior.
type Option func(*Logger)

// WithLevel sets the maximum log-level for the logger.
//
// Log messages with levels higher than the specified level will be ignored
// (not passed to the underlying adapter).
func WithLevel(level int) Option {
	return func(l *Logger) {
		l.maxLevel = level
	}
}

// WithCallerAtLevel enables automatic caller information capture for log
// entries with level less or equal passed threshold.
//
// When enabled, the logger will automatically capture and include caller
// information (formatted as "/path/to/file:line") as a field in log entries whose
// level is less or equal the threshold. For example, WithCallerAtLevel(LevelWarning)
// will add caller information to Error and Warning logs, but no others.
//
// By default, automatic caller capture is disabled. Use this option when you
// need to track the source location of important log messages like errors and
// warnings.
func WithCallerAtLevel(level int) Option {
	return func(l *Logger) {
		l.callerMaxLevel = level
	}
}

// mappers holds mapping configurations for [Logger].
type mappers struct {
	name       func(name string) fields.Field
	error      func(err error) fields.Field
	stackTrace func(st *stacktrace.Stack) fields.Field
	caller     func(frame stacktrace.Frame) fields.Field
}

type mapperOption func(*mappers)

// withMapperOptions applies mapper options to the logger.
func withMapperOption(opt mapperOption) Option {
	return func(l *Logger) {
		opt(l.mappers)
	}
}

// DefaultNameMapper is the default mapper for converting logger names to fields.
// It creates a field with key "logger-name" and the name as value.
func DefaultNameMapper(name string) fields.Field {
	return fields.F("logger-name", name)
}

// WithNameMapper sets a custom name mapper for the logger.
//
// The name mapper controls how logger names (set via [Logger.WithName]) are
// converted to fields. By default, [DefaultNameMapper] is used, which creates
// fields with the key "logger-name".
func WithNameMapper(fn func(name string) fields.Field) Option {
	return withMapperOption(func(cfg *mappers) {
		cfg.name = fn
	})
}

// NameFormatterHierarchical is the default formatter for combining new logger name
// with existing, using `:` separator (e.g., "parent:child").
func NameFormatterHierarchical(prev, next string) string {
	const sep = ":"

	if prev == "" {
		return next
	}

	return prev + sep + next
}

// NameFormatterReplaced is a name formatter that always replaces the previous
// name with the new name.
func NameFormatterReplaced(_, next string) string {
	return next
}

// WithNameFormatter sets a custom name formatter for the logger.
//
// The name formatter controls how logger names are combined when
// [Logger.WithName] is called multiple times on a logger chain. By default,
// [NameFormatterHierarchical] is used.
//
// Custom formatters can implement hierarchical naming (e.g., "parent:child") or
// other naming strategies. The formatter receives the previous name (empty
// string if no name was set) and the next name, and returns the final name.
func WithNameFormatter(fn func(prev, next string) string) Option {
	return func(l *Logger) {
		l.nameFormatter = fn
	}
}

// DefaultErrorMapper is the default mapper for converting errors to fields.
// It creates a field with key "error" and the error message string as value.
//
// The mapper calls err.Error() to obtain the string representation, which may
// panic for improperly implemented error types (e.g., nil pointer with value
// receiver). When this occurs, the mapper recovers and logs either "<nil>" or
// "<PANIC=...>" with the panic details.
//
// Example output:
//
//	err := nil                    // "<nil>"
//	err := errors.New("failed")   // "failed"
//	err := (*CustomErr)(nil)      // "<nil>" (typed nil)
//	err := &badErr{}              // "<PANIC=...>" if Error() panics
func DefaultErrorMapper(err error) (f fields.Field) {
	const nilErr = "<nil>"

	// we have a possibility of panicking while calling err.Error(), so in order to
	// keep logging alive and sane we have to recover here.
	defer func() {
		if re := recover(); re != nil {
			// most likely provided error is a nil pointer for value receiver type (famous
			// nil != nil problem) in this case we just log "<nil>" as error text
			if v := reflect.ValueOf(err); v.Kind() == reflect.Ptr && v.IsNil() {
				f.V = nilErr
				return
			}

			// otherwise we log the panic info and error value
			f.V = fmt.Sprintf("<PANIC=%v> %#v", re, err)
		}
	}()

	f.K = "error"
	f.V = nilErr

	if err != nil {
		f.V = err.Error()
	}

	return f
}

// WithErrorMapper sets a custom error mapper for the logger.
//
// The error mapper controls how errors passed to logging methods are converted
// to fields. By default, [DefaultErrorMapper] is used, which creates fields with
// the key "error" and calls err.Error() to get the string representation.
//
// Custom mappers can change the field key, format errors differently, or add
// additional error information (such as error types or stack traces from error
// implementations that include them).
func WithErrorMapper(fn func(err error) fields.Field) Option {
	return withMapperOption(func(cfg *mappers) {
		cfg.error = fn
	})
}

// DefaultStackTraceMapper is the default mapper for converting stack traces to
// fields. It creates a field with key "stacktrace" and the stack trace string
// representation as value.
//
// The mapper calls st.String() to obtain the formatted stack trace string.
func DefaultStackTraceMapper(st *stacktrace.Stack) fields.Field {
	return fields.F("stacktrace", st.String())
}

// WithStackTraceMapper sets a custom stack trace mapper for the logger.
//
// The stack trace mapper controls how stack traces (captured via
// [Logger.WithStackTrace]) are converted to fields. By default,
// [DefaultStackTraceMapper] is used, which creates fields with the key
// "stacktrace" and formats the stack as a string.
//
// Custom mappers can change the field key, format stack traces differently, or
// limit the number of frames included.
func WithStackTraceMapper(fn func(st *stacktrace.Stack) fields.Field) Option {
	return withMapperOption(func(cfg *mappers) {
		cfg.stackTrace = fn
	})
}

// DefaultCallerMapper is the default mapper for converting caller frames to
// fields. It creates a field with key "caller" and formats the frame as
// "file:line".
//
// Example output: "/path/to/logger.go:123".
func DefaultCallerMapper(frame stacktrace.Frame) fields.Field {
	return fields.F("caller", frame.FullPath())
}

// WithCallerMapper sets a custom caller mapper for the logger.
//
// The caller mapper controls how caller frames (captured via automatic caller
// capture with [WithCallerAtLevel]) are converted to fields. By default,
// [DefaultCallerMapper] is used, which creates fields with the key "caller"
// and formats frames as "file:line".
//
// Custom mappers can change the field key, format caller information
// differently, or include additional frame details like function names.
func WithCallerMapper(fn func(frame stacktrace.Frame) fields.Field) Option {
	return withMapperOption(func(cfg *mappers) {
		cfg.caller = fn
	})
}

// defaultMappers returns a mappers structure initialized with default mappers.
func defaultMappers() *mappers {
	return &mappers{
		name:       DefaultNameMapper,
		error:      DefaultErrorMapper,
		stackTrace: DefaultStackTraceMapper,
		caller:     DefaultCallerMapper,
	}
}

type Logger struct {
	// maxLevel is a maximum log-level of logger, assuming that log-levels are
	// ordered from the most important to the least important.
	//
	// In case log-level is higher than defined maximum, operation will be no-op.
	maxLevel int

	adapter Adapter

	// in opposition to logger itself, mappers are used by-pointer since it never
	// changes and there is no need to copy it on every method call.
	mappers *mappers

	name          string
	nameFormatter func(prev, next string) string

	// callerMaxLevel is the maximum log-level at which caller information is
	// automatically captured and added to log entries. Levels at or below this
	// threshold will include caller information. Set to -1 to disable.
	callerMaxLevel int
}

// New creates new [Logger] with maximum log-level set to LevelInfo and default
// mappers. To change log level and other behaviours use options.
//
// Panics if adapter is nil. To create no-op logger use [NewNop].
func New(adapter Adapter, opts ...Option) Logger {
	if adapter == nil {
		panic("logger adapter cannot be nil, use logger.NewNop() to create no-op logger")
	}

	l := Logger{
		maxLevel:       LevelInfo,
		adapter:        adapter,
		mappers:        defaultMappers(),
		name:           "",
		nameFormatter:  NameFormatterHierarchical,
		callerMaxLevel: math.MinInt,
	}

	lp := &l
	for _, opt := range opts {
		opt(lp)
	}

	return l
}

// NewNop creates a new logger that does nothing. Methods creating child-loggers
// will also create no-op loggers.
func NewNop() Logger {
	return Logger{
		maxLevel:       math.MaxInt,
		adapter:        nil,
		mappers:        nil,
		name:           "",
		nameFormatter:  nil,
		callerMaxLevel: math.MinInt,
	}
}

// IsNop returns true if the logger is a no-op logger.
//
// A no-op logger is created via [NewNop] and performs no operations. All
// logging methods and child logger creation methods return immediately without
// calling the underlying adapter. This is useful for testing or when a logger
// is required but logging is not desired.
func (l Logger) IsNop() bool {
	return l.adapter == nil
}

// Error logs a message with the [LevelError] log-level.
//
// Use Error to log any unrecoverable error, such as a database query failure
// where the application cannot continue. It's OK to pass nil as the error.
// To attach a stack trace, use [Logger.WithStackTrace].
func (l Logger) Error(msg string, err error, fs ...fields.Field) {
	l.log(LevelError, msg, err, fs...)
}

// Warning logs a message with the [LevelWarning] log-level.
//
// Use Warning to log any recoverable error or concerning situation that doesn't
// prevent the application from continuing, such as a deprecated API usage or
// a retry-able failure. For warnings with an error, use [Logger.WarningE].
func (l Logger) Warning(msg string, fs ...fields.Field) {
	l.log(LevelWarning, msg, nil, fs...)
}

// WarningE logs a message with the [LevelWarning] log-level and the provided
// error.
//
// Use WarningE to log any recoverable error, such as an error during a remote
// API call where the service did not respond and the application will retry.
func (l Logger) WarningE(msg string, err error, fs ...fields.Field) {
	l.log(LevelWarning, msg, err, fs...)
}

// Info logs a message with the [LevelInfo] log-level.
//
// Use Info to log informational messages that highlight the progress of the
// application.
func (l Logger) Info(msg string, fs ...fields.Field) {
	l.log(LevelInfo, msg, nil, fs...)
}

// InfoE logs a message with the [LevelInfo] log-level and the provided error.
//
// Use InfoE to log informational messages that highlight the progress of the
// application along with an error.
func (l Logger) InfoE(msg string, err error, fs ...fields.Field) {
	l.log(LevelInfo, msg, err, fs...)
}

// Debug logs a message with the [LevelDebug] log-level.
//
// Use Debug to log detailed information that is useful during development and
// debugging.
func (l Logger) Debug(msg string, fs ...fields.Field) {
	l.log(LevelDebug, msg, nil, fs...)
}

// DebugE logs a message with the [LevelDebug] log-level and the provided error.
//
// Use DebugE to log detailed information that is useful during development and
// debugging along with an error.
func (l Logger) DebugE(msg string, err error, fs ...fields.Field) {
	l.log(LevelDebug, msg, err, fs...)
}

// Trace logs a message with the [LevelTrace] log-level.
//
// Use Trace to log very detailed information, typically of interest only when
// diagnosing problems.
func (l Logger) Trace(msg string, fs ...fields.Field) {
	l.log(LevelTrace, msg, nil, fs...)
}

// TraceE logs a message with the [LevelTrace] log-level and the provided error.
//
// Use TraceE to log very detailed information, typically of interest only when
// diagnosing problems along with an error.
func (l Logger) TraceE(msg string, err error, fs ...fields.Field) {
	l.log(LevelTrace, msg, err, fs...)
}

// Log logs a message with the given log-level, optional error, and fields.
//
// The level parameter expected to be one of the logger level constants
// (LevelError, LevelWarning, LevelInfo, LevelDebug, LevelTrace). Usage of custom
// levels is also supported, but heavily depends on the underlying adapter
// capabilities.
//
// If the level is higher than the logger's maximum level (set via [WithLevel]),
// the message is not logged. If err is not nil, it is converted to a field using
// the error mapper and appended to the provided fields.
//
// For no-op loggers, this method returns immediately without any operation.
func (l Logger) Log(level int, msg string, err error, fs ...fields.Field) {
	l.log(level, msg, err, fs...)
}

// log is the internal logging method, its sole reason to exist is to uniformly
// catch caller frame in case it is required.
//
//revive:disable-next-line:confusing-naming
func (l Logger) log(level int, msg string, err error, fs ...fields.Field) {
	if level > l.maxLevel || l.IsNop() {
		return
	}

	if level <= l.callerMaxLevel {
		const callerSkip = 2 // skip log and the calling method (Error, Info, etc.)

		frame := stacktrace.CaptureCaller(callerSkip)

		fs = append(fs, l.mappers.caller(frame))
	}

	if l.name != "" {
		fs = append(fs, l.mappers.name(l.name))
	}

	if err != nil {
		fs = append(fs, l.mappers.error(err))
	}

	l.adapter.Log(level, msg, fs...)
}

// WithFields returns a new child logger with the given fields attached to it.
//
// The returned child logger will include these fields in all subsequent log
// entries, in addition to any fields already attached to the parent logger.
// The parent logger remains unaffected.
//
// For no-op loggers, this method returns the same no-op logger.
func (l Logger) WithFields(fs ...fields.Field) Logger {
	if l.IsNop() {
		return l
	}

	//revive:disable-next-line:modifies-value-receiver
	l.adapter = l.adapter.WithFields(fs...)

	return l
}

const stackTraceDepth = 32

// WithStackTrace returns a new child-logger with the stack trace attached to it.
// The skip parameter defines how many stack frames to skip when capturing the
// stack trace. skip=0 means the caller of WithStackTrace is included in the
// stack trace.
//
// For no-op loggers, this method returns the same no-op logger.
func (l Logger) WithStackTrace(skip int) Logger {
	if l.IsNop() {
		return l
	}

	l.adapter = l.adapter.WithFields(l.mappers.stackTrace(
		stacktrace.CaptureStack(skip+1, stackTraceDepth),
	))

	return l
}

// WithName returns a new child logger with the given name assigned to it.
//
// The name is converted to a field using the name mapper (by default, creates
// a field with key "logger-name") and attached to all subsequent log entries
// from this logger. This is useful for identifying which component or module
// generated a log entry.
//
// When called multiple times on a logger chain, the name formatter (set via
// [WithNameFormatter]) determines how names are combined. By default, names
// are concatenated with a `:` separator (e.g., "parent:child").
//
// The parent logger remains unaffected. For no-op loggers, this method returns
// the same no-op logger.
func (l Logger) WithName(name string) Logger {
	if l.IsNop() {
		return l
	}

	//revive:disable-next-line:modifies-value-receiver
	l.name = l.nameFormatter(l.name, name)

	return l
}

// Flush flushes the underlying logger adapter, allowing buffered adapters to
// write logs to the output.
//
// It is the application's responsibility to call [Logger.Flush] before exiting.
func (l Logger) Flush() error {
	if l.IsNop() {
		return nil
	}

	return l.adapter.Flush() //nolint:wrapcheck
}

// IsZero returns true if the logger is a zero-value structure.
//
// A zero-value logger is an uninitialized Logger struct that has not been
// created via [New] or [NewNop]. Zero-value loggers should not be used, but in
// rare cases it is required to check if value is not initialized.
func (l Logger) IsZero() bool {
	// properly created loggers have adapters set and callerMaxLevel initialized,
	// therefore check for those fields should suffice.
	return l.adapter == nil && l.callerMaxLevel == 0
}

// formatterIsEqual checks if f1 is equal to f2.
func formatterIsEqual(f1, f2 func(string, string) string) bool {
	if f1 == nil && f2 == nil {
		return true
	}

	if f1 == nil || f2 == nil {
		return false
	}

	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer()
}

// IsEqual returns true if two loggers are functionally equal. Two loggers are
// considered equal if they have the same maxLevel, adapter, mappers, name, and
// nameFormatter.
func IsEqual(l1, l2 Logger) bool {
	return l1.maxLevel == l2.maxLevel &&
		l1.adapter == l2.adapter &&
		l1.name == l2.name &&
		l1.mappers == l2.mappers &&
		formatterIsEqual(l1.nameFormatter, l2.nameFormatter)
}

// NewErrorLogger creates a new [e.ErrorLogger] that logs errors with the
// given log-level.
//
// This function is useful for scenarios where the log level is configured after
// logger creation, such as in middleware or when adapting to interfaces that
// expect an error logging function rather than a full logger.
func NewErrorLogger(lgr Logger, level int) e.ErrorLogger {
	return func(msg string, err error, fs ...fields.Field) {
		lgr.log(level, msg, err, fs...)
	}
}
