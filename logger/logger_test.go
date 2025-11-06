package logger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
	"dev.gaijin.team/go/golib/logger/bufferadapter"
	"dev.gaijin.team/go/golib/stacktrace"
)

func TestMappers(t *testing.T) {
	t.Parallel()

	// Track whether each mapper was called
	var (
		nameMapperCalled       bool
		errorMapperCalled      bool
		stackTraceMapperCalled bool
	)

	// Create custom mappers that wrap defaults and track calls
	customNameMapper := func(name string) fields.Field {
		nameMapperCalled = true
		return logger.DefaultNameMapper(name)
	}

	customErrorMapper := func(err error) fields.Field {
		errorMapperCalled = true
		return logger.DefaultErrorMapper(err)
	}

	customStackTraceMapper := func(st *stacktrace.Stack) fields.Field {
		stackTraceMapperCalled = true
		return logger.DefaultStackTraceMapper(st)
	}

	// Create logger with all custom mappers
	adapter, buff := bufferadapter.New()
	lgr := logger.New(
		adapter,
		logger.WithNameMapper(customNameMapper),
		logger.WithErrorMapper(customErrorMapper),
		logger.WithStackTraceMapper(customStackTraceMapper),
	)

	// Trigger name mapper
	lgr = lgr.WithName("test-logger")

	// Trigger error mapper
	lgr.Error("test message", e.New("test error"))

	// Trigger stack trace mapper
	lgr = lgr.WithStackTrace(0)
	lgr.Info("test with stack")

	// Verify all mappers were called
	assert.True(t, nameMapperCalled, "name mapper should have been called")
	assert.True(t, errorMapperCalled, "error mapper should have been called")
	assert.True(t, stackTraceMapperCalled, "stack trace mapper should have been called")

	// Verify the fields were actually added to logs
	entries := buff.GetAll()
	require.Len(t, entries, 2, "should have 2 log entries")

	// First entry: Error with error and name
	entry1 := entries[0]
	assert.Equal(t, logger.LevelError, entry1.Level)
	assert.Equal(t, "test message", entry1.Msg)
	assert.Contains(t, entry1.Fields, fields.F("logger-name", "test-logger"))
	assert.Contains(t, entry1.Fields, fields.F("error", "test error"))

	// Second entry: Info with name and stack trace
	entry2 := entries[1]
	assert.Equal(t, logger.LevelInfo, entry2.Level)
	assert.Equal(t, "test with stack", entry2.Msg)
	assert.Contains(t, entry2.Fields, fields.F("logger-name", "test-logger"))

	// Check stack trace field exists (value will vary, just check key)
	assert.NotZero(t, entry2.Fields.ToDict()["stacktrace"], "should have stacktrace field")
}

func TestNopLogger(t *testing.T) {
	t.Parallel()

	lgr := logger.NewNop()

	assert.True(t, lgr.IsNop(), "NewNop() should create a nop logger")
	assert.False(t, lgr.IsZero(), "nop logger should not be zero-value")

	assert.NotPanics(t, func() {
		// This test verifies that nop logger truly does nothing by ensuring no
		// interactions with any adapter occur as we know that nop logger has nil adapter
		// - interaction would cause panic
		lgr := logger.NewNop()

		// Log at all levels with various methods
		lgr.Error("error msg", e.New("error"))
		lgr.Warning("warning msg")
		lgr.WarningE("warning msg", e.New("error"))
		lgr.Info("info msg", fields.F("key", "value"))
		lgr.InfoE("info msg", e.New("error"))
		lgr.Debug("debug msg")
		lgr.DebugE("debug msg", e.New("error"))
		lgr.Trace("trace msg")
		lgr.TraceE("trace msg", e.New("error"))

		// Add child loggers and log
		lgr.WithName("test").Info("should not log")
		lgr.WithFields(fields.F("k", "v")).Error("should not log", e.New("err"))
		lgr.WithStackTrace(0).Debug("should not log")

		_ = lgr.Flush()
	})

	assert.NoError(t, lgr.Flush())
	assert.True(t, lgr.WithFields(fields.F("key", "value")).IsNop(),
		"WithFields should return nop logger")
	assert.True(t, lgr.WithName("test-name").IsNop(),
		"WithName should return nop logger")
	assert.True(t, lgr.WithStackTrace(0).IsNop(),
		"WithStackTrace should return nop logger")
}

