package logger

import (
	"dev.gaijin.team/go/golib/fields"
)

// Adapter is an interface that allows to encapsulate any logger backend inside
// [Logger].
//
// Adapters provide a bridge between the logger abstraction and concrete logging
// implementations (zap, logrus, slog, etc.). They handle the translation of
// logger's generic API calls to backend-specific operations.
//
// For examples of adapter implementations, see the adapters subpackages.
type Adapter interface {
	// Log logs a message with the provided level, message, and fields.
	//
	// The level parameter is one of the logger level constants (LevelError,
	// LevelWarning, LevelInfo, LevelDebug, LevelTrace) or a custom level value.
	//
	// The fs parameter contains zero or more fields that should be attached to the
	// log entry. This may include both fields from WithFields calls and fields
	// passed directly to the Log call.
	Log(level int, msg string, fs ...fields.Field)

	// WithFields returns a new adapter instance with the given fields attached.
	// The returned adapter should include these fields in all subsequent Log calls.
	//
	// This method must return a new adapter instance and must not modify the
	// original adapter.
	//
	// Rationale:
	//
	// Although it might look unnecessary to alter the adapter instead of carrying
	// fields in the logger itself, this design allows adapters to optimize field
	// handling based on their backend capabilities.
	//
	// Some logging backends may support efficient field storage and reuse, while
	// others may require fields to be passed with each log call. For example zap
	// pre-encodes fields and then reuses encoded values, instead of encoding them
	// each time a log entry is created.
	//
	// By attaching fields to the adapter, we enable adapters to implement the
	// most efficient strategy for their backend.
	WithFields(fs ...fields.Field) Adapter

	// Flush flushes any buffered log entries to the output. This is a no-op for
	// non-buffered loggers.
	//
	// It is the application's responsibility to call [Logger.Flush] before exiting
	// to ensure all log entries are written. Adapters should return any errors
	// that occur during flushing.
	Flush() error
}
