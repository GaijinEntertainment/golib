// Package logger provides a structured logging abstraction with support for
// multiple backends, customizable field mappers, and hierarchical logger contexts.
//
// The logger package is designed around three core concepts: loggers, adapters,
// and mappers. Loggers provide the user-facing API, adapters bridge to concrete
// logging implementations, and mappers control how special values (names, errors,
// stack traces) are converted to log fields.
//
// # Basic Usage
//
// Create a logger with an adapter and start logging:
//
//	import (
//		"dev.gaijin.team/go/golib/logger"
//		"dev.gaijin.team/go/golib/logger/bufferadapter"
//		"dev.gaijin.team/go/golib/fields"
//	)
//
//	func main() {
//		adapter, _ := bufferadapter.New()
//		lgr := logger.New(adapter)
//
//		lgr.Info("application started", fields.F("version", "1.0.0"))
//		lgr.Error("operation failed", err, fields.F("operation", "save"))
//
//		lgr.Flush()
//	}
//
// # Log Levels
//
// The logger supports five standard log levels, ordered from most to least
// important:
//
//	LevelError   - Unrecoverable errors that prevent operation completion
//	LevelWarning - Recoverable errors or concerning situations
//	LevelInfo    - Informational messages about application progress (default)
//	LevelDebug   - Detailed information useful during development
//	LevelTrace   - Very detailed information for diagnosing problems
//
// By default, loggers are created with a maximum level of LevelInfo, meaning
// Debug and Trace messages are filtered out. Use WithLevel to change this:
//
//	lgr := logger.New(adapter, logger.WithLevel(logger.LevelDebug))
//
// # Automatic Caller Capture
//
// The logger can automatically capture and include caller information (file and
// line number) for log entries at or below a specified level threshold. This is
// useful for tracking the source location of important log messages like errors
// and warnings without manually adding caller information.
//
// Enable automatic caller capture using WithCallerAtLevel:
//
//	// Capture caller for Error and Warning logs only
//	lgr := logger.New(adapter, logger.WithCallerAtLevel(logger.LevelWarning))
//	lgr.Error("operation failed", err)  // includes caller: "/path/to/file.go:123"
//	lgr.Info("processing")               // no caller information
//
// By default, automatic caller capture is disabled. When enabled, the caller
// information is formatted and added as a field to log entries using the caller
// mapper (see Customization section below).
//
// Each level has two methods: one without an error parameter (Info, Warning,
// Debug, Trace) and one with an error parameter (InfoE, WarningE, DebugE,
// TraceE). The Error is, obviously, singular and has the error parameter, though
// it is okay to provide nil as the error.
//
// # Child Loggers and Context
//
// Child loggers allow you to build up contextual information that's automatically
// included in all subsequent log entries. The parent logger remains unaffected:
//
//	// Add fields to a child logger
//	requestLogger := lgr.WithFields(
//		fields.F("request_id", "abc123"),
//		fields.F("user_id", 42),
//	)
//	requestLogger.Info("processing request")  // includes request_id and user_id
//	lgr.Info("other operation")               // does not include request fields
//
//	// Add a logger name for component identification
//	dbLogger := lgr.WithName("database")
//	dbLogger.Info("connection established")  // includes logger-name: database
//
//	// Build hierarchical names by chaining WithName calls
//	// By default, names are joined with ":" separator
//	serviceLogger := lgr.WithName("service")
//	handlerLogger := serviceLogger.WithName("handler")
//	handlerLogger.Info("processing")  // includes logger-name: service:handler
//
//	// Attach stack trace to debug issues
//	lgr.WithStackTrace(0).Error("unexpected error", err)
//
// Child logger methods can be chained to combine multiple contexts:
//
//	apiLogger := lgr.
//		WithName("api").
//		WithFields(fields.F("version", "v2")).
//		WithStackTrace(0)
//
// # Adapters
//
// Adapters bridge the logger abstraction to concrete logging implementations.
// The logger package doesn't include any output mechanism itself - all log
// output is delegated to adapters.
//
// Common adapters:
//
//   - bufferadapter: In-memory buffer for testing
//   - zapadapter: Integration with uber-go/zap
//   - slogadapter: Integration with log/slog (Go 1.21+)
//
// To create a custom adapter, implement the [logger.Adapter] interface.
//
// # Customization with Mappers and Formatters
//
// Though the logger supports concepts of logger names, errors logging and stack
// traces, etc. - different logging backends may have different conventions for
// how these values are represented in logs or lack of such functionality at all.
// For that reason logger abstracts these concepts passing it to the adapters as
// fields, with the help of mappers.
//
// The logger provides four types of mappers:
//
// Name Mapper: Converts logger names (from WithName) to fields. By default,
// creates a field with key "logger-name":
//
//	lgr := logger.New(adapter, logger.WithNameMapper(func(name string) fields.Field {
//		return fields.F("component", name)  // Use "component" instead of "logger-name"
//	}))
//
// Name Formatter: Controls how logger names are combined when WithName is called
// multiple times. By default, uses NameFormatterHierarchical which joins names
// with ":" separator. Use NameFormatterReplaced to replace names instead:
//
//	// Hierarchical naming (default): "service:handler:method"
//	lgr := logger.New(adapter)
//	lgr.WithName("service").WithName("handler").WithName("method")
//
//	// Replacement naming: only "method"
//	lgr := logger.New(adapter, logger.WithNameFormatter(logger.NameFormatterReplaced))
//	lgr.WithName("service").WithName("handler").WithName("method")
//
// Error Mapper: Converts errors to fields. By default, creates a field with
// key "error" and calls err.Error() to get the string representation. The
// default mapper includes panic recovery for improperly implemented error types:
//
//	lgr := logger.New(adapter, logger.WithErrorMapper(func(err error) fields.Field {
//		return fields.F("error", fmt.Sprintf("%+v", err))  // Use value of error with %+v
//	}))
//
// Stack Trace Mapper: Converts stack traces to fields. By default, creates a
// field with key "stacktrace":
//
//	lgr := logger.New(adapter, logger.WithStackTraceMapper(func(st *stacktrace.Stack) fields.Field {
//		return fields.F("stack", st.String())  // Use "stack" instead of "stacktrace"
//	}))
//
// Caller Mapper: Converts caller frames (from automatic caller capture) to
// fields. By default, creates a field with key "caller" formatted as
// "file:line":
//
//	lgr := logger.New(adapter,
//		logger.WithCallerAtLevel(logger.LevelError),  // Enable automatic capture
//		logger.WithCallerMapper(func(frame stacktrace.Frame) fields.Field {
//			return fields.F("source", frame.ShortPath())  // Use "source" and short path
//		}),
//	)
//
// All mappers and formatters can be customized independently during logger
// creation.
//
// # No-Op Logger
//
// For scenarios where logging is not desired (testing, optional logging), use
// a no-op logger that performs no operations:
//
//	lgr := logger.NewNop()
//	lgr.Info("this does nothing")
//
//	// Check if a logger is no-op
//	if lgr.IsNop() {
//		// Skip expensive logging preparation
//	}
//
// No-op loggers propagate through child logger creation - calling WithFields,
// WithName, or WithStackTrace on a no-op logger returns another no-op logger.
//
// # Context Integration
//
// The logger provides context.Context integration for passing loggers through
// your application:
//
//	import "context"
//
//	// Store logger in context
//	ctx := logger.ToCtx(context.Background(), lgr)
//
//	// Retrieve logger from context
//	lgr, ok := logger.FromCtx(ctx)
//	if !ok {
//		// Handle missing logger
//	}
//
//	// Retrieve logger or get no-op if missing
//	lgr := logger.FromCtxOrNop(ctx)
//
// This allows passing loggers to functions without explicit parameters while
// maintaining type safety.
//
// # Error Logging Adapter
//
// For integration with interfaces that expect a simple error logging function,
// use NewErrorLogger:
//
//	import "dev.gaijin.team/go/golib/e"
//
//	// Create error logger with specific level
//	errorLogger := logger.NewErrorLogger(lgr, logger.LevelError)
//
//	// Use as e.ErrorLogger
//	var logFn e.ErrorLogger = errorLogger
//	logFn("operation failed", err, fields.F("operation", "save"))
//
// This is useful when working with libraries that accept error logging callbacks
// or when implementing interfaces with logging requirements.
//
// # Flushing
//
// Some logging backends are buffered and require flushing in order to ensure all
// log entries are written out. Logger provides a Flush method to trigger backend
// flush, though it is application responsibility to call it at appropriate times:
//
//	defer lgr.Flush()
//
// For unbuffered adapters, Flush is a no-op. It's safe-by-contract to call Flush
// even if the adapter doesn't buffer its output.
package logger