func TestWithLevel(t *testing.T) {
	t.Parallel()

	t.Run("default log level", func(t *testing.T) {
		t.Parallel()

		adapter, buff := bufferadapter.New()
		lgr := logger.New(adapter)

		// Log at all levels
		lgr.Error("error message", nil)
		lgr.Warning("warning message")
		lgr.Info("info message")
		lgr.Debug("debug message")
		lgr.Trace("trace message")

		entries := buff.GetAll()

		// Default level is LevelInfo, so Error, Warning, and Info should be logged
		// but Debug and Trace should be filtered out
		require.Len(t, entries, 3, "should have 3 log entries (Error, Warning, Info)")

		assert.Equal(t, logger.LevelError, entries[0].Level)
		assert.Equal(t, "error message", entries[0].Msg)

		assert.Equal(t, logger.LevelWarning, entries[1].Level)
		assert.Equal(t, "warning message", entries[1].Msg)

		assert.Equal(t, logger.LevelInfo, entries[2].Level)
		assert.Equal(t, "info message", entries[2].Msg)
	})

	t.Run("custom log level", func(t *testing.T) {
		t.Parallel()

		adapter, buff := bufferadapter.New()
		lgr := logger.New(adapter, logger.WithLevel(logger.LevelTrace))

		// Log at all levels
		lgr.Error("error message", nil)
		lgr.Warning("warning message")
		lgr.Info("info message")
		lgr.Debug("debug message")
		lgr.Trace("trace message")

		entries := buff.GetAll()

		// All levels should be logged
		require.Len(t, entries, 5, "should have 5 log entries (all levels)")

		assert.Equal(t, logger.LevelError, entries[0].Level)
		assert.Equal(t, "error message", entries[0].Msg)

		assert.Equal(t, logger.LevelWarning, entries[1].Level)
		assert.Equal(t, "warning message", entries[1].Msg)

		assert.Equal(t, logger.LevelInfo, entries[2].Level)
		assert.Equal(t, "info message", entries[2].Msg)

		assert.Equal(t, logger.LevelDebug, entries[3].Level)
		assert.Equal(t, "debug message", entries[3].Msg)

		assert.Equal(t, logger.LevelTrace, entries[4].Level)
		assert.Equal(t, "trace message", entries[4].Msg)
	})
}

// errWithValueReceiver is a custom error type with a value receiver that will
// panic when Error() is called on a nil pointer.
type errWithValueReceiver struct {
	msg string
}

func (e errWithValueReceiver) Error() string {
	return e.msg
}

// errThatPanics is a custom error type that always panics in Error() method.
type errThatPanics struct{}

func (*errThatPanics) Error() string {
	panic("intentional panic in Error()")
}

func TestErrorMapperInLogging(t *testing.T) {
	t.Parallel()

	t.Run("normal_error_logging", func(t *testing.T) {
		t.Parallel()

		adapter, buff := bufferadapter.New()
		lgr := logger.New(adapter)

		lgr.Error("test message", e.New("test error"))

		entries := buff.GetAll()
		require.Len(t, entries, 1)

		assert.Equal(t, "test message", entries[0].Msg)
		assert.Contains(t, entries[0].Fields, fields.F("error", "test error"))
	})

	t.Run("typed_nil_error_logging", func(t *testing.T) {
		t.Parallel()

		adapter, buff := bufferadapter.New()
		lgr := logger.New(adapter)

		var err error = (*errWithValueReceiver)(nil)
		lgr.Error("test message", err)

		entries := buff.GetAll()
		require.Len(t, entries, 1)

		assert.Equal(t, "test message", entries[0].Msg)
		assert.Contains(t, entries[0].Fields, fields.F("error", "<nil>"))
	})

	t.Run("panicking_error_logging", func(t *testing.T) {
		t.Parallel()

		adapter, buff := bufferadapter.New()
		lgr := logger.New(adapter)

		var err error = &errThatPanics{}

		// Should not panic when logging
		assert.NotPanics(t, func() {
			lgr.Error("test message", err)
		})

		entries := buff.GetAll()
		require.Len(t, entries, 1)

		assert.Equal(t, "test message", entries[0].Msg)

		// Check that error field exists and contains panic info
		errorField := entries[0].Fields.ToDict()["error"]
		assert.NotNil(t, errorField)
		assert.Contains(t, errorField, "<PANIC=")
		assert.Contains(t, errorField, "intentional panic in Error()")
	})
}

