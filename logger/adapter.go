package logger

import (
	"dev.gaijin.team/go/golib/fields"
)

// Adapter is an interface that allows to encapsulate any logger inside [Logger].
type Adapter interface {
	// Log logs a message with the provided level, message, optional error, and an
	// arbitrary number of fields.
	Log(level int, msg string, err error, fs ...fields.Field)

	// WithFields returns a logger adapter with the given fields attached to it.
	WithFields(fs ...fields.Field) Adapter

	// WithName returns a logger adapter with the given name attached to it.
	WithName(name string) Adapter

	// WithStackTrace returns a logger adapter with the stack trace attached to it.
	WithStackTrace(trace string) Adapter

	// Flush all logs to the output. It is a no-op for non-buffered loggers.
	//
	// It is application's responsibility to call this method before it exits.
	Flush() error
}

// NopAdapter is a no-op logger adapter, it does nothing even if any of its
// methods called.
type NopAdapter struct{}

func (NopAdapter) Log(int, string, error, ...fields.Field) {}

func (a NopAdapter) WithFields(...fields.Field) Adapter { return a }

func (a NopAdapter) WithName(string) Adapter { return a }

func (a NopAdapter) WithStackTrace(string) Adapter { return a }

func (NopAdapter) Flush() error { return nil }
