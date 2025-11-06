//nolint:err113
package zapadapter_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
	"dev.gaijin.team/go/golib/logger/zapadapter"
)

type discardingWriter struct{}

func (*discardingWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (*discardingWriter) Sync() (err error) {
	return nil
}

func zapLogger() (*zap.Logger, *observer.ObservedLogs) {
	discardingCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		&discardingWriter{},
		zapcore.DebugLevel,
	)

	testCore, logs := observer.New(zapcore.DebugLevel)

	return zap.New(zapcore.NewTee(discardingCore, testCore)), logs
}

func newAdapter(opts ...zapadapter.Option) (logger.Adapter, *observer.ObservedLogs) {
	lgr, logs := zapLogger()

	return zapadapter.New(lgr, opts...), logs
}

func TestAdapter(t *testing.T) {
	t.Parallel()

	t.Run(".Log() with standard levels", func(t *testing.T) {
		t.Parallel()

		adapter, logs := newAdapter()

		tt := []struct {
			level    int
			zapLevel zapcore.Level
			msg      string
			err      error
		}{
			{logger.LevelDebug, zapcore.DebugLevel, "debug", errors.New("debug")},
			{logger.LevelInfo, zapcore.InfoLevel, "info", errors.New("info")},
			{logger.LevelWarning, zapcore.WarnLevel, "warning", errors.New("warning")},
			{logger.LevelError, zapcore.ErrorLevel, "error", errors.New("error")},
			{logger.LevelTrace, zapcore.DebugLevel, "trace", errors.New("trace")},
		}

		for _, tc := range tt {
			logs.TakeAll()

			adapter.Log(tc.level, tc.msg, fields.F("error", tc.err))

			require.Equal(t, 1, logs.Len())
			assert.Equal(t, tc.zapLevel, logs.All()[0].Level)
			assert.Equal(t, tc.msg, logs.All()[0].Message)
			assert.Equal(t, tc.err.Error(), logs.All()[0].ContextMap()["error"])
		}
	})

	t.Run(".Log() with unknown level", func(t *testing.T) {
		t.Parallel()

		adapter, logs := newAdapter()

		adapter.Log(42, "unknown level")

		// Should map to InfoLevel by default, creating only ONE entry
		require.Equal(t, 1, logs.Len())
		assert.Equal(t, zapcore.InfoLevel, logs.All()[0].Level)
		assert.Equal(t, "unknown level", logs.All()[0].Message)
		assert.Nil(t, logs.All()[0].ContextMap()["error"])
	})

	t.Run(".WithFields()", func(t *testing.T) {
		t.Parallel()

		adapter, logs := newAdapter()

		adapter = adapter.WithFields(fields.F("foo", "bar"), fields.F("baz", 42))

		adapter.Log(logger.LevelInfo, "test")
		adapter.Log(logger.LevelDebug, "test", fields.F("bux", "qux"))
		adapter.Log(logger.LevelDebug, "test", fields.F("baz", "bar"))

		require.Equal(t, 3, logs.Len())

		assert.Equal(t, "bar", logs.All()[0].ContextMap()["foo"])
		assert.Equal(t, int64(42), logs.All()[0].ContextMap()["baz"])

		assert.Equal(t, "bar", logs.All()[1].ContextMap()["foo"])
		assert.Equal(t, int64(42), logs.All()[1].ContextMap()["baz"])
		assert.Equal(t, "qux", logs.All()[1].ContextMap()["bux"])

		assert.Equal(t, "bar", logs.All()[2].ContextMap()["foo"])
		assert.Equal(t, "bar", logs.All()[2].ContextMap()["baz"])
	})

	t.Run("WithLogLevelMapper option", func(t *testing.T) {
		t.Parallel()

		// Custom mapper that always returns ErrorLevel
		customMapper := func(level int) zapcore.Level {
			return zapcore.ErrorLevel
		}

		adapter, logs := newAdapter(zapadapter.WithLogLevelMapper(customMapper))

		// Log at Info level, but should appear as Error due to custom mapper
		adapter.Log(logger.LevelInfo, "test")

		require.Equal(t, 1, logs.Len())
		assert.Equal(t, zapcore.ErrorLevel, logs.All()[0].Level)
		assert.Equal(t, "test", logs.All()[0].Message)
	})

	t.Run("custom mapper preserved in derived adapters", func(t *testing.T) {
		t.Parallel()

		// Custom mapper that maps everything to WarnLevel
		customMapper := func(level int) zapcore.Level {
			return zapcore.WarnLevel
		}

		adapter, logs := newAdapter(zapadapter.WithLogLevelMapper(customMapper))
		derived := adapter.WithFields(fields.F("foo", "bar"))

		// Both should use custom mapper
		adapter.Log(logger.LevelInfo, "original")
		derived.Log(logger.LevelDebug, "derived")

		require.Equal(t, 2, logs.Len())
		assert.Equal(t, zapcore.WarnLevel, logs.All()[0].Level)
		assert.Equal(t, zapcore.WarnLevel, logs.All()[1].Level)
		assert.Equal(t, "bar", logs.All()[1].ContextMap()["foo"])
	})

	t.Run("NewErrorLogger()", func(t *testing.T) {
		t.Parallel()

		adapter, logs := newAdapter()
		lgr := logger.New(adapter, logger.WithLevel(logger.LevelError))
		errorLogger := logger.NewErrorLogger(lgr, logger.LevelError)

		errorLogger("test", errors.New("test"))

		require.Equal(t, 1, logs.Len())

		assert.Equal(t, zapcore.ErrorLevel, logs.All()[0].Level)
		assert.Equal(t, "test", logs.All()[0].Message)
		assert.Equal(t, "test", logs.All()[0].ContextMap()["error"])
	})
}