// trackingAdapter is a test adapter that records all method calls for verification.
type trackingAdapter struct {
	shared *trackingAdapterShared
	fields []fields.Field
}

type trackingAdapterShared struct {
	logCalls        []logCall
	withFieldsCalls []withFieldsCall
	flushCalls      int
}

type logCall struct {
	level  int
	msg    string
	fields []fields.Field
}

type withFieldsCall struct {
	fields []fields.Field
}

func newTrackingAdapter() *trackingAdapter {
	return &trackingAdapter{
		shared: &trackingAdapterShared{
			logCalls:        make([]logCall, 0),
			withFieldsCalls: make([]withFieldsCall, 0),
			flushCalls:      0,
		},
		fields: make([]fields.Field, 0),
	}
}

func (ta *trackingAdapter) Log(level int, msg string, fs ...fields.Field) {
	allFields := make([]fields.Field, 0, len(ta.fields)+len(fs))

	allFields = append(allFields, ta.fields...)
	allFields = append(allFields, fs...)

	ta.shared.logCalls = append(ta.shared.logCalls, logCall{
		level:  level,
		msg:    msg,
		fields: allFields,
	})
}

func (ta *trackingAdapter) WithFields(fs ...fields.Field) logger.Adapter {
	ta.shared.withFieldsCalls = append(ta.shared.withFieldsCalls, withFieldsCall{
		fields: fs,
	})

	// Return new adapter instance with accumulated fields
	newFields := make([]fields.Field, 0, len(ta.fields)+len(fs))

	newFields = append(newFields, ta.fields...)
	newFields = append(newFields, fs...)

	newAdapter := &trackingAdapter{
		shared: ta.shared,
		fields: newFields,
	}

	return newAdapter
}

func (ta *trackingAdapter) Flush() error {
	ta.shared.flushCalls++

	return nil
}

func TestLogger(t *testing.T) {
	t.Parallel()

	assert.PanicsWithValue(t,
		"logger adapter cannot be nil, use logger.NewNop() to create no-op logger",
		func() { logger.New(nil) })

	t.Run("log_passes_correct_level_and_message", func(t *testing.T) {
		t.Parallel()

		adapter := newTrackingAdapter()
		lgr := logger.New(adapter)

		lgr.Error("error msg", nil)
		lgr.Warning("warning msg")
		lgr.Info("info msg")
		lgr.Debug("debug msg")

		require.Len(t, adapter.shared.logCalls, 3, "should have 3 log calls (default level is Info)")

		assert.Equal(t, logger.LevelError, adapter.shared.logCalls[0].level)
		assert.Equal(t, "error msg", adapter.shared.logCalls[0].msg)

		assert.Equal(t, logger.LevelWarning, adapter.shared.logCalls[1].level)
		assert.Equal(t, "warning msg", adapter.shared.logCalls[1].msg)

		assert.Equal(t, logger.LevelInfo, adapter.shared.logCalls[2].level)
		assert.Equal(t, "info msg", adapter.shared.logCalls[2].msg)
	})

	t.Run("log_with_fields_passes_fields", func(t *testing.T) {
		t.Parallel()

		adapter := newTrackingAdapter()
		lgr := logger.New(adapter)

		lgr.Info("test", fields.F("key1", "value1"), fields.F("key2", 42))

		require.Len(t, adapter.shared.logCalls, 1)
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("key1", "value1"))
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("key2", 42))
	})

	t.Run("log_with_error_includes_error_field", func(t *testing.T) {
		t.Parallel()

		adapter := newTrackingAdapter()
		lgr := logger.New(adapter)

		lgr.Error("error occurred", e.New("test error"), fields.F("context", "test"))

		require.Len(t, adapter.shared.logCalls, 1)
		assert.Equal(t, "error occurred", adapter.shared.logCalls[0].msg)
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("error", "test error"))
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("context", "test"))
	})

	t.Run("with_fields_calls_adapter_with_fields", func(t *testing.T) {
		t.Parallel()

		adapter := newTrackingAdapter()
		lgr := logger.New(adapter)

		childLgr := lgr.WithFields(fields.F("key1", "value1"), fields.F("key2", 42))
		childLgr.Info("test message")

		require.Len(t, adapter.shared.withFieldsCalls, 1, "should have 1 WithFields call")
		assert.Contains(t, adapter.shared.withFieldsCalls[0].fields, fields.F("key1", "value1"))
		assert.Contains(t, adapter.shared.withFieldsCalls[0].fields, fields.F("key2", 42))

		require.Len(t, adapter.shared.logCalls, 1)
		assert.Equal(t, "test message", adapter.shared.logCalls[0].msg)
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("key1", "value1"))
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("key2", 42))
	})

	t.Run("with_name_calls_adapter_with_name_field", func(t *testing.T) {
		t.Parallel()

		adapter := newTrackingAdapter()
		lgr := logger.New(adapter)

		childLgr := lgr.WithName("test-logger")
		childLgr.Info("test message")

		require.Len(t, adapter.shared.withFieldsCalls, 1, "should have 1 WithFields call for name")
		assert.Contains(t, adapter.shared.withFieldsCalls[0].fields, fields.F("logger-name", "test-logger"))

		require.Len(t, adapter.shared.logCalls, 1)
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("logger-name", "test-logger"))
	})

	t.Run("with_stack_trace_calls_adapter_with_stacktrace_field", func(t *testing.T) {
		t.Parallel()

		adapter := newTrackingAdapter()
		lgr := logger.New(adapter)

		childLgr := lgr.WithStackTrace(0)
		childLgr.Info("test message")

		require.Len(t, adapter.shared.withFieldsCalls, 1, "should have 1 WithFields call for stacktrace")
		require.Len(t, adapter.shared.withFieldsCalls[0].fields, 1)
		assert.Equal(t, "stacktrace", adapter.shared.withFieldsCalls[0].fields[0].K)

		require.Len(t, adapter.shared.logCalls, 1)

		// Verify stacktrace field exists in log call
		found := false

		for _, field := range adapter.shared.logCalls[0].fields {
			if field.K == "stacktrace" {
				found = true

				assert.NotEmpty(t, field.V, "stacktrace should not be empty")

				break
			}
		}

		assert.True(t, found, "stacktrace field should be present")
	})

	t.Run("chained_fields_accumulate", func(t *testing.T) {
		t.Parallel()

		adapter := newTrackingAdapter()
		lgr := logger.New(adapter)

		childLgr := lgr.WithFields(fields.F("key1", "value1")).
			WithName("test-logger").
			WithFields(fields.F("key2", "value2"))

		childLgr.Info("test message", fields.F("key3", "value3"))

		require.Len(t, adapter.shared.withFieldsCalls, 3, "should have 3 WithFields calls")

		require.Len(t, adapter.shared.logCalls, 1)
		// All fields should be accumulated
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("key1", "value1"))
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("logger-name", "test-logger"))
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("key2", "value2"))
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("key3", "value3"))
	})

	t.Run("flush_calls_adapter_flush", func(t *testing.T) {
		t.Parallel()

		adapter := newTrackingAdapter()
		lgr := logger.New(adapter)

		err := lgr.Flush()

		assert.NoError(t, err)
		assert.Equal(t, 1, adapter.shared.flushCalls, "should have called Flush once")
	})

	t.Run("parent_logger_unaffected_by_child", func(t *testing.T) {
		t.Parallel()

		adapter := newTrackingAdapter()
		lgr := logger.New(adapter)

		childLgr := lgr.WithFields(fields.F("child", "field"))

		// Clear previous calls
		adapter.shared.logCalls = nil

		// Log from parent
		lgr.Info("parent message", fields.F("parent", "field"))

		// Log from child
		childLgr.Info("child message", fields.F("another", "field"))

		require.Len(t, adapter.shared.logCalls, 2)

		// Parent log should only have parent field
		assert.Contains(t, adapter.shared.logCalls[0].fields, fields.F("parent", "field"))
		assert.NotContains(t, adapter.shared.logCalls[0].fields, fields.F("child", "field"))

		// Child log should have both child and another field
		assert.Contains(t, adapter.shared.logCalls[1].fields, fields.F("child", "field"))
		assert.Contains(t, adapter.shared.logCalls[1].fields, fields.F("another", "field"))
	})
}
